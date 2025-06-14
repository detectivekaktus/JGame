package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/detectivekaktus/JGame/internal/config"
	"github.com/detectivekaktus/JGame/internal/database"
	"github.com/detectivekaktus/JGame/internal/handler"
	"github.com/detectivekaktus/JGame/internal/httputils"
)

// Establishes a connection to the database passed with the request via
// context to the next handler. Retrieves the user session and checks:
// 1. Does the session exist (both valid and invalid)?
// 2. Is the session expired?
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn := database.GetConnection()
		defer conn.Close(context.Background())

		session, err := handler.GetUserSession(conn, r)
		if err != nil {
			if errors.Is(http.ErrNoCookie, err) {
				httputils.SendErrorMessage(w, http.StatusUnauthorized, "Unauthorized",
					"You must be logged in to perform this action.")
				return
			} else if err.Error() == "session expired" {
				cookie, _ := handler.DeleteUserSession(conn, session.Id)
				http.SetCookie(w, cookie)
				httputils.SendErrorMessage(w, http.StatusForbidden, "Forbidden",
					"Your session has expired.")
				return
			} else if err.Error() == "no rows in result set" {
				httputils.SendErrorMessage(w, http.StatusUnauthorized, "Unauthorized",
					"Invalid session id.")
				return
			}

			fmt.Fprintf(os.Stderr, "AuthMiddleware error: %v\n", err)
			httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
				"Could not authenticate the user.")
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "session", session))
		r = r.WithContext(context.WithValue(r.Context(), "db_connection", conn))
		next.ServeHTTP(w, r)
	})
}

// The requests with this middleware will not be allowed if they have body.
func RejectBodyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if httputils.HasContent(r) {
			httputils.SendErrorMessage(w, http.StatusBadRequest, "Request body not allowed",
				"This endpoint does not accept a request body.")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// The requests with this middleware will not be allowed if they don't have body.
func RequireBodyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !httputils.HasContent(r) {
			httputils.SendErrorMessage(w, http.StatusBadRequest, "No content provided",
				"Expected content inside the request body, got nothing.")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// The requests with this middleware will not be allowed if their Content-Type
// isn't application/json
func RequireJsonContentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !httputils.IsContentType(w, r, "application/json") {
			httputils.SendErrorMessage(w, http.StatusBadRequest, "Content-Type mismatch",
				"Expected Content-Type to be application/json.")
			return
		}
		next.ServeHTTP(w, r)
	})
}

var AllowedCorsOrigins = map[string]bool {
	"https://127.0.0.1:5173": true,
	"https://localhost:5173": true,
	"https://" + config.AppConfig.LocalIp + ":5173": true,
}

// Sets up Cross-Origin Resource Sharing mechanism workarounds to accept requests
// only from the frontend (localhost:5173).
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if AllowedCorsOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
