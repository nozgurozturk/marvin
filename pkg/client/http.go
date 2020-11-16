
package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type HTTPClient interface {
	Get(endpoint string, header map[string]string) ([]byte, error) // GET request
	Post(endpoint string, header map[string]string, body map[string]interface{}) ([]byte, error) // POST request
}

type Http struct {
	baseUrl string
}

// Creates custom http client for http request
func New(baseUrl string) HTTPClient {
	return &Http{
		baseUrl: baseUrl,
	}
}

func (c *Http) Get(endpoint string, header map[string]string) ([]byte, error) {

	timeout := 10 * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", c.baseUrl+endpoint, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range header {
		request.Header.Add(key, value)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	data, _ := ioutil.ReadAll(response.Body)

	return data, nil
}

// POST request
func (c *Http) Post(endpoint string, header map[string]string, body map[string]interface{}) ([]byte, error) {
	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	timeout := 10 * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("POST", c.baseUrl+endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	for key, value := range header {
		request.Header.Add(key, value)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	data, _ := ioutil.ReadAll(response.Body)

	return data, nil
}
