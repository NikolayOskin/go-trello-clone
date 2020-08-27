package auth

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/dgrijalva/jwt-go"
)

type JWTService struct {
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
	jwtTTLHours time.Duration
}

func NewJWTService(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, jwtTTLHours time.Duration) (*JWTService, error) {
	if privateKey == nil {
		return nil, errors.New("PrivateKey cannot be nil")
	}
	if publicKey == nil {
		return nil, errors.New("PublicKey cannot be nil")
	}

	return &JWTService{
		privateKey,
		publicKey,
		jwtTTLHours,
	}, nil
}

func (a *JWTService) GenerateToken(user model.User) (string, error) {
	claims := model.JWTClaims{
		User: model.User{
			ID:    user.ID,
			Email: user.Email,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(a.jwtTTLHours * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := token.SignedString(a.privateKey)

	return tokenStr, err
}

func (a *JWTService) ValidateToken(tokenStr string) (*model.JWTClaims, error) {
	claims := model.JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return a.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(*model.JWTClaims); ok && !token.Valid {
		return nil, err
	}
	return &claims, nil
}
