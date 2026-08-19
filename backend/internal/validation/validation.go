package validation

import (
	"errors"
	"math"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"
)

type CategoriaInput struct {
	Nombre string `json:"nombre"`
}

type GastoInput struct {
	Descripcion string  `json:"descripcion"`
	Monto       float64 `json:"monto"`
	Fecha       string  `json:"fecha"`
	CategoriaID uint    `json:"categoriaId"`
}

type RegisterInput struct {
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ValidateRegister(input *RegisterInput) error {
	input.Nombre = strings.TrimSpace(input.Nombre)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if length := utf8.RuneCountInString(input.Nombre); length < 2 || length > 100 {
		return errors.New("El nombre debe tener entre 2 y 100 caracteres.")
	}
	if err := validateEmail(input.Email); err != nil {
		return err
	}
	if err := validatePassword(input.Password); err != nil {
		return err
	}
	return nil
}

func ValidateLogin(input *LoginInput) error {
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if err := validateEmail(input.Email); err != nil {
		return err
	}
	if input.Password == "" {
		return errors.New("La contraseña es obligatoria.")
	}
	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return errors.New("El email es obligatorio.")
	}
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email {
		return errors.New("El email no tiene un formato válido.")
	}
	return nil
}

func validatePassword(password string) error {
	length := utf8.RuneCountInString(password)
	if length == 0 {
		return errors.New("La contraseña es obligatoria.")
	}
	if length < 8 || len(password) > 72 {
		return errors.New("La contraseña debe tener entre 8 y 72 caracteres.")
	}
	return nil
}

func ValidateCategoria(input *CategoriaInput) error {
	input.Nombre = strings.TrimSpace(input.Nombre)
	length := utf8.RuneCountInString(input.Nombre)

	if length == 0 {
		return errors.New("El nombre de la categoría es obligatorio.")
	}
	if length < 2 || length > 50 {
		return errors.New("El nombre de la categoría debe tener entre 2 y 50 caracteres.")
	}
	return nil
}

func ValidateGasto(input *GastoInput) (time.Time, error) {
	input.Descripcion = strings.TrimSpace(input.Descripcion)
	length := utf8.RuneCountInString(input.Descripcion)

	if length == 0 {
		return time.Time{}, errors.New("La descripción es obligatoria.")
	}
	if length < 3 || length > 200 {
		return time.Time{}, errors.New("La descripción debe tener entre 3 y 200 caracteres.")
	}
	if input.Monto <= 0 {
		return time.Time{}, errors.New("El monto debe ser mayor que cero.")
	}
	if math.Abs(input.Monto*100-math.Round(input.Monto*100)) > 0.000001 {
		return time.Time{}, errors.New("El monto puede tener como máximo dos decimales.")
	}
	if input.Fecha == "" {
		return time.Time{}, errors.New("La fecha es obligatoria.")
	}
	fecha, err := time.Parse("2006-01-02", input.Fecha)
	if err != nil {
		return time.Time{}, errors.New("La fecha debe tener el formato YYYY-MM-DD.")
	}
	if input.CategoriaID == 0 {
		return time.Time{}, errors.New("La categoría es obligatoria.")
	}
	return fecha, nil
}

func ParseDate(value string) (time.Time, error) {
	fecha, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, errors.New("La fecha debe tener el formato YYYY-MM-DD.")
	}
	return fecha, nil
}
