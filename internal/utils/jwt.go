package utils

import (
	"converse/internal/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

func GenerateJWT(user *models.User) (string, time.Time, error) {
    jwtSecret := []byte(os.Getenv("JWT_SECRET"))

    expiresAt := time.Now().Add(24 * time.Hour)

    claims := &Claims{
        UserID: user.UserID,
        Email:  user.Email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "converse",
            Subject:   user.UserID,
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        return "", time.Time{}, err
    }

    return tokenString, expiresAt, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
    jwtSecret := []byte(os.Getenv("JWT_SECRET"))

    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, jwt.ErrSignatureInvalid
}