package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Article 文章模型
type Article struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"size:255;not null" json:"title"`
	Slug        string     `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	Content     string     `gorm:"type:text" json:"content"`
	Summary     string     `gorm:"size:500" json:"summary"`
	CategoryID  *uint      `json:"category_id"`
	Category    *Category  `json:"category,omitempty"`
	CategoryName string    `gorm:"size:100" json:"category_name"` // 冗余存储避免JOIN
	CategorySlug string    `gorm:"size:100" json:"category_slug"` // 冗余存储分类Slug
	AuthorID    uint       `json:"author_id"`
	Status      int8       `gorm:"default:0;index" json:"status"` // 0:草稿 1:已发布 2:回收站
	ViewCount   int        `gorm:"default:0" json:"view_count"`
	IsTop       bool       `gorm:"default:false" json:"is_top"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	LegacyURL   string     `gorm:"size:500;index" json:"legacy_url"` // 旧路径支持重定向
	Tags        []Tag      `gorm:"many2many:article_tags;" json:"tags,omitempty"`
}

// TableName 指定表名
func (Article) TableName() string {
	return "articles"
}

// Category 分类模型
type Category struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Name         string `gorm:"size:100;not null" json:"name"`
	Slug         string `gorm:"size:100;uniqueIndex;not null" json:"slug"`
	ParentID     *uint  `json:"parent_id"`
	ArticleCount int    `gorm:"default:0" json:"article_count"`
	SortOrder    int    `gorm:"default:0" json:"sort_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}

// Tag 标签模型
type Tag struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Name         string `gorm:"size:50;not null" json:"name"`
	Slug         string `gorm:"size:50;uniqueIndex;not null" json:"slug"`
	ArticleCount int    `gorm:"default:0" json:"article_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}

// ArticleTag 文章标签关联
type ArticleTag struct {
	ArticleID uint `gorm:"primaryKey"`
	TagID     uint `gorm:"primaryKey"`
}

// TableName 指定表名
func (ArticleTag) TableName() string {
	return "article_tags"
}

// User 用户模型
type User struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	Username string  `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password string  `gorm:"size:255;not null" json:"-"`
	Nickname string  `gorm:"size:50" json:"nickname"`
	Email    *string `gorm:"size:100;uniqueIndex" json:"email"`
	Role     string  `gorm:"size:20;default:'editor'" json:"role"` // admin, editor
	Status   int8    `gorm:"default:1" json:"status"` // 0:禁用 1:启用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// Media 媒体模型
type Media struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Filename     string    `gorm:"size:255" json:"filename"`
	OriginalName string    `gorm:"size:255" json:"original_name"`
	MimeType     string    `gorm:"size:100" json:"mime_type"`
	FileSize     int       `json:"file_size"`
	Width        int       `json:"width"`
	Height       int       `json:"height"`
	StoragePath  string    `gorm:"size:500" json:"storage_path"`
	ThumbPath    string    `gorm:"size:500" json:"thumb_path"`
	UploaderID   uint      `json:"uploader_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName 指定表名
func (Media) TableName() string {
	return "media"
}

// ArticleStatus 文章状态常量
const (
	ArticleStatusDraft     int8 = 0
	ArticleStatusPublished int8 = 1
	ArticleStatusTrashed   int8 = 2
)

// Menu 菜单模型
type Menu struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	URL       string    `gorm:"size:255;not null" json:"url"`
	Icon      string    `gorm:"size:100" json:"icon"`
	Target    string    `gorm:"size:20;default:'_self'" json:"target"`
	ParentID  *uint     `gorm:"index" json:"parent_id"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	Position  string    `gorm:"size:50;default:'header'" json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menus"
}

// Setting 设置模型
type Setting struct {
	Key       string    `gorm:"primaryKey;size:100;column:key" json:"key"`
	Value     string    `gorm:"type:text;column:value" json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 设置表名
func (Setting) TableName() string {
	return "settings"
}


