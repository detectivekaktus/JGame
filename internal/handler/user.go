package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/detectivekaktus/JGame/internal/database"
	"github.com/detectivekaktus/JGame/internal/httputil"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func hashPassword(passwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	return string(hash), err
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !httputil.IsContentType(w, r, "application/json") || !httputil.HasContent(w, r) {
		return
	}

	if _, err := r.Cookie("session_id"); err == nil {
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

	conn := database.GetConnection()
	defer conn.Close(context.Background())

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
	json.NewEncoder(w).Encode(UserResponse{
		Id: id,
		Name: user.Name,
		Email: user.Email,
	})
}
