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

func (a *authRepo) GetUser(input *model.AuthTable) (*model.AuthTable, error) {
	var user model.AuthTable

	result := a.db.Model(&model.AuthTable{}).First(&user, "name = ?", input.Name)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (a *authRepo) CreateUser(input *model.AuthTable) (*model.AuthTable, error) {
	result := a.db.Model(&model.AuthTable{}).Create(input)

	if result.Error != nil {
		return nil, result.Error
	}

	return input, nil
}

func (a *authRepo) UpdateUser(input *model.AuthTable) error {
	result := a.db.Model(&model.AuthTable{}).Where("name = ?", input.Name).Update("password", input.Password)

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
