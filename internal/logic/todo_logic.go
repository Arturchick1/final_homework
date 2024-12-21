package logic

import (
	"context"
	"errors"
	"final-project/internal/entity"
	"final-project/internal/repository"
	_ "final-project/internal/repository"
	"strconv"
	"strings"
	"time"
)

type TodoLogic struct {
	todoRepo *repository.TodoRepository
}

func New(repo *repository.TodoRepository) *TodoLogic {
	return &TodoLogic{todoRepo: repo}
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat is empty")
	}
	repeatTime, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	splitRepeatStr := strings.Split(repeat, " ")
	switch splitRepeatStr[0] {
	case "d":
		if len(splitRepeatStr) == 1 {
			return "", errors.New("invalid count day")
		}

		dayCount, err := strconv.Atoi(splitRepeatStr[1])
		if err != nil {
			return "", err
		}
		if dayCount > 400 {
			return "", errors.New("invalid count day")
		}

		for {
			repeatTime = repeatTime.AddDate(0, 0, dayCount)
			if repeatTime.Compare(now) > 0 {
				break
			}
		}
	case "y":
		repeatTime = repeatTime.AddDate(1, 0, 0)
		for repeatTime.Compare(now) < 0 {
			repeatTime = repeatTime.AddDate(1, 0, 0)
		}
	default:
		return "", errors.New("invalid argument")
	}

	return repeatTime.Format("20060102"), nil
}

func (t *TodoLogic) DeleteTodo(ctx context.Context, id int64) error {
	err := t.todoRepo.DeleteTodo(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoLogic) GetTodo(ctx context.Context, id int) (entity.TodoLogic, error) {
	todoRepo, err := t.todoRepo.GetTodo(ctx, id)
	if err != nil {
		return entity.TodoLogic{}, err
	}

	todoLogic := entity.TodoLogic{
		Id:      todoRepo.Id,
		Date:    todoRepo.Date,
		Title:   todoRepo.Title,
		Comment: todoRepo.Comment,
		Repeat:  todoRepo.Repeat,
	}

	return todoLogic, nil
}

func (t *TodoLogic) GetList(ctx context.Context) ([]entity.TodoLogic, error) {
	todosRepository, err := t.todoRepo.GetListTodo(ctx)
	if err != nil {
		return nil, err
	}
	if len(todosRepository) == 0 {
		return nil, nil
	}

	todosLogic := make([]entity.TodoLogic, 0, len(todosRepository))
	for _, v := range todosRepository {
		todo := entity.TodoLogic{
			Id:      v.Id,
			Date:    v.Date,
			Title:   v.Title,
			Comment: v.Comment,
			Repeat:  v.Repeat,
		}

		todosLogic = append(todosLogic, todo)
	}

	return todosLogic, nil
}

func (t *TodoLogic) UpdateTodo(ctx context.Context, todoHandler entity.TodoHandler) error {
	idInt, err := strconv.ParseInt(todoHandler.Id, 10, 64)
	if err != nil {
		return err
	}

	todo := entity.TodoLogic{
		Id:      idInt,
		Date:    todoHandler.Date,
		Title:   todoHandler.Title,
		Comment: todoHandler.Comment,
		Repeat:  todoHandler.Repeat,
	}

	err = t.todoRepo.UpdateTodo(ctx, todo)
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoLogic) CreateTodo(ctx context.Context, todoHandler entity.TodoHandler) (int64, error) {
	todo := entity.TodoLogic{
		Date:    todoHandler.Date,
		Title:   todoHandler.Title,
		Comment: todoHandler.Comment,
		Repeat:  todoHandler.Repeat,
	}

	idLastCreated, err := t.todoRepo.CreateTodo(ctx, todo)
	if err != nil {
		return 0, err
	}

	return idLastCreated, nil
}
