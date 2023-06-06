package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	key := []byte("mysecret")
	plaintext := []byte("hello world")

	encryptionResult, err := Encrypt(plaintext, key)
	assert.NoError(t, err)

	decryptedPlaintext, err := Decrypt(key, encryptionResult.CipherText, encryptionResult.AuthTag, encryptionResult.Nonce)
	assert.NoError(t, err)

	assert.Equal(t, plaintext, decryptedPlaintext)
}

func TestDecrypt(t *testing.T) {
	key := []byte("mysecret")
	plaintext := []byte("hello world")

	encryptionResult, err := Encrypt(plaintext, key)
	assert.NoError(t, err)

	decryptedPlaintext, err := Decrypt(key, encryptionResult.CipherText, encryptionResult.AuthTag, encryptionResult.Nonce)
	assert.NoError(t, err)

	assert.Equal(t, plaintext, decryptedPlaintext)
}
