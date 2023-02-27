package repository

import (
	model "simple-backend/internal/domain/todo"
	"simple-backend/internal/utils/databases"
	errorUtils "simple-backend/internal/utils/error"

	"gorm.io/gorm"
)

type todoRepo struct {
	db *databases.TMysql
}

func Init(db *databases.TMysql) model.TodoRepoInterface {
	return &todoRepo{db: db}
}

func (t *todoRepo) WithTrx(trx *gorm.DB) model.TodoRepoInterface {
	return &todoRepo{db: t.db.WithTrx(trx)}
}

func (t *todoRepo) GetAllTodo(field *model.TodoQueries) (*int64, []*model.TodoTable, *errorUtils.CustomError) {
	var totalCount int64
	allTodo := make([]*model.TodoTable, 0)

	err := t.db.Execute(func(DB *gorm.DB) error {
		// query := DB.Preload(clause.Associations)

		if field.Title != nil {
			DB.Where("title = ?", field.Title)
		}

		if field.IsCompleted != nil {
			DB.Where("is_completed = ?", field.IsCompleted)
		}

		if field.CreatedAt != nil {
			DB.Where("created_at >= ?", field.CreatedAt)
		}

		if field.UpdatedAt != nil {
			DB.Where("updated_at >= ?", field.UpdatedAt)
		}

		if field.OrderBy != nil {
			DB.Order("created_at " + *field.OrderBy)
		}

		if field.CurrentPage < 1 {
			field.CurrentPage = 1
		}

		if field.PerPage < 1 {
			field.PerPage = 15
		}

		return DB.Count(&totalCount).Offset(int((field.CurrentPage - 1) * field.PerPage)).Limit(int(field.PerPage)).Find(&allTodo).Error
	}, allTodo)

	if err != nil {
		return nil, nil, errorUtils.CheckRepoError(err)
	}

	return &totalCount, allTodo, nil
}

func (t *todoRepo) GetTodo(todo *model.TodoTable) (*model.TodoTable, *errorUtils.CustomError) {

	err := t.db.Execute(func(DB *gorm.DB) error {
		return DB.Where("user_id = ?", todo.UserId).First(todo).Error
	}, todo)

	if err != nil {
		return nil, errorUtils.CheckRepoError(err)
	}

	return todo, nil
}

func (t *todoRepo) CreateTodo(todo *model.TodoTable) *errorUtils.CustomError {
	err := t.db.Execute(func(DB *gorm.DB) error {
		return DB.Create(todo).Error
	}, todo)

	if err != nil {
		return errorUtils.CheckRepoError(err)
	}

	return nil
}

func (t *todoRepo) UpdateTodo(newTodo *model.TodoTable) (*model.TodoTable, *errorUtils.CustomError) {
	err := t.db.Execute(func(DB *gorm.DB) error {
		query := DB.Where("user_id = ?", newTodo.UserId)

		return query.Select("title", "description", "is_completed").Updates(newTodo).Error
	}, newTodo)

	if err != nil {
		return nil, errorUtils.CheckRepoError(err)
	}

	return newTodo, nil
}

func (t *todoRepo) UpdateTodoCompleted(updatedTodo *model.TodoTable) *errorUtils.CustomError {
	err := t.db.Execute(func(DB *gorm.DB) error {
		result := DB.Where("user_id = ?", updatedTodo.UserId).First(updatedTodo)

		if result.Error != nil {
			return result.Error
		}

		return result.Update("is_completed", 1).Error
	}, updatedTodo)

	if err != nil {
		return errorUtils.CheckRepoError(err)
	}

	return nil
}

func (t *todoRepo) DeleteTodo(deletedTodo *model.TodoTable) (*model.TodoTable, *errorUtils.CustomError) {
	err := t.db.Execute(func(DB *gorm.DB) error {
		query := DB.Where("user_id = ?", deletedTodo.UserId)

		return query.First(deletedTodo).Delete(deletedTodo).Error
	}, deletedTodo)

	if err != nil {
		return nil, errorUtils.CheckRepoError(err)
	}

	return deletedTodo, nil
}
