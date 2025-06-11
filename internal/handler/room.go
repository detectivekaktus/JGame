package handler

import (
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/detectivekaktus/JGame/internal/database"
	"github.com/detectivekaktus/JGame/internal/httputils"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

// -- DESING NOTE -------------
// I made decision to store rooms and room related information
// directly in RAM to optimize read and write operations on the
// objects. To prevent the overuse of RAM resources, there's
// MAX_ROOMS constant which defines the maximum number of rooms
// that can exist at the same point in time. The O(n) operations
// on arrays (slices) are not the problem, since there are not
// that many elements.
//
// A potential drawback of this decision is the complete loss of data
// and progress related to rooms, in case of server going down or rebooting.

type Room struct {
	Id           int    `json:"room_id"`
	Name         string `json:"name"`
	PackId       int    `json:"pack_id"`
	UserId       int    `json:"user_id"`
	Users        []int  `json:"users"`
	CurrentUsers int    `json:"current_users"`
	MaxUsers     int    `json:"max_users"`
	BannedUsers  []int  `json:"-"`
	Password		 string `json:"password"`
}

// same as the one above, but without Password and BannedUsers fields
type RoomResponse struct {
	Id           int    `json:"room_id"`
	Name         string `json:"name"`
	PackId       int    `json:"pack_id"`
	UserId       int    `json:"user_id"`
	Users        []int  `json:"users"`
	CurrentUsers int    `json:"current_users"`
	MaxUsers     int    `json:"max_users"`
}

type RoomStatusResponse struct {
	Message string `json:"message"`
}

const (
	MAX_USERS_IN_ROOM  = 2 << 3
	MAX_ROOMS          = 2 << 9
	MAX_ROOMS_RESPONSE = 2 << 5
)

var usersInGame map[int]bool = make(map[int]bool)
var rooms []*Room = make([]*Room, 0, MAX_ROOMS)

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	if len(rooms) >= MAX_ROOMS {
		httputils.SendErrorMessage(w, http.StatusServiceUnavailable, "Service unavailable",
			"Room limit reached.")
		return
	}

	conn := r.Context().Value("db_connection").(*pgx.Conn)
	session := r.Context().Value("session").(*Session)

	if usersInGame[session.UserId] {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Already in game.")
		return
	}

	var requestedRoom Room
	err := json.NewDecoder(r.Body).Decode(&requestedRoom)
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

	roomId := 1
	for i, room := range rooms {
		if room == nil {
			roomId = i + 1
			break
		}
	}

	room := &Room{
		Id: roomId,
		Name: requestedRoom.Name,
		PackId: requestedRoom.PackId,
		UserId: session.UserId,
		Users: make([]int, 0, MAX_USERS_IN_ROOM),
		CurrentUsers: 1,
		MaxUsers: MAX_USERS_IN_ROOM,
		BannedUsers: make([]int, 0),
		Password: requestedRoom.Password,
	}
	room.Users = append(room.Users, session.UserId)
	rooms = append(rooms, room)
	usersInGame[session.UserId] = true

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(RoomResponse{
		Id: room.Id,
		Name: room.Name,
		PackId: room.PackId,
		UserId: room.UserId,
		Users: room.Users,
		CurrentUsers: room.CurrentUsers,
		MaxUsers: room.MaxUsers,
	})
}

func PutRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	intid, _ := strconv.Atoi(id)

	if intid >= MAX_ROOMS || intid > len(rooms) {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Invalid room id")
		return
	}

	var requestedRoom Room
	err := json.NewDecoder(r.Body).Decode(&requestedRoom)
	if err != nil {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Could not process the body of the request.")
		return
	}

	if strings.TrimSpace(requestedRoom.Name) == "" || requestedRoom.PackId == 0 || 
		requestedRoom.Password == "" {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Missing fields",
			"name and pack_id fields must be specified on PUT request.")
		return
	}

	if requestedRoom.UserId != 0 || requestedRoom.Id != 0 ||
		len(requestedRoom.Users) != 0 || requestedRoom.CurrentUsers != 0 ||
		requestedRoom.MaxUsers != 0 || len(requestedRoom.BannedUsers) != 0 {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Modifying non-editable fields",
			"user_id, users, current_users, max_users, and banned_users fields can't be changed via PUT request.")
		return
	}

	room := rooms[intid - 1]
	session := r.Context().Value("session").(*Session)
	conn := r.Context().Value("db_connection").(*pgx.Conn)

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

	room.Name = requestedRoom.Name
	room.PackId = requestedRoom.PackId
	room.Password = requestedRoom.Password

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(RoomResponse{
		Id: room.Id,
		Name: room.Name,
		PackId: room.PackId,
		UserId: room.UserId,
		Users: room.Users,
		CurrentUsers: room.CurrentUsers,
		MaxUsers: room.MaxUsers,
	})
}

func PatchRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	intid, _ := strconv.Atoi(id)

	if intid >= MAX_ROOMS || intid > len(rooms) {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Invalid room id")
		return
	}

	var requestedRoom Room
	err := json.NewDecoder(r.Body).Decode(&requestedRoom)
	if err != nil {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Could not process the body of the request.")
		return
	}

	if requestedRoom.UserId != 0 || requestedRoom.Id != 0 ||
		len(requestedRoom.Users) != 0 || requestedRoom.CurrentUsers != 0 ||
		requestedRoom.MaxUsers != 0 || len(requestedRoom.BannedUsers) != 0 {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Modifying non-editable fields",
			"room_id, user_id, users, current_users, max_users, and banned_users fields can't be changed via PATCH request.")
		return
	}

	room := rooms[intid - 1]
	session := r.Context().Value("session").(*Session)

	if session.UserId != room.UserId {
		httputils.SendErrorMessage(w, http.StatusForbidden, "Forbidden",
			"Can't modify a room that is not owned by themselves.")
		return
	}

	if requestedRoom.Name != "" {
		room.Name = requestedRoom.Name
	}

	if requestedRoom.Password != "" {
		room.Password = requestedRoom.Password
	}

	if requestedRoom.PackId != 0 {
		conn := r.Context().Value("db_connection").(*pgx.Conn)

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

		room.PackId = requestedRoom.PackId
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(RoomResponse{
		Id: room.Id,
		Name: room.Name,
		PackId: room.PackId,
		UserId: room.UserId,
		Users: room.Users,
		CurrentUsers: room.CurrentUsers,
		MaxUsers: room.MaxUsers,
	})
}

func DeleteRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	intid, _ := strconv.Atoi(id)
	if intid >= MAX_ROOMS || intid > len(rooms) {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Invalid room id")
		return
	}

	room := rooms[intid - 1]
	session := r.Context().Value("session").(*Session)

	if session.UserId != room.UserId {
		httputils.SendErrorMessage(w, http.StatusForbidden, "Forbidden",
			"Can't delete a room that is not owned by themselves.")
		return
	}

	for _, userId := range room.Users {
		delete(usersInGame, userId)
	}

	// deletes item from intid - 1 (included) up to intid (not included)
	// basically deleting the room and leaving free space for other rooms.
	rooms = slices.Delete(rooms, intid - 1, intid)
	w.WriteHeader(http.StatusNoContent)
}

func GetRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	intid, _ := strconv.Atoi(id)
	if intid >= MAX_ROOMS || intid > len(rooms) {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Invalid room id")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	room := rooms[intid - 1]
	json.NewEncoder(w).Encode(RoomResponse{
		Id: room.Id,
		Name: room.Name,
		PackId: room.PackId,
		UserId: room.UserId,
		Users: room.Users,
		CurrentUsers: room.CurrentUsers,
		MaxUsers: room.MaxUsers,
	})
}

