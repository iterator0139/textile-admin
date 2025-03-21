package main

import (
	"log"
	"os"
	"textile-admin/internal/config"
	"textile-admin/internal/domain/entity"
	"textile-admin/internal/handler"
	"textile-admin/internal/middleware"
	"textile-admin/internal/repository"
	"textile-admin/internal/service"
	"textile-admin/pkg/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	// Load application configuration
	cfg := config.LoadConfig()

	// Initialize all components
	router := initializeApp(cfg)

	// Start the server
	log.Printf("Starting server on %s", cfg.ServerAddress)
	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initializeApp initializes all components of the application
func initializeApp(cfg config.Config) *gin.Engine {
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
		})
	})

	return router
}

// ensureUploadDirExists ensures that the upload directory exists
func ensureUploadDirExists(uploadDir string) {
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			log.Fatalf("Failed to create upload directory: %v", err)
		}
		log.Printf("Created upload directory: %s", uploadDir)
	}
}

// connectDatabase initializes the database connection
func connectDatabase(cfg config.Config) *gorm.DB {
	// Initialize database connection
	dbConn, err := db.NewGormDBConnection(cfg.DBConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection successful")
	return dbConn
}

// migrateDatabase runs auto-migrations for database schema
func migrateDatabase(db *gorm.DB) {
	log.Println("Running database migrations...")
	err := db.AutoMigrate(&entity.User{}, &entity.ReadingTask{})
	if err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	log.Println("Database migrations completed successfully")
} 