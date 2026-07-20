// ./backend/shared/config.go

package config

import (
	"fmt"

	env "github.com/Enziofael/nutrigo/shared/enviroment"
)

func GetDebugLogging() bool {
	return env.GetBool("DEBUG_LOGGING")
}

// ====================================================
// BOT
// ====================================================

func GetBotToken() string {
	return env.GetString("TELEGRAM_BOT_API_TOKEN")
}

// ====================================================
// DB
// ====================================================

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

// ====================================================
// Backend
// ====================================================

func GetBackendAPIToken() string {
	return env.GetString("API_TOKEN")
}

func GetBackendHost() string {
	return env.GetString("BACKEND_HOST")
}

func GetBackendPort() string {
	return env.GetString("BACKEND_PORT")
}

func GetBackendProtocolPrefix() string {
	return env.GetString("BACKEND_PROTOCOL_PREFIX")
}

func GetBackendURL() string {
	return fmt.Sprintf("%v:%v", GetBackendHost(), GetBackendPort())
}

func GetBackendFullURL() string {
	return fmt.Sprintf("%v://%v:%v", GetBackendProtocolPrefix(), GetBackendHost(), GetBackendPort())
}
