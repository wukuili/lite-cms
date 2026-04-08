-- LiteCMS 初始化迁移

-- 创建分类表（必须先创建，因为文章表依赖它）
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    CONSTRAINT uni_categories_slug UNIQUE (slug),
    parent_id INTEGER REFERENCES categories(id),
    article_count INTEGER DEFAULT 0,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建标签表
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    slug VARCHAR(50) NOT NULL,
    CONSTRAINT uni_tags_slug UNIQUE (slug),
    article_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建文章表
CREATE TABLE IF NOT EXISTS articles (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    CONSTRAINT uni_articles_slug UNIQUE (slug),
    content TEXT,
    summary VARCHAR(500),
    category_id INTEGER REFERENCES categories(id),
    category_name VARCHAR(100),
    category_slug VARCHAR(100),
    author_id INTEGER NOT NULL,
    status SMALLINT DEFAULT 0,
    view_count INTEGER DEFAULT 0,
    is_top BOOLEAN DEFAULT FALSE,
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建文章标签关联表
CREATE TABLE IF NOT EXISTS article_tags (
    article_id BIGINT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (article_id, tag_id)
);

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(50),
    email VARCHAR(100),
    CONSTRAINT uni_users_username UNIQUE (username),
    CONSTRAINT uni_users_email UNIQUE (email),
    role VARCHAR(20) DEFAULT 'editor',
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建媒体表
CREATE TABLE IF NOT EXISTS media (
    id BIGSERIAL PRIMARY KEY,
    filename VARCHAR(255),
    original_name VARCHAR(255),
    mime_type VARCHAR(100),
    file_size INTEGER,
    width INTEGER,
    height INTEGER,
    storage_path VARCHAR(500),
    thumb_path VARCHAR(500),
    uploader_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引（针对百万级数据优化）
-- 文章表索引
CREATE INDEX IF NOT EXISTS idx_articles_status_published ON articles(status, published_at DESC)
    WHERE status = 1 AND deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_articles_category ON articles(category_id, published_at DESC)
    WHERE status = 1 AND deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_articles_slug ON articles(slug)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_articles_author ON articles(author_id);

-- 覆盖索引优化列表查询
CREATE INDEX IF NOT EXISTS idx_articles_list ON articles(id, title, summary, category_name, published_at, view_count)
    WHERE status = 1 AND deleted_at IS NULL;

-- 标签关联索引
CREATE INDEX IF NOT EXISTS idx_article_tags_tag ON article_tags(tag_id, article_id);

-- 全文搜索索引
CREATE INDEX IF NOT EXISTS idx_articles_content_search ON articles
    USING gin(to_tsvector('simple', title || ' ' || COALESCE(content, '')));

-- 分类排序索引
CREATE INDEX IF NOT EXISTS idx_categories_sort ON categories(sort_order);

-- 插入默认管理员（密码：admin123，实际部署需修改）
INSERT INTO users (username, password, nickname, role, status)
VALUES ('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH', 'Administrator', 'admin', 1)
ON CONFLICT (username) DO NOTHING;
