package models

import (
	"encoding/json"
	"time"
)

type Gasto struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Descripcion string    `json:"descripcion" gorm:"not null;size:200"`
	Monto       float64   `json:"monto" gorm:"type:numeric(12,2);not null"`
	Fecha       time.Time `json:"fecha" gorm:"not null"`
	CategoriaID uint      `json:"categoriaId" gorm:"not null;index"`
	Categoria   Categoria `json:"categoria" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	UsuarioID   uint      `json:"-" gorm:"index"`
	Usuario     Usuario   `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (Gasto) TableName() string {
	return "gastos"
}

// MarshalJSON keeps the public API date-only as required by the REST contract.
func (g Gasto) MarshalJSON() ([]byte, error) {
	type gastoJSON struct {
		ID          uint      `json:"id"`
		Descripcion string    `json:"descripcion"`
		Monto       float64   `json:"monto"`
		Fecha       string    `json:"fecha"`
		CategoriaID uint      `json:"categoriaId"`
		Categoria   Categoria `json:"categoria"`
	}
	return json.Marshal(gastoJSON{
		ID:          g.ID,
		Descripcion: g.Descripcion,
		Monto:       g.Monto,
		Fecha:       g.Fecha.Format("2006-01-02"),
		CategoriaID: g.CategoriaID,
		Categoria:   g.Categoria,
	})
}
