package config

import (
	"fmt"
	"log"
	"os"
	"strings"
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
	url      string
	username string
	password string
}

func (d *DatabaseConfig) ConnectionString() string {
	s := strings.Split(d.url, "postgresql://")
	return fmt.Sprintf("postgresql://%s:%s@%s", d.username, d.password, s[1])
}

func Load() *Config {
	cfg := &Config{
		Server: ServerConfig{
			Addr:         getEnv("server_addr", ":8080"),
			ReadTimeout:  getEnvTime("server_read_timeout", 1*time.Second),
			WriteTimeout: getEnvTime("server_write_timeout", 1*time.Second),
			IdleTimeout:  getEnvTime("server_idle_timeout", 30*time.Second),
		},
		Database: DatabaseConfig{
			url:      os.Getenv("postgres_url"),
			username: os.Getenv("postgres_username"),
			password: os.Getenv("postgres_password"),
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
