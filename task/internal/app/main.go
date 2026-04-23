package app

import (
	"context"
	"net/http"
	"task/internal/config"
	"task/internal/database"
	"task/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
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
	logger, err := zap.NewProduction()
	ctx := context.WithValue(context.Background(), "logger", logger)

	store := database.NewTaskStorage(pool)
	taskcrud := service.TaskCRUD{Storage: *store, Ctx: ctx}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/{id}", taskcrud.GetTask)
	mux.HandleFunc("POST /tasks", taskcrud.PostTask)
	mux.HandleFunc("DELETE /tasks/{id}", taskcrud.DelTask)

	http.ListenAndServe(":"+cnfg.TASKPORT, mux)

}
