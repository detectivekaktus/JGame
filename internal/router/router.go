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

	r.Handle("/register",
		chainMiddlewares(http.HandlerFunc(handler.RegisterUser),
			middleware.RequireBodyMiddleware,
			middleware.RequireJsonContentMiddleware)).
		Methods("POST", "OPTIONS")
	r.Handle("/login",
		chainMiddlewares(http.HandlerFunc(handler.Login),
			middleware.RequireBodyMiddleware,
			middleware.RequireJsonContentMiddleware)).
		Methods("POST", "OPTIONS")
	r.Handle("/logout",
		chainMiddlewares(http.HandlerFunc(handler.Logout),
			middleware.RejectBodyMiddleware)).
		Methods("POST", "OPTIONS")

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
		Methods("PUT", "OPTIONS")
	users.Handle("/me",
		chainMiddlewares(http.HandlerFunc(handler.PatchCurrentUser),
			middleware.RequireBodyMiddleware,
			middleware.RequireJsonContentMiddleware)).
		Methods("PATCH", "OPTIONS")
	users.Handle("/me",
		chainMiddlewares(http.HandlerFunc(handler.DeleteCurrentUser),
			middleware.RejectBodyMiddleware)).
		Methods("DELETE", "OPTIONS")
	// Available without auth
	r.Handle("/users/{id:[0-9]+}",
		chainMiddlewares(http.HandlerFunc(handler.GetUser),
			middleware.RejectBodyMiddleware)).
		Methods("GET")

	packs := r.PathPrefix("/packs").Subrouter()
	packs.Use(middleware.AuthMiddleware)
	packs.Handle("",
		chainMiddlewares(http.HandlerFunc(handler.CreatePack),
			middleware.RequireBodyMiddleware,
			middleware.RequireJsonContentMiddleware)).
		Methods("POST", "OPTIONS")
	packs.Handle("/{id:[0-9]+}",
		chainMiddlewares(http.HandlerFunc(handler.PutPack),
			middleware.RequireBodyMiddleware,
			middleware.RequireJsonContentMiddleware)).
		Methods("PUT", "OPTIONS")
	packs.Handle("/{id:[0-9]+}",
		chainMiddlewares(http.HandlerFunc(handler.PatchPack),
			middleware.RequireBodyMiddleware,
			middleware.RequireJsonContentMiddleware)).
		Methods("PATCH", "OPTIONS")
	packs.Handle("/{id:[0-9]+}",
		chainMiddlewares(http.HandlerFunc(handler.DeletePack),
			middleware.RejectBodyMiddleware)).
		Methods("DELETE", "OPTIONS")
	// Available without auth
	r.Handle("/packs",
		chainMiddlewares(http.HandlerFunc(handler.GetPacks),
			middleware.RejectBodyMiddleware)).
		Methods("GET")
	r.Handle("/packs/{id:[0-9]+}",
		chainMiddlewares(http.HandlerFunc(handler.GetPack),
			middleware.RejectBodyMiddleware)).
		Methods("GET")

	return r
}
