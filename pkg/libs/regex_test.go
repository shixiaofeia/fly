package libs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckPwd(t *testing.T) {
	assert.Equal(t, CheckPwd("Fly123!@#$"), true)
	assert.Equal(t, CheckPwd("Fly123333"), true)
	assert.Equal(t, CheckPwd("Ff23456789"), true)
	assert.Equal(t, CheckPwd("123456789f!"), true)
	assert.Equal(t, CheckPwd("f123456789!"), true)
	assert.Equal(t, CheckPwd("FlyWelcome"), false)
	assert.Equal(t, CheckPwd("123456789"), false)
	assert.Equal(t, CheckPwd("F23456789"), false)
	assert.Equal(t, CheckPwd("Ff2啊56789"), false)
}
