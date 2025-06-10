package main

import (
	"log"
	"net/http"

	"github.com/detectivekaktus/JGame/internal/config"
	"github.com/detectivekaktus/JGame/internal/router"
)

func main() {
	r := router.NewRouter()
	
	log.Fatal(http.ListenAndServeTLS(":8080", config.AppConfig.SslCertPath, config.AppConfig.SslKeyPath, r))
}
