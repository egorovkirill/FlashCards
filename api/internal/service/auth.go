package service

import (
	"ToDo/internal/repository"
	"ToDo/pkg/entities"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type AuthService struct {
	repo repository.Authenthication
}

type tokenClaims struct {
	jwt.StandardClaims
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
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			IssuedAt:  time.Now().Unix(),
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

func NewAuthService(repo repository.Authenthication) *AuthService {
	return &AuthService{repo: repo}
}
