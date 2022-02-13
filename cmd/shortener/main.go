package main

import (
	"net/http"

	"github.com/abayken/shorten-url/internal/app/handlers"
	"github.com/abayken/shorten-url/internal/app/storage"
)

func main() {
	handler := handlers.URLHandler{Storage: storage.MapURLStorage{}}
	http.HandleFunc("/", handler.ServerHTTP)
	http.ListenAndServe(":8080", nil)
}
