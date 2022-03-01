package app

import (
	"math/rand"
	"strconv"
	"time"
)

type URLShortener interface {
	ID() string
}

type RealURLShortener struct {
	URL string
}

func (shortener RealURLShortener) ID() string {
	rand.Seed(time.Now().UnixNano())
	randomID := rand.Intn((99999 - 10000) + 10000)

	return strconv.Itoa(randomID)
}
