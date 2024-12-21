package main

import (
	"final-project/internal/http"
	"final-project/internal/logic"
	"final-project/internal/repository"
	"final-project/internal/storage"
	"final-project/tests"
	"os"
	"strconv"
)

func main() {

	httpPort := os.Getenv("TODO_PORT")
	if httpPort == "" {
		httpPort = strconv.Itoa(tests.Port)
	}
	pathDatabase := os.Getenv("TODO_DBFILE")
	if pathDatabase == "" {
		pathDatabase = tests.DBFile
	}

	storage, err := storage.New(pathDatabase)
	if err != nil {
		panic(err)
	}
	defer storage.DB.Close()

	todoRepository := repository.New(storage.DB)
	todoLogic := logic.New(todoRepository)
	todoHandler := http.NewTodoHandler(todoLogic)

	todoHandler.Start(httpPort)
}
