package jwt

import (
	"github.com/golang-jwt/jwt/v4"
)

// CreateJWTToken 生成token.
func CreateJWTToken(data map[string]interface{}) (string, error) {
	var (
		token  = jwt.New(jwt.SigningMethodHS256)
		claims = make(jwt.MapClaims)
	)

	for index, val := range data {
		claims[index] = val
	}
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(SecretKey))
	return tokenString, err
}

// ParseToken 解析token.
func ParseToken(token string) (userId string, err error) {
	var tokenInfo *jwt.Token

	if tokenInfo, err = jwt.Parse(token, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(SecretKey), err
	}); err != nil {
		return
	}
	if err = tokenInfo.Claims.Valid(); err != nil {
		return
	}
	tokenMap := tokenInfo.Claims.(jwt.MapClaims)
	if val, ok := tokenMap["user_id"]; ok {
		if uid, ok := val.(string); ok {
			userId = uid
		}
	}
	return
}
