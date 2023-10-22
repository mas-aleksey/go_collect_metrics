package clients

import (
	"bytes"
	"compress/gzip"
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
	baseURL   string
	client    *http.Client
	rateLimit int
}

func NewBaseClient(baseURL string, timeout time.Duration, rateLimit int) *BaseClient {
	if !strings.HasPrefix(baseURL, "http") {
		baseURL = "http://" + baseURL
	}
	return &BaseClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: timeout,
		},
		rateLimit: rateLimit,
	}
}

func (c *BaseClient) MakeURL(URL string) string {
	baseURL := strings.TrimRight(c.baseURL, "/")
	path := strings.TrimLeft(URL, "/")
	return fmt.Sprint(baseURL, "/", path)
}

func (c *BaseClient) DoRequest(r *Request) (Response, error) {
	client := &http.Client{}
	var requestBody bytes.Buffer

	_, ok := r.Headers["Content-Encoding"]
	if ok {
		gz := gzip.NewWriter(&requestBody)
		if _, err := gz.Write(r.Body); err != nil {
			return Response{}, err
		}
		if err := gz.Close(); err != nil {
			return Response{}, err
		}
	} else {
		requestBody = *bytes.NewBuffer(r.Body)
	}

	req, err := http.NewRequest(r.Method, r.URL, &requestBody)
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
