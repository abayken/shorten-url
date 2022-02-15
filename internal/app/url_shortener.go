package app

import (
	"math/rand"
	"strconv"
)

type URLShortener interface {
	ID() string
}

type RealURLShortener struct {
	URL string
}

func (shortener RealURLShortener) ID() string {
	randomID := rand.Intn((99999 - 10000) + 10000)

	return strconv.Itoa(randomID)
}
