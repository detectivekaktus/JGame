package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/detectivekaktus/JGame/internal/config"
	"github.com/detectivekaktus/JGame/internal/database"
	"github.com/detectivekaktus/JGame/internal/router"
)

func main() {
	r := router.NewRouter()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Caught SIGINT signal. Cleaning up database before exiting...")

		conn := database.GetConnection()
		database.Execute(conn, "DELETE FROM rooms.player")
		database.Execute(conn, "DELETE FROM rooms.room")
		conn.Close(context.Background())

		os.Exit(0)
	}()
	
	log.Fatal(http.ListenAndServeTLS(":8080", config.AppConfig.SslCertPath, config.AppConfig.SslKeyPath, r))
}
