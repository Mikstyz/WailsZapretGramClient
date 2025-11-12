package main

import (
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
		var i int = 0

		for true {
			//time.Sleep(1 * time.Second)
			i++
			fmt.Printf("ping :%d \n", i)
			test.InServer()
		}
	}
}
