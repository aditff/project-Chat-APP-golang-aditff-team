package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RedisClient *redis.Client
var Ctx = context.Background()

// LoadEnv loads environment variables from .env file
func LoadEnv() {
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: .env file not found, using system environment variables")
    }
}

// InitPostgres initializes PostgreSQL connection
func InitPostgres() {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        getEnv("POSTGRES_HOST", "localhost"),
        getEnv("POSTGRES_USER", "postgres"),
        getEnv("POSTGRES_PASSWORD", "password"),
        getEnv("POSTGRES_DB", "chatapp"),
        getEnv("POSTGRES_PORT", "5432"),
    )

    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to PostgreSQL: %v", err)
    }

    log.Println("PostgreSQL connected successfully")
}

// InitRedis initializes Redis client
func InitRedis() {
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     getEnv("REDIS_HOST", "localhost") + ":" + getEnv("REDIS_PORT", "6379"),
        Password: getEnv("REDIS_PASSWORD", ""), 
        DB:       0,
    })

    _, err := RedisClient.Ping(Ctx).Result()
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }

    log.Println("Redis connected successfully")
}

// getEnv reads an environment variable or returns default value
func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}
