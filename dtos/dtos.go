package dtos

type Data struct {
	Id            int      `json:"id"`
	WantedNews    []string `json:"wanted_news"`
	OmittedTopics []string `json:"omitted_topics"`
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

type Schedule struct {
	HourID    int         `json:"id"`
	UsersInfo map[int]int `json:"users_info"`
}

type UserInfo struct {
	UserID int
	ChatID int
}

func (s Schedule) GetPrimaryKey() int {
	return s.HourID
}
