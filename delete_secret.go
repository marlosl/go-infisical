package infisical

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

func (c *InfisicalClient) DeleteSecret(secretKey string) error {
	payload := map[string]interface{}{
		"workspaceId": c.workspace,
		"environment": c.environment,
		"type":        c.secretType,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("DELETE", BaseURL+"/api/v2/secrets/"+secretKey, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Authorization", "Bearer "+c.serviceToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("failed to delete secret")
	}

	return nil
}
