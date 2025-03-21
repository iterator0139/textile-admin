package entity

import "time"

// ReadingTask represents a user's reading task and its associated file
type ReadingTask struct {
	ID        int64     `json:"task_id" gorm:"primaryKey;column:id;autoIncrement"`
	UserID    int64     `json:"user_id" gorm:"column:user_id;not null;index"`
	FileName  string    `json:"file_name" gorm:"column:file_name;not null;size:255"`
	FilePath  string    `json:"file_path" gorm:"column:file_path;not null;size:512"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	Status    string    `json:"status" gorm:"column:status;not null;default:pending;type:enum('pending','processing','completed','failed')"`
}

// TableName specifies the table name for ReadingTask
func (ReadingTask) TableName() string {
	return "reading_tasks"
}

// TaskResponse represents the response for a reading task
type TaskResponse struct {
	TaskID    int64     `json:"task_id"`
	UserID    int64     `json:"user_id"`
	FileName  string    `json:"file_name"`
	FileURL   string    `json:"file_url"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// UploadResponse represents the response for a file upload
type UploadResponse struct {
	TaskID   int64  `json:"task_id"`
	FileName string `json:"file_name"`
	FileURL  string `json:"file_url"`
} 