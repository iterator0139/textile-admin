package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"textile-admin/internal/service"
	"textile-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// ReadingHandler handles HTTP requests for reading tasks
type ReadingHandler struct {
	service     *service.ReadingService
	uploadDir   string
}

// NewReadingHandler creates a new instance of ReadingHandler
func NewReadingHandler(service *service.ReadingService, uploadDir string) *ReadingHandler {
	return &ReadingHandler{
		service:     service,
		uploadDir:   uploadDir,
	}
}

// RegisterRoutes registers the routes for reading tasks
func (h *ReadingHandler) RegisterRoutes(router *gin.Engine) {
	readingGroup := router.Group("/api/reading")
	{
		readingGroup.POST("/upload", h.UploadFile)
		readingGroup.GET("/task/:task_id", h.GetTask)
		readingGroup.GET("/tasks/user/:user_id", h.GetUserTasks)
		readingGroup.PUT("/task/:task_id/status", h.UpdateTaskStatus)
	}

	// Route for file download
	router.GET("/files/:file_name", h.DownloadFile)
}

// UploadFile handles the file upload and creation of reading task
func (h *ReadingHandler) UploadFile(c *gin.Context) {
	// Parse user ID from form data
	userIDStr := c.PostForm("user_id")
	if userIDStr == "" {
		response.BadRequest(c, "User ID is required")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID format")
		return
	}

	// Get the file from form data
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "File is required")
		return
	}

	// Validate file size (example: max 50MB)
	if file.Size > 50*1024*1024 {
		response.BadRequest(c, "File size exceeds the limit (50MB)")
		return
	}

	// Validate file type if needed
	// ...

	// Create the reading task
	result, err := h.service.CreateTask(userID, file)
	if err != nil {
		response.InternalServerError(c, "Failed to create reading task: "+err.Error())
		return
	}

	response.Success(c, "阅读任务创建成功", result)
}

// GetTask handles the retrieval of a reading task by ID
func (h *ReadingHandler) GetTask(c *gin.Context) {
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid task ID format")
		return
	}

	task, err := h.service.GetTaskByID(taskID)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve task: "+err.Error())
		return
	}

	if task == nil {
		response.NotFound(c, "Task not found")
		return
	}

	response.Success(c, "查询成功", task)
}

// GetUserTasks handles the retrieval of all reading tasks for a user
func (h *ReadingHandler) GetUserTasks(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID format")
		return
	}

	tasks, err := h.service.GetTasksByUserID(userID)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve tasks: "+err.Error())
		return
	}

	response.Success(c, "查询成功", tasks)
}

// UpdateTaskStatus handles updating the status of a reading task
func (h *ReadingHandler) UpdateTaskStatus(c *gin.Context) {
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid task ID format")
		return
	}

	// Get the status from JSON body
	var requestBody struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		"pending":    true,
		"processing": true,
		"completed":  true,
		"failed":     true,
	}

	if !validStatuses[requestBody.Status] {
		response.BadRequest(c, "Invalid status value. Must be one of: pending, processing, completed, failed")
		return
	}

	// Update the task status
	err = h.service.UpdateTaskStatus(taskID, requestBody.Status)
	if err != nil {
		response.InternalServerError(c, "Failed to update task status: "+err.Error())
		return
	}

	response.Success(c, "状态更新成功", nil)
}

// DownloadFile handles file download requests
func (h *ReadingHandler) DownloadFile(c *gin.Context) {
	fileName := c.Param("file_name")
	
	// For security reasons, let's sanitize the filename to prevent directory traversal
	fileName = filepath.Base(fileName)
	
	// Construct the file path
	filePath := filepath.Join(h.uploadDir, fileName)
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		response.NotFound(c, "File not found")
		return
	}
	
	// Set the appropriate content disposition and serve the file
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.File(filePath)
} 