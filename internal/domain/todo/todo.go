package domain

import (
	"simple-backend/internal/interactor/page"
	"time"

	"gorm.io/gorm"
)

// Database struct
type TodoTable struct {
	Id          uint           `gorm:"column:id;primaryKey" json:"id" form:"id"`
	CreatedAt   time.Time      `gorm:"column:created_at;" json:"created_at" form:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;" json:"updated_at" form:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;" json:"deleted_at" form:"deleted_at"`
	Title       string         `gorm:"column:title;type:char(10);not null;" json:"title" form:"title" binding:"required"`
	Description string         `gorm:"column:description;type:mediumtext;" json:"description" form:"description"`
	IsCompleted uint           `gorm:"column:is_completed;type:tinyint(3);default:0" json:"is_completed" form:"is_completed" binding:"oneof=1 0"`
	CompletedAt *time.Time     `gorm:"column:completed_at;type:datetime;default:null" json:"completed_at,omitempty"`
}

// query string for search
type TodoQuery struct {
	Title       *string    `json:"title" form:"title"`
	IsCompleted *string    `json:"is_completed" form:"is_completed"`
	CreatedAt   *time.Time `json:"created_at" form:"created_at" time_format:"2006-01-02 15:04:05" time_utc:"1"`
	UpdatedAt   *time.Time `json:"updated_at" form:"updated_at" time_format:"2006-01-02 15:04:05" time_utc:"1"`
	OrderBy     *string    `json:"order_by" form:"order_by"`
}

type TodoQueries struct {
	TodoQuery
	page.Pagination
}

type TodoRepoInterface interface {
	GetAllTodo(field *TodoQueries) (*int64, []*TodoTable, error)
	GetTodo(id int) (*TodoTable, error)
	CreateTodo(todo *TodoTable) error
	UpdateTodo(id int, newTodo *TodoTable) (*TodoTable, error)
	UpdateTodoCompleted(id int) (*int64, error)
	DeleteTodo(id int) (*TodoTable, error)
}

type TodoServiceInterface interface {
	GetAllTodo(field *TodoQueries) (*int64, []*TodoTable, error)
	GetTodo(id int) (*TodoTable, error)
	CreateTodo(input *TodoTable) error
	UpdateTodo(id int, newTodo *TodoTable) (*TodoTable, error)
	UpdateTodoCompleted(id int) (*int64, error)
	DeleteTodo(id int) (*TodoTable, error)
}

func (t *TodoTable) TableName() string {
	return "todos"
}
