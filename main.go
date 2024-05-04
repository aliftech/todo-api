package main

import (
	"github.com/aliftech/todo-api/controllers"
	"github.com/aliftech/todo-api/utils"
	"github.com/gin-gonic/gin"
)

func init() {
	utils.Setup()
	utils.ConnectDB()
}

func main() {
	r := gin.Default()

	r.GET("/task", controllers.GetTask)
	r.GET("/task/:id", controllers.ShowTask)
	r.POST("/task", controllers.CreateTask)
	r.PUT("/task/:id", controllers.UpdateTask)
	r.DELETE("/task/:id", controllers.DeleteTask)

	r.Run()
}
