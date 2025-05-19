package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DevMode 	bool
	DbUrl 		string
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

	c := &Config{
		DevMode: len(os.Getenv("BACKEND_DEV")) > 0,
		DbUrl: dbUrl,
	}
	return c
}
