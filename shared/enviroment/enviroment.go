// ./backend/shared/enviroment.go

// Helper package to work with godotenv. Allow to get parsed values
package enviroment

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

const DOTENV_SEARCHING_DEPTH int = 5

// Loads .env file in godotenv
func init() {

	envPath := ".env"

	for i := 0; i < DOTENV_SEARCHING_DEPTH; i++ {

		log.Printf("Loading .env attempt №%v", i)
		err := godotenv.Load(envPath)
		if err == nil {
			log.Printf("- Loaded succesfully after attempt #%v with path %v", i, envPath)
			return
		} else {
			if i == DOTENV_SEARCHING_DEPTH-1 {
				log.Fatalf("Fatal error loading .env file after %v attempts: %s", i, err)
			}
		}
		envPath = filepath.Join("..", ".env")
	}
}

// Returns string env value
func GetString(key string) string {
	return os.Getenv(key)
}

// Returns bool env value. Panics if parse failed
func GetBool(key string) bool {
	b, err := strconv.ParseBool(GetString(key))
	if err != nil {
		log.Panicf("Invalid env value for %s: \"%s\". Should be bool", key, GetString(key))
	}

	return b
}

// Returns int env value. Panics if parse failed
func GetInt(key string) int {
	i, err := strconv.Atoi(GetString(key))
	if err != nil {
		log.Panicf("Invalid env value for %s: \"%s\". Should be int", key, GetString(key))
	}
	return i
}

// Returns float64 env value. Panics if parse failed
func GetDouble(key string) float64 {
	f, err := strconv.ParseFloat(GetString(key), 64)
	if err != nil {
		log.Panicf("Invalid env value for %s: \"%s\". Should be float64", key, GetString(key))
	}
	return f
}

// Returns float32 env value. Panics if parse failed
func GetFloat(key string) float32 {
	f, err := strconv.ParseFloat(GetString(key), 32)
	if err != nil {
		log.Panicf("Invalid env value for %s: \"%s\". Should be float32", key, GetString(key))
	}
	return float32(f)
}
