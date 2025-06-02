package handler

// The quiz game packs are stored inside packs.pack PostgreSQL table
// Currently the table looks like this, considering all the migrations
// done in the past:
// packs.pack(
//   pack_id primary int,
//   user_id int (from users.user)
//   body    json
//   name    varchar(32)
// )
//
// The body is defined within /api/pack_schema.json schema file and all
// the operations on the packs that require creation or modification on
// packs are compared against the schema.

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/detectivekaktus/JGame/internal/database"
	"github.com/detectivekaktus/JGame/internal/httputils"
	"github.com/detectivekaktus/JGame/internal/validation"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const MAX_PACKS_RESPONSE = 2 << 7

type Pack struct {
	Id      int             `json:"id"`
	UserId  int             `json:"user_id"`
	Name    string          `json:"name"`
	Body    json.RawMessage `json:"body"`
}

func CreatePack(w http.ResponseWriter, r *http.Request) {
	conn := r.Context().Value("db_connection").(*pgx.Conn)

	var pack Pack
	err := json.NewDecoder(r.Body).Decode(&pack)
	if err != nil {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Could not process the body of the request.")
		return
	}

	if !validation.ValidateAgainstSchema(validation.PACK_SCHEMA, pack.Body) {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"The body parameter does not satisfy the schema.")
		return
	}

	err = database.QueryRow(conn, "INSERT INTO packs.pack (user_id, name, body) VALUES ($1, $2, $3) RETURNING pack_id",
		pack.UserId, pack.Name, pack.Body).
		Scan(&pack.Id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == database.ForeignKeyViolation {
			httputils.SendErrorMessage(w, http.StatusBadRequest, "User error",
				"This user can't own a pack.")
			return
		}

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

func PutPack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	conn := r.Context().Value("db_connection").(*pgx.Conn)
	
	var pack Pack 
	err := json.NewDecoder(r.Body).Decode(&pack)
	if err != nil {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Could not process the body of the request.")
		return
	}

	if pack.Id != 0 {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Can't modify id of a pack.")
		return
	}

	if strings.TrimSpace(pack.Name) == "" || pack.UserId == 0 {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Missing fields",
			"name, user_id, and body fields must be specified on PUT request.")
		return
	}

	if !validation.ValidateAgainstSchema(validation.PACK_SCHEMA, pack.Body) {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"The body parameter does not satisfy the schema.")
		return
	}

	_, err = database.Execute(conn, "UPDATE packs.pack SET user_id = $1, body = $2, name = $3",
		pack.UserId, pack.Body, pack.Name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not update pack at PUT /api/packs/id: %v", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not update pack.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	intid, _ := strconv.Atoi(id)
	json.NewEncoder(w).Encode(Pack{
		Id: intid,
		UserId: pack.UserId,
		Body: pack.Body,
		Name: pack.Name,
	})
}

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

