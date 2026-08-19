package models

import "time"

type Usuario struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Nombre       string    `json:"nombre" gorm:"not null;size:100"`
	Email        string    `json:"email" gorm:"not null;uniqueIndex;size:255"`
	PasswordHash string    `json:"-" gorm:"not null"`
	CreatedAt    time.Time `json:"createdAt"`
	Gastos       []Gasto   `json:"-"`
}

func (Usuario) TableName() string {
	return "usuarios"
}
