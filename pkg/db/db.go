package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/OlegLuppov/go_final_project/models"
	"github.com/OlegLuppov/go_final_project/pkg/dateutil"
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

// Добавить задачу в БД
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

// Получить список задач
func (db *SchedulerDb) GetTasks(limit int, stringSearch string) (*models.TasklList, error) {
	var selectStr string
	var whereSelect string
	baseSelect := `SELECT id, date, title, comment, repeat FROM scheduler`
	paramsSelect := `ORDER BY date LIMIT 0, :limit`

	if len(stringSearch) > 0 {
		stringSearch = strings.TrimSpace(stringSearch)
		date, err := time.Parse(dateutil.DateLayoutDMY, stringSearch)

		if err != nil {
			whereSelect = `WHERE title LIKE '%' || :search || '%' OR comment LIKE '%' || :search || '%'`
		} else {
			stringSearch = date.Format(dateutil.DateLayoutYMD)
			whereSelect = `WHERE date = :search`
		}

		selectStr = fmt.Sprintf("%s %s %s", baseSelect, whereSelect, paramsSelect)
	} else {
		selectStr = fmt.Sprintf("%s %s", baseSelect, paramsSelect)
	}

	rows, err := db.Db.Query(
		selectStr,
		sql.Named("limit", limit),
		sql.Named("search", stringSearch),
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		task := models.Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return &models.TasklList{Tasks: tasks}, nil
}

// Получить задачу по id
func (db *SchedulerDb) GetTaskById(id string) (*models.Task, error) {
	row := db.Db.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id`,
		sql.Named("id", id),
	)

	var task models.Task

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

// Обновить задачу
func (db *SchedulerDb) UpdateTask(task *models.Task) error {

	res, err := db.Db.Exec(
		`UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID),
	)

	if err != nil {
		return err
	}

	countRows, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if countRows == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}

	return nil
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
