package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"textile-admin/pkg/db"

	"gopkg.in/yaml.v3"
)

// Config holds all the application configuration
type Config struct {
	// Server configuration
	ServerAddress  string
	UploadDir      string
	FileURLPrefix  string

	// Database configuration
	DBConfig db.DBConfig

	// Logging configuration
	LogLevelValue  string
	LogFormatValue string

	// Reference to the original YAML config
	yamlConfig *YAMLConfig
}

// ServerConfig represents server configuration in YAML
type ServerConfig struct {
	Address       string `yaml:"address"`
	UploadDir     string `yaml:"upload_dir"`
	FileURLPrefix string `yaml:"file_url_prefix"`
}

// DatabaseConfig represents database configuration in YAML
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

// LogConfig represents logging configuration in YAML
type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// YAMLConfig represents the root configuration structure in YAML
type YAMLConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Log      LogConfig      `yaml:"log"`
}

// LogLevel returns the configured log level
func (c *Config) LogLevel() string {
	if c.yamlConfig != nil && c.yamlConfig.Log.Level != "" {
		return c.yamlConfig.Log.Level
	}
	return c.LogLevelValue
}

// LogFormat returns the configured log format
func (c *Config) LogFormat() string {
	if c.yamlConfig != nil && c.yamlConfig.Log.Format != "" {
		return c.yamlConfig.Log.Format
	}
	return c.LogFormatValue
}

// LoadConfig loads the application configuration from YAML and environment variables
func LoadConfig() Config {
	// Determine environment: dev or prod
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev" // Default to development
	}

	// Load configuration from YAML file
	yamlConfig, err := loadYAMLConfig(env)
	if err != nil {
		fmt.Printf("Warning: Failed to load YAML config: %v, using defaults\n", err)
	}

	// Create and initialize the config with default values
	cfg := Config{
		ServerAddress: ":8080",
		UploadDir:     "uploads",
		FileURLPrefix: "http://localhost:8080/files",
		DBConfig: db.DBConfig{
			Host:     "localhost",
			Port:     3306,
			User:     "root",
			Password: "",
			DBName:   "textile_admin",
		},
		LogLevelValue:  "info",
		LogFormatValue: "text",
		yamlConfig:     yamlConfig,
	}

	// If YAML config was loaded successfully, use its values
	if yamlConfig != nil {
		// Set server config
		if yamlConfig.Server.Address != "" {
			cfg.ServerAddress = yamlConfig.Server.Address
		}
		if yamlConfig.Server.UploadDir != "" {
			cfg.UploadDir = yamlConfig.Server.UploadDir
		}
		if yamlConfig.Server.FileURLPrefix != "" {
			cfg.FileURLPrefix = yamlConfig.Server.FileURLPrefix
		}

		// Set database config
		if yamlConfig.Database.Host != "" {
			cfg.DBConfig.Host = yamlConfig.Database.Host
		}
		if yamlConfig.Database.Port != 0 {
			cfg.DBConfig.Port = yamlConfig.Database.Port
		}
		if yamlConfig.Database.User != "" {
			cfg.DBConfig.User = yamlConfig.Database.User
		}
		if yamlConfig.Database.Password != "" {
			cfg.DBConfig.Password = yamlConfig.Database.Password
		}
		if yamlConfig.Database.DBName != "" {
			cfg.DBConfig.DBName = yamlConfig.Database.DBName
		}

		// Set logging config
		if yamlConfig.Log.Level != "" {
			cfg.LogLevelValue = yamlConfig.Log.Level
		}
		if yamlConfig.Log.Format != "" {
			cfg.LogFormatValue = yamlConfig.Log.Format
		}
	}

	// Override with environment variables if they exist
	processEnvVars(&cfg)

	// Ensure upload directory exists and is absolute
	absUploadDir, err := filepath.Abs(cfg.UploadDir)
	if err == nil {
		cfg.UploadDir = absUploadDir
	}

	return cfg
}

// loadYAMLConfig loads configuration from the appropriate YAML file
func loadYAMLConfig(env string) (*YAMLConfig, error) {
	configFile := fmt.Sprintf("configs/config.%s.yaml", env)
	
	// Read YAML file
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %v", err)
	}

	// Parse YAML
	var config YAMLConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("could not parse config file: %v", err)
	}

	return &config, nil
}

// processEnvVars overrides configuration with environment variables
func processEnvVars(cfg *Config) {
	// Process environment variables for server settings
	if val := os.Getenv("SERVER_ADDRESS"); val != "" {
		cfg.ServerAddress = val
	}
	if val := os.Getenv("UPLOAD_DIR"); val != "" {
		cfg.UploadDir = val
	}
	if val := os.Getenv("FILE_URL_PREFIX"); val != "" {
		cfg.FileURLPrefix = val
	}

	// Process environment variables for database settings
	if val := os.Getenv("DB_HOST"); val != "" {
		cfg.DBConfig.Host = val
	}
	if val := os.Getenv("DB_PORT"); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			cfg.DBConfig.Port = port
		}
	}
	if val := os.Getenv("DB_USER"); val != "" {
		cfg.DBConfig.User = val
	}
	if val := os.Getenv("DB_PASSWORD"); val != "" {
		cfg.DBConfig.Password = val
	}
	if val := os.Getenv("DB_NAME"); val != "" {
		cfg.DBConfig.DBName = val
	}

	// Process environment variables for logging settings
	if val := os.Getenv("LOG_LEVEL"); val != "" {
		cfg.LogLevelValue = val
	}
	if val := os.Getenv("LOG_FORMAT"); val != "" {
		cfg.LogFormatValue = val
	}

	// Process environment variables contained in values
	overrideFromEnv(cfg)
}

// overrideFromEnv replaces ${ENV_VAR} patterns in config values with environment variable values
func overrideFromEnv(cfg *Config) {
	// Replace ${ENV_VAR} in database password
	cfg.DBConfig.Password = replaceEnvVars(cfg.DBConfig.Password)

	// Replace other values as needed
	cfg.UploadDir = replaceEnvVars(cfg.UploadDir)
	cfg.FileURLPrefix = replaceEnvVars(cfg.FileURLPrefix)
}

// replaceEnvVars replaces ${ENV_VAR} patterns in the input string with environment variable values
func replaceEnvVars(input string) string {
	result := input

	// Find all ${...} patterns
	for {
		start := strings.Index(result, "${")
		if start == -1 {
			break
		}

		end := strings.Index(result[start:], "}")
		if end == -1 {
			break
		}
		end += start

		// Extract the environment variable name
		envVarName := result[start+2 : end]
		
		// Get the environment variable value
		envVarValue := os.Getenv(envVarName)
		
		// Replace the pattern with the value
		result = result[:start] + envVarValue + result[end+1:]
	}

	return result
} 