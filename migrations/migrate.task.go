package migrations

import (
	"github.com/aliftech/todo-api/models"
	"github.com/aliftech/todo-api/utils"
)

func init() {
	utils.Setup()
	utils.ConnectDB()
}

func TaskMigration() {
	utils.DB.AutoMigrate(&models.Tasks{})
}
