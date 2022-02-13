package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/abayken/shorten-url/internal/app"
	"github.com/abayken/shorten-url/internal/app/storage"
)

type URLHandler struct {
	Storage storage.URLStorage
}

func (handler *URLHandler) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		defer r.Body.Close()

		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		url := string(body)

		urlShortener := app.URLShortener{URL: url}

		shortURLID := urlShortener.ID()

		handler.Storage.Save(shortURLID, url)

		w.WriteHeader(http.StatusCreated)

		w.Write([]byte("http://localhost:8080/" + shortURLID))
	case http.MethodGet:
		fmt.Println("called")
		shortURLID := r.URL.Path[1:]

		fullURL := handler.Storage.Get(shortURLID)

		if fullURL != "" {
			fmt.Println("if tag: " + fullURL)
			w.Header().Set("Location", fullURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			fmt.Println("else tag")
			w.WriteHeader(http.StatusBadRequest)
		}
	default:
		break
	}
}
