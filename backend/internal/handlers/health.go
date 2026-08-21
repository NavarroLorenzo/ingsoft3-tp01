package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Health(c *gin.Context) {
	if h.DB == nil {
		writeError(c, http.StatusServiceUnavailable, "El servicio no está saludable.")
		return
	}
	sqlDB, err := h.DB.DB()
	if err != nil || sqlDB.PingContext(c.Request.Context()) != nil {
		writeError(c, http.StatusServiceUnavailable, "El servicio no está saludable.")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}
