package repository

import (
	model "simple-backend/internal/domain/auth"

	"gorm.io/gorm"
)

type authRepo struct {
	db *gorm.DB
}

func Init(db *gorm.DB) model.AuthRepoInterface {
	return &authRepo{db: db}
}

func (a *authRepo) GetUser(input *model.UserTable) (*model.UserTable, error) {
	var user model.UserTable

	result := a.db.Model(&model.UserTable{}).First(&user, "name = ?", input.Name)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (a *authRepo) CreateUser(input *model.UserTable) (*model.UserTable, error) {
	result := a.db.Model(&model.UserTable{}).Create(input)

	if result.Error != nil {
		return nil, result.Error
	}

	return input, nil
}

func (a *authRepo) UpdateUser(input *model.UserTable) error {
	result := a.db.Model(&model.UserTable{}).Where("name = ?", input.Name).Update("password", input.Password)

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
