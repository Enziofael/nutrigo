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

func DBHost() string {
	return env.GetString("DB_HOST")
}

func DBPort() string {
	return env.GetString("DB_PORT")
}

func DBUser() string {
	return env.GetString("DB_USER")
}

func DBPwd() string {
	return env.GetString("DB_PASSWORD")
}

func DBName() string {
	return env.GetString("DB_NAME")
}

func DBSslmode() string {
	return env.GetString("DB_SSL_MODE")
}

func DBConString() string {
	return fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v", DBHost(), DBPort(), DBUser(), DBPwd(), DBName(), DBSslmode())
}

func GetBackendAPIToken() string {
	return env.GetString("API_TOKEN")
}
