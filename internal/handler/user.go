package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/detectivekaktus/JGame/internal/database"
	"github.com/detectivekaktus/JGame/internal/httputil"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Session struct {
	Id        int
	UserId    int
	CreatedAt time.Time
	ExpiresAt time.Time
}

type VerifiedUserResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UnverifiedUserResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
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

	if time.Now().After(session.ExpiresAt) {
		return nil, err
	}

	return &session, nil
}

func hashPassword(passwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	return string(hash), err
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !httputil.IsContentType(w, r, "application/json") {
		return
	}

	if !httputil.HasContent(r) {
		httputil.SendErrorMessage(w, http.StatusBadRequest, "Content error",
			"Expected content inside the request body, got nothing.")
		return
	}

	conn := database.GetConnection()
	defer conn.Close(context.Background())

	if session, _ := getUserSession(conn, r); session != nil {
		httputil.SendErrorMessage(w, http.StatusForbidden, "Authentication error",
			"Cannot POST a user while logged in.")
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not process JSON body for POST /api/users: %v\n", err)
		httputil.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not process the body of the request.")
		return
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not generate password hash for POST /api/users: %v\n", err)
		httputil.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not process user password.")
		return
	}
	user.Password = hash

	var id int
	err = database.QueryRow(conn, "INSERT INTO users.\"user\" (name, email, password) VALUES ($1, $2, $3) RETURNING user_id",
		user.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == database.UniqueViolation {
			fmt.Fprintf(os.Stderr, "Unique email constraint violation when inserting new user for POST /api/users: %v\n", err)
			httputil.SendErrorMessage(w, http.StatusConflict, "Internal error",
				"A user with this email address already exists.")
			return
		}

		fmt.Fprintf(os.Stderr, "Could not insert new user in the table for POST /api/users: %v\n", err)
		httputil.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not create user.")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(VerifiedUserResponse{
		Id: id,
		Name: user.Name,
		Email: user.Email,
	})
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if httputil.HasContent(r) {
		httputil.SendErrorMessage(w, http.StatusBadRequest, "Request body not allowed",
			"This endpoint does not accept a request body.")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	conn := database.GetConnection()
	defer conn.Close(context.Background())

	var user User
	err := database.QueryRow(conn, "SELECT user_id, email, name FROM users.\"user\" WHERE user_id = $1", id).
		Scan(&user.Id, &user.Email, &user.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputil.SendErrorMessage(w, http.StatusNotFound, "Not found",
				"No user with given id exists.")
			return
		}

		fmt.Fprintf(os.Stderr, "Could not get user from database for GET /api/users/id: %v\n", err)
		httputil.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not get the user with the given id.")
		return
	}

	w.WriteHeader(http.StatusOK)

	if session, _ := getUserSession(conn, r); session != nil {
		if session.UserId == user.Id {
			json.NewEncoder(w).Encode(VerifiedUserResponse{
				Id: user.Id,
				Email: user.Email,
				Name: user.Name,
			})
			return
		}
	}

	json.NewEncoder(w).Encode(UnverifiedUserResponse{
		Id: user.Id,
		Name: user.Name,
	})
}
