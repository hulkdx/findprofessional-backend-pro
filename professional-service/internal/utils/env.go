package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func GetEnv(key string, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	return value
}

func GetEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Environment variables %s not found", key))
	}
	return value
}

func GetEnvIntOrPanic(key string) int {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Environment variables %s not found", key))
	}
	a, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("Environment variables %s is not int", key))
	}
	return a
}

func GetEnvTime(key string, def time.Duration) time.Duration {
	str := os.Getenv(key)
	if str == "" {
		return def
	}
	duration, err := time.ParseDuration(str)
	if err != nil {
		panic(fmt.Sprintf("Environment variables %s err=%s", key, err.Error()))
	}
	return duration
}

func GetEnvTimeOrPanic(key string) time.Duration {
	str := os.Getenv(key)
	if str == "" {
		panic(fmt.Sprintf("Environment variables %s not found", key))
	}
	duration, err := time.ParseDuration(str)
	if err != nil {
		panic(fmt.Sprintf("Environment variables %s err=%s", key, err.Error()))
	}
	return duration
}
