package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/may20xx/booking/pkg/log"
)

type Config struct {
	Port string

	DBPort     string
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret           string
	JWTRefreshSecret    string
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string

	MailHost string
	MailPort int
	MailUser string
	MailPass string
}

var config *Config

func init() {
	if err := godotenv.Load(); err != nil {
		log.Msg.Errorf("Error loading .env file: %s\n", err)
	}

	config = loadConfig()
}

func loadConfig() *Config {
	return &Config{
		Port:                getEnv("PORT", "8080"),
		DBPort:              getEnv("DB_PORT", "5432"),
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBUser:              getEnv("DB_USERNAME", "postgres"),
		DBPassword:          getEnv("DB_PASSWORD", "postgres"),
		DBName:              getEnv("DB_NAME", "postgres"),
		JWTSecret:           getEnvMustExist("JWT_SECRET"),
		JWTRefreshSecret:    getEnvMustExist("JWT_REFRESH_SECRET"),
		CloudinaryCloudName: getEnvMustExist("CLOUDINARY_CLOUD_NAME"),
		CloudinaryAPIKey:    getEnvMustExist("CLOUDINARY_API_KEY"),
		CloudinaryAPISecret: getEnvMustExist("CLOUDINARY_API_SECRET"),
		MailHost:            getEnvMustExist("MAIL_HOST"),
		MailPort:            587,
		MailUser:            getEnvMustExist("MAIL_USER"),
		MailPass:            getEnvMustExist("MAIL_PASS"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvMustExist(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Msg.Fatal(fmt.Sprintf("%s must be set", key))
	}
	return value
}

func GetConfig() *Config {
	return config
}
