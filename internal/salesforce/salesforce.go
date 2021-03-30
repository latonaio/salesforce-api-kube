package salesforce

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"bitbucket.org/latonaio/aion-core/pkg/log"
	"bitbucket.org/latonaio/salesforce-api-kube/internal/str"
	"bitbucket.org/latonaio/salesforce-api-kube/pkg"
)

// client struct
type client struct {
	client          *http.Client
	requiredHeaders http.Header
}

func NewClient() *client {
	requiredHeaders := http.Header{}
	requiredHeaders.Add("Content-Type", "application/json")
	requiredHeaders.Add("Accept-Encoding", "gzip")
	return &client{
		client:          &http.Client{},
		requiredHeaders: requiredHeaders,
	}
}

func (c *client) buildRequest(method, endpoint, accessToken string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch response: %v", err)
	}
	req.Header = c.requiredHeaders
	req.Header.Add("Authorization", "Bearer "+accessToken)
	return req, nil
}

// Do does http request
func (c *client) Do(r *http.Request) (string, error) {
	resp, err := c.client.Do(r)
	if err != nil {
		return "", fmt.Errorf("failed to fetch response: %v", err)
	}
	defer pkg.SafeClose(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("HTTP %s: failed to read response body: %v", resp.Status, err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %s: %s", resp.Status, body)
	}
	return string(body), err
}

// DoRequest does http request and return http response body.
func DoRequest(metadata map[string]interface{}, oauthClient *OAuthClient) (string, error) {
	sfclient := NewClient()

	// Parse metadata json(Check nil and convert)
	objectIF, ok := metadata["object"]
	if !ok {
		return "", errors.New("invalid metadata: object not found")
	}
	object, ok := objectIF.(string)
	if !ok {
		return "", errors.New("failed to convert interface{} to string")
	}
	methodIF, ok := metadata["method"]
	if !ok {
		return "", errors.New("invalid metadata: method not found")
	}
	method, ok := methodIF.(string)
	if !ok {
		return "", errors.New("failed to convert interface{} to string")
	}
	method = strings.ToLower(method) // salesforce api only accepts lowercase methods.

	// Get InstanceURL and AccessToken(Bearer)
	info, err := oauthClient.GetOAuthInfo()
	if err != nil {
		return "", fmt.Errorf("failed to salesforce authorization: %v", err)
	}

	// Build URL
	const servicesUrl = "/services/apexrest"
	u, err := url.Parse(info.InstanceUrl + servicesUrl)
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %v", err)
	}

	u.Path = path.Join(u.Path, object, "do"+str.ToFirstUppercase(method)+object)
	if accountIDIF, exist := metadata["account_id"]; exist {
		log.Printf("account_id is exist: %v\n", accountIDIF)
		accountID, ok := accountIDIF.(string)
		if !ok {
			return "", errors.New("failed to convert account_id to string")
		}
		u.Path = path.Join(u.Path, accountID)
	}
	if pathParamIF, exist := metadata["path_param"]; exist {
		log.Printf("path_param is exist: %v\n", pathParamIF)
		pathParam, ok := pathParamIF.(string)
		if !ok {
			return "", errors.New("failed to convert path_param to string")
		}
		u.Path = path.Join(u.Path, pathParam)
	}
	if queryParamsIF, exist := metadata["query_params"]; exist {
		log.Printf("query_params is exist: %v\n", queryParamsIF)
		queryParams, ok := queryParamsIF.(map[string]string)
		if !ok {
			return "", errors.New("failed to convert query_params to map[string]string")
		}
		for k, v := range queryParams {
			u.Query().Set(k, v)
		}
		u.RawQuery = u.Query().Encode()
	}

	req, err := sfclient.buildRequest(method, u.String(), info.AccessToken, nil)
	if err != nil {
		return "", fmt.Errorf("failed to build request: %v", err)
	}
	log.Printf("send request: %v", req)

	respBody, err := sfclient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch response: %v", err)
	}
	return respBody, nil
}
