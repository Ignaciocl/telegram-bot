package bot

import (
	"fmt"
	objs "github.com/SakoDroid/telego/objects"
	log "github.com/sirupsen/logrus"
	"telegram-bot/dtos"
	"time"
)

func (nb *NewsBot) StartGoRoutines() error {
	m := StartHandlersOperations(nb)
	nb.channels = m
	return nil
}

func (nb *NewsBot) StartHandlers() error {
	bot := nb.TelegramBot

	bot.AddHandler("/start", func(u *objs.Update) {
		kb := bot.CreateInlineKeyboard()
		di := dtos.Data{
			FreeTimes: []dtos.Times{},
			Id:        u.Message.From.Id,
			Timezone:  "Argentina",
			Name:      u.Message.From.FirstName,
		}
		kb.AddCallbackButtonHandler("Agregar horario libre", "smt", 1, func(update *objs.Update) {
			nb.channels[addTimeAvailable] <- dtos.GetInformation{
				Id:       u.Message.From.Id,
				ToAnswer: u.Message.Chat.Id,
			}
		})
		kb.AddCallbackButtonHandler("Cambiar zona horaria", "smt1", 2, func(update *objs.Update) {
			nb.channels[changeTimeZone] <- dtos.GetInformation{
				Id:       u.Message.From.Id,
				ToAnswer: u.Message.Chat.Id,
			}
		})
		kb.AddCallbackButtonHandler("Buscar calendario de todos", "smt2", 2, func(update *objs.Update) {
			nb.channels[getSchedules] <- dtos.GetInformation{
				Id:       u.Message.From.Id,
				ToAnswer: u.Message.Chat.Id,
			}
		})
		zone, _ := time.LoadLocation("America/Argentina/Buenos_Aires")
		log.Infof("%v, \nentire message is: %+v", u.Message.Date, u.Message)
		log.Infof("date is: %v", time.Unix(int64(u.Message.Date), 0).In(zone))

		nb.DB.Insert(di)

		//Sends the message to the chat that the message has been received from. The message will be a reply to the received message.
		_, err := bot.AdvancedMode().ASendMessage(u.Message.Chat.Id, "Selecciona una de las opciones de abajo (se asume que tu region es argentina).", "", u.Message.MessageId, false, false, nil, false, false, kb)
		if err != nil {
			fmt.Printf("error happened, %v\n", err)
		}
	}, "private")
	bot.AddHandler("/add", func(u *objs.Update) {
		nb.channels[addTimeAvailable] <- dtos.GetInformation{
			Id:       u.Message.From.Id,
			ToAnswer: u.Message.Chat.Id,
		}
	}, "private")
	bot.AddHandler("/schedules", func(u *objs.Update) {
		nb.channels[getSchedules] <- dtos.GetInformation{
			Id:       u.Message.From.Id,
			ToAnswer: u.Message.Chat.Id,
		}
	}, "all")
	//bot.AddHandler(".*", func(u *objs.Update) {
	//	fmt.Println("pepe")
	//	bot.SendMessage(u.Message.Chat.Id, fmt.Sprintf("you sent: %s", u.Message.Text), "", u.Message.MessageId, false, false)
	//}, "private")
	return nil
}
