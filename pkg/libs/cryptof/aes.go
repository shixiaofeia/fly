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
func (slf *Aes) CFBEncrypt(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cipher.NewCFBEncrypter(block, iv).XORKeyStream(cipherText[aes.BlockSize:], []byte(plaintext))

	return hex.EncodeToString(cipherText), nil
}

// CFBDecrypt 解密.
func (slf *Aes) CFBDecrypt(ciphertext string, key []byte) (string, error) {
	cipherText, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("cipher ciphertext too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	cipher.NewCFBDecrypter(block, iv).XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

// CBCEncrypt 加密.
func (slf *Aes) CBCEncrypt(plaintext string, key, iv []byte) (string, error) {
	plaintextByte := []byte(plaintext)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	if iv == nil {
		iv = key[:blockSize]
	}
	plaintextByte = slf.PKCS7Padding(plaintextByte, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(plaintextByte))
	blockMode.CryptBlocks(crypted, plaintextByte)

	return base64.StdEncoding.EncodeToString(crypted), nil
}

// CBCDecrypt 解密.
func (slf *Aes) CBCDecrypt(ciphertext string, key, iv []byte) (string, error) {
	ciphertextByte, _ := base64.StdEncoding.DecodeString(ciphertext)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	if iv == nil {
		iv = key[:blockSize]
	}
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData := make([]byte, len(ciphertextByte))
	blockMode.CryptBlocks(origData, ciphertextByte)
	origData = slf.PKCS7UnPadding(origData)

	return string(origData), nil
}

func (slf *Aes) PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (slf *Aes) PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
