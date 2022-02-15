package main

import (
	"github.com/abayken/shorten-url/internal/app/router"
)

func main() {
	router := router.GetRouter()
	router.Run(":8080")
}
