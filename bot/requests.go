package bot

import (
	"bot-telegram/db"
	"bot-telegram/dtos"
	"bot-telegram/services/news"
	"bot-telegram/utils"
	"errors"
	"fmt"
	teleBot "github.com/SakoDroid/telego"
	objs "github.com/SakoDroid/telego/objects"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

const (
	CategoriesWanted    = "categoriesWanted"
	DeleteData          = "deleteData"
	RemoveCategories    = "removeCategories"
	AddCategories       = "addCategories"
	GetNew              = "getNews"
	scheduleNewsChannel = "scheduleNews"
)

var PossibleCategories = []string{"business", "entertainment", "environment", "food", "health", "politics", "science", "sports", "technology", "top", "tourism", "world"}

func StartHandlersOperations(newsBot *NewsBot) map[string]chan dtos.GetInformation {
	deleteChan := getNewChan()
	getWhitelistChan := getNewChan()
	removeCategoryChan := getNewChan()
	addCategoryChan := getNewChan()
	getNewsChan := getNewChan()
	scheduleNewsChan := getNewChan()
	go deleteData(deleteChan, newsBot.DB, newsBot.TelegramBot)
	go getLikedTopics(getWhitelistChan, newsBot.DB, newsBot.TelegramBot)
	go removeCategory(removeCategoryChan, newsBot.DB, newsBot.TelegramBot)
	go addCategory(addCategoryChan, newsBot.DB, newsBot.TelegramBot)
	go getWantedNews(getNewsChan, newsBot.DB, newsBot.TelegramBot, newsBot.GPTService)
	go scheduleNews(scheduleNewsChan, newsBot.ScheduleDB, newsBot.TelegramBot)

	// ToDo there HAS to be a better way to do this
	return map[string]chan dtos.GetInformation{
		DeleteData:          deleteChan,
		CategoriesWanted:    getWhitelistChan,
		RemoveCategories:    removeCategoryChan,
		AddCategories:       addCategoryChan,
		GetNew:              getNewsChan,
		scheduleNewsChannel: scheduleNewsChan,
	}
}

func getNewChan() chan dtos.GetInformation {
	return make(chan dtos.GetInformation, 100)
}

func deleteData(c chan dtos.GetInformation, database db.DB[dtos.Data], bot *teleBot.Bot) {
	for {
		toRemove := <-c
		fmt.Printf("delete was called with id: %d\n", toRemove)
		database.Delete(toRemove.Id)
		bot.SendMessage(toRemove.ToAnswer, "your data was deleted, to start using our bot again send /start", "", 0, false, false)
	}
}

// getLikedTopics returns the liked topics for the given user
func getLikedTopics(c chan dtos.GetInformation, database db.DB[dtos.Data], bot *teleBot.Bot) {
	for {
		userInfo := <-c
		fmt.Printf("get liked topics was called with ID: %d\n", userInfo)
		data, _ := database.Get(userInfo.Id)
		formattedMessage := utils.GetItemsMessage("Your liked topics are", data.WantedNews)
		_, err := bot.SendMessage(userInfo.ToAnswer, formattedMessage+"\nPress /start to see menu again", "markdown", 0, false, false)
		if err != nil {
			fmt.Printf("error sending liked topics for user %v: %v\n", userInfo.Id, err)
		}
	}
}

func addCategory(c chan dtos.GetInformation, db db.DB[dtos.Data], bot *teleBot.Bot) {
	for {
		toAddCategory := <-c
		log.Infof("checking of adding category, %d", toAddCategory.Id)
		go addCategoryForChat(toAddCategory, db, bot)
		bot.SendMessage(toAddCategory.ToAnswer, "processing request", "", 0, true, false)
	}
}

func addCategoryForChat(info dtos.GetInformation, db db.DB[dtos.Data], bot *teleBot.Bot) {
	data, err := db.Get(info.Id)
	if err != nil {
		log.Errorf("couldn't get info for user, err: %v", err)
	}
	kb := bot.CreateInlineKeyboard()
	for i, k := range PossibleCategories {
		if !utils.Contains(k, data.WantedNews) {
			kb.AddCallbackButtonHandler(k, k, i%4+1, func(update *objs.Update) {
				d, err := db.Get(info.Id)
				key := update.CallbackQuery.Data
				if err != nil {
					log.Errorf("couldn't get info for user, err: %v", err)
					return
				}
				if utils.Contains(key, d.WantedNews) {
					bot.SendMessage(info.ToAnswer, fmt.Sprintf("you already added this category to your wishlist: %s", key), "", 0, false, false)
				} else if len(d.WantedNews) > 4 {
					bot.SendMessage(info.ToAnswer, "you have too many categories already", "", 0, false, false)
				} else {
					d.WantedNews = append(d.WantedNews, key)
					db.Update(d)

					message := fmt.Sprintf("added %s to the categories wanted. Current liked categories are", key)
					formattedMessage := utils.GetItemsMessage(message, d.WantedNews)
					bot.SendMessage(info.ToAnswer, formattedMessage, "", 0, false, false)
				}
			})
		}
	}
	bot.AdvancedMode().ASendMessage(info.ToAnswer, fmt.Sprintf("Select options below to add categories, current selected are: %v", data.WantedNews), "", 0, false, false, nil, false, false, kb)
}

func removeCategory(c chan dtos.GetInformation, db db.DB[dtos.Data], bot *teleBot.Bot) {
	for {
		toAddCategory := <-c
		log.Infof("checking of removing category, %d", toAddCategory.Id)
		go removeCategoryForChat(toAddCategory, db, bot)
		bot.SendMessage(toAddCategory.ToAnswer, "processing request", "", 0, true, false)
	}
}

func removeCategoryForChat(info dtos.GetInformation, db db.DB[dtos.Data], bot *teleBot.Bot) {
	data, err := db.Get(info.Id)
	if err != nil {
		log.Errorf("couldn't get info for user, err: %v", err)
	}
	kb := bot.CreateInlineKeyboard()
	for i, k := range PossibleCategories {
		if utils.Contains(k, data.WantedNews) {
			kb.AddCallbackButtonHandler(k, k, i+1, func(update *objs.Update) {
				d, err := db.Get(info.Id)
				key := update.CallbackQuery.Data
				if err != nil {
					log.Errorf("couldn't get info for user, err: %v", err)
					return
				}
				if !utils.Contains(key, d.WantedNews) {
					message := fmt.Sprintf("You already have removed this category from your wanted topics: %s. Your current wanted topics are", key)
					formattedMessage := utils.GetItemsMessage(message, d.WantedNews)
					bot.SendMessage(info.ToAnswer, formattedMessage, "Markdown", 0, false, false)
				} else {
					s := make([]string, 0, len(d.WantedNews)-1)
					for _, k := range d.WantedNews {
						if k != key {
							s = append(s, k)
						}
					}
					d.WantedNews = s
					db.Update(d)
					message := fmt.Sprintf("Category '%s' was removed correctly!. Your current wanted topics are", key)
					formattedMessage := utils.GetItemsMessage(message, d.WantedNews)
					bot.SendMessage(info.ToAnswer, formattedMessage, "", 0, false, false)
				}
			})
		}
	}
	message := fmt.Sprintf("Select options below to remove categories, current selected are")
	formattedMessage := utils.GetItemsMessage(message, data.WantedNews)
	bot.AdvancedMode().ASendMessage(info.ToAnswer, formattedMessage, "", 0, false, false, nil, false, false, kb)
}

func getWantedNews(c chan dtos.GetInformation, db db.DB[dtos.Data], bot *teleBot.Bot, gptService news.Provider) {
	for {
		toReturnNews := <-c
		log.Infof("user asking for news, %d", toReturnNews.Id)
		go getWantedNewsForUser(toReturnNews, db, bot, gptService)
		bot.SendMessage(toReturnNews.ToAnswer, "processing request", "", 0, true, false)
	}
}

func getWantedNewsForUser(info dtos.GetInformation, db db.DB[dtos.Data], bot *teleBot.Bot, gptService news.Provider) {
	data, err := db.Get(info.Id)
	topicsWanted := data.WantedNews
	if err != nil || len(topicsWanted) == 0 {
		bot.SendMessage(info.ToAnswer, "Topic not selected, searching for no category", "", 0, false, false)
	}
	n, err := news.GetNew(strings.Join(topicsWanted, ","))
	if err != nil {
		bot.SendMessage(info.ToAnswer, "An error has occurred, please try again later.", "", 0, false, false)
	} else {
		kb := bot.CreateInlineKeyboard()
		message := news.GetSummarizedMessage(n, gptService)
		kb.AddURLButton("Go to the news", n.Url, 1)
		bot.AdvancedMode().ASendMessage(info.ToAnswer, message, "markdown", 0, false, false, nil, false, false, kb)
	}
}

// scheduleNews schedules the hour at which a user want news
func scheduleNews(scheduleChannel chan dtos.GetInformation, scheduleDB db.DB[dtos.Schedule], bot *teleBot.Bot) {
	for {
		userInfo := <-scheduleChannel
		hoursKeyboard := bot.CreateInlineKeyboard()
		rowIdx := 1
		for hour := 0; hour <= 23; hour++ {
			if hour%6 == 0 {
				rowIdx += 1
			}
			hoursKeyboard.AddCallbackButtonHandler(
				fmt.Sprintf("%v", hour),
				fmt.Sprintf("%v", hour),
				rowIdx,
				func(update *objs.Update) {
					saveScheduleNewsHour(userInfo, update.CallbackQuery.Data, scheduleDB, bot)
				},
			)
		}
		_, err := bot.AdvancedMode().ASendMessage(userInfo.ToAnswer, "Select an hour to schedule the news", "", 0, false, false, nil, false, false, hoursKeyboard)
		if err != nil {
			panic(fmt.Sprintf("error sending hours keyboard: %v", err))
		}
	}
}

// saveScheduleNewsHour saves in the database the hour that the user have selected
func saveScheduleNewsHour(userInfo dtos.GetInformation, selectedHour string, scheduleDB db.DB[dtos.Schedule], bot *teleBot.Bot) {
	hour, err := strconv.Atoi(selectedHour)
	if err != nil {
		panic("Invalid format for hour")
	}

	scheduleData, err := scheduleDB.Get(hour)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		panic(fmt.Sprintf("Error getting data from Schedule DB: %v", err))
	}

	// We have to insert the new row
	if errors.Is(err, db.ErrNotFound) {
		newScheduleData := dtos.Schedule{
			HourID:    hour,
			UsersInfo: map[int]int{userInfo.Id: userInfo.ToAnswer},
		}
		scheduleDB.Insert(newScheduleData)
		_, err = bot.SendMessage(userInfo.ToAnswer, fmt.Sprintf("Hour scheduled correctly! You will receive news at %v hs", selectedHour), "", 0, false, false)
		if err != nil {
			panic(fmt.Sprintf("error sending response about scheduled hour: %v", err))
		}
		return
	}

	chatID, ok := scheduleData.UsersInfo[userInfo.Id]
	if !ok {
		scheduleData.UsersInfo[userInfo.Id] = userInfo.ToAnswer
		scheduleDB.Update(scheduleData)
		_, err = bot.SendMessage(userInfo.ToAnswer, fmt.Sprintf("Hour scheduled correctly! You will receive news at %v hs", selectedHour), "", 0, false, false)
		if err != nil {
			panic(fmt.Sprintf("error sending response about scheduled hour: %v", err))
		}
		return
	}

	_, err = bot.SendMessage(chatID, "You already have scheduled this hour "+selectedHour, "", 0, false, false)
	if err != nil {
		panic(fmt.Sprintf("error sending response about already scheduled hour: %v", err))
	}
}
