package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Environment struct {
	// Server configs
	ServerPort string
	ServerHost string

	// Database configs
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Application configs
	Environment string
	LogLevel    string
}

// LoadEnv carga las variables de entorno desde el archivo .env si existe
func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		// No retornamos error si el archivo no existe, ya que podríamos estar en producción
		// usando variables de entorno del sistema
		if !os.IsNotExist(err) {
			return fmt.Errorf("error loading .env file: %w", err)
		}
	}
	return nil
}

// NewEnvironment crea una nueva instancia de Environment con los valores por defecto
func NewEnvironment() (*Environment, error) {
	if err := LoadEnv(); err != nil {
		return nil, err
	}

	dbPort, err := strconv.Atoi(getEnvOrDefault("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	return &Environment{
		// Server configs
		ServerPort: getEnvOrDefault("SERVER_PORT", "8080"),
		ServerHost: getEnvOrDefault("SERVER_HOST", "0.0.0.0"),

		// Database configs
		DBHost:     getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnvOrDefault("DB_USER", "postgres"),
		DBPassword: getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:     getEnvOrDefault("DB_NAME", "car_service"),
		DBSSLMode:  getEnvOrDefault("DB_SSL_MODE", "disable"),

		// Application configs
		Environment: getEnvOrDefault("APP_ENV", "development"),
		LogLevel:    getEnvOrDefault("LOG_LEVEL", "info"),
	}, nil
}

// GetDSN retorna el Data Source Name para la conexión a la base de datos
func (e *Environment) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		e.DBHost,
		e.DBPort,
		e.DBUser,
		e.DBPassword,
		e.DBName,
		e.DBSSLMode,
	)
}

// getEnvOrDefault obtiene una variable de entorno o retorna un valor por defecto
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 