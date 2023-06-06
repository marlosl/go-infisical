package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

const (
	BlocksizeBytes = 16
)

type EncryptionResult struct {
	CipherText []byte
	AuthTag    []byte
	Nonce      []byte
}

func Decrypt(key []byte, cipherText []byte, tag []byte, iv []byte) ([]byte, error) {

	if len(cipherText) == 0 && len(tag) == 0 && len(iv) == 0 {
		return []byte{}, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCMWithNonceSize(block, len(iv))
	if err != nil {
		return nil, err
	}

	var nonce = iv
	var ciphertext = append(cipherText, tag...)

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func Encrypt(plaintext []byte, key []byte) (result EncryptionResult, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return EncryptionResult{}, err
	}

	aesgcm, err := cipher.NewGCMWithNonceSize(block, BlocksizeBytes)
	if err != nil {
		return EncryptionResult{}, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	ciphertextOnly := ciphertext[:len(ciphertext)-BlocksizeBytes]

	authTag := ciphertext[len(ciphertext)-BlocksizeBytes:]

	return EncryptionResult{
		CipherText: ciphertextOnly,
		AuthTag:    authTag,
		Nonce:      nonce,
	}, nil
}

func GenerateNewKey() (newKey []byte, keyErr error) {
	key := make([]byte, BlocksizeBytes)
	_, err := rand.Read(key)
	return key, err
}
