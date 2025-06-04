package handler

// -- DESING NOTE -------------
// I made decision to store rooms and room related information
// directly in RAM to optimize read and write operations on the
// objects. To prevent the overuse of RAM resources, there's
// MAX_ROOMS constant which defines the maximum number of rooms
// that can exist at the same point in time. The O(n) operations
// on arrays (slices) are not the problem, since there are not
// that many elements.

type Room struct {
	Id           int    `json:"room_id"`
	Name         string `json:"name"`
	PackId       int    `json:"pack_id"`
	UserId       int    `json:"user_id"`
	Users        []int  `json:"users"`
	CurrentUsers int    `json:"current_users"`
	MaxUsers     int    `json:"max_users"`
}

const (
	MAX_USERS_IN_ROOM = 2 << 4
	MAX_ROOMS         = 2 << 10
)

var usersInGame map[int]bool
var rooms []Room = make([]Room, 0, MAX_ROOMS)

