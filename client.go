package deloBankClient

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

type DELOBankClient struct {
	BaseUrl   string
	PayoutApi string
	StatusApi string
	Client    *http.Client
}

type transferResponse struct {
	RequestId string `json:"requestID"`
}

type getStateResponse struct {
	Status    string
	RequestId string `json:"requestID"`
	Results   []struct {
		TransferId   string `json:"transferID"`
		Result       string `json:"result"`
		Status       string `json:"status"`
		IsFatal      bool   `json:"isFatal"`
		CardAuthInfo struct {
			AuthCode *string `json:"authCode"`
			RefNum   *string `json:"refNum"`
		} `json:"cardAuthInfo"`
	} `json:"results"`
}

func (deloBank DELOBankClient) Payout(
	logger logrus.FieldLogger,
	ClientId string,
	ClientSecret string,
	jsonBody map[string]interface{},
) (*transferResponse, error) {
	url := fmt.Sprintf("%s/%s", deloBank.BaseUrl, deloBank.PayoutApi)
	headers := requestHeaders{ClientId: ClientId, ClientSecret: ClientSecret}
	result := &transferResponse{}
	err := NewRequestBuilder(deloBank.Client, logger).
		PostJson(url, jsonBody, headers).
		Send().
		UnmarshalJsonResponse(result)

	if err != nil {
		logger.Errorf("DeloBank. InitPayout: HTTP-request failed: %v", err)
		return nil, err
	}

	return result, nil
}

func (deloBank DELOBankClient) GetState(
	logger logrus.FieldLogger,
	requestId string,
	ClientId string,
	ClientSecret string,
) (*getStateResponse, error) {
	result := &getStateResponse{}
	url := fmt.Sprintf("%s/%s/%s", deloBank.BaseUrl, deloBank.StatusApi, requestId)
	headers := requestHeaders{ClientId: ClientId, ClientSecret: ClientSecret}
	err := NewRequestBuilder(deloBank.Client, logger).
		Get(url, headers).
		Send().
		UnmarshalJsonResponse(result)
	if err != nil {
		logger.Errorf("DeloBank. GetState: HTTP-request failed")
		return nil, err
	}
	return result, nil
}
