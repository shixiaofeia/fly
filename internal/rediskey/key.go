package rediskey

import "fmt"

type Key struct {
	Prefix string
}

// NewKey
func NewKey() *Key {
	return &Key{Prefix: Prefix}
}

// GetUserToken 获取用户token.
func (k *Key) GetUserToken(userId uint32) string {
	return fmt.Sprintf("%s:%s:%d", k.Prefix, UserToken, userId)
}
