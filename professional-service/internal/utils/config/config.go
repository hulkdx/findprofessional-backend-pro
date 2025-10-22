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
	DEFAULT_SERVER_PORT          = "8081"
	DEFAULT_SERVER_READ_TIMEOUT  = 10 * time.Second
	DEFAULT_SERVER_WRITE_TIMEOUT = 20 * time.Second
	DEFAULT_SERVER_IDLE_TIMEOUT  = 30 * time.Second
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
			Port:         utils.GetEnv("server_port", DEFAULT_SERVER_PORT),
			ReadTimeout:  utils.GetEnvTime("server_read_timeout", DEFAULT_SERVER_READ_TIMEOUT),
			WriteTimeout: utils.GetEnvTime("server_write_timeout", DEFAULT_SERVER_WRITE_TIMEOUT),
			IdleTimeout:  utils.GetEnvTime("server_idle_timeout", DEFAULT_SERVER_IDLE_TIMEOUT),
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
