package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"golang.org/x/crypto/nacl/box"
)

type SymmetricEncryptionResult struct {
	CipherText []byte
	AuthTag    []byte
	Nonce      []byte
}

func DecryptSymmetric(key []byte, cipherText []byte, tag []byte, iv []byte) ([]byte, error) {

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

func GenerateNewKey() (newKey []byte, keyErr error) {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	return key, err
}

func EncryptSymmetric(plaintext []byte, key []byte) (result SymmetricEncryptionResult, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return SymmetricEncryptionResult{}, err
	}

	aesgcm, err := cipher.NewGCMWithNonceSize(block, 16)
	if err != nil {
		return SymmetricEncryptionResult{}, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	ciphertextOnly := ciphertext[:len(ciphertext)-16]

	authTag := ciphertext[len(ciphertext)-16:]

	return SymmetricEncryptionResult{
		CipherText: ciphertextOnly,
		AuthTag:    authTag,
		Nonce:      nonce,
	}, nil
}

func DecryptAsymmetric(ciphertext []byte, nonce []byte, publicKey []byte, privateKey []byte) (plainText []byte) {
	plainTextToReturn, _ := box.Open(nil, ciphertext, (*[24]byte)(nonce), (*[32]byte)(publicKey), (*[32]byte)(privateKey))
	return plainTextToReturn
}

func EncryptAssymmetric(message []byte, nonce []byte, publicKey []byte, privateKey []byte) (encryptedMessage []byte) {
	encryptedPlainText := box.Seal(nil, message, (*[24]byte)(nonce), (*[32]byte)(publicKey), (*[32]byte)(privateKey))
	return encryptedPlainText
}
