package main

import (
	"os"

	"github.com/abayken/shorten-url/internal/app"
	"github.com/abayken/shorten-url/internal/app/router"
	"github.com/abayken/shorten-url/internal/app/storage"
)

func main() {
	setupEnv()
	router := router.GetRouter(storage.NewMapURLStorage(make(map[string]string)), app.RealURLShortener{})
	router.Run(os.Getenv("SERVER_ADDRESS"))
}

func setupEnv() {
	os.Setenv("SERVER_ADDRESS", ":8080")
	os.Setenv("BASE_URL", "http://localhost:8080/")
}
