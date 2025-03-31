package server

import (
	"context"
	"net/http"
	"starPivot/internal/model"
	"starPivot/internal/service/chat"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// handleChat 处理聊天请求
func (s *Server) handleChat(c echo.Context) error {

	var req model.Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request format: " + err.Error(),
		})
	}

	username := c.Request().Header.Get("X-Username")

	chatID := req.ChatID
	if chatID == "" {
		chatID = uuid.New().String()
	}

	if req.Messages == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "messages cannot be empty",
		})
	}

	chatHistory, err := s.ChatHistory.GetChatHistory(username, chatID)
	if err != nil && err != model.ErrChatHistoryNotFound {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to get chat history: " + err.Error(),
		})
	}

	// 创建聊天消息并保存到历史记录中
	chatMsg := chat.CreateMessageFromTemplate(req.Messages, chatHistory)

	result := chat.Generate(context.Background(), s.chatModel, chatMsg)

	return c.JSON(http.StatusOK, model.Response{
		Messages: result.String(),
		ChatID:   chatID,
	})
}

func (s *Server) handleListChatIDs(c echo.Context) error {
	username := c.Request().Header.Get("X-Username")

	ids, err := s.ChatHistory.ListChatIDByUsername(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to get chat history: " + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, ids)
}

func (s *Server) handleDeleteChat(c echo.Context) error {
	username := c.Request().Header.Get("X-Username")
	chatID := c.Param("chatID")

	err := s.ChatHistory.DeleteChatHistory(username, chatID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to delete chat history: " + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "chat history deleted successfully",
	})
}

func (s *Server) handleGetChatHistory(c echo.Context) error {
	username := c.Request().Header.Get("X-Username")
	chatID := c.Param("chatID")

	history, err := s.ChatHistory.GetChatHistory(username, chatID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to get chat history: " + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, history)
}
