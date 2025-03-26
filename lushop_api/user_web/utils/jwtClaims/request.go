package jwtClaims

import "github.com/golang-jwt/jwt"

type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint
	jwt.StandardClaims
}
