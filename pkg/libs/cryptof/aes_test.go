package cryptof

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAes_CFB(t *testing.T) {
	var (
		factory = NewAes()
		key     = []byte("fly6666666666666")
		data    = "what are you 弄啥嘞"
	)

	ciphertext, err := factory.CFBEncrypt([]byte(data), key)
	if err != nil {
		t.Error(err)
	}

	plaintext, err := factory.CFBDecrypt(ciphertext, key)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, string(plaintext), data)
}

func TestAes_CBC(t *testing.T) {
	var (
		factory = NewAes()
		key     = []byte("fly6666666666666")
		data    = "what are you 弄啥嘞"
	)

	ciphertext, err := factory.CBCEncrypt([]byte(data), key, nil)
	if err != nil {
		t.Error(err)
	}

	plaintext, err := factory.CBCDecrypt(ciphertext, key, nil)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, string(plaintext), data)
}
