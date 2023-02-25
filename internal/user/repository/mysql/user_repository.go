package repository

import (
	model "simple-backend/internal/domain/user"
	errorUtils "simple-backend/internal/utils/error"

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

func (u *userRepo) GetUsers() ([]*model.UserTable, *errorUtils.CustomError) {
	users := make([]*model.UserTable, 0)

	result := u.db.Model(&model.UserTable{}).Find(&users)
	if result.Error != nil {
		return nil, errorUtils.CheckGormError(result.Error)
	}

	return users, nil
}

func (u *userRepo) GetUser(user *model.UserTable) (*model.UserTable, *errorUtils.CustomError) {
	result := u.db.Model(&model.UserTable{}).First(user)

	if result.Error != nil {
		return nil, errorUtils.CheckGormError(result.Error)
	}

	return user, nil
}

func (u *userRepo) CreateUser(user *model.UserTable) (*model.UserTable, *errorUtils.CustomError) {
	result := u.db.Model(&model.UserTable{}).Create(user)

	if result.Error != nil {
		return nil, errorUtils.CheckGormError(result.Error)
	}

	return user, nil
}

func (u *userRepo) UpdateUser(user *model.UserTable) *errorUtils.CustomError {
	result := u.db.Model(&model.UserTable{}).Where("id = ?", user.Id).Updates(user)

	// TODO: May should return updated user data
	if result.Error != nil {
		return errorUtils.CheckGormError(result.Error)
	}

	return nil
}

func (u *userRepo) DeleteUser(user *model.UserTable) *errorUtils.CustomError {
	result := u.db.Model(&model.UserTable{}).Delete(user)

	// TODO: May should return deleted user data
	if result.Error != nil {
		return errorUtils.CheckGormError(result.Error)
	}

	return nil
}
