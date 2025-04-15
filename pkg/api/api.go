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
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
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

// Обработчик возвращает список задач
func (taskHandler *TaskHandler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	stringSearch := r.URL.Query().Get("search")

	data, err := taskHandler.db.GetTasks(50, stringSearch)

	if err != nil {
		setErrResponse(w, http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})

		return
	}

	if data.Tasks == nil {
		setSuccessfulGetResponse(w, http.StatusOK, models.TasklList{Tasks: []models.Task{}})
		return
	}

	setSuccessfulGetResponse(w, http.StatusOK, *data)
}

// Отправляет в ответ ошибку в формате json
func setErrResponse(w http.ResponseWriter, statusCode int, err models.ErrorResponse) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	errEncode := json.NewEncoder(w).Encode(err)

	if errEncode != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, errEncode.Error(), http.StatusInternalServerError)
	}
}

// Отправляет в ответ id задачи в формате json
func setSuccessfulPostResponse(w http.ResponseWriter, statusCode int, data models.SuccessfullyСreatedResponse) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	errEncode := json.NewEncoder(w).Encode(data)

	if errEncode != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, errEncode.Error(), http.StatusInternalServerError)
	}
}

// Отправляет в ответ id задачи в формате json
func setSuccessfulGetResponse(w http.ResponseWriter, statusCode int, data models.TasklList) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	errEncode := json.NewEncoder(w).Encode(data)

	if errEncode != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, errEncode.Error(), http.StatusInternalServerError)
	}
}

// Проверка даты
func checkDate(task *models.Task) error {
	now := time.Now()

	if len(task.Date) == 0 {
		task.Date = now.Format(dateutil.DateLayoutYMD)
	}

	parseDate, err := time.Parse(dateutil.DateLayoutYMD, task.Date)

	if err != nil {
		return err
	}

	if parseDate.Before(now) && task.Date != now.Format(dateutil.DateLayoutYMD) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(dateutil.DateLayoutYMD)
		} else {
			task.Date, err = dateutil.NextDate(now.Format(dateutil.DateLayoutYMD), task.Date, task.Repeat)

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
	router.Get("/api/tasks", taskHandler.GetTasksHandler)
	router.Post("/api/task", taskHandler.PostTaskHandler)

	return router
}
