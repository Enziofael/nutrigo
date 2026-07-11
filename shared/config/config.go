package config

import (
	env "github.com/Enziofael/nutrigo/shared/enviroment"
)

func GetBotToken() string {
	return env.GetString("TELEGRAM_BOT_API_TOKEN")
}

func GetDebugLogging() bool {
	return env.GetBool("DEBUG_LOGGING")
}