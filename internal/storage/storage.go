package storage

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

type Storage struct {
	DB *sql.DB
}

func New(pathDatabase string) (Storage, error) {
	db, err := sql.Open("sqlite", pathDatabase)
	if err != nil {
		return Storage{}, err
	}

	q, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS scheduler(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL,
	title VARCHAR(256) NOT NULL,
	comment TEXT NOT NULL,
	repeat VARCHAR(128) NOT NULL);
	CREATE INDEX IF NOT EXISTS todo_date ON scheduler(id);
	CREATE INDEX IF NOT EXISTS todo_date ON scheduler(date);`)
	if err != nil {
		return Storage{}, err
	}

	_, err = q.Exec()
	if err != nil {
		return Storage{}, err
	}

	return Storage{DB: db}, nil
}
