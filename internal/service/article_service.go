package service

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/lite-cms/cms/internal/model"
	"github.com/lite-cms/cms/internal/pkg/cache"
	"github.com/lite-cms/cms/internal/pkg/counter"
	"github.com/lite-cms/cms/internal/repository"
	"github.com/mozillazg/go-pinyin"
)

// 预编译正则，避免每次 GenerateSlug 调用时重复编译
var (
	slugNonAlphaNumRe = regexp.MustCompile(`[^a-z0-9\-]+`)
	slugMultiDashRe   = regexp.MustCompile(`\-+`)
)

// ArticleService 文章服务
type ArticleService struct {
	repo        *repository.ArticleRepository
	catRepo     *repository.CategoryRepository
	tagRepo     *repository.TagRepository
	viewCounter *counter.ViewCounter
}

// NewArticleService 创建文章服务
func NewArticleService(
	repo *repository.ArticleRepository,
	catRepo *repository.CategoryRepository,
	tagRepo *repository.TagRepository,
	vc *counter.ViewCounter,
) *ArticleService {
	return &ArticleService{
		repo:        repo,
		catRepo:     catRepo,
		tagRepo:     tagRepo,
		viewCounter: vc,
	}
}

// ArticleInput 文章输入
type ArticleInput struct {
	Title      string `json:"title" binding:"required"`
	Slug       string `json:"slug"`
	Content    string `json:"content"`
	Summary    string `json:"summary"`
	CategoryID *uint  `json:"category_id"`
	TagIDs     []uint `json:"tag_ids"`
	Status     int8   `json:"status"`
	IsTop      bool   `json:"is_top"`
	LegacyURL  string `json:"legacy_url"`
}

// ListOutput 列表输出
type ListOutput struct {
	Items      []model.Article `json:"items"`
	NextCursor string          `json:"next_cursor"`
	HasMore    bool            `json:"has_more"`
}

// List 获取文章列表（带缓存）
func (s *ArticleService) List(ctx context.Context, cursor string, limit int) (*ListOutput, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("article:list:%s:%d", cursor, limit)
	if cached, ok := cache.Get(cacheKey); ok {
		if result, ok := cached.(*ListOutput); ok {
			return result, nil
		}
	}

	// 从数据库获取
	result, err := s.repo.ListPublished(ctx, cursor, limit)
	if err != nil {
		return nil, err
	}

	output := &ListOutput{
		Items:      result.Items,
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	}

	// 写入缓存
	cache.SetWithTTL(cacheKey, output, 1, time.Minute)

	return output, nil
}

