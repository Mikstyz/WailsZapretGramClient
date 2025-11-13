package ethernet

import (
	"ZapretGram/backend/Core/Tools"
	Model "ZapretGram/backend/Core/ethernet/Model"
	model "ZapretGram/backend/Core/ethernet/Model"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
)

type TcpRequest struct {
	Tcp *Tcp
}

func (t *TcpRequest) RiderTcp(Pubkey string) {
	panic("unimplemented")
}

func NewRequest(tcp *Tcp, key string) *TcpRequest {
	fmt.Printf("conn: %+v\n", tcp)
	if tcp.Conn == nil {
		fmt.Println("Невозможно отправить запрос, соедененеие закрыто")
		return nil
	}

	tcp.Key = Tools.NewKey(key)

	var Tcpclient = &TcpRequest{
		Tcp: tcp,
	}

	return Tcpclient
}

func (tcp *Tcp) sendRequest(b []byte) {
	// Отправляем сначала длину сообщения, потом сами данные
	length := uint32(len(b))

	fmt.Printf("len message: %d\n", length)

	buf := make([]byte, 4+len(b))
	binary.BigEndian.PutUint32(buf[:4], length)
	copy(buf[4:], b)

	if _, err := tcp.Conn.Write(buf); err != nil {
		fmt.Printf("Ошибка при отправке сообщения: %v\n", err)
		return
	}
}

func (tcp *Tcp) readingAnsfer() *Model.ResponseTcp {
	// Чтение ответа
	var lenBuf [4]byte
	if _, err := io.ReadFull(tcp.Conn, lenBuf[:]); err != nil {
		fmt.Printf("Ошибка чтения длины ответа: %v\n", err)
		return nil
	}

	respLen := binary.BigEndian.Uint32(lenBuf[:])
	respBuf := make([]byte, respLen)
	if _, err := io.ReadFull(tcp.Conn, respBuf); err != nil {
		fmt.Printf("Ошибка чтения ответа: %v\n", err)
		return nil
	}

	var res Model.ResponseTcp
	if err := tcp.Key.DecPublicKey(respBuf, &res); err != nil {
		fmt.Printf("Ошибка расшифровки ответа: %v\n", err)
		return nil
	}

	return &res
}

func (t *Tcp) Ping() bool {
	ping := Model.RequestTcp{
		Action: "ping",
	}

	b, err := t.Key.EncPublicKey(ping)
	if err != nil {
		fmt.Errorf("Ошибка шифрования: %v\n", err)
	}
	t.sendRequest(b)

	res := t.readingAnsfer()

	fmt.Printf("res server: %+d\n", res)
	return res.Status == "ok"
}

func (t *TcpRequest) Auth(log string, pass string, action string) error {
	if action != "register" && action != "login" {
		fmt.Println("неверный тип запроса only register of login")
		return fmt.Errorf("неверный тип запроса only register of login\n")
	}

	//модель юзера
	user := Model.RequestTcp{
		Action:   action,
		DateTime: time.Now().Format(time.RFC3339),
		Data: Model.RequestAuthData{
			UserIn:     log,
			PasswordIn: pass,
		},
	}

	//Шифруем сообщение

	b, err := t.Tcp.Key.EncPublicKey(user)

	fmt.Printf("b: %d", b)

	if err != nil {
		fmt.Errorf("ошибка шифрования: %v\n", err)
	}

	//Отправляем длину сообщения
	t.Tcp.sendRequest(b)

	//читаем ответ
	dataBytes := t.Tcp.readingAnsfer()

	fmt.Printf("res server: %+v\n", dataBytes)

	// Сериализуем поле Data обратно в json
	dataRaw, err := json.Marshal(dataBytes.Data)
	if err != nil {
		return fmt.Errorf("ошибка маршала Data: %v", err)
	}

	// Парсим уже конкретный тип auth
	var responseData Model.ResponseAuthData
	if err := json.Unmarshal(dataRaw, &responseData); err != nil {
		return fmt.Errorf("ошибка парсинга auth data: %v", err)
	}
	// Сохраняем значения
	t.Tcp.Name = responseData.UserName
	t.Tcp.UserId = responseData.UserId
	t.Tcp.Token = responseData.Token
	t.Tcp.Chats = responseData.Chats

	fmt.Printf("Auth %s: user=%s id=%d token=%s, len chats:%d \n", dataBytes.Status, responseData.UserName, responseData.UserId, responseData.Token, len(responseData.Chats))

	fmt.Print("\n====chats===================\n")
	for key, chat := range t.Tcp.Chats {
		fmt.Printf("user=%s, chat=%+v\n", key, chat)
	}
	fmt.Print("============================\n")

	return nil
}

