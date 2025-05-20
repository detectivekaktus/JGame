package router

import (
	"github.com/detectivekaktus/JGame/internal/handler"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/users", handler.PostUser).Methods("POST")
	r.HandleFunc("/api/users/{id:[0-9]+}", handler.GetUser).Methods("GET")
	r.HandleFunc("/api/users/{id:[0-9]+}", handler.PutUser).Methods("PUT")
	r.HandleFunc("/api/users/{id:[0-9]+}", handler.PatchUser).Methods("PATCH")
	r.HandleFunc("/api/users/{id:[0-9]+}", handler.DeleteUser).Methods("DELETE")

	r.HandleFunc("/api/login", handler.Login).Methods("POST")

	return r
}
