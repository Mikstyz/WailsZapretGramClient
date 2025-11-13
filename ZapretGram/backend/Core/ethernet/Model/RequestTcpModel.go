package model

type RequestTcp struct {
	// request info
	Action   string      `json:"action"`             //Тип запроса [chat, register, login]
	DateTime string      `json:"datetime,omitempty"` //Время запроса
	CorrId   string      `json:"correlation_id,omitempty"`
	Data     interface{} `json:"data,omitempty"` //Нужные нам данные
}

type RequestChatData struct {
	//chat
	ChatId int64  `json:"chatid,omitempty"` //Чат
	Text   string `json:"text,omitempty"`   //Сообщение
}

type RequestMessage struct {
	UserId  int64  `json:"usreid"`
	ChatId  int64  `json:"chatid"`
	Message string `json:"message"`
}

type RequestNewChata struct {
	//new chat
	CratorId int64  `json:"cratorid,omitempty"`
	UserName string `json:"UserName,omitempty"`
}

type RequestAuthData struct {
	//Регистрация и вход
	UserIn     string `json:"userin,omitempty"`   //регистрация и вход пользователя
	PasswordIn string `json:"password,omitempty"` //пароль для входи или регестрации
}
