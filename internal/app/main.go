package app

import (
	"net/http"
	"to-do-list/internal/database"
	"to-do-list/internal/transport"
)

func StartApp() {

	store := database.New()
	taskcrud := transport.TaskCRUD{Storage: *store}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/{id}", taskcrud.GetTask)
	mux.HandleFunc("POST /tasks", taskcrud.PostTask)
	mux.HandleFunc("DELETE /tasks/{id}", taskcrud.DelTask)

	http.ListenAndServe(":8080", mux)

}
