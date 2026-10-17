package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	JWT      JWTConfig
	Redis    RedisConfig
	Database DatabaseConfig
	Ngrok    NgrokConfig
}

type ServerConfig struct {
	Port string
}
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}
type NgrokConfig struct {
	Upstream  string
	URL       string
	Authtoken string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadConfig() *Config {
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5435"))
	jwtExpHours, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT-SECRET", "your-secret-code"),
			Expiration: time.Duration(jwtExpHours) * time.Hour,
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		Ngrok: NgrokConfig{
			Upstream:  getEnv("NGROK_UPSTREAM", "http://localhost:8080"),
			URL:       getEnv("NGROK_URL", ""),
			Authtoken: getEnv("NGROK_AUTHTOKEN", "Your Key"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "secret"),
			DBName:   getEnv("DB_NAME", "demo"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
