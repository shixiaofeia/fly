package utils

import (
	"testing"
)

func TestGoogleAuth(t *testing.T) {
	auth := NewGoogleAuth()
	// 获取google秘钥
	secret := auth.GetSecret()
	t.Log(secret)
}
