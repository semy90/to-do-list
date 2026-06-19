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

/*
todo:
pagination
swagger
kafka
*/

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
	mux.HandleFunc("GET /tasks/get/{id}", taskcrud.GetTask)
	mux.HandleFunc("GET /tasks/range", taskcrud.GetTaskFromTo)
	mux.HandleFunc("POST /tasks/post", taskcrud.PostTask)
	mux.HandleFunc("DELETE /tasks/delete/{id}", taskcrud.DelTask)
	mux.HandleFunc("PUT /tasks/update", taskcrud.UpdateTask)
	authMux := taskcrud.CheckAuth(mux)
	http.ListenAndServe(":"+cnfg.PORT, authMux)
}
