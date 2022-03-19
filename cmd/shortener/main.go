package main

import (
	"flag"
	"log"

	"github.com/abayken/shorten-url/internal/app"
	"github.com/abayken/shorten-url/internal/app/handlers"
	"github.com/abayken/shorten-url/internal/app/storage"
	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
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

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "Адресс сервера")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "BaseURL сокращенного урла")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "Путь до файла где хранятся урлы")

	flag.Parse()

	router := GetRouter(storage.FileURLStorage{Path: cfg.FileStoragePath}, app.RealURLShortener{}, cfg)
	router.Run(cfg.ServerAddress)
}

func GetRouter(storage storage.URLStorage, urlShortener app.URLShortener, cfg Config) *gin.Engine {
	handler := handlers.URLHandler{Storage: storage, URLShortener: urlShortener, BaseURL: cfg.BaseURL}
	router := gin.New()
	router.Use(Compress())
	router.Use(Unpack())
	router.Use(Tokenize())
	router.GET("/:id", handler.GetFullURL)
	router.POST("/", handler.PostFullURL)
	router.POST("/api/shorten", handler.PostAPIFullURL)

	return router
}
