package repository

import (
	"context"
	authModel "simple-backend/internal/domain/auth"
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

func (a *authRedisCacheRepo) GetUser(dest interface{}) error {
	return a.rdb.Get(a.ctx, "user", dest)
}

func (a *authRedisCacheRepo) SetUser(value interface{}) error {
	return a.rdb.Set(a.ctx, "user", value, 0)
}
