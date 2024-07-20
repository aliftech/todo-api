package migrations

func MainMigration() {
	TaskMigration()
	UserMigrate()
}
