package repository

import (
	authModel "simple-backend/internal/domain/auth"
	userModel "simple-backend/internal/domain/user"
	errorUtils "simple-backend/internal/utils/error"

	"gorm.io/gorm"
)

type authRepo struct {
	db *gorm.DB
}

func Init(db *gorm.DB) authModel.AuthRepoInterface {
	return &authRepo{db: db}
}

func getModelQuery(db *gorm.DB) *gorm.DB {
	return db.Model(&userModel.UserTable{})
}

func (a *authRepo) WithTx(trx *gorm.DB) authModel.AuthRepoInterface {
	return &authRepo{db: trx}
}

func (a *authRepo) GetUser(input *userModel.UserTable) (*userModel.UserTable, *errorUtils.CustomError) {
	result := getModelQuery(a.db).First(input, "email = ?", input.Email)

	if result.Error != nil {
		return nil, errorUtils.CheckGormError(result.Error)
	}

	return input, nil
}
