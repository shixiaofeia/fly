package stringf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetInviteCode(t *testing.T) {
	assert.Equal(t, NewInviteCode().IdToCode(12345678)[:6], "LE338B")
}

func TestDecodeInviteCode(t *testing.T) {
	assert.Equal(t, NewInviteCode().CodeToId("LE338B2D"), uint32(12345678))
}
