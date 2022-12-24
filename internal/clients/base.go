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

func NewBaseClient(baseURL string, timeout time.Duration) *BaseClient {
	if !strings.HasPrefix(baseURL, "http") {
		baseURL = "http://" + baseURL
	}
	return &BaseClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: timeout,
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
		return Response{}, fmt.Errorf("error: %s details: %s", resp.Status, body)
	}

	return Response{
		Body:       body,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}, nil
}
