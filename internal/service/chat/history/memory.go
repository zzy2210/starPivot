package history

import (
	"fmt"
	"starPivot/internal/model"

	"github.com/cloudwego/eino/schema"
)

type MemoryChatHistory struct {
	chatHistoryWithID map[string][]*schema.Message
	UserWithChatID    map[string]string
}

func NewMemoryChatHistory() *MemoryChatHistory {
	return &MemoryChatHistory{
		chatHistoryWithID: make(map[string][]*schema.Message),
		UserWithChatID:    make(map[string]string),
	}
}

func (h *MemoryChatHistory) GetChatHistory(username string, chatID string) ([]*schema.Message, error) {
	chatHistory, ok := h.chatHistoryWithID[chatID]
	if !ok {
		return nil, model.ErrChatHistoryNotFound
	}
	return chatHistory, nil
}

func (h *MemoryChatHistory) AddChatHistory(username string, chatID string, message *schema.Message) error {
	h.chatHistoryWithID[chatID] = append(h.chatHistoryWithID[chatID], message)
	return nil
}

func (h *MemoryChatHistory) DeleteChatHistory(username string, chatID string) error {
	delete(h.chatHistoryWithID, chatID)
	return nil
}

func (h *MemoryChatHistory) ListChatIDByUsername(username string) ([]string, error) {
	chatIDs := make([]string, 0, len(h.chatHistoryWithID))
	for chatID := range h.chatHistoryWithID {
		chatIDs = append(chatIDs, chatID)
	}
	return chatIDs, nil
}

func (h *MemoryChatHistory) ListChatHistoryByUsername(username string) ([]*schema.Message, error) {
	chatHistory, ok := h.chatHistoryWithID[username]
	if !ok {
		return nil, fmt.Errorf("chat history not found")
	}
	return chatHistory, nil
}
