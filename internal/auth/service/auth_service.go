package service

import (
	repo "simple-backend/internal/auth/repository/mysql"
	model "simple-backend/internal/domain/auth"

	"gorm.io/gorm"
)

type authService struct {
	Repository model.AuthRepoInterface
}

func Init(db *gorm.DB) model.AuthServiceInterface {
	return &authService{
		Repository: repo.Init(db),
	}
}

func (a *authService) GetUser(input *model.AuthTable) (*model.AuthTable, error) {
	user, err := a.Repository.GetUser(input)

	return user, err
}
