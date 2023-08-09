package bot

import (
	"bot-telegram/dtos"
	"fmt"
	objs "github.com/SakoDroid/telego/objects"
)

func (nb *NewsBot) StartGoRoutines() error {
	m := StartHandlersOperations(nb)
	nb.channels = m
	return nil
}

func (nb *NewsBot) StartHandlers() error {
	bot := nb.TelegramBot
	bot.AddHandler("/hi", func(u *objs.Update) {
		fmt.Println("hi was called")
		bot.SendMessage(u.Message.Chat.Id, "hello", "", u.Message.MessageId, false, false)
	}, "private")

	bot.AddHandler("/start", func(u *objs.Update) {
		kb := bot.CreateInlineKeyboard()
		kb.AddCallbackButtonHandler("Delete all your data", "/hi", 2, func(update *objs.Update) {
			fmt.Println("delete was clicked")
			toRemove := dtos.GetInformation{
				Id:       u.Message.From.Id,
				ToAnswer: u.Message.Chat.Id,
			}
			bot.SendMessage(u.Message.Chat.Id, "Your data is being removed...", "", u.Message.MessageId, false, false)
			nb.channels[DeleteData] <- toRemove
		})
		kb.AddCallbackButtonHandler("Check categories added", "/hi1", 3, func(update *objs.Update) {
			toCheck := dtos.GetInformation{
				Id:       u.Message.From.Id,
				ToAnswer: u.Message.Chat.Id,
			}
			nb.channels[CategoriesWanted] <- toCheck
		})
		kb.AddCallbackButtonHandler("Remove categories", "/hi2", 5, func(update *objs.Update) {
			toRemove := dtos.GetInformation{
				Id:       u.Message.From.Id,
				ToAnswer: u.Message.Chat.Id,
			}
			nb.channels[RemoveCategories] <- toRemove
		})
		kb.AddCallbackButtonHandler("Get a summarized news about your interests", "/getWantedNews", 4, func(update *objs.Update) {
			toRetrieveNews := dtos.GetInformation{
				Id:       u.Message.From.Id,
				ToAnswer: u.Message.Chat.Id,
			}
			nb.channels[GetNew] <- toRetrieveNews
		})
		kb.AddCallbackButtonHandler("Add categories", "/summarize", 5, func(update *objs.Update) {
			toRemove := dtos.GetInformation{
				Id:       u.Message.From.Id,
				ToAnswer: u.Message.Chat.Id,
			}
			nb.channels[AddCategories] <- toRemove
		})
		kb.AddCallbackButtonHandler("Schedule News", "/scheduleNews", 6, func(update *objs.Update) {
			userInfo := dtos.GetInformation{
				Id:       u.Message.From.Id,
				ToAnswer: u.Message.Chat.Id,
			}
			nb.channels[scheduleNewsChannel] <- userInfo
		})

		di := dtos.Data{
			OmittedTopics: []string{},
			WantedNews:    []string{},
			Id:            u.Message.From.Id,
		}

		nb.DB.Insert(di)

		//Sends the message to the chat that the message has been received from. The message will be a reply to the received message.
		_, err := bot.AdvancedMode().ASendMessage(u.Message.Chat.Id, "Please select one of the options below.", "", u.Message.MessageId, false, false, nil, false, false, kb)
		if err != nil {
			fmt.Printf("error happened, %v\n", err)
		}
	}, "private")
	return nil
}
