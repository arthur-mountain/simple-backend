package service

import (
	"context"
	"errors"
	repo "simple-backend/internal/auth/repository/mysql"
	cache "simple-backend/internal/auth/repository/redis"
	model "simple-backend/internal/domain/auth"
	authUtils "simple-backend/internal/utils/auth"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type authService struct {
	Repository model.AuthRepoInterface
	Cache      model.AuthCacheRepoInterface
}

func Init(db *gorm.DB, redis *redis.Client, redisCtx context.Context) model.AuthServiceInterface {
	return &authService{
		Repository: repo.Init(db),
		Cache:      cache.Init(redis, redisCtx),
	}
}

func (a *authService) GetUser(input *model.UserBody) (*model.UserTable, error) {
	user, err := a.Repository.GetUser(&model.UserTable{
		Name:     input.Name,
		Password: input.Password,
	})

	isPassed := authUtils.IsPasswordPassed(user.Password, input.Password)

	if !isPassed {
		return nil, errors.New("password is not correct")
	}

	return user, err
}

func (a *authService) CreateUser(input *model.UserBody) (*model.UserTable, error) {
	user, err := a.Repository.CreateUser(&model.UserTable{
		Name:     input.Name,
		Password: authUtils.GetPasswordHashed(input.Password),
	})

	return user, err
}

func (a *authService) UpdateUser(input *model.UserBody) error {
	err := a.Repository.UpdateUser(&model.UserTable{
		Name:     input.Name,
		Password: authUtils.GetPasswordHashed(input.Password),
	})

	return err
}

func (a *authService) ForgotPassword(input *model.UserBody) error {
	_, err := a.Repository.GetUser(&model.UserTable{Name: input.Name})

	if err != nil {
		return err
	}

	return nil
}
