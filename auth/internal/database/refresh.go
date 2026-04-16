package database

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func NewRefreshCache(rds *redis.Client) *RefreshCache {
	return &RefreshCache{rds: rds}
}

type RefreshCache struct {
	rds *redis.Client
}

func (rs *RefreshCache) Get(token string) (int, error) {
	id, err := rs.rds.Get(token).Result()
	if err != nil {
		return -1, err
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return -1, err
	}
	return idInt, nil
}

func (rs *RefreshCache) Add(id int, token string) error {
	if err := rs.rds.Set(token, fmt.Sprintf("user:%d", id), time.Hour*24*7).Err(); err != nil {
		return err
	}
	return nil
}

func (rs *RefreshCache) Del(token string) error {
	if err := rs.rds.Del(token).Err(); err != nil {
		return err
	}
	return nil
}
