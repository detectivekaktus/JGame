package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/detectivekaktus/JGame/internal/database"
	"github.com/detectivekaktus/JGame/internal/handler"
	"github.com/detectivekaktus/JGame/internal/httputils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn := database.GetConnection()
		defer conn.Close(context.Background())

		session, err := handler.GetUserSession(conn, r)
		if err != nil {
			if errors.Is(http.ErrNoCookie, err) {
				httputils.SendErrorMessage(w, http.StatusForbidden, "Forbidden",
					"You must be logged in to perform this action.")
				return
			} else if err.Error() == "session expired" {
				httputils.SendErrorMessage(w, http.StatusForbidden, "Forbidden",
					"Your session has expired.")
				return
			} else if err.Error() == "no rows in result set" {
				httputils.SendErrorMessage(w, http.StatusForbidden, "Forbidden",
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
