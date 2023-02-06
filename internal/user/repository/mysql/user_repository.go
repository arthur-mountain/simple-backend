package repository

import (
	model "simple-backend/internal/domain/user"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func Init(db *gorm.DB) model.UserRepoInterface {
	return &userRepo{db: db}
}

func (u *userRepo) WithTx(trx *gorm.DB) model.UserRepoInterface {
	return &userRepo{db: trx}
}

func (u *userRepo) GetUsers() ([]*model.UserTable, error) {
	users := make([]*model.UserTable, 0)

	result := u.db.Model(&model.UserTable{}).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (u *userRepo) GetUser(user *model.UserTable) (*model.UserTable, error) {
	result := u.db.Model(&model.UserTable{}).First(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (u *userRepo) CreateUser(user *model.UserTable) (*model.UserTable, error) {
	result := u.db.Model(&model.UserTable{}).Create(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (u *userRepo) UpdateUser(user *model.UserTable) error {
	result := u.db.Model(&model.UserTable{}).Where("id = ?", user.Id).Updates(user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (u *userRepo) DeleteUser(user *model.UserTable) error {
	result := u.db.Model(&model.UserTable{}).Delete(user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
