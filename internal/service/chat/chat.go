package chat

import (
	"context"
	"log"
	"starPivot/internal/model"

	"github.com/cloudwego/eino-ext/components/model/openai"
	einoModel "github.com/cloudwego/eino/components/model"
)

func CreateOPenAIChatModel(ctx context.Context, config *model.BaseAIConfig) einoModel.ChatModel {

	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: config.BaseURL,
		Model:   config.Model,
		APIKey:  config.APIKey,
	})
	if err != nil {
		log.Fatalf("create failed: err= %v", err)
	}
	return chatModel
}
