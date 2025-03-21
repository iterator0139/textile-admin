package repository

import (
	"log"
	"textile-admin/internal/domain/entity"

	"gorm.io/gorm"
)

// ReadingRepository handles database operations for reading tasks
type ReadingRepository struct {
	db *gorm.DB
}

// NewReadingRepository creates a new instance of ReadingRepository
func NewReadingRepository(db *gorm.DB) *ReadingRepository {
	return &ReadingRepository{db: db}
}

// CreateTask creates a new reading task in the database
func (r *ReadingRepository) CreateTask(userID int64, fileName, filePath string) (int64, error) {
	task := entity.ReadingTask{
		UserID:   userID,
		FileName: fileName,
		FilePath: filePath,
		Status:   "pending",
	}

	result := r.db.Create(&task)
	if result.Error != nil {
		log.Printf("Error creating reading task: %v", result.Error)
		return 0, result.Error
	}

	return task.ID, nil
}

// GetTaskByID retrieves a reading task by its ID
func (r *ReadingRepository) GetTaskByID(taskID int64) (*entity.ReadingTask, error) {
	var task entity.ReadingTask
	
	result := r.db.First(&task, taskID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // No task found
		}
		log.Printf("Error querying task by ID: %v", result.Error)
		return nil, result.Error
	}

	return &task, nil
}

// GetTasksByUserID retrieves all reading tasks for a given user
func (r *ReadingRepository) GetTasksByUserID(userID int64) ([]*entity.ReadingTask, error) {
	var tasks []*entity.ReadingTask
	
	result := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&tasks)
	if result.Error != nil {
		log.Printf("Error querying tasks by user ID: %v", result.Error)
		return nil, result.Error
	}

	return tasks, nil
}

// UpdateTaskStatus updates the status of a reading task
func (r *ReadingRepository) UpdateTaskStatus(taskID int64, status string) error {
	result := r.db.Model(&entity.ReadingTask{}).Where("id = ?", taskID).Update("status", status)
	if result.Error != nil {
		log.Printf("Error updating task status: %v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
} 