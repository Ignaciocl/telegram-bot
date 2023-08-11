package bot

import (
	"fmt"
	teleBot "github.com/SakoDroid/telego"
	objs "github.com/SakoDroid/telego/objects"
	log "github.com/sirupsen/logrus"
	"math"
	"strconv"
	"telegram-bot/db"
	"telegram-bot/dtos"
)

const (
	addTimeAvailable = "AddTimeAvailable"
	changeTimeZone   = "changeTimeZone"
	getSchedules     = "getSchedules"
)

var TimezoneMap = map[string]string{
	"Argentina":             "America/Argentina/Buenos_Aires",
	"Australia (Melbourne)": "Australia/Melbourne",
}

var PossibleDays = []string{"lunes", "martes", "miercoles", "jueves", "viernes", "sabado", "domingo"}

func StartHandlersOperations(tBot *NewsBot) map[string]chan dtos.GetInformation {
	at := getNewChan()
	cz := getNewChan()
	gs := getNewChan()
	go changeTimezone(cz, tBot.DB, tBot.TelegramBot)
	go addTimeAvailability(at, tBot.DB, tBot.TelegramBot)
	go getSchedulesSet(gs, tBot.DB, tBot.TelegramBot)
	// ToDo there HAS to be a better way to do this
	return map[string]chan dtos.GetInformation{
		addTimeAvailable: at,
		changeTimeZone:   cz,
		getSchedules:     gs,
	}
}

func getNewChan() chan dtos.GetInformation {
	return make(chan dtos.GetInformation, 100)
}

func changeTimezone(c chan dtos.GetInformation, database db.DB[dtos.Data], bot *teleBot.Bot) {
	for {
		toAdd := <-c
		go addCategoryForChat(toAdd, database, bot)
	}
}

func addCategoryForChat(info dtos.GetInformation, db db.DB[dtos.Data], bot *teleBot.Bot) {
	data, err := db.Get(info.Id)
	if err != nil {
		log.Errorf("couldn't get info for user, err: %v", err)
	}
	kb := bot.CreateInlineKeyboard()
	counter := 1
	for k := range TimezoneMap {
		kb.AddCallbackButtonHandler(k, k, counter%len(TimezoneMap)+1, func(u *objs.Update) {
			id := u.Message.From.Id
			toAnswer := u.Message.Chat.Id
			d, err := db.Get(id)
			key := u.CallbackQuery.Data
			if err != nil {
				log.Errorf("couldn't get info for user, err: %v", err)
				return
			}
			if key == d.Timezone {
				bot.SendMessage(toAnswer, fmt.Sprintf("Ya tenias el timezone de: %s, apreta /start para seguir", key), "", 0, false, false)
			} else {
				d.Timezone = k
				db.Update(d)

				message := fmt.Sprintf("ya tenes %s como region horaria. Apreta /start para seguir", k)
				bot.SendMessage(toAnswer, message, "", 0, false, false)
			}
		})
	}
	bot.AdvancedMode().ASendMessage(info.ToAnswer, fmt.Sprintf("Selecciona el que quieras cambiar de abajo, la que tenes seleccionada ya es: %v", data.Timezone), "", 0, false, false, nil, false, false, kb)
}

func addTimeAvailability(c chan dtos.GetInformation, database db.DB[dtos.Data], bot *teleBot.Bot) {
	for {
		toAdd := <-c
		go chooseDay(toAdd, database, bot)
	}
}

func chooseDay(info dtos.GetInformation, db db.DB[dtos.Data], bot *teleBot.Bot) {
	kb := bot.CreateInlineKeyboard()
	for i, k := range PossibleDays {
		kb.AddCallbackButtonHandler(k, k, i%4+1, func(u *objs.Update) {
			d, err := db.Get(u.Message.From.Id)
			key := u.CallbackQuery.Data
			if err != nil {
				log.Errorf("couldn't get info for user, err: %v", err)
				return
			}
			addTimeToDay(key, d.FreeTimes, u.Message.Chat.Id, bot, db, d)
		})
	}
	bot.AdvancedMode().ASendMessage(info.ToAnswer, fmt.Sprintf("Elegi el dia que quieras modificar o apreta /start para volver"), "", 0, false, false, nil, false, false, kb)
}

func addTimeToDay(day string, currentTimes []dtos.Times, toAnswer int, bot *teleBot.Bot, db db.DB[dtos.Data], data dtos.Data) {
	t := dtos.Times{
		TimesOfDay: []dtos.Time{},
		Day:        day,
	}
	found := false
	for _, k := range currentTimes {
		if k.Day == day {
			t = k
			found = true
			break
		}
	}
	kb := bot.CreateInlineKeyboard()

	for i := 1; i < 25; i += 1 {
		kb.AddCallbackButtonHandler(strconv.Itoa(i), strconv.Itoa(i), int(math.Floor(float64(i/6))+1), func(u *objs.Update) {
			toAnswer := u.Message.Chat.Id
			from, _ := strconv.Atoi(u.CallbackQuery.Data)
			innerKb := bot.CreateInlineKeyboard()
			for i := from; i < 25; i += 1 {
				innerKb.AddCallbackButtonHandler(strconv.Itoa(i), strconv.Itoa(i), int(math.Floor(float64(i/6))+1), func(u *objs.Update) {
					toAnswer := u.Message.Chat.Id
					to, _ := strconv.Atoi(u.CallbackQuery.Data)
					toAdd := []dtos.Time{
						{
							StartingTime: from,
							FinishTime:   to,
						},
					}
					if found {
						for _, k := range currentTimes {
							if k.Day == day {
								k.TimesOfDay = toAdd
								break
							}
						}
					} else {
						data.FreeTimes = append(data.FreeTimes, dtos.Times{
							TimesOfDay: toAdd,
							Day:        day,
						})
					}
					db.Update(data)
					bot.SendMessage(toAnswer, fmt.Sprintf("se guardo que el %s estas disponible desde %d hasta %d, apreta /start para volver o /add para agregar mas horarios", day, from, to), "", 0, false, false)
				})
			}
			bot.AdvancedMode().ASendMessage(toAnswer, "Elegi hasta donde queres empezar a estar disponible o apreta /start para volver", "", 0, false, false, nil, false, false, innerKb)

		})

	}
	var message string

	if found {
		log.Infof("data is: %+v", t)
		message = dtos.GetItemsMessage(fmt.Sprintf("tus tiempos para el dia %s son", day), []dtos.Times{t})
	} else {
		message = fmt.Sprintf("No tenes nada para el dia: %s", day)
	}
	bot.AdvancedMode().ASendMessage(toAnswer, fmt.Sprintf("%s\nElegi desde donde queres empezar a estar disponible o apreta /start para volver", message), "", 0, false, false, nil, false, false, kb)
}

func getSchedulesSet(c chan dtos.GetInformation, database db.DB[dtos.Data], bot *teleBot.Bot) {
	for {
		toAdd := <-c
		data, err := database.GetAll()
		if err != nil {
			log.Errorf("couldn't get all values, it really seems like an error, should be checked")
			continue
		}
		msg := "Calendario de todos los que dijeron hasta ahora"
		for _, d := range data {
			msg = fmt.Sprintf("%s\n%s", msg, d.String())
		}
		bot.SendMessage(toAdd.ToAnswer, msg, "", 0, false, false)
	}
}
