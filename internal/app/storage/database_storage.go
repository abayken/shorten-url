package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
)

type DatabaseStorage struct {
	Url string
}

func (storage DatabaseStorage) InitTablesIfNeeded() {
	conn, err := pgx.Connect(context.Background(), storage.Url)

	if err != nil {
		log.Fatal(err)

		return
	}

	_, err = conn.Exec(context.Background(), "create table if not exists url (user_id varchar(100), short_url_id varchar(300), full_url varchar(1000));")

	if err != nil {
		log.Fatal(err)

		return
	}
}

func (storage DatabaseStorage) initDB() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), storage.Url)

	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func (storage DatabaseStorage) Save(shortURLID, fullURL, userID string) {
	var db = storage.initDB()
	db.Exec(context.Background(), "insert into url values ($1, $2, $3)", userID, shortURLID, fullURL)
}

func (storage DatabaseStorage) Get(shortURLID string) string {
	var db = storage.initDB()

	var fullURL string
	db.QueryRow(context.Background(), "select full_url from url where short_url_id = $1", shortURLID).Scan(&fullURL)

	return fullURL
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
