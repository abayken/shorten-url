package main

import (
	"github.com/abayken/shorten-url/internal/app"
	"github.com/abayken/shorten-url/internal/app/router"
	"github.com/abayken/shorten-url/internal/app/storage"
)

func main() {
	router := router.GetRouter(storage.NewMapURLStorage(make(map[string]string)), app.RealURLShortener{})
	router.Run(":8080")
}
