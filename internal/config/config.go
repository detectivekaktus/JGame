package config

import (
	"fmt"
	"os"
)

type Config struct {
	DevMode 	bool
	DbUrl 		string
}

var AppConfig = load()

func load() *Config {
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
