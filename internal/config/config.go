package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	JWTSecret   string
	Environment string
}

func Load() *Config {


	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %s", err)
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnvRequire("DB_USER"),
		DBPassword:  getEnvRequire("DB_PASSWORD"),
		DBName:      getEnvRequire("DB_NAME"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
		JWTSecret:   getEnvRequire("JWT_SECRET"),
		Environment: getEnv("ENV", "development"),
	}
}

// 기본값 있는 환경 변수
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// 필수 환경 변수
func getEnvRequire(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("필수 환경 변수가 설정되지 않았습니다: %s", key)
	}
	return value
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort, c.DBSSLMode)
}
