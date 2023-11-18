package clients

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseClient_MakeURL(t *testing.T) {
	baseClient1 := BaseClient{baseURL: "localhost", client: &http.Client{}}
	baseClient2 := BaseClient{baseURL: "localhost/", client: &http.Client{}}

	assert.Equal(t, "localhost/api/v1", baseClient1.MakeURL("api/v1"))
	assert.Equal(t, "localhost/api/v1", baseClient2.MakeURL("api/v1"))
	assert.Equal(t, "localhost/api/v1", baseClient1.MakeURL("///api/v1"))
	assert.Equal(t, "localhost/api/v1", baseClient2.MakeURL("///api/v1"))
	assert.Equal(t, "localhost/api/v2/", baseClient1.MakeURL("api/v2/"))
	assert.Equal(t, "localhost/api/v2/", baseClient2.MakeURL("api/v2/"))
}

func TestBaseClient_DoRequest(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.URL.Path, "/endpoint/")
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		defer r.Body.Close()

		body, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		assert.Equal(t, `{"msg": "ping"}`, string(body))
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("pong"))
	}))
	defer svr.Close()
	baseClient := BaseClient{baseURL: svr.URL, client: &http.Client{}}
	request := Request{
		Method:       http.MethodPost,
		URL:          baseClient.MakeURL("endpoint/"),
		Headers:      map[string]string{"Content-Type": "application/json"},
		Body:         []byte(`{"msg": "ping"}`),
		OkStatusCode: http.StatusCreated,
	}
	resp, err := baseClient.DoRequest(&request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "text/plain", resp.Headers.Get("Content-Type"))
	assert.Equal(t, "pong", string(resp.Body))
}

func TestBaseClient_DoRequest_Failed(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg": "Something went wrong"}`))
	}))
	defer svr.Close()
	baseClient := BaseClient{baseURL: svr.URL, client: &http.Client{}}
	request := Request{
		Method:       http.MethodGet,
		URL:          baseClient.MakeURL("endpoint/"),
		Headers:      map[string]string{"Content-Type": "application/json"},
		Body:         []byte(`{"msg": "ping"}`),
		OkStatusCode: http.StatusCreated,
	}
	resp, err := baseClient.DoRequest(&request)
	assert.Equal(t, Response{}, resp)
	assert.NotNil(t, err)
	assert.Equal(t, "error: 400 Bad Request details: {\"msg\": \"Something went wrong\"}", err.Error())
}
