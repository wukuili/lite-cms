-- LiteCMS 002_add_menus 迁移 (Up)

-- 创建菜单表
CREATE TABLE IF NOT EXISTS menus (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    url VARCHAR(255) NOT NULL,
    icon VARCHAR(100),
    target VARCHAR(20) DEFAULT '_self',
    parent_id INTEGER REFERENCES menus(id) ON DELETE SET NULL,
    sort_order INTEGER DEFAULT 0,
    position VARCHAR(50) DEFAULT 'header',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引以优化查询
CREATE INDEX IF NOT EXISTS idx_menus_parent_id ON menus(parent_id);
CREATE INDEX IF NOT EXISTS idx_menus_position_sort ON menus(position, sort_order);

-- 插入一些演示数据（可选：根据需求决定是否保留）
INSERT INTO menus (name, url, position, sort_order) VALUES ('首页', '/', 'header', 1) ON CONFLICT DO NOTHING;
INSERT INTO menus (name, url, position, sort_order) VALUES ('文章', '/articles', 'header', 2) ON CONFLICT DO NOTHING;
INSERT INTO menus (name, url, position, sort_order) VALUES ('关于', '/about', 'header', 3) ON CONFLICT DO NOTHING;
