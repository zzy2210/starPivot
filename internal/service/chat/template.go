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
		schema.SystemMessage("你是一个{role}。你需要用{style}的语气回答问题,你的目标是帮助用户保持积极乐观的心态，提供用户建议的时候需要关注心理健康"),
		// 插入历史对话
		schema.MessagesPlaceholder("chat_history", true),

		// 用户消息模板
		schema.UserMessage("你好，我是{role}，我有一个问题想请教你，你能帮我解答一下吗？"),
	)
	return template
}

func CreateMessageFromTemplate() []*schema.Message {
	template := CreateTemplate()
	messages, err := template.Format(context.Background(), map[string]any{
		"role":      "程序员鼓励师",
		"style":     "积极乐观",
		"questions": "我写的代码依托，用了ai后又看不懂怎么办",
		//手动构造历史对话
		"chat_history": []*schema.Message{
			schema.UserMessage("你好，我的朋友"),
			schema.AssistantMessage("嘿！我是你的程序员鼓励师！记住，每个优秀的程序员都是从 Debug 中成长起来的。有什么我可以帮你的吗？", nil),
		},
	})
	if err != nil {
		log.Fatalf("create msg from template failed: err= %v", err)
	}
	return messages
}
