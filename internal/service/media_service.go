package service

import (
	"context"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/lite-cms/cms/internal/model"
	"github.com/lite-cms/cms/internal/pkg/storage"
	"github.com/lite-cms/cms/internal/repository"
)

// MediaService 媒体服务
type MediaService struct {
	repo    *repository.MediaRepository
	storage *storage.Storage
}

// NewMediaService 创建媒体服务
func NewMediaService(repo *repository.MediaRepository, st *storage.Storage) *MediaService {
	return &MediaService{repo: repo, storage: st}
}

// Upload 上传媒体
func (s *MediaService) Upload(ctx context.Context, file *multipart.FileHeader, uploaderID uint) (*model.Media, error) {
	path, err := s.storage.Save(file)
	if err != nil {
		return nil, err
	}

	media := &model.Media{
		Filename:     file.Filename,
		OriginalName: file.Filename,
		MimeType:     file.Header.Get("Content-Type"),
		FileSize:     int(file.Size),
		StoragePath:  path,
		UploaderID:   uploaderID,
	}

	if err := s.repo.Create(ctx, media); err != nil {
		return nil, err
	}

	return media, nil
}

// List 列表媒体
func (s *MediaService) List(ctx context.Context, page, pageSize int) ([]model.Media, int64, error) {
	return s.repo.GetAll(ctx, page, pageSize)
}

// Delete 删除媒体
func (s *MediaService) Delete(ctx context.Context, id uint) error {
	media, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// 删除磁盘上的文件 (StoragePath: "/uploads/yyyy/mm/file.ext")
	// 本地存放路径对应的是 ./static/uploads/...
	physicalPath := filepath.Join(".", "static", media.StoragePath)
	os.Remove(physicalPath)

	return nil
}
