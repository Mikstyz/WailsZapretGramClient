package main

import (
	"ZapretGram/backend/Core/ethernet"
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

		// err := Tools.Ping(ip, port)
		// if err != nil {
		// 	log.Print("false")
		// }

		client, _ := ethernet.NewTcpClient(ip, port, cfg.DBConn)
		tcp := ethernet.NewRequest(client, Pubkey)
		fmt.Printf("%s\n", tcp.Tcp.Key.MyKey())
		// status := tcp.Tcp.Ping()

		// fmt.Print(tcp.Tcp.Key)
		// fmt.Print(status)

		//client, _ := ethernet.NewTcpClient("26.69.104.210", "9000")

		err := tcp.Auth("slut2", "imSLUT", "login")
		fmt.Printf("My user id: %d\n", tcp.Tcp.UserId)

		if err := client.RiderTcp(tcp.Tcp.Key); err != nil {
			log.Fatalf("Ошибка RiderTcp: %v", err)
		}

		fmt.Println("\nКлиент запущен, ожидание сообщений от сервера...")

		//tcp.NewChat("slut")

		fmt.Println("отправили сообщение")
		err = tcp.NewMessage(1, "ruslan huesos")

		if err != nil {
			fmt.Print("ошибка авторизации\n")
		}

		//tcp.NewChat("slut2")
		select {}
	}
}
