package cryptof

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
)

type Hash struct{}

var _hash = new(Hash)

func NewHash() *Hash {
	return _hash
}

func (slf *Hash) Md5(data string) string {
	h := md5.New()
	return slf.toString(h, data)
}

func (slf *Hash) Sha1(data string) string {
	h := sha1.New()
	return slf.toString(h, data)
}

func (slf *Hash) Sha224(data string) string {
	h := sha256.New224()
	return slf.toString(h, data)
}

func (slf *Hash) Sha256(data string) string {
	h := sha256.New()
	return slf.toString(h, data)
}

func (slf *Hash) Sha384(data string) string {
	h := sha512.New384()
	return slf.toString(h, data)
}

func (slf *Hash) Sha512(data string) string {
	h := sha512.New()
	return slf.toString(h, data)
}

func (*Hash) toString(h hash.Hash, data string) string {
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
