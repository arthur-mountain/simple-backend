package domain

import (
	user "simple-backend/internal/domain/user"
)

type AuthServiceInterface interface {
	Login(input *user.UserBody) (string, error)
	ForgotPassword(input *user.UserBody) error
}

type AuthRepoInterface interface {
	GetUser(user *user.UserTable) (*user.UserTable, error)
}

type AuthCacheRepoInterface interface {
}
