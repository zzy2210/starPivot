package chat

import (
	"context"
	"log"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func CreateTemplate() prompt.ChatTemplate {
	template := prompt.FromMessages(schema.FString,
		// 系统消息目标
		&schema.Message{
			Role:    schema.System,
			Content: "你是一个{role}。你需要用{style}的语气回答问题",
		},
		// 插入历史对话
		schema.MessagesPlaceholder("chat_history", true),

		// 用户消息模板
		&schema.Message{
			Role:    schema.User,
			Content: "{task}",
		},
	)
	return template
}

// CreateMessageFromTemplate 根据消息内容和聊天历史创建消息
func CreateMessageFromTemplate(msg string, chatHistory []*schema.Message) []*schema.Message {
	template := CreateTemplate()

	// 确保chatHistory不为nil
	if chatHistory == nil {
		chatHistory = make([]*schema.Message, 0)
	}

	messages, err := template.Format(context.Background(), map[string]any{
		"role":         "用户助理",
		"style":        "积极乐观",
		"task":         msg,
		"chat_history": chatHistory,
	})
	if err != nil {
		log.Fatalf("从模板创建消息失败: err= %v", err)
	}
	return messages
}
