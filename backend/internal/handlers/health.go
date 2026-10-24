package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BuildCommit identifica el código del binario; Docker lo fija al compilar.
// El valor local permite seguir usando Compose sin pasar argumentos de build.
var BuildCommit = "desconocido"

func (h *Handler) Health(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	if h.DB == nil {
		writeError(c, http.StatusServiceUnavailable, "El servicio no está saludable.")
		return
	}
	sqlDB, err := h.DB.DB()
	if err != nil || sqlDB.PingContext(c.Request.Context()) != nil {
		writeError(c, http.StatusServiceUnavailable, "El servicio no está saludable.")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "commit": BuildCommit})
}
