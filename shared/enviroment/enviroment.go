package enviroment

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
)

var (
	once sync.Once
)

func load() {
	envPath := filepath.Join("..", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
}

func Get(key string) string {
	once.Do(load)
	return os.Getenv(key)
}
