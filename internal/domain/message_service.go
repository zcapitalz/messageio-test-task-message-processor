package domain

import (
	"log/slog"
	"message-processor/internal/utils/slogutils"
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
)

type MessageService struct {
	messageStorage         ServiceMessageStorage
	messageProcessingQueue MessageProcessingQueueWriter
}

type ServiceMessageStorage interface {
	SaveMessage(message *Message) error
	DeleteMessageByID(messageID ksuid.KSUID) error
	GetMessageProcessingStats(startTime, endTime time.Time) (*MessageProcessingStats, error)
}

type MessageProcessingQueueWriter interface {
	Write(message *Message) error
}

func NewMessageService(
	messageStorage ServiceMessageStorage,
	messageProcessingQueue MessageProcessingQueueWriter) *MessageService {
	return &MessageService{
		messageStorage:         messageStorage,
		messageProcessingQueue: messageProcessingQueue,
	}
}

func (s *MessageService) SaveMessage(messageDTO *SaveMessageDTO) error {
	message := &Message{
		ID:     ksuid.New(),
		Text:   messageDTO.Text,
		Status: MessageStatusNotProcessed,
	}

	err := s.messageStorage.SaveMessage(message)
	if err != nil {
		err = errors.Wrap(err, "save message")
		slog.Error("", slogutils.ErrorAttr(err))
		return err
	}

	err = s.messageProcessingQueue.Write(message)
	if err != nil {
		err = errors.Wrap(err, "send message to queue")
		slog.Error("", slogutils.ErrorAttr(err))

		err1 := s.messageStorage.DeleteMessageByID(message.ID)
		if err1 != nil {
			slog.Error("", slogutils.ErrorAttrWrap(err, "delete message because could not put it into processing queue"))
		}

		return err
	}

	return err
}

func (s *MessageService) GetMessageProcessingStats(startTime, endTime time.Time) (*MessageProcessingStats, error) {
	stats, err := s.messageStorage.GetMessageProcessingStats(startTime, endTime)
	if err != nil {
		err = errors.Wrap(err, "send message to queue")
		slog.Error("", slogutils.ErrorAttr(err))
		return nil, err
	}

	return stats, nil
}
