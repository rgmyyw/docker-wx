package Algorithm

import (
	"crypto/aes"
	"crypto/cipher"
)

func NewAES_GCMEncrypter(key []byte, plaintext []byte, nonce []byte, aad []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil
	}
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, aad)
	return ciphertext
}

func NewAES_GCMDecrypter(key []byte, ciphertext []byte, nonce []byte, aad []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, aad)
	if err != nil {
		return nil
	}
	return plaintext
}
