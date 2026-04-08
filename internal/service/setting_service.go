package service

import (
	"context"

	"github.com/lite-cms/cms/internal/model"
	"github.com/lite-cms/cms/internal/repository"
)

// SettingService 设置服务
type SettingService struct {
	repo *repository.SettingRepository
}

// NewSettingService 创建设置服务
func NewSettingService(repo *repository.SettingRepository) *SettingService {
	return &SettingService{repo: repo}
}

// GetAll 获取所有设置返回为数组
func (s *SettingService) GetAll(ctx context.Context) ([]model.Setting, error) {
	return s.repo.GetAll(ctx)
}

// GetMap 获取所有设置返回为Map
func (s *SettingService) GetMap(ctx context.Context) (map[string]string, error) {
	settings, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	res := make(map[string]string)
	for _, setting := range settings {
		res[setting.Key] = setting.Value
	}
	return res, nil
}

// GetValue 获取单个配置值
func (s *SettingService) GetValue(ctx context.Context, key string, defaultVal string) string {
	setting, err := s.repo.GetByKey(ctx, key)
	if err != nil {
		return defaultVal
	}
	return setting.Value
}

// SaveBatch 批量保存设置
func (s *SettingService) SaveBatch(ctx context.Context, data map[string]string) error {
	var settings []model.Setting
	for k, v := range data {
		settings = append(settings, model.Setting{
			Key:   k,
			Value: v,
		})
	}
	return s.repo.SaveBatch(ctx, settings)
}
