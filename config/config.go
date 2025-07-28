package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var once sync.Once

// Config fetches environment variables, loading .env only once
func Config(key string) string {
	once.Do(func() {
		err := godotenv.Load(".env")
		if err != nil {
			fmt.Println("Warning: .env file not loaded")
		}
	})
	return os.Getenv(key)
}
