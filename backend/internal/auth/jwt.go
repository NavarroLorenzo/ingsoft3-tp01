package auth

import (
	"fmt"
	"time"

	"gestor-gastos/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

const TokenDuration = 24 * time.Hour

type Claims struct {
	UserID uint   `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(user models.Usuario, secret string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(TokenDuration)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func ParseToken(tokenString, secret string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma no permitido")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid || claims.UserID == 0 {
		if err == nil {
			err = fmt.Errorf("token inválido")
		}
		return nil, err
	}
	return claims, nil
}