func (t *TcpRequest) NewChat(recipient string) error {
	fmt.Println("newchat request")
	request := model.RequestTcp{
		Action:   "newchat",
		DateTime: time.Now().Format(time.RFC3339),
		Data: model.RequestNewChata{
			CratorId: t.Tcp.UserId,
			UserName: recipient,
		},
	}

	fmt.Println("создана новая сущность чата")

	b, err := t.Tcp.Key.EncPublicKey(request)

	if err != nil {
		fmt.Errorf("ошибка создания нового чата: %w\n", err)
		return err
	}

	//запрос на сервер
	fmt.Println("запрос на сервер")
	t.Tcp.sendRequest(b)
	//читаем ответ
	fmt.Println("читаем ответ")
	dataBytes := t.Tcp.readingAnsfer()

	fmt.Printf("res server: %+v\n", dataBytes)

	// Сериализуем поле Data обратно в json
	dataRaw, err := json.Marshal(dataBytes.Data)
	if err != nil {
		return fmt.Errorf("ошибка маршала Data: %v\n", err)
	}

	var responseData model.ResponseNewChata
	if err := json.Unmarshal(dataRaw, &responseData); err != nil {
		return fmt.Errorf("ошибка парсинга auth data: %v\n", err)
	}

	if responseData.ChatId == 0 {
		return fmt.Errorf("ошибка при создании нового чата, chatid is nil\n")
	}

	if t.Tcp.Chats == nil {
		t.Tcp.Chats = make(map[string]model.Chat)
	}

	t.Tcp.Chats[strconv.FormatInt(responseData.ChatId, 10)] = model.Chat{
		Id: responseData.ChatId,
	}

	return nil
}

func (t *TcpRequest) NewMessage(ChatId int64, message string) error {
	fmt.Println("newmessage request")
	request := model.RequestTcp{
		Action:   "message",
		CorrId:   Tools.GenerateUUID(),
		DateTime: time.Now().Format(time.RFC3339),
		Data: model.MessageInChat{
			UserId:  t.Tcp.UserId,
			ChatId:  ChatId,
			Message: message,
		},
	}

	b, err := t.Tcp.Key.EncPublicKey(request)

	if err != nil {
		fmt.Errorf("ошибка создания нового чата: %w\n", err)
		return err
	}

	//запрос на сервер
	fmt.Println("запрос на сервер")
	t.Tcp.sendRequest(b)

	//читаем ответ
	fmt.Println("читаем ответ")
	dataBytes := t.Tcp.readingAnsfer()

	fmt.Printf("res server: %+v\n", dataBytes)

	// Сериализуем поле Data обратно в json
	dataRaw, err := json.Marshal(dataBytes.Data)
	if err != nil {
		return fmt.Errorf("ошибка маршала Data: %v\n", err)
	}

	var responsedata model.ResponseChat
	if err := json.Unmarshal(dataRaw, &responsedata); err != nil {
		return fmt.Errorf("ошибка парсинга message data: %v\n", err)
	}

	//тут будет кидать в буффер чата

	//вывод сообщения в чат
	fmt.Printf("chatId%d, messageid:%d\n", responsedata.ChatId, responsedata.MessageId)

	return nil
}
