package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"starPivot/internal/config"
	"starPivot/internal/server"
)

func main() {
	// 创建上下文
	ctx := context.Background()

	// 创建服务器配置
	// 服务器配置从 confing.ini 文件中读取
	cfg, err := config.LoadConfig("config.ini")
	if err != nil {
		log.Fatalf("Failed to load server config: %v", err)
	}
	// 创建服务器
	srv, err := server.NewServer(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// 启动服务器
	go func() {
		if err := srv.Start(); err != nil {
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
