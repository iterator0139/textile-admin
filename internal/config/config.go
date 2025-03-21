package config

import (
	"os"
	"path/filepath"
	"strconv"

	"textile-admin/pkg/db"
)

// Config holds the application configuration
type Config struct {
	ServerAddress string
	UploadDir     string
	DBConfig      db.DBConfig
	FileURLPrefix string
}

// LoadConfig loads the application configuration from environment variables
func LoadConfig() Config {
	// Default values
	cfg := Config{
		ServerAddress: ":8080",
		UploadDir:     "uploads",
		DBConfig: db.DBConfig{
			Host:     "localhost",
			Port:     3306,
			User:     "root",
			Password: "password",
			DBName:   "textile_admin",
		},
		FileURLPrefix: "http://localhost:8080/files",
	}

	// Override with environment variables if they exist
	if serverAddr := os.Getenv("SERVER_ADDRESS"); serverAddr != "" {
		cfg.ServerAddress = serverAddr
	}

	if uploadDir := os.Getenv("UPLOAD_DIR"); uploadDir != "" {
		cfg.UploadDir = uploadDir
	}

	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		cfg.DBConfig.Host = dbHost
	}

	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		if port, err := strconv.Atoi(dbPort); err == nil {
			cfg.DBConfig.Port = port
		}
	}

	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		cfg.DBConfig.User = dbUser
	}

	if dbPass := os.Getenv("DB_PASSWORD"); dbPass != "" {
		cfg.DBConfig.Password = dbPass
	}

	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		cfg.DBConfig.DBName = dbName
	}

	if fileURLPrefix := os.Getenv("FILE_URL_PREFIX"); fileURLPrefix != "" {
		cfg.FileURLPrefix = fileURLPrefix
	}

	// Ensure upload directory exists
	if _, err := os.Stat(cfg.UploadDir); os.IsNotExist(err) {
		os.MkdirAll(cfg.UploadDir, 0755)
	}

	// Always convert upload path to absolute path
	absUploadDir, err := filepath.Abs(cfg.UploadDir)
	if err == nil {
		cfg.UploadDir = absUploadDir
	}

	return cfg
} 