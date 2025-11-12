package main

import (
	"ZapretGram/backend/Core/ethernet"
	"ZapretGram/backend/conf"
	"ZapretGram/backend/test"
	"time"

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
		var i int = 0

		client, _ := ethernet.NewTcpClient("26.69.104.210", "9000")

		tcp := ethernet.NewRequest(client, "wsdfvbndfghbjnmklrftghjkrtfghm348etvfghnj4567zsxdcfgvhbjjSDFGHRFGHSDFGVXDFGFGBHKJMLLTRFYGHUJK")

		for true {
			time.Sleep(500 * time.Millisecond)
			i++
			fmt.Printf("ping :%d \n", i)

			err := tcp.Auth("slut", "imSLUT", "register")

			if err != nil {
				fmt.Print("ошибка авторизации")
			}
		}
	}
}
