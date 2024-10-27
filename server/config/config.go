package config

import (
	"log"
	"os"
)

type Config struct {
	GotifyURL   string
	GotifyToken string

	PGConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		DB       string
	}
}

var conf *Config

func LoadConfig() *Config {
	if conf != nil {
		return conf
	}
	gotifyURL := os.Getenv("GOTIFY_URL")
	gotifyToken := os.Getenv("GOTIFY_TOKEN")

	if gotifyURL == "" || gotifyToken == "" {
		log.Fatal("GOTIFY_URL and GOTIFY_TOKEN environment variables must be set")
	}
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	dbname := os.Getenv("PG_DB")

	conf = &Config{}
	conf.GotifyURL = gotifyURL
	conf.GotifyToken = gotifyToken
	conf.PGConfig.DB = dbname
	conf.PGConfig.Host = host
	conf.PGConfig.Port = port
	conf.PGConfig.User = user
	conf.PGConfig.Password = password
	conf.PGConfig.DB = dbname

	return conf
}
