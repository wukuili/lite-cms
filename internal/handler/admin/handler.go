package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lite-cms/cms/internal/middleware"
	"github.com/lite-cms/cms/internal/model"
	"github.com/lite-cms/cms/internal/pkg/database"
	"github.com/lite-cms/cms/internal/pkg/llm"
	"github.com/lite-cms/cms/internal/repository"
	"github.com/lite-cms/cms/internal/service"
	"golang.org/x/crypto/bcrypt"
)

// Handler 后台处理器
type Handler struct {
	articleSvc  *service.ArticleService
	categorySvc *service.CategoryService
	tagSvc      *service.TagService
	userRepo    *repository.UserRepository
	userSvc     *service.UserService
	mediaSvc    *service.MediaService
	menuSvc     *service.MenuService
	settingSvc  *service.SettingService
	llmClient   *llm.Client
}

// NewHandler 创建后台处理器
func NewHandler(articleSvc *service.ArticleService, categorySvc *service.CategoryService, tagSvc *service.TagService, userRepo *repository.UserRepository, userSvc *service.UserService, mediaSvc *service.MediaService, menuSvc *service.MenuService, settingSvc *service.SettingService, llmCli *llm.Client) *Handler {
	return &Handler{
		articleSvc:  articleSvc,
		categorySvc: categorySvc,
		tagSvc:      tagSvc,
		userRepo:    userRepo,
		userSvc:     userSvc,
		mediaSvc:    mediaSvc,
		menuSvc:     menuSvc,
		settingSvc:  settingSvc,
		llmClient:   llmCli,
	}
}

// RegisterRoutes 注册后台路由
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/logout", middleware.AuthRequired(), h.Logout)
		auth.GET("/profile", middleware.AuthRequired(), h.Profile)
	}

	dashboard := r.Group("/dashboard")
	dashboard.Use(middleware.AuthRequired())
	{
		dashboard.GET("/stats", h.DashboardStats)
	}

	articles := r.Group("/articles")
	articles.Use(middleware.AuthRequired())
	{
		articles.GET("", h.ListArticles)
		articles.POST("", h.CreateArticle)
		articles.GET("/:id", h.GetArticle)
		articles.PUT("/:id", h.UpdateArticle)
		articles.DELETE("/:id", h.DeleteArticle)
		articles.POST("/ai-suggest", h.AISuggest)
		articles.POST("/generate-slug", h.GenerateSlugAPI)
	}

	categories := r.Group("/categories")
	categories.Use(middleware.AuthRequired())
	{
		categories.GET("", h.ListCategories)
		categories.POST("", h.CreateCategory)
		categories.PUT("/:id", h.UpdateCategory)
		categories.DELETE("/:id", h.DeleteCategory)
		categories.POST("/batch-delete", h.BatchDeleteCategories)
	}

	tags := r.Group("/tags")
	tags.Use(middleware.AuthRequired())
	{
		tags.GET("", h.ListTags)
		tags.POST("", h.CreateTag)
		tags.PUT("/:id", h.UpdateTag)
		tags.DELETE("/:id", h.DeleteTag)
		tags.POST("/batch-delete", h.BatchDeleteTags)
	}

	media := r.Group("/media")
	media.Use(middleware.AuthRequired())
	{
		media.GET("", h.ListMedia)
		media.POST("/upload", h.UploadMedia)
		media.DELETE("/:id", h.DeleteMedia)
	}

	menus := r.Group("/menus")
	menus.Use(middleware.AuthRequired(), middleware.AdminRequired())
	{
		menus.GET("", h.ListMenus)
		menus.POST("", h.CreateMenu)
		menus.PUT("/:id", h.UpdateMenu)
		menus.DELETE("/:id", h.DeleteMenu)
	}

	users := r.Group("/users")
	users.Use(middleware.AuthRequired(), middleware.AdminRequired())
	{
		users.GET("", h.ListUsers)
		users.POST("", h.CreateUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
	}

	settings := r.Group("/settings")
	settings.Use(middleware.AuthRequired(), middleware.AdminRequired())
	{
		settings.GET("", h.ListSettings)
		settings.POST("", h.SaveSettings)
	}
}

// Login 登录
func (h *Handler) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	user, err := h.userRepo.GetByUsername(c.Request.Context(), input.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	if user.Status != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "账号已被禁用"})
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	// 设置Cookie，用于浏览器静默认证（如页面导航）
	isSecure := c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https"
	c.SetCookie("token", token, 3600*24, "/admin", "", isSecure, true)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// Logout 退出登录
func (h *Handler) Logout(c *gin.Context) {
	// 清除Cookie
	c.SetCookie("token", "", -1, "/admin", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "退出成功"})
}

// Profile 获取当前用户信息
func (h *Handler) Profile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"id":       userID,
		"username": username,
		"role":     role,
	})
}

