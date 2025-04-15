package models

// Структура задачи
type Task struct {
	ID      string `json:"id"`      // id задачи
	Date    string `json:"date"`    // дата
	Title   string `json:"title"`   // заголовок задачи
	Comment string `json:"comment"` // коментарий к задаче
	Repeat  string `json:"repeat"`  // Правило повторения
}

// Ответ при успешном получении задач
type TasklList struct {
	Tasks []Task `json:"tasks"` // список задач
}

// Ответ при успешном создании задачи
type SuccessfullyСreatedResponse struct {
	Id string `json:"id"`
}

// Ответ при успешном обновлении задачи
type SuccessfullyUpdateResponse struct{}

// Ответ при ошибке
type ErrorResponse struct {
	Error string `json:"error"`
}
