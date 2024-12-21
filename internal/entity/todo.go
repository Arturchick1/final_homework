package entity

type ResultCreate struct {
	Id    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

type ResultGetList struct {
	Tasks []TodoHandler `json:"tasks"`
	Error string        `json:"error,omitempty"`
}

type ResultGetUpdateTodo struct {
	Task  TodoHandler `json:"tasks,omitempty"`
	Error string      `json:"error,omitempty"`
}

type TodoRepository struct {
	Id      int64
	Date    string
	Title   string
	Comment string
	Repeat  string
}

type TodoLogic struct {
	Id      int64  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TodoHandler struct {
	Id      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
	Error   string `json:"error,omitempty"`
}
