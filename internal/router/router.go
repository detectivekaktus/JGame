package router

import (
	"github.com/detectivekaktus/JGame/internal/handler"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/users", handler.PostUser).Methods("POST")

	return r
}
