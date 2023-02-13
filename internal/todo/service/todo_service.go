package service

import (
	"errors"
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

func (t *todoService) GetTodo(id uint, uid string) (*model.TodoTable, error) {
	var err error
	todo := new(model.TodoTable)
	todo.Id = id
	todo.UserId = uid

	todo, err = t.Repository.GetTodo(todo)

	return todo, err
}

func (t *todoService) CreateTodo(input *model.TodoCreate) error {
	table := new(model.TodoTable)
	table.UserId = input.UserId
	table.Title = input.Title
	table.Description = input.Description

	err := t.Repository.CreateTodo(table)

	return err
}

func (t *todoService) UpdateTodo(input *model.TodoUpdate) (*model.TodoTable, error) {
	table := new(model.TodoTable)
	table.Id = input.Id
	table.UserId = input.UserId

	if input.Title != "" {
		table.Title = input.Title
	}

	if input.Description != "" {
		table.Description = input.Description
	}

	if input.IsCompleted != nil {
		table.IsCompleted = *input.IsCompleted
	}

	updatedTodo, err := t.Repository.UpdateTodo(table)

	return updatedTodo, err
}

func (t *todoService) UpdateTodoCompleted(input *model.TodoUpdate) error {
	var err error
	updatedTodo := new(model.TodoTable)
	updatedTodo.Id = input.Id
	updatedTodo.UserId = input.UserId

	updatedTodo, err = t.Repository.GetTodo(updatedTodo)
	if err != nil {
		return err
	}

	if updatedTodo.UserId != input.UserId {
		return errors.New("identify error")
	}

	if updatedTodo.IsCompleted == 1 {
		return errors.New("todo is completed")
	}

	err = t.Repository.UpdateTodoCompleted(updatedTodo)

	return err
}

func (t *todoService) DeleteTodo(id uint, uid string) (*model.TodoTable, error) {
	var err error
	deletedTodo := new(model.TodoTable)
	deletedTodo.Id = id
	deletedTodo.UserId = uid

	deletedTodo, err = t.Repository.DeleteTodo(deletedTodo)

	return deletedTodo, err
}