// ListArticles 文章列表
func (h *Handler) ListArticles(c *gin.Context) {
	cursor := c.Query("cursor")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.articleSvc.ListAll(c.Request.Context(), cursor, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetArticle 获取文章
func (h *Handler) GetArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	// 通过 repository 直接获取（admin 需要获取所有状态的文章）
	article, err := h.articleSvc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}
	c.JSON(http.StatusOK, article)
}

// CreateArticle 创建文章
func (h *Handler) CreateArticle(c *gin.Context) {
	var input service.ArticleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	
	uid, ok := userID.(uint)
	if !ok {
		if fuid, ok := userID.(float64); ok {
			uid = uint(fuid)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的用户ID"})
			return
		}
	}

	article, err := h.articleSvc.Create(c.Request.Context(), uid, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, article)
}

// UpdateArticle 更新文章
func (h *Handler) UpdateArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var input service.ArticleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	article, err := h.articleSvc.Update(c.Request.Context(), uint(id), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, article)
}

// DeleteArticle 删除文章
func (h *Handler) DeleteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := h.articleSvc.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ListCategories 分类列表
func (h *Handler) ListCategories(c *gin.Context) {
	categories, err := h.categorySvc.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// CreateCategory 创建分类
func (h *Handler) CreateCategory(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Slug     string `json:"slug" binding:"required"`
		ParentID *uint  `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	category, err := h.categorySvc.Create(c.Request.Context(), input.Name, input.Slug, input.ParentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// UpdateCategory 更新分类
func (h *Handler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var input struct {
		Name     string `json:"name" binding:"required"`
		Slug     string `json:"slug" binding:"required"`
		ParentID *uint  `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	category, err := h.categorySvc.Update(c.Request.Context(), uint(id), input.Name, input.Slug, input.ParentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory 删除分类
func (h *Handler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := h.categorySvc.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// BatchDeleteCategories 批量删除分类
func (h *Handler) BatchDeleteCategories(c *gin.Context) {
	var input struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.categorySvc.BatchDelete(c.Request.Context(), input.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ListTags 标签列表
func (h *Handler) ListTags(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	tags, total, err := h.tagSvc.List(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": tags,
		"total": total,
	})
}

// CreateTag 创建标签
func (h *Handler) CreateTag(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
		Slug string `json:"slug" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	tag, err := h.tagSvc.Create(c.Request.Context(), input.Name, input.Slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// UpdateTag 更新标签
func (h *Handler) UpdateTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var input struct {
		Name string `json:"name" binding:"required"`
		Slug string `json:"slug" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	tag, err := h.tagSvc.Update(c.Request.Context(), uint(id), input.Name, input.Slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// DeleteTag 删除标签
func (h *Handler) DeleteTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := h.tagSvc.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// BatchDeleteTags 批量删除标签
func (h *Handler) BatchDeleteTags(c *gin.Context) {
	var input struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.tagSvc.BatchDelete(c.Request.Context(), input.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// === Helper functions ===

// ListMedia 媒体列表
func (h *Handler) ListMedia(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	
	medias, total, err := h.mediaSvc.List(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": medias,
		"total": total,
	})
}

// UploadMedia 上传媒体
func (h *Handler) UploadMedia(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件"})
		return
	}

	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	var uid uint
	switch v := userID.(type) {
	case uint:
		uid = v
	case float64:
		uid = uint(v)
	}

	media, err := h.mediaSvc.Upload(c.Request.Context(), file, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, media)
}

// DeleteMedia 删除媒体
func (h *Handler) DeleteMedia(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := h.mediaSvc.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// DashboardStats 获取仪表盘统计数据
func (h *Handler) DashboardStats(c *gin.Context) {
	db := database.GetDB()
	var articleCount int64
	var viewsCount int64
	var categoryCount int64
	var tagCount int64

	db.Model(&model.Article{}).Count(&articleCount)
	row := db.Model(&model.Article{}).Select("COALESCE(SUM(view_count), 0)").Row()
	row.Scan(&viewsCount)
	db.Model(&model.Category{}).Count(&categoryCount)
	db.Model(&model.Tag{}).Count(&tagCount)

	c.JSON(http.StatusOK, gin.H{
		"articles":   articleCount,
		"views":      viewsCount,
		"categories": categoryCount,
		"tags":       tagCount,
	})
}

// AISuggest AI智能建议
func (h *Handler) AISuggest(c *gin.Context) {
	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "正文内容不能为空"})
		return
	}

	if h.llmClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "系统未配置 AI 助手"})
		return
	}

	result, err := h.llmClient.GenerateSummaryAndTags(input.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI 生成失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GenerateSlugAPI 生成文章URL别名
func (h *Handler) GenerateSlugAPI(c *gin.Context) {
	var input struct {
		Title string `json:"title" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题不能为空"})
		return
	}

	slug := service.GenerateSlug(input.Title)
	c.JSON(http.StatusOK, gin.H{"slug": slug})
}
