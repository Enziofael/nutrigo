// ./backend/shared/config.go

package config

import (
	"fmt"

	env "github.com/Enziofael/nutrigo/shared/enviroment"
)

func GetBotToken() string {
	return env.GetString("TELEGRAM_BOT_API_TOKEN")
}

func GetDebugLogging() bool {
	return env.GetBool("DEBUG_LOGGING")
}

func GetBackendPort() string {
	return fmt.Sprintf(":%v", env.GetString("BACKEND_PORT"))
}

func DbHost() string {
	return env.GetString("DB_HOST")
}

func DbPort() string {
	return env.GetString("DB_PORT")
}

func DbUser() string {
	return env.GetString("DB_USER")
}

func DbName() string {
	return env.GetString("DB_NAME")
}

func DbSslmode() string {
	return env.GetString("DB_SSL_MODE")
}

func DbConString() string {
	return fmt.Sprintf("host=%v port=%v user=%v dbname=%v sslmode=%v", DbHost(), DbPort(), DbUser(), DbName(), DbSslmode())
}
