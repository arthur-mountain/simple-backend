package repository

import (
	authModel "simple-backend/internal/domain/auth"
	userModel "simple-backend/internal/domain/user"
	"simple-backend/internal/utils/databases"
	errorUtils "simple-backend/internal/utils/error"

	"gorm.io/gorm"
)

type authRepo struct {
	db *databases.TMysql
}

func Init(db *databases.TMysql) authModel.AuthRepoInterface {
	return &authRepo{db: db}
}

func (a *authRepo) WithTx(trx *gorm.DB) authModel.AuthRepoInterface {
	return &authRepo{db: a.db.WithTrx(trx)}
}

func (a *authRepo) GetUser(user *userModel.UserTable) (*userModel.UserTable, *errorUtils.CustomError) {
	err := a.db.Execute(func(DB *gorm.DB) error {
		return DB.Where("email = ?", user.Email).First(user).Error
	}, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}
