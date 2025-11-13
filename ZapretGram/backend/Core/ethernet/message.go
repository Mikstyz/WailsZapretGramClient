package ethernet

import (
	"ZapretGram/backend/Core/db/repo"
	model "ZapretGram/backend/Core/ethernet/Model"
)

func (c *Tcp) AddMessageInDb(msg model.MessageInChat) error {
	msgRepo := repo.NewMessagesRepo(c.DB)

	msgRepo.AddMessage(msg)
	// дальше работаешь с msgRepo
	return nil
}

func (c *Tcp) GetMessages(chatid int64, offset int) ([]model.MessageInChat, error) {
	msgRepo := repo.NewMessagesRepo(c.DB)

	msg, err := msgRepo.GetMessages(chatid, offset, 25)

	if err != nil {
		return nil, err
	}

	if msg == nil {
		return []model.MessageInChat{}, nil
	}

	return msg, nil
}
