package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DevMode     bool
	DbUrl       string
	SslCertPath string
	SslKeyPath  string
}

var AppConfig = load()

func load() *Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load `.env` file: %v\n", err)
		os.Exit(1)
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		fmt.Fprintf(os.Stderr, "No environment variable DATABASE_URL found. It's either not defined or blank.\n")
		os.Exit(1)
	}

	mode := os.Getenv("MODE")
	if mode == "" {
		fmt.Fprintf(os.Stderr, "No environment variable MODE found. Assuming default `dev`.\n")
		mode = "dev"
	}

	sslCertificatePath := os.Getenv("SSL_CERT_PATH")
	sslKeyPath := os.Getenv("SSL_KEY_PATH")

	if sslCertificatePath == "" || sslKeyPath == "" {
		fmt.Fprintf(os.Stderr, "No SSL certificate or key found. The dev environment must run with HTTPS. Obtain a SSL certificate.")
		os.Exit(1)
	}

	c := &Config{
		DevMode: mode == "dev",
		DbUrl: dbUrl,
		SslCertPath: sslCertificatePath,
		SslKeyPath: sslKeyPath,
	}
	return c
}
