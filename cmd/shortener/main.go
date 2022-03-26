package main

import (
	"context"
	"flag"
	"log"

	"github.com/abayken/shorten-url/internal/app"
	"github.com/abayken/shorten-url/internal/app/handlers"
	"github.com/abayken/shorten-url/internal/app/storage"
	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"urls.json"`
	DatabaseURL     string `env:"DATABASE_DSN"`
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
	flag.StringVar(&cfg.DatabaseURL, "d", cfg.DatabaseURL, "Урл базы данных")

	flag.Parse()

	var storage = getDatabaseStorage(cfg.DatabaseURL)

	router := GetRouter(storage, app.RealURLShortener{}, cfg)
	router.Run(cfg.ServerAddress)
}

func getDatabaseStorage(url string) storage.DatabaseStorage {
	conn, err := pgx.Connect(context.Background(), url)

	if err != nil {
		log.Fatal(err)
	}

	var storage = storage.DatabaseStorage{DB: conn}
	storage.InitTablesIfNeeded()

	return storage
}

func GetRouter(storage storage.URLStorage, urlShortener app.URLShortener, cfg Config) *gin.Engine {
	handler := handlers.URLHandler{Storage: storage, URLShortener: urlShortener, BaseURL: cfg.BaseURL}
	router := gin.New()
	router.Use(gzip.Gzip(gzip.BestSpeed, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	router.Use(Tokenize())
	router.GET("/:id", handler.GetFullURL)
	router.POST("/", handler.PostFullURL)
	router.POST("/api/shorten", handler.PostAPIFullURL)
	router.POST("/api/shorten/batch", handler.BatchURLS)
	router.GET("/api/user/urls", handler.GetUserURLs)

	health := handlers.Health{DatabaseURL: cfg.DatabaseURL}
	router.GET("/ping", health.CheckDatabase)

	return router
}
