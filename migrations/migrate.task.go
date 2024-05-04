package main

import (
	"github.com/aliftech/todo-api/models"
	"github.com/aliftech/todo-api/utils"
)

func init() {
	utils.Setup()
	utils.ConnectDB()
}

func main() {
	utils.DB.AutoMigrate(&models.Tasks{})
}
