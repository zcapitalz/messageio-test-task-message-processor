package processor

import (
	"context"
	"log/slog"
	"message-processor/internal/config"
	"message-processor/internal/db/kafka"
	"message-processor/internal/db/postgres"
	"message-processor/internal/domain"
	"message-processor/internal/storages"
	"message-processor/internal/utils/slogutils"
)

func Process(cfg *config.ProcessorConfig) {
	logger := slogutils.MustNewLogger(cfg.Env)
	slog.SetDefault(logger)

	slog.Info("Setting up processor dependencies")

	postgresClient, err := postgres.NewClient(cfg.Postgres)
	if err != nil {
		slog.Error("create Postgres client", slogutils.ErrorAttr(err))
		return
	}
	messageProcessingQueueKafkaReader, err :=
		kafka.NewReader(&cfg.KafkaReaderConfig)
	if err != nil {
		slog.Error("create Kafka writer", slogutils.ErrorAttr(err))
		return
	}
	defer messageProcessingQueueKafkaReader.Close()

	messageStorage := storages.NewMessageStorage(postgresClient)
	messageProcessingQueueReader :=
		storages.NewMessageProcessingQueueReader(messageProcessingQueueKafkaReader)

	messageProcessor, err := domain.NewMessageProcessor(
		messageStorage, messageProcessingQueueReader, cfg.MaxProcessors)
	if err != nil {
		slog.Error("create message processor", slogutils.ErrorAttr(err))
		return
	}

	slog.Info("Starting processor")
	messageProcessor.ProcessMessages(context.Background())
}
