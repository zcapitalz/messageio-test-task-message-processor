package storages

import (
	"context"
	"encoding/json"
	"message-processor/internal/domain"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type MessageProcessingQueueReader struct {
	reader *kafka.Reader
}

func NewMessageProcessingQueueReader(reader *kafka.Reader) *MessageProcessingQueueReader {
	return &MessageProcessingQueueReader{
		reader: reader,
	}
}

func (q *MessageProcessingQueueReader) Read(ctx context.Context) (*domain.Message, domain.CommitMessagesFunc, error) {
	kafkaMessage, err := q.reader.ReadMessage(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "read message from kafka reader")
	}

	var message domain.Message
	err = json.Unmarshal(kafkaMessage.Value, &message)
	if err != nil {
		return nil, nil, errors.Wrap(err, "parse message")
	}

	commitMessages := func(ctx context.Context) error {
		return q.reader.CommitMessages(ctx, kafkaMessage)
	}

	return &message, commitMessages, nil
}
