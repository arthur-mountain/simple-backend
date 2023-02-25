package repository

import (
	"net/http"
	model "simple-backend/internal/domain/todo"
	errorUtils "simple-backend/internal/utils/error"

	"gorm.io/gorm"
)

type todoRepo struct {
	db *gorm.DB
}

func Init(db *gorm.DB) model.TodoRepoInterface {
	return &todoRepo{db: db}
}

func (t *todoRepo) GetAllTodo(field *model.TodoQueries) (*int64, []*model.TodoTable, *errorUtils.CustomError) {
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
		return nil, nil, errorUtils.CheckGormError(result.Error)
	}

	return &totalCount, allTodo, nil
}

func (t *todoRepo) GetTodo(todo *model.TodoTable) (*model.TodoTable, *errorUtils.CustomError) {
	result := t.db.Model(&model.TodoTable{}).Where("user_id = ?", todo.UserId).First(&todo, todo.Id)
	if result.Error != nil {
		return nil, errorUtils.CheckGormError(result.Error)
	}

	return todo, nil
}

func (t *todoRepo) CreateTodo(todo *model.TodoTable) *errorUtils.CustomError {
	result := t.db.Model(&model.TodoTable{}).Create(&todo)

	if result.Error != nil {
		return errorUtils.CheckGormError(result.Error)
	}

	return nil
}

func (t *todoRepo) UpdateTodo(newTodo *model.TodoTable) (*model.TodoTable, *errorUtils.CustomError) {
	query := t.db.Model(&model.TodoTable{}).Where("user_id = ?", newTodo.UserId)

	result := query.Select("title", "description").Save(&newTodo)
	if result.Error != nil {
		return nil, errorUtils.CheckGormError(result.Error)
	}

	return newTodo, nil
}

func (t *todoRepo) UpdateTodoCompleted(updatedTodo *model.TodoTable) *errorUtils.CustomError {
	query := t.db.Model(&model.TodoTable{}).Where("user_id = ?", updatedTodo.UserId)

	result := query.First(&updatedTodo).Update("is_completed", 1)
	if result.Error != nil {
		return errorUtils.CheckGormError(result.Error)
	}

	if result.RowsAffected == 0 {
		return errorUtils.NewCustomError(http.StatusNotFound, gorm.ErrRecordNotFound.Error(), nil)
	}

	return nil
}

func (t *todoRepo) DeleteTodo(deletedTodo *model.TodoTable) (*model.TodoTable, *errorUtils.CustomError) {
	query := t.db.Model(&model.TodoTable{}).Where("user_id = ?", deletedTodo.UserId)

	result := query.First(&deletedTodo).Delete(&deletedTodo)
	if result.Error != nil {
		return nil, errorUtils.CheckGormError(result.Error)
	}

	return deletedTodo, nil
}
