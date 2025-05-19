package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/detectivekaktus/JGame/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

func hashPassword(passwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	return string(hash), err
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if ctnType := r.Header.Get("Content-Type"); ctnType != "application/json" {
		SendErrorMessage(w, http.StatusBadRequest, "Content error",
			fmt.Sprintf("Excepted to find Content-Type `application/json`, found `%s`.", ctnType))
		return
	}

	if r.ContentLength == 0 {
		SendErrorMessage(w, http.StatusBadRequest, "Content error",
			"Expected a valid user JSON in the body, got nothing.")
		return
	}

	if _, err := r.Cookie("session_id"); err == nil {
		SendErrorMessage(w, http.StatusForbidden, "Authentication error",
			"Cannot POST a user while logged in.")
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not process JSON body for POST /api/user: %v\n", err)
		SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not process the body of the request.")
		return
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not generate password hash for POST /api/user: %v\n", err)
		SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
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
		fmt.Fprintf(os.Stderr, "Could not insert new user in the table for POST /api/user: %v\n", err)
		SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
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
