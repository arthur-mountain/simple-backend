package repository

import (
	"context"
	model "simple-backend/internal/domain/auth"
	"simple-backend/internal/utils/databases"
)

type authRedisCacheRepo struct {
	rdb *databases.MyRedis
	ctx *context.Context
}

func Init(rdb *databases.MyRedis) model.AuthCacheRepoInterface {
	return &authRedisCacheRepo{rdb: rdb}
}

func InitWithCtx(rdb *databases.MyRedis, ctx *context.Context) model.AuthCacheRepoInterface {
	return &authRedisCacheRepo{rdb: rdb, ctx: ctx}
}

func (a *authRedisCacheRepo) GetUser() (*model.UserTable, error) {
	user, err := a.rdb.Get(a.ctx, "user")

	if err != nil {
		return nil, err
	}

	return user.(*model.UserTable), nil
}

func (a *authRedisCacheRepo) SetUser(value interface{}) error {
	err := a.rdb.Set(a.ctx, "user", value, 0)

	if err != nil {
		return err
	}

	return nil
}
