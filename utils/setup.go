package utils

import (
	"log"
	"os"
	"runtime"

	"github.com/joho/godotenv"
)

func Setup() {
	// Check if running in App Engine environment
	if os.Getenv("GAE_ENV") == "" && runtime.GOOS != "appengine" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}
