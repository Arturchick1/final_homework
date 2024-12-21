package repository

import (
	"context"
	"database/sql"
	"errors"
	"final-project/internal/entity"
)

type TodoRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (t *TodoRepository) DeleteTodo(ctx context.Context, id int64) error {
	res, err := t.db.ExecContext(ctx, `DELETE FROM scheduler WHERE id = :id`, sql.Named("id", id))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("there is no such todo")
	}

	return nil
}

func (t *TodoRepository) GetTodo(ctx context.Context, id int) (entity.TodoRepository, error) {
	if id == 0 {
		row := t.db.QueryRowContext(ctx, `SELECT last_insert_rowid()`)
		err := row.Scan(&id)
		if err != nil {
			return entity.TodoRepository{}, err
		}
	}
	row := t.db.QueryRowContext(ctx, `SELECT id, date, title, comment, repeat
		FROM scheduler WHERE id = :id`, sql.Named("id", id))

	var todo entity.TodoRepository
	err := row.Scan(
		&todo.Id,
		&todo.Date,
		&todo.Title,
		&todo.Comment,
		&todo.Repeat,
	)
	if err != nil {
		return entity.TodoRepository{}, err
	}

	return todo, nil
}

func (t *TodoRepository) GetListTodo(ctx context.Context) ([]entity.TodoRepository, error) {
	rows, err := t.db.QueryContext(ctx, `SELECT id, date, title, comment, repeat
		FROM scheduler ORDER BY date LIMIT 50`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]entity.TodoRepository, 0)
	for rows.Next() {
		todo := entity.TodoRepository{}

		err = rows.Scan(
			&todo.Id,
			&todo.Date,
			&todo.Title,
			&todo.Comment,
			&todo.Repeat,
		)
		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (t *TodoRepository) UpdateTodo(ctx context.Context, todoLogic entity.TodoLogic) error {
	res, err := t.db.ExecContext(ctx, `UPDATE scheduler SET date = :date, title = :title, 
			comment = :comment, repeat = :repeat
			WHERE id = :id`,
		sql.Named("id", todoLogic.Id),
		sql.Named("date", todoLogic.Date),
		sql.Named("title", todoLogic.Title),
		sql.Named("comment", todoLogic.Comment),
		sql.Named("repeat", todoLogic.Repeat))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("there is no such todo")
	}

	return nil
}

func (t *TodoRepository) CreateTodo(ctx context.Context, todoLogic entity.TodoLogic) (int64, error) {
	res, err := t.db.ExecContext(ctx, "INSERT INTO scheduler(date, title, comment, repeat) "+
		"VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", todoLogic.Date),
		sql.Named("title", todoLogic.Title),
		sql.Named("comment", todoLogic.Comment),
		sql.Named("repeat", todoLogic.Repeat))
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
