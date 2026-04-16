package app

import (
	"context"
	"net/http"
	"task/internal/config"
	"task/internal/database"
	"task/internal/transport"

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
	store := database.NewTaskStorage(pool)
	taskcrud := transport.TaskCRUD{Storage: *store}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/{id}", taskcrud.GetTask)
	mux.HandleFunc("POST /tasks", taskcrud.PostTask)
	mux.HandleFunc("DELETE /tasks/{id}", taskcrud.DelTask)

	http.ListenAndServe(":"+cnfg.TASKPORT, mux)

}
