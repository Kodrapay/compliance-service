package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	ServicePort string
	DatabaseURL string
	RedisAddr   string
	RedisPassword string
	RedisDB int
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // Default Redis address
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	
	redisDBStr := os.Getenv("REDIS_DB")
	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		redisDB = 0 // Default Redis DB
	}

	return &Config{
		ServicePort: port,
		DatabaseURL: dbURL,
		RedisAddr:   redisAddr,
		RedisPassword: redisPassword,
		RedisDB: redisDB,
	}
}