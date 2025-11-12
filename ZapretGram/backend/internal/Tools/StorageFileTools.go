package Tools

import (
	model "ZapretGram/backend/internal/service/Model"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func storageFilePath() string {
	return filepath.Join("..", "internal", "Storage", "servers.json")
}

func SaveToStorage(storage model.StorageServer) error {
	path := storageFilePath()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("ошибка создания каталога: %v", err)
	}

	data, err := json.Marshal(storage)
	if err != nil {
		return fmt.Errorf("ошибка маршалинга: %v", err)
	}

	// Записываем с правами только для владельца
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("ошибка записи: %v", err)
	}

	return nil
}

func LoadFromStorage() (model.StorageServer, error) {
	var storage model.StorageServer
	path := storageFilePath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return storage, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return storage, fmt.Errorf("ошибка чтения: %v", err)
	}

	if err := json.Unmarshal(data, &storage); err != nil {
		return storage, fmt.Errorf("ошибка парсинга: %v", err)
	}

	return storage, nil
}
