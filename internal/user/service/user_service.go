package service

import (
	model "simple-backend/internal/domain/user"
	repo "simple-backend/internal/user/repository/mysql"
	redisCache "simple-backend/internal/user/repository/redis"
	authUtils "simple-backend/internal/utils/auth"
	"simple-backend/internal/utils/databases"

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

func (a *userService) GetUsers() ([]*model.UserTable, error) {
	users, err := a.Repository.GetUsers()

	return users, err
}

func (a *userService) GetUser(id uint) (*model.UserTable, error) {
	var err error
	user := &model.UserTable{}
	user.Id = id

	user, err = a.Repository.GetUser(user)

	return user, err
}

func (a *userService) CreateUser(input *model.UserBody) (*model.UserTable, error) {
	user, err := a.Repository.CreateUser(&model.UserTable{
		Name:     input.Name,
		Password: authUtils.GetPasswordHashed(input.Password),
	})

	return user, err
}

func (a *userService) UpdateUser(id uint, input *model.UserBody) error {
	updatedUser := &model.UserTable{
		Name:     input.Name,
		Password: authUtils.GetPasswordHashed(input.Password),
	}
	updatedUser.Id = id

	err := a.Repository.UpdateUser(updatedUser)

	return err
}

func (a *userService) DeleteUser(id uint) error {
	deletedUser := &model.UserTable{}
	deletedUser.Id = id

	err := a.Repository.DeleteUser(deletedUser)

	return err
}
