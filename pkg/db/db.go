package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title TEXT NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat TEXT CHECK(LENGTH(repeat) <= 128) NOT NULL DEFAULT ""
	);

	CREATE INDEX IF NOT EXISTS date_idx ON scheduler (date);
`

func Connect(dbFile string) (*sql.DB, error) {

	_, err := os.Stat(dbFile)

	checkExist := os.IsNotExist(err)

	if err != nil && checkExist {

		file, err := os.Create(dbFile)

		if err != nil {
			return nil, err
		}

		defer file.Close()

	}

	db, err := sql.Open("sqlite", dbFile)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(schema)

	if err != nil {
		return nil, err
	}

	return db, nil
}
