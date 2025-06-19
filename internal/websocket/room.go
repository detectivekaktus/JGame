package websocket

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"

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
	ROOM_DELETED   ActionType = "room_deleted"

	START_GAME     ActionType = "start_game"
	GAME_STARTED   ActionType = "game_started"

	GET_USERS      ActionType = "get_users"
	USERS_LIST     ActionType = "users_list"

	GET_GAME_STATE ActionType = "get_game_state"
	GAME_STATE     ActionType = "game_state"

	NEXT_QUESTION  ActionType = "next_question"
	QUESTION       ActionType = "question"
	QUESTIONS_DONE ActionType = "questions_done"

	ANSWER         ActionType = "answer"

	ERROR          ActionType = "error"
)

type WSMessage struct {
	Type    ActionType     `json:"type"`
	Payload map[string]any `json:"payload"`
}

type User struct {
	Id     int      `json:"id"`
	Name   string   `json:"name"`
	Role   UserRole `json:"role"`
	Score  int      `json:"score"`

	RoomId int      `json:"room_id"`
}

type PackQuestion struct {
	Title   string       `json:"title"`
	ImgUrl  string       `json:"image_url"`
	Value   int          `json:"value"`
	Answers []PackAnswer `json:"answers"`
}

type PackAnswer struct {
	Text    string `json:"text"`
	Correct bool   `json:"correct"`
}

type Pack struct {
	Title           string         `json:"title"`
	Questions       []PackQuestion `json:"questions"`
	CurrentQuestion int
}

