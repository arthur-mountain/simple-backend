package connect

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr     string
	Password string
	Db       int
}

func (r *RedisConfig) NewClient(ctx context.Context) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password, // no password set
		DB:       r.Db,       // use default DB
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	fmt.Println(pong)
	return redisClient, nil
}

// type MyRedis struct {
// 	redisCtx context.Context
// 	redis    *redis.Client
// }

// func RedisInit(
// 	addr string,
// 	password string,
// 	db int,
// 	ctx context.Context,
// ) (*MyRedis, error) {
// 	redisClient := redis.NewClient(&redis.Options{
// 		Addr:     addr,
// 		Password: password, // no password set
// 		DB:       db,       // use default DB
// 	})

// 	pong, err := redisClient.Ping(ctx).Result()
// 	if err != nil {
// 		return nil, err
// 	}

// 	fmt.Println(pong)

// 	return &MyRedis{redisCtx: ctx, redis: redisClient}, nil
// }

// func (r *MyRedis) Get(ctx context.Context, key string) (*string, error) {
// 	getterCtx := r.redisCtx
// 	if ctx != nil {
// 		getterCtx = ctx
// 	}

// 	val, err := r.redis.Get(getterCtx, key).Result()
// 	if err == redis.Nil {
// 		return nil, errors.New("val does not exist")
// 	} else if err != nil {
// 		return nil, errors.New("redis get failed:" + err.Error())
// 	} else {
// 		return &val, nil
// 	}
// }

// func (r *MyRedis) Set(ctx context.Context, key string, value interface{}, expired time.Duration) error {
// 	setterCtx := r.redisCtx
// 	if ctx != nil {
// 		setterCtx = ctx
// 	}

// 	if err := r.redis.Set(setterCtx, key, value, expired).Err(); err != nil {
// 		return errors.New("client.Set failed" + err.Error())
// 	}

// 	return nil
// }
