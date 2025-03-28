package chat

import (
	"context"
	"log"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func CreateTemplate(msg string) prompt.ChatTemplate {
	template := prompt.FromMessages(schema.FString,
		// 系统消息目标
		schema.SystemMessage("你是一个{role}。你需要用{style}的语气回答问题,你的目标是帮助用户保持积极乐观的心态，提供用户建议的时候需要关注心理健康"),
		// 插入历史对话
		schema.MessagesPlaceholder("chat_history", true),

		// 用户消息模板
		schema.UserMessage(msg),
	)
	return template
}

// CreateDefaultMessages 创建默认的消息（无参数版本）
func CreateDefaultMessages() []*schema.Message {
	template := CreateTemplate("你好，我需要一些鼓励")
	messages, err := template.Format(context.Background(), map[string]any{
		"role":         "程序员鼓励师",
		"style":        "积极乐观",
		"chat_history": make([]*schema.Message, 0),
	})
	if err != nil {
		log.Fatalf("从模板创建消息失败: err= %v", err)
	}
	return messages
}

// CreateMessageFromTemplate 根据消息内容和聊天历史创建消息
func CreateMessageFromTemplate(msg string, chatHistory []*schema.Message) []*schema.Message {
	template := CreateTemplate(msg)

	// 确保chatHistory不为nil
	if chatHistory == nil {
		chatHistory = make([]*schema.Message, 0)
	}

	messages, err := template.Format(context.Background(), map[string]any{
		"role":         "",
		"style":        "",
		"chat_history": chatHistory,
	})
	if err != nil {
		log.Fatalf("从模板创建消息失败: err= %v", err)
	}
	return messages
}
