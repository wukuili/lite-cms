package service

import (
	"context"
	"time"

	"github.com/lite-cms/cms/internal/model"
	"github.com/lite-cms/cms/internal/pkg/cache"
	"github.com/lite-cms/cms/internal/repository"
)

// CategoryService 分类服务
type CategoryService struct {
	repo *repository.CategoryRepository
}

// NewCategoryService 创建分类服务
func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// GetAll 获取所有分类（带缓存）
func (s *CategoryService) GetAll(ctx context.Context) ([]model.Category, error) {
	cacheKey := "category:list"
	if cached, ok := cache.Get(cacheKey); ok {
		if categories, ok := cached.([]model.Category); ok {
			return categories, nil
		}
	}

	categories, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	cache.SetWithTTL(cacheKey, categories, 1, 10*time.Minute)
	return categories, nil
}

// GetBySlug 根据Slug获取分类
func (s *CategoryService) GetBySlug(ctx context.Context, slug string) (*model.Category, error) {
	return s.repo.GetBySlug(ctx, slug)
}

// Create 创建分类
func (s *CategoryService) Create(ctx context.Context, name, slug string, parentID *uint) (*model.Category, error) {
	category := &model.Category{
		Name:     name,
		Slug:     slug,
		ParentID: parentID,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	// 清除缓存
	cache.Del("category:list")

	return category, nil
}

// Update 更新分类
func (s *CategoryService) Update(ctx context.Context, id uint, name, slug string, parentID *uint) (*model.Category, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	category.Name = name
	category.Slug = slug
	category.ParentID = parentID

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, err
	}

	// 清除缓存
	cache.Del("category:list")

	return category, nil
}

// Delete 删除分类
func (s *CategoryService) Delete(ctx context.Context, id uint) error {
	// 清理关联文章的 category_id（通过 Repository）
	s.repo.ClearCategoryFromArticles(ctx, id)

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	cache.Del("category:list")
	return nil
}

// BatchDelete 批量删除分类
func (s *CategoryService) BatchDelete(ctx context.Context, ids []uint) error {
	if err := s.repo.BatchDelete(ctx, ids); err != nil {
		return err
	}
	cache.Del("category:list")
	return nil
}
