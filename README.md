# LiteCMS

轻量级CMS系统，专为低配置VPS优化，支持百万级文章。

## 特性

- **低内存占用**：应用内存控制在256MB以内
- **百万级支持**：游标分页、优化的索引策略
- **高性能缓存**：内置ristretto缓存
- **服务端渲染**：SEO友好，无需前端框架

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑.env文件，配置数据库连接
```

### 3. 初始化数据库

```bash
psql -U postgres -d litecms -f migrations/001_init.up.sql
```

### 4. 运行

```bash
go run cmd/server/main.go
```

## 项目结构

```
lite-cms/
├── cmd/server/main.go      # 入口
├── internal/
│   ├── config/             # 配置
│   ├── handler/            # 处理器
│   │   ├── admin/          # 后台API
│   │   └── web/            # 前台API
│   ├── middleware/         # 中间件
│   ├── model/              # 数据模型
│   ├── repository/         # 数据访问
│   ├── service/            # 业务逻辑
│   └── pkg/
│       ├── cache/          # 缓存
│       ├── database/       # 数据库
│       └── storage/        # 存储
├── migrations/             # 数据库迁移
├── static/                 # 静态资源
└── templates/              # 模板文件
```

## API

### 前台
- `GET /` - 首页
- `GET /articles` - 文章列表
- `GET /article/:slug` - 文章详情
- `GET /category/:slug` - 分类文章
- `GET /tag/:slug` - 标签文章
- `GET /search?q=keyword` - 搜索

### 后台
- `POST /api/admin/auth/login` - 登录
- `GET /api/admin/articles` - 文章列表
- `POST /api/admin/articles` - 创建文章
- `PUT /api/admin/articles/:id` - 更新文章
- `DELETE /api/admin/articles/:id` - 删除文章

## 性能优化

1. **游标分页**：替代OFFSET，O(1)复杂度
2. **覆盖索引**：避免回表查询
3. **部分索引**：只索引需要的记录
4. **冗余存储**：category_name避免JOIN
	5. **内置缓存**：ristretto 10MB
6. **批量刷新**：浏览量在内存中累积后异步批量更新，缩减DB压力

## 内存和资源控制 (低配VPS优化目标 256MB)

- 数据库连接池：5 open / 2 idle （降低连接内存消耗）
- 缓存：10MB / Max 10K keys
- 编译文件体积优化：通过去除符号表可以减少约30%的二进制大小

### 安装及编译构建

对生产环境的极简部署，请在构建时带上 `ldflags` 参数以压缩文件体积：

```bash
# 压缩体积且不包含调试信息编译
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o server ./cmd/server

# 可选使用 UPX 进一步二次压缩 
# upx --best --lzma server
```

## License

MIT
