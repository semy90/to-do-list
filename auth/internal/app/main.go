package app

import (
	"auth/internal/config"
	"auth/internal/database"
	"auth/internal/transport"
	"context"
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v5/pgxpool"
)

func StartApp() {
	cnfg, err := config.NewСonfig()
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.New(context.Background(), cnfg.DATABASE_URL)
	if err != nil {
		panic(err)
	}
	opt, err := redis.ParseURL(cnfg.REDIS_URL)
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(opt)
	store := database.NewUserStorage(pool)
	cache := database.NewRefreshCache(redisClient)
	auth := transport.Auth{Storage: *store, Cache: *cache}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", auth.Register)
	mux.HandleFunc("POST /login", auth.Login)
	mux.HandleFunc("DELETE /logout", auth.Logout)
	fmt.Println(cnfg.PORT)
	http.ListenAndServe(":"+cnfg.PORT, mux)

}
