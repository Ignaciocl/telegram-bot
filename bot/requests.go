package bot

import (
	"telegram-bot/dtos"
)

func StartHandlersOperations(newsBot *NewsBot) map[string]chan dtos.GetInformation {

	// ToDo there HAS to be a better way to do this
	return map[string]chan dtos.GetInformation{}
}

func getNewChan() chan dtos.GetInformation {
	return make(chan dtos.GetInformation, 100)
}
