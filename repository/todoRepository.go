package repository

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	error2 "todo/error"
	"todo/logger"
	"todo/model"
)

type TodoRepository struct {
	client *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db}
}

func (r *TodoRepository) Create(todo *model.Todo) (*model.Todo, *model.Error) {
	query := "INSERT INTO todos (title, status) VALUES (?, ?)"
	result, err := r.client.Exec(query, todo.Title, "pending")
	if err != nil {
		return nil, error2.NewInternalServerError("Internal Server Error(Database exec)")
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, error2.NewInternalServerError("Internal Server Error(Database exec)")
	}
	todo.ID = int(lastId)
	return todo, nil
}

func (r *TodoRepository) Update(id int, todo *model.Todo) *model.Error {
	query := "UPDATE todos SET title = ?, status = ? WHERE id = ?"
	_, err := r.client.Exec(query, todo.Title, todo.Status, id)
	if err != nil {
		return error2.NewInternalServerError("Internal Server Error(Database exec)")
	}
	return nil
}

func (r *TodoRepository) Delete(id int) *model.Error {
	query := "DELETE FROM todos WHERE id = ?"
	_, err := r.client.Exec(query, id)
	if err != nil {
		return error2.NewInternalServerError("Internal Server Error(Database exec)")
	}
	return nil
}

func (r *TodoRepository) GetAll() ([]*model.Todo, *model.Error) {
	rows, err := r.client.Query("SELECT id, title, status FROM todos")
	if err != nil {
		return nil, error2.NewInternalServerError("Internal Server Error(Database exec)")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Log.Error("Error closing rows in database")
		}
	}(rows)
	var dtos []*model.Todo
	err = sqlx.StructScan(rows, &dtos)
	if err != nil {
		return nil, error2.NewInternalServerError("Internal Server Error(Database exec)")
	}
	return dtos, nil
}

func (r *TodoRepository) GetById(id int) (*model.Todo, *model.Error) {
	var existingTodo model.Todo
	row := r.client.QueryRow("SELECT id, title, status FROM todos WHERE id = ?", id)
	err := row.Scan(&existingTodo.ID, &existingTodo.Title, &existingTodo.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, error2.NewNotFoundError("Todo not found")
		}
		return nil, error2.NewInternalServerError("Internal Server Error(Database exec)")
	}
	return &existingTodo, nil
}
