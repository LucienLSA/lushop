package jwtClaims

import "github.com/golang-jwt/jwt"

type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint
	DeviceID    string // 设备指纹字段
	*jwt.StandardClaims
}
