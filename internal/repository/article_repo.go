package repository

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/lite-cms/cms/internal/model"
	"gorm.io/gorm"
)

// ArticleRepository 文章仓储
type ArticleRepository struct {
	db *gorm.DB
}

// NewArticleRepository 创建文章仓储
func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// CursorResult 游标分页结果
type CursorResult struct {
	Items      []model.Article `json:"items"`
	NextCursor string          `json:"next_cursor"`
	HasMore    bool            `json:"has_more"`
}

// CursorDecode 解码游标
type CursorDecode struct {
	ID          uint      `json:"id"`
	PublishedAt time.Time `json:"published_at"`
}

// ListPublished 获取已发布文章列表（游标分页）
func (r *ArticleRepository) ListPublished(ctx context.Context, cursor string, limit int) (*CursorResult, error) {
	var articles []model.Article
	query := r.db.WithContext(ctx).
		Preload("Tags").
		Where("status = ? AND deleted_at IS NULL", model.ArticleStatusPublished).
		Order("id DESC").
		Limit(limit + 1) // 多取一条判断是否有下一页

	if cursor != "" {
		decoded, err := r.decodeCursor(cursor)
		if err == nil {
			query = query.Where("id < ?", decoded.ID)
		}
	}

	if err := query.Find(&articles).Error; err != nil {
		return nil, err
	}

	hasMore := len(articles) > limit
	if hasMore {
		articles = articles[:limit]
	}

	nextCursor := ""
	if hasMore && len(articles) > 0 {
		nextCursor = r.encodeCursor(articles[len(articles)-1].ID, time.Now())
	}

	return &CursorResult{
		Items:      articles,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

// ListAll 获取所有文章列表（后台管理用，包含草稿）
func (r *ArticleRepository) ListAll(ctx context.Context, cursor string, limit int) (*CursorResult, error) {
	var articles []model.Article
	query := r.db.WithContext(ctx).
		Preload("Tags").
		Where("deleted_at IS NULL").
		Order("id DESC").
		Limit(limit + 1)

	if cursor != "" {
		decoded, err := r.decodeCursor(cursor)
		if err == nil {
			query = query.Where("id < ?", decoded.ID)
		}
	}

	if err := query.Find(&articles).Error; err != nil {
		return nil, err
	}

	hasMore := len(articles) > limit
	if hasMore {
		articles = articles[:limit]
	}

	nextCursor := ""
	if hasMore && len(articles) > 0 {
		nextCursor = r.encodeCursor(articles[len(articles)-1].ID, time.Now())
	}

	return &CursorResult{
		Items:      articles,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

// ListByCategory 按分类获取文章列表
func (r *ArticleRepository) ListByCategory(ctx context.Context, categoryID uint, cursor string, limit int) (*CursorResult, error) {
	var articles []model.Article
	query := r.db.WithContext(ctx).
		Where("status = ? AND category_id = ? AND deleted_at IS NULL", model.ArticleStatusPublished, categoryID).
		Order("id DESC").
		Limit(limit + 1)

	if cursor != "" {
		decoded, err := r.decodeCursor(cursor)
		if err == nil {
			query = query.Where("id < ?", decoded.ID)
		}
	}

	if err := query.Find(&articles).Error; err != nil {
		return nil, err
	}

	hasMore := len(articles) > limit
	if hasMore {
		articles = articles[:limit]
	}

	nextCursor := ""
	if hasMore && len(articles) > 0 {
		nextCursor = r.encodeCursor(articles[len(articles)-1].ID, time.Now())
	}

	return &CursorResult{
		Items:      articles,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

// ListByTag 按标签获取文章列表
func (r *ArticleRepository) ListByTag(ctx context.Context, tagID uint, cursor string, limit int) (*CursorResult, error) {
	var articles []model.Article
	query := r.db.WithContext(ctx).
		Joins("JOIN article_tags ON article_tags.article_id = articles.id").
		Where("articles.status = ? AND article_tags.tag_id = ? AND articles.deleted_at IS NULL", model.ArticleStatusPublished, tagID).
		Order("articles.id DESC").
		Limit(limit + 1)

	if cursor != "" {
		decoded, err := r.decodeCursor(cursor)
		if err == nil {
			query = query.Where("articles.id < ?", decoded.ID)
		}
	}

	if err := query.Find(&articles).Error; err != nil {
		return nil, err
	}

	hasMore := len(articles) > limit
	if hasMore {
		articles = articles[:limit]
	}

	nextCursor := ""
	if hasMore && len(articles) > 0 {
		nextCursor = r.encodeCursor(articles[len(articles)-1].ID, time.Now())
	}

	return &CursorResult{
		Items:      articles,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

// GetByID 根据ID获取文章
func (r *ArticleRepository) GetByID(ctx context.Context, id uint) (*model.Article, error) {
	var article model.Article
	err := r.db.WithContext(ctx).
		Preload("Tags").
		Preload("Category").
		First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// GetBySlug 根据Slug获取文章
func (r *ArticleRepository) GetBySlug(ctx context.Context, slug string) (*model.Article, error) {
	var article model.Article
	err := r.db.WithContext(ctx).
		Preload("Tags").
		Preload("Category").
		Where("slug = ? AND deleted_at IS NULL", slug).
		First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// GetByLegacyURL 根据旧路径获取文章
func (r *ArticleRepository) GetByLegacyURL(ctx context.Context, url string) (*model.Article, error) {
	var article model.Article
	err := r.db.WithContext(ctx).
		Where("legacy_url = ? AND deleted_at IS NULL", url).
		First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// Create 创建文章
func (r *ArticleRepository) Create(ctx context.Context, article *model.Article) error {
	return r.db.WithContext(ctx).Create(article).Error
}

// Update 更新文章
func (r *ArticleRepository) Update(ctx context.Context, article *model.Article) error {
	return r.db.WithContext(ctx).Save(article).Error
}

// SoftDelete 软删除文章
func (r *ArticleRepository) SoftDelete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&model.Article{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

// Search 全文搜索
func (r *ArticleRepository) Search(ctx context.Context, keyword string, cursor string, limit int) (*CursorResult, error) {
	var articles []model.Article

	// 采用 ILIKE 进行模糊搜索，因为 PostgreSQL 的 simple 解析器对中文分词支持不佳，会导致无法匹配子串
	searchKey := "%" + keyword + "%"
	query := r.db.WithContext(ctx).
		Where("status = ? AND deleted_at IS NULL", model.ArticleStatusPublished).
		Where("title ILIKE ? OR content ILIKE ?", searchKey, searchKey).
		Order("id DESC").
		Limit(limit + 1)

	if cursor != "" {
		decoded, err := r.decodeCursor(cursor)
		if err == nil {
			query = query.Where("id < ?", decoded.ID)
		}
	}

	if err := query.Find(&articles).Error; err != nil {
		return nil, err
	}

	hasMore := len(articles) > limit
	if hasMore {
		articles = articles[:limit]
	}

	nextCursor := ""
	if hasMore && len(articles) > 0 {
		nextCursor = r.encodeCursor(articles[len(articles)-1].ID, time.Now())
	}

	return &CursorResult{
		Items:      articles,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

// UpdateColumn 更新单个列
func (r *ArticleRepository) UpdateColumn(ctx context.Context, id uint, column string, value interface{}) error {
	return r.db.WithContext(ctx).
		Model(&model.Article{}).
		Where("id = ?", id).
		UpdateColumn(column, value).Error
}

// IncrementViewCount 增加浏览量
func (r *ArticleRepository) IncrementViewCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&model.Article{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// encodeCursor 编码游标
func (r *ArticleRepository) encodeCursor(id uint, publishedAt time.Time) string {
	data := CursorDecode{ID: id, PublishedAt: publishedAt}
	jsonData, _ := json.Marshal(data)
	return base64.URLEncoding.EncodeToString(jsonData)
}

// decodeCursor 解码游标
func (r *ArticleRepository) decodeCursor(cursor string) (*CursorDecode, error) {
	jsonData, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}
	var data CursorDecode
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
