package routes

import (
	"github.com/aliftech/todo-api/controllers"
	"github.com/aliftech/todo-api/utils"
	"github.com/gin-gonic/gin"
)

func MainRouter() {
	r := gin.Default()

	// Auth Router
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)

	// Task Router
	r.GET("/task", utils.ValidateAuth, controllers.GetTask)
	r.GET("/task/:id", utils.ValidateAuth, controllers.ShowTask)
	r.POST("/task", utils.ValidateAuth, controllers.CreateTask)
	r.PUT("/task/:id", utils.ValidateAuth, controllers.UpdateTask)
	r.DELETE("/task/:id", utils.ValidateAuth, controllers.DeleteTask)

	r.Run()
}
