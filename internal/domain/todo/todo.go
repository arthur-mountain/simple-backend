package domain

import (
	userModel "simple-backend/internal/domain/user"
	"simple-backend/internal/interactor/page"
	"simple-backend/internal/interactor/special"
	errorUtils "simple-backend/internal/utils/error"
	"time"
)

// TodoTable "belongs to" UserTable
// UserTable "Has many" TodoTable
// Database struct
type TodoTable struct {
	special.Base
	Title       string              `gorm:"column:title;type:char(10);not null;" json:"title"`
	Description string              `gorm:"column:description;type:mediumtext;" json:"description" `
	IsCompleted uint                `gorm:"column:is_completed;type:tinyint(3);default:0" json:"is_completed"`
	CompletedAt *time.Time          `gorm:"column:completed_at;type:datetime;default:null" json:"completed_at,omitempty"`
	UserId      string              `gorm:"column:user_id;type:varchar(255);not null;" json:"-"`
	User        userModel.UserTable `gorm:"foreignKey:UserId;references:IdentityId" json:"-"`
}

type TodoCreate struct {
	UserId      string `json:"-"`
	Title       string `json:"title" form:"title" binding:"required"`
	Description string `json:"description" form:"description"`
}

type TodoUpdate struct {
	Id          uint   `json:"-"`
	UserId      string `json:"-"`
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
	IsCompleted *uint  `json:"is_completed" form:"is_completed"`
}

type TodoDelete struct {
	Id     uint   `json:"-"`
	UserId string `json:"-"`
}

// query string for search
type TodoQuery struct {
	UserId      string     `json:"-" form:"-"`
	Title       *string    `json:"title" form:"title"`
	IsCompleted *uint      `json:"is_completed" form:"is_completed"`
	CreatedAt   *time.Time `json:"created_at" form:"created_at" time_format:"2006-01-02 15:04:05" time_utc:"1"`
	UpdatedAt   *time.Time `json:"updated_at" form:"updated_at" time_format:"2006-01-02 15:04:05" time_utc:"1"`
	OrderBy     *string    `json:"order_by" form:"order_by" binding:"omitempty,oneof=asc desc"`
}

type TodoQueries struct {
	TodoQuery
	page.Pagination
}

type TodoRepoInterface interface {
	GetAllTodo(field *TodoQueries) (*int64, []*TodoTable, *errorUtils.CustomError)
	GetTodo(todo *TodoTable) (*TodoTable, *errorUtils.CustomError)
	CreateTodo(todo *TodoTable) *errorUtils.CustomError
	UpdateTodo(todo *TodoTable) (*TodoTable, *errorUtils.CustomError)
	UpdateTodoCompleted(todo *TodoTable) *errorUtils.CustomError
	DeleteTodo(todo *TodoTable) (*TodoTable, *errorUtils.CustomError)
}

type TodoServiceInterface interface {
	GetAllTodo(field *TodoQueries) (*int64, []*TodoTable, *errorUtils.CustomError)
	GetTodo(id uint, uid string) (*TodoTable, *errorUtils.CustomError)
	CreateTodo(input *TodoCreate) *errorUtils.CustomError
	UpdateTodo(input *TodoUpdate) (*TodoTable, *errorUtils.CustomError)
	UpdateTodoCompleted(input *TodoUpdate) *errorUtils.CustomError
	DeleteTodo(id uint, uid string) (*TodoTable, *errorUtils.CustomError)
}

func (t *TodoTable) TableName() string {
	return "todos"
}
