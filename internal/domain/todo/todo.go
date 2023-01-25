package domain

import (
	authModel "simple-backend/internal/domain/auth"
	"simple-backend/internal/interactor/page"
	"simple-backend/internal/interactor/special"
	"time"
)

// TodoTable "belongs to" UserTable
// UserTable "Has many" TodoTable
// Database struct
type TodoTable struct {
	special.Base
	Title       string              `gorm:"column:title;type:char(10);not null;" json:"title" form:"title" binding:"required"`
	Description string              `gorm:"column:description;type:mediumtext;" json:"description" form:"description"`
	IsCompleted uint                `gorm:"column:is_completed;type:tinyint(3);default:0" json:"is_completed" form:"is_completed" binding:"oneof=1 0"`
	CompletedAt *time.Time          `gorm:"column:completed_at;type:datetime;default:null" json:"completed_at,omitempty"`
	UserId      string              `gorm:"column:user_id;type:varchar(255);not null;" json:"-"`
	User        authModel.UserTable `gorm:"foreignKey:UserId;references:IdentityId" json:"-"`
}

// query string for search
type TodoQuery struct {
	UserId      string     `json:"-" form:"-"`
	Title       *string    `json:"title" form:"title"`
	IsCompleted *string    `json:"is_completed" form:"is_completed"`
	CreatedAt   *time.Time `json:"created_at" form:"created_at" time_format:"2006-01-02 15:04:05" time_utc:"1"`
	UpdatedAt   *time.Time `json:"updated_at" form:"updated_at" time_format:"2006-01-02 15:04:05" time_utc:"1"`
	OrderBy     *string    `json:"order_by" form:"order_by" binding:"omitempty,oneof=asc desc"`
}

type TodoQueries struct {
	TodoQuery
	page.Pagination
}

type TodoRepoInterface interface {
	GetAllTodo(field *TodoQueries) (*int64, []*TodoTable, error)
	GetTodo(todo *TodoTable) (*TodoTable, error)
	CreateTodo(todo *TodoTable) error
	UpdateTodo(id int, newTodo *TodoTable) (*TodoTable, error)
	UpdateTodoCompleted(updatedTodo *TodoTable) (*int64, error)
	DeleteTodo(updatedTodo *TodoTable) (*TodoTable, error)
}

type TodoServiceInterface interface {
	GetAllTodo(field *TodoQueries) (*int64, []*TodoTable, error)
	GetTodo(todo *TodoTable) (*TodoTable, error)
	CreateTodo(input *TodoTable) error
	UpdateTodo(id int, newTodo *TodoTable) (*TodoTable, error)
	UpdateTodoCompleted(updatedTodo *TodoTable) (*int64, error)
	DeleteTodo(updatedTodo *TodoTable) (*TodoTable, error)
}

func (t *TodoTable) TableName() string {
	return "todos"
}
