package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/abayken/shorten-url/internal/app"
	"github.com/abayken/shorten-url/internal/app/storage"
	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	Storage storage.URLStorage
}

func (handler *URLHandler) GetFullURL(context *gin.Context) {
	shortURLID := context.Param("id")

	fullURL := handler.Storage.Get(shortURLID)

	if fullURL != "" {
		context.Header("Location", fullURL)
		context.Status(http.StatusTemporaryRedirect)
	} else {
		context.Status(http.StatusBadRequest)
	}
}

func (handler *URLHandler) PostFullURL(context *gin.Context) {
	fullURLByte, err := ioutil.ReadAll(context.Request.Body)

	if err != nil {
		context.Status(http.StatusBadRequest)

		return
	}

	url := string(fullURLByte)

	defer context.Request.Body.Close()

	urlShortener := app.URLShortener{URL: url}

	shortURLID := urlShortener.ID()

	handler.Storage.Save(shortURLID, url)

	context.String(http.StatusCreated, "http://localhost:8080/"+shortURLID)
}
