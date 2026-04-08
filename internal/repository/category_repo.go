package repository

import (
	"context"

	"github.com/lite-cms/cms/internal/model"
	"gorm.io/gorm"
)

// CategoryRepository 分类仓储
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓储
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// GetAll 获取所有分类
func (r *CategoryRepository) GetAll(ctx context.Context) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.WithContext(ctx).
		Order("sort_order ASC, id ASC").
		Find(&categories).Error
	return categories, err
}

// GetByID 根据ID获取分类
func (r *CategoryRepository) GetByID(ctx context.Context, id uint) (*model.Category, error) {
	var category model.Category
	err := r.db.WithContext(ctx).First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetBySlug 根据Slug获取分类
func (r *CategoryRepository) GetBySlug(ctx context.Context, slug string) (*model.Category, error) {
	var category model.Category
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Create 创建分类
func (r *CategoryRepository) Create(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// Update 更新分类
func (r *CategoryRepository) Update(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// ClearCategoryFromArticles 清理指定分类下文章的 category_id
func (r *CategoryRepository) ClearCategoryFromArticles(ctx context.Context, categoryID uint) error {
	return r.db.WithContext(ctx).
		Model(&model.Article{}).
		Where("category_id = ?", categoryID).
		Updates(map[string]interface{}{
			"category_id":   nil,
			"category_name": "",
			"category_slug": "",
		}).Error
}

// Delete 删除分类
func (r *CategoryRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Category{}, id).Error
}

// BatchDelete 批量删除分类
func (r *CategoryRepository) BatchDelete(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 清理相关文章的分类信息
		if err := tx.Model(&model.Article{}).
			Where("category_id IN ?", ids).
			Updates(map[string]interface{}{
				"category_id":   nil,
				"category_name": "",
				"category_slug": "",
			}).Error; err != nil {
			return err
		}
		// 删除分类
		return tx.Delete(&model.Category{}, ids).Error
	})
}

// UpdateArticleCount 更新文章计数
func (r *CategoryRepository) UpdateArticleCount(ctx context.Context, id uint) error {
	var count int64
	r.db.WithContext(ctx).
		Model(&model.Article{}).
		Where("category_id = ? AND status = ? AND deleted_at IS NULL", id, model.ArticleStatusPublished).
		Count(&count)

	return r.db.WithContext(ctx).
		Model(&model.Category{}).
		Where("id = ?", id).
		Update("article_count", count).Error
}
