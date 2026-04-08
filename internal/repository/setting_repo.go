package repository

import (
	"context"

	"github.com/lite-cms/cms/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SettingRepository 设置仓储
type SettingRepository struct {
	db *gorm.DB
}

// NewSettingRepository 创建设置仓储
func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

// GetAll 获取所有设置
func (r *SettingRepository) GetAll(ctx context.Context) ([]model.Setting, error) {
	var settings []model.Setting
	err := r.db.WithContext(ctx).Find(&settings).Error
	return settings, err
}

// GetByKey 根据键获取设置
func (r *SettingRepository) GetByKey(ctx context.Context, key string) (*model.Setting, error) {
	var setting model.Setting
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// SaveBatch 批量保存设置
func (r *SettingRepository) SaveBatch(ctx context.Context, settings []model.Setting) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(&settings).Error
}
