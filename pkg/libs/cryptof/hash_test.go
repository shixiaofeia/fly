package cryptof

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewHash(t *testing.T) {
	var (
		factory = NewHash()
		text    = "想你想你想我"
	)

	assert.Equal(t, factory.Md5(text), "aa587fd8f966e3ab1ca24737da504cba")
	assert.Equal(t, factory.Sha1(text), "16eb06bbcbf56911ba16e904599fcebfaee5e10f")
	assert.Equal(t, factory.Sha224(text), "cb4cc8520a03112ee0433579e9b02875a10616c9097e68c3ce81145c")
	assert.Equal(t, factory.Sha256(text), "0d434250c0a26a5034205b78810b376436ec1e2d0bba2b49974912f03b8d90b7")
	assert.Equal(t, factory.Sha384(text), "08ed2d6323b6a7346e72e9cd806956ac1a41fe79021980a80e00f7e2eaaf58f482b888df5aeba4fcbf26cb56f009d1de")
	assert.Equal(t, factory.Sha512(text), "71b71e7e52a07849b7fd7be07b6b3a4b0fd7850c85b9b601b7cffb1f9e1ba43799b128ce7a3d7bf170b4c02d51ad6387acd0b3ee26cde1dcc62fd5bb35a9e7d1")
}
