package repository

import (
	model "simple-backend/internal/domain/todo"

	"gorm.io/gorm"
)

type todoRepo struct {
	db *gorm.DB
}

func Init(db *gorm.DB) model.TodoRepoInterface {
	return &todoRepo{db: db}
}

func (t *todoRepo) GetAllTodo(field *model.TodoQueries) (*int64, []*model.TodoTable, error) {
	var totalCount int64
	allTodo := make([]*model.TodoTable, 0)
	query := t.db.Model(&model.TodoTable{})
	// query := t.db.Model(&model.TodoTable{}).Preload(clause.Associations)

	if field.Title != nil {
		query.Where("title = ?", field.Title)
	}

	if field.IsCompleted != nil {
		query.Where("is_completed = ?", field.IsCompleted)
	}

	if field.CreatedAt != nil {
		query.Where("created_at >= ?", field.CreatedAt)
	}

	if field.UpdatedAt != nil {
		query.Where("updated_at >= ?", field.UpdatedAt)
	}

	if field.OrderBy != nil {
		query.Order("created_at " + *field.OrderBy)
	}

	if field.CurrentPage < 1 {
		field.CurrentPage = 1
	}

	if field.PerPage < 1 {
		field.PerPage = 15
	}

	result := query.Count(&totalCount).Offset(int((field.CurrentPage - 1) * field.PerPage)).Limit(int(field.PerPage)).Find(&allTodo)

	if result.Error != nil {
		return nil, nil, result.Error
	}

	return &totalCount, allTodo, nil
}

func (t *todoRepo) GetTodo(todo *model.TodoTable) (*model.TodoTable, error) {
	result := t.db.Model(&model.TodoTable{}).Where("user_id = ?", todo.UserId).First(&todo, todo.Id)
	if result.Error != nil {
		return nil, result.Error
	}

	return todo, nil
}

func (t *todoRepo) CreateTodo(todo *model.TodoTable) error {
	result := t.db.Model(&model.TodoTable{}).Create(&todo)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (t *todoRepo) UpdateTodo(id int, newTodo *model.TodoTable) (*model.TodoTable, error) {
	query := t.db.Model(&model.TodoTable{}).Where("user_id = ?", newTodo.UserId)
	updatedTodo := &model.TodoTable{}

	result := query.First(&updatedTodo, id)
	if result.Error != nil {
		return nil, result.Error
	}

	updatedTodo.Title = newTodo.Title
	updatedTodo.Description = newTodo.Description
	updatedTodo.IsCompleted = newTodo.IsCompleted

	result = query.Select("title", "description").Save(&updatedTodo)
	if result.Error != nil {
		return nil, result.Error
	}

	return updatedTodo, nil
}

func (t *todoRepo) UpdateTodoCompleted(updatedTodo *model.TodoTable) (*int64, error) {
	query := t.db.Model(&model.TodoTable{}).Where("user_id = ?", updatedTodo.UserId)

	result := query.First(&updatedTodo).Update("is_completed", 1)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &result.RowsAffected, nil
}

func (t *todoRepo) DeleteTodo(deletedTodo *model.TodoTable) (*model.TodoTable, error) {
	query := t.db.Model(&model.TodoTable{}).Where("user_id = ?", deletedTodo.UserId)

	result := query.First(&deletedTodo).Delete(&deletedTodo)
	if result.Error != nil {
		return nil, result.Error
	}

	return deletedTodo, nil
}
