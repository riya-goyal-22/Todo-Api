package interfaces

import "todo/model"

type TodoRepoInterface interface {
	Create(todo *model.Todo) (*model.Todo, *model.Error)
	Update(id int, todo *model.Todo) *model.Error
	Delete(id int) *model.Error
	GetAll() ([]*model.Todo, *model.Error)
	GetById(id int) (*model.Todo, *model.Error)
}

type TodoServiceInterface interface {
	CreateTodo(dto *model.DTO) (*model.DTO, *model.Error)
	UpdateTodo(id int, dto *model.DTO) *model.Error
	DeleteTodo(id int) *model.Error
	GetTodos() ([]*model.DTO, *model.Error)
	GetTodoById(id int) (*model.DTO, *model.Error)
}
