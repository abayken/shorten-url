package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/abayken/shorten-url/internal/app/storage"
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

	err = handler.Storage.Save(shortURLID, model.URL, userID)

	if err != nil {
		var duplicateError *storage.DuplicateURLError

		if errors.As(err, &duplicateError) {
			responseModel := PostAPIURLResponse{Result: handler.BaseURL + "/" + duplicateError.ShortURLID}

			jsonResponse, err := json.Marshal(responseModel)

			if err != nil {
				ctx.Status(http.StatusInternalServerError)

				return
			}

			ctx.String(http.StatusConflict, string(jsonResponse))

			return
		}
	}

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

		ctx.Data(http.StatusOK, "application/json", jsonResponse)
	} else {
		ctx.Status(http.StatusNoContent)
	}
}

func (handler *URLHandler) BatchURLS(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)

	if err != nil {
		ctx.Status(http.StatusBadRequest)

		return
	}

	userID := ctx.GetString("token")

	/// Формат тело запроса
	type URLBatchRequest struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	/// Формат ответа
	type URLBatchResponse struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}

	var urls []URLBatchRequest

	err = json.Unmarshal(body, &urls)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}

	/// массив урлов которые добавятся в базу
	var batchURLs []storage.BatchURL

	/// массив сокращенных урлов для ответа
	var batchURLsResponse []URLBatchResponse

	for _, item := range urls {
		shortURLID := handler.URLShortener.ID()

		batchURLs = append(batchURLs,
			storage.BatchURL{
				UserID:     userID,
				ShortURLID: shortURLID,
				FullURL:    item.OriginalURL,
			})

		batchURLsResponse = append(
			batchURLsResponse,
			URLBatchResponse{
				CorrelationID: item.CorrelationID,
				ShortURL:      handler.BaseURL + "/" + shortURLID,
			})
	}

	err = handler.Storage.BatchURLs(batchURLs)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}

	jsonResponse, err := json.Marshal(batchURLsResponse)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}

	ctx.Data(http.StatusCreated, "application/json", jsonResponse)
}
