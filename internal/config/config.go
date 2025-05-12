package config

import "os"

type Config struct {
	DevMode 	bool
	UserDbUrl string
	PackDbUrl string
}

var AppConfig = load()

func load() *Config {
	c := &Config{
		DevMode: len(os.Getenv("BACKEND_DEV")) > 0,
		UserDbUrl: os.Getenv("USER_DB_URL"),
		PackDbUrl: os.Getenv("PACK_DB_URL"),
	}
	return c
}
