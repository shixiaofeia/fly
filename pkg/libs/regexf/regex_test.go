package regexf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var factory = NewRegex()

func TestRegex_Pwd(t *testing.T) {
	assert.Equal(t, factory.Pwd("Fly123!@#$"), true)
	assert.Equal(t, factory.Pwd("Fly123333"), true)
	assert.Equal(t, factory.Pwd("Ff23456789"), true)
	assert.Equal(t, factory.Pwd("123456789f!"), true)
	assert.Equal(t, factory.Pwd("f123456789!"), true)
	assert.Equal(t, factory.Pwd("FlyWelcome"), false)
	assert.Equal(t, factory.Pwd("123456789"), false)
	assert.Equal(t, factory.Pwd("F23456789"), false)
	assert.Equal(t, factory.Pwd("Ff2啊56789"), false)
}

func TestRegex_Email(t *testing.T) {
	assert.Equal(t, factory.Email("fly@test.com"), true)
	assert.Equal(t, factory.Email("fly@test"), false)
}
