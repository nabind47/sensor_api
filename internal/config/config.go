package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Primary Primary
	Server  ServerConfig
	Auth    AuthConfig
}

type Primary struct {
	Env string
}

type ServerConfig struct {
	Port               int
	ReadTimeout        int
	WriteTimeout       int
	IdleTimeout        int
	LogLevel           string
	CorsAllowedOrigins []string
}

type AuthConfig struct {
	ClientID     string
	ClientSecret string
	ClientExpiry time.Duration
}

func GetString(key string, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return valAsInt
}

func GetSecret(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return "", fmt.Errorf("missing required secret: %s", key)
	}
	return val, nil
}

func Load() (*Config, error) {
	clientID, err := GetSecret("AUTH_CLIENT_ID")
	if err != nil {
		return nil, err
	}

	clientSecret, err := GetSecret("AUTH_CLIENT_SECRET")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Primary: Primary{
			Env: GetString("APP_ENV", "development"),
		},
		Server: ServerConfig{
			Port:               GetInt("PORT", 8080),
			ReadTimeout:        GetInt("SERVER_READ_TIMEOUT", 10),
			WriteTimeout:       GetInt("SERVER_WRITE_TIMEOUT", 10),
			IdleTimeout:        GetInt("SERVER_IDLE_TIMEOUT", 30),
			LogLevel:           GetString("SERVER_LOG_LEVEL", "info"),
			CorsAllowedOrigins: strings.Split(GetString("CORS_ALLOWED_ORIGINS", ""), ","),
		},
		Auth: AuthConfig{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			// ClientExpiry: GetInt("AUTH_TOKEN_EXPIRY_IN_SECONDS", 120) * time.Second,
			ClientExpiry: time.Duration(GetInt("AUTH_TOKEN_EXPIRY_IN_SECONDS", 120)) * time.Second,
		},
	}

	return cfg, nil
}
