package router

import (
	"net/http"

	"github.com/detectivekaktus/JGame/internal/handler"
	"github.com/detectivekaktus/JGame/internal/middleware"
	"github.com/gorilla/mux"
)

func chainMiddlewares(h http.Handler, middlewares ...func(http.Handler)http.Handler) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func NewRouter() *mux.Router {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.Use(middleware.CorsMiddleware)

	// /api/users
	r.Handle("/register",
		chainMiddlewares(http.HandlerFunc(handler.RegisterUser),
			middleware.RequireBodyMiddleware,
			middleware.RequireJsonContentMiddleware)).
		Methods("POST")
	r.Handle("/login",
		chainMiddlewares(http.HandlerFunc(handler.Login),
			middleware.RequireBodyMiddleware,
			middleware.RequireJsonContentMiddleware)).
		Methods("POST")
	r.Handle("/logout",
		chainMiddlewares(http.HandlerFunc(handler.Logout),
			middleware.RejectBodyMiddleware)).
		Methods("POST")
	r.Handle("/users/{id:[0-9]+}",
		chainMiddlewares(http.HandlerFunc(handler.GetUser),
			middleware.RejectBodyMiddleware)).
		Methods("GET")

	// /api/users
	users := r.PathPrefix("/users").Subrouter()
	users.Use(middleware.AuthMiddleware)

	users.Handle("/me",
		chainMiddlewares(http.HandlerFunc(handler.GetCurrentUser),
			middleware.RejectBodyMiddleware)).
		Methods("GET")
	users.Handle("/me",
		chainMiddlewares(http.HandlerFunc(handler.PutCurrentUser),
			middleware.RequireBodyMiddleware,
			middleware.RequireJsonContentMiddleware)).
		Methods("PUT")
	users.Handle("/me",
		chainMiddlewares(http.HandlerFunc(handler.PatchCurrentUser),
			middleware.RequireBodyMiddleware,
			middleware.RequireJsonContentMiddleware)).
		Methods("PATCH")
	users.Handle("/me",
		chainMiddlewares(http.HandlerFunc(handler.DeleteCurrentUser),
			middleware.RejectBodyMiddleware)).
		Methods("DELETE")

	return r
}
