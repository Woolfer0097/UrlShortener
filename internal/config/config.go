package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	StoragePostgres = "postgresql"
	StorageInMemory = "in_memory"
)

type Config struct {
	API      ApiConfig
	Storage  StorageConfig
	Database DatabaseConfig
	Logger   LoggerConfig
}

type StorageConfig struct {
	Type string
}

type ApiConfig struct {
	AppName string
	AppHost string
	AppPort string
	BaseURL string
}

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

type LoggerConfig struct {
	Level string
	Path  string
}

var cfg *Config

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading environment variables")
	}

	viper.AutomaticEnv()

	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_HOST", "0.0.0.0")
	viper.SetDefault("BASE_URL", "http://localhost:8080")
	viper.SetDefault("STORAGE_TYPE", StoragePostgres)

	cfg = &Config{
		API: ApiConfig{
			AppName: viper.GetString("APP_NAME"),
			AppHost: viper.GetString("APP_HOST"),
			AppPort: viper.GetString("APP_PORT"),
			BaseURL: viper.GetString("BASE_URL"),
		},
		Storage: StorageConfig{
			Type: viper.GetString("STORAGE_TYPE"),
		},
		Database: DatabaseConfig{
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			Name:     viper.GetString("DB_NAME"),
		},
		Logger: LoggerConfig{
			Level: viper.GetString("LOG_LEVEL"),
			Path:  viper.GetString("LOG_PATH"),
		},
	}
}

func GetConfig() *Config {
	if cfg == nil {
		Init()
	}
	return cfg
}

func GetPostgresDSN() string {
	c := GetConfig()
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name)
}
