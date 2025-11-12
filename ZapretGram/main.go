package main

import (
	"ZapretGram/backend/conf"
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	_ "modernc.org/sqlite"
)

//go:embed:../frontend/dist/*
var assets embed.FS

func main() {
	// Загружаем конфиг
	cfg, err := conf.LoadConfig()
	if err != nil {
		log.Fatalf("[main] Ошибка загрузки конфига: %v", err)
	}
	defer cfg.DBConn.Close()

	log.Printf("[main] Конфиг загружен успешно")

	// Создаем экземпляр приложения
	app := NewApp(cfg) // Передаем конфиг в приложение

	// Запускаем Wails приложение
	err = wails.Run(&options.App{
		Title:  "ZapretGram",
		Width:  1200,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{
			R: 27,
			G: 38,
			B: 54,
			A: 1,
		},
		OnStartup:  app.Startup,
		OnShutdown: app.Shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatalf("[main] Ошибка запуска приложения: %v", err)
	}

	log.Printf("[main] Приложение завершило работу")
}
