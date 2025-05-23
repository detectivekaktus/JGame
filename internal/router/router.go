package router

import (
	"net/http"

	"github.com/detectivekaktus/JGame/internal/handler"
	"github.com/detectivekaktus/JGame/internal/middleware"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	// /api/users
	r.Handle("/register", middleware.RequireBodyMiddleware(http.HandlerFunc(handler.RegisterUser))).Methods("POST")
	r.Handle("/login", middleware.RequireBodyMiddleware(http.HandlerFunc(handler.Login))).Methods("POST")
	r.Handle("/logout", middleware.RejectBodyMiddleware(http.HandlerFunc(handler.Logout))).Methods("POST")
	r.Handle("/users/{id:[0-9]+}", middleware.RejectBodyMiddleware(http.HandlerFunc(handler.GetUser))).Methods("GET")

	// /api/users
	users := r.PathPrefix("/users").Subrouter()
	users.Use(middleware.AuthMiddleware)

	users.Handle("/me", middleware.RejectBodyMiddleware(http.HandlerFunc(handler.GetCurrentUser))).Methods("GET")
	users.Handle("/me", middleware.RequireBodyMiddleware(http.HandlerFunc(handler.PutCurrentUser))).Methods("PUT")
	users.Handle("/me", middleware.RequireBodyMiddleware(http.HandlerFunc(handler.PatchCurrentUser))).Methods("PATCH")
	users.Handle("/me", middleware.RejectBodyMiddleware(http.HandlerFunc(handler.DeleteCurrentUser))).Methods("DELETE")

	return r
}
