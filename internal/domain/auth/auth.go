package domain

import (
	"simple-backend/internal/interactor/page"
	"simple-backend/internal/interactor/special"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthTable struct {
	special.Base
	Name       string `gorm:"uniqueIndex;column:name;type:varchar(50);not null;" json:"name"`
	Password   string `gorm:"column:password;type:varchar(255);not null;" json:"-"`
	IdentityId string `gorm:"uniqueIndex;column:identity_id;type:varchar(255);not null;" json:"-"`
}

type AuthBody struct {
	Name            string `json:"name" form:"name" binding:"required"`
	Password        string `json:"password" form:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"omitempty,eqfield=Password"`
}

type AuthQuery struct {
	Name      *string    `json:"name" form:"name"`
	CreatedAt *time.Time `json:"created_at" form:"created_at" time_format:"2006-01-02 15:04:05" time_utc:"1"`
	UpdatedAt *time.Time `json:"updated_at" form:"updated_at" time_format:"2006-01-02 15:04:05" time_utc:"1"`
	OrderBy   *string    `json:"order_by" form:"order_by" binding:"omitempty,oneof=asc desc"`
}

type AuthQueries struct {
	AuthQuery
	page.Pagination
}

type AuthServiceInterface interface {
	GetUser(input *AuthBody) (*AuthTable, error)
	CreateUser(input *AuthBody) (*AuthTable, error)
	UpdateUser(input *AuthBody) error
	ForgotPassword(input *AuthBody) error
}

type AuthRepoInterface interface {
	GetUser(user *AuthTable) (*AuthTable, error)
	CreateUser(user *AuthTable) (*AuthTable, error)
	UpdateUser(user *AuthTable) error
}

func (t *AuthTable) TableName() string {
	return "users"
}

func (t *AuthTable) BeforeCreate(tx *gorm.DB) error {
	t.IdentityId = uuid.New().String()

	return nil
}