func GetRooms(w http.ResponseWriter, r *http.Request) {
	name := strings.ToLower(r.URL.Query().Get("name"))

	var responseRooms []*RoomResponse = make([]*RoomResponse, 0, MAX_ROOMS_RESPONSE)
	for _, room := range rooms {
		if room == nil {
			continue
		}

		if name != "" {
			if strings.ToLower(room.Name) == name {
				responseRooms = append(responseRooms, &RoomResponse{
					Id: room.Id,
					Name: room.Name,
					PackId: room.PackId,
					UserId: room.UserId,
					Users: room.Users,
					CurrentUsers: room.CurrentUsers,
					MaxUsers: room.MaxUsers,
				})
			}
			continue
		}
		responseRooms = append(responseRooms, &RoomResponse{
			Id: room.Id,
			Name: room.Name,
			PackId: room.PackId,
			UserId: room.UserId,
			Users: room.Users,
			CurrentUsers: room.CurrentUsers,
			MaxUsers: room.MaxUsers,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(responseRooms)
}

func JoinRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	intid, _ := strconv.Atoi(id)

	if intid >= MAX_ROOMS || intid > len(rooms) {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Invalid room id")
		return
	}

	session := r.Context().Value("session").(*Session)

	if usersInGame[session.UserId] {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Already in game.")
		return
	}

	var requestedRoom Room
	err := json.NewDecoder(r.Body).Decode(&requestedRoom)
	if err != nil {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Could not process the body of the request.")
		return
	}

	room := rooms[intid - 1]
	if requestedRoom.Password != room.Password {
		httputils.SendErrorMessage(w, http.StatusForbidden, "Forbidden",
			"Wrong password credential.")
		return
	}

	if slices.Contains(room.BannedUsers, session.UserId) {
		httputils.SendErrorMessage(w, http.StatusForbidden, "Forbidden",
			"You were banned from this room.")
		return
	}

	if room.CurrentUsers >= MAX_USERS_IN_ROOM {
		httputils.SendErrorMessage(w, http.StatusServiceUnavailable, "Service unavailable",
			"Users limit reached.")
		return
	}

	room.Users = append(room.Users, session.UserId)
	room.CurrentUsers++
	usersInGame[session.UserId] = true

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(RoomStatusResponse{
		Message: "Joined room.",
	})
}

func LeaveRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	intid, _ := strconv.Atoi(id)

	if intid >= MAX_ROOMS || intid > len(rooms) {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Invalid room id")
		return
	}

	session := r.Context().Value("session").(*Session)

	if !usersInGame[session.UserId] {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Not in game.")
		return
	}

	room := rooms[intid - 1]
	if !slices.Contains(room.Users, session.UserId) {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Not in room.")
		return
	}

	index := slices.Index(room.Users, session.UserId) // may possibly be -1?
	room.Users = slices.Delete(room.Users, index, index + 1)
	delete(usersInGame, session.UserId)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(RoomStatusResponse{
		Message: "Left room.",
	})
}

func BanUserInRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	intid, _ := strconv.Atoi(id)

	if intid >= MAX_ROOMS || intid > len(rooms) {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Invalid room id")
		return
	}

	session := r.Context().Value("session").(*Session)

	if !usersInGame[session.UserId] {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Not in game.")
		return
	}

	room := rooms[intid - 1]
	if session.UserId != room.UserId {
		httputils.SendErrorMessage(w, http.StatusForbidden, "Forbidden",
			"Can't ban a user from a room that is not owned by themselves.")
		return
	}

	var bannedUser User
	err := json.NewDecoder(r.Body).Decode(&bannedUser)
	if err != nil {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Could not process the body of the request.")
		return
	}

	if bannedUser.Id == 0 {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Expected user_id of the user to be banned.")
		return
	}

	if !slices.Contains(room.Users, bannedUser.Id) {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Can't ban a user that's not in room.")
		return
	}

	index := slices.Index(room.Users, bannedUser.Id)
	room.Users = slices.Delete(room.Users, index, index + 1)
	room.BannedUsers = append(room.BannedUsers, bannedUser.Id)
	delete(usersInGame, bannedUser.Id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(RoomStatusResponse{
		Message: "Banned user.",
	})
}

func UnbanUserInRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	intid, _ := strconv.Atoi(id)

	if intid >= MAX_ROOMS || intid > len(rooms) {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Invalid room id")
		return
	}

	session := r.Context().Value("session").(*Session)

	if !usersInGame[session.UserId] {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Not in game.")
		return
	}

	room := rooms[intid - 1]
	if session.UserId != room.UserId {
		httputils.SendErrorMessage(w, http.StatusForbidden, "Forbidden",
			"Can't unban a user from a room that is not owned by themselves.")
		return
	}

	var bannedUser User
	err := json.NewDecoder(r.Body).Decode(&bannedUser)
	if err != nil {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Could not process the body of the request.")
		return
	}

	if bannedUser.Id == 0 {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"Expected user_id of the user to be unbanned.")
		return
	}

	if !slices.Contains(room.BannedUsers, bannedUser.Id) {
		httputils.SendErrorMessage(w, http.StatusBadRequest, "Malformatted request",
			"User is not banned.")
		return
	}

	index := slices.Index(room.BannedUsers, bannedUser.Id)
	room.BannedUsers = slices.Delete(room.BannedUsers, index, index + 1)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(RoomStatusResponse{
		Message: "Unbanned user.",
	})
}
