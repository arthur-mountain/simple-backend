package databases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type MyRedis struct {
	redisCtx context.Context
	redis    *redis.Client
}

func RedisInit(
	addr string,
	password string,
	db int,
	ctx context.Context,
) (*MyRedis, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	fmt.Println(pong)

	return &MyRedis{redisCtx: ctx, redis: redisClient}, nil
}

func (r *MyRedis) GetRedisClient() *redis.Client {
	return r.redis
}

func (r *MyRedis) GetRedisDefaultCtx() context.Context {
	return r.redisCtx
}

func (r *MyRedis) Get(ctx *context.Context, key string, dest interface{}) error {
	getterCtx := r.redisCtx
	if ctx != nil {
		getterCtx = *ctx
	}

	// Redis error
	val, err := r.redis.Get(getterCtx, key).Result()
	if err == redis.Nil {
		return errors.New("val does not exist")
	}
	if err != nil {
		return errors.New("client.Get failed:" + err.Error())
	}

	return json.Unmarshal([]byte(val), dest)
}

func (r *MyRedis) Set(ctx *context.Context, key string, value interface{}, expired time.Duration) error {
	setterCtx := r.redisCtx
	if ctx != nil {
		setterCtx = *ctx
	}

	// Json error
	jsonVal, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Redis error
	if err = r.redis.Set(setterCtx, key, jsonVal, expired).Err(); err != nil {
		return errors.New("client.Set failed" + err.Error())
	}

	return nil
}
