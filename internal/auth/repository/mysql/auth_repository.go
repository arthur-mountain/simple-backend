package repository

import (
	authModel "simple-backend/internal/domain/auth"
	userModel "simple-backend/internal/domain/user"

	"gorm.io/gorm"
)

type authRepo struct {
	db *gorm.DB
}

func Init(db *gorm.DB) authModel.AuthRepoInterface {
	return &authRepo{db: db}
}

func (a *authRepo) WithTx(trx *gorm.DB) authModel.AuthRepoInterface {
	return &authRepo{db: trx}
}

func (a *authRepo) GetUser(input *userModel.UserTable) (*userModel.UserTable, error) {
	result := a.db.Model(&userModel.UserTable{}).First(input, "email = ?", input.Email)
	if result.Error != nil {
		return nil, result.Error
	}

	return input, nil
}
