package main

import (
	"ZapretGram/backend/Core/ethernet"
	model "ZapretGram/backend/Core/ethernet/Model"
	"ZapretGram/backend/Core/service"
	"ZapretGram/backend/conf"
	"context"
	"database/sql"
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

func (a *App) GetUserInfo(userID int) map[string]interface{} {
	return map[string]interface{}{
		"id":   userID,
		"name": "Test User",
	}
}

func (a *App) ConnectServer(ip string, port string, Pubkey string) error {
	return nil
}

func (a *App) Auth(log string, pass string, action string) map[string]model.Chat {
	return map[string]model.Chat{}
}

func (a *App) NewChat(recipient string) map[string]model.Chat {
	return map[string]model.Chat{}
}

func (a *App) OpenChat(chatid int64) error {
	return nil
}

func (a *App) NewMessage(ChatId int64, message string) error {
	return nil
}

func (a *App) UpdateChat(chatId int64) map[string]model.Chat {
	return make(map[string]model.Chat)
}

// Дополнительные методы для работы с сообщениями
func (a *App) GetMessageBuffer() []model.MessageInChat {
	return a.msgService.GetBuffer()
}

func (a *App) ClearMessageBuffer() {
	if a.msgService != nil {
		a.msgService.ClearBuffer()
	}
}
