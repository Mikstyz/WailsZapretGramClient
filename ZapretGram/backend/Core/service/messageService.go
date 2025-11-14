package service

import (
	"ZapretGram/backend/Core/db/repo"
	model "ZapretGram/backend/Core/ethernet/Model"
	"context"
	"database/sql"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type MessageService struct {
	chatid    int64
	msgBuffer []model.MessageInChat
	repo      *repo.MessagesRepo
	ctx       context.Context
}

func NewMessageService(db *sql.DB, chatid int64) *MessageService {
	return &MessageService{
		msgBuffer: make([]model.MessageInChat, 0, 100), // емкость 100
		repo:      repo.NewMessagesRepo(db),
		chatid:    chatid,
	}
}

func (s *MessageService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// GetBuffer возвращает текущий буфер сообщений
func (s *MessageService) GetBuffer() []model.MessageInChat {
	if len(s.msgBuffer) == 0 {
		return []model.MessageInChat{}
	}
	return s.msgBuffer
}

// ViewBuffer выводит полный буфер сообщений в консоль
func (s *MessageService) ViewBuffer() {
	fmt.Println("=== ПОЛНЫЙ БУФЕР СООБЩЕНИЙ ===")
	fmt.Printf("Размер буфера: %d/%d\n", len(s.msgBuffer), cap(s.msgBuffer))
	fmt.Println("------------------------------")

	if len(s.msgBuffer) == 0 {
		fmt.Println("Буфер пуст")
		return
	}

	for i, msg := range s.msgBuffer {
		fmt.Printf("[%d] ID: %d, ChatID: %d, Content: %s\n",
			i, msg.Id, msg.ChatId, msg.Message)
	}

	fmt.Println("==============================")
}

// emitBufferUpdate отправляет событие с обновленным буфером
func (s *MessageService) emitBufferUpdate() {
	if s.ctx != nil {
		runtime.EventsEmit(s.ctx, "bufferUpdated", s.msgBuffer)
	}
}

func (s *MessageService) inBuffer(msg model.MessageInChat) error {
	if len(s.msgBuffer) >= 100 {
		s.msgBuffer = s.msgBuffer[1:] // удаляем самый старый элемент
	}
	s.msgBuffer = append(s.msgBuffer, msg) // добавляем новый в конец
	return nil
}

func (s *MessageService) AddMessage(msg model.MessageInChat) error {
	// id, err := s.repo.AddMessage(msg)
	// if err != nil {
	// 	return fmt.Errorf("ошибка при сохранении сообщения в бд: %w", err)
	// }

	err := s.inBuffer(msg)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении сообщения в буфер: %w", err)
	}

	// Отправляем событие об обновлении буфера
	s.emitBufferUpdate()

	fmt.Printf("Добавлено новое сообщение - lenbuff:%d \n", len(s.msgBuffer))
	return nil
}

func (s *MessageService) LoadMessages(chatid int64, lastMessageID int) error {
	msgs, err := s.repo.GetMessages(chatid, lastMessageID, 25)
	if err != nil {
		return err
	}

	// Добавляем загруженные сообщения в начало буфера
	s.msgBuffer = append(msgs, s.msgBuffer...)

	// Если буфер превысил максимальный размер (100), обрезаем с конца
	if len(s.msgBuffer) > 100 {
		s.msgBuffer = s.msgBuffer[:100]
	}

	// Отправляем событие об обновлении буфера
	s.emitBufferUpdate()

	return nil
}

func (s *MessageService) LoadLastMessageInChat(chatId int64) error {
	// Загружаем последние сообщения из чата
	msgs, err := s.repo.GetMessages(chatId, 0, 25) // Предполагаем, что такой метод есть
	if err != nil {
		return err
	}

	// Заменяем текущий буфер загруженными сообщениями
	s.msgBuffer = msgs

	// Отправляем событие об обновлении буфера
	s.emitBufferUpdate()

	return nil
}

// ClearBuffer очищает буфер и уведомляет фронт
func (s *MessageService) ClearBuffer() {
	s.msgBuffer = make([]model.MessageInChat, 0, 100)
	s.emitBufferUpdate()
}

// GetBufferSize возвращает текущий размер буфера
func (s *MessageService) GetBufferSize() int {
	return len(s.msgBuffer)
}
