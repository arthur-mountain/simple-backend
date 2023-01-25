package service

import (
	model "simple-backend/internal/domain/todo"
	repo "simple-backend/internal/todo/repository/mysql"

	"gorm.io/gorm"
)

type todoService struct {
	Repository model.TodoRepoInterface
}

func Init(db *gorm.DB) model.TodoServiceInterface {
	return &todoService{
		Repository: repo.Init(db),
	}
}

func (t *todoService) GetAllTodo(field *model.TodoQueries) (*int64, []*model.TodoTable, error) {
	totalCount, allTodo, err := t.Repository.GetAllTodo(field)

	return totalCount, allTodo, err
}

func (t *todoService) GetTodo(todo *model.TodoTable) (*model.TodoTable, error) {
	todo, err := t.Repository.GetTodo(todo)

	return todo, err
}

func (t *todoService) CreateTodo(input *model.TodoTable) error {
	err := t.Repository.CreateTodo(input)

	return err
}

func (t *todoService) UpdateTodo(id int, newTodo *model.TodoTable) (*model.TodoTable, error) {
	updatedTodo, err := t.Repository.UpdateTodo(id, newTodo)

	return updatedTodo, err
}

func (t *todoService) UpdateTodoCompleted(updatedTodo *model.TodoTable) (*int64, error) {
	isSuccess, err := t.Repository.UpdateTodoCompleted(updatedTodo)

	return isSuccess, err
}

func (t *todoService) DeleteTodo(deletedTodo *model.TodoTable) (*model.TodoTable, error) {
	deletedTodo, err := t.Repository.DeleteTodo(deletedTodo)

	return deletedTodo, err
}
