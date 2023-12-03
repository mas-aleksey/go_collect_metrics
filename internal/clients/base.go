// Package clients - функционал для отправки API запросов.
package clients

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/tiraill/go_collect_metrics/internal/utils"
)

// Request - структура описывает API запрос.
type Request struct {
	Method       string            // метод запроса
	URL          string            // URL запроса
	Headers      map[string]string // заголовки запроса
	Body         []byte            // тело запроса
	OkStatusCode int               // ожидаемый код ответа
}

// Response - структура описывает API ответ.
type Response struct {
	Body       []byte      // тело ответа
	StatusCode int         // статус код ответа
	Headers    http.Header // заголовки ответа
}

// BaseClient - структура описывает базового клиента.
type BaseClient struct {
	baseURL   string
	client    *http.Client
	rateLimit int
	publicKey *utils.PublicKey
}

// NewBaseClient - метод для создания базового клиента.
func NewBaseClient(baseURL string, timeout time.Duration, rateLimit int, publicKeyPath string) *BaseClient {
	if !strings.HasPrefix(baseURL, "http") {
		baseURL = "http://" + baseURL
	}
	publicKey, err := utils.LoadPublicKey(publicKeyPath)
	if err != nil {
		log.Fatal("Failed to load public key:", err)
	}
	return &BaseClient{
		baseURL:   baseURL,
		client:    &http.Client{Timeout: timeout},
		rateLimit: rateLimit,
		publicKey: publicKey,
	}
}

// MakeURL - метод формирует url для запроса.
func (c *BaseClient) MakeURL(url string) string {
	baseURL := strings.TrimRight(c.baseURL, "/")
	path := strings.TrimLeft(url, "/")
	return fmt.Sprint(baseURL, "/", path)
}

// DoRequest - метод выполняет API запрос.
func (c *BaseClient) DoRequest(r *Request) (Response, error) {
	client := &http.Client{}
	var requestBody bytes.Buffer

	if c.publicKey != nil {
		encryptedBody, err := c.publicKey.Encrypt(r.Body)
		if err != nil {
			log.Println("Failed to encrypt body:", err)
			return Response{}, err
		}
		r.Body = encryptedBody
	}

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
