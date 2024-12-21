package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"final-project/internal/entity"
	"final-project/internal/logic"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

type TodoHandlers struct {
	todoLogic *logic.TodoLogic
}

func NewTodoHandler(logic *logic.TodoLogic) *TodoHandlers {
	return &TodoHandlers{todoLogic: logic}
}

func (t *TodoHandlers) newGroupTodo(eg *echo.Group) {
	eg.PUT("", t.update)
	eg.POST("/done", t.done)
	eg.POST("", t.create)
	eg.GET("", t.get)
	eg.DELETE("", t.delete)
}

func (t *TodoHandlers) Start(port string) *echo.Echo {
	e := echo.New()

	e.Static("/", "web")
	e.GET("/api/nextdate", nextDate)
	e.GET("/api/tasks", t.getList)
	t.newGroupTodo(e.Group("/api/task"))

	err := e.Start(":" + port)
	if err != http.ErrServerClosed {
		panic(err)
	}
	return e
}

func (t *TodoHandlers) delete(c echo.Context) error {
	result := entity.ResultCreate{}
	idStr := c.QueryParam("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		result.Error = err.Error()
		return c.JSON(http.StatusBadRequest, result)
	}

	ctx := c.Request().Context()
	err = t.todoLogic.DeleteTodo(ctx, id)
	if err != nil {
		result.Error = err.Error()
		return c.JSON(http.StatusBadRequest, result)
	}

	return c.JSON(http.StatusOK, result)
}

func (t *TodoHandlers) update(c echo.Context) error {
	result := entity.ResultCreate{}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(c.Request().Body)
	if err != nil {
		result.Error = err.Error()
		return c.JSON(http.StatusBadRequest, result)
	}

	todoJSON := entity.TodoHandler{}
	if err = json.Unmarshal(buf.Bytes(), &todoJSON); err != nil {
		result.Error = err.Error()
		return c.JSON(http.StatusBadRequest, result)
	}

	todo := entity.TodoHandler{
		Id:      todoJSON.Id,
		Date:    todoJSON.Date,
		Title:   todoJSON.Title,
		Comment: todoJSON.Comment,
		Repeat:  todoJSON.Repeat,
	}

	if todo.Title == "" {
		result.Error = "Title is required"
		return c.JSON(http.StatusBadRequest, result)
	}

	todo.Date, err = checkDate(todo.Date, todo.Repeat)
	if err != nil {
		result.Error = err.Error()
		return c.JSON(http.StatusBadRequest, result)
	}

	ctx := c.Request().Context()

	err = t.todoLogic.UpdateTodo(ctx, todo)
	if err != nil {
		result.Error = err.Error()
		return c.JSON(http.StatusBadRequest, result)
	}

	return c.JSON(http.StatusOK, result)
}

func (t *TodoHandlers) done(c echo.Context) error {
	result := entity.ResultCreate{}
	idStr := c.QueryParam("id")
	if idStr == "" {
		result.Error = "param id is empty"
		return c.JSON(http.StatusBadRequest, result)
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		result.Error = err.Error()
		return c.JSON(http.StatusBadRequest, result)
	}

	ctx := c.Request().Context()
	todoLogic, err := t.todoLogic.GetTodo(ctx, int(id))
	if err != nil {
		result.Error = err.Error()
		return c.JSON(http.StatusBadRequest, result)
	}

	if todoLogic.Repeat == "" {
		err = t.todoLogic.DeleteTodo(ctx, todoLogic.Id)
		if err != nil {
			result.Error = err.Error()
			return c.JSON(http.StatusBadRequest, result)
		}

		return c.JSON(http.StatusOK, result)
	}

	now := time.Now()
	todoLogic.Date, err = logic.NextDate(now, todoLogic.Date, todoLogic.Repeat)
	if err != nil {
		result.Error = err.Error()
		return c.JSON(http.StatusBadRequest, result)
	}
	todoHandler := entity.TodoHandler{
		Id:      idStr,
		Date:    todoLogic.Date,
		Title:   todoLogic.Title,
		Comment: todoLogic.Comment,
		Repeat:  todoLogic.Repeat,
	}

	err = t.todoLogic.UpdateTodo(ctx, todoHandler)
	if err != nil {
		result.Error = err.Error()
		return c.JSON(http.StatusBadRequest, result)
	}

	return c.JSON(http.StatusOK, result)
}

