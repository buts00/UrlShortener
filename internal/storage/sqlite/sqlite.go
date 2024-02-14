package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/buts00/UrlShortener/internal/storage"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
	    id INTEGER PRIMARY key,
	    alias TEXT NOT NULL UNIQUE,
	    url TEXT NOT NUll);
	CREATE INDEX IF NOT EXISTS idx_alias on url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url (url,alias) VALUES (?,?)")
	if err != nil {
		return 0, fmt.Errorf("%s : %w", op, err)
	}
	res, err := stmt.Exec(urlToSave, alias)

	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
		return 0, fmt.Errorf("%s : %w", op, storage.ErrURLExists)
	}
	if err != nil {
		return 0, fmt.Errorf("%s : %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s : %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	op := "storage.sqlite.GetUrl"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")

	if err != nil {
		return "", fmt.Errorf("%s : %w", op, err)
	}
	var newAlias string
	err = stmt.QueryRow(alias).Scan(&newAlias)

	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%s : %w", op, storage.ErrURLNotFound)
	}
	if err != nil {
		return "", fmt.Errorf("%s : %w", op, err)
	}

	return newAlias, nil
}
