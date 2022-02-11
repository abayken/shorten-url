package app

import (
	"math/rand"
	"strconv"
)

type URLShortener struct {
	URL string
}

func (shortener URLShortener) Id() string {
	randomID := rand.Intn((99999 - 10000) + 10000)

	return strconv.Itoa(randomID)
}
