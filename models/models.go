package models

type Task struct {
	ID      string `json:"id"`      //id задачи
	Date    string `json:"date"`    // дата
	Title   string `json:"title"`   // заголовок задачи
	Comment string `json:"comment"` //коментарий к задаче
	Repeat  string `json:"repeat"`  //Правило повторения
}

type SuccessfullyСreatedResponse struct {
	Id string `json:"id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
