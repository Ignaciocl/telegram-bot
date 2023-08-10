package bot

import (
	"fmt"
	teleBot "github.com/SakoDroid/telego"
	telegramConfig "github.com/SakoDroid/telego/configs"
	log "github.com/sirupsen/logrus"
	"os"
	"telegram-bot/db"
	"telegram-bot/dtos"
)

const (
	botTokenEnv = "TELEGRAM_BOT_TOKEN"
	tableName   = "pepe"
	dbUrl       = "DB_URL"
)

// CreateBot returns a NewsBot with all the services it requires initialized
func CreateBot() (TelegramBotInterface, error) {

	url := os.Getenv(dbUrl)
	if url == "" {
		return nil, fmt.Errorf("error bot url missing")
	}
	log.Infof("url is: %s", url)
	database, err := db.CreateDB[dtos.Data](tableName, url)
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
