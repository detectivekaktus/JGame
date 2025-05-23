package router

import (
	"github.com/detectivekaktus/JGame/internal/handler"
	"github.com/detectivekaktus/JGame/internal/middleware"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.HandleFunc("/register", handler.RegisterUser).Methods("POST")
	r.HandleFunc("/login", handler.Login).Methods("POST")
	r.HandleFunc("/logout", handler.Logout).Methods("POST")

	// TODO: Handle it differently?
	r.HandleFunc("/users/{id:[0-9]+}", handler.GetUser).Methods("GET")

	users := r.PathPrefix("/users").Subrouter()
	users.Use(middleware.AuthMiddleware)
	
	users.HandleFunc("/me", handler.GetCurrentUser).Methods("GET")
	users.HandleFunc("/me", handler.PutCurrentUser).Methods("PUT")
	users.HandleFunc("/me", handler.PatchCurrentUser).Methods("PATCH")
	users.HandleFunc("/me", handler.DeleteCurrentUser).Methods("DELETE")

	return r
}
