package entity

import "time"

// User represents a user in the system
type User struct {
	ID        int64     `json:"user_id" gorm:"primaryKey;column:id;autoIncrement"`
	Username  string    `json:"username" gorm:"column:username;not null;size:255"`
	Email     string    `json:"email" gorm:"column:email;not null;uniqueIndex;size:255"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
} 