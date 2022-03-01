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
}

func (storage FileURLStorage) Save(shortURLID, fullURL string) {
	file, err := os.OpenFile(storage.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fileModel := FileModel{ShortURLID: shortURLID, FullURL: fullURL}
	bytes, err := json.Marshal(fileModel)

	if err != nil {
		log.Fatal(err)
	}

	writer := bufio.NewWriter(file)
	writer.Write(bytes)
	writer.WriteByte('\n')
	writer.Flush()
}

func (storage FileURLStorage) Get(shortURLID string) string {
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
			return item.FullURL
		}
	}

	return ""
}
