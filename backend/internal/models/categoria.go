package models

type Categoria struct {
	ID     uint    `json:"id" gorm:"primaryKey"`
	Nombre string  `json:"nombre" gorm:"not null;size:50"`
	Gastos []Gasto `json:"-"`
}

func (Categoria) TableName() string {
	return "categorias"
}