// ListAll 获取所有文章列表（后台管理用，包含草稿）
func (s *ArticleService) ListAll(ctx context.Context, keyword string, cursor string, limit int) (*ListOutput, error) {
	result, err := s.repo.ListAll(ctx, keyword, cursor, limit)
	if err != nil {
		return nil, err
	}

	return &ListOutput{
		Items:      result.Items,
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}

// GetByID 根据ID获取文章
func (s *ArticleService) GetByID(ctx context.Context, id uint) (*model.Article, error) {
	return s.repo.GetByID(ctx, id)
}

// GetBySlug 根据Slug获取文章（带缓存）
func (s *ArticleService) GetBySlug(ctx context.Context, slug string) (*model.Article, error) {
	cacheKey := fmt.Sprintf("article:slug:%s", slug)
	if cached, ok := cache.Get(cacheKey); ok {
		if article, ok := cached.(*model.Article); ok {
			// 批量计数器累积浏览量
			s.viewCounter.Increment(article.ID)
			return article, nil
		}
	}

	article, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	cache.SetWithTTL(cacheKey, article, 1, 10*time.Minute)

	// 批量计数器累积浏览量
	s.viewCounter.Increment(article.ID)

	return article, nil
}

// GetByLegacyURL 根据旧路径获取文章
func (s *ArticleService) GetByLegacyURL(ctx context.Context, url string) (*model.Article, error) {
	return s.repo.GetByLegacyURL(ctx, url)
}

// Create 创建文章
func (s *ArticleService) Create(ctx context.Context, authorID uint, input *ArticleInput) (*model.Article, error) {
	article := &model.Article{
		Title:    input.Title,
		Content:  input.Content,
		Summary:  input.Summary,
		AuthorID: authorID,
		Status:   input.Status,
		IsTop:    input.IsTop,
		LegacyURL: input.LegacyURL,
	}

	// 生成Slug
	if input.Slug != "" {
		article.Slug = input.Slug
	} else {
		article.Slug = GenerateSlug(input.Title)
	}

	// 设置分类
	if input.CategoryID != nil {
		article.CategoryID = input.CategoryID
		cat, err := s.catRepo.GetByID(ctx, *input.CategoryID)
		if err == nil {
			article.CategoryName = cat.Name
			article.CategorySlug = cat.Slug
		}
	}

	// 设置发布时间
	if input.Status == model.ArticleStatusPublished {
		now := time.Now()
		article.PublishedAt = &now
	}

	// 创建文章
	if err := s.repo.Create(ctx, article); err != nil {
		return nil, err
	}

	// 同步标签
	if len(input.TagIDs) > 0 {
		s.tagRepo.SyncArticleTags(ctx, article.ID, input.TagIDs)
		for _, tagID := range input.TagIDs {
			go func(tid uint) {
				if err := s.tagRepo.UpdateArticleCount(context.Background(), tid); err != nil {
					log.Printf("异步更新标签文章计数失败 tag_id=%d: %v", tid, err)
				}
			}(tagID)
		}
	}

	// 更新分类文章计数
	if article.CategoryID != nil {
		go func(catID uint) {
			if err := s.catRepo.UpdateArticleCount(context.Background(), catID); err != nil {
				log.Printf("异步更新分类文章计数失败(create) cat_id=%d: %v", catID, err)
			}
		}(*article.CategoryID)
	}

	return article, nil
}

// Update 更新文章
func (s *ArticleService) Update(ctx context.Context, id uint, input *ArticleInput) (*model.Article, error) {
	article, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	oldCategoryID := article.CategoryID
	oldStatus := article.Status

	article.Title = input.Title
	article.Content = input.Content
	article.Summary = input.Summary
	article.CategoryID = input.CategoryID
	article.Status = input.Status
	article.IsTop = input.IsTop
	article.LegacyURL = input.LegacyURL

	if input.Slug != "" {
		article.Slug = input.Slug
	}

	// 设置分类名称
	if input.CategoryID != nil {
		cat, err := s.catRepo.GetByID(ctx, *input.CategoryID)
		if err == nil {
			article.CategoryName = cat.Name
			article.CategorySlug = cat.Slug
		}
	}

	// 状态变更时设置发布时间
	if oldStatus != model.ArticleStatusPublished && input.Status == model.ArticleStatusPublished {
		now := time.Now()
		article.PublishedAt = &now
	}

	// 更新文章
	if err := s.repo.Update(ctx, article); err != nil {
		return nil, err
	}

	// 同步标签
	if input.TagIDs != nil {
		s.tagRepo.SyncArticleTags(ctx, article.ID, input.TagIDs)
	}

	// 清除缓存
	cache.Del(fmt.Sprintf("article:slug:%s", article.Slug))

	// 更新分类文章计数
	categoryChanged := false
	if oldCategoryID == nil && article.CategoryID != nil {
		categoryChanged = true
	} else if oldCategoryID != nil && article.CategoryID == nil {
		categoryChanged = true
	} else if oldCategoryID != nil && article.CategoryID != nil && *oldCategoryID != *article.CategoryID {
		categoryChanged = true
	}

	if categoryChanged {
		if oldCategoryID != nil {
			go func(catID uint) {
				if err := s.catRepo.UpdateArticleCount(context.Background(), catID); err != nil {
					log.Printf("异步更新旧分类文章计数失败 cat_id=%d: %v", catID, err)
				}
			}(*oldCategoryID)
		}
		if article.CategoryID != nil {
			go func(catID uint) {
				if err := s.catRepo.UpdateArticleCount(context.Background(), catID); err != nil {
					log.Printf("异步更新新分类文章计数失败 cat_id=%d: %v", catID, err)
				}
			}(*article.CategoryID)
		}
	}

	return article, nil
}

// Delete 删除文章
func (s *ArticleService) Delete(ctx context.Context, id uint) error {
	article, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.SoftDelete(ctx, id); err != nil {
		return err
	}

	// 清除缓存
	cache.Del(fmt.Sprintf("article:slug:%s", article.Slug))

	// 更新分类文章计数
	if article.CategoryID != nil {
		go func(catID uint) {
			if err := s.catRepo.UpdateArticleCount(context.Background(), catID); err != nil {
				log.Printf("异步更新分类文章计数失败(delete) cat_id=%d: %v", catID, err)
			}
		}(*article.CategoryID)
	}

	return nil
}

// Search 搜索文章
func (s *ArticleService) Search(ctx context.Context, keyword, cursor string, limit int) (*ListOutput, error) {
	result, err := s.repo.Search(ctx, keyword, cursor, limit)
	if err != nil {
		return nil, err
	}

	return &ListOutput{
		Items:      result.Items,
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}

// UpdatePublishedAt 更新发布时间（内部或导入工具使用）
func (s *ArticleService) UpdatePublishedAt(ctx context.Context, id uint, t time.Time) error {
	return s.repo.UpdateColumn(ctx, id, "published_at", t)
}

// GenerateSlug 生成文章URL别名（带随机后缀防止冲突）
func GenerateSlug(title string) string {
	return generateSlugInternal(title, true)
}

// GeneratePureSlug 生成纯拼音别名（不带后缀，用于分类和标签）
func GeneratePureSlug(title string) string {
	return generateSlugInternal(title, false)
}

func generateSlugInternal(title string, withSuffix bool) string {
	a := pinyin.NewArgs()
	a.Fallback = func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	}
	py := pinyin.Pinyin(title, a)

	var sb strings.Builder
	for _, p := range py {
		if len(p) > 0 {
			sb.WriteString(p[0])
			sb.WriteString("-")
		}
	}

	slug := strings.TrimRight(sb.String(), "-")
	slug = strings.ToLower(slug)

	// 只保留字母、数字和连字符
	slug = slugNonAlphaNumRe.ReplaceAllString(slug, "-")

	// 处理连续的连字符
	slug = slugMultiDashRe.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	if slug == "" {
		slug = fmt.Sprintf("%d", time.Now().Unix())
	}

	if withSuffix {
		// 添加短随机字符串防止冲突
		return fmt.Sprintf("%s-%d", slug, time.Now().Unix()%10000)
	}

	return slug
}
