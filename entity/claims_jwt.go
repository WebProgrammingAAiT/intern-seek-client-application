package entity

import "github.com/dgrijalva/jwt-go"

//Claims is used for setting jwt token payload
type Claims struct {
	UserID uint
	Name   string
	jwt.StandardClaims
}
