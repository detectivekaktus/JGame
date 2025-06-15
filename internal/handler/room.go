package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/detectivekaktus/JGame/internal/database"
	"github.com/detectivekaktus/JGame/internal/httputils"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type Room struct {
	Id           int    `json:"room_id"`
	Name         string `json:"name"`
	PackId       int    `json:"pack_id"`
	UserId       int    `json:"user_id"`
	CurrentUsers int    `json:"current_users"`
	MaxUsers     int    `json:"max_users"`
	Password		 string `json:"password"`
}

// same as the one above, but without Password fields
type RoomResponse struct {
	Id           int    `json:"room_id"`
	Name         string `json:"name"`
	PackId       int    `json:"pack_id"`
	UserId       int    `json:"user_id"`
	CurrentUsers int    `json:"current_users"`
	MaxUsers     int    `json:"max_users"`
}

type RoomStatusResponse struct {
	Message string `json:"message"`
}

const (
	MAX_USERS_IN_ROOM  = 2 << 3
	MAX_ROOMS_RESPONSE = 2 << 5
)

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	conn := r.Context().Value("db_connection").(*pgx.Conn)
	session := r.Context().Value("session").(*Session)

	var inGame bool
	err := database.QueryRow(conn, "SELECT EXISTS (SELECT * FROM rooms.player WHERE user_id = $1)", session.UserId).
		Scan(&inGame)
	if err != nil {
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not verify user status.")
		return
	}

	if inGame {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Already in game.")
		return
	}

	var requestedRoom Room
	err = json.NewDecoder(r.Body).Decode(&requestedRoom)
	if err != nil {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Could not process the body of the request.")
		return
	}

	var packExists bool
	err = database.QueryRow(conn, "SELECT EXISTS (SELECT * FROM packs.pack WHERE pack_id = $1)", requestedRoom.PackId).
		Scan(&packExists)
	if err != nil {
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not retrieve pack associated with the room.")
		return
	}

	if !packExists {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"No pack with the given id exists.")
		return
	}

	room := &Room{
		Id: rand.Intn(2 << 15),
		Name: requestedRoom.Name,
		PackId: requestedRoom.PackId,
		UserId: session.UserId,
		CurrentUsers: 1,
		MaxUsers: MAX_USERS_IN_ROOM,
		Password: requestedRoom.Password,
	}

	_, err = database.Execute(conn, "INSERT INTO rooms.room (room_id, user_id, name, pack_id, current_users, max_users, password) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		room.Id, room.UserId, room.Name, room.PackId, room.CurrentUsers, room.MaxUsers, room.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create room POST /api/rooms: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not create room.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(RoomResponse{
		Id: room.Id,
		Name: room.Name,
		PackId: room.PackId,
		UserId: room.UserId,
		CurrentUsers: 0,
		MaxUsers: room.MaxUsers,
	})
}

func PutRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var requestedRoom Room
	err := json.NewDecoder(r.Body).Decode(&requestedRoom)
	if err != nil {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Could not process the body of the request.")
		return
	}

	if strings.TrimSpace(requestedRoom.Name) == "" || requestedRoom.PackId == 0 ||
		 strings.TrimSpace(requestedRoom.Password) == "" {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Missing fields",
			"name, password and pack_id fields must be specified on PUT request.")
		return
	}

	if requestedRoom.UserId != 0 || requestedRoom.Id != 0 ||
		requestedRoom.CurrentUsers != 0 || requestedRoom.MaxUsers != 0 {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Modifying non-editable fields",
			"user_id, current_users, max_users fields can't be changed via PUT request.")
		return
	}

	session := r.Context().Value("session").(*Session)
	conn := r.Context().Value("db_connection").(*pgx.Conn)

	var room Room
	err = database.QueryRow(conn, "SELECT * FROM rooms.room WHERE room_id = $1", id).
		Scan(&room.Id, &room.UserId, &room.Name, &room.PackId, &room.CurrentUsers, &room.MaxUsers, &room.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputils.SendErrorMessage(w, http.StatusNotFound, "Not found",
				"No room with given id exists.")
			return
		}
		fmt.Fprintf(os.Stderr, "Could not get room from database PUT /api/room/id: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not get the room with the given id.")
		return
	}

	if session.UserId != room.UserId {
		httputils.SendErrorMessage(w, http.StatusForbidden, "Forbidden",
			"Can't modify a room that is not owned by themselves.")
		return
	}

	var packExists bool
	err = database.QueryRow(conn, "SELECT EXISTS (SELECT * FROM packs.pack WHERE pack_id = $1)", requestedRoom.PackId).
		Scan(&packExists)
	if err != nil {
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not retrieve pack associated with the room.")
		return
	}

	if !packExists {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"No pack with the given id exists.")
		return
	}

	err = database.QueryRow(conn, "UPDATE rooms.room SET name = $1, pack_id = $2, password = $3 WHERE room_id = $4 RETURNING name, pack_id",
		requestedRoom.Name, requestedRoom.PackId, requestedRoom.Password, room.Id).
	Scan(&room.Name, &room.PackId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not update room from database PUT /api/room/id: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not update the room with the given id.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(RoomResponse{
		Id: room.Id,
		Name: room.Name,
		PackId: room.PackId,
		UserId: room.UserId,
		CurrentUsers: room.CurrentUsers,
		MaxUsers: room.MaxUsers,
	})
}

func PatchRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var requestedRoom Room
	err := json.NewDecoder(r.Body).Decode(&requestedRoom)
	if err != nil {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Could not process the body of the request.")
		return
	}

	if requestedRoom.UserId != 0 || requestedRoom.Id != 0 ||
		requestedRoom.CurrentUsers != 0 || requestedRoom.MaxUsers != 0 {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Modifying non-editable fields",
			"room_id, user_id, current_users, max_users fields can't be changed via PATCH request.")
		return
	}

	conn := r.Context().Value("db_connection").(*pgx.Conn)
	session := r.Context().Value("session").(*Session)

	var room Room
	err = database.QueryRow(conn, "SELECT * FROM rooms.room WHERE room_id = $1", id).
		Scan(&room.Id, &room.UserId, &room.Name, &room.PackId, &room.CurrentUsers, &room.MaxUsers, &room.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputils.SendErrorMessage(w, http.StatusNotFound, "Not found",
				"No room with given id exists.")
			return
		}
		fmt.Fprintf(os.Stderr, "Could not get room from database PUT /api/room/id: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not get the room with the given id.")
		return
	}

	if session.UserId != room.UserId {
		httputils.SendErrorMessage(w, http.StatusForbidden, "Forbidden",
			"Can't modify a room that is not owned by themselves.")
		return
	}

	var fieldsSb strings.Builder
	var args []any // check handler/user.go:221
	fieldsSb.WriteString("UPDATE rooms.room SET ")

	if strings.TrimSpace(requestedRoom.Name) != "" {
		args = append(args, requestedRoom.Name)
		fieldsSb.WriteString(fmt.Sprintf("name = $%d", len(args)))
	}

	if requestedRoom.PackId != 0 {
		var packExists bool
		err = database.QueryRow(conn, "SELECT EXISTS (SELECT * FROM packs.pack WHERE pack_id = $1)", requestedRoom.PackId).
			Scan(&packExists)
		if err != nil {
			httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
				"Could not retrieve pack associated with the room.")
			return
		}

		if !packExists {
			httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
				"No pack with the given id exists.")
			return
		}

		if len(args) != 0 {
			fieldsSb.WriteString(", ")
		}

		args = append(args, requestedRoom.PackId)
		fieldsSb.WriteString(fmt.Sprintf("pack_id = $%d", len(args)))
	}

	if strings.ToUpper(strings.TrimSpace(requestedRoom.Password)) != "PASSWORD_UNSET" {
		if requestedRoom.Password == "" {
			httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
				"Use PASSWORD_UNSET constant to set no password.")
			return
	}

		if len(args) != 0 {
			fieldsSb.WriteString(", ")
		}

		args = append(args, requestedRoom.Password)
		fieldsSb.WriteString(fmt.Sprintf("password = $%d", len(args)))
	}

	args = append(args, id)
	fieldsSb.WriteString(fmt.Sprintf(" WHERE room_id = $%d RETURNING name, pack_id", len(args)))
	err = database.QueryRow(conn, fieldsSb.String(), args...).
		Scan(&room.Name, &room.PackId)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(RoomResponse{
		Id: room.Id,
		Name: room.Name,
		PackId: room.PackId,
		UserId: room.UserId,
		CurrentUsers: room.CurrentUsers,
		MaxUsers: room.MaxUsers,
	})
}

func DeleteRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	session := r.Context().Value("session").(*Session)
	conn := r.Context().Value("db_connection").(*pgx.Conn)

	var room Room
	err := database.QueryRow(conn, "SELECT * FROM rooms.room WHERE room_id = $1", id).
		Scan(&room.Id, &room.UserId, &room.Name, &room.PackId, &room.CurrentUsers, &room.MaxUsers, &room.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputils.SendErrorMessage(w, http.StatusNotFound, "Not found",
				"No room with given id exists.")
			return
		}
		fmt.Fprintf(os.Stderr, "Could not get room from database PUT /api/room/id: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not get the room with the given id.")
		return
	}

	if session.UserId != room.UserId {
		httputils.SendErrorMessage(w, http.StatusForbidden, "Forbidden",
			"Can't delete a room that is not owned by themselves.")
		return
	}

	_, err = database.Execute(conn, "DELETE FROM rooms.player WHERE room_id = $1", room.Id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not delete user status DELETE /api/rooms/id: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not delete user status.")
		return
	}

	_, err = database.Execute(conn, "DELETE FROM rooms.room WHERE room_id = $1", room.Id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not delete room DELETE /api/rooms/id: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not delete room.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	conn := database.GetConnection()
	defer conn.Close(context.Background())

	var room Room
	err := database.QueryRow(conn, "SELECT * FROM rooms.room WHERE room_id = $1", id).
		Scan(&room.Id, &room.UserId, &room.Name, &room.PackId, &room.CurrentUsers, &room.MaxUsers, &room.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputils.SendErrorMessage(w, http.StatusNotFound, "Not found",
				"No room with given id exists.")
			return
		}
		fmt.Fprintf(os.Stderr, "Could not get room from database PUT /api/room/id: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not get the room with the given id.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(RoomResponse{
		Id: room.Id,
		Name: room.Name,
		PackId: room.PackId,
		UserId: room.UserId,
		CurrentUsers: room.CurrentUsers,
		MaxUsers: room.MaxUsers,
	})
}

func GetRooms(w http.ResponseWriter, r *http.Request) {
	name := "%" + strings.ToLower(r.URL.Query().Get("name")) + "%"

	conn := database.GetConnection()
	defer conn.Close(context.Background())

	var rows pgx.Rows
	if name == "" {
		rows = database.QueryRows(conn, "SELECT * FROM rooms.room LIMIT $1", MAX_PACKS_RESPONSE)
	} else {
		rows = database.QueryRows(conn, "SELECT * FROM rooms.room WHERE LOWER(name) ILIKE $1 LIMIT $2", name, MAX_ROOMS_RESPONSE)
	}
	defer rows.Close()

	var rooms []RoomResponse
	for rows.Next() {
		var room Room
		err := rows.Scan(&room.Id, &room.UserId, &room.Name, &room.PackId, &room.CurrentUsers, &room.MaxUsers, &room.Password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not read rooms at GET /api/rooms: %v", err)
			httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
				"Could not read rooms")
			return
		}

		rooms = append(rooms, RoomResponse{
			Id: room.Id,
			UserId: room.UserId,
			Name: room.Name,
			PackId: room.PackId,
			CurrentUsers: room.CurrentUsers,
			MaxUsers: room.MaxUsers,
		})
	}

	err := rows.Err()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not iterate rooms at GET /api/rooms: %v", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not iterate rooms")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(rooms)
}

