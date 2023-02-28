package service

import (
	"api/internal/repository/postgres"
	"api/pkg/entities"
	"errors"
	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type AuthService struct {
	repo postgres.Authenthication
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

func (c *AuthService) CreateUser(user entities.User) (int, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return 0, err
	}
	user.Password = string(password)

	return c.repo.CreateUser(user)
}

func (c *AuthService) GenerateToken(user entities.User) (string, error) {
	userCRD, _ := c.repo.ValidateUser(user)
	err := bcrypt.CompareHashAndPassword([]byte(userCRD.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		userCRD.Id,
	})

	return token.SignedString([]byte(os.Getenv("privateKey")))
}

func (c *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(os.Getenv("privateKey")), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func NewAuthService(repo postgres.Authenthication) *AuthService {
	return &AuthService{repo: repo}
}
