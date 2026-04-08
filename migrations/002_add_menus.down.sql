-- LiteCMS 002_add_menus 迁移 (Down)

-- 删除索引
DROP INDEX IF EXISTS idx_menus_position_sort;
DROP INDEX IF EXISTS idx_menus_parent_id;

-- 删除菜单表
DROP TABLE IF EXISTS menus;
