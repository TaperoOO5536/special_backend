package env

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	err := godotenv.Load("pkg/env/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	DSN := os.Getenv("DSN")
	if DSN == "" {
		return fmt.Errorf("DSN is not set")
	}
	token := os.Getenv("TOKEN")
	if token == "" {
		return fmt.Errorf("TOKEN is not set")
	}
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		return fmt.Errorf("GRPC_PORT is not set")
	}
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		return fmt.Errorf("HTTP_PORT is not set")
	}
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		return fmt.Errorf("KAFKA_BORKERS is not set")
	}
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		return fmt.Errorf("ALLOWED_ORIGINS is not set")
	}
	allowedMethods := os.Getenv("ALLOWED_METHODS")
	if allowedMethods == "" {
		return fmt.Errorf("ALLOWED_METHODS is not set")
	}
	allowedHeaders := os.Getenv("ALLOWED_HEADERS")
	if allowedHeaders == "" {
		return fmt.Errorf("ALLOWED_HEADERS is not set")
	}

	return nil
}

func GetDsn() string {
	return os.Getenv("DSN")
}

func GetToken() string {
	return os.Getenv("TOKEN")
}

func GetGRPCPort() string {
	return GetEnvDefault("GRPC_PORT", "8080")
}

func GetHTTPPort() string {
	return GetEnvDefault("HTTP_PORT", "8081")
}

func GetKafkaBrokers() []string {
	return strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
}

func GetAllowedOrigins() []string {
	return strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
}

func GetAllowedMethods() []string {
	return strings.Split(os.Getenv("ALLOWED_METHODS"), ",")
}

func GetAllowedHeaders() []string {
	return strings.Split(os.Getenv("ALLOWED_HEADERS"), ",")
}

func GetEnvDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
