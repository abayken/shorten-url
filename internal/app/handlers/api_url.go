package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

/// Здесь лежат хэндлеры для апишки

/// Формат body для метода /api/shorten от клиента
type PostAPIURLRequest struct {
	URL string `json:"url"`
}

/// Формат ответа для метода /api/shorten
type PostAPIURLResponse struct {
	Result string `json:"result"`
}

/// Метод который возвращает сокращенный URL
/// Отвечает в виде JSON
func (handler *URLHandler) PostAPIFullURL(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)

	if err != nil {
		ctx.Status(http.StatusBadRequest)

		return
	}

	var model PostAPIURLRequest

	err = json.Unmarshal(body, &model)

	if err != nil {
		ctx.Status(http.StatusBadRequest)

		return
	}

	defer ctx.Request.Body.Close()

	shortURLID := handler.URLShortener.ID()
	userID := ctx.GetString("token")

	handler.Storage.Save(shortURLID, model.URL, userID)

	responseModel := PostAPIURLResponse{Result: handler.BaseURL + "/" + shortURLID}

	jsonResponse, err := json.Marshal(responseModel)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}

	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.String(http.StatusCreated, string(jsonResponse))
}

func (handler *URLHandler) GetUserURLs(ctx *gin.Context) {
	userID := ctx.GetString("token")

	urls := handler.Storage.FetchUserURLs(userID)

	if len(urls) > 0 {
		for index, url := range urls {
			urls[index] = url.BaseURLAppended(handler.BaseURL)
		}

		jsonResponse, err := json.Marshal(urls)

		if err != nil {
			ctx.Status(http.StatusInternalServerError)

			return
		}

		ctx.String(http.StatusOK, string(jsonResponse))
	} else {
		ctx.Status(http.StatusNoContent)
	}
}
