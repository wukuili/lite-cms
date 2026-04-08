package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/lite-cms/cms/internal/config"
	adminHandler "github.com/lite-cms/cms/internal/handler/admin"
	webHandler "github.com/lite-cms/cms/internal/handler/web"
	"github.com/lite-cms/cms/internal/middleware"
	"github.com/lite-cms/cms/internal/model"
	"github.com/lite-cms/cms/internal/pkg/cache"
	"github.com/lite-cms/cms/internal/pkg/counter"
	"github.com/lite-cms/cms/internal/pkg/database"
	"github.com/lite-cms/cms/internal/pkg/llm"
	"github.com/lite-cms/cms/internal/pkg/storage"
	"github.com/lite-cms/cms/internal/repository"
	"github.com/lite-cms/cms/internal/service"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 设置运行模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	log.Println("数据库连接成功")

	// 自动迁移（开发环境）
	// 始终自动迁移以确保表结构正确
	if err := autoMigrate(db); err != nil {
		log.Printf("自动迁移失败: %v", err)
	}

	// 初始化缓存
	_, err = cache.Init(&cfg.Cache)
	if err != nil {
		log.Printf("缓存初始化失败: %v", err)
	} else {
		log.Println("缓存初始化成功")
	}

	// 创建Repository
	articleRepo := repository.NewArticleRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	tagRepo := repository.NewTagRepository(db)
	userRepo := repository.NewUserRepository(db)
	mediaRepo := repository.NewMediaRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	settingRepo := repository.NewSettingRepository(db)

	st := storage.New(cfg.Storage.UploadPath, cfg.Storage.MaxSize)

	// 初始化浏览量批量计数器
	viewCounter := counter.New(db)
	defer viewCounter.Stop()

	// 创建Service
	articleSvc := service.NewArticleService(articleRepo, categoryRepo, tagRepo, viewCounter)
	categorySvc := service.NewCategoryService(categoryRepo)
	tagSvc := service.NewTagService(tagRepo)
	mediaSvc := service.NewMediaService(mediaRepo, st)
	menuSvc := service.NewMenuService(menuRepo)
	userSvc := service.NewUserService(userRepo)
	settingSvc := service.NewSettingService(settingRepo)
	llmClient := llm.NewClient(cfg.LLM.BaseURL, cfg.LLM.APIKey, cfg.LLM.Model)

	// 创建Handler
	adminHdl := adminHandler.NewHandler(articleSvc, categorySvc, tagSvc, userRepo, userSvc, mediaSvc, menuSvc, settingSvc, llmClient)
	webHdl := webHandler.NewHandler(articleSvc, categorySvc, tagSvc, articleRepo, menuSvc, settingSvc)

	// 创建路由
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.RateLimit(cfg.Server.RateLimit, cfg.Server.RateBurst))

	// 静态文件
	r.Static("/static", "./static")
	r.Static("/uploads", "./static/uploads")

	// 后台管理页面
	r.GET("/admin/login", func(c *gin.Context) {
		c.File("./static/admin/login.html")
	})

	admin := r.Group("/admin", middleware.AuthRequired())
	{
		admin.GET("/", func(c *gin.Context) {
			c.File("./static/admin/index.html")
		})
		admin.GET("/articles", func(c *gin.Context) {
			c.File("./static/admin/articles.html")
		})
		admin.GET("/categories", func(c *gin.Context) {
			c.File("./static/admin/categories.html")
		})
		admin.GET("/tags", func(c *gin.Context) {
			c.File("./static/admin/tags.html")
		})
		admin.GET("/media", func(c *gin.Context) {
			c.File("./static/admin/media.html")
		})

		// 仅限管理员页面
		adminOnly := admin.Group("/", middleware.AdminRequired())
		{
			adminOnly.GET("/menus", func(c *gin.Context) {
				c.File("./static/admin/menus.html")
			})
			adminOnly.GET("/users", func(c *gin.Context) {
				c.File("./static/admin/users.html")
			})
			adminOnly.GET("/settings", func(c *gin.Context) {
				c.File("./static/admin/settings.html")
			})
		}
	}

	// 模板
	r.HTMLRender = createRenderer()

	// 前台路由
	webHdl.RegisterRoutes(r)

	// 后台API
	api := r.Group("/api/admin")
	adminHdl.RegisterRoutes(api)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("服务器启动在 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// autoMigrate 自动迁移数据库表结构
func autoMigrate(db *gorm.DB) error {
	log.Println("正在开始数据库自动迁移...")
	// 使用GORM自动迁移
	err := db.AutoMigrate(
		&model.Category{},
		&model.Tag{},
		&model.Article{},
		&model.ArticleTag{},
		&model.User{},
		&model.Media{},
		&model.Menu{},
		&model.Setting{},
	)
	if err != nil {
		log.Printf("自动迁移记录错误: %v", err)
	} else {
		log.Println("自动迁移成功完成")
	}
	return err
}

// createRenderer 创建模板渲染器
func createRenderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	// 复用 Goldmark 实例，增加 WithUnsafe 选项以支持导入的 HTML 内容
	md := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	funcMap := template.FuncMap{
		"safe": func(x string) template.HTML { return template.HTML(x) },
		"default": func(val, defaultVal string) string {
			if val == "" {
				return defaultVal
			}
			return val
		},
		"renderMarkdown": func(content string) template.HTML {
			var buf bytes.Buffer
			if err := md.Convert([]byte(content), &buf); err != nil {
				return template.HTML(template.HTMLEscapeString(content))
			}
			return template.HTML(buf.String())
		},
	}

	// 前台页面模板
	r.AddFromFilesFuncs("index.html", funcMap, "templates/base.html", "templates/index.html")
	r.AddFromFilesFuncs("articles.html", funcMap, "templates/base.html", "templates/articles.html")
	r.AddFromFilesFuncs("article.html", funcMap, "templates/base.html", "templates/article.html")
	r.AddFromFilesFuncs("category.html", funcMap, "templates/base.html", "templates/category.html")
	r.AddFromFilesFuncs("tag.html", funcMap, "templates/base.html", "templates/tag.html")
	r.AddFromFilesFuncs("search.html", funcMap, "templates/base.html", "templates/search.html")
	r.AddFromFilesFuncs("error.html", funcMap, "templates/base.html", "templates/error.html")

	return r
}
