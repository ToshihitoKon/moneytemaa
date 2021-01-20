package constants

import (
	"log"
	"os"
)

var (
	DbUser            = NewEnv("MONEYTEMAA_DB_USER", true)
	DbPass            = NewEnv("MONEYTEMAA_DB_PASS", true)
	DbName            = NewEnv("MONEYTEMAA_DB_NAME", true)
	SlackSigninSecret = NewEnv("MONEYTEMAA_SLACK_SIGNIN_SECRET", true)
	SlackBotUserId    = NewEnv("MONEYTEMAA_SLACK_BOT_USER_ID", true)
	SlackBotToken     = NewEnv("MONEYTEMAA_SLACK_BOT_TOKEN", true)
)

func NewEnv(key string, required bool) string {
	value := os.Getenv(key)
	if required && value == "" {
		log.Fatalf("%s is required", key)
	}
	return value
}
