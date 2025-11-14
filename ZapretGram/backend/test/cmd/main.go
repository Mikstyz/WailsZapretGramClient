package main

import (
	"ZapretGram/backend/Core/ethernet"
	"ZapretGram/backend/Core/service"
	"ZapretGram/backend/conf"
	"ZapretGram/backend/test"

	"fmt"
	"log"
	"strings"

	_ "modernc.org/sqlite"
)

// main file
func main() {
	cfg, err := conf.LoadConfig()
	if err != nil {
		log.Fatalf("[main] ошибка загрузки конфига: %v", err)
	}

	defer cfg.DBConn.Close()
	log.Printf("[main] successful run")

	//err = cfg.ServersStorage.RemoveServer("aboba")
	//fmt.Printf("%w", err)

	// Пример запуска: либо слушаем, либо шлём

	var mode string
	fmt.Print("type: ")
	fmt.Scanln(&mode)
	mode = strings.TrimSpace(strings.ToLower(mode))
	if mode == "s" {
		test.OutServer()
	}

	if mode == "c" {
		ip := "26.69.104.210"
		port := "9000"
		Pubkey := "wsdfvbndfghbjnmklrftghjkrtfghm348etvfghnj4567zsxdcfgvhbjjSDFGHRFGHSDFGVXDFGFGBHKJMLLTRFYGHUJK"

		client, _ := ethernet.NewTcpClient(ip, port, cfg.DBConn)
		tcp := ethernet.NewRequest(client, Pubkey)

		// СОХРАНИТЕ MessageService в переменную и установите контекст
		messageService := service.NewMessageService(tcp.Tcp.DB, 2)
		messageService.SetContext(tcp.Tcp.Ctx) // ctx должен быть доступен в этой функции

		fmt.Printf("%s\n", tcp.Tcp.Key.MyKey())

		err := tcp.Auth("RyslanDayn3", "RyslanDayn3", "login")
		fmt.Printf("My user id: %d\n", tcp.Tcp.UserId)

		if err := client.RiderTcp(tcp.Tcp.Key); err != nil {
			log.Fatalf("Ошибка RiderTcp: %v", err)
		}

		fmt.Println("\nКлиент запущен, ожидание сообщений от сервера...")
		fmt.Println("отправили сообщение")
		tcp.NewMessage(2, "jopa")

		if err != nil {
			fmt.Print("ошибка авторизации\n")
		}

		select {}
	}
}
