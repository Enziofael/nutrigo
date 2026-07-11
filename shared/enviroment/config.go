package enviroment

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func Load() {
	envPath := filepath.Join("..", "..", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func Get(key string) string {
	return os.Getenv(key)
}
