package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/detectivekaktus/JGame/internal/database"
	"github.com/detectivekaktus/JGame/internal/httputils"
	"github.com/jackc/pgx/v5"
)

type Pack struct {
	Id 			int 						`json:"id"`
	UserId 	int 						`json:"user_id"`
	Name 		string 					`json:"name"`
	Body 		json.RawMessage `json:"body"`
}

// Expected body:
// {
//   "name": ... (32 chars max)
//   "user_id": ... (owner_id)
//   "body": ... (json)
// }
func CreatePack(w http.ResponseWriter, r *http.Request) {
	conn := r.Context().Value("db_connection").(*pgx.Conn)

	var pack Pack
	err := json.NewDecoder(r.Body).Decode(&pack)
	if err != nil {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Could not process the body of the request.")
		return
	}

	err = database.QueryRow(conn, "INSERT INTO packs.pack (user_id, name, body) VALUES ($1, $2, $3) RETURNING pack_id",
		pack.UserId, pack.Name, pack.Body).
		Scan(&pack.Id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create a new pack at POST /api/packs: %v", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not create pack.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(pack)
}

func GetPack(w http.ResponseWriter, r *http.Request) {
	
}

// Can apply `name` filter to the result.
func GetPacks(w http.ResponseWriter, r *http.Request) {
	
}

// Expected body:
// {
//   "name": ... (32 chars max)
//   "user_id": ... (owner_id)
//   "body": ... (json)
// }
func PutPack(w http.ResponseWriter, r *http.Request) {
	
}

// Expected body:
// {
//   "name": ... (32 chars max)
//   "user_id": ... (owner_id)
//   "body": ... (json)
// }
func PatchPack(w http.ResponseWriter, r *http.Request) {
	
}

func DeletePack(w http.ResponseWriter, r *http.Request) {
	
}

