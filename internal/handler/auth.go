package handler

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/detectivekaktus/JGame/internal/database"
	"github.com/detectivekaktus/JGame/internal/httputils"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	Id        int
	UserId    int
	CreatedAt time.Time
	ExpiresAt time.Time
}

// Return the user session stored in the database. The session is retrieved
// via `session_id` cookie attached to the request. If the cookie is not set
// nil is returned. If the session expired nil is returned.
func getUserSession(conn *pgx.Conn, r *http.Request) (*Session, error) {
	session_cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}

	var session Session
	err = database.QueryRow(conn, "SELECT * FROM users.user_session WHERE session_id = $1", session_cookie.Value).
		Scan(&session.Id, &session.UserId, &session.CreatedAt, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}

	if time.Now().UTC().After(session.ExpiresAt) {
		return nil, err
	}

	return &session, nil
}

func Login(w http.ResponseWriter, r *http.Request) {
	if !httputils.IsContentType(w, r, "application/json") {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Content-Type mismatch",
			"Expected Content-Type to be application/json.")
		return
	}

	if !httputils.HasContent(r) {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "No content provided",
			"Expected content inside the request body, got nothing.")
		return
	}

	conn := database.GetConnection()
	defer conn.Close(context.Background())

	if session, _ := getUserSession(conn, r); session != nil {
		httputils.SendErrorMessage(w, http.StatusForbidden, "Forbidden",
			"Can't log in when already logged in.")
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Could not process the body of the request.")
		return
	}

	if strings.TrimSpace(user.Email) == "" || user.Password == "" {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Must specify valid email and password.")
		return
	}

	var hashedPassword string
	err = database.QueryRow(conn, "SELECT user_id, password FROM users.\"user\" WHERE email = $1", user.Email).
		Scan(&user.Id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputils.SendErrorMessage(w, http.StatusForbidden, "Forbidden",
				"Invalid email credential.")
			return
		}

		fmt.Fprintf(os.Stderr, "Could not retrieve the user password for POST /api/login: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not authenticate the user.")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
		httputils.SendErrorMessage(w, http.StatusForbidden, "Forbidden",
			"Invalid password credential.")
		return
	}

	max := new(big.Int).Lsh(big.NewInt(1), 128)
	sessionId, err := rand.Int(rand.Reader, max)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not generate secure id for POST /api/login: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not create session.")
		return
	}

	now := time.Now().UTC()
	expires := now.Add(720 * time.Hour)
	_, err = database.Execute(conn, "INSERT INTO users.user_session (session_id, user_id, created_at, expires_at) VALUES ($1, $2, $3, $4)",
		sessionId.String(), user.Id, now, expires) // expires after 30 days
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not insert new session into database for POST /api/login: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not create session.")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "session_id",
		Value: sessionId.String(),
		Path: "/",
		Expires: expires,
		HttpOnly: true,
		Secure: true,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusNoContent)
}

