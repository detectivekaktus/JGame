package main

import (
	"log"
	"net/http"

	"github.com/detectivekaktus/JGame/internal/router"
)

func main() {
	r := router.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", r))
}
