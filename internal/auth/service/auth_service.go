package service

import (
	"errors"
	repo "simple-backend/internal/auth/repository/mysql"
	redisCache "simple-backend/internal/auth/repository/redis"
	authModel "simple-backend/internal/domain/auth"
	userModel "simple-backend/internal/domain/user"
	authUtils "simple-backend/internal/utils/auth"
	"simple-backend/internal/utils/databases"

	"gorm.io/gorm"
)

type authService struct {
	Repository authModel.AuthRepoInterface
	Cache      authModel.AuthCacheRepoInterface
}

func Init(db *gorm.DB, redis *databases.MyRedis) authModel.AuthServiceInterface {
	return &authService{
		Repository: repo.Init(db),
		Cache:      redisCache.Init(redis),
	}
}

func (a *authService) Login(input *userModel.UserBody) (string, error) {
	user, err := a.Repository.GetUser(&userModel.UserTable{
		Name:     input.Name,
		Password: input.Password,
	})

	if err != nil {
		return "", err
	}

	isPassed := authUtils.IsPasswordPassed(user.Password, input.Password)

	if !isPassed {
		return "", errors.New("password is not correct")
	}

	token, err := authUtils.GenerateToken(map[string]interface{}{
		"uid":      user.IdentityId,
		"userName": user.Name,
	})

	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *authService) ForgotPassword(input *userModel.UserBody) error {
	_, err := a.Repository.GetUser(&userModel.UserTable{Name: input.Name})

	if err != nil {
		return err
	}

	return nil
}
