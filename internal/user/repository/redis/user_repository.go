package repository

import (
	"context"
	model "simple-backend/internal/domain/user"
	"simple-backend/internal/utils/databases"
)

type userRedisCacheRepo struct {
	rdb *databases.MyRedis
	ctx *context.Context
}

func Init(rdb *databases.MyRedis) model.UserCacheRepoInterface {
	return &userRedisCacheRepo{rdb: rdb}
}

func InitWithCtx(rdb *databases.MyRedis, ctx *context.Context) model.UserCacheRepoInterface {
	return &userRedisCacheRepo{rdb: rdb, ctx: ctx}
}

func (a *userRedisCacheRepo) GetUser() (*model.UserTable, error) {
	user, err := a.rdb.Get(a.ctx, "user")

	if err != nil {
		return nil, err
	}

	return user.(*model.UserTable), nil
}

func (a *userRedisCacheRepo) SetUser(value interface{}) error {
	err := a.rdb.Set(a.ctx, "user", value, 0)

	if err != nil {
		return err
	}

	return nil
}
