package conf

import (
	"ZapretGram/backend/internal/service"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	ClientDB       string
	ServersStorage *service.ServiceStorage
	DBConn         *sql.DB
}

// Загрузка env файла
func envInc() (string, error) {
	log.Print("[env] loading env file")

	pathenv := filepath.Join("backend", "cmd", "conf.env")
	if err := godotenv.Load(pathenv); err != nil {
		return "", fmt.Errorf("ошибка загрузки .env: %v", err)
	}

	dbPath := os.Getenv("client_db")
	if dbPath == "" {
		return "", fmt.Errorf("переменная client_db не указана в conf.env")
	}

	log.Printf("[env] env db path loaded successfully: %s", dbPath)
	return dbPath, nil
}

// Подключение к базе данных
func dbInc(dbPath string) (*sql.DB, error) {
	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка преобразования пути: %v", err)
	}

	// Создаём папку, если её нет
	if err := os.MkdirAll(filepath.Dir(absPath), os.ModePerm); err != nil {
		return nil, fmt.Errorf("не удалось создать папку для БД: %v", err)
	}

	log.Printf("[db] connecting to db: %s", absPath)

	conn, err := sql.Open("sqlite", absPath) // modernc.org/sqlite драйвер
	if err != nil {
		return nil, fmt.Errorf("[db] ошибка открытия соединения: %v", err)
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("[db] ping error: %v", err)
	}

	log.Print("[db] connection established successfully")
	return conn, nil
}

// Инициализация ServersStorage
func loadServersStorage() (*service.ServiceStorage, error) {
	st := service.NewStorage()
	if st == nil {
		return nil, fmt.Errorf("не удалось инициализировать storage")
	}
	return st, nil
}

// Основная загрузка конфига
func LoadConfig() (*Config, error) {
	dbPath, err := envInc()
	if err != nil {
		return nil, err
	}

	conn, err := dbInc(dbPath)
	if err != nil {
		return nil, err
	}

	st, err := loadServersStorage()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		ClientDB:       dbPath,
		ServersStorage: st,
		DBConn:         conn,
	}

	log.Print("[conf] конфигурация успешно загружена")
	return cfg, nil
}
