package bot

import (
	"bot-telegram/db"
	"bot-telegram/dtos"
	"bot-telegram/services/news"
	"fmt"
	teleBot "github.com/SakoDroid/telego"
)

type NewsBotInterface interface {
	Run() error
	StartGoRoutines() error
	StartHandlers() error
	SendScheduledNews(hour int) error
}

type NewsBot struct {
	TelegramBot *teleBot.Bot
	DB          db.DB[dtos.Data]
	ScheduleDB  db.DB[dtos.Schedule]
	GPTService  news.Provider
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
