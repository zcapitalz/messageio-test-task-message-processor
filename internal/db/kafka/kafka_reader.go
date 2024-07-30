package kafka

import (
	"message-processor/internal/config"
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

func NewReader(cfg *config.KafkaReaderConfig) (*kafka.Reader, error) {
	tlsConfig, err := newTLSConfig(cfg.KafkaSSLConfig)
	if err != nil {
		return nil, errors.Wrap(err, "create TLS config")
	}
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS:       tlsConfig,
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.Brokers,
		Topic:       cfg.Topic,
		GroupID:     cfg.ConsumerGroupID,
		Dialer:      dialer,
		Logger:      newLogger("Kafka reader: "),
		ErrorLogger: newErrorLogger("Kafka reader: "),
	})
	return reader, nil
}
