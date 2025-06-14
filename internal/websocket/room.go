package websocket

import (
	"fmt"
	"net/http"
	"os"

	"github.com/detectivekaktus/JGame/internal/httputils"
	"github.com/detectivekaktus/JGame/internal/middleware"
	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize: 2 << 9,
	WriteBufferSize: 2 << 9,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return middleware.AllowedCorsOrigins[origin]
	},
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
		_, msg, err := conn.ReadMessage()
		if err != nil && websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			fmt.Fprintf(os.Stderr, "Reading websocket message went wrong: %v\n", err)
			return
		}
		fmt.Print(msg)
		
		err = conn.WriteMessage(websocket.TextMessage, msg);
		if err != nil && websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			fmt.Fprintf(os.Stderr, "Writing websocket message went wrong: %v\n", err)
			return
		}
	}
}

