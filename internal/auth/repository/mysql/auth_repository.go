package repository

import (
	"errors"
	model "simple-backend/internal/domain/auth"
	authUtils "simple-backend/internal/utils/auth"

	"gorm.io/gorm"
)

type authRepo struct {
	db *gorm.DB
}

func Init(db *gorm.DB) model.AuthRepoInterface {
	return &authRepo{db: db}
}

func (a *authRepo) GetUser(input *model.AuthTable) (*model.AuthTable, error) {
	var user *model.AuthTable

	result := a.db.Model(&model.AuthTable{}).First(&user, "name = ?", input.Name)
	if result.Error != nil {
		return nil, result.Error
	}

	pwd := input.Password
	hashedPwd := user.Password
	isPassed := authUtils.IsPasswordPassed(hashedPwd, pwd)

	if !isPassed {
		return nil, errors.New("password is not correct")
	}

	return user, nil
}
