package config

import "os"

type Config struct {
	DevMode 	bool
	DbUrl 		string
}

var AppConfig = load()

func load() *Config {
	c := &Config{
		DevMode: len(os.Getenv("BACKEND_DEV")) > 0,
		DbUrl: os.Getenv("DATABASE_URL"),
	}
	return c
}
