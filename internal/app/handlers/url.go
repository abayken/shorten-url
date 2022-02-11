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

		urlShortener := app.UrlShortener{Url: url}

		shortUrl := urlShortener.AsShort()

		urlsMap[shortUrl] = url

		w.WriteHeader(http.StatusCreated)

		w.Write([]byte(shortUrl))
	case http.MethodGet:
		shortUrl := r.URL.Path[1:]

		if fullUrl, ok := urlsMap[shortUrl]; ok {
			w.Header().Set("Location", fullUrl)
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	default:
		break
	}
}
