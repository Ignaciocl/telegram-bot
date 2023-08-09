package dtos

type Times struct {
	StartingTime int `json:"starting_time"`
	FinishTime   int `json:"finish_time"`
}

type Data struct {
	Id        int     `json:"id"`
	FreeTimes []Times `json:"free_times"`
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
