package infisical

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/marlosl/go-infisical/cache"
)

type InfisicalClient struct {
	serviceToken       string
	serviceTokenSecret string
	encryptedKey       string
	iv                 string
	tag                string
	projectKey         string
	workspace          string
	environment        string
	secretType         string
	httpClient         *http.Client
	cache              *cache.Cache
}

func NewClient(serviceToken string) (*InfisicalClient, error) {
	var err error
	c := &InfisicalClient{}

	c.serviceToken = serviceToken
	c.httpClient = &http.Client{}

	c.cache, err = cache.NewCache()
	if err != nil {
		return nil, err
	}

	c.InitServiceTokenSecret()

	err = c.InitService()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *InfisicalClient) InitServiceTokenSecret() {
	if len(c.serviceToken) > 0 {
		return
	}
	c.serviceTokenSecret = c.serviceToken[strings.LastIndex(c.serviceToken, ".")+1:]
}

func (c *InfisicalClient) InitService() error {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", BaseURL+"/api/v2/service-token", nil)
	req.Header.Set("Authorization", "Bearer "+c.serviceToken)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	serviceTokenData := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&serviceTokenData)
	if err != nil {
		return err
	}

	c.encryptedKey = serviceTokenData["encryptedKey"].(string)
	c.iv = serviceTokenData["iv"].(string)
	c.tag = serviceTokenData["tag"].(string)
	c.workspace = serviceTokenData["workspace"].(string)
	c.environment = serviceTokenData["environment"].(string)

	projectKey, err := c.decrypt(
		c.encryptedKey,
		c.iv,
		c.tag,
		c.serviceTokenSecret,
	)
	if err != nil {
		return err
	}

	c.projectKey = projectKey
	c.secretType = SecretType

	return nil
}
