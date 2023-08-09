package bot

import (
	"fmt"
	teleBot "github.com/SakoDroid/telego"
	telegramConfig "github.com/SakoDroid/telego/configs"
	env "github.com/joho/godotenv"
	"os"
	"telegram-bot/db"
	"telegram-bot/dtos"
)

const (
	botTokenEnv = "TELEGRAM_BOT_TOKEN"
	timeout     = 1
	tableName   = "pepe"
)

// CreateBot returns a NewsBot with all the services it requires initialized
func CreateBot() (TelegramBotInterface, error) {

	if err := env.Load(); err != nil {
		fmt.Println("error loading environment variables")
		return nil, err
	}

	database, err := db.CreateDB[dtos.Data]("postgres", tableName, os.Getenv("DB_URL"))
	if err != nil {
		return nil, fmt.Errorf("error creating DB: %v", err)
	}
	token := os.Getenv(botTokenEnv)
	if token == "" {
		return nil, fmt.Errorf("error bot token missing")
	}

	updateConfiguration := telegramConfig.DefaultUpdateConfigs()
	botConfig := telegramConfig.BotConfigs{
		BotAPI: telegramConfig.DefaultBotAPI,
		APIKey: token, UpdateConfigs: updateConfiguration,
		Webhook:        false,
		LogFileAddress: telegramConfig.DefaultLogFile,
	}

	bot, err := teleBot.NewBot(&botConfig)
	if err != nil {
		fmt.Printf("error creating telegram bot: %v", err)
		return nil, err
	}

	return &NewsBot{
		TelegramBot: bot,
		DB:          database,
	}, nil
}
