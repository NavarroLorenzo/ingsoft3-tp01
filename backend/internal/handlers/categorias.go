package handlers

import (
	"errors"
	"net/http"

	"gestor-gastos/backend/internal/models"
	"gestor-gastos/backend/internal/validation"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *Handler) GetCategorias(c *gin.Context) {
	var categorias []models.Categoria
	if err := h.DB.Order("nombre ASC").Find(&categorias).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "No se pudieron obtener las categorías.")
		return
	}
	c.JSON(http.StatusOK, categorias)
}

func (h *Handler) GetCategoria(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var categoria models.Categoria
	if err := h.DB.First(&categoria, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "Categoría no encontrada.")
			return
		}
		writeError(c, http.StatusInternalServerError, "No se pudo obtener la categoría.")
		return
	}
	c.JSON(http.StatusOK, categoria)
}

func (h *Handler) CreateCategoria(c *gin.Context) {
	var input validation.CategoriaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, "El cuerpo de la solicitud es inválido.")
		return
	}
	if err := validation.ValidateCategoria(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	if h.categoryNameExists(input.Nombre, 0) {
		writeError(c, http.StatusConflict, "Ya existe una categoría con ese nombre.")
		return
	}

	categoria := models.Categoria{Nombre: input.Nombre}
	if err := h.DB.Create(&categoria).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "No se pudo crear la categoría.")
		return
	}
	c.JSON(http.StatusCreated, categoria)
}

func (h *Handler) UpdateCategoria(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var input validation.CategoriaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, "El cuerpo de la solicitud es inválido.")
		return
	}
	if err := validation.ValidateCategoria(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var categoria models.Categoria
	if err := h.DB.First(&categoria, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "Categoría no encontrada.")
			return
		}
		writeError(c, http.StatusInternalServerError, "No se pudo actualizar la categoría.")
		return
	}
	if h.categoryNameExists(input.Nombre, id) {
		writeError(c, http.StatusConflict, "Ya existe una categoría con ese nombre.")
		return
	}

	categoria.Nombre = input.Nombre
	if err := h.DB.Save(&categoria).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "No se pudo actualizar la categoría.")
		return
	}
	c.JSON(http.StatusOK, categoria)
}

func (h *Handler) DeleteCategoria(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var categoria models.Categoria
	if err := h.DB.First(&categoria, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "Categoría no encontrada.")
			return
		}
		writeError(c, http.StatusInternalServerError, "No se pudo eliminar la categoría.")
		return
	}
	var count int64
	if err := h.DB.Model(&models.Gasto{}).Where("categoria_id = ?", id).Count(&count).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "No se pudo eliminar la categoría.")
		return
	}
	if count > 0 {
		writeError(c, http.StatusConflict, "No se puede eliminar la categoría porque tiene gastos asociados.")
		return
	}
	if err := h.DB.Delete(&categoria).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "No se pudo eliminar la categoría.")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) categoryNameExists(nombre string, excludeID uint) bool {
	var count int64
	query := h.DB.Model(&models.Categoria{}).Where("LOWER(nombre) = LOWER(?)", nombre)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	return query.Count(&count).Error == nil && count > 0
}