type Room struct {
	handler.Room
	Started         bool
	Finished        bool
	Users           map[int]*User
	BannedUsers     map[int]*User

	Pack            Pack

	Connections     map[int]*websocket.Conn
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
		fmt.Println(msg)
		roomIdFloat, ok := msg.Payload["room_id"].(float64)
		if !ok {
			err = sendError(conn, 400, "missing room_id")
			if err != nil {
				return
			}
			continue
		}
		roomId := int(roomIdFloat)

		switch msg.Type {
		case JOIN_ROOM: {
			dbConn := r.Context().Value("db_connection").(*pgx.Conn)
			session := r.Context().Value("session").(*handler.Session)

			var user User
			err = database.QueryRow(dbConn, "SELECT * FROM rooms.player WHERE user_id = $1", session.UserId).
				Scan(&user.Id, &user.RoomId)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					fmt.Fprintf(os.Stderr, "Could not check user game status: %v\n", err)
					err = sendError(conn, 500, "internal server error")
					if err != nil {
						return
					}
					conn.Close()
					return
				}
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

				var rawPackBody json.RawMessage
				err = database.QueryRow(dbConn, "SELECT body FROM packs.pack WHERE pack_id = $1", room.PackId).
					Scan(&rawPackBody)

				json.Unmarshal(rawPackBody, &room.Pack)
				room.Pack.CurrentQuestion = 0

				room.Users = make(map[int]*User)
				room.Users[session.UserId] = &User{
					RoomId: roomId,
					Id: session.UserId,
					Role: OWNER,
				}

				room.Connections = make(map[int]*websocket.Conn)
				room.Connections[session.UserId] = conn

				room.Finished = false

				room.BannedUsers = make(map[int]*User)

				rooms[roomId] = &room

				_, err = database.Execute(dbConn, "UPDATE rooms.room SET current_users = $1", len(room.Users))
				if err != nil {
					fmt.Fprintf(os.Stderr, "Could not update room: %v\n", err)
					err = sendError(conn, 500, "internal server error")
					if err != nil {
						return
					}
					conn.Close()
					return
				}

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

				rows := database.QueryRows(dbConn, "SELECT u.user_id, u.name FROM rooms.player p JOIN users.\"user\" u ON p.user_id = u.user_id WHERE p.room_id = $1", roomId)

				var users []User
				for rows.Next() {
					var user User
					err := rows.Scan(&user.Id, &user.Name)
					if err != nil {
						sendError(conn, 500, "could not get users")
						return
					}

					user.Role = room.Users[user.Id].Role
					user.Score = room.Users[user.Id].Score
					user.RoomId = roomId

					users = append(users, user)
				}

				for _, c := range room.Connections {
					sendMessage(c, WSMessage{
						Type: USERS_LIST,
						Payload: map[string]any{
							"users": users,
						},
					})
				}
			} else {
				if room.CurrentUsers >= handler.MAX_USERS_IN_ROOM {
					err = sendError(conn, 503, "max users reached")
					if err != nil {
						return
					}
				}

				user, ok := room.Users[session.UserId]
				if ok {
					sendMessage(conn, WSMessage{
						Type: JOINED_ROOM,
						Payload: map[string]any{
							"role": user.Role,
							"user_id": user.Id,
						},
					})
					continue
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

				room.Users[session.UserId] = &User{
					RoomId: roomId,
					Id: session.UserId,
					Role: PLAYER,
				}
				room.Connections[session.UserId] = conn

				_, err = database.Execute(dbConn, "UPDATE rooms.room SET current_users = $1", len(room.Users))
				if err != nil {
					fmt.Fprintf(os.Stderr, "Could not update room: %v\n", err)
					err = sendError(conn, 500, "internal server error")
					if err != nil {
						return
					}
					conn.Close()
					return
				}

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

				rows := database.QueryRows(dbConn, "SELECT u.user_id, u.name FROM rooms.player p JOIN users.\"user\" u ON p.user_id = u.user_id WHERE p.room_id = $1", roomId)

				var users []User
				for rows.Next() {
					var user User
					err := rows.Scan(&user.Id, &user.Name)
					if err != nil {
						sendError(conn, 500, "could not get users")
						return
					}

					user.Role = room.Users[user.Id].Role
					user.Score = room.Users[user.Id].Score
					user.RoomId = roomId

					users = append(users, user)
				}
				for _, c := range room.Connections {
					sendMessage(c, WSMessage{
						Type: USERS_LIST,
						Payload: map[string]any{
							"users": users,
						},
					})
				}
			}
		}

		case LEAVE_ROOM: {
			dbConn := r.Context().Value("db_connection").(*pgx.Conn)
			session := r.Context().Value("session").(*handler.Session)

			room := rooms[roomId]

			if room.Users[session.UserId].Role == OWNER {
				_, err = database.Execute(dbConn, "DELETE FROM rooms.player WHERE room_id = $1", roomId)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Could not delete user state: %v\n", err)
					err = sendError(conn, 500, "internal server error")
					if err != nil {
						return
					}
					conn.Close()
					return
				}

				_, err = database.Execute(dbConn, "DELETE FROM rooms.room WHERE room_id = $1", roomId)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Could not delete room: %v\n", err)
					err = sendError(conn, 500, "internal server error")
					if err != nil {
						return
					}
					conn.Close()
					return
				}

				for _, c := range room.Connections {
					sendMessage(c, WSMessage{
						Type: ROOM_DELETED,
					})
					c.Close()
				}
				delete(rooms, roomId)
			} else {
				_, err = database.Execute(dbConn, "DELETE FROM rooms.player WHERE user_id = $1", session.UserId)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Could not delete player state: %v\n", err)
					err = sendError(conn, 500, "internal server error")
					if err != nil {
						return
					}
					conn.Close()
					return
				}

				delete(room.Users, session.UserId)

				_, err = database.Execute(dbConn, "UPDATE rooms.room SET current_users = $1", len(room.Users))
				if err != nil {
					fmt.Fprintf(os.Stderr, "Could not update room: %v\n", err)
					err = sendError(conn, 500, "internal server error")
					if err != nil {
						return
					}
					conn.Close()
					return
				}

				sendMessage(conn, WSMessage{ Type: LEFT_ROOM, })

				rows := database.QueryRows(dbConn, "SELECT u.user_id, u.name FROM rooms.player p JOIN users.\"user\" u ON p.user_id = u.user_id WHERE p.room_id = $1", roomId)

				var users []User
				for rows.Next() {
					var user User
					err := rows.Scan(&user.Id, &user.Name)
					if err != nil {
						sendError(conn, 500, "could not get users")
						return
					}

					user.Role = room.Users[user.Id].Role
					user.Score = room.Users[user.Id].Score
					user.RoomId = roomId

					users = append(users, user)
				}

				for _, c := range room.Connections {
					sendMessage(c, WSMessage{
						Type: USERS_LIST,
						Payload: map[string]any{
							"users": users,
						},
					})
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

			room.Started = true

			for _, c := range room.Connections {
				err := sendMessage(c, WSMessage{ Type: GAME_STARTED, })
				if err != nil {
					return
				}
			}
		}

		case GET_USERS: {
			room := rooms[roomId]

			dbConn := r.Context().Value("db_connection").(*pgx.Conn)
			rows := database.QueryRows(dbConn, "SELECT u.user_id, u.name FROM rooms.player p JOIN users.\"user\" u ON p.user_id = u.user_id WHERE p.room_id = $1", roomId)

			var users []User
			for rows.Next() {
				var user User
				err := rows.Scan(&user.Id, &user.Name)
				if err != nil {
					sendError(conn, 500, "could not get users")
					return
				}

				user.Role = room.Users[user.Id].Role
				user.Score = room.Users[user.Id].Score
				user.RoomId = roomId

				users = append(users, user)
			}

			err := sendMessage(conn, WSMessage{
				Type: USERS_LIST,
				Payload: map[string]any{
					"users": users,
				},
			})
			if err != nil {
				return
			}
		}

		case GET_GAME_STATE: {
			room := rooms[roomId]

			err := sendMessage(conn, WSMessage{
				Type: GAME_STATE,
				Payload: map[string]any{
					"started": room.Started,
					"finished": room.Finished,
				},
			})
			if err != nil {
				return
			}
		}

		case NEXT_QUESTION: {
			room := rooms[roomId]
			dbConn := r.Context().Value("db_connection").(*pgx.Conn)
			session := r.Context().Value("session").(*handler.Session)

			if room.UserId != session.UserId {
				err := sendError(conn, 403, "next question can be invoked only by the room owner")
				if err != nil {
					return
				}
			}

			if room.Pack.CurrentQuestion >= len(room.Pack.Questions) {
				room.Finished = true
				for _, c := range room.Connections {
					err := sendMessage(c, WSMessage{ Type: QUESTIONS_DONE, })
					if err != nil {
						return
					}
				}
				rows := database.QueryRows(dbConn, "SELECT u.user_id, u.name FROM rooms.player p JOIN users.\"user\" u ON p.user_id = u.user_id WHERE p.room_id = $1", roomId)

				var users []User
				for rows.Next() {
					var user User
					err := rows.Scan(&user.Id, &user.Name)
					if err != nil {
						sendError(conn, 500, "could not get users")
						return
					}

					user.Role = room.Users[user.Id].Role
					user.Score = room.Users[user.Id].Score
					user.RoomId = roomId

					users = append(users, user)
				}
				sort.Slice(users, func(i, j int) bool {
					return users[i].Score > users[j].Score
				})

				winner := users[0]
				
				_, err := database.Execute(dbConn, "UPDATE users.\"user\" SET matches_won = matches_won + 1 WHERE user_id = $1", winner.Id)
				if err != nil {
					fmt.Printf("Could not register a win: %v", err)
					return
				}

				for _, user := range users {
					database.Execute(dbConn, "UPDATE users.\"user\" SET matches_played = matches_played + 1 WHERE user_id = $1", user.Id)
				}

				continue
			}

			question := room.Pack.Questions[room.Pack.CurrentQuestion]
			room.Pack.CurrentQuestion++

			for _, c := range room.Connections {
				err := sendMessage(c, WSMessage{
					Type: QUESTION,
					Payload: map[string]any{
						"question": question,
					},
				})
				if err != nil {
					return
				}
			}
		}

		case ANSWER: {
			room := rooms[roomId]
			dbConn := r.Context().Value("db_connection").(*pgx.Conn)
			session := r.Context().Value("session").(*handler.Session)

			floatAnswer, ok := msg.Payload["answer"].(float64)
			if !ok {
				err := sendError(conn, 400, "expected answer to be given.")
				if err != nil {
					return
				}
			}
			answer := int(floatAnswer)

			if room.Pack.CurrentQuestion == 0 {
				err := sendError(conn, 400, "no question has been successfully played yet.")
				if err != nil {
					return
				}
			}

			question := room.Pack.Questions[room.Pack.CurrentQuestion - 1]
			for i, a := range question.Answers {
				if answer == i && a.Correct {
					room.Users[session.UserId].Score += question.Value
				}
			}

			rows := database.QueryRows(dbConn, "SELECT u.user_id, u.name FROM rooms.player p JOIN users.\"user\" u ON p.user_id = u.user_id WHERE p.room_id = $1", roomId)

			var users []User
			for rows.Next() {
				var user User
				err := rows.Scan(&user.Id, &user.Name)
				if err != nil {
					sendError(conn, 500, "could not get users")
					return
				}

				user.Role = room.Users[user.Id].Role
				user.Score = room.Users[user.Id].Score
				user.RoomId = roomId

				users = append(users, user)
			}

			for _, c := range room.Connections {
				sendMessage(c, WSMessage{
					Type: USERS_LIST,
					Payload: map[string]any{
						"users": users,
					},
				})
			}
		}

		default: {
			fmt.Println(msg)
			err = sendError(conn, 400, "unknown action")
			if err != nil {
				return
			}
		}
		}
	}
}

