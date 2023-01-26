package domain

import (
	"simple-backend/internal/interactor/page"
	"simple-backend/internal/interactor/special"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserTable struct {
	special.Base
	Name       string `gorm:"index;column:name;type:varchar(50);not null;" json:"name"`
	Password   string `gorm:"column:password;type:varchar(255);not null;" json:"-"`
	IdentityId string `gorm:"uniqueIndex;column:identity_id;type:varchar(255);not null;" json:"-"`
}

type UserBody struct {
	Name            string `json:"name" form:"name" binding:"required"`
	Password        string `json:"password" form:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"omitempty,eqfield=Password"`
}

type UserQuery struct {
	Name      *string    `json:"name" form:"name"`
	CreatedAt *time.Time `json:"created_at" form:"created_at" time_format:"2006-01-02 15:04:05" time_utc:"1"`
	UpdatedAt *time.Time `json:"updated_at" form:"updated_at" time_format:"2006-01-02 15:04:05" time_utc:"1"`
	OrderBy   *string    `json:"order_by" form:"order_by" binding:"omitempty,oneof=asc desc"`
}

type UserQueries struct {
	UserQuery
	page.Pagination
}

type AuthServiceInterface interface {
	GetUser(input *UserBody) (*UserTable, error)
	CreateUser(input *UserBody) (*UserTable, error)
	UpdateUser(input *UserBody) error
	ForgotPassword(input *UserBody) error
}

type AuthRepoInterface interface {
	GetUser(user *UserTable) (*UserTable, error)
	CreateUser(user *UserTable) (*UserTable, error)
	UpdateUser(user *UserTable) error
}

type AuthCacheRepoInterface interface {
}

func (t *UserTable) TableName() string {
	return "users"
}

func (t *UserTable) BeforeCreate(tx *gorm.DB) error {
	t.IdentityId = uuid.New().String()

	return nil
}
