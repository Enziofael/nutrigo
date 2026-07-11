package enviroment

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
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

func GetString(key string) string {
	once.Do(load)
	return os.Getenv(key)
}

func GetBool(key string) bool {
	once.Do(load)

	b, err := strconv.ParseBool(GetString(key))
	if err != nil {
		log.Panicf("Invalid env value for %s: \"%s\". Should be bool", key, GetString(key))
	}

	return b
}

func GetInt(key string) int {
	once.Do(load)

	i, err := strconv.Atoi(GetString(key))
	if err != nil {
		log.Panicf("Invalid env value for %s: \"%s\". Should be int", key, GetString(key))
	}

	return i
}

func GetDouble(key string) float64 {
	once.Do(load)

	f, err := strconv.ParseFloat(GetString(key), 64)
	if err != nil {
		log.Panicf("Invalid env value for %s: \"%s\". Should be float64", key, GetString(key))
	}

	return f
}

func GetFloat(key string) float32 {
	once.Do(load)

	f, err := strconv.ParseFloat(GetString(key), 32)
	if err != nil {
		log.Panicf("Invalid env value for %s: \"%s\". Should be float32", key, GetString(key))
	}

	return float32(f)
}
