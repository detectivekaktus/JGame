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
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	Id        string
	UserId    int
	CreatedAt time.Time
	ExpiresAt time.Time
}

type LoginResponse struct {
	Message string	             `json:"message"`
	User    VerifiedUserResponse `json:"user"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

// Return the user session stored in the database. The session is retrieved
// via `session_id` cookie attached to the request. If the cookie is not set
// nil is returned. If the session expired nil is returned.
func GetUserSession(conn *pgx.Conn, r *http.Request) (*Session, error) {
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

	if time.Now().UTC().After(session.ExpiresAt.UTC()) {
		database.Execute(conn, "DELETE FROM users.user_session WHERE session_id = $1", session_cookie.Value)
		return nil, errors.New("session expired")
	}

	return &session, nil
}

// Deletes the user session from the database and returns the updated
// `session_id` cookie to send to the client.
func DeleteUserSession(conn *pgx.Conn, id string) (*http.Cookie, error) {
	_, err := database.Execute(conn, "DELETE FROM users.user_session WHERE session_id = $1", id)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name: "session_id",
		Value: "",
		Path: "/",
		Expires: time.Unix(0, 0).UTC(),
		HttpOnly: true,
		Secure: true,
		SameSite: http.SameSiteNoneMode,
		Partitioned: true,
	}, nil
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	conn := database.GetConnection()
	defer conn.Close(context.Background())

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Could not process the body of the request.")
		return
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not generate password hash POST /api/register: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not process user password.")
		return
	}
	user.Password = hash

	var id int
	err = database.QueryRow(conn, "INSERT INTO users.\"user\" (name, email, password, matches_played, matches_won) VALUES ($1, $2, $3, $4, $5) RETURNING user_id",
		user.Name, user.Email, user.Password, 0, 0).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == database.UniqueViolation {
			httputils.SendErrorMessage(w, http.StatusConflict, "Conflict",
				"A user with this email address already exists.")
			return
		}

		fmt.Fprintf(os.Stderr, "Could not insert new user in the table POST /api/register: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not create user.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(VerifiedUserResponse{
		Id: id,
		Name: user.Name,
		Email: user.Email,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	conn := database.GetConnection()
	defer conn.Close(context.Background())

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
	err = database.QueryRow(conn, "SELECT user_id, name, password FROM users.\"user\" WHERE email = $1", user.Email).
		Scan(&user.Id, &user.Name, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputils.SendErrorMessage(w, http.StatusUnauthorized, "Unauthorized",
				"Invalid email credential.")
			return
		}

		fmt.Fprintf(os.Stderr, "Could not retrieve the user password POST /api/login: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not authenticate the user.")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
		httputils.SendErrorMessage(w, http.StatusUnauthorized, "Unauthorized",
			"Invalid password credential.")
		return
	}

	max := new(big.Int).Lsh(big.NewInt(1), 128)
	sessionId, err := rand.Int(rand.Reader, max)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not generate secure id POST /api/login: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not create session.")
		return
	}

	now := time.Now().UTC()
	expires := now.Add(720 * time.Hour) // expires after 30 days
	_, err = database.Execute(conn, "INSERT INTO users.user_session (session_id, user_id, created_at, expires_at) VALUES ($1, $2, $3, $4)",
		sessionId.String(), user.Id, now, expires)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not insert new session into database POST /api/login: %v\n", err)
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
		SameSite: http.SameSiteNoneMode,
		Partitioned: true,
	})

	err = database.QueryRow(conn, "SELECT user_id, email, name, matches_played, matches_won FROM users.\"user\" WHERE user_id = $1", user.Id).
		Scan(&user.Id, &user.Email, &user.Name, &user.MatchesPlayed, &user.MatchesWon)
	if err != nil {
		fmt.Println(err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not get the user.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(LoginResponse{
		Message: "Login successful",
		User: VerifiedUserResponse{
			Id: user.Id,
			Name: user.Name,
			Email: user.Email,
			MatchesPlayed: user.MatchesPlayed,
			MatchesWon: user.MatchesWon,
		},
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	conn := database.GetConnection()
	defer conn.Close(context.Background())

	session, err := GetUserSession(conn, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get the current user session POST /api/logout: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Session retrieval error",
			"Could not get the original session via cookie.")
		return
	}

	cookie, err := DeleteUserSession(conn, session.Id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not delete the current user session POST /api/logout: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Internal error",
			"Could not delete the provided session.")
		return
	}

	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(LogoutResponse{
		Message: "Logout successful",
	})
}

