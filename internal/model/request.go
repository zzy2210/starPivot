package model

import "github.com/cloudwego/eino/schema"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type SendChatRequest struct {
	Messages string `json:"Messages"`
	ChatID   string `json:"ChatID"`
}

type SendChatResponse struct {
	Messages string `json:"Messages"`
	ChatID   string `json:"ChatID"`
}

type ListChatIDsRequest struct {
	Username string `json:"Username"`
}

type ListChatIDsResponse struct {
	ChatIDs []string `json:"ChatIDs"`
}

type DeleteChatRequest struct {
	Username string `json:"Username"`
	ChatID   string `json:"ChatID"`
}

type DeleteChatResponse struct {
	Message string `json:"Message"`
}

type GetChatHistoryRequest struct {
	Username string `json:"Username"`
	ChatID   string `json:"ChatID"`
}

// TODO Y1nui: 直接返回eino 的 message 并不合理
// 需要返回一个更友好的格式“
type GetChatHistoryResponse struct {
	Messages []*schema.Message `json:"Messages"`
}

type NewChatRequest struct {
	Username string `json:"Username"`
}

type NewChatResponse struct {
	ChatID string `json:"ChatID"`
}

type ConfigRequest struct {
	Model   string `json:"Model"`
	BaseURL string `json:"BaseURL"`
	APIKey  string `json:"APIKey"`
	// 服务提供商,当前只支持 openai，后续考虑 ds，ollama 等
	ModelType string `json:"ModelType"`
}
