package router

import (
	"github.com/abayken/shorten-url/internal/app"
	"github.com/abayken/shorten-url/internal/app/handlers"
	"github.com/abayken/shorten-url/internal/app/storage"
	"github.com/gin-gonic/gin"
)

func GetRouter(storage storage.URLStorage, urlShortener app.URLShortener) *gin.Engine {
	handler := handlers.URLHandler{Storage: storage, URLShortener: urlShortener}
	router := gin.New()
	router.GET("/:id", handler.GetFullURL)
	router.POST("/", handler.PostFullURL)

	return router
}
