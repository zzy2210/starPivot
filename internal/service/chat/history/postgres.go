package history

import (
	"starPivot/internal/data"
	"starPivot/internal/model"

	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 使用gorm

type PostgresChatHistory struct {
	db *gorm.DB
}

func NewDatabaseChatHistory(db *gorm.DB) *PostgresChatHistory {
	return &PostgresChatHistory{db: db}
}

func (h *PostgresChatHistory) GetChatHistory(username string, chatID string) ([]*schema.Message, error) {
	chatHistory := make([]*schema.Message, 0)
	dialogue := data.Dialogue{}
	err := h.db.Where("username = ? AND id = ?", username, chatID).First(&dialogue).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, model.ErrChatHistoryNotFound
	}

	messages := []data.Message{}
	err = h.db.Where("dialogue_id = ?", dialogue.ID).Order("seq_num ASC").Find(&messages).Error
	if err != nil {
		return nil, err
	}

	for _, message := range messages {
		chatHistory = append(chatHistory, &schema.Message{
			Role:    schema.RoleType(message.Role),
			Content: message.Content,
		})
	}

	return chatHistory, nil
}

func (h *PostgresChatHistory) AddChatHistory(username string, chatID string, message *schema.Message) error {
	// 如果 dialogue 不存在则创建
	dialogue := data.Dialogue{}
	err := h.db.Where("username = ? AND id = ?", username, chatID).First(&dialogue).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		dialogue.ID = chatID
		dialogue.Username = username
		err = h.db.Create(&dialogue).Error
		if err != nil {
			return err
		}
	}

	// 查询同一对话的最新message的seq_num
	var latestMessage data.Message
	err = h.db.Where("dialogue_id = ?", dialogue.ID).Order("seq_num DESC").First(&latestMessage).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		latestMessage.SeqNum = 0
	}

	msg := data.Message{
		ID:         uuid.New().String(),
		DialogueID: dialogue.ID,
		Role:       string(message.Role),
		Content:    message.Content,
		SeqNum:     latestMessage.SeqNum + 1,
	}
	return h.db.Create(&msg).Error
}

func (h *PostgresChatHistory) DeleteChatHistory(username string, chatID string) error {
	// 删除 dialogue 和 message
	if err := h.db.Where("username = ? AND id = ?", username, chatID).Delete(&data.Dialogue{}).Error; err != nil {
		return err
	}
	if err := h.db.Where("dialogue_id = ?", chatID).Delete(&data.Message{}).Error; err != nil {
		return err
	}
	return nil
}

func (h *PostgresChatHistory) ListChatIDByUsername(username string) ([]string, error) {
	chatIDs := make([]string, 0)
	err := h.db.Model(&data.Dialogue{}).Where("username = ?", username).Distinct().Pluck("id", &chatIDs).Error
	if err != nil {
		return nil, err
	}
	return chatIDs, nil
}
