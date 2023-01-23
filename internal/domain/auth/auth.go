package domain

import (
	"simple-backend/internal/interactor/page"
	"time"

	"gorm.io/gorm"
)

type AuthTable struct {
	Id        uint           `gorm:"column:id;primaryKey" json:"id" form:"id"`
	CreatedAt time.Time      `gorm:"column:created_at;" json:"created_at" form:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;" json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;" json:"deleted_at" form:"deleted_at"`
	Name      string         `gorm:"uniqueIndex;column:name;type:varchar(50);not null;" json:"name"`
	Password  string         `gorm:"column:password;type:varchar(255);not null;" json:"password"`
}

type AuthQuery struct {
	Name      *string    `json:"name" form:"name"`
	CreatedAt *time.Time `json:"created_at" form:"created_at" time_format:"2006-01-02 15:04:05" time_utc:"1"`
	UpdatedAt *time.Time `json:"updated_at" form:"updated_at" time_format:"2006-01-02 15:04:05" time_utc:"1"`
	OrderBy   *string    `json:"order_by" form:"order_by"`
}

type AuthQueries struct {
	AuthQuery
	page.Pagination
}

type AuthServiceInterface interface {
	GetUser(input *AuthTable) (*AuthTable, error)
}

type AuthRepoInterface interface {
	GetUser(user *AuthTable) (*AuthTable, error)
}

func (t *AuthTable) TableName() string {
	return "users"
}
