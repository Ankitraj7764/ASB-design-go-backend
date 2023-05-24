package models

import "github.com/golang-jwt/jwt"

type JWTData struct {
	jwt.StandardClaims
	CustmClaims map[string]string `json:"custom-claims"`
}
