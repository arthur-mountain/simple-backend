package service

import (
	model "simple-backend/internal/domain/user"
	repo "simple-backend/internal/user/repository/mysql"
	redisCache "simple-backend/internal/user/repository/redis"
	authUtils "simple-backend/internal/utils/auth"
	"simple-backend/internal/utils/databases"
	errorUtils "simple-backend/internal/utils/error"
)

type userService struct {
	Repository model.UserRepoInterface
	Cache      model.UserCacheRepoInterface
}

func Init(db *databases.TMysql, redis *databases.MyRedis) model.UserServiceInterface {
	return &userService{
		Repository: repo.Init(db),
		Cache:      redisCache.Init(redis),
	}
}

func (a *userService) GetUsers() ([]*model.UserTable, *errorUtils.CustomError) {
	return a.Repository.GetUsers()
}

func (a *userService) GetUser(id uint) (*model.UserTable, *errorUtils.CustomError) {
	user := &model.UserTable{}
	user.Id = id

	return a.Repository.GetUser(user)
}

func (a *userService) CreateUser(input *model.UserCreate) (*model.UserTable, *errorUtils.CustomError) {
	return a.Repository.CreateUser(&model.UserTable{
		Name:     input.Name,
		Email:    input.Email,
		Password: authUtils.GetPasswordHashed(input.Password),
	})
}

func (a *userService) UpdateUser(id uint, input *model.UserUpdate) (*model.UserTable, *errorUtils.CustomError) {
	var err *errorUtils.CustomError
	user := new(model.UserTable)
	user.Id = id

	// find user is exists
	if user, err = a.Repository.GetUser(user); err != nil {
		return nil, err
	}

	if input.Name != "" {
		user.Name = input.Name
	}

	if input.Email != "" {
		user.Email = input.Email
	}

	return a.Repository.UpdateUser(user)
}

func (a *userService) DeleteUser(id uint) (*model.UserTable, *errorUtils.CustomError) {
	deletedUser := &model.UserTable{}
	deletedUser.Id = id

	return a.Repository.DeleteUser(deletedUser)
}
