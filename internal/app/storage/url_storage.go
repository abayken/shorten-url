package storage

type URLStorage interface {
	Save(shortURLID, fullURL string)
	Get(shortURLID string) string
}
