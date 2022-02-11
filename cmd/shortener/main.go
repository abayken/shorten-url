package main

import (
	"net/http"

	"github.com/abayken/shorten-url/internal/app/handlers"
)

func main() {
	http.HandleFunc("/", handlers.CreateShortURL)
	http.ListenAndServe(":8080", nil)
}
