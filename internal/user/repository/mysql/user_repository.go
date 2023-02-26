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
	}, users)

	if err != nil {
		return nil, errorUtils.CheckRepoError(err)
	}

	return users, nil
}

func (u *userRepo) GetUser(user *model.UserTable) (*model.UserTable, *errorUtils.CustomError) {
	err := u.db.Execute(func(DB *gorm.DB) error {
		return DB.First(user).Error
	}, user)

	if err != nil {
		return nil, errorUtils.CheckRepoError(err)
	}

	return user, nil
}

func (u *userRepo) CreateUser(user *model.UserTable) (*model.UserTable, *errorUtils.CustomError) {
	err := u.db.Execute(func(DB *gorm.DB) error {
		return DB.Create(user).Error
	}, user)

	if err != nil {
		return nil, errorUtils.CheckRepoError(err)
	}

	return user, nil
}

func (u *userRepo) UpdateUser(user *model.UserTable) (*model.UserTable, *errorUtils.CustomError) {
	err := u.db.Execute(func(DB *gorm.DB) error {
		return DB.Select("name", "email").Updates(user).Error
	}, user)

	if err != nil {
		return nil, errorUtils.CheckRepoError(err)
	}

	return user, nil
}

func (u *userRepo) DeleteUser(user *model.UserTable) (*model.UserTable, *errorUtils.CustomError) {
	err := u.db.Execute(func(DB *gorm.DB) error {
		return DB.First(user).Delete(user).Error
	}, user)

	if err != nil {
		return nil, errorUtils.CheckRepoError(err)
	}

	return user, nil
}
