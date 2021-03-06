package storage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type FileURLStorage struct {
	Path string
}

type FileModel struct {
	ShortURLID string `json:"short_url_id"`
	FullURL    string `json:"full_url"`
	UserID     string `json:"user_id"`
}

func (storage FileURLStorage) Save(shortURLID, fullURL, userID string) error {
	file, err := os.OpenFile(storage.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)

	if err != nil {
		return err
	}

	defer file.Close()

	fileModel := FileModel{ShortURLID: shortURLID, FullURL: fullURL, UserID: userID}
	bytes, err := json.Marshal(fileModel)

	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	writer.Write(bytes)
	writer.WriteByte('\n')
	writer.Flush()

	return nil
}

func (storage FileURLStorage) Get(shortURLID string) (string, error) {
	file, err := os.OpenFile(storage.Path, os.O_RDONLY|os.O_CREATE, 0777)

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		bytes := []byte(scanner.Text())
		var item FileModel
		err = json.Unmarshal(bytes, &item)
		if err == nil && item.ShortURLID == shortURLID {
			return item.FullURL, nil
		}
	}

	return "", nil
}

func (storage FileURLStorage) FetchUserURLs(userID string) []UserURL {
	file, err := os.OpenFile(storage.Path, os.O_RDONLY|os.O_CREATE, 0777)

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	var urls []UserURL

	for scanner.Scan() {
		bytes := []byte(scanner.Text())
		var item FileModel
		err = json.Unmarshal(bytes, &item)
		if err == nil && item.UserID == userID {
			urls = append(urls, UserURL{Short: item.ShortURLID, Original: item.FullURL})
		}
	}

	return urls
}

func (storage FileURLStorage) BatchURLs(urls []BatchURL) error {
	log.Fatal("Данный метод не имеет реализацию")

	return nil
}

func (storage FileURLStorage) DeleteURLs(urlIDs []string, userID string) error {
	return nil
}
