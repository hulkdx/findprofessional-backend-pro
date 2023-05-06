package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func Load() *Config {
	cfg := &Config{
		Server: ServerConfig{
			Addr:         getEnv("SERVER_ADDR", ":8080"),
			ReadTimeout:  getEnvTime("SERVER_READ_TIMEOUT", 1*time.Second),
			WriteTimeout: getEnvTime("SERVER_WRITE_TIMEOUT", 1*time.Second),
			IdleTimeout:  getEnvTime("SERVER_IDLE_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			Driver:   os.Getenv("DB_DRIVER"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
	}
	return cfg
}

func getEnv(key string, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	return value
}

func getEnvTime(key string, def time.Duration) time.Duration {
	str := os.Getenv(key)
	if str == "" {
		return def
	}
	duration, err := time.ParseDuration(str)
	if err != nil {
		log.Fatal("Unable to parse time", err)
	}
	return duration
}
