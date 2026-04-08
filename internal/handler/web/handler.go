package web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lite-cms/cms/internal/pkg/cache"
	"github.com/lite-cms/cms/internal/repository"
	"github.com/lite-cms/cms/internal/service"
)

// Handler 前台处理器
type Handler struct {
	articleSvc  *service.ArticleService
	categorySvc *service.CategoryService
	tagSvc      *service.TagService
	articleRepo *repository.ArticleRepository
	menuSvc     *service.MenuService
	settingSvc  *service.SettingService
}

// NewHandler 创建前台处理器
func NewHandler(
	articleSvc *service.ArticleService,
	categorySvc *service.CategoryService,
	tagSvc *service.TagService,
	articleRepo *repository.ArticleRepository,
	menuSvc *service.MenuService,
	settingSvc *service.SettingService,
) *Handler {
	return &Handler{
		articleSvc:  articleSvc,
		categorySvc: categorySvc,
		tagSvc:      tagSvc,
		articleRepo: articleRepo,
		menuSvc:     menuSvc,
		settingSvc:  settingSvc,
	}
}

// RegisterRoutes 注册前台路由
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/", h.Index)
	r.GET("/posts", h.Articles)
	r.GET("/posts/:slug", h.ArticleDetail)
	r.GET("/category/:slug", h.CategoryArticles)
	r.GET("/tags/:slug", h.TagArticles)
	r.GET("/search", h.Search)
	r.GET("/sitemap.xml", h.Sitemap)
	r.GET("/feed.xml", h.Feed)

	// 捕获所有未匹配路由，尝试处理旧路径重定向
	r.NoRoute(h.RedirectLegacyURL)
}

// Index 首页
func (h *Handler) Index(c *gin.Context) {
	result, err := h.articleSvc.List(c.Request.Context(), "", 10)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"articles": result.Items,
		"has_more": result.HasMore,
		"next":     result.NextCursor,
		"menus":    h.getMenus(c),
		"settings": h.getSettings(c),
	})
}

// Articles 文章列表
func (h *Handler) Articles(c *gin.Context) {
	cursor := c.Query("cursor")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.articleSvc.List(c.Request.Context(), cursor, limit)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "articles.html", gin.H{
		"articles":   result.Items,
		"has_more":   result.HasMore,
		"next":       result.NextCursor,
		"categories": h.getCategories(c),
		"menus":      h.getMenus(c),
		"settings":   h.getSettings(c),
	})
}

// ArticleDetail 文章详情
func (h *Handler) ArticleDetail(c *gin.Context) {
	slug := c.Param("slug")

	article, err := h.articleSvc.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "文章不存在"})
		return
	}

	c.HTML(http.StatusOK, "article.html", gin.H{
		"article":  article,
		"menus":    h.getMenus(c),
		"settings": h.getSettings(c),
	})
}

// RedirectLegacyURL 处理旧路径重定向
func (h *Handler) RedirectLegacyURL(c *gin.Context) {
	path := c.Request.URL.Path
	if path == "" || path == "/" {
		return
	}

	// 尝试从数据库查找匹配的旧路径
	article, err := h.articleSvc.GetByLegacyURL(c.Request.Context(), path)
	if err == nil && article != nil {
		// 找到匹配项，执行 301 重定向到新路径
		newURL := fmt.Sprintf("/posts/%s", article.Slug)
		c.Redirect(http.StatusMovedPermanently, newURL)
		return
	}

	// 没找到，返回 404
	c.HTML(http.StatusNotFound, "error.html", gin.H{
		"error":    "您访问的页面不存在",
		"menus":    h.getMenus(c),
		"settings": h.getSettings(c),
	})
}

// CategoryArticles 分类文章
func (h *Handler) CategoryArticles(c *gin.Context) {
	slug := c.Param("slug")
	cursor := c.Query("cursor")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	category, err := h.categorySvc.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "分类不存在"})
		return
	}

	result, err := h.articleRepo.ListByCategory(c.Request.Context(), category.ID, cursor, limit)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "category.html", gin.H{
		"category": category,
		"articles": result.Items,
		"has_more": result.HasMore,
		"next":     result.NextCursor,
		"menus":    h.getMenus(c),
		"settings": h.getSettings(c),
	})
}

// TagArticles 标签文章
func (h *Handler) TagArticles(c *gin.Context) {
	slug := c.Param("slug")
	cursor := c.Query("cursor")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	tag, err := h.tagSvc.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "标签不存在"})
		return
	}

	// 获取该标签的文章列表
	result, err := h.articleRepo.ListByTag(c.Request.Context(), tag.ID, cursor, limit)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "tag.html", gin.H{
		"tag":      tag,
		"articles": result.Items,
		"has_more": result.HasMore,
		"next":     result.NextCursor,
		"menus":    h.getMenus(c),
		"settings": h.getSettings(c),
	})
}

// Search 搜索（增加关键词长度限制）
func (h *Handler) Search(c *gin.Context) {
	keyword := c.Query("q")
	cursor := c.Query("cursor")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if keyword == "" {
		c.HTML(http.StatusOK, "search.html", gin.H{
			"menus":    h.getMenus(c),
			"settings": h.getSettings(c),
		})
		return
	}

	// 搜索关键词长度限制：至少2个字符，最多50个字符
	keywordRunes := []rune(keyword)
	if len(keywordRunes) < 2 {
		c.HTML(http.StatusOK, "search.html", gin.H{
			"keyword":  keyword,
			"error":    "搜索关键词至少需要2个字符",
			"menus":    h.getMenus(c),
			"settings": h.getSettings(c),
		})
		return
	}
	if len(keywordRunes) > 50 {
		keyword = string(keywordRunes[:50])
	}

	result, err := h.articleSvc.Search(c.Request.Context(), keyword, cursor, limit)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "search.html", gin.H{
		"keyword":  keyword,
		"articles": result.Items,
		"has_more": result.HasMore,
		"next":     result.NextCursor,
		"menus":    h.getMenus(c),
		"settings": h.getSettings(c),
	})
}

