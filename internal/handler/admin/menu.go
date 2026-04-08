package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lite-cms/cms/internal/pkg/cache"
)

// ListMenus 菜单列表
func (h *Handler) ListMenus(c *gin.Context) {
	menus, err := h.menuSvc.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menus)
}

// CreateMenu 创建菜单
func (h *Handler) CreateMenu(c *gin.Context) {
	var input struct {
		Name      string `json:"name" binding:"required"`
		URL       string `json:"url" binding:"required"`
		Icon      string `json:"icon"`
		Target    string `json:"target"`
		Position  string `json:"position"`
		SortOrder int    `json:"sort_order"`
		ParentID  *uint  `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	menu, err := h.menuSvc.Create(c.Request.Context(), input.Name, input.URL, input.Icon, input.Target, input.Position, input.SortOrder, input.ParentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, menu)
	cache.Del("web:menus")
}

// UpdateMenu 更新菜单
func (h *Handler) UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var input struct {
		Name      string `json:"name" binding:"required"`
		URL       string `json:"url" binding:"required"`
		Icon      string `json:"icon"`
		Target    string `json:"target"`
		Position  string `json:"position"`
		SortOrder int    `json:"sort_order"`
		ParentID  *uint  `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	menu, err := h.menuSvc.Update(c.Request.Context(), uint(id), input.Name, input.URL, input.Icon, input.Target, input.Position, input.SortOrder, input.ParentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
	cache.Del("web:menus")
}

// DeleteMenu 删除菜单
func (h *Handler) DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := h.menuSvc.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cache.Del("web:menus")
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
