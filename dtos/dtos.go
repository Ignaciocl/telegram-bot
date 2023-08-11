package dtos

import (
	"fmt"
	"strings"
)

type Times struct {
	TimesOfDay []Time `json:"times"`
	Day        string `json:"day"`
}

type Time struct {
	StartingTime int `json:"starting_time"`
	FinishTime   int `json:"finish_time"`
}

type Data struct {
	Id        int     `json:"id"`
	FreeTimes []Times `json:"free_times"`
	Timezone  string  `json:"timezone"`
	Name      string  `json:"name"`
}

func (d Data) GetPrimaryKey() int {
	return d.Id
}

type Chat struct {
	Id       int // Id from whom the request started
	ToAnswer int //Chat to answer
}

type DeleteDataInformation struct {
	Id       int
	ToAnswer int
}

type GetInformation struct {
	Id       int
	ToAnswer int
}

type UserInfo struct {
	UserID int
	ChatID int
}

func (d Data) String() string {
	return GetItemsMessage(fmt.Sprintf("Los tiempos para %s son", d.Name), d.FreeTimes)
}

func GetItemsMessage(message string, dates []Times) string {
	formattedMessage := fmt.Sprintf("%s:", message)
	datesFormatted := make([]string, 0)
	for _, k := range dates {
		day := k.Day
		times := k.TimesOfDay
		for _, t := range times {
			from, to := t.StartingTime, t.FinishTime
			day = fmt.Sprintf("%s\n\t\tdesde: %d hasta: %d", day, from, to)
		}
		datesFormatted = append(datesFormatted, day)
	}
	formattedItems := strings.Join(datesFormatted, "\n\t+ ")
	formattedItems += "\n"
	return formattedMessage + formattedItems
}
