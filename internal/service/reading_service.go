package service

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"textile-admin/internal/domain/entity"
	"textile-admin/internal/repository"

	"github.com/google/uuid"
)

// ReadingService handles the business logic for reading tasks
type ReadingService struct {
	repo          *repository.ReadingRepository
	uploadDir     string
	fileURLPrefix string
}

// NewReadingService creates a new instance of ReadingService
func NewReadingService(repo *repository.ReadingRepository, uploadDir, fileURLPrefix string) *ReadingService {
	return &ReadingService{
		repo:          repo,
		uploadDir:     uploadDir,
		fileURLPrefix: fileURLPrefix,
	}
}

// CreateTask creates a new reading task and saves the uploaded file
func (s *ReadingService) CreateTask(userID int64, file *multipart.FileHeader) (*entity.UploadResponse, error) {
	// Generate a unique filename to prevent collisions
	originalFilename := filepath.Base(file.Filename)
	uniqueFilename := generateUniqueFilename(originalFilename)
	
	// Define the file path
	filePath := filepath.Join(s.uploadDir, uniqueFilename)
	
	// Save the file
	if err := saveUploadedFile(file, filePath); err != nil {
		log.Printf("Error saving file: %v", err)
		return nil, err
	}
	
	// Create task in database
	taskID, err := s.repo.CreateTask(userID, originalFilename, filePath)
	if err != nil {
		// Attempt to delete the file if database operation fails
		os.Remove(filePath)
		return nil, err
	}
	
	// Build the file URL
	fileURL := fmt.Sprintf("%s/%s", s.fileURLPrefix, uniqueFilename)
	
	return &entity.UploadResponse{
		TaskID:   taskID,
		FileName: originalFilename,
		FileURL:  fileURL,
	}, nil
}

// GetTaskByID retrieves a task by its ID and converts it to response format
func (s *ReadingService) GetTaskByID(taskID int64) (*entity.TaskResponse, error) {
	task, err := s.repo.GetTaskByID(taskID)
	if err != nil {
		return nil, err
	}
	
	if task == nil {
		return nil, nil
	}
	
	// Extract the filename from the file path
	filename := filepath.Base(task.FilePath)
	fileURL := fmt.Sprintf("%s/%s", s.fileURLPrefix, filename)
	
	return &entity.TaskResponse{
		TaskID:    task.ID,
		UserID:    task.UserID,
		FileName:  task.FileName,
		FileURL:   fileURL,
		Status:    task.Status,
		CreatedAt: task.CreatedAt,
	}, nil
}

// GetTasksByUserID retrieves all tasks for a user and converts them to response format
func (s *ReadingService) GetTasksByUserID(userID int64) ([]*entity.TaskResponse, error) {
	tasks, err := s.repo.GetTasksByUserID(userID)
	if err != nil {
		return nil, err
	}
	
	var responses []*entity.TaskResponse
	for _, task := range tasks {
		filename := filepath.Base(task.FilePath)
		fileURL := fmt.Sprintf("%s/%s", s.fileURLPrefix, filename)
		
		responses = append(responses, &entity.TaskResponse{
			TaskID:    task.ID,
			UserID:    task.UserID,
			FileName:  task.FileName,
			FileURL:   fileURL,
			Status:    task.Status,
			CreatedAt: task.CreatedAt,
		})
	}
	
	return responses, nil
}

// UpdateTaskStatus updates the status of a reading task
func (s *ReadingService) UpdateTaskStatus(taskID int64, status string) error {
	return s.repo.UpdateTaskStatus(taskID, status)
}

// GetFilePath returns the actual file path for a given task
func (s *ReadingService) GetFilePath(taskID int64) (string, error) {
	task, err := s.repo.GetTaskByID(taskID)
	if err != nil {
		return "", err
	}
	
	if task == nil {
		return "", fmt.Errorf("task not found")
	}
	
	return task.FilePath, nil
}

// saveUploadedFile saves the uploaded file to the specified destination
func saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	
	_, err = io.Copy(out, src)
	return err
}

// generateUniqueFilename creates a unique filename by adding a UUID
func generateUniqueFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	name := originalName[:len(originalName)-len(ext)]
	uuid := uuid.New().String()
	return fmt.Sprintf("%s_%s%s", name, uuid, ext)
} 