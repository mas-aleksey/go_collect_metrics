package clients

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Method       string
	URL          string
	Headers      map[string]string
	Body         []byte
	OkStatusCode int
}

type Response struct {
	Body       []byte
	StatusCode int
	Headers    http.Header
}

type BaseClient struct {
	baseURL string
	client  *http.Client
}

type ClientConfig struct {
	BaseURL string
	Timeout time.Duration
}

func NewClientConfig(baseURL string, timeout time.Duration) ClientConfig {
	return ClientConfig{
		BaseURL: baseURL,
		Timeout: timeout,
	}
}

func NewBaseClient(config ClientConfig) *BaseClient {
	return &BaseClient{
		baseURL: config.BaseURL,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

func (c *BaseClient) MakeURL(URL string) string {
	baseURL := strings.TrimRight(c.baseURL, "/")
	path := strings.TrimLeft(URL, "/")
	return fmt.Sprint(baseURL, "/", path)
}

func (c *BaseClient) DoRequest(r Request) (Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(r.Method, r.URL, bytes.NewBuffer(r.Body))
	if err != nil {
		return Response{}, err
	}

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	if resp.StatusCode != r.OkStatusCode {
		return Response{}, fmt.Errorf("Error: %s details: %s\n", resp.Status, body)
	}

	return Response{
		Body:       body,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}, nil
}
