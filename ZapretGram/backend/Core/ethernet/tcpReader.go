package ethernet

import (
	tools "ZapretGram/backend/Core/Tools"
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type ServerMessage struct {
	Action string      `json:"action"`
	CorrId string      `json:"correlation_id,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

type RequestMessage struct {
	Action string      `json:"action"`           // Тип ответа (соответствует Action запроса)
	Status string      `json:"status,omitempty"` // ok / error / fail
	Data   interface{} `json:"data,omitempty"`   // payload, зависит от Action
}

func (c *TcpClient) RiderTcp(pubkey *tools.Pubkey) error {
	if c.Conn == nil {
		return fmt.Errorf("conn not found")
	}

	go c.readerLoop(pubkey)
	return nil
}

func (c *TcpClient) processIncomigMessage(pubkey *tools.Pubkey, data []byte) {
	var message ServerMessage
	var b []byte

	pubkey.DecPublicKey(data, &b)

	if err := json.Unmarshal(b, &message); err == nil {
		fmt.Errorf("ошибка расшифровки")
		return
	}

	if message.Action == "" {
		fmt.Errorf("не задан тип сообщения")
		return
	}

	switch message.Action {
	case "register":
		log.Print("register message in server")
	case "response", "error":
		c.routeToPedingRequest(&message)

	default:
		log.Print("abiba message")
	}

}

func (c *TcpClient) routeToPedingRequest(message *ServerMessage) {
	if message.CorrId == "" {
		return
	}

	c.mu.RLock()
	responseChat, exitis := c.PendingReqs[message.CorrId]
	c.mu.Unlock()

	if exitis {
		responseChat <- message

		c.mu.Lock()
		delete(c.PendingReqs, message.CorrId)
		c.mu.Unlock()
	}
}

func (c *TcpClient) readerLoop(pubkey *tools.Pubkey) {
	reader := bufio.NewReader(c.Conn)

	for {
		select {
		case <-c.done:
			return
		default:
			// Читаем длину сообщения
			var lenBuf [4]byte
			if _, err := io.ReadFull(reader, lenBuf[:]); err != nil {
				if err != io.EOF {
					fmt.Printf("Ошибка чтения длины от %s: %v\n", c.Conn.RemoteAddr(), err)
				}
				c.Conn.Close()
				return
			}

			msgLen := binary.BigEndian.Uint32(lenBuf[:])

			// Проверяем максимальный размер сообщения
			if msgLen > 10*1024*1024*100 { // 1gb максимум
				fmt.Printf("Сообщение слишком большое: %d bytes\n", msgLen)
				c.Conn.Close()
				return
			}

			// Читаем само сообщение
			msgBuf := make([]byte, msgLen)
			if _, err := io.ReadFull(reader, msgBuf); err != nil {
				fmt.Printf("Ошибка чтения данных от %s: %v\n", c.Conn.RemoteAddr(), err)
				c.Conn.Close()
				return
			}

			go c.processIncomigMessage(pubkey, msgBuf)
		}
	}
}
