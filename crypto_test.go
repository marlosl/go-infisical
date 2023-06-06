package infisical

import (
	"testing"

	"github.com/marlosl/go-infisical/crypto"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	client := &InfisicalClient{}
	text := "hello world"

	newKey, _ := crypto.GenerateNewKey()
	secret := string(newKey)

	encryptionData, err := client.encrypt(text, string(secret))
	assert.NoError(t, err)

	decryptedText, err := client.decrypt(encryptionData.Ciphertext, encryptionData.IV, encryptionData.Tag, secret)
	assert.NoError(t, err)

	assert.Equal(t, text, decryptedText)
}

func TestDecrypt(t *testing.T) {
	client := &InfisicalClient{}
	text := "hello world"

	newKey, _ := crypto.GenerateNewKey()
	secret := string(newKey)

	encryptionData, err := client.encrypt(text, secret)
	assert.NoError(t, err)

	decryptedText, err := client.decrypt(encryptionData.Ciphertext, encryptionData.IV, encryptionData.Tag, secret)
	assert.NoError(t, err)

	assert.Equal(t, text, decryptedText)
}
