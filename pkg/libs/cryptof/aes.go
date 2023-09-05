package cryptof

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

type Aes struct{}

var _aes = new(Aes)

func NewAes() *Aes {
	return _aes
}

// CFBEncrypt 加密.
func (slf *Aes) CFBEncrypt(plaintext, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cipher.NewCFBEncrypter(block, iv).XORKeyStream(cipherText[aes.BlockSize:], plaintext)

	return hex.EncodeToString(cipherText), nil
}

// CFBDecrypt 解密.
func (slf *Aes) CFBDecrypt(ciphertext string, key []byte) ([]byte, error) {
	cipherText, err := hex.DecodeString(ciphertext)
	if err != nil {
		return []byte{}, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	if len(cipherText) < aes.BlockSize {
		return []byte{}, fmt.Errorf("cipher ciphertext too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	cipher.NewCFBDecrypter(block, iv).XORKeyStream(cipherText, cipherText)

	return cipherText, nil
}

// CBCEncrypt 加密.
func (slf *Aes) CBCEncrypt(plaintext, key, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	if iv == nil {
		iv = key[:blockSize]
	}
	plaintext = slf.PKCS7Padding(plaintext, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(plaintext))
	blockMode.CryptBlocks(encrypted, plaintext)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// CBCDecrypt 解密.
func (slf *Aes) CBCDecrypt(ciphertext string, key, iv []byte) ([]byte, error) {
	ciphertextByte, _ := base64.StdEncoding.DecodeString(ciphertext)
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	blockSize := block.BlockSize()
	if iv == nil {
		iv = key[:blockSize]
	}
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData := make([]byte, len(ciphertextByte))
	blockMode.CryptBlocks(origData, ciphertextByte)
	origData = slf.PKCS7UnPadding(origData)

	return origData, nil
}

func (slf *Aes) PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func (slf *Aes) PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
