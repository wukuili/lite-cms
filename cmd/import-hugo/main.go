package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/pelletier/go-toml/v2"
	"github.com/lite-cms/cms/internal/config"
	"github.com/lite-cms/cms/internal/model"
	"github.com/lite-cms/cms/internal/pkg/database"
	"github.com/lite-cms/cms/internal/repository"
	"github.com/lite-cms/cms/internal/service"
)

type HugoFrontMatter struct {
	Title      string   `yaml:"title" toml:"title"`
	Date       any      `yaml:"date" toml:"date"` // Can be time.Time or string
	Tags       []string `yaml:"tags" toml:"tags"`
	Categories []string `yaml:"categories" toml:"categories"`
	Draft      bool     `yaml:"draft" toml:"draft"`
	Summary    string   `yaml:"summary" toml:"summary"`
	Slug       string   `yaml:"slug" toml:"slug"`
	URL        string   `yaml:"url" toml:"url"`
}

func main() {
	importPath := flag.String("path", "", "Path to Hugo content directory")
	authorID := flag.Uint("author", 0, "Default author ID (0 to use first admin)")
	flag.Parse()

	if *importPath == "" {
		log.Fatal("Please specify Hugo content path using -path flag")
	}

	// 1. Load config and DB
	cfg := config.Load()
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 2. Setup Repos and Services
	articleRepo := repository.NewArticleRepository(db)
	catRepo := repository.NewCategoryRepository(db)
	tagRepo := repository.NewTagRepository(db)
	articleService := service.NewArticleService(articleRepo, catRepo, tagRepo, nil) // No view counter needed for import

	// 3. Get Default Author
	if *authorID == 0 {
		var firstUser model.User
		if err := db.Where("role = ?", "admin").First(&firstUser).Error; err != nil {
			log.Fatalf("No admin user found to assign articles. Please create one first.")
		}
		*authorID = firstUser.ID
		log.Printf("Assigning articles to default admin: %s (ID: %d)", firstUser.Username, firstUser.ID)
	}

	// 4. Walk directory
	err = filepath.WalkDir(*importPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		log.Printf("Processing %s...", path)
		if err := processFile(path, *authorID, articleService, catRepo, tagRepo); err != nil {
			log.Printf("Error processing %s: %v", path, err)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Walk failed: %v", err)
	}

	fmt.Println("\nHugo import completed successfully!")
}

func processFile(path string, authorID uint, svc *service.ArticleService, catRepo *repository.CategoryRepository, tagRepo *repository.TagRepository) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	fm, body, format, err := parseHugoFile(content)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	// Clean basic metadata
	fm.Title = cleanString(fm.Title)
	fm.Slug = cleanString(fm.Slug)
	fm.Summary = cleanString(fm.Summary)

	// Map Categories
	var categoryID *uint
	if len(fm.Categories) > 0 {
		catName := truncateString(cleanString(fm.Categories[0]), 50)
		catSlug := service.GeneratePureSlug(catName)
		if len(catSlug) > 50 {
			catSlug = catSlug[:50]
		}
		cat, err := catRepo.GetBySlug(context.Background(), catSlug)
		if err != nil {
			// Create category
			cat = &model.Category{Name: catName, Slug: catSlug}
			if err := catRepo.Create(context.Background(), cat); err != nil {
				log.Printf("Warning: failed to create category %s: %v", catName, err)
			} else {
				categoryID = &cat.ID
			}
		} else {
			categoryID = &cat.ID
		}
	}

	// Map Tags
	var tagIDs []uint
	for _, tagName := range fm.Tags {
		tagName = truncateString(cleanString(tagName), 50)
		tagSlug := service.GeneratePureSlug(tagName)
		if len(tagSlug) > 50 {
			tagSlug = tagSlug[:50]
		}
		tag, err := tagRepo.GetBySlug(context.Background(), tagSlug)
		if err != nil {
			tag = &model.Tag{Name: tagName, Slug: tagSlug}
			if err := tagRepo.Create(context.Background(), tag); err != nil {
				log.Printf("Warning: failed to create tag %s: %v", tagName, err)
			} else {
				tagIDs = append(tagIDs, tag.ID)
			}
		} else {
			tagIDs = append(tagIDs, tag.ID)
		}
	}

	// Prepare Article Input
	status := model.ArticleStatusPublished
	if fm.Draft {
		status = model.ArticleStatusDraft
	}

	summary := fm.Summary
	if summary == "" {
		// Auto-generate summary from body
		runes := []rune(body)
		if len(runes) > 200 {
			summary = string(runes[:200]) + "..."
		} else {
			summary = string(runes)
		}
	}
	if len(summary) > 500 {
		summary = truncateString(summary, 500)
	}

	input := &service.ArticleInput{
		Title:      truncateString(fm.Title, 255),
		Slug:       truncateString(fm.Slug, 255),
		Content:    body,
		Summary:    summary,
		CategoryID: categoryID,
		TagIDs:     tagIDs,
		Status:     status,
		LegacyURL:  fm.URL,
	}

	// Create Article
	article, err := svc.Create(context.Background(), authorID, input)
	if err != nil {
		return fmt.Errorf("failed to create article: %w", err)
	}

	// Fix PublishedAt if provided
	if fm.Date != nil {
		var pubTime time.Time
		switch t := fm.Date.(type) {
		case time.Time:
			pubTime = t
		case string:
			formats := []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02"}
			for _, f := range formats {
				if parsed, err := time.Parse(f, t); err == nil {
					pubTime = parsed
					break
				}
			}
		}
		if !pubTime.IsZero() {
			// Update the database record directly for PublishedAt
			if err := svc.UpdatePublishedAt(context.Background(), article.ID, pubTime); err != nil {
				log.Printf("Warning: failed to update PublishedAt for %s: %v", fm.Title, err)
			}
		}
	}

	log.Printf("Imported: %s (Format: %s)", fm.Title, format)
	return nil
}

func cleanString(s string) string {
	s = strings.ReplaceAll(s, "\r", "")
	return strings.TrimSpace(s)
}

func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) > maxLen {
		if maxLen > 10 {
			return string(runes[:maxLen-3]) + "..."
		}
		return string(runes[:maxLen])
	}
	return s
}

func parseHugoFile(data []byte) (*HugoFrontMatter, string, string, error) {
	// Handle UTF-8 BOM
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))
	
	str := string(data)
	str = strings.TrimSpace(str)
	var fmStr string
	var body string
	var format string

	if strings.HasPrefix(str, "---") {
		parts := strings.SplitN(str, "---", 3)
		if len(parts) < 3 {
			return nil, "", "", fmt.Errorf("invalid YAML front matter")
		}
		fmStr = parts[1]
		body = strings.TrimSpace(parts[2])
		format = "YAML"
	} else if strings.HasPrefix(str, "+++") {
		parts := strings.SplitN(str, "+++", 3)
		if len(parts) < 3 {
			return nil, "", "", fmt.Errorf("invalid TOML front matter")
		}
		fmStr = parts[1]
		body = strings.TrimSpace(parts[2])
		format = "TOML"
	} else {
		return nil, "", "", fmt.Errorf("no front matter found (starts with: %q)", getPeek(str))
	}

	var fm HugoFrontMatter
	if format == "YAML" {
		if err := yaml.Unmarshal([]byte(fmStr), &fm); err != nil {
			return nil, "", "", err
		}
	} else {
		if err := toml.Unmarshal([]byte(fmStr), &fm); err != nil {
			return nil, "", "", err
		}
	}

	return &fm, body, format, nil
}

func getPeek(s string) string {
	if len(s) > 20 {
		return s[:20]
	}
	return s
}
