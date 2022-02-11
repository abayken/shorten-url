package app

import (
	"math/rand"
	"strconv"
)

type UrlShortener struct {
	Url string
}

func (shortener UrlShortener) AsShort() string {
	randomId := rand.Intn((99999 - 10000) + 10000)

	return "bit.ly" + strconv.Itoa(randomId)
}
