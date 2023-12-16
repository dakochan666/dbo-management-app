package service

import (
	"dbo-management-app/models"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Exp      int64  `json:"exp"`
	jwt.StandardClaims
}

var jwtKey = []byte(os.Getenv("SECRET"))

func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hashedPassword, err
}

func ComparePassword(hashPassword []byte, password []byte) error {
	err := bcrypt.CompareHashAndPassword(hashPassword, password)
	return err
}

func SignToken(user models.User) (string, error) {
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
		Exp:   time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign token
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

func VerifyToken(tokenString string) (*Claims, error) {
	// Verify Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
