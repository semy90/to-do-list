package app

import (
	"auth/internal/config"
	"auth/internal/database"
	"auth/internal/service"
	"context"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func StartApp() {
	cnfg, err := config.NewСonfig()
	if err != nil {
		panic(err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	ctx := context.WithValue(context.Background(), "logger", logger)

	pool, err := pgxpool.New(context.Background(), cnfg.DATABASE_URL)
	if err != nil {
		panic(err)
	}
	opt, err := redis.ParseURL(cnfg.REDIS_URL)
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(opt)
	_, err = redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
	store := database.NewUserStorage(pool)
	cache := database.NewRefreshCache(redisClient)
	auth := service.Auth{Storage: *store, Cache: *cache, Ctx: ctx}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", auth.Register)
	mux.HandleFunc("POST /login", auth.Login)
	mux.HandleFunc("DELETE /logout", auth.Logout)
	mux.HandleFunc("GET /validate", auth.CheckAuth)
	logger.Info("auth service start succesfully!")
	if err := http.ListenAndServe(":"+cnfg.PORT, mux); err != nil {
	}

}
