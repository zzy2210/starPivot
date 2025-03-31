package server

import (
	"context"
	"fmt"
	"net/http"

	"starPivot/internal/service/chat"

	eModel "github.com/cloudwego/eino/components/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

// Server 封装了 Echo 服务器和 ChatModel
type Server struct {
	echo        *echo.Echo
	chatModel   eModel.ChatModel
	logger      *logrus.Logger
	ChatHistory chat.ChatHistoryFactory
}

// Config 服务器配置选项
type Config struct {
	Port      string
	ChatModel eModel.ChatModel
	LogLevel  logrus.Level
}

// NewServer 创建一个新的服务器实例
func NewServer(ctx context.Context, cfg *Config) (*Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if cfg.ChatModel == nil {
		return nil, fmt.Errorf("chat model cannot be nil")
	}

	// 初始化日志
	logger := logrus.New()

	// 设置日志级别，如果未指定则默认为Info级别
	if cfg.LogLevel != 0 {
		logger.SetLevel(cfg.LogLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	chatHistory := chat.NewMemoryChatHistory()

	// 设置日志格式为JSON
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.Info("start init server")

	srv := &Server{
		echo:        echo.New(),
		chatModel:   cfg.ChatModel,
		logger:      logger,
		ChatHistory: chatHistory,
	}

	// 配置中间件
	srv.echo.Use(middleware.Logger())
	srv.echo.Use(middleware.Recover())

	// 配置CORS，允许前端访问
	srv.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// 注册路由
	srv.registerRoutes()

	logger.Info("server init success")
	return srv, nil
}

// registerRoutes 注册所有路由
func (s *Server) registerRoutes() {
	s.logger.Info("start register routes")

	// 健康检查
	s.echo.GET("/health", func(c echo.Context) error {
		s.logger.Debug("health check")
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// 聊天接口
	chatRouter := s.echo.Group("/chat")
	chatRouter.POST("/chat", s.handleChat)
	chatRouter.GET("/chat/ids", s.handleListChatIDs)
	chatRouter.DELETE("/chat/:chatID", s.handleDeleteChat)
	chatRouter.GET("/chat/:chatID", s.handleGetChatHistory)

	s.logger.Info("routes registered")
}

// Start 启动服务器
func (s *Server) Start(port string) error {
	s.logger.WithField("port", port).Info("server start")
	return s.echo.Start(port)
}

// Shutdown 优雅关闭服务器
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("server shutdown")
	err := s.echo.Shutdown(ctx)
	if err != nil {
		s.logger.WithError(err).Error("server close error")
	} else {
		s.logger.Info("server close success")
	}
	return err
}
