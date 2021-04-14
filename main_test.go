package main

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"

	"bitbucket.org/latonaio/salesforce-api-kube/internal/salesforce"
)

func NewRequest(t *testing.T, method, url string, body io.Reader, header http.Header) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("failed to construct request: %v", err)
	}
	req.Header = header
	return req
}

func TestBuildRequest(t *testing.T) {

	// setting header
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("Accept-Encoding", "gzip")
	header.Add("Authorization", "Bearer "+mockAccessToken)

	tests := []struct {
		fromMetadata    map[string]interface{}
		expectedRequest *http.Request
	}{
		// Account
		{
			map[string]interface{}{
				"method": "get",
				"object": "AccountRelatedList",
			},
			NewRequest(t,
				"get",
				mockInstanceUrl+"/services/apexrest/AccountRelatedList/doGetAccountRelatedList",
				nil,
				header,
			),
		},
		{
			map[string]interface{}{
				"method": "get",
				"object": "Account",
			},
			NewRequest(t,
				"get",
				mockInstanceUrl+"/services/apexrest/Account/doGetAccount",
				nil,
				header,
			),
		},
		{
			map[string]interface{}{
				"method":     "get",
				"object":     "Account",
				"path_param": "test_id",
			},
			NewRequest(t,
				"get",
				mockInstanceUrl+"/services/apexrest/Account/doGetAccount/test_id",
				nil,
				header,
			),
		},

		// Contract
		{
			map[string]interface{}{
				"method": "get",
				"object": "ContractRelatedList",
			},
			NewRequest(t,
				"get",
				mockInstanceUrl+"/services/apexrest/ContractRelatedList/doGetContractRelatedList",
				nil,
				header,
			),
		},
		{
			map[string]interface{}{
				"method":     "get",
				"object":     "Contract",
				"path_param": "test_id",
			},
			NewRequest(t,
				"get",
				mockInstanceUrl+"/services/apexrest/Contract/doGetContract/test_id",
				nil,
				header,
			),
		},
		{
			map[string]interface{}{
				"method":       "get",
				"object":       "Contract",
				"query_params": map[string]string{"AccountId": "test_id"},
			},
			NewRequest(t,
				"get",
				mockInstanceUrl+"/services/apexrest/Contract/doGetContract?AccountId=test_id",
				nil,
				header,
			),
		},
		{
			map[string]interface{}{
				"method": "post",
				"object": "Contract",
				"body":   `{"AccountId":"xxxxxxx","EndDate":"2020-12-12","Name":"test","StartDate":"2020-11-11"}`,
			},
			NewRequest(t,
				"post",
				mockInstanceUrl+"/services/apexrest/Contract/doPostContract",
				bytes.NewBufferString(`{"AccountId":"xxxxxxx","EndDate":"2020-12-12","Name":"test","StartDate":"2020-11-11"}`),
				header,
			),
		},
		{
			map[string]interface{}{
				"method": "put",
				"object": "Contract",
				"body":   `{"AccountId":"xxxxxxx","EndDate":"2020-12-12","Name":"test","StartDate":"2020-11-11"}`,
			},
			NewRequest(t,
				"put",
				mockInstanceUrl+"/services/apexrest/Contract/doPutContract",
				bytes.NewBufferString(`{"AccountId":"xxxxxxx","EndDate":"2020-12-12","Name":"test","StartDate":"2020-11-11"}`),
				header,
			),
		},
		{
			map[string]interface{}{
				"method": "put",
				"object": "ContractPDF",
				"path_param": "0000",
				"query_params": map[string]string{"pdfName": "aaaa_0000.pdf"},
				"body":   `xxx`,
			},
			NewRequest(t,
				"put",
				mockInstanceUrl+"/services/apexrest/ContractPDF/doPutContractPDF/0000?pdfName=aaaa_0000.pdf",
				bytes.NewBufferString(`xxx`),
				header,
			),
		},
	}
	for i, tt := range tests {
		oauthClient, err := NewMockOAuthClient()
		if err != nil {
			panic(err)
		}
		got, err := salesforce.BuildRequest(tt.fromMetadata, oauthClient)
		if err != nil {
			t.Errorf("failed to build request: %v", err)
		}

		if got.Method != tt.expectedRequest.Method {
			t.Errorf("%d# Method got: %v\n, want: %v\n", i, got.Method, tt.expectedRequest.Method)
		}
		if got.URL.String() != tt.expectedRequest.URL.String() {
			t.Errorf("%d# URL got: %v\n, want: %v\n", i, got.URL.String(), tt.expectedRequest.URL.String())
		}
		if !reflect.DeepEqual(got.Header, tt.expectedRequest.Header) {
			t.Errorf("%d# Header got: %v\n want: %v\n", i, got.Header, tt.expectedRequest.Header)
		}
		if !reflect.DeepEqual(got.Body, tt.expectedRequest.Body) {
			t.Errorf("%d# body got: %v\n want: %v\n", i, got.Body, tt.expectedRequest.Body)
		}
	}
}
