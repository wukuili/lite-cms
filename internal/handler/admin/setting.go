package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lite-cms/cms/internal/pkg/cache"
)

// ListSettings 获取所有设置
func (h *Handler) ListSettings(c *gin.Context) {
	settings, err := h.settingSvc.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, settings)
}

// SaveSettings 批量保存设置
func (h *Handler) SaveSettings(c *gin.Context) {
	var input map[string]string
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.settingSvc.SaveBatch(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 清除前台设置缓存，使改动立即生效
	cache.Del("web:settings")

	c.JSON(http.StatusOK, gin.H{"message": "设置已保存"})
}
