package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
)

type DatabaseStorage struct {
	Url string
}

func (storage DatabaseStorage) initDB() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), storage.Url)

	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func (storage DatabaseStorage) Save(shortURLID, fullURL, userID string) {

}

func (storage DatabaseStorage) Get(shortURLID string) string {
	return "todo"
}

func (storage DatabaseStorage) FetchUserURLs(userID string) []UserURL {
	var db = storage.initDB()

	rows, err := db.Query(context.Background(), "select * from url where user_id = $1", userID)

	if err != nil {
		log.Fatal(err)
	}

	var urls []UserURL

	for rows.Next() {
		var userURL UserURL

		var userID string
		err = rows.Scan(&userID, &userURL.Short, &userURL.Original)

		if err == nil {
			urls = append(urls, userURL)
		}
	}

	return urls
}
