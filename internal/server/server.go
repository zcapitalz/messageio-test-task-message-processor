package server

import (
	"context"
	"log/slog"
	_ "message-processor/api/server"
	"message-processor/internal/config"
	messagecontroller "message-processor/internal/controllers/message"
	"message-processor/internal/db/kafka"
	"message-processor/internal/db/postgres"
	"message-processor/internal/domain"
	"message-processor/internal/storages"
	"message-processor/internal/utils/slogutils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

// @title           Message processing http server
// @version         1.0
// @BasePath  /
func Serve(cfg config.ServerConfig) {
	logger := slogutils.MustNewLogger(cfg.Env)
	slog.SetDefault(logger)

	slog.Info("Setting up server dependencies")

	postgresClient, err := postgres.NewClient(cfg.Postgres)
	if err != nil {
		slog.Error("create Postgres client", slogutils.ErrorAttr(err))
		return
	}
	messageProcessingQueueKafkaWriter, err := kafka.NewWriter(&cfg.KafkaWriterConfig)
	if err != nil {
		slog.Error("create Kafka writer", slogutils.ErrorAttr(err))
		return
	}
	defer messageProcessingQueueKafkaWriter.Close()

	messageStorage := storages.NewMessageStorage(postgresClient)
	messageProcessingQueueWriter := storages.NewMessageProcessingQueueWriter(messageProcessingQueueKafkaWriter)

	messageService := domain.NewMessageService(messageStorage, messageProcessingQueueWriter)

	messageController := messagecontroller.NewMessageController(messageService)

	switch cfg.Env {
	case config.EnvLocal:
		gin.SetMode(gin.DebugMode)
	case config.EnvDev, config.EnvProd:
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(sloggin.New(logger))
	engine.Use(gin.Recovery())
	engine.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
	messageController.RegisterRoutes(engine)

	srv := &http.Server{
		Addr:    cfg.HTTPServer.IpAddress + ":" + cfg.HTTPServer.Port,
		Handler: engine.Handler(),
	}

	slog.Info("Starting server ...")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server listen: %s\n", slogutils.ErrorAttr(err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown:", slogutils.ErrorAttr(err))
		os.Exit(1)
	}

	select {
	case <-ctx.Done():
		slog.Info("timeout of 5 seconds.")
	}
	slog.Info("Server exiting")
}
