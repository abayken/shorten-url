package storage

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v4"
)

type DatabaseStorage struct {
	DB *pgx.Conn
}

type DuplicateURLError struct {
	FullURL    string
	ShortURLID string
}

func (error *DuplicateURLError) Error() string {
	return "There is already url like this"
}

type DeletedURLError struct {
}

func (error *DeletedURLError) Error() string {
	return "It is deleted URL"
}

func (storage DatabaseStorage) InitTablesIfNeeded() {
	_, err := storage.DB.Exec(context.Background(), "create table if not exists url (user_id varchar(100), short_url_id varchar(300), full_url varchar(1000), deleted boolean default false);")

	if err != nil {
		log.Fatal(err)

		return
	}

	_, err = storage.DB.Exec(context.Background(), "create unique index if not exists full_url_index on url (full_url);")

	if err != nil {
		log.Fatal(err)

		return
	}
}

func (storage DatabaseStorage) Save(shortURLID, fullURL, userID string) error {
	_, err := storage.DB.Exec(context.Background(), "insert into url values ($1, $2, $3)", userID, shortURLID, fullURL)

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
	var shortURLID string
	err := storage.DB.QueryRow(context.Background(), "select short_url_id from url where full_url = $1", fullURL).Scan(&shortURLID)

	return shortURLID, err
}

func (storage DatabaseStorage) Get(shortURLID string) (string, error) {
	var fullURL string
	var deleted bool
	storage.DB.QueryRow(context.Background(), "select full_url, deleted from url where short_url_id = $1", shortURLID).Scan(&fullURL, &deleted)

	if deleted {
		return "", &DeletedURLError{}
	}

	return fullURL, nil
}

func (storage DatabaseStorage) FetchUserURLs(userID string) []UserURL {
	rows, err := storage.DB.Query(context.Background(), "select * from url where user_id = $1", userID)

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
	rows := make([][]interface{}, 0)

	for _, url := range urls {
		row := []interface{}{url.UserID, url.ShortURLID, url.FullURL}
		rows = append(rows, row)
	}

	_, err := storage.DB.CopyFrom(context.Background(),
		pgx.Identifier{"url"},
		[]string{"user_id", "short_url_id", "full_url"},
		pgx.CopyFromRows(rows),
	)

	return err
}

func (storage DatabaseStorage) DeleteURLs(urlIDs []string, userID string) error {
	IDs := &pgtype.TextArray{}
	err := IDs.Set(urlIDs)

	if err != nil {
		return err
	}

	_, err = storage.DB.Exec(context.Background(), "update url set deleted = true where user_id = $1 and short_url_id = any($2)", userID, IDs)

	return err
}
