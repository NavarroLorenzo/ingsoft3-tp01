package handlers

import (
	"net/http"

	"gestor-gastos/backend/internal/models"
	"github.com/gin-gonic/gin"
)

type resumenCategoria struct {
	CategoriaID uint    `json:"categoriaId"`
	Categoria   string  `json:"categoria"`
	Total       float64 `json:"total"`
}

type resumenResponse struct {
	Total          float64            `json:"total"`
	CantidadGastos int64              `json:"cantidadGastos"`
	PorCategoria   []resumenCategoria `json:"porCategoria"`
}

func (h *Handler) GetResumen(c *gin.Context) {
	query, ok := dateRange(c, h.DB.Where("usuario_id = ?", currentUserID(c)).Model(&models.Gasto{}))
	if !ok {
		return
	}

	var result struct {
		Total          float64
		CantidadGastos int64
	}
	if err := query.Select("COALESCE(SUM(monto), 0) AS total, COUNT(*) AS cantidad_gastos").Scan(&result).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "No se pudo obtener el resumen.")
		return
	}

	categoryQuery, ok := dateRange(c, h.DB.Where("gastos.usuario_id = ?", currentUserID(c)).Model(&models.Gasto{}))
	if !ok {
		return
	}
	var porCategoria []resumenCategoria
	if err := categoryQuery.
		Select("gastos.categoria_id AS categoria_id, categorias.nombre AS categoria, COALESCE(SUM(gastos.monto), 0) AS total").
		Joins("JOIN categorias ON categorias.id = gastos.categoria_id").
		Group("gastos.categoria_id, categorias.nombre").
		Order("categorias.nombre ASC").
		Scan(&porCategoria).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "No se pudo obtener el resumen.")
		return
	}

	if porCategoria == nil {
		porCategoria = []resumenCategoria{}
	}
	c.JSON(http.StatusOK, resumenResponse{Total: result.Total, CantidadGastos: result.CantidadGastos, PorCategoria: porCategoria})
}
