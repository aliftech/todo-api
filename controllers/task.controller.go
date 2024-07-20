package controllers

import (
	"net/http"

	"github.com/aliftech/todo-api/models"
	"github.com/aliftech/todo-api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateTask(c *gin.Context) {
	// Validate authorization
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized",
			"data":    nil,
		})
		return
	}

	// Bind and validate the request body
	var taskRequest struct {
		Title       string
		Description string
		Due         string
		Status      string
	}

	if err := c.ShouldBindJSON(&taskRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request data",
			"data":    nil,
		})
		return
	}

	// Create a new task
	task := models.Tasks{
		Userid:      userID.(uint), // Assuming userID is of type uint
		Title:       taskRequest.Title,
		Description: taskRequest.Description,
		Due:         taskRequest.Due,
		Status:      taskRequest.Status,
	}
	if result := utils.DB.Create(&task); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to create task",
			"data":    nil,
		})
		return
	}

	// Return the created task
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "New task have been added.",
		"data":    task,
	})
}

func GetTask(c *gin.Context) {
	// Validate authorization
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized",
			"data":    nil,
		})
		return
	}

	// Ensure userID is of type int
	userIDInt, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Invalid user ID",
			"data":    nil,
		})
		return
	}

	// Get the task data for the authenticated user
	var tasks []models.Tasks
	if result := utils.DB.Where("userid = ?", userIDInt).Find(&tasks); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to retrieve tasks",
			"data":    nil,
		})
		return
	}

	// Return the result in response
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Get all user tasks.",
		"data":    tasks,
	})
}

func ShowTask(c *gin.Context) {
	// Validate authorization
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized",
			"data":    nil,
		})
		return
	}

	// Ensure userID is of type int
	userIDInt, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Invalid user ID",
			"data":    nil,
		})
		return
	}
	// Get ID from url
	task_id := c.Param("id")

	// Get the task data
	var task models.Tasks
	if result := utils.DB.Where("userid = ? AND id = ?", userIDInt, task_id).First(&task); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   true,
				"message": "Task not found",
				"data":    nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Failed to retrieve task",
				"data":    nil,
			})
		}
		return
	}

	// Return the result in response
	c.JSON(200, gin.H{
		"error":   false,
		"message": "Show detail task.",
		"data":    task,
	})
}

func UpdateTask(c *gin.Context) {
	// Validate authorization
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized",
			"data":    nil,
		})
		return
	}

	// Ensure userID is of type int
	userIDInt, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Invalid user ID",
			"data":    nil,
		})
		return
	}

	// Get ID from url
	task_id := c.Param("id")

	// Get data from the request body
	var task_request struct {
		Title       string
		Description string
		Due         string
		Status      string
	}

	c.Bind(&task_request)

	// Show the task that gonna be updated
	var task models.Tasks
	if result := utils.DB.Where("userid = ? AND id = ?", userIDInt, task_id).Find(&task); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   true,
				"message": "Task not found",
				"data":    nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Failed to retrieve task",
				"data":    nil,
			})
		}
		return
	}

	// Update task
	utils.DB.Model(&task).Updates(models.Tasks{
		Title:       task_request.Title,
		Description: task_request.Description,
		Due:         task_request.Due,
		Status:      task_request.Status,
	})

	// Return a response
	c.JSON(200, gin.H{
		"error":   false,
		"message": "Task have been updated.",
		"data":    task,
	})
}

func DeleteTask(c *gin.Context) {
	// Validate authorization
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized",
			"data":    nil,
		})
		return
	}

	// Ensure userID is of type int
	userIDInt, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Invalid user ID",
			"data":    nil,
		})
		return
	}

	// Get task ID from URL
	taskID := c.Param("id")

	// Retrieve the task to check if it exists and belongs to the user
	var task models.Tasks
	if result := utils.DB.Where("userid = ? AND id = ?", userIDInt, taskID).First(&task); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   true,
				"message": "Task not found",
				"data":    nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Failed to retrieve task",
				"data":    nil,
			})
		}
		return
	}

	// Delete the task
	if result := utils.DB.Delete(&task); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to delete task",
			"data":    nil,
		})
		return
	}

	// Response the message
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Task has been deleted.",
		"data":    nil,
	})
}
