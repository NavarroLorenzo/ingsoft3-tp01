package handlers

import (
	"errors"
	"net/http"

	"gestor-gastos/backend/internal/auth"
	"gestor-gastos/backend/internal/models"
	"gestor-gastos/backend/internal/validation"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type authResponse struct {
	User  models.Usuario `json:"user"`
	Token string         `json:"token"`
}

func (h *Handler) Register(c *gin.Context) {
	var input validation.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, "El cuerpo de la solicitud es inválido.")
		return
	}
	if err := validation.ValidateRegister(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var existing models.Usuario
	err := h.DB.Where("LOWER(email) = LOWER(?)", input.Email).First(&existing).Error
	if err == nil {
		writeError(c, http.StatusConflict, "Ya existe un usuario con ese email.")
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		writeError(c, http.StatusInternalServerError, "No se pudo registrar el usuario.")
		return
	}

	hash, err := auth.HashPassword(input.Password)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "No se pudo registrar el usuario.")
		return
	}
	user := models.Usuario{Nombre: input.Nombre, Email: input.Email, PasswordHash: hash}
	if err := h.DB.Create(&user).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "No se pudo registrar el usuario.")
		return
	}
	token, err := auth.GenerateToken(user, h.JWTSecret)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "No se pudo generar el token.")
		return
	}
	c.JSON(http.StatusCreated, authResponse{User: user, Token: token})
}

func (h *Handler) Login(c *gin.Context) {
	var input validation.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, "El cuerpo de la solicitud es inválido.")
		return
	}
	if err := validation.ValidateLogin(&input); err != nil {
		writeError(c, http.StatusUnauthorized, "Email o contraseña incorrectos.")
		return
	}

	var user models.Usuario
	if err := h.DB.Where("LOWER(email) = LOWER(?)", input.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusUnauthorized, "Email o contraseña incorrectos.")
			return
		}
		writeError(c, http.StatusInternalServerError, "No se pudo iniciar sesión.")
		return
	}
	if auth.ComparePassword(user.PasswordHash, input.Password) != nil {
		writeError(c, http.StatusUnauthorized, "Email o contraseña incorrectos.")
		return
	}
	token, err := auth.GenerateToken(user, h.JWTSecret)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "No se pudo generar el token.")
		return
	}
	c.JSON(http.StatusOK, authResponse{User: user, Token: token})
}

func (h *Handler) Me(c *gin.Context) {
	userID := currentUserID(c)
	var user models.Usuario
	if err := h.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusUnauthorized, "Usuario no encontrado.")
			return
		}
		writeError(c, http.StatusInternalServerError, "No se pudo obtener el usuario.")
		return
	}
	c.JSON(http.StatusOK, user)
}
