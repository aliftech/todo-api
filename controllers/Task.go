package controllers

import (
	"github.com/aliftech/todo-api/models"
	"github.com/aliftech/todo-api/utils"
	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	// Get data of request body
	var task_request struct {
		Title       string
		Description string
		Status      string
	}

	c.Bind(&task_request)

	// Create a post
	task := models.Tasks{Title: task_request.Title, Description: task_request.Description, Status: task_request.Status}
	result := utils.DB.Create(&task)

	if result.Error != nil {
		c.Status(400)
		return
	}
	// Return it
	c.JSON(200, gin.H{
		"data": task,
	})
}

func GetTask(c *gin.Context) {
	// Get the task data
	var tasks []models.Tasks
	utils.DB.Find(&tasks)

	// Return the result in response
	c.JSON(200, gin.H{
		"data": tasks,
	})
}

func ShowTask(c *gin.Context) {
	// Get ID from url
	task_id := c.Param("id")

	// Get the task data
	var task models.Tasks
	utils.DB.First(&task, task_id)

	// Return the result in response
	c.JSON(200, gin.H{
		"data": task,
	})
}

func UpdateTask(c *gin.Context) {
	// Get ID from url
	task_id := c.Param("id")

	// Get data from the request body
	var task_request struct {
		Title       string
		Description string
		Status      string
	}

	c.Bind(&task_request)

	// Show the task that gonna be updated
	var task models.Tasks
	utils.DB.First(&task, task_id)

	// Update task
	utils.DB.Model(&task).Updates(models.Tasks{
		Title:       task_request.Title,
		Description: task.Description,
		Status:      task_request.Status,
	})

	// Return a response
	c.JSON(200, gin.H{
		"data": task,
	})
}

func DeleteTask(c *gin.Context) {
	// Get id from url
	task_id := c.Param("id")

	// Delete the task
	utils.DB.Delete(&models.Tasks{}, task_id)

	// Response the message
	c.JSON(200, gin.H{
		"message": "A task have been deleted.",
	})
}
