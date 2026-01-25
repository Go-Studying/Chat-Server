package config

import (
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
	godotenv.Load()

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
		log.Fatalf("환경 별수 설정하세요!! %s", key)
	}
	return value
}

func (c *Config) GetDSN() string {
	return "host=" + c.DBHost +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" port=" + c.DBPort +
		" sslmode=" + c.DBSSLMode
}
