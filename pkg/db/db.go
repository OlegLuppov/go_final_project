package db

import (
	"database/sql"
	"os"

	"github.com/OlegLuppov/go_final_project/models"
	_ "modernc.org/sqlite"
)

type SchedulerDb struct {
	Db *sql.DB
}

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

func (db *SchedulerDb) AddTask(task *models.Task) (int64, error) {
	var id int64

	res, err := db.Db.Exec(
		`INSERT INTO scheduler (date,title,comment,repeat) VALUES (:date,:title,:comment,:repeat)`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)

	if err != nil {
		return id, err
	}

	id, err = res.LastInsertId()

	if err != nil {
		return id, err
	}

	return id, nil
}

// Подключение к БД
func Connect(dbFile string) (*SchedulerDb, error) {

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

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec(schema)

	if err != nil {
		return nil, err
	}

	return &SchedulerDb{Db: db}, nil
}
