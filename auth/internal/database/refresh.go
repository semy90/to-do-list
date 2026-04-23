package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func NewRefreshCache(rds *redis.Client) *RefreshCache {
	return &RefreshCache{rds: rds}
}

type RefreshCache struct {
	rds *redis.Client
}

func (rs *RefreshCache) Get(ctx context.Context, id int) (string, error) {
	token, err := rs.rds.Get(fmt.Sprintf("user:%d", id)).Result()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (rs *RefreshCache) Add(context context.Context, id int, token string) error {
	if err := rs.rds.Set(fmt.Sprintf("user:%d", id), token, time.Hour*24*7).Err(); err != nil {
		return err
	}
	return nil
}

func (rs *RefreshCache) Del(context context.Context, id int) error {
	if err := rs.rds.Del(fmt.Sprintf("user:%d", id)).Err(); err != nil {
		return err
	}
	return nil
}
