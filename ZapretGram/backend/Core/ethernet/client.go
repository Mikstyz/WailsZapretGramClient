package ethernet

import (
	tools "ZapretGram/backend/Core/Tools"
	Model "ZapretGram/backend/Core/ethernet/Model"
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
)

type Tcp struct {
	//Подключение и рутина
	Conn net.Conn
	mu   sync.RWMutex

	//база данных подключения
	DB *sql.DB

	//чтение соощений асинх
	messageChan chan *Model.RequestTcp
	done        chan struct{}

	//сопосост запр-ответ
	PendingReqs map[string]chan *Model.RequestTcp

	//id and port server
	IP   string
	Port string

	//ключ шифрования сервера public
	Key *tools.Pubkey

	//данные юзера
	UserId int64
	Name   string
	Token  string

	//chats
	Chats map[string]Model.Chat
}

func NewTcpClient(ip, port string, db *sql.DB) (*Tcp, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		return nil, err
	}

	return &Tcp{
		Conn: conn,
		IP:   ip,
		Port: port,
		DB:   db,
	}, nil
}

func (c *Tcp) Disconnect() {
	if c.Conn != nil {
		c.Conn.Close()
		c.Conn = nil
	}
}

func (c *Tcp) Reconnect() error {
	c.Disconnect()
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", c.IP, c.Port))
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}

func (c *Tcp) RequestTcp(req Model.RequestTcp, pubkey *tools.Pubkey, result interface{}) error {
	if c.Conn == nil {
		return fmt.Errorf("conn not found")
	}

	// Шифруем запрос
	b, err := pubkey.EncPublicKey(req)
	if err != nil {
		return err
	}

	log.Printf("%d", b)

	// Отправляем в сокет с \n, чтобы клиент мог читать ReadBytes('\n')
	_, err = c.Conn.Write(append(b, '\n'))
	if err != nil {
		return err
	}

	// Читаем ответ
	reader := bufio.NewReader(c.Conn)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return err
	}

	// Расшифровываем
	var resp Model.ResponseTcp
	if err := pubkey.DecPublicKey(line[:len(line)-1], &resp); err != nil {
		return err
	}

	// Если есть result, разбираем Data
	if result != nil && resp.Data != nil {
		dataBytes, err := json.Marshal(resp.Data)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(dataBytes, result); err != nil {
			return err
		}
	}

	return nil
}

func (c *Tcp) ListenMessages(pubkey *tools.Pubkey) {
	reader := bufio.NewReader(c.Conn)
	go func() {
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				fmt.Println("Соединение закрыто:", err)
				return
			}

			var resp Model.ResponseTcp
			if err := pubkey.DecPublicKey(line[:len(line)-1], &resp); err != nil {
				fmt.Println("Ошибка расшифровки:", err)
				continue
			}

			// Обработка по Action
			switch resp.Action {

			//chat
			case "chat":
				var chat Model.ResponseChatData
				dataBytes, _ := json.Marshal(resp.Data)
				json.Unmarshal(dataBytes, &chat)
				fmt.Printf("Новое сообщение от %s: %s\n", chat.ChatId, chat.Text)

			//Успешная авторизация
			case "auth":
				var auth Model.ResponseAuthData
				dataBytes, _ := json.Marshal(resp.Data)
				json.Unmarshal(dataBytes, &auth)

				if resp.Status != "ok" {
					fmt.Printf("не удачный вход на сервер, не верный login or password")
				}
				fmt.Printf("успешный вход на сервер. ID: %s\n", auth.Token)
				fmt.Printf("auth status: %d", resp.Status)

			case "error":
				var errData Model.ResponseErrorData
				dataBytes, _ := json.Marshal(resp.Data)
				json.Unmarshal(dataBytes, &errData)
				fmt.Printf("Ошибка %d: %s\n", errData.ErrorCode, errData.Details)
			default:
				fmt.Println("Неизвестный action:", resp.Action)
			}
		}
	}()
}
