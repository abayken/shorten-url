package storage

/// Абстракция storage-a
/// Есть фейковая имлементация которая сохраняет в словарь
type URLStorage interface {
	Save(shortURLID, fullURL string)
	Get(shortURLID string) string
}
