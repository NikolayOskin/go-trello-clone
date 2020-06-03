package model

import "github.com/dgrijalva/jwt-go"

type JWTClaims struct {
	User User
	jwt.StandardClaims
}
