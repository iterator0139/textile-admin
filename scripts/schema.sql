-- Create database if not exists
CREATE DATABASE IF NOT EXISTS textile_admin;

-- Use the database
USE textile_admin;

-- Create users table if needed (assumed to exist based on foreign key relationship)
CREATE TABLE IF NOT EXISTS users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create reading_tasks table
CREATE TABLE IF NOT EXISTS reading_tasks (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  file_name VARCHAR(255) NOT NULL,
  file_path VARCHAR(512) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  status ENUM('pending', 'processing', 'completed', 'failed') NOT NULL DEFAULT 'pending',
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create index for faster lookup of reading tasks by user_id
CREATE INDEX idx_reading_tasks_user_id ON reading_tasks(user_id); 