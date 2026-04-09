package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

func init() {
	Load()
}

type Config struct {
	AppEnv    string
	AppName   string
	AppPort   string
	DBHost    string
	DBName    string
	DBPort    string
	RedisHost string
	RedisPort string
	JWTSecret string
}

var (
	AppConfig *Config
	once      sync.Once
)

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func Load() {
	once.Do(func() {
		if AppConfig != nil {
			return // already loaded in main
		}

		// Load .env file
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found.")
		}

		AppConfig = &Config{
			AppEnv:    getEnv("APP_ENV", "development"),
			AppName:   getEnv("APP_NAME", "Banking Application"),
			AppPort:   getEnv("APP_PORT", "8080"),
			DBHost:    getEnv("DB_HOST", "localhost"),
			DBName:    getEnv("DB_NAME", "test"),
			DBPort:    getEnv("DB_PORT", "27017"),
			RedisHost: getEnv("REDIS_HOST", "localhost"),
			RedisPort: getEnv("REDIS_PORT", "6379"),
			JWTSecret: getEnv("JWT_SECRET", "secret"),
		}
	})

	log.Println("Config loaded")
}
