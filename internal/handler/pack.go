package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/detectivekaktus/JGame/internal/database"
	"github.com/detectivekaktus/JGame/internal/httputils"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

const MAX_PACKS_RESPONSE = 2 << 7

type Pack struct {
	Id      int             `json:"id"`
	UserId  int             `json:"user_id"`
	Name    string          `json:"name"`
	Body    json.RawMessage `json:"body"`
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
	vars := mux.Vars(r)
	id := vars["id"]

	conn := database.GetConnection()
	defer conn.Close(context.Background())

	var pack Pack
	err := database.QueryRow(conn, "SELECT * FROM packs.pack WHERE pack_id = $1", id).
		Scan(&pack.Id, &pack.UserId, &pack.Body, &pack.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputils.SendErrorMessage(w, http.StatusNotFound, "Not found",
				"No pack with given id exists.")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(pack)
}

// Can apply `name` filter to the result. Returns max MAX_PACKS_RESPONSE packs.
func GetPacks(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	conn := database.GetConnection()
	defer conn.Close(context.Background())

	var rows pgx.Rows
	if name == "" {
		rows = database.QueryRows(conn, "SELECT * FROM packs.pack LIMIT $1", MAX_PACKS_RESPONSE)
		defer rows.Close()
	} else {
		rows = database.QueryRows(conn, "SELECT * FROM packs.pack WHERE name = $1 LIMIT $2", name, MAX_PACKS_RESPONSE)
		defer rows.Close()
	}

	var packs []Pack
	for rows.Next() {
		var pack Pack
		err := rows.Scan(&pack.Id, &pack.UserId, &pack.Body, &pack.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not read packs at GET /api/packs: %v", err)
			httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
				"Could not read packs")
			return
		}
		packs = append(packs, pack)
	}

	err := rows.Err()
	if err != nil {
			fmt.Fprintf(os.Stderr, "Could not iterate packs at GET /api/packs: %v", err)
			httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
				"Could not iterate packs")
			return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(packs)
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
	vars := mux.Vars(r)
	id := vars["id"]

	conn := r.Context().Value("db_connection").(*pgx.Conn)

	_, err := database.Execute(conn, "DELETE FROM packs.pack WHERE pack_id = $1", id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not delete pack at /api/packs/id: %v", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not delete pack with the given id.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

