package router

import (
	"github.com/abayken/shorten-url/internal/app/handlers"
	"github.com/abayken/shorten-url/internal/app/storage"
	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	handler := handlers.URLHandler{Storage: storage.MapURLStorage{}}
	router := gin.New()
	router.GET("/:id", handler.GetFullURL)
	router.POST("/", handler.PostFullURL)

	return router
}
