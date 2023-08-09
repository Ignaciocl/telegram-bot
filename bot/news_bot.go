package bot

import (
	"fmt"
	teleBot "github.com/SakoDroid/telego"
	"telegram-bot/db"
	"telegram-bot/dtos"
)

type TelegramBotInterface interface {
	Run() error
	StartGoRoutines() error
	StartHandlers() error
}

type NewsBot struct {
	TelegramBot *teleBot.Bot
	DB          db.DB[dtos.Data]
	channels    map[string]chan dtos.GetInformation
}

func (nb *NewsBot) Run() error {
	err := nb.TelegramBot.Run()
	if err != nil {
		fmt.Printf("error initializating telegram bot: %v", err)
		return err
	}

	return nil
}
