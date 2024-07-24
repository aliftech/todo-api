package routes

import (
	"github.com/aliftech/todo-api/controllers"
	"github.com/aliftech/todo-api/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func MainRouter() {
	r := gin.Default()

	// Configure CORS options
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8000"}, // Replace with your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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
