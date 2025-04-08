package server

import (
	"context"
	"net/http"
	"starPivot/internal/model"
	"starPivot/internal/service/chat"

	"github.com/labstack/echo/v4"
)

func (s *Server) handleConfig(c echo.Context) error {
	var req model.ConfigRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid request format: " + err.Error(),
			Data:    nil,
		})
	}

	switch req.ModelType {
	case "openai":
		s.ChatConfig = &model.BaseAIConfig{
			BaseURL: req.BaseURL,
			Model:   req.Model,
			APIKey:  req.APIKey,
		}
	default:
		return c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid model type: " + req.ModelType,
			Data:    nil,
		})
	}

	s.chatModel = chat.CreateOPenAIChatModel(context.Background(), s.ChatConfig)

	return c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "config updated",
		Data:    nil,
	})
}
