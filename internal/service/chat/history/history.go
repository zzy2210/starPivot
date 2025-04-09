package history

import (
	"github.com/cloudwego/eino/schema"
)

type ChatHistoryFactory interface {
	GetChatHistory(username string, chatID string) ([]*schema.Message, error)
	AddChatHistory(username string, chatID string, message *schema.Message) error
	DeleteChatHistory(username string, chatID string) error
	ListChatIDByUsername(username string) ([]string, error)
}
