package infisical

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (c *InfisicalClient) UpdateSecret(secretKey, secretValue, secretComment string) error {
	encryptedKey, err := Encrypt(secretKey, c.projectKey)
	if err != nil {
		return err
	}

	encryptedValue, err := Encrypt(secretValue, c.projectKey)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"workspaceId":           c.workspace,
		"environment":           c.environment,
		"type":                  c.secretType,
		"secretKeyCiphertext":   encryptedKey.Ciphertext,
		"secretKeyIV":           encryptedKey.IV,
		"secretValueCiphertext": encryptedValue.Ciphertext,
		"secretValueIV":         encryptedValue.IV,
	}

	if len(secretComment) > 0 {
		encryptedCommentData, err := Encrypt(secretComment, c.projectKey)
		if err != nil {
			panic(err)
		}

		payload["secretCommentCiphertext"] = encryptedCommentData.Ciphertext
		payload["secretCommentIV"] = encryptedCommentData.IV
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("PATCH", BaseURL+"/api/v3/secrets/"+secretKey, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Authorization", "Bearer "+c.serviceToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("failed to update secret")
	}

	fmt.Println("Secret updated successfully.")
	return nil
}
