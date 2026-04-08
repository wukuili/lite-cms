package repository

import (
	"context"

	"github.com/lite-cms/cms/internal/model"
	"gorm.io/gorm"
)

// TagRepository 标签仓储
type TagRepository struct {
	db *gorm.DB
}

// NewTagRepository 创建标签仓储
func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

// GetAll 获取所有标签
func (r *TagRepository) GetAll(ctx context.Context) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.WithContext(ctx).
		Order("article_count DESC, id ASC").
		Find(&tags).Error
	return tags, err
}

// List 获取分页标签
func (r *TagRepository) List(ctx context.Context, page, pageSize int) ([]model.Tag, int64, error) {
	var tags []model.Tag
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Tag{})
	db.Count(&total)

	offset := (page - 1) * pageSize
	err := db.Order("article_count DESC, id ASC").
		Offset(offset).Limit(pageSize).
		Find(&tags).Error

	return tags, total, err
}

// GetByID 根据ID获取标签
func (r *TagRepository) GetByID(ctx context.Context, id uint) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.WithContext(ctx).First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetBySlug 根据Slug获取标签
func (r *TagRepository) GetBySlug(ctx context.Context, slug string) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// Create 创建标签
func (r *TagRepository) Create(ctx context.Context, tag *model.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

// Update 更新标签
func (r *TagRepository) Update(ctx context.Context, tag *model.Tag) error {
	return r.db.WithContext(ctx).Save(tag).Error
}

// Delete 删除标签
func (r *TagRepository) Delete(ctx context.Context, id uint) error {
	// 先清理关联的 article_tags 记录
	if err := r.db.WithContext(ctx).
		Where("tag_id = ?", id).
		Delete(&model.ArticleTag{}).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Delete(&model.Tag{}, id).Error
}

// BatchDelete 批量删除标签
func (r *TagRepository) BatchDelete(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 清理关联的文章标签记录
		if err := tx.Where("tag_id IN ?", ids).Delete(&model.ArticleTag{}).Error; err != nil {
			return err
		}
		// 删除标签
		return tx.Delete(&model.Tag{}, ids).Error
	})
}

// GetByArticleID 获取文章的标签
func (r *TagRepository) GetByArticleID(ctx context.Context, articleID uint) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.WithContext(ctx).
		Joins("JOIN article_tags ON article_tags.tag_id = tags.id").
		Where("article_tags.article_id = ?", articleID).
		Find(&tags).Error
	return tags, err
}

// SyncArticleTags 同步文章标签
func (r *TagRepository) SyncArticleTags(ctx context.Context, articleID uint, tagIDs []uint) error {
	// 删除现有关联
	if err := r.db.WithContext(ctx).
		Where("article_id = ?", articleID).
		Delete(&model.ArticleTag{}).Error; err != nil {
		return err
	}

	// 创建新关联
	if len(tagIDs) == 0 {
		return nil
	}

	var articleTags []model.ArticleTag
	for _, tagID := range tagIDs {
		articleTags = append(articleTags, model.ArticleTag{
			ArticleID: articleID,
			TagID:     tagID,
		})
	}

	return r.db.WithContext(ctx).Create(&articleTags).Error
}

// UpdateArticleCount 更新文章计数
func (r *TagRepository) UpdateArticleCount(ctx context.Context, id uint) error {
	var count int64
	r.db.WithContext(ctx).
		Model(&model.ArticleTag{}).
		Where("tag_id = ?", id).
		Count(&count)

	return r.db.WithContext(ctx).
		Model(&model.Tag{}).
		Where("id = ?", id).
		Update("article_count", count).Error
}
