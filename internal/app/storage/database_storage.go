package storage

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

type DatabaseStorage struct {
	URL string
}

type DuplicateURLError struct {
	FullURL    string
	ShortURLID string
}

func (error *DuplicateURLError) Error() string {
	return "There is already url like this"
}

func (storage DatabaseStorage) InitTablesIfNeeded() {
	conn, err := pgx.Connect(context.Background(), storage.URL)

	if err != nil {
		log.Fatal(err)

		return
	}

	_, err = conn.Exec(context.Background(), "create table if not exists url (user_id varchar(100), short_url_id varchar(300), full_url varchar(1000));")

	if err != nil {
		log.Fatal(err)

		return
	}

	_, err = conn.Exec(context.Background(), "create unique index if not exists full_url_index on url (full_url);")

	if err != nil {
		log.Fatal(err)

		return
	}
}

func (storage DatabaseStorage) initDB() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), storage.URL)

	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func (storage DatabaseStorage) Save(shortURLID, fullURL, userID string) error {
	var db = storage.initDB()
	_, err := db.Exec(context.Background(), "insert into url values ($1, $2, $3)", userID, shortURLID, fullURL)

	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				shortURLID, err = storage.getShortURLIDByFullURL(fullURL)

				if err == nil {
					return &DuplicateURLError{FullURL: fullURL, ShortURLID: shortURLID}
				} else {
					return err
				}
			}
		}

		return err
	}

	return nil
}

func (storage DatabaseStorage) getShortURLIDByFullURL(fullURL string) (string, error) {
	var db = storage.initDB()

	var shortURLID string
	err := db.QueryRow(context.Background(), "select short_url_id from url where full_url = $1", fullURL).Scan(&shortURLID)

	return shortURLID, err
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

func (storage DatabaseStorage) BatchURLs(urls []BatchURL) error {
	var db = storage.initDB()

	rows := make([][]interface{}, 0)

	for _, url := range urls {
		row := []interface{}{url.UserID, url.ShortURLID, url.FullURL}
		rows = append(rows, row)
	}

	_, err := db.CopyFrom(context.Background(),
		pgx.Identifier{"url"},
		[]string{"user_id", "short_url_id", "full_url"},
		pgx.CopyFromRows(rows),
	)

	return err
}
