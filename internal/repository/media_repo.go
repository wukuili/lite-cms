package repository

import (
	"context"

	"github.com/lite-cms/cms/internal/model"
	"gorm.io/gorm"
)

// MediaRepository 媒体仓储
type MediaRepository struct {
	db *gorm.DB
}

// NewMediaRepository 创建媒体仓储
func NewMediaRepository(db *gorm.DB) *MediaRepository {
	return &MediaRepository{db: db}
}

// Create 创建媒体记录
func (r *MediaRepository) Create(ctx context.Context, media *model.Media) error {
	return r.db.WithContext(ctx).Create(media).Error
}

// GetAll 获取所有媒体
func (r *MediaRepository) GetAll(ctx context.Context, page, pageSize int) ([]model.Media, int64, error) {
	var medias []model.Media
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Media{})
	query.Count(&total)

	err := query.Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&medias).Error

	return medias, total, err
}

// Delete 删除媒体
func (r *MediaRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Media{}, id).Error
}

// GetByID 获取单个媒体
func (r *MediaRepository) GetByID(ctx context.Context, id uint) (*model.Media, error) {
	var media model.Media
	err := r.db.WithContext(ctx).First(&media, id).Error
	if err != nil {
		return nil, err
	}
	return &media, nil
}
