package repository

import (
	model "simple-backend/internal/domain/user"
	"simple-backend/internal/utils/databases"
	errorUtils "simple-backend/internal/utils/error"

	"gorm.io/gorm"
)

type userRepo struct {
	db *databases.TMysql
}

func Init(db *databases.TMysql) model.UserRepoInterface {
	return &userRepo{db: db}
}

func (u *userRepo) WithTx(trx *gorm.DB) model.UserRepoInterface {
	return &userRepo{db: u.db.WithTrx(trx)}
}

func (u *userRepo) GetUsers() ([]*model.UserTable, *errorUtils.CustomError) {
	users := make([]*model.UserTable, 0)

	err := u.db.Execute(func(DB *gorm.DB) error {
		return DB.Find(&users).Error
	}, &model.UserTable{})
	if err != nil {
		return nil, errorUtils.CheckRepoError(err)
	}

	return users, nil
}

func (u *userRepo) GetUser(user *model.UserTable) (*model.UserTable, *errorUtils.CustomError) {

	err := u.db.Execute(func(DB *gorm.DB) error {
		return DB.First(user).Error
	}, &model.UserTable{})

	if err != nil {
		return nil, errorUtils.CheckRepoError(err)
	}

	return user, nil
}

func (u *userRepo) CreateUser(user *model.UserTable) (*model.UserTable, *errorUtils.CustomError) {
	err := u.db.Execute(func(DB *gorm.DB) error {
		return DB.Create(user).Error
	}, &model.UserTable{})

	if err != nil {
		return nil, errorUtils.CheckRepoError(err)
	}

	return user, nil
}

func (u *userRepo) UpdateUser(user *model.UserTable) *errorUtils.CustomError {
	err := u.db.Execute(func(DB *gorm.DB) error {
		return DB.Where("id = ?", user.Id).Updates(user).Error
	}, &model.UserTable{})

	// TODO: May should return updated user data
	if err != nil {
		return errorUtils.CheckRepoError(err)
	}

	return nil
}

func (u *userRepo) DeleteUser(user *model.UserTable) *errorUtils.CustomError {
	err := u.db.Execute(func(DB *gorm.DB) error {
		return DB.Delete(user).Error
	}, &model.UserTable{})

	// TODO: May should return deleted user data
	if err != nil {
		return errorUtils.CheckRepoError(err)
	}

	return nil
}
