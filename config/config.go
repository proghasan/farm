package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret string
	AppPort   string
}

func Load() *Config {
	godotenv.Load() // ignore error if .env doesn't exist
	port, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))
	return &Config{
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     port,
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "farm"),
		JWTSecret:  getEnv("JWT_SECRET", "change-me-in-production"),
		AppPort:    getEnv("APP_PORT", ":8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
