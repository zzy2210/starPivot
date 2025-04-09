package server

import (
	"context"
	"fmt"
	"net/http"

	"starPivot/internal/config"
	"starPivot/internal/model"
	"starPivot/internal/service/chat/history"

	eModel "github.com/cloudwego/eino/components/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Server 封装了 Echo 服务器和 ChatModel
type Server struct {
	echo       *echo.Echo
	logger     *logrus.Logger
	cfg        *config.Config
	ChatConfig *model.BaseAIConfig

	chatModel   eModel.ChatModel
	ChatHistory history.ChatHistoryFactory
}

// NewServer 创建一个新的服务器实例
func NewServer(ctx context.Context, cfg *config.Config) (*Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// 初始化日志
	logger := logrus.New()

	// 设置日志级别，如果未指定则默认为Info级别
	if cfg.ServerConfig.LogLevel != 0 {
		logger.SetLevel(cfg.ServerConfig.LogLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
	// 设置日志格式为JSON
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.Info("start init server")

	srv := &Server{
		echo:   echo.New(),
		logger: logger,
		cfg:    cfg,
	}

	if err := srv.SetHistoryStorage(); err != nil {
		return nil, err
	}

	// 配置中间件
	srv.echo.Use(middleware.Logger())
	srv.echo.Use(middleware.Recover())

	//配置CORS，允许前端访问
	srv.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://192.168.50.241:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "X-Username"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
	}))

	// 注册路由
	srv.registerRoutes()

	logger.Info("server init success")
	return srv, nil
}

func (s *Server) SetHistoryStorage() error {
	switch s.cfg.ServerConfig.HistoryStorage {
	case "memory":
		s.ChatHistory = history.NewMemoryChatHistory()
	case "database":
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			s.cfg.DataBaseConfig.Host,
			s.cfg.DataBaseConfig.Port,
			s.cfg.DataBaseConfig.User,
			s.cfg.DataBaseConfig.Password,
			s.cfg.DataBaseConfig.DBName,
		)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("failed to connect database: %v", err)
		}
		s.ChatHistory = history.NewDatabaseChatHistory(db)
	default:
		return fmt.Errorf("invalid history storage: %s", s.cfg.ServerConfig.HistoryStorage)
	}
	return nil
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
	chatRouter.Use(s.chatMiddleware)
	chatRouter.POST("/chat", s.handleChat)
	chatRouter.GET("/ids", s.handleListChatIDs)
	chatRouter.DELETE("/:chatID", s.handleDeleteChat)
	chatRouter.GET("/:chatID", s.handleGetChatHistory)

	// 配置接口
	configRouter := s.echo.Group("/config")
	configRouter.POST("/model", s.handleConfig)

	s.logger.Info("routes registered")
}

// 中间件，/chat开头的路由，需要验证是否配置了 chatModel
func (s *Server) chatMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if s.chatModel == nil {
			return c.JSON(http.StatusOK, model.Response{
				Code:    http.StatusInternalServerError,
				Message: "chat model not configured",
				Data:    nil,
			})
		}
		return next(c)
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	s.logger.WithField("port", s.cfg.ServerConfig.Port).Info("server start")
	return s.echo.Start(fmt.Sprintf(":%s", s.cfg.ServerConfig.Port))
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
