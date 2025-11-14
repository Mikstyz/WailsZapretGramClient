package model

type ResponseTcp struct {
	Action string      `json:"action"`           // Тип ответа (соответствует Action запроса)
	Status string      `json:"status,omitempty"` // ok / error / fail
	CorrId string      `json:"correlation_id,omitempty"`
	Data   interface{} `json:"data,omitempty"` // payload, зависит от Action
}

// ответ от чата
type ResponseChatData struct {
	ChatId int64  `json:"chatid,omitempty"` // chatid
	Text   string `json:"text,omitempty"`   // Текст сообщения
}

type ResponseNewChata struct {
	//new chat
	ChatId int64 `json:"chatid,omitempty"`
}

// ответ от регестрации или логина
type ResponseAuthData struct {
	Token    string          `json:"token,omitempty"`
	UserName string          `json:"username,omitempty"`
	UserId   int64           `json:"userid,omitempty"`
	Chats    map[string]Chat `json:"chats,omitempty"`
}

type ResponseNewChat struct {
	ChatId int64 `json:"chatid,omitempty"`
}

// Ошибка от сервера
type ResponseErrorData struct {
	ErrorCode int    `json:"errorcode,omitempty"` //код ошибки
	Details   string `json:"details,omitempty"`   //сообщение об ошибке
}

// Ответ от сервера
type ResponseServerData struct {
	Status  int    `json:"status,omitempty"`  //код ошибки
	Message string `json:"message,omitempty"` //сообщение об ошибке
}

// Просто чтобы было что скопирывать
type Response struct {
}
