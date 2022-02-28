package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/abayken/shorten-url/internal/app/handlers"
	"github.com/abayken/shorten-url/internal/app/storage"
	"github.com/caarlos0/env/v6"
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

/// Тест который проверяет сокращения урла через POST запрос
func TestURLSave(testing *testing.T) {
	result := request(http.MethodPost, "/", strings.NewReader(fullURL))

	assert.Equal(testing, 201, result.StatusCode)

	bodyResult, err := ioutil.ReadAll(result.Body)
	require.NoError(testing, err)
	err = result.Body.Close()
	require.NoError(testing, err)

	shortURL := string(bodyResult[:])
	assert.Equal(testing, baseURL+fakeID, shortURL)
}

func TestURLGet(testing *testing.T) {
	result := request(http.MethodPost, "/", strings.NewReader(fullURL))

	bodyResult, err := ioutil.ReadAll(result.Body)
	require.NoError(testing, err)
	err = result.Body.Close()
	require.NoError(testing, err)

	shortURL := string(bodyResult[:])
	fmt.Println(shortURL)
	/// Делаем GET запрос и проверяем результат
	getMethodResult := request(http.MethodGet, shortURL, nil)
	getMethodResult.Body.Close()

	assert.Equal(testing, fullURL, getMethodResult.Header.Get("Location"))
}

/// Тест на метод /api/shorten
func TestURLApiPost(testing *testing.T) {
	requestModel := handlers.PostAPIURLRequest{URL: fullURL}
	requestBody, _ := json.Marshal(requestModel)

	result := request(http.MethodPost, "/api/shorten", bytes.NewReader(requestBody))

	defer result.Body.Close()

	// проверка статус кода
	assert.Equal(testing, 201, result.StatusCode)

	/// проверка сокращенного урла
	bodyResult, _ := ioutil.ReadAll(result.Body)

	var responseModel handlers.PostAPIURLResponse

	_ = json.Unmarshal(bodyResult, &responseModel)
	assert.Equal(testing, baseURL+fakeID, responseModel.Result)
}

/// стучится в некий ендпойнт
func request(method string, url string, body io.Reader) *http.Response {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	router := GetRouter(storage.NewMapURLStorage(make(map[string]string)), FakeURLShortener{}, cfg)

	request := httptest.NewRequest(method, url, body)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	result := recorder.Result()

	return result
}
