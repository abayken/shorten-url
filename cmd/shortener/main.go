package main

import (
	"github.com/abayken/shorten-url/internal/app"
	"github.com/abayken/shorten-url/internal/app/router"
	"github.com/abayken/shorten-url/internal/app/storage"
)

func main() {
	router := router.GetRouter(storage.MapURLStorage{}, app.RealURLShortener{})
	router.Run(":8080")
}
