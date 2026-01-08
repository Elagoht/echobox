package config

import (
	"os"
	"strconv"
)

const (
	DefaultPort         = "5867"
	DefaultReadTimeout  = 30
	DefaultWriteTimeout = 30
)

type Server struct {
	Port         string
	ReadTimeout  int
	WriteTimeout int
}

func Load() *Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	readTimeout := getEnvInt("READ_TIMEOUT", DefaultReadTimeout)
	writeTimeout := getEnvInt("WRITE_TIMEOUT", DefaultWriteTimeout)

	return &Server{
		Port:         port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

func getEnvInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}
