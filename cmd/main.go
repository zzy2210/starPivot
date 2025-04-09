package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"starPivot/internal/config"
	"starPivot/internal/server"

	"github.com/spf13/cobra"
)

var (
	configPath string
	logLevel   string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "starPivot",
		Short: "StarPivot server application",
		Run: func(cmd *cobra.Command, args []string) {
			startServer()
		},
	}

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "internal/config/config.ini", "Path to config file")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Log level (debug|info|warn|error)")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Command execution failed: %v", err)
	}
}

func startServer() {
	ctx := context.Background()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load server config: %v", err)
	}

	// 设置日志级别
	if err := setLogLevel(logLevel); err != nil {
		log.Printf("Invalid log level: %v, using default", err)
	}

	srv, err := server.NewServer(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
}

func setLogLevel(level string) error {
	// TODO: 实现日志级别设置逻辑
	return nil
}
