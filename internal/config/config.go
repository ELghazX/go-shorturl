package config

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	DB    *sql.DB
	Redis *redis.Client
}

func New() (*Config, error) {
	_ = godotenv.Load()
	
	db, err := setupDatabase()
	if err != nil {
		return nil, err
	}

	redisClient := setupRedis()

	return &Config{
		DB:    db,
		Redis: redisClient,
	}, nil
}

func setupDatabase() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "password"),
		getEnv("DB_NAME", "urlshortener"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		short_code VARCHAR(10) UNIQUE NOT NULL,
		long_url TEXT NOT NULL,
		clicks INTEGER DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err = db.Exec(query)
	if err != nil {
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func setupRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: getEnv("REDIS_HOST", "localhost") + ":" + getEnv("REDIS_PORT", "6379"),
	})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
