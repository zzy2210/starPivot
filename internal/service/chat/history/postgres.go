package history

import (
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
	err := h.db.Where("username = ? AND chat_id = ?", username, chatID).Find(&chatHistory).Error
	if err != nil {
		return nil, err
	}
	return chatHistory, nil
}

func (h *PostgresChatHistory) AddChatHistory(username string, chatID string, message *schema.Message) error {
	return h.db.Create(message).Error
}

func (h *PostgresChatHistory) DeleteChatHistory(username string, chatID string) error {
	return h.db.Where("username = ? AND chat_id = ?", username, chatID).Delete(&schema.Message{}).Error
}

func (h *PostgresChatHistory) ListChatIDByUsername(username string) ([]string, error) {
	chatIDs := make([]string, 0)
	err := h.db.Model(&schema.Message{}).Where("username = ?", username).Distinct().Pluck("chat_id", &chatIDs).Error
	if err != nil {
		return nil, err
	}
	return chatIDs, nil
}

func (h *PostgresChatHistory) ListChatHistoryByUsername(username string) ([]*schema.Message, error) {
	chatHistory := make([]*schema.Message, 0)
	err := h.db.Where("username = ?", username).Find(&chatHistory).Error
	if err != nil {
		return nil, err
	}
	return chatHistory, nil
}
