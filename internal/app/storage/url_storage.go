package storage

/// Абстракция storage-a
/// Есть фейковая имлементация которая сохраняет в словарь
type URLStorage interface {
	Save(shortURLID, fullURL, userID string)
	Get(shortURLID string) string
	FetchUserURLs(userID string) []UserURL
}

/// Урл определенного юзера
type UserURL struct {
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}

func (url UserURL) BaseURLAppended(baseURL string) UserURL {
	return UserURL{Short: baseURL + "/" + url.Short, Original: url.Original}
}