// Sitemap 站点地图（带缓存，避免频繁查库）
func (h *Handler) Sitemap(c *gin.Context) {
	type URL struct {
		Loc        string `xml:"loc"`
		LastMod    string `xml:"lastmod"`
		ChangeFreq string `xml:"changefreq"`
		Priority   string `xml:"priority"`
	}

	type SitemapData struct {
		XMLName string `xml:"urlset"`
		Xmlns   string `xml:"xmlns,attr"`
		Urls    []URL  `xml:"url"`
	}

	// 尝试从缓存获取已生成的sitemap
	cacheKey := fmt.Sprintf("web:sitemap:%s", c.Request.Host)
	if cached, ok := cache.Get(cacheKey); ok {
		if sitemap, ok := cached.(*SitemapData); ok {
			c.XML(http.StatusOK, sitemap)
			return
		}
	}

	result, _ := h.articleSvc.List(c.Request.Context(), "", 500)

	var urls []URL
	urls = append(urls, URL{
		Loc:        "http://" + c.Request.Host + "/",
		LastMod:    time.Now().Format("2006-01-02"),
		ChangeFreq: "daily",
		Priority:   "1.0",
	})

	for _, art := range result.Items {
		lastMod := art.UpdatedAt.Format("2006-01-02")
		if art.PublishedAt != nil {
			lastMod = art.PublishedAt.Format("2006-01-02")
		}
		urls = append(urls, URL{
			Loc:        "http://" + c.Request.Host + "/posts/" + art.Slug,
			LastMod:    lastMod,
			ChangeFreq: "weekly",
			Priority:   "0.8",
		})
	}

	sitemap := &SitemapData{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Urls:  urls,
	}
	// 缓存1小时
	cache.SetWithTTL(cacheKey, sitemap, 1, time.Hour)
	c.XML(http.StatusOK, sitemap)
}

// Feed RSS订阅
func (h *Handler) Feed(c *gin.Context) {
	type Item struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		PubDate     string `xml:"pubDate"`
		Guid        string `xml:"guid"`
	}

	type Channel struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Items       []Item `xml:"item"`
	}

	type Rss struct {
		XMLName string  `xml:"rss"`
		Version string  `xml:"version,attr"`
		Channel Channel `xml:"channel"`
	}

	result, _ := h.articleSvc.List(c.Request.Context(), "", 20)
    
	var items []Item
	for _, art := range result.Items {
		pubDate := art.CreatedAt.Format(time.RFC1123Z)
		if art.PublishedAt != nil {
			pubDate = art.PublishedAt.Format(time.RFC1123Z)
		}
		items = append(items, Item{
			Title:       art.Title,
			Link:        "http://" + c.Request.Host + "/posts/" + art.Slug,
			Description: art.Summary,
			PubDate:     pubDate,
			Guid:        "http://" + c.Request.Host + "/posts/" + art.Slug,
		})
	}
	
	rss := Rss{
		Version: "2.0",
		Channel: Channel{
			Title:       "LiteCMS",
			Link:        "http://" + c.Request.Host + "/",
			Description: "基于 LiteCMS 构建的轻量级站点",
			Items:       items,
		},
	}
	c.XML(http.StatusOK, rss)
}

// getCategories 获取分类列表
func (h *Handler) getCategories(c *gin.Context) interface{} {
	categories, _ := h.categorySvc.GetAll(c.Request.Context())
	return categories
}

// getMenus 获取菜单列表并且区分顶部和底部（带缓存）
func (h *Handler) getMenus(c *gin.Context) map[string]interface{} {
	cacheKey := "web:menus"
	if cached, ok := cache.Get(cacheKey); ok {
		if result, ok := cached.(map[string]interface{}); ok {
			return result
		}
	}

	menus, _ := h.menuSvc.GetAll(c.Request.Context())
	
	var headerMenus []interface{}
	var footerMenus []interface{}
	
	for _, m := range menus {
		if m.Position == "footer" {
			footerMenus = append(footerMenus, m)
		} else {
			headerMenus = append(headerMenus, m)
		}
	}
	
	result := map[string]interface{}{
		"header": headerMenus,
		"footer": footerMenus,
	}
	cache.SetWithTTL(cacheKey, result, 1, 5*time.Minute)
	return result
}

// getSettings 获取系统设置（带缓存）
func (h *Handler) getSettings(c *gin.Context) map[string]string {
	cacheKey := "web:settings"
	if cached, ok := cache.Get(cacheKey); ok {
		if result, ok := cached.(map[string]string); ok {
			log.Printf("DEBUG: 缓存命中, site_name=%s", result["site_name"])
			return result
		}
	}

	settings, err := h.settingSvc.GetMap(c.Request.Context())
	if err != nil {
		log.Printf("DEBUG: 从数据库加载设置失败: %v", err)
	} else {
		log.Printf("DEBUG: 从数据库加载设置成功, site_name=%s, 键数量=%d", settings["site_name"], len(settings))
		for k, v := range settings {
			log.Printf("  - %s: %s", k, v)
		}
	}
	
	cache.SetWithTTL(cacheKey, settings, 1, 5*time.Minute)
	return settings
}
