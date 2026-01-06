package utils

import (
	"fmt"
	"time"

	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/config"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/models"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(userID int, email string) (string, error) {
	claims := models.Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.Cfg.JWTSecret)
}

func ValidateJWT(tokenString string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(t *jwt.Token) (interface{}, error) {
		return config.Cfg.JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
