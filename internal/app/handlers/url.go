package handlers

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/abayken/shorten-url/internal/app"
	"github.com/abayken/shorten-url/internal/app/storage"
	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	Storage      storage.URLStorage
	URLShortener app.URLShortener
	BaseURL      string
}

func (handler *URLHandler) GetFullURL(ctx *gin.Context) {
	shortURLID := ctx.Param("id")

	fullURL := handler.Storage.Get(shortURLID)

	if fullURL != "" {
		ctx.Header("Location", fullURL)
		ctx.Status(http.StatusTemporaryRedirect)
	} else {
		ctx.Status(http.StatusBadRequest)
	}
}

func (handler *URLHandler) PostFullURL(ctx *gin.Context) {
	fullURLByte, err := ioutil.ReadAll(ctx.Request.Body)

	if err != nil {
		ctx.Status(http.StatusBadRequest)

		return
	}

	url := string(fullURLByte)

	defer ctx.Request.Body.Close()

	shortURLID := handler.URLShortener.ID()

	userID := ctx.GetString("token")

	err = handler.Storage.Save(shortURLID, url, userID)

	if err != nil {
		var duplicateError *storage.DuplicateURLError

		if errors.As(err, &duplicateError) {
			ctx.String(http.StatusConflict, handler.BaseURL+"/"+duplicateError.ShortURLID)

			return
		}
	}

	ctx.String(http.StatusCreated, handler.BaseURL+"/"+shortURLID)
}
