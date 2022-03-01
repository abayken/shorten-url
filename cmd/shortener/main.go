package main

import (
	"log"

	"github.com/abayken/shorten-url/internal/app"
	"github.com/abayken/shorten-url/internal/app/handlers"
	"github.com/abayken/shorten-url/internal/app/storage"
	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"urls.json"`
}

func main() {
	/// получаем переменные окружения
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	var storageStrategy storage.URLStorage

	if cfg.FileStoragePath == "" {
		storageStrategy = storage.NewMapURLStorage(make(map[string]string))
	} else {
		storageStrategy = storage.FileURLStorage{Path: cfg.FileStoragePath}
	}

	router := GetRouter(storageStrategy, app.RealURLShortener{}, cfg)
	router.Run(cfg.ServerAddress)
}

func GetRouter(storage storage.URLStorage, urlShortener app.URLShortener, cfg Config) *gin.Engine {
	handler := handlers.URLHandler{Storage: storage, URLShortener: urlShortener, BaseURL: cfg.BaseURL}
	router := gin.New()
	router.GET("/:id", handler.GetFullURL)
	router.POST("/", handler.PostFullURL)
	router.POST("/api/shorten", handler.PostAPIFullURL)

	return router
}
