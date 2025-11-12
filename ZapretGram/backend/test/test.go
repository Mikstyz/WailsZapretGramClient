package test

import (
	tool "ZapretGram/backend/Core/Tools"
	"ZapretGram/backend/Core/ethernet"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	_ "modernc.org/sqlite"

	modelE "ZapretGram/backend/Core/ethernet/Model"
)

func InServer() {
	key := tool.NewKey("wsdfvbndfghbjnmklrftghjkrtfghm348etvfghnj4567zsxdcfgvhbjjSDFGHRFGHSDFGVXDFGFGBHKJMLLTRFYGHUJK")

	tcp, err := ethernet.NewTcpClient("26.69.104.210", "9000")
	if err != nil {
		fmt.Printf("Ошибка подключения: %v\n", err)
		return
	}
	defer tcp.Disconnect()

	req := modelE.RequestTcp{
		Action:   "register",
		DateTime: time.Now().Format(time.RFC3339),
		Data: modelE.RequestAuthData{
			UserIn:     "aoba",
			PasswordIn: "pass",
		},
	}

	// Шифруем сообщение
	b, err := key.EncPublicKey(req)
	if err != nil {
		fmt.Printf("Ошибка шифрования: %v\n", err)
		return
	}

	// Отправляем сначала длину сообщения, потом сами данные
	length := uint32(len(b))
	buf := make([]byte, 4+len(b))
	binary.BigEndian.PutUint32(buf[:4], length)
	copy(buf[4:], b)

	if _, err := tcp.Conn.Write(buf); err != nil {
		fmt.Printf("Ошибка при отправке сообщения: %v\n", err)
		return
	}

	// Чтение ответа
	var lenBuf [4]byte
	if _, err := io.ReadFull(tcp.Conn, lenBuf[:]); err != nil {
		fmt.Printf("Ошибка чтения длины ответа: %v\n", err)
		return
	}

	respLen := binary.BigEndian.Uint32(lenBuf[:])
	respBuf := make([]byte, respLen)
	if _, err := io.ReadFull(tcp.Conn, respBuf); err != nil {
		fmt.Printf("Ошибка чтения ответа: %v\n", err)
		return
	}

	var res modelE.ResponseTcp
	if err := key.DecPublicKey(respBuf, &res); err != nil {
		fmt.Printf("Ошибка расшифровки ответа: %v\n", err)
		return
	}

	fmt.Printf("Ответ от сервера: %+v\n", res)
}

func OutServer() {
	addr := "26.162.220.63:9000"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
		return
	}
	defer ln.Close()
	fmt.Printf("Сервер запущен на %s\n", addr)

	key := tool.NewKey("wsdfvbndfghbjnmklrftghjkrtfghm348etvfghnj4567zsxdcfgvhbjjSDFGHRFGHSDFGVXDFGFGBHKJMLLTRFYGHUJK")
	clients := make(map[net.Conn]bool)
	var mu sync.Mutex

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Ошибка при подключении: %v\n", err)
			continue
		}

		mu.Lock()
		clients[conn] = true
		mu.Unlock()

		fmt.Printf("Новое соединение: %s\n", conn.RemoteAddr())

		go func(c net.Conn) {
			defer func() {
				c.Close()
				mu.Lock()
				delete(clients, c)
				mu.Unlock()
				fmt.Printf("Отключение: %s\n", c.RemoteAddr())
			}()

			for {
				time.Sleep(1 * time.Second)
				// Читаем длину сообщения
				var lenBuf [4]byte
				if _, err := io.ReadFull(c, lenBuf[:]); err != nil {
					fmt.Printf("Ошибка чтения длины от %s: %v\n", c.RemoteAddr(), err)
					return
				}
				msgLen := binary.BigEndian.Uint32(lenBuf[:])

				// Читаем само сообщение
				msgBuf := make([]byte, msgLen)
				if _, err := io.ReadFull(c, msgBuf); err != nil {
					fmt.Printf("Ошибка чтения данных от %s: %v\n", c.RemoteAddr(), err)
					return
				}

				var req modelE.RequestTcp
				if err := key.DecPublicKey(msgBuf, &req); err != nil {
					fmt.Printf("Ошибка расшифровки от %s: %v\n", c.RemoteAddr(), err)
					continue
				}

				fmt.Printf("Сообщение от %s: Action=%s\n", c.RemoteAddr(), req.Action)

				// Отправляем простой эхо-ответ
				resp := modelE.ResponseTcp{
					Action: "echo",
					Status: "ok",
				}

				respBytes, _ := key.EncPublicKey(resp)
				length := uint32(len(respBytes))
				buf := make([]byte, 4+len(respBytes))
				binary.BigEndian.PutUint32(buf[:4], length)
				copy(buf[4:], respBytes)

				if _, err := c.Write(buf); err != nil {
					fmt.Printf("Ошибка отправки ответа %s: %v\n", c.RemoteAddr(), err)
					return
				}
			}
		}(conn)
	}
}
