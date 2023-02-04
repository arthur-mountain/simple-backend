package domain

import (
	user "simple-backend/internal/domain/user"
)

type LoginBody struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type ResetPasswordBody struct {
	// VerifyCode      string `json:"verify_code" form:"verify_code" binding:"required"`
	Password        string `json:"password" form:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"required,eqfield=Password"`
}

type AuthServiceInterface interface {
	Login(input *LoginBody) (string, error)
	ForgotPassword(input *LoginBody) error
}

type AuthRepoInterface interface {
	GetUser(user *user.UserTable) (*user.UserTable, error)
}

type AuthCacheRepoInterface interface {
	GetUser(dest interface{}) error
	SetUser(value interface{}) error
}
