package domain

import (
	"simple-backend/internal/interactor/page"

	"gorm.io/gorm"
)

type AuthTable struct {
	gorm.Model
	Name     string `gorm:"uniqueIndex;column:name;type:varchar(50);not null;" json:"name"`
	Password string `gorm:"column:password;type:varchar(255);not null;" json:"password"`
}

type AuthQuery struct {
	Name *string `json:"name" form:"name"`
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
