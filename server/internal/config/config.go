package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	redis "github.com/redis/go-redis/v9"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBURL      string
	RedisHost  string
	RedisPort  string
	Redis      *redis.Client
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading environment variables directly...")
	}

	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "wallet_user"),
		DBPassword: getEnv("DB_PASSWORD", "wallet_pass"),
		DBName:     getEnv("DB_NAME", "wallet"),
		RedisHost:  getEnv("REDIS_HOST", "localhost"),
		RedisPort:  getEnv("REDIS_PORT", "6379"),
	}

	cfg.DBURL = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func (c *Config) InitRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr: c.RedisHost + ":" + c.RedisPort,
	})
	c.Redis = rdb
}
