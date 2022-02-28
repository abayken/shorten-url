package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/abayken/shorten-url/internal/app/handlers"
	"github.com/abayken/shorten-url/internal/app/storage"
	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FakeURLShortener struct {
}

func (shortener FakeURLShortener) ID() string {
	return fakeID
}

const (
	fullURL = "https://hello.com/23213213123"
	fakeID  = "12345"
	baseURL = "http://localhost:8080/"
)

func TestEndpoints(t *testing.T) {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	router := GetRouter(storage.NewMapURLStorage(make(map[string]string)), FakeURLShortener{}, cfg)

	type want struct {
		status         int
		locationHeader string
		bodyChecker    func(body []byte) bool
	}

	type test struct {
		name   string
		url    string
		method string
		body   string
		want   want
	}

	var tests = []test{
		{
			name:   "Test POST to /",
			url:    "/",
			method: http.MethodPost,
			body:   fullURL,
			want: want{
				status: http.StatusCreated,
				bodyChecker: func(body []byte) bool {
					response := string(body)
					return response == baseURL+fakeID
				},
			},
		},
		{
			name:   "Test GET short url",
			url:    baseURL + fakeID,
			method: http.MethodGet,
			want: want{
				status:         http.StatusTemporaryRedirect,
				locationHeader: fullURL,
			},
		},
		{
			name:   "Test /api/shorten",
			url:    "/api/shorten",
			method: http.MethodPost,
			body:   fmt.Sprintf(`{"result":"%s"}`, fullURL),
			want: want{
				status: http.StatusCreated,
				bodyChecker: func(body []byte) bool {
					var response handlers.PostAPIURLResponse

					json.Unmarshal(body, &response)
					return response.Result == baseURL+fakeID
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := request(tt.method, tt.url, tt.body, router)
			bodyResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.status, result.StatusCode)
			assert.Equal(t, tt.want.locationHeader, result.Header.Get("Location"))

			if tt.want.bodyChecker != nil {
				assert.Equal(t, tt.want.bodyChecker(bodyResult), true)
			}

			defer result.Body.Close()
		})
	}
}

/// стучится в некий ендпойнт
func request(method string, url string, body string, router *gin.Engine) *http.Response {
	request := httptest.NewRequest(method, url, strings.NewReader(body))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	result := recorder.Result()
	return result
}
