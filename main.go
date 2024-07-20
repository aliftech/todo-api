package main

import (
	"github.com/aliftech/todo-api/migrations"
	"github.com/aliftech/todo-api/routes"
	"github.com/aliftech/todo-api/utils"
)

func init() {
	utils.Setup()
	utils.ConnectDB()
	migrations.MainMigration()
}

func main() {
	routes.MainRouter()
}
