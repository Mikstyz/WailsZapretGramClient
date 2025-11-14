package main

import (
	"ZapretGram/backend/Core/Tools"
	"ZapretGram/backend/Core/ethernet"
	model "ZapretGram/backend/Core/ethernet/Model"
	"ZapretGram/backend/Core/service"
	"ZapretGram/backend/conf"
	"context"
	"database/sql"
	"fmt"
	"log"
)

// app.go
type App struct {
	ctx context.Context
	cfg *conf.Config

	tcp            *ethernet.TcpRequest
	ServersStorage *service.ServiceStorage
	DBConn         *sql.DB
	// Добавь другие зависимости (TCP клиент, база данных и т.д.)
}

func NewApp(cfg *conf.Config) *App {
	return &App{
		cfg:    cfg,
		tcp:    nil,
		DBConn: cfg.DBConn,
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

func (a *App) ConnectServer(ip string, port string, Pubkey string) error {
	err := Tools.Ping(ip, port)
	if err != nil {
		return err
	}

	client, err := ethernet.NewTcpClient(ip, port, a.DBConn)

	tcp := ethernet.NewRequest(client, Pubkey)

	if tcp == nil {
		return fmt.Errorf("tcp is nil")
	}

	a.tcp = tcp
	return nil
}

func (a *App) Auth(log string, pass string, action string) error {
	fmt.Print("auth in aboba")
	fmt.Printf(log, pass, action)
	err := a.tcp.Auth(log, pass, action)

	if err != nil {
		return err
	}

	return nil
}

func (a *App) NewChat(recipient string) error {
	err := a.tcp.NewChat(recipient)

	if err != nil {
		return err
	}

	return nil

}

func (a *App) NewMessage(ChatId int64, message string) error {
	err := a.tcp.NewMessage(ChatId, message)

	if err != nil {
		return err
	}

	return nil
}

func (a *App) UpdateChat(chatId int64) map[string]model.Chat {

	m := a.tcp.Tcp.Chats

	if m == nil {
		return make(map[string]model.Chat)
	}

	return m
}
