package websocket

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/detectivekaktus/JGame/internal/database"
	"github.com/detectivekaktus/JGame/internal/handler"
	"github.com/detectivekaktus/JGame/internal/httputils"
	"github.com/detectivekaktus/JGame/internal/middleware"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
)

type ActionType string
type UserRole   string

const (
	OWNER  UserRole = "owner"
	PLAYER UserRole = "player"
)

const (
	JOIN_ROOM      ActionType = "join_room"
	JOINED_ROOM    ActionType = "joined_room"

	LEAVE_ROOM     ActionType = "leave_room"
	LEFT_ROOM      ActionType = "left_room"

	DELETE_ROOM    ActionType = "delete_room"
	ROOM_DELETED   ActionType = "room_deleted"

	START_GAME     ActionType = "start_game"
	GAME_STARTED   ActionType = "game_started"

	GET_USERS      ActionType = "get_users"
	USERS_LIST     ActionType = "users_list"

	ERROR          ActionType = "error"
)

type WSMessage struct {
	Type    ActionType     `json:"type"`
	Payload map[string]any `json:"payload"`
}

type User struct {
	Id     int
	Role   UserRole
	RoomId int
}

type Room struct {
	handler.Room
	Started     bool
	Users       []*User
	BannedUsers []*User
}

var rooms map[int]*Room = make(map[int]*Room)

var upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize: 2 << 9,
	WriteBufferSize: 2 << 9,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return middleware.AllowedCorsOrigins[origin]
	},
}

func sendMessage(conn *websocket.Conn, msg WSMessage) error {
	out, _ := json.Marshal(msg)
	err := conn.WriteMessage(websocket.TextMessage, out)
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			fmt.Fprintf(os.Stderr, "Writing websocket message went wrong: %v\n", err)
		}
		return err
	}
	return nil
}

func sendError(conn *websocket.Conn, code int, msg string) error {
	return sendMessage(conn, WSMessage{
		Type: ERROR,
		Payload: map[string]any{
			"code": code,
			"message": msg,
		},
	})
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not upgrade connection to a websocket: %v\n", err)
		httputils.SendErrorMessage(w, http.StatusInternalServerError, "Internal error",
			"Could not upgrade the request to a websocket.")
		return
	}
	defer conn.Close()

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Fprintf(os.Stderr, "Reading websocket message went wrong: %v\n", err)
			} 
			return	
		}
		var msg WSMessage
		json.Unmarshal(raw, &msg)
		roomId, ok := msg.Payload["room_id"].(int)
		if !ok {
			err = sendError(conn, 400, "missing room_id")
			if err != nil {
				return
			}
			continue
		}

		switch msg.Type {
		case JOIN_ROOM: {
			dbConn := r.Context().Value("db_connection").(*pgx.Conn)
			session := r.Context().Value("session").(*handler.Session)

			var user User
			err = database.QueryRow(dbConn, "SELECT * FROM rooms.player WHERE user_id = $1", session.UserId).
				Scan(&user.Id, &user.RoomId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not check user game status: %v\n", err)
				err = sendError(conn, 500, "internal server error")
				if err != nil {
					return
				}
				conn.Close()
				return
			}

			if user.RoomId != 0 && user.RoomId != roomId {
				err = sendError(conn, 400, "already in game")
				if err != nil {
					return
				}
				conn.Close()
				return
			}

			room, ok := rooms[roomId]
			if !ok {
				var room Room
				err = database.QueryRow(dbConn, "SELECT room_id, user_id, name, pack_id, current_users, max_users FROM rooms.room WHERE room_id = $1", roomId).
					Scan(&room.Id, &room.UserId, &room.Name, &room.PackId, &room.CurrentUsers, &room.MaxUsers)
				if err != nil {
					if err == sql.ErrNoRows {
						err = sendError(conn, 404, "no room with this id exists.")
						if err != nil {
							return
						}
						conn.Close()
						return
					}
					fmt.Fprintf(os.Stderr, "Could not get the room: %v\n", err)
					err = sendError(conn, 500, "internal server error")
					if err != nil {
						return
					}
					conn.Close()
					return
				}

				_, err = database.Execute(dbConn, "INSERT INTO rooms.player (user_id, room_id) VALUES ($1, $2)", session.UserId, roomId)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Could not insert into the player table: %v\n", err)
					err = sendError(conn, 500, "internal server error")
					if err != nil {
						return
					}
					conn.Close()
					return
				}

				room.Users = append(room.Users, &User{
					RoomId: roomId,
					Id: session.UserId,
					Role: OWNER,
				})
				rooms[roomId] = &room
				err := sendMessage(conn, WSMessage{
					Type: JOINED_ROOM,
					Payload: map[string]any{
						"user_id": session.UserId,
						"role": OWNER,
					},
				})
				if err != nil {
					return
				}
			} else {
				_, err = database.Execute(dbConn, "INSERT INTO rooms.player (user_id, room_id) VALUES ($1, $2)", session.UserId, roomId)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Could not insert into the player table: %v\n", err)
					err = sendError(conn, 500, "internal server error")
					if err != nil {
						return
					}
					conn.Close()
					return
				}

				room.Users = append(room.Users, &User{
					RoomId: roomId,
					Id: session.UserId,
					Role: PLAYER,
				})
				err := sendMessage(conn, WSMessage{
					Type: JOINED_ROOM,
					Payload: map[string]any{
						"user_id": session.UserId,
						"role": PLAYER,
					},
				})
				if err != nil {
					return
				}
			}
		}

		case START_GAME: {
			session := r.Context().Value("session").(*handler.Session)
			room := rooms[roomId]

			if session.UserId != room.UserId {
				err = sendError(conn, 403, "only owner can start the game")
				if err != nil {
					return
				}
			}

			err := sendMessage(conn, WSMessage{ Type: GAME_STARTED, })
			if err != nil {
				return
			}
		}

		case GET_USERS: {
			err := sendMessage(conn, WSMessage{
				Type: USERS_LIST,
				Payload: map[string]any{
					"users": rooms[roomId].Users,
				},
			})
			if err != nil {
				return
			}
		}

		default: {
			err = sendError(conn, 400, "unknown action")
			if err != nil {
				return
			}
		}
		}
	}
}

