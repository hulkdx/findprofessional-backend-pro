package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

// TODO: Change timeouts for production
const (
	DefaultServerPort         = "8081"
	DefaultServerReadTimeout  = 10 * time.Second
	DefaultServerWriteTimeout = 20 * time.Second
	DefaultServerIdleTimeout  = 30 * time.Second
)

type Config struct {
	Server   *ServerConfig
	Database *DatabaseConfig
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
		Server: &ServerConfig{
			Port:         utils.GetEnv("server_port", DefaultServerPort),
			ReadTimeout:  utils.GetEnvTime("server_read_timeout", DefaultServerReadTimeout),
			WriteTimeout: utils.GetEnvTime("server_write_timeout", DefaultServerWriteTimeout),
			IdleTimeout:  utils.GetEnvTime("server_idle_timeout", DefaultServerIdleTimeout),
		},
		Database: LoadDataBaseConfig(),
	}
	return cfg
}

func LoadDataBaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		getDsn(
			os.Getenv("postgres_url"),
			os.Getenv("postgres_username"),
			os.Getenv("postgres_password"),
		),
	}
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
