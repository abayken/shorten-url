package handlers

import (
	"io"
	"net/http"

	"github.com/abayken/shorten-url/internal/app"
)

var urlsMap = make(map[string]string)

func CreateShortURL(w http.ResponseWriter, r *http.Request) {
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

		shortUrlId := urlShortener.Id()

		urlsMap[shortUrlId] = url

		w.WriteHeader(http.StatusCreated)

		w.Write([]byte("http://localhost:8080/" + shortUrlId))
	case http.MethodGet:
		shortUrlId := r.URL.Path[1:]
		if fullURL, ok := urlsMap[shortUrlId]; ok {
			w.Header().Set("Location", fullURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	default:
		break
	}
}
