package storage

/// Абстракция storage-a
/// Есть фейковая имлементация которая сохраняет в словарь
type URLStorage interface {
	Save(shortURLID, fullURL, userID string) error
	Get(shortURLID string) string
	FetchUserURLs(userID string) []UserURL
	BatchURLs(urls []BatchURL) error
}

/// Урл определенного юзера
type UserURL struct {
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}

func (url UserURL) BaseURLAppended(baseURL string) UserURL {
	return UserURL{Short: baseURL + "/" + url.Short, Original: url.Original}
}

/// Формат по которому пачкой добавляются урлы в базу
type BatchURL struct {
	UserID     string
	ShortURLID string
	FullURL    string
}
