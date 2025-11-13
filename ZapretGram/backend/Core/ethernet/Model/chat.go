package model

type Chat struct {
	Id int64 `json:"Id"`
}

type MessageInChat struct {
	Id      int64  `json:"meesageid"`
	UserId  int64  `json:"userId"`
	ChatId  int64  `json:"chatid"`
	Message string `json:"message"`
}

type ResponseChat struct {
	ChatId    int64 `json:"chatid"`
	MessageId int64 `json:"messageId"`
}
