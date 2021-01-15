package service

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/leor-w/ihome/user-service/model"
	"github.com/leor-w/ihome/user-service/repo"
)

type CustomClaims struct {
	User *model.User
	jwt.StandardClaims
}

type AuthInterface interface {
	Decode(token string) (*CustomClaims, error)
	Encode(userModel *model.User) (string, error)
}

type Config struct {
	TokenKey string
	Issuer   string
}

type AuthService struct {
	Conf Config
	Repo repo.UserRepositoryInterface
}

// Decode 将 jwt 字符串校验并解码为 model.User 模型
func (svc *AuthService) Decode(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return svc.Conf.TokenKey, nil
	})

	if token == nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// Encode 将 model.User 模型编码为 jwt 字符串
func (svc *AuthService) Encode(user *model.User) (string, error) {
	expireToken := time.Now().Add(time.Hour * 72).Unix()
	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    svc.Conf.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(svc.Conf.TokenKey)
}
