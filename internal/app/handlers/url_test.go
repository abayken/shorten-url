package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/abayken/shorten-url/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/// Тест который проверяет сокращения урла через POST запрос
func TestURLSave(testing *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://hello.com/23213213123"))
	recorder := httptest.NewRecorder()
	customURLHandler := URLHandler{Storage: storage.MapURLStorage{}}
	h := http.HandlerFunc(customURLHandler.ServerHTTP)
	h.ServeHTTP(recorder, request)
	result := recorder.Result()

	assert.Equal(testing, 201, result.StatusCode)

	bodyResult, err := ioutil.ReadAll(result.Body)
	require.NoError(testing, err)
	err = result.Body.Close()
	require.NoError(testing, err)

	shortURL := string(bodyResult[:])
	assert.True(testing, shortURL != "")
}

func TestURLGet(testing *testing.T) {
	fullURL := "https://hello.com/23213213123"
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fullURL))
	recorder := httptest.NewRecorder()
	customURLHandler := URLHandler{Storage: storage.MapURLStorage{}}
	h := http.HandlerFunc(customURLHandler.ServerHTTP)
	h.ServeHTTP(recorder, request)
	result := recorder.Result()

	bodyResult, err := ioutil.ReadAll(result.Body)
	require.NoError(testing, err)
	err = result.Body.Close()
	require.NoError(testing, err)

	shortURL := string(bodyResult[:])

	request = httptest.NewRequest(http.MethodGet, shortURL, nil)
	getMethodRecorder := httptest.NewRecorder()
	h.ServeHTTP(getMethodRecorder, request)
	getMethodResult := getMethodRecorder.Result()
	assert.Equal(testing, fullURL, getMethodResult.Header.Get("Location"))
}
