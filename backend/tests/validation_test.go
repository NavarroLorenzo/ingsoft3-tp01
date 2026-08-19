package tests

import (
	"testing"

	"gestor-gastos/backend/internal/validation"
)

func validInput() validation.GastoInput {
	return validation.GastoInput{
		Descripcion: "Supermercado semanal",
		Monto:       2500.50,
		Fecha:       "2026-08-12",
		CategoriaID: 1,
	}
}

func TestValidateGasto(t *testing.T) {
	tests := []struct {
		name  string
		input validation.GastoInput
		valid bool
	}{
		{name: "monto negativo", input: func() validation.GastoInput { input := validInput(); input.Monto = -1; return input }(), valid: false},
		{name: "monto cero", input: func() validation.GastoInput { input := validInput(); input.Monto = 0; return input }(), valid: false},
		{name: "descripcion vacia", input: func() validation.GastoInput { input := validInput(); input.Descripcion = ""; return input }(), valid: false},
		{name: "descripcion corta", input: func() validation.GastoInput { input := validInput(); input.Descripcion = "ab"; return input }(), valid: false},
		{name: "monto con tres decimales", input: func() validation.GastoInput { input := validInput(); input.Monto = 1.123; return input }(), valid: false},
		{name: "gasto valido", input: validInput(), valid: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validation.ValidateGasto(&tt.input)
			if (err == nil) != tt.valid {
				t.Fatalf("ValidateGasto() error = %v, valid = %v", err, tt.valid)
			}
		})
	}
}
