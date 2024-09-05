package service

import (
	"todo/interfaces"
	"todo/model"
)

type TodoService struct {
	repo interfaces.TodoRepoInterface
}

func NewTodoService(repo interfaces.TodoRepoInterface) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) CreateTodo(dto *model.DTO) (*model.DTO, *model.Error) {
	var todo = &model.Todo{
		Title:  dto.Title,
		Status: dto.Status,
	}
	todoNew, err := s.repo.Create(todo)
	if err != nil {
		return nil, err
	}
	var response = &model.DTO{
		ID:     todoNew.ID,
		Title:  todoNew.Title,
		Status: todoNew.Status,
	}
	return response, nil
}

func (s *TodoService) UpdateTodo(id int, dto *model.DTO) *model.Error {
	var todo = &model.Todo{
		Title:  dto.Title,
		Status: dto.Status,
	}
	err := s.repo.Update(id, todo)
	if err != nil {
		return err
	}
	return nil
}

func (s *TodoService) DeleteTodo(id int) *model.Error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *TodoService) GetTodos() ([]*model.DTO, *model.Error) {
	todos, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	var dtos []*model.DTO
	for i, _ := range todos {
		var dto = &model.DTO{
			ID:     todos[i].ID,
			Title:  todos[i].Title,
			Status: todos[i].Status,
		}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}

func (s *TodoService) GetTodoById(id int) (*model.DTO, *model.Error) {
	todo, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	var dto = &model.DTO{
		ID:     todo.ID,
		Title:  todo.Title,
		Status: todo.Status,
	}
	return dto, nil
}
