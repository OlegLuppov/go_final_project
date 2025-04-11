package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/OlegLuppov/go_final_project/pkg/dateutil"

	"github.com/OlegLuppov/go_final_project/models"
	"github.com/OlegLuppov/go_final_project/pkg/db"

	"github.com/go-chi/chi"
)

type TaskHandler struct {
	db *db.SchedulerDb
}

// Обработчик возвращает следубщую дату
func (taskHandler *TaskHandler) NextDateHandler(w http.ResponseWriter, r *http.Request) {
	queryNow := r.FormValue("now")
	queryDate := r.FormValue(("date"))
	queryRepeat := r.FormValue(("repeat"))

	nextDate, err := dateutil.NextDate(queryNow, queryDate, queryRepeat)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}

// Обработчик создает задачу и возвращает id при успешном создании, иначе ошибку
func (taskHandler *TaskHandler) PostTaskHandler(w http.ResponseWriter, r *http.Request) {

	task := new(models.Task)
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)

	if err != nil {
		setErrResponse(w, http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})

		return
	}

	if err := json.Unmarshal(buf.Bytes(), task); err != nil {
		setErrResponse(w, http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})

		return
	}

	if len(task.Title) == 0 {
		setErrResponse(w, http.StatusBadRequest, models.ErrorResponse{
			Error: "the title expected non-empty, but got empty",
		})

		return
	}

	err = checkDate(task)

	if err != nil {
		setErrResponse(w, http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})

		return
	}

	id, err := taskHandler.db.AddTask(task)

	if err != nil {
		setErrResponse(w, http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})

		return
	}

	setSuccessfulPostResponse(w, http.StatusCreated, models.SuccessfullyСreatedResponse{
		Id: strconv.Itoa(int(id)),
	})
}

// Отправляет в ответ ошибку в формате json
func setErrResponse(w http.ResponseWriter, statusCode int, err models.ErrorResponse) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(err)
}

func setSuccessfulPostResponse(w http.ResponseWriter, statusCode int, data models.SuccessfullyСreatedResponse) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

// Проверка даты
func checkDate(task *models.Task) error {
	now := time.Now()

	if len(task.Date) == 0 {
		task.Date = now.Format(dateutil.DateLayout)
	}

	parseDate, err := time.Parse(dateutil.DateLayout, task.Date)

	if err != nil {
		return err
	}

	if parseDate.Before(now) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(dateutil.DateLayout)
		} else {
			task.Date, err = dateutil.NextDate(now.Format(dateutil.DateLayout), task.Date, task.Repeat)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Регистрация Обработчиков
func RegisterHandlers(db *db.SchedulerDb) *chi.Mux {
	router := chi.NewMux()

	taskHandler := TaskHandler{db: db}

	router.Get("/api/nextdate", taskHandler.NextDateHandler)
	router.Post("/api/task", taskHandler.PostTaskHandler)

	return router
}
