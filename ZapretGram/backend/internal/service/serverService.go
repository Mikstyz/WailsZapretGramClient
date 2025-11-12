package service

import (
	tools "ZapretGram/backend/internal/Tools"
	model "ZapretGram/backend/internal/service/Model"
	"fmt"
	"log"
)

// ServiceStorage управляет сохранением и загрузкой списка серверов.
type ServiceStorage struct {
	storage *model.StorageServer
}

// NewStorage создаёт новый объект ServiceStorage и сразу пытается загрузить данные из файла.
func NewStorage() *ServiceStorage {
	s := &ServiceStorage{}

	loaded, err := tools.LoadFromStorage()
	if err != nil {
		log.Printf("[Storage] Ошибка загрузки: %v", err)
		s.storage = &model.StorageServer{
			Servers: make(map[string]model.Server),
		}
	} else {
		s.storage = &loaded
	}

	// Если структура пуста — инициализируем карту.
	if s.storage.Servers == nil {
		s.storage.Servers = make(map[string]model.Server)
	}

	return s
}

// saveServer сохраняет текущее состояние хранилища в файл.
func (s *ServiceStorage) saveServer() error {
	if s.storage == nil {
		return fmt.Errorf("storage не инициализирован")
	}

	if err := tools.SaveToStorage(*s.storage); err != nil {
		return fmt.Errorf("ошибка сохранения: %v", err)
	}

	return nil
}

// AddServer добавляет или обновляет сервер в хранилище.
func (s *ServiceStorage) AddServer(name, address, publicKey string) error {
	if s == nil || s.storage == nil {
		return fmt.Errorf("storage не инициализирован")
	}

	s.storage.Servers[name] = model.Server{
		Address:   address,
		PublicKey: publicKey,
	}

	if err := s.saveServer(); err != nil {
		return fmt.Errorf("ошибка при добавлении сервера: %v", err)
	}

	return nil
}

// RemoveServer удаляет сервер из хранилища по имени.
func (s *ServiceStorage) RemoveServer(name string) error {
	if s == nil || s.storage == nil {
		return fmt.Errorf("storage не инициализирован")
	}

	delete(s.storage.Servers, name)

	if err := s.saveServer(); err != nil {
		return fmt.Errorf("ошибка при удалении сервера: %v", err)
	}

	return nil
}

// GetServer возвращает сервер по имени.
func (s *ServiceStorage) GetServer(name string) (model.Server, bool) {
	if s == nil || s.storage == nil {
		return model.Server{}, false
	}

	server, exists := s.storage.Servers[name]
	return server, exists
}

// GetAllServers возвращает карту всех серверов.
func (s *ServiceStorage) GetAllServers() map[string]model.Server {
	if s == nil || s.storage == nil {
		return make(map[string]model.Server)
	}
	return s.storage.Servers
}
