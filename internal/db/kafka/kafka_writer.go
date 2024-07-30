package kafka

import (
	"message-processor/internal/config"
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

func NewWriter(cfg *config.KafkaWriterConfig) (*kafka.Writer, error) {
	tlsConfig, err := newTLSConfig(cfg.KafkaSSLConfig)
	if err != nil {
		return nil, errors.Wrap(err, "create TLS config")
	}
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS:       tlsConfig,
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:     cfg.Brokers,
		Topic:       cfg.Topic,
		Balancer:    &kafka.LeastBytes{},
		Dialer:      dialer,
		Logger:      newLogger("Kafka writer: "),
		ErrorLogger: newErrorLogger("Kafka writer: "),
	})
	writer.AllowAutoTopicCreation = true
	return writer, nil
}
