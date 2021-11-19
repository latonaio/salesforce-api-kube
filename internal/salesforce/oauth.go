package salesforce

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/latonaio/salesforce-api-kube/pkg"
)

type OAuthClientIF interface {
	GetOAuthInfo() (*OAuthInfo, error)
}

// oauthClient struct Http client to get OAuth info
type OAuthClient struct {
	client   *http.Client
	endpoint string
	info     loginInfo
}

// NewOAuthClient return http client to get OAuth info
func NewOAuthClient() (*OAuthClient, error) {
	info, err := getInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get login info: %v", err)
	}
	var endpoint = "https://" + info.Host + "/services/oauth2/token"
	return &OAuthClient{
		client:   &http.Client{},
		endpoint: endpoint,
		info:     *info,
	}, nil
}

type OAuthInfo struct {
	AccessToken string `json:"access_token"`
	InstanceUrl string `json:"instance_url"`
}

// GetOAuthInfo return &oauthResp(instanceURL, accessToken(Bearer))
func (c *OAuthClient) GetOAuthInfo() (*OAuthInfo, error) {
	// Build form
	form := url.Values{}
	form.Add("grant_type", "password")
	form.Add("client_id", c.info.ClientId)
	form.Add("client_secret", c.info.ClientSecret)
	form.Add("username", c.info.UserName)
	form.Add("password", c.info.Password)

	// Build request
	resp, err := http.PostForm(c.endpoint, form)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch response: %v", err)
	}
	defer pkg.SafeClose(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTP %s: failed to read response body: %v", resp.Status, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %s: %s", resp.Status, body)
	}

	// Parse response
	var oauthResp OAuthInfo
	if err := json.Unmarshal(body, &oauthResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json to response struct: %v", err)
	}

	return &oauthResp, nil
}
