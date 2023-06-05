package infisical

import (
	"encoding/json"
	"net/http"
)

func (c *InfisicalClient) GetSecret(secretKey string) (string, error) {

	req, _ := http.NewRequest("GET", BaseURL+"/api/v3/secrets/"+secretKey, nil)
	req.Header.Set("Authorization", "Bearer "+c.serviceToken)
	q := req.URL.Query()
	q.Add("environment", c.environment)
	q.Add("workspaceId", c.workspace)
	q.Add("type", c.secretType)

	req.URL.RawQuery = q.Encode()

	res, err := c.client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	secretData := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&secretData)
	if err != nil {
		return "", err
	}

	encryptedSecret := secretData["secret"].(map[string]interface{})

	secretValue, err := Decrypt(
		encryptedSecret["secretValueCiphertext"].(string),
		encryptedSecret["secretValueIV"].(string),
		encryptedSecret["secretValueTag"].(string),
		c.projectKey,
	)
	if err != nil {
		return "", err
	}

	return secretValue, nil
}

func (c *InfisicalClient) GetCachedSecret(secretKey string) (string, error) {
	value, err := c.cache.Read(secretKey)
	if err != nil {
		value, err = c.GetSecret(secretKey)
		if err != nil {
			return "", err
		}
		c.cache.Update(secretKey, value)
	}
	return value, nil
}

func (c *InfisicalClient) GetSecrets() (map[string]string, error) {

	secrets := make(map[string]string)
	req, _ := http.NewRequest("GET", BaseURL+"/api/v3/secrets", nil)
	req.Header.Set("Authorization", "Bearer "+c.serviceToken)
	q := req.URL.Query()
	q.Add("environment", c.environment)
	q.Add("workspaceId", c.workspace)
	q.Add("type", c.secretType)

	req.URL.RawQuery = q.Encode()

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	secretData := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&secretData)
	if err != nil {
		return nil, err
	}

	encryptedSecrets := secretData["secrets"].([]interface{})
	for _, encryptedSecret := range encryptedSecrets {
		secret := encryptedSecret.(map[string]interface{})

		secretKey, err := Decrypt(
			secret["secretKeyCiphertext"].(string),
			secret["secretKeyIV"].(string),
			secret["secretKeyTag"].(string),
			c.projectKey,
		)
		if err != nil {
			return nil, err
		}

		secretValue, err := Decrypt(
			secret["secretValueCiphertext"].(string),
			secret["secretValueIV"].(string),
			secret["secretValueTag"].(string),
			c.projectKey,
		)
		if err != nil {
			return nil, err
		}

		secrets[secretKey] = secretValue
	}

	return secrets, nil
}
