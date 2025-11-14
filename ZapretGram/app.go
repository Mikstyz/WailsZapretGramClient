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

type App struct {
	ctx            context.Context
	cfg            *conf.Config
	tcp            *ethernet.TcpRequest
	msgService     *service.MessageService
	ServersStorage *service.ServiceStorage
	DBConn         *sql.DB
}

func NewApp(cfg *conf.Config) *App {
	return &App{
		cfg:        cfg,
		tcp:        nil,
		msgService: nil, // Явно инициализируем как nil
		DBConn:     cfg.DBConn,
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// УБЕРИТЕ эту строку - msgService еще не инициализирован!
	// a.msgService.SetContext(a.ctx)

	log.Println("[App] Startup completed")

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

func (a *App) ConnectServer(ip string, port string, Pubkey string) error {
	err := Tools.Ping(ip, port)
	if err != nil {
		return err
	}

	client, err := ethernet.NewTcpClient(ip, port, a.DBConn)
	if err != nil {
		return fmt.Errorf("ошибка создания TCP клиента: %v", err)
	}

	tcp := ethernet.NewRequest(client, Pubkey)
	if tcp == nil {
		return fmt.Errorf("tcp is nil")
	}

	a.tcp = tcp
	return nil
}

func (a *App) Auth(log string, pass string, action string) map[string]model.Chat {
	fmt.Print("auth in aboba")
	fmt.Printf(log, pass, action)

	if a.tcp == nil {
		fmt.Println("TCP клиент не инициализирован")
		return map[string]model.Chat{}
	}

	chats := a.tcp.Auth(log, pass, action)
	return chats
}

func (a *App) NewChat(recipient string) map[string]model.Chat {
	if a.tcp == nil {
		fmt.Println("TCP клиент не инициализирован")
		return map[string]model.Chat{}
	}

	datachat := a.tcp.NewChat(recipient)
	if datachat == nil {
		return map[string]model.Chat{}
	}
	return datachat
}

func (a *App) OpenChat(chatid int64) error {
	a.msgService = service.NewMessageService(a.DBConn, chatid)
	if a.msgService == nil {
		return fmt.Errorf("не удалось зарегистрировать буфер сообщений в чате")
	}

	// ТЕПЕРЬ устанавливаем контекст, когда msgService инициализирован
	a.msgService.SetContext(a.ctx)

	return nil
}

func (a *App) NewMessage(ChatId int64, message string) error {
	if a.tcp == nil {
		return fmt.Errorf("TCP клиент не инициализирован")
	}

	err := a.tcp.NewMessage(ChatId, message)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) UpdateChat(chatId int64) map[string]model.Chat {
	if a.tcp == nil || a.tcp.Tcp == nil {
		return make(map[string]model.Chat)
	}

	m := a.tcp.Tcp.Chats
	if m == nil {
		return make(map[string]model.Chat)
	}
	return m
}

// Дополнительные методы для работы с сообщениями
func (a *App) GetMessageBuffer() []model.MessageInChat {
	if a.msgService == nil {
		return []model.MessageInChat{}
	}
	return a.msgService.GetBuffer()
}

func (a *App) ClearMessageBuffer() {
	if a.msgService != nil {
		a.msgService.ClearBuffer()
	}
}
