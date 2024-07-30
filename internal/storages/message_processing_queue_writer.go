package storages

import (
	"context"
	"encoding/json"
	"message-processor/internal/domain"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type MessageProcessingQueueWriter struct {
	writer *kafka.Writer
}

func NewMessageProcessingQueueWriter(writer *kafka.Writer) *MessageProcessingQueueWriter {
	return &MessageProcessingQueueWriter{
		writer: writer,
	}
}

func (q *MessageProcessingQueueWriter) Write(message *domain.Message) error {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "serialize message before sending to queue")
	}

	err = q.writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   message.ID.Bytes(),
			Value: messageJSON,
		})
	if err != nil {
		return errors.Wrap(err, "write message to kafka writer")
	}

	return nil
}