func (t *TodoHandlers) get(c echo.Context) error {
	todoResult := entity.TodoHandler{}
	ctx := c.Request().Context()
	idStr := c.QueryParam("id")
	if idStr == "" {
		todoResult.Error = errors.New("param id is empty").Error()
		return c.JSON(http.StatusBadRequest, todoResult)
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		todoResult.Error = err.Error()
		return c.JSON(http.StatusBadRequest, todoResult)
	}

	todoLogic, err := t.todoLogic.GetTodo(ctx, int(id))
	if err != nil {
		todoResult.Error = err.Error()
		return c.JSON(http.StatusBadRequest, todoResult)
	}

	todoResult = entity.TodoHandler{
		Id:      idStr,
		Date:    todoLogic.Date,
		Title:   todoLogic.Title,
		Comment: todoLogic.Comment,
		Repeat:  todoLogic.Repeat,
	}

	return c.JSON(http.StatusOK, todoResult)
}

func (t *TodoHandlers) getList(c echo.Context) error {
	result := entity.ResultGetList{
		Tasks: []entity.TodoHandler{},
	}
	ctx := c.Request().Context()

	todosLogic, err := t.todoLogic.GetList(ctx)
	if err != nil {
		result.Error = err.Error()
		return c.JSON(http.StatusBadRequest, result)
	}
	if todosLogic == nil {
		return c.JSON(http.StatusOK, result)
	}

	todoHandlers := make([]entity.TodoHandler, 0, len(todosLogic))
	for _, v := range todosLogic {
		idStr := strconv.FormatInt(v.Id, 10)
		todo := entity.TodoHandler{
			Id:      idStr,
			Date:    v.Date,
			Title:   v.Title,
			Comment: v.Comment,
			Repeat:  v.Repeat,
		}

		todoHandlers = append(todoHandlers, todo)
	}
	result.Tasks = todoHandlers

	return c.JSON(http.StatusOK, result)
}

func checkDate(date string, repeat string) (string, error) {
	now := time.Now().Truncate(24 * time.Hour)
	if date == "" {
		date = now.Format("20060102")
		return date, nil
	}

	dateParse, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	if dateParse.Compare(now) < 0 {
		if repeat == "" {
			date = now.Format("20060102")
			return date, nil
		} else {
			date, err = logic.NextDate(now, date, repeat)
			if err != nil {
				return "", err
			}
		}
	}

	return date, nil
}

func (t *TodoHandlers) create(c echo.Context) error {
	result := entity.ResultCreate{}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(c.Request().Body)
	if err != nil {
		result.Error = err.Error()
		return c.JSON(http.StatusBadRequest, result)
	}

	todo := entity.TodoHandler{}
	if err = json.Unmarshal(buf.Bytes(), &todo); err != nil {
		result.Error = err.Error()
		return c.JSON(http.StatusBadRequest, result)
	}

	if todo.Title == "" {
		result.Error = "Title is required"
		return c.JSON(http.StatusBadRequest, result)
	}

	todo.Date, err = checkDate(todo.Date, todo.Repeat)
	if err != nil {
		result.Error = err.Error()
		return c.JSON(http.StatusBadRequest, result)
	}

	ctx := c.Request().Context()

	idLastCreated, err := t.todoLogic.CreateTodo(ctx, todo)
	if err != nil {
		result.Error = err.Error()
		return c.JSON(http.StatusBadRequest, result)
	}

	c.Response().Header().Set("Content-Type", "application/json; charset=UTF-8")

	result = entity.ResultCreate{Id: idLastCreated}
	return c.JSON(http.StatusCreated, result)
}

func nextDate(c echo.Context) error {
	nowStr := c.FormValue("now")
	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	date := c.FormValue("date")
	repeat := c.FormValue("repeat")

	nextDate, err := logic.NextDate(now, date, repeat)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, nextDate)
}
