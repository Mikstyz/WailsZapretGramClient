package repo

import (
	model "ZapretGram/backend/Core/ethernet/Model"
	"database/sql"
)

type MessagesRepo struct {
	db *sql.DB
}

func (r *MessagesRepo) NewMessagesRepo(db *sql.DB) *MessagesRepo {
	return &MessagesRepo{
		db: db,
	}
}

// offset = сколько сообщений пропустить
// limit = сколько вернуть (обычно 10)
func (r *MessagesRepo) GetMessages(chatID int64, offset int, limit int) ([]model.MessageInChat, error) {
	rows, err := r.db.Query(`
        SELECT id, user_id, chat_id, message
        FROM messages
        WHERE chat_id = ?
        ORDER BY id DESC
        LIMIT ? OFFSET ?
    `, chatID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.MessageInChat

	for rows.Next() {
		var m model.MessageInChat
		if err := rows.Scan(&m.Id, &m.UserId, &m.ChatId, &m.Message); err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	// разворачиваем, чтобы вернуть в нормальном порядке
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result, nil
}
