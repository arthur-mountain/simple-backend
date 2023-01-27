package repository

import (
	"context"
	authModel "simple-backend/internal/domain/auth"
	userModel "simple-backend/internal/domain/user"
	"simple-backend/internal/utils/databases"
)

type authRedisCacheRepo struct {
	rdb *databases.MyRedis
	ctx *context.Context
}

func Init(rdb *databases.MyRedis) authModel.AuthCacheRepoInterface {
	return &authRedisCacheRepo{rdb: rdb}
}

func InitWithCtx(rdb *databases.MyRedis, ctx *context.Context) authModel.AuthCacheRepoInterface {
	return &authRedisCacheRepo{rdb: rdb, ctx: ctx}
}

func (a *authRedisCacheRepo) GetUser() (*userModel.UserTable, error) {
	user, err := a.rdb.Get(a.ctx, "user")

	if err != nil {
		return nil, err
	}

	return user.(*userModel.UserTable), nil
}

func (a *authRedisCacheRepo) SetUser(value interface{}) error {
	err := a.rdb.Set(a.ctx, "user", value, 0)

	if err != nil {
		return err
	}

	return nil
}
