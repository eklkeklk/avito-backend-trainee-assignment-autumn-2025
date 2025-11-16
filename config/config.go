package config

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewConnection(cfg *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func GetConfig() *Config {
	return &Config{
		Host:     getEnvInfo("DB_HOST", "localhost"),
		Port:     getEnvInfo("DB_PORT", "5432"),
		User:     getEnvInfo("DB_USER", "postgres"),
		Password: getEnvInfo("DB_PASSWORD", "password"),
		DBName:   getEnvInfo("DB_NAME", "avito_db"),
		SSLMode:  getEnvInfo("SSL_MODE", "disable"),
	}
}

func GetPort() string {
	return ":" + getEnvInfo("PORT", "8080")
}

func getEnvInfo(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
