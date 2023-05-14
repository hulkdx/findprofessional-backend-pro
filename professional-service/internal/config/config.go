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
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func (s *ServerConfig) Addr() string {
	return fmt.Sprintf(":%s", s.Port)
}

type DatabaseConfig struct {
	url      string
	username string
	password string
}

func (d *DatabaseConfig) Dsn() string {
	url := d.url
	if url == "" {
		panic("Url is not provided")
	}
	split := strings.Split(url, "postgresql://")[1]
	hasSsl := strings.Contains(split, "sslmode=")
	var restUrl string
	if hasSsl {
		restUrl = split
	} else {
		restUrl = fmt.Sprintf("%s?sslmode=disable", split)
	}
	return fmt.Sprintf("postgresql://%s:%s@%s", d.username, d.password, restUrl)
}

func Load() *Config {
	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("server_port", "8081"),
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
