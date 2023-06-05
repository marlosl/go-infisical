package infisical

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/marlosl/go-infisical/cache"
	"github.com/marlosl/go-infisical/consts"
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
	client             *http.Client
	cache              *cache.Cache
}

func CreateInfisicalClient() (*InfisicalClient, error) {
	var err error
	c := &InfisicalClient{}

	c.client = &http.Client{}

	c.cache, err = cache.NewCache()
	if err != nil {
		return nil, err
	}

	c.InitServiceToken()

	err = c.InitService()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *InfisicalClient) InitServiceToken() {
	c.serviceToken = serviceToken
	token := os.Getenv(consts.InfisicalServiceToken)

	if len(token) > 0 {
		c.serviceToken = token
	}

	if len(c.serviceToken) > 0 {
		c.serviceTokenSecret = c.serviceToken[strings.LastIndex(c.serviceToken, ".")+1:]
	}
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

	projectKey, err := Decrypt(
		c.encryptedKey,
		c.iv,
		c.tag,
		c.serviceTokenSecret,
	)
	if err != nil {
		return err
	}

	c.projectKey = projectKey
	c.secretType = PersonalSecret

	return nil
}
