package main

import (
	"os"
	"textile-admin/internal/config"
	"textile-admin/internal/domain/entity"
	"textile-admin/internal/handler"
	"textile-admin/internal/middleware"
	"textile-admin/internal/repository"
	"textile-admin/internal/service"
	"textile-admin/pkg/db"
	"textile-admin/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	// Load application configuration
	cfg := config.LoadConfig()

	// Initialize logger based on configuration
	if cfg.LogFormat() == "json" {
		logger.InitJSONLogger(cfg.LogLevel())
	} else {
		logger.InitTextLogger(cfg.LogLevel())
	}

	logger.Info("Starting application with environment: " + getEnv())
	logger.Info("Server will listen on " + cfg.ServerAddress)

	// Initialize all components
	router := initializeApp(cfg)

	// Start the server
	logger.Info("Starting server on " + cfg.ServerAddress)
	if err := router.Run(cfg.ServerAddress); err != nil {
		logger.Fatal("Failed to start server: " + err.Error())
	}
}

// getEnv returns the current environment
func getEnv() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev" // Default to development
	}
	return env
}

// initializeApp initializes all components of the application
func initializeApp(cfg config.Config) *gin.Engine {
	// Set Gin mode based on environment
	if getEnv() == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Ensure upload directory exists
	ensureUploadDirExists(cfg.UploadDir)

	// Initialize database connection
	dbConn := connectDatabase(cfg)

	// Auto migrate database schema if needed
	migrateDatabase(dbConn)

	// Initialize components
	readingRepo := repository.NewReadingRepository(dbConn)
	readingService := service.NewReadingService(readingRepo, cfg.UploadDir, cfg.FileURLPrefix)
	readingHandler := handler.NewReadingHandler(readingService, cfg.UploadDir)

	// Initialize Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(middleware.CORSMiddleware())

	// Register routes
	readingHandler.RegisterRoutes(router)

	// Add a health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
			"env":    getEnv(),
		})
	})

	return router
}

// ensureUploadDirExists ensures that the upload directory exists
func ensureUploadDirExists(uploadDir string) {
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			logger.Fatal("Failed to create upload directory: " + err.Error())
		}
		logger.Info("Created upload directory: " + uploadDir)
	}
}

// connectDatabase initializes the database connection
func connectDatabase(cfg config.Config) *gorm.DB {
	// Initialize database connection
	dbConn, err := db.NewGormDBConnection(cfg.DBConfig)
	if err != nil {
		logger.Fatal("Failed to connect to database: " + err.Error())
	}

	logger.Info("Database connection successful")
	return dbConn
}

// migrateDatabase runs auto-migrations for database schema
func migrateDatabase(db *gorm.DB) {
	logger.Info("Running database migrations...")
	err := db.AutoMigrate(&entity.User{}, &entity.ReadingTask{})
	if err != nil {
		logger.Fatal("Failed to run database migrations: " + err.Error())
	}
	logger.Info("Database migrations completed successfully")
} 