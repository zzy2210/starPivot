package data

type ChatHistory struct {
	ID          string `json:"id" gorm:"primary_key column:id"`
	Username    string `json:"username" gorm:"column:username"`
	ChatMessage string `json:"chat_message" gorm:"column:chat_message"` // 聊天消息序列化
}
