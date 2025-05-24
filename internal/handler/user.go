package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/detectivekaktus/JGame/internal/database"
	"github.com/detectivekaktus/JGame/internal/httputils"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
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

func hashPassword(passwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	return string(hash), err
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	conn := r.Context().Value("db_connection").(*pgx.Conn)
	session := r.Context().Value("session").(*Session)

	var user User
	err := database.QueryRow(conn, "SELECT user_id, email, name FROM users.\"user\" WHERE user_id = $1", session.UserId).
		Scan(&user.Id, &user.Email, &user.Name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get logged in user for GET /api/users/me: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not get the user.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(VerifiedUserResponse{
		Id: user.Id,
		Name: user.Name,
		Email: user.Email,
	})
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	conn := database.GetConnection()
	defer conn.Close(context.Background())

	var user User
	err := database.QueryRow(conn, "SELECT user_id, name FROM users.\"user\" WHERE user_id = $1", id).
		Scan(&user.Id, &user.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputils.SendErrorMessage(w, http.StatusNotFound, "Not found",
				"No user with given id exists.")
			return
		}

		fmt.Fprintf(os.Stderr, "Could not get user from database for GET /api/users/id: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not get the user with the given id.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(UnverifiedUserResponse{
		Id: user.Id,
		Name: user.Name,
	})
}

func DeleteCurrentUser(w http.ResponseWriter, r *http.Request) {
	conn := r.Context().Value("db_connection").(*pgx.Conn)
	session := r.Context().Value("session").(*Session)

	cookie, err := deleteUserSession(conn, session.Id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not delete user session for DELETE /api/users/id: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not delete user session for the user with the given id.")
		return
	}

	_, err = database.Execute(conn, "DELETE FROM users.\"user\" WHERE user_id = $1", session.UserId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not delete user for DELETE /api/users/id: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not delete user with the given id.")
		return
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusNoContent)
}

func PutCurrentUser(w http.ResponseWriter, r *http.Request) {
	conn := r.Context().Value("db_connection").(*pgx.Conn)
	session := r.Context().Value("session").(*Session)

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Could not process the body of the request.")
		return
	}

	if user.Id != 0 {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Can't modify id of a user.")
		return
	}

	if strings.TrimSpace(user.Name) == "" || strings.TrimSpace(user.Email) == "" || user.Password == "" {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Missing fields",
			"name, email, and password fields must be specified on PUT request.")
		return
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not generate password hash for PUT /api/users/id: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not process user password.")
		return
	}
	user.Password = hash

	_, err = database.Execute(conn, "UPDATE users.\"user\" SET name = $1, email = $2, password = $3 WHERE user_id = $4",
		user.Name, user.Email, user.Password, session.UserId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not update user for PUT /api/users/id: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not update user.")
		return
	}

	cookie, err := deleteUserSession(conn, session.Id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not delete user session for PUT /api/users/id: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not update user session.")
		return
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusNoContent)
}

func PatchCurrentUser(w http.ResponseWriter, r *http.Request) {
	conn := r.Context().Value("db_connection").(*pgx.Conn)
	session := r.Context().Value("session").(*Session)

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Could not process the body of the request.")
		return
	}

	if user.Id != 0 {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Can't modify id of a user.")
		return
	}

	var fieldsSb strings.Builder
	// Even though I know that all the values inside args are of type string
	// I still need to define it as []any, so I can use it later on on database.Execute
	var args []any
	fieldsSb.WriteString("UPDATE users.\"user\" SET ")

	if strings.TrimSpace(user.Email) != "" {
		args = append(args, user.Email)
		fieldsSb.WriteString(fmt.Sprintf("email = $%d", len(args)))
	}

	if strings.TrimSpace(user.Name) != "" {
		if len(args) != 0 {
			fieldsSb.WriteString(", ")
		}
		args = append(args, user.Name)
		fieldsSb.WriteString(fmt.Sprintf("name = $%d", len(args)))
	}

	if user.Password != "" {
		if len(args) != 0 {
			fieldsSb.WriteString(", ")
		}

		hash, err := hashPassword(user.Password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not generate password hash for PATCH /api/users/id: %v\n", err)
			httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
				"Could not process user password.")
			return
		}
		user.Password = hash

		args = append(args, user.Password)
		fieldsSb.WriteString(fmt.Sprintf("password = $%d", len(args)))
	}
	args = append(args, session.UserId)
	fieldsSb.WriteString(fmt.Sprintf(" WHERE user_id = $%d", len(args)))

	_, err = database.Execute(conn, fieldsSb.String(), args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not update user for PATCH /api/users/id: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not update user.")
		return
	}

	if user.Password != "" {
		cookie, err := deleteUserSession(conn, session.Id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not log out user for PATCH /api/users/id: %v\n", err)
			httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
				"Could not update user.")
			return
		}

		http.SetCookie(w, cookie)
	}
	
	w.WriteHeader(http.StatusNoContent)
}

