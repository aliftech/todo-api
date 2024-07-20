package routes

import (
	"github.com/aliftech/todo-api/controllers"
	"github.com/gin-gonic/gin"
)

func TaskRouter() {
	r := gin.Default()

	r.GET("/task", controllers.GetTask)
	r.GET("/task/:id", controllers.ShowTask)
	r.POST("/task", controllers.CreateTask)
	r.PUT("/task/:id", controllers.UpdateTask)
	r.DELETE("/task/:id", controllers.DeleteTask)

	r.Run()
}
