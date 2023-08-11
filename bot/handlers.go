package bot

import (
	"fmt"
	objs "github.com/SakoDroid/telego/objects"
	log "github.com/sirupsen/logrus"
	"strings"
	"telegram-bot/dtos"
	"telegram-bot/utils"
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
	bot.AddHandler("/allSubscribed", func(u *objs.Update) {
		d, _ := nb.DB.GetAll()
		names := make([]string, 0)
		for _, i := range d {
			names = append(names, i.Name)
		}
		bot.SendMessage(u.Message.Chat.Id, fmt.Sprintf("people is: %s", strings.Join(names, ", ")), "", 0, false, false)
	}, "all")
	bot.AddHandler("/changeName", func(u *objs.Update) {
		data := strings.Split(u.Message.Text, " ")
		if len(data) < 2 {
			bot.SendMessage(u.Message.Chat.Id, fmt.Sprintf("message received is: %s\nCouldn't change name", u.Message.Text), "", 0, false, false)
			return
		}
		p, _ := nb.DB.Get(u.Message.From.Id)
		p.Name = data[1]
		nb.DB.Update(p)
		bot.SendMessage(u.Message.Chat.Id, fmt.Sprintf("NAME CHANGED TO %s", p.Name), "", 0, false, false)
	}, "all")
	bot.AddHandler("/removeDay", func(u *objs.Update) {
		data := strings.Split(u.Message.Text, " ")
		if len(data) < 2 {
			bot.SendMessage(u.Message.Chat.Id, fmt.Sprintf("message received is: %s\nCouldn't change name", u.Message.Text), "", 0, false, false)
			return
		}
		day := strings.ToLower(data[1])
		if !utils.Contains(day, []string{"lunes", "martes", "miercoles", "jueves", "viernes", "sabado", "domingo"}) {
			bot.SendMessage(u.Message.Chat.Id, fmt.Sprintf("tiene que ser un dia de la semana, recibido: %s", u.Message.Text), "", 0, false, false)
			return
		}
		p, _ := nb.DB.Get(u.Message.From.Id)
		times := make([]dtos.Times, 0)
		found := false
		for _, k := range p.FreeTimes {
			if k.Day == day {
				found = true
			} else {
				times = append(times, k)
			}
		}
		if found {
			nb.DB.Update(p)
			bot.SendMessage(u.Message.Chat.Id, fmt.Sprintf("Se saco el dia: %s\n", day), "", 0, false, false)
		} else {
			bot.SendMessage(u.Message.Chat.Id, fmt.Sprintf("No tenias el dia: %s\n", day), "", 0, false, false)
		}
	})
	//bot.AddHandler(".*", func(u *objs.Update) {
	//	fmt.Println("pepe")
	//	bot.SendMessage(u.Message.Chat.Id, fmt.Sprintf("you sent: %s", u.Message.Text), "", u.Message.MessageId, false, false)
	//}, "private")
	return nil
}
