package history

import (
	"encoding/json"
	"starPivot/internal/data"

	"github.com/cloudwego/eino/schema"
	"gorm.io/gorm"
)

// 使用gorm

type PostgresChatHistory struct {
	db *gorm.DB
}

func NewPostgresChatHistory(db *gorm.DB) *PostgresChatHistory {
	return &PostgresChatHistory{db: db}
}

func (h *PostgresChatHistory) GetChatHistory(username string, chatID string) ([]*schema.Message, error) {
	chatHistory := make([]*schema.Message, 0)

	dbHistory := data.ChatHistory{}
	err := h.db.Where("username = ? AND id = ?", username, chatID).Find(&dbHistory).Error
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(dbHistory.ChatMessage), &chatHistory)
	if err != nil {
		return nil, err
	}

	return chatHistory, nil
}

func (h *PostgresChatHistory) AddChatHistory(username string, chatID string, message *schema.Message) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	dbHistory := data.ChatHistory{
		ID:          chatID,
		Username:    username,
		ChatMessage: string(messageBytes),
	}
	return h.db.Create(&dbHistory).Error
}

func (h *PostgresChatHistory) DeleteChatHistory(username string, chatID string) error {
	return h.db.Where("username = ? AND id = ?", username, chatID).Delete(&data.ChatHistory{}).Error
}

func (h *PostgresChatHistory) ListChatIDByUsername(username string) ([]string, error) {
	chatIDs := make([]string, 0)
	err := h.db.Model(&data.ChatHistory{}).Where("username = ?", username).Distinct().Pluck("id", &chatIDs).Error
	if err != nil {
		return nil, err
	}
	return chatIDs, nil
}

func (h *PostgresChatHistory) ListChatHistoryByUsername(username string) ([]*schema.Message, error) {
	chatHistory := make([]*schema.Message, 0)
	dbHistory := make([]data.ChatHistory, 0)
	err := h.db.Where("username = ?", username).Find(&dbHistory).Error
	if err != nil {
		return nil, err
	}
	for _, history := range dbHistory {
		chat := schema.Message{}
		err = json.Unmarshal([]byte(history.ChatMessage), &chat)
		if err != nil {
			return nil, err
		}
		chatHistory = append(chatHistory, &chat)
	}
	return chatHistory, nil
}
