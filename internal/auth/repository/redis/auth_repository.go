package repository

import (
	"context"
	model "simple-backend/internal/domain/auth"

	"github.com/redis/go-redis/v9"
)

type authRepo struct {
	rdb *redis.Client
	ctx context.Context
}

func Init(rdb *redis.Client, ctx context.Context) model.AuthCacheRepoInterface {
	return &authRepo{rdb: rdb, ctx: ctx}
}
