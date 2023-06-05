package infisical

import (
	"encoding/base64"
	"fmt"

	"github.com/marlosl/go-infisical/crypto"
)

const (
	BaseURL        = "https://app.infisical.com"
	BlocksizeBytes = 16
	SharedSecret   = "shared"
	PersonalSecret = "personal"
	secretType     = "personal"
)

var (
	serviceToken       = ""
	serviceTokenSecret = ""
)

type EncryptionData struct {
	Ciphertext string `json:"ciphertext"`
	Tag        string `json:"tag"`
	IV         string `json:"iv"`
}

func Encrypt(text string, secret string) (*EncryptionData, error) {
	key := []byte(secret)
	plaintext := []byte(text)

	cipher, err := crypto.Encrypt(plaintext, key)
	if err != nil {
		panic(err.Error())
	}

	encryptionData := &EncryptionData{
		Ciphertext: base64.StdEncoding.EncodeToString(cipher.CipherText),
		Tag:        base64.StdEncoding.EncodeToString(cipher.AuthTag),
		IV:         base64.StdEncoding.EncodeToString(cipher.Nonce),
	}

	return encryptionData, nil
}

func Decrypt(ciphertext string, iv string, tag string, secret string) (string, error) {
	secretKey := []byte(secret)

	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	decodedIV, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return "", err
	}

	decodedTag, err := base64.StdEncoding.DecodeString(tag)
	if err != nil {
		return "", err
	}

	plaintext, err := crypto.Decrypt(secretKey, decodedCiphertext, decodedTag, decodedIV)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
