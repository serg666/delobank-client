package deloBankClient

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"
)

type deloBankRB struct {
	client *http.Client
	logger logrus.FieldLogger
	req    *http.Request
	body   []byte
	Error  error
}

type requestHeaders struct {
	ClientId     string
	ClientSecret string
}

func (rb *deloBankRB) PostJson(
	url string,
	jsonBody map[string]interface{},
	headers requestHeaders,
) *deloBankRB {
	return rb.runIfNoError(func() error {
		jsonValue, err := json.Marshal(jsonBody)
		if err != nil {
			rb.logger.Error("DeloBank. Body is not JSON")
			return err
		}

		rb.logger.Printf("DeloBank. PostJson: url=%s, body=%s", url, rb.maskBody(string(jsonValue)))

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
		if err != nil {
			rb.logger.Error("DeloBank. Request not created")
			return err
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("clientID", headers.ClientId)
		req.Header.Add("clientSecret", headers.ClientSecret)
		rb.req = req

		return nil
	})
}

func (rb *deloBankRB) Get(url string, headers requestHeaders) *deloBankRB {
	return rb.runIfNoError(func() error {
		rb.logger.Printf("DeloBank. Get: url=%s", url)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			rb.logger.Error("DeloBank. Request not created")
			return err
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("clientID", headers.ClientId)
		req.Header.Add("clientSecret", headers.ClientSecret)
		rb.req = req

		rb.req = req
		return nil
	})
}

func (rb *deloBankRB) Send() *deloBankRB {
	return rb.runIfNoError(func() error {
		if rb.req == nil {
			return errors.New("No request")
		}

		res, err := rb.client.Do(rb.req)
		if err != nil {
			rb.logger.Error("DeloBank. Request send failed")
			return err
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			rb.logger.Error("DeloBank. Read response failed")
			return err
		}

		rb.logger.Printf("DeloBank. Response: statusCode=%d, body=%s", res.StatusCode, string(body))

		if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted {
			return errors.New("Bad status code")
		}

		rb.body = body

		return nil
	})
}

func (rb *deloBankRB) UnmarshalJsonResponse(v interface{}) error {
	rb.runIfNoError(func() error {
		if rb.body == nil {
			return errors.New("Response body not set")
		}

		return json.Unmarshal(rb.body, v)
	})

	return rb.Error
}

func (rb *deloBankRB) runIfNoError(fn func() error) *deloBankRB {
	if rb.Error == nil {
		rb.Error = fn()
	}

	return rb
}

func (rb deloBankRB) maskBody(data string) string {
	result := regexp.MustCompile(`"pan":"[^"]+([^"]{4})"`).
		ReplaceAllString(data, `"pan":"******$1"`)
	return regexp.MustCompile(`"csc":"([^"]+)"`).
		ReplaceAllString(result, `"csc":"***"`)
}

func NewRequestBuilder(client *http.Client, logger logrus.FieldLogger) *deloBankRB {
	return &deloBankRB{
		client: client,
		logger: logger,
	}
}
