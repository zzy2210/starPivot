package data

import "time"

type Message struct {
	ID         string    `json:"id" gorm:"primary_key column:id"`
	DialogueID string    `json:"dialogue_id" gorm:"column:dialogue_id"`
	Role       string    `json:"role" gorm:"column:role"`
	Content    string    `json:"content" gorm:"column:content"`
	SeqNum     int       `json:"seq_num" gorm:"column:seq_num"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
}

type Dialogue struct {
	ID          string    `json:"id" gorm:"primary_key column:id"`
	Username    string    `json:"username" gorm:"column:username"`
	CreateTime  time.Time `json:"create_time" gorm:"column:create_time"`
	LastUpdated time.Time `json:"last_updated" gorm:"column:last_updated"`
}
