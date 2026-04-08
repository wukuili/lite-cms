package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Storage 存储服务
type Storage struct {
	UploadPath string
	MaxSize    int64
}

// New 创建存储服务
func New(uploadPath string, maxSize int64) *Storage {
	// 确保目录存在
	os.MkdirAll(uploadPath, 0755)
	return &Storage{
		UploadPath: uploadPath,
		MaxSize:    maxSize,
	}
}

// Save 上传文件
func (s *Storage) Save(file *multipart.FileHeader) (string, error) {
	if file.Size > s.MaxSize {
		return "", fmt.Errorf("文件大小超过限制，最大允许 %d 字节", s.MaxSize)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
		return "", fmt.Errorf("不支持的文件格式: %s", ext)
	}

	// 生成新文件名：年/月/时间戳_系统原名.ext
	now := time.Now()
	datePath := fmt.Sprintf("%04d/%02d", now.Year(), now.Month())
	fullDir := filepath.Join(s.UploadPath, datePath)
	
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return "", err
	}

	newFilename := fmt.Sprintf("%d%s", now.UnixNano(), ext)
	relativePath := filepath.Join(datePath, newFilename)
	dstPath := filepath.Join(s.UploadPath, relativePath)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return "", err
	}

	// 返回可以通过URL访问的路径
	return "/uploads/" + strings.ReplaceAll(relativePath, "\\", "/"), nil
}
