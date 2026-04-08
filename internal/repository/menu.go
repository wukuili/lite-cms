package repository

import (
	"context"

	"github.com/lite-cms/cms/internal/model"
	"gorm.io/gorm"
)

// MenuRepository 菜单仓储
type MenuRepository struct {
	db *gorm.DB
}

// NewMenuRepository 创建菜单仓储
func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

// GetAll 获取所有菜单（默认按 Position 和 SortOrder 排序）
func (r *MenuRepository) GetAll(ctx context.Context) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.WithContext(ctx).
		Order("position ASC, sort_order ASC, id ASC").
		Find(&menus).Error
	return menus, err
}

// GetByID 根据ID获取菜单
func (r *MenuRepository) GetByID(ctx context.Context, id uint) (*model.Menu, error) {
	var menu model.Menu
	err := r.db.WithContext(ctx).First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// Create 创建菜单
func (r *MenuRepository) Create(ctx context.Context, menu *model.Menu) error {
	return r.db.WithContext(ctx).Create(menu).Error
}

// Update 更新菜单
func (r *MenuRepository) Update(ctx context.Context, menu *model.Menu) error {
	return r.db.WithContext(ctx).Save(menu).Error
}

// Delete 删除菜单
func (r *MenuRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Menu{}, id).Error
}
