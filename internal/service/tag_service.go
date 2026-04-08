package service

import (
	"context"
	"time"

	"github.com/lite-cms/cms/internal/model"
	"github.com/lite-cms/cms/internal/pkg/cache"
	"github.com/lite-cms/cms/internal/repository"
)

// TagService 标签服务
type TagService struct {
	repo *repository.TagRepository
}

// NewTagService 创建标签服务
func NewTagService(repo *repository.TagRepository) *TagService {
	return &TagService{repo: repo}
}

// GetAll 获取所有标签（带缓存）
func (s *TagService) GetAll(ctx context.Context) ([]model.Tag, error) {
	cacheKey := "tag:list"
	if cached, ok := cache.Get(cacheKey); ok {
		if tags, ok := cached.([]model.Tag); ok {
			return tags, nil
		}
	}

	tags, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	cache.SetWithTTL(cacheKey, tags, 1, 10*time.Minute)
	return tags, nil
}

// List 获取分页标签
func (s *TagService) List(ctx context.Context, page, pageSize int) ([]model.Tag, int64, error) {
	return s.repo.List(ctx, page, pageSize)
}

// GetBySlug 根据Slug获取标签
func (s *TagService) GetBySlug(ctx context.Context, slug string) (*model.Tag, error) {
	return s.repo.GetBySlug(ctx, slug)
}

// Create 创建标签
func (s *TagService) Create(ctx context.Context, name, slug string) (*model.Tag, error) {
	tag := &model.Tag{
		Name: name,
		Slug: slug,
	}

	if err := s.repo.Create(ctx, tag); err != nil {
		return nil, err
	}

	cache.Del("tag:list")
	return tag, nil
}

// Update 更新标签
func (s *TagService) Update(ctx context.Context, id uint, name, slug string) (*model.Tag, error) {
	tag, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	tag.Name = name
	tag.Slug = slug

	if err := s.repo.Update(ctx, tag); err != nil {
		return nil, err
	}

	cache.Del("tag:list")
	return tag, nil
}

// Delete 删除标签
func (s *TagService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	cache.Del("tag:list")
	return nil
}

// BatchDelete 批量删除标签
func (s *TagService) BatchDelete(ctx context.Context, ids []uint) error {
	if err := s.repo.BatchDelete(ctx, ids); err != nil {
		return err
	}
	cache.Del("tag:list")
	return nil
}
