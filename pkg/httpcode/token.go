package httpcode

import (
	"fly/internal/const"
	"github.com/dgrijalva/jwt-go"
)

// CreateJWTToken 生成token
func CreateJWTToken(data map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	for index, val := range data {
		claims[index] = val
	}
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(constants.JwtSecretKey))
	return tokenString, err
}

// ParseToken 解析token
func ParseToken(token string) (userId uint32, err error) {
	tokenInfo, _ := jwt.Parse(token, func(token *jwt.Token) (i interface{}, e error) {
		return constants.JwtSecretKey, nil
	})
	err = tokenInfo.Claims.Valid()
	if err != nil {
		return
	}
	tokenMap := tokenInfo.Claims.(jwt.MapClaims)
	userId = uint32(tokenMap["user_id"].(float64))
	return
}
