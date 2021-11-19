package main

import (
	"github.com/latonaio/salesforce-api-kube/internal/salesforce"
)

const (
	mockAccessToken = "mock_access_token"
	mockInstanceUrl = "https://mock_instance"
)

type MockOAuthClient struct {
}

func NewMockOAuthClient() (*MockOAuthClient, error) {
	return &MockOAuthClient{}, nil
}
func (c *MockOAuthClient) GetOAuthInfo() (*salesforce.OAuthInfo, error) {
	return &salesforce.OAuthInfo{
		AccessToken: mockAccessToken,
		InstanceUrl: mockInstanceUrl,
	}, nil
}
