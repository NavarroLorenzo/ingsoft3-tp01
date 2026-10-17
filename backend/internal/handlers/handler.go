package handlers

import (
	"net/http"
	"strconv"

	"gestor-gastos/backend/internal/middleware"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB        *gorm.DB
	JWTSecret string
}

func NewRouter(handler *Handler) *gin.Engine {
	router := gin.Default()
	// Mide cada ruta y captura panics; Gin conserva la respuesta 500 y sigue atendiendo.
	router.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
	router.GET("/health", handler.Health)
	// El monitor puede usar HEAD; conserva el mismo ping a la base que hace GET.
	router.HEAD("/health", handler.Health)

	api := router.Group("/api")
	{
		authRoutes := api.Group("/auth")
		authRoutes.POST("/register", handler.Register)
		authRoutes.POST("/login", handler.Login)
		authRoutes.GET("/me", middleware.AuthRequired(handler.JWTSecret), handler.Me)

		protected := api.Group("/")
		protected.Use(middleware.AuthRequired(handler.JWTSecret))
		protected.GET("/categorias", handler.GetCategorias)
		protected.GET("/categorias/:id", handler.GetCategoria)
		protected.POST("/categorias", handler.CreateCategoria)
		protected.PUT("/categorias/:id", handler.UpdateCategoria)
		protected.DELETE("/categorias/:id", handler.DeleteCategoria)
		protected.GET("/gastos", handler.GetGastos)
		protected.GET("/gastos/:id", handler.GetGasto)
		protected.POST("/gastos", handler.CreateGasto)
		protected.PUT("/gastos/:id", handler.UpdateGasto)
		protected.DELETE("/gastos/:id", handler.DeleteGasto)
		protected.GET("/resumen", handler.GetResumen)
	}
	return router
}

func currentUserID(c *gin.Context) uint {
	userID, _ := c.Get("userID")
	return userID.(uint)
}

func writeError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func parseID(c *gin.Context) (uint, bool) {
	// IntSize coincide con uint: evita truncar IDs si se compila para 32 bits.
	id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
	if err != nil || id == 0 {
		writeError(c, http.StatusBadRequest, "El identificador es inválido.")
		return 0, false
	}
	return uint(id), true
}
