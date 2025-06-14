package router

import (
	"net/http"

	"github.com/detectivekaktus/JGame/internal/handler"
	"github.com/detectivekaktus/JGame/internal/middleware"
	"github.com/detectivekaktus/JGame/internal/websocket"
	"github.com/gorilla/mux"
)

func chainMiddlewares(h http.Handler, middlewares ...func(http.Handler)http.Handler) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func NewRouter() *mux.Router {
	root := mux.NewRouter()
	root.Use(middleware.CorsMiddleware)

	api := root.PathPrefix("/api").Subrouter()

	api.Handle("/register",
		chainMiddlewares(http.HandlerFunc(handler.RegisterUser),
			middleware.RequireBodyMiddleware,
			middleware.RequireJsonContentMiddleware)).
		Methods("POST", "OPTIONS")
	api.Handle("/login",
		chainMiddlewares(http.HandlerFunc(handler.Login),
			middleware.RequireBodyMiddleware,
			middleware.RequireJsonContentMiddleware)).
		Methods("POST", "OPTIONS")
	api.Handle("/logout",
		chainMiddlewares(http.HandlerFunc(handler.Logout),
			middleware.RejectBodyMiddleware)).
		Methods("POST", "OPTIONS")

	users := api.PathPrefix("/users").Subrouter()
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
	api.Handle("/users/{id:[0-9]+}",
		chainMiddlewares(http.HandlerFunc(handler.GetUser),
			middleware.RejectBodyMiddleware)).
		Methods("GET")

	packs := api.PathPrefix("/packs").Subrouter()
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
	api.Handle("/packs",
		chainMiddlewares(http.HandlerFunc(handler.GetPacks),
			middleware.RejectBodyMiddleware)).
		Methods("GET")
	api.Handle("/packs/{id:[0-9]+}",
		chainMiddlewares(http.HandlerFunc(handler.GetPack),
			middleware.RejectBodyMiddleware)).
		Methods("GET")

	rooms := api.PathPrefix("/rooms").Subrouter()
	rooms.Use(middleware.AuthMiddleware)
	rooms.Handle("",
		chainMiddlewares(http.HandlerFunc(handler.CreateRoom),
			middleware.RequireBodyMiddleware,
			middleware.RequireJsonContentMiddleware)).
	Methods("POST", "OPTIONS")
	rooms.Handle("/{id:[0-9]+}",
		chainMiddlewares(http.HandlerFunc(handler.PutRoom),
			middleware.RequireBodyMiddleware,
			middleware.RequireJsonContentMiddleware)).
	Methods("PUT", "OPTIONS")
	rooms.Handle("/{id:[0-9]+}",
		chainMiddlewares(http.HandlerFunc(handler.PatchRoom),
			middleware.RequireBodyMiddleware,
			middleware.RequireJsonContentMiddleware)).
	Methods("PATCH", "OPTIONS")
	rooms.Handle("/{id:[0-9]+}",
		chainMiddlewares(http.HandlerFunc(handler.DeleteRoom),
			middleware.RejectBodyMiddleware)).
	Methods("DELETE", "OPTIONS")
	// Available without auth
	api.Handle("/rooms",
		chainMiddlewares(http.HandlerFunc(handler.GetRooms),
			middleware.RejectBodyMiddleware)).
		Methods("GET")
	api.Handle("/rooms/{id:[0-9]+}",
		chainMiddlewares(http.HandlerFunc(handler.GetRoom),
			middleware.RejectBodyMiddleware)).
		Methods("GET")

	ws := root.PathPrefix("/ws").Subrouter()
	ws.Use(middleware.AuthMiddleware)
	ws.Handle("",
		chainMiddlewares(http.HandlerFunc(websocket.WebsocketHandler))).
		Methods("GET")

	return root
}
