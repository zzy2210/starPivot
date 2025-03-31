package server

import (
	"context"
	"net/http"
	"starPivot/internal/model"
	"starPivot/internal/service/chat"

	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// handleChat 处理聊天请求
func (s *Server) handleChat(c echo.Context) error {

	var req model.SendChatRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			Code:    400,
			Message: "invalid request format: " + err.Error(),
			Data:    nil,
		})
	}

	username := c.Request().Header.Get("X-Username")

	chatID := req.ChatID
	if chatID == "" {
		chatID = uuid.New().String()
	}

	if req.Messages == "" {
		return c.JSON(http.StatusBadRequest, model.Response{
			Code:    400,
			Message: "messages cannot be empty",
			Data:    nil,
		})
	}

	chatHistory, err := s.ChatHistory.GetChatHistory(username, chatID)
	if err != nil && err != model.ErrChatHistoryNotFound {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Code:    500,
			Message: "failed to get chat history: " + err.Error(),
			Data:    nil,
		})
	}

	// 创建聊天消息并保存到历史记录中
	chatMsg := chat.CreateMessageFromTemplate(req.Messages, chatHistory)

	result := chat.Generate(context.Background(), s.chatModel, chatMsg)

	return c.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "success",
		Data: model.SendChatResponse{
			Messages: result.String(),
			ChatID:   chatID,
		},
	})
}

func (s *Server) handleNewChat(c echo.Context) error {
	username := c.Request().Header.Get("X-Username")
	chatID := uuid.New().String()

	err := s.ChatHistory.AddChatHistory(username, chatID, &schema.Message{
		Role:    "system",
		Content: "You are a helpful assistant.",
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Code:    500,
			Message: "failed to add chat history: " + err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "success",
		Data: model.NewChatResponse{
			ChatID: chatID,
		},
	})
}

func (s *Server) handleListChatIDs(c echo.Context) error {
	username := c.Request().Header.Get("X-Username")

	ids, err := s.ChatHistory.ListChatIDByUsername(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Code:    500,
			Message: "failed to get chat history: " + err.Error(),
			Data:    nil,
		})
	}
	return c.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "success",
		Data: model.ListChatIDsResponse{
			ChatIDs: ids,
		},
	})
}

func (s *Server) handleDeleteChat(c echo.Context) error {
	username := c.Request().Header.Get("X-Username")
	chatID := c.Param("chatID")

	err := s.ChatHistory.DeleteChatHistory(username, chatID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Code:    500,
			Message: "failed to delete chat history: " + err.Error(),
			Data:    nil,
		})
	}
	return c.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "success",
		Data: model.DeleteChatResponse{
			Message: "chat history deleted successfully",
		},
	})
}

func (s *Server) handleGetChatHistory(c echo.Context) error {
	username := c.Request().Header.Get("X-Username")
	chatID := c.Param("chatID")

	history, err := s.ChatHistory.GetChatHistory(username, chatID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Code:    500,
			Message: "failed to get chat history: " + err.Error(),
			Data:    nil,
		})
	}
	return c.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "success",
		Data: model.GetChatHistoryResponse{
			Messages: history,
		},
	})
}
