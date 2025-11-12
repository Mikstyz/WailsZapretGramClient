package ethernet

import (
	"ZapretGram/backend/Core/Tools"
	model "ZapretGram/backend/Core/ethernet/Model"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"time"
)

type TcpRequest struct {
	Tcp *TcpClient
}

func NewRequest(tcp *TcpClient, key string) *TcpRequest {
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

func (tcp *TcpClient) sendRequest(b []byte) {
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

func (tcp *TcpClient) readingAnsfer() *model.ResponseTcp {
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

	var res model.ResponseTcp
	if err := tcp.Key.DecPublicKey(respBuf, &res); err != nil {
		fmt.Printf("Ошибка расшифровки ответа: %v\n", err)
		return nil
	}

	return &res
}

func (t *TcpRequest) Auth(log string, pass string, action string) error {
	if action != "register" && action != "login" {
		fmt.Println("неверный тип запроса only register of login")
		return fmt.Errorf("неверный тип запроса only register of login\n")
	}

	//модель юзера
	user := model.RequestTcp{
		Action:   action,
		DateTime: time.Now().Format(time.RFC3339),
		Data: model.RequestAuthData{
			UserIn:     log,
			PasswordIn: pass,
		},
	}

	//Шифруем сообщение

	b, err := t.Tcp.Key.EncPublicKey(user)

	fmt.Printf("b: %d", b)

	if err != nil {
		fmt.Errorf("Ошибка шифрования: %v\n", err)
	}

	//Отправляем длину сообщения
	t.Tcp.sendRequest(b)

	//читаем ответ
	res := t.Tcp.readingAnsfer()

	fmt.Printf("res server: %+v\n", res)

	return nil
}

func (t *TcpClient) Ping() bool {
	ping := model.RequestTcp{
		Action: "ping",
	}

	b, err := t.Key.EncPublicKey(ping)
	if err != nil {
		fmt.Errorf("Ошибка шифрования: %v\n", err)
	}
	t.sendRequest(b)

	res := t.readingAnsfer()

	fmt.Printf("res server: %+v\n", res)
	s, _ := strconv.ParseBool(res.Status)
	return s
}
