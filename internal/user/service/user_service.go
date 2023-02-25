package service

import (
	model "simple-backend/internal/domain/user"
	repo "simple-backend/internal/user/repository/mysql"
	redisCache "simple-backend/internal/user/repository/redis"
	authUtils "simple-backend/internal/utils/auth"
	"simple-backend/internal/utils/databases"
	errorUtils "simple-backend/internal/utils/error"

	"gorm.io/gorm"
)

type userService struct {
	Repository model.UserRepoInterface
	Cache      model.UserCacheRepoInterface
}

func Init(db *gorm.DB, redis *databases.MyRedis) model.UserServiceInterface {
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

func (a *userService) UpdateUser(id uint, input *model.UserUpdate) *errorUtils.CustomError {
	user := model.UserTable{}
	user.Id = id

	// find user is exists
	updatedUser, err := a.Repository.GetUser(&user)
	if err != nil {
		return err
	}

	if input.Name != "" {
		updatedUser.Name = input.Name
	}

	if input.Email != "" {
		updatedUser.Email = input.Email
	}

	return a.Repository.UpdateUser(updatedUser)
}

func (a *userService) DeleteUser(id uint) *errorUtils.CustomError {
	deletedUser := &model.UserTable{}
	deletedUser.Id = id

	return a.Repository.DeleteUser(deletedUser)
}
