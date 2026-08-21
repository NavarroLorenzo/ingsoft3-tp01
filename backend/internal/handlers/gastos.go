package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"gestor-gastos/backend/internal/models"
	"gestor-gastos/backend/internal/validation"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *Handler) GetGastos(c *gin.Context) {
	query := h.DB.Preload("Categoria").Where("usuario_id = ?", currentUserID(c)).Model(&models.Gasto{})

	if value := c.Query("categoriaId"); value != "" {
		categoriaID, err := strconv.ParseUint(value, 10, 64)
		if err != nil || categoriaID == 0 {
			writeError(c, http.StatusBadRequest, "El filtro categoriaId es inválido.")
			return
		}
		query = query.Where("categoria_id = ?", categoriaID)
	}
	if value := c.Query("desde"); value != "" {
		fecha, err := validation.ParseDate(value)
		if err != nil {
			writeError(c, http.StatusBadRequest, "El filtro desde es inválido.")
			return
		}
		query = query.Where("fecha >= ?", fecha)
	}
	if value := c.Query("hasta"); value != "" {
		fecha, err := validation.ParseDate(value)
		if err != nil {
			writeError(c, http.StatusBadRequest, "El filtro hasta es inválido.")
			return
		}
		query = query.Where("fecha < ?", fecha.AddDate(0, 0, 1))
	}
	if value := strings.TrimSpace(c.Query("texto")); value != "" {
		query = query.Where("descripcion ILIKE ?", "%"+value+"%")
	}

	var gastos []models.Gasto
	if err := query.Order("fecha DESC, id DESC").Find(&gastos).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "No se pudieron obtener los gastos.")
		return
	}
	c.JSON(http.StatusOK, gastos)
}

func (h *Handler) GetGasto(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var gasto models.Gasto
	if err := h.DB.Preload("Categoria").Where("id = ? AND usuario_id = ?", id, currentUserID(c)).First(&gasto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "Gasto no encontrado.")
			return
		}
		writeError(c, http.StatusInternalServerError, "No se pudo obtener el gasto.")
		return
	}
	c.JSON(http.StatusOK, gasto)
}

func (h *Handler) CreateGasto(c *gin.Context) {
	var input validation.GastoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, "El cuerpo de la solicitud es inválido.")
		return
	}
	fecha, err := validation.ValidateGasto(&input)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	if !h.categoryExists(input.CategoriaID) {
		writeError(c, http.StatusBadRequest, "La categoría indicada no existe.")
		return
	}

	gasto := models.Gasto{Descripcion: input.Descripcion, Monto: input.Monto, Fecha: fecha, CategoriaID: input.CategoriaID, UsuarioID: currentUserID(c)}
	if err := h.DB.Create(&gasto).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "No se pudo crear el gasto.")
		return
	}
	if err := h.DB.Preload("Categoria").First(&gasto, gasto.ID).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "El gasto fue creado pero no se pudo recuperar.")
		return
	}
	c.JSON(http.StatusCreated, gasto)
}

func (h *Handler) UpdateGasto(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var input validation.GastoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, "El cuerpo de la solicitud es inválido.")
		return
	}
	fecha, err := validation.ValidateGasto(&input)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var gasto models.Gasto
	if err := h.DB.Where("id = ? AND usuario_id = ?", id, currentUserID(c)).First(&gasto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "Gasto no encontrado.")
			return
		}
		writeError(c, http.StatusInternalServerError, "No se pudo actualizar el gasto.")
		return
	}
	if !h.categoryExists(input.CategoriaID) {
		writeError(c, http.StatusBadRequest, "La categoría indicada no existe.")
		return
	}

	gasto.Descripcion = input.Descripcion
	gasto.Monto = input.Monto
	gasto.Fecha = fecha
	gasto.CategoriaID = input.CategoriaID
	if err := h.DB.Save(&gasto).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "No se pudo actualizar el gasto.")
		return
	}
	if err := h.DB.Preload("Categoria").First(&gasto, gasto.ID).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "El gasto fue actualizado pero no se pudo recuperar.")
		return
	}
	c.JSON(http.StatusOK, gasto)
}

func (h *Handler) DeleteGasto(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	result := h.DB.Where("id = ? AND usuario_id = ?", id, currentUserID(c)).Delete(&models.Gasto{})
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "No se pudo eliminar el gasto.")
		return
	}
	if result.RowsAffected == 0 {
		writeError(c, http.StatusNotFound, "Gasto no encontrado.")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) categoryExists(id uint) bool {
	var count int64
	return h.DB.Model(&models.Categoria{}).Where("id = ?", id).Count(&count).Error == nil && count > 0
}

// dateRange adds the same optional date filters used by the summary endpoint.
func dateRange(c *gin.Context, query *gorm.DB) (*gorm.DB, bool) {
	if value := c.Query("desde"); value != "" {
		fecha, err := validation.ParseDate(value)
		if err != nil {
			writeError(c, http.StatusBadRequest, "El filtro desde es inválido.")
			return nil, false
		}
		query = query.Where("fecha >= ?", fecha)
	}
	if value := c.Query("hasta"); value != "" {
		fecha, err := validation.ParseDate(value)
		if err != nil {
			writeError(c, http.StatusBadRequest, "El filtro hasta es inválido.")
			return nil, false
		}
		query = query.Where("fecha < ?", fecha.AddDate(0, 0, 1))
	}
	return query, true
}
