package infisical

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *InfisicalClient) CreateSecret(secretKey, secretValue, secretComment string) error {
	encryptedKeyData, err := c.encrypt(secretKey, c.projectKey)
	if err != nil {
		return err
	}

	encryptedValueData, err := c.encrypt(secretValue, c.projectKey)
	if err != nil {
		return err
	}

	secretData := map[string]interface{}{
		"workspaceId":           c.workspace,
		"environment":           c.environment,
		"type":                  c.secretType,
		"secretKeyCiphertext":   encryptedKeyData.Ciphertext,
		"secretKeyIV":           encryptedKeyData.IV,
		"secretKeyTag":          encryptedKeyData.Tag,
		"secretValueCiphertext": encryptedValueData.Ciphertext,
		"secretValueIV":         encryptedValueData.IV,
		"secretValueTag":        encryptedValueData.Tag,
	}

	if len(secretComment) > 0 {
		encryptedCommentData, err := c.encrypt(secretComment, c.projectKey)
		if err != nil {
			panic(err)
		}

		secretData["secretCommentCiphertext"] = encryptedCommentData.Ciphertext
		secretData["secretCommentIV"] = encryptedCommentData.IV
		secretData["secretCommentTag"] = encryptedCommentData.Tag
	}

	secretJSON, err := json.Marshal(secretData)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", BaseURL+"/api/v3/secrets/"+secretKey, bytes.NewBuffer(secretJSON))
	req.Header.Set("Authorization", "Bearer "+c.serviceToken)
	req.Header.Set("Content-Type", "application/json")

	_, err = c.httpClient.Do(req)
	if err != nil {
		return err
	}

	fmt.Println("Secret created successfully.")
	return nil
}
