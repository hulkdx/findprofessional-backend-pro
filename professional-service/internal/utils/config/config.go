package config

import (
	"fmt"
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
	Dsn string
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
			getDsn(
				os.Getenv("postgres_url"),
				os.Getenv("postgres_username"),
				os.Getenv("postgres_password"),
			),
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
		panic(err)
	}
	return duration
}

func getDsn(url string, username string, password string) string {
	if url == "" {
		panic("Url is not provided")
	}
	if !strings.Contains(url, "sslmode=") {
		url = url + "?sslmode=disable"
	}
	restUrl := strings.Split(url, "postgresql://")[1]
	return fmt.Sprintf("postgresql://%s:%s@%s", username, password, restUrl)
}

func IsDebug() bool {
	return os.Getenv("DEBUG") == "true"
}
