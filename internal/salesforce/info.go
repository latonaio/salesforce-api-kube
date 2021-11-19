package salesforce

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// loginInfo struct (Login Info to Salesforce)
type loginInfo struct {
	Host         string `json:"host"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	UserName     string `json:"username"`
	Password     string `json:"password"`
}

func getInfo() (*loginInfo, error) {
	filePath := "./config.json"
	if os.Getenv("DEV") == "true" || os.Getenv("DEV") == "True" {
		filePath = "./config.test.json"
	}
	j, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
	var i loginInfo
	if err = json.Unmarshal(j, &i); err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	return &i, nil
}
