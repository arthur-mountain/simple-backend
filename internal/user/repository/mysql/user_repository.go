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

func (a *userRepo) GetUsers() ([]*model.UserTable, error) {
	users := make([]*model.UserTable, 0)

	result := a.db.Model(&model.UserTable{}).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (a *userRepo) GetUser(user *model.UserTable) (*model.UserTable, error) {
	result := a.db.Model(&model.UserTable{}).First(&user, "name = ?", user.Name)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (a *userRepo) CreateUser(user *model.UserTable) (*model.UserTable, error) {
	result := a.db.Model(&model.UserTable{}).Create(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (a *userRepo) UpdateUser(user *model.UserTable) error {
	result := a.db.Model(&model.UserTable{}).Where("name = ?", user.Name).Update("password", user.Password)

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (a *userRepo) DeleteUser(user *model.UserTable) error {
	result := a.db.Model(&model.UserTable{}).Delete(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
