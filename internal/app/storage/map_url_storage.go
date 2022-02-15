package storage

/// Storage который сохраняет урлы в словарик
type MapURLStorage struct {
}

var urlsMap = make(map[string]string)

func (storage MapURLStorage) Save(shortURLID, fullURL string) {
	urlsMap[shortURLID] = fullURL
}

func (storage MapURLStorage) Get(shortURLID string) string {
	return urlsMap[shortURLID]
}
