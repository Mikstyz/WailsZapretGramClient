package ethernet

import (
	tools "ZapretGram/backend/Core/Tools"
	model "ZapretGram/backend/Core/ethernet/Model"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
)

type TcpClient struct {
	//Подключение и рутина
	Conn net.Conn
	mu   sync.RWMutex

	//чтение соощений асинх
	messageChan chan *ServerMessage
	done        chan struct{}

	//сопосост запр-ответ
	PendingReqs map[string]chan *ServerMessage

	//id and port server
	IP   string
	Port string

	//ключ шифрования сервера public
	Key *tools.Pubkey
}

func NewTcpClient(ip, port string) (*TcpClient, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		return nil, err
	}

	return &TcpClient{
		Conn: conn,
		IP:   ip,
		Port: port,
	}, nil
}

func (c *TcpClient) Disconnect() {
	if c.Conn != nil {
		c.Conn.Close()
		c.Conn = nil
	}
}

func (c *TcpClient) Reconnect() error {
	c.Disconnect()
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", c.IP, c.Port))
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}

func (c *TcpClient) RequestTcp(req model.RequestTcp, pubkey *tools.Pubkey, result interface{}) error {
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
	var resp model.ResponseTcp
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

func (c *TcpClient) ListenMessages(pubkey *tools.Pubkey) {
	reader := bufio.NewReader(c.Conn)
	go func() {
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				fmt.Println("Соединение закрыто:", err)
				return
			}

			var resp model.ResponseTcp
			if err := pubkey.DecPublicKey(line[:len(line)-1], &resp); err != nil {
				fmt.Println("Ошибка расшифровки:", err)
				continue
			}

			// Обработка по Action
			switch resp.Action {

			//chat
			case "chat":
				var chat model.ResponseChatData
				dataBytes, _ := json.Marshal(resp.Data)
				json.Unmarshal(dataBytes, &chat)
				fmt.Printf("Новое сообщение от %s: %s\n", chat.ChatId, chat.Text)

			//Успешная авторизация
			case "auth":
				var auth model.ResponseAuthData
				dataBytes, _ := json.Marshal(resp.Data)
				json.Unmarshal(dataBytes, &auth)

				if !auth.Status {
					fmt.Printf("не удачный вход на сервер, не верный login or password")
				}
				fmt.Printf("успешный вход на сервер. ID: %s\n", auth.Token)
				fmt.Printf("auth status: %d", auth.Status)

			case "error":
				var errData model.ResponseErrorData
				dataBytes, _ := json.Marshal(resp.Data)
				json.Unmarshal(dataBytes, &errData)
				fmt.Printf("Ошибка %d: %s\n", errData.ErrorCode, errData.Details)
			default:
				fmt.Println("Неизвестный action:", resp.Action)
			}
		}
	}()
}
