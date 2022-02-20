package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

/// Здесь лежат хэндлеры для апишки

/// Формат body для метода /api/shorten от клиента
type PostApiURLRequest struct {
	URL string `json:"url"`
}

/// Формат ответа для метода /api/shorten
type PostApiURLResponse struct {
	Result string `json:"result"`
}

const (
	BASE_URL = "http://localhost:8080/"
)

/// Метод который возвращает сокращенный URL
/// Отвечает в виде JSON
func (handler *URLHandler) PostApiFullURL(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)

	if err != nil {
		ctx.Status(http.StatusBadRequest)

		return
	}

	var model PostApiURLRequest

	err = json.Unmarshal(body, &model)

	if err != nil {
		ctx.Status(http.StatusBadRequest)

		return
	}

	defer ctx.Request.Body.Close()

	shortURLID := handler.URLShortener.ID()
	handler.Storage.Save(shortURLID, model.URL)

	responseModel := PostApiURLResponse{Result: BASE_URL + shortURLID}

	jsonResponse, err := json.Marshal(responseModel)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}

	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.String(http.StatusOK, string(jsonResponse))
}
