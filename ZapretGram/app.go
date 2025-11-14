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
	msgService     *service.MessageService
	ServersStorage *service.ServiceStorage
	DBConn         *sql.DB
	// Добавь другие зависимости (TCP клиент, база данных и т.д.)
}

func NewApp(cfg *conf.Config) *App {
	return &App{
		cfg:        cfg,
		tcp:        nil,
		msgService: nil,
		DBConn:     cfg.DBConn,
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	a.msgService.SetContext(a.ctx)

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

<<<<<<< HEAD
func (a *App) Auth(log string, pass string, action string) (map[string]model.Chat, error) {
=======
func (a *App) Auth(log string, pass string, action string) map[string]model.Chat {
>>>>>>> b45bacd57b688e7d2ab741613d97cbf8fb6d7ea2
	fmt.Print("auth in aboba")
	fmt.Printf(log, pass, action)
	chats := a.tcp.Auth(log, pass, action)

<<<<<<< HEAD
	if err != nil {
		return nil, err
	}

	// Return the chats map after successful auth
	m := a.tcp.Tcp.Chats
	if m == nil {
		return make(map[string]model.Chat), nil
	}

	return m, nil
}

func (a *App) NewChat(recipient string) (map[string]model.Chat, error) {
	err := a.tcp.NewChat(recipient)

	if err != nil {
		return nil, err
	}

	// Return the updated chats map
	m := a.tcp.Tcp.Chats
	if m == nil {
		return make(map[string]model.Chat), nil
	}
=======
	return chats
}

func (a *App) NewChat(recipient string) map[string]model.Chat {
	datachat := a.tcp.NewChat(recipient)

	if datachat == nil {
		return map[string]model.Chat{}
	}

	return datachat
>>>>>>> b45bacd57b688e7d2ab741613d97cbf8fb6d7ea2

	return m, nil
}

func (a *App) OpenChat(chatid int64) error {
	a.msgService = service.NewMessageService(a.DBConn, chatid)

	if a.msgService == nil {
		return fmt.Errorf("не удалось зарегестрировать буффер сообщений в чате")
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
