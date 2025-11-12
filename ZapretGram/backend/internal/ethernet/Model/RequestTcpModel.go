package model

type RequestTcp struct {
	// request info
	Action   string      `json:"Action"`             //Тип запроса [chat, register, login]
	DateTime string      `json:"datetime,omitempty"` //Время запроса
	Data     interface{} `json:"data,omitempty"`     //Нужные нам данные
}

type RequestChatData struct {
	//chat
	ChatId int64  `json:"chatid,omitempty"` //Чат
	Text   string `json:"text,omitempty"`   //Сообщение
}

type RequestAuthData struct {
	//Регистрация и вход
	UserIn     string `json:"userin,omitempty"`   //регистрация и вход пользователя
	PasswordIn string `json:"password,omitempty"` //пароль для входи или регестрации
}
