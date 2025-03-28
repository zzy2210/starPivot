package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"starPivot/internal/server"
	"starPivot/internal/service/chat"
)

func main() {
	// 创建上下文
	ctx := context.Background()

	// 创建聊天模型
	chatModel := chat.CreateOPenAIChatModel(ctx)

	// 创建服务器配置
	cfg := &server.Config{
		Port:      ":8080",
		ChatModel: chatModel,
	}

	// 创建服务器
	srv, err := server.NewServer(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// 启动服务器
	go func() {
		if err := srv.Start(cfg.Port); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 优雅关闭
	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
