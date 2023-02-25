package service

import (
	"errors"
	"net/http"
	repo "simple-backend/internal/auth/repository/mysql"
	redisCache "simple-backend/internal/auth/repository/redis"
	authModel "simple-backend/internal/domain/auth"
	userModel "simple-backend/internal/domain/user"
	authUtils "simple-backend/internal/utils/auth"
	"simple-backend/internal/utils/databases"
	errorUtils "simple-backend/internal/utils/error"
)

type authService struct {
	Repository authModel.AuthRepoInterface
	Cache      authModel.AuthCacheRepoInterface
}

func Init(db *databases.TMysql, redis *databases.MyRedis) authModel.AuthServiceInterface {
	return &authService{
		Repository: repo.Init(db),
		Cache:      redisCache.Init(redis),
	}
}

func (a *authService) Login(input *authModel.LoginBody) (*string, *errorUtils.CustomError) {
	user, err := a.Repository.GetUser(&userModel.UserTable{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		return nil, err
	}

	isPassed := authUtils.IsPasswordPassed(user.Password, input.Password)
	if !isPassed {
		return nil, errorUtils.NewCustomError(
			errors.New("password is not correct"),
			http.StatusUnauthorized,
		)
	}

	token, tokenErr := authUtils.GenerateToken(map[string]interface{}{
		"uid":      user.IdentityId,
		"userName": user.Name,
	})

	if err != nil {
		return nil, errorUtils.NewCustomError(
			tokenErr,
			http.StatusUnauthorized,
		)
	}

	return &token, nil
}

func (a *authService) ForgotPassword(input *authModel.LoginBody) *errorUtils.CustomError {
	_, err := a.Repository.GetUser(&userModel.UserTable{Email: input.Email})

	return err
}
