package ethernet

import (
	tools "ZapretGram/backend/Core/Tools"
	Model "ZapretGram/backend/Core/ethernet/Model"
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
)

func (c *TcpClient) RiderTcp(pubkey *tools.Pubkey) error {
	if c.Conn == nil {
		return fmt.Errorf("conn not found")
	}

	go c.readerLoop(pubkey)
	return nil
}

func (c *TcpClient) processIncomigMessage(pubkey *tools.Pubkey, data []byte) {
	var message Model.RequestTcp
	var b []byte

	fmt.Printf("data: %d\n", data)
	pubkey.DecPublicKey(data, &b)

	// обрезаем возможные нулевые байты после расшифровки
	b = bytes.Trim(b, "\x00")

	if len(b) == 0 {
		fmt.Println("Пустое сообщение после расшифровки")
		return
	}

	// теперь анмаршалим **только расшифрованные данные**
	if err := json.Unmarshal(b, &message); err != nil {
		fmt.Printf("Ошибка расшифровки основной сущности: %v\n", err)
		return
	}

	if message.Action == "" {
		fmt.Println("не задан тип сообщения")
		return
	}

	fmt.Printf("===========\n[reader] получено сообщение от сервера %s\n===========\n", message.Action)

	switch strings.ToLower(message.Action) {
	case "register":
		log.Print("register message in server")
	case "message":
		var msg Model.MessageInChat

		// маршалим Data обратно в json
		dataBytes, err := json.Marshal(message.Data)
		if err != nil {
			fmt.Printf("Ошибка маршала Data: %v\n", err)
			return
		}

		// анмаршалим в конкретную структуру
		if err := json.Unmarshal(dataBytes, &msg); err != nil {
			fmt.Printf("Ошибка анмаршала MessageInChat: %v\n", err)
			return
		}

		fmt.Printf("Сообщение - chatid %d, messageid=%d,userId=%d, текст='%s'\n",
			msg.ChatId, msg.Id, msg.UserId, msg.Message)

		// Отправляем обратно в ожидание ответа (если нужно)
		c.routeToPendingRequest(&message)

	default:
		log.Print("неизвестное сообщение")
	}
}

func (c *TcpClient) routeToPendingRequest(message *Model.RequestTcp) {
	if message.CorrId == "" {
		fmt.Println("corr error")
		return
	}

	c.mu.Lock() // блокируем на всю операцию
	responseChat, exists := c.PendingReqs[message.CorrId]
	if exists {
		delete(c.PendingReqs, message.CorrId)
	}
	c.mu.Unlock()

	if exists {
		responseChat <- message
	}
}

func (c *TcpClient) readerLoop(pubkey *tools.Pubkey) {
	reader := bufio.NewReader(c.Conn)
	const maxMsgLen = 100 * 1024 * 1024 // 100 MB

	for {
		select {
		case <-c.done:
			return
		default:
			var lenBuf [4]byte
			if _, err := io.ReadFull(reader, lenBuf[:]); err != nil {
				if err != io.EOF {
					fmt.Printf("Ошибка чтения длины от %s: %v\n", c.Conn.RemoteAddr(), err)
				}
				c.Conn.Close()
				c.Conn = nil
				return
			}

			msgLen := binary.BigEndian.Uint32(lenBuf[:])
			if msgLen == 0 || msgLen > maxMsgLen {
				fmt.Printf("Сообщение слишком большое: %d bytes\n", msgLen)
				c.Conn.Close()
				c.Conn = nil
				return
			}

			msgBuf := make([]byte, msgLen)
			if _, err := io.ReadFull(reader, msgBuf); err != nil {
				fmt.Printf("Ошибка чтения данных от %s: %v\n", c.Conn.RemoteAddr(), err)
				c.Conn.Close()
				c.Conn = nil
				return
			}

			// обработка без go, чтобы не ломать очередь сообщений
			c.processIncomigMessage(pubkey, msgBuf)
		}
	}
}
