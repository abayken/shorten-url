package storage

import "log"

/// Storage который сохраняет урлы в словарик
type MapURLStorage struct {
	urlsMap map[string]string
}

func NewMapURLStorage(urls map[string]string) *MapURLStorage {
	return &MapURLStorage{urlsMap: urls}
}

func (storage MapURLStorage) Save(shortURLID, fullURL, userID string) error {
	storage.urlsMap[shortURLID] = fullURL

	return nil
}

func (storage MapURLStorage) Get(shortURLID string) string {
	return storage.urlsMap[shortURLID]
}

func (storage MapURLStorage) FetchUserURLs(userID string) []UserURL {
	return []UserURL{}
}

func (storage MapURLStorage) BatchURLs(urls []BatchURL) error {
	log.Fatal("Данный метод не имеет реализацию")

	return nil
}
