package service

import (
	"context"

	"github.com/lite-cms/cms/internal/model"
	"github.com/lite-cms/cms/internal/repository"
)

// MenuService 菜单服务
type MenuService struct {
	repo *repository.MenuRepository
}

// NewMenuService 创建菜单服务
func NewMenuService(repo *repository.MenuRepository) *MenuService {
	return &MenuService{repo: repo}
}

// GetAll 获取所有菜单
func (s *MenuService) GetAll(ctx context.Context) ([]model.Menu, error) {
	return s.repo.GetAll(ctx)
}

// Create 创建菜单
func (s *MenuService) Create(ctx context.Context, name, url, icon, target, position string, sortOrder int, parentID *uint) (*model.Menu, error) {
	if target == "" {
		target = "_self"
	}
	if position == "" {
		position = "header"
	}
	menu := &model.Menu{
		Name:      name,
		URL:       url,
		Icon:      icon,
		Target:    target,
		Position:  position,
		SortOrder: sortOrder,
		ParentID:  parentID,
	}

	if err := s.repo.Create(ctx, menu); err != nil {
		return nil, err
	}

	return menu, nil
}

// Update 更新菜单
func (s *MenuService) Update(ctx context.Context, id uint, name, url, icon, target, position string, sortOrder int, parentID *uint) (*model.Menu, error) {
	menu, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if target == "" {
		target = "_self"
	}
	if position == "" {
		position = "header"
	}

	menu.Name = name
	menu.URL = url
	menu.Icon = icon
	menu.Target = target
	menu.Position = position
	menu.SortOrder = sortOrder
	menu.ParentID = parentID

	if err := s.repo.Update(ctx, menu); err != nil {
		return nil, err
	}

	return menu, nil
}

// Delete 删除菜单
func (s *MenuService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
