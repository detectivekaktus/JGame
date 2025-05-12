package main

import (
	"log"
	"net/http"

	"github.com/detectivekaktus/JGame/internal/router"
)

// TODO: Based on the BACKEND_DEV environment variable, set up
// port :80 and enable static file serving from `./website/dist`
// directory for html, css and js files.
//
// Currently the backend API runs on a spearate server from the
// frontend, so no html, css and js is served, only core CRUD logic
// around components (users, rooms and packs).
func main() {
	r := router.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", r))
}
