package storage

/// Storage который сохраняет урлы в словарик
type MapURLStorage struct {
	urlsMap map[string]string
}

func NewMapURLStorage(urls map[string]string) *MapURLStorage {
	return &MapURLStorage{urlsMap: urls}
}

func (storage MapURLStorage) Save(shortURLID, fullURL string) {
	storage.urlsMap[shortURLID] = fullURL
}

func (storage MapURLStorage) Get(shortURLID string) string {
	return storage.urlsMap[shortURLID]
}
