package main

import (
	"ZapretGram/backend/Core/Tools"
	"ZapretGram/backend/Core/ethernet"
	"ZapretGram/backend/conf"
	"context"
	"log"
)

// app.go
type App struct {
	ctx context.Context
	cfg *conf.Config

	// Добавь другие зависимости (TCP клиент, база данных и т.д.)
}

func NewApp(cfg *conf.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	log.Println("[App] Startup completed")

	// Здесь можно запустить TCP клиент и другие сервисы
	if err := a.startServices(); err != nil {
		log.Printf("[App] Ошибка запуска сервисов: %v", err)
	}
}

func (a *App) Shutdown(ctx context.Context) {
	log.Println("[App] Shutting down...")
	// Graceful shutdown всех сервисов
}

func (a *App) startServices() error {
	// Запуск TCP клиента и других сервисов
	return nil
}

func (a *App) GetMessage() string {
	return "hello"
}

func (a *App) Greet(name string) string {
	return "Hello " + name + "!"
}

func (a *App) GetUserInfo(userID int) map[string]interface{} {
	return map[string]interface{}{
		"id":   userID,
		"name": "Test User",
	}
}

//=======================================================================

func (a *App) ConnectServer(ip string, port string, Pubkey string) bool {
	err := Tools.Ping(ip, port)
	if err != nil {
		return false
	}

	client, err := ethernet.NewTcpClient(ip, port)

	tcp := ethernet.NewRequest(client, Pubkey)

	if tcp == nil {
		return false
	}

	return true
}
