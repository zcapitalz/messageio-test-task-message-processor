package domain

import (
	"context"
	"fmt"
	"log/slog"
	"message-processor/internal/utils/slogutils"
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/memorystore"
)

const (
	processorCreationLimitPerMinute = 60
	processorCreationLimitWaitTime  = time.Millisecond * 500
)

type MessageProcessor struct {
	maxProcessors          int
	messageStorage         ProcessorMessageStorage
	messageProcessingQueue MessageProcessingQueueReader
}

type ProcessorMessageStorage interface {
	SetMessageStatusProcessed(messageID ksuid.KSUID, status MessageStatus) error
}

type MessageProcessingQueueReader interface {
	Read(ctx context.Context) (*Message, CommitMessagesFunc, error)
}

type CommitMessagesFunc func(context.Context) error

func NewMessageProcessor(
	messageStorage ProcessorMessageStorage,
	messageProcessingQueue MessageProcessingQueueReader,
	maxProcessors int) (*MessageProcessor, error) {

	if maxProcessors < 0 {
		return nil, fmt.Errorf("maxProcessors should be > 0")
	}

	return &MessageProcessor{
		messageProcessingQueue: messageProcessingQueue,
		messageStorage:         messageStorage,
		maxProcessors:          maxProcessors,
	}, nil
}

func (p *MessageProcessor) ProcessMessages(ctx context.Context) error {
	if p.maxProcessors == 0 {
		slog.Info("maxProcessors is set to 0, exiting...")
		return nil
	}

	processors := make(chan struct{}, p.maxProcessors)
	processorCreationLimiter, err := newProcessorCreationLimiter()
	if err != nil {
		err = errors.Wrap(err, "create processor creation limiter")
		slog.Error(err.Error())
		return err
	}

	for {
		_, _, _, ok, err := processorCreationLimiter.Take(ctx, "processor-creation")
		if err != nil {
			return errors.Wrap(err, "take from processor creation limiter")
		}
		if !ok {
			time.Sleep(processorCreationLimitWaitTime)
			continue
		}

		p.createProcessor(ctx, processors)
	}
}

func (p *MessageProcessor) createProcessor(ctx context.Context, processors chan struct{}) {
	processors <- struct{}{}
	processorID := ksuid.New()
	processorIDSlogAttr := slog.Attr{
		Key:   "processorID",
		Value: slog.StringValue(processorID.String()),
	}
	slog.Info("Processor created", processorIDSlogAttr)

	go func() {
		defer func() {
			<-processors

			logArgs := []any{processorIDSlogAttr}
			if v := recover(); v != nil {
				logArgs = append(logArgs, slogutils.PanicAttr(v))
			}
			slog.Info("Processor stoped", logArgs...)
		}()

		err := p.processMessages(ctx)
		if err != nil {
			slog.Error("process messages", processorIDSlogAttr, slogutils.ErrorAttr(err))
		}
	}()
}

func (p *MessageProcessor) processMessages(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		msg, commitMessagesFunc, err := p.messageProcessingQueue.Read(ctx)
		if err != nil {
			errMsg := "receive message from queue"
			slog.Error(errMsg, slogutils.ErrorAttr(err))
			return errors.Wrap(err, errMsg)
		}

		err = p.messageStorage.SetMessageStatusProcessed(msg.ID, MessageStatusProcessed)
		if err != nil {
			return errors.Wrap(err, "update message status")
		}
		slog.Info("Processed message", "messageID", msg.ID)

		err = commitMessagesFunc(ctx)
		if err != nil {
			return errors.Wrap(err, "commit message")
		}
	}
}

func newProcessorCreationLimiter() (limiter.Store, error) {
	return memorystore.New(&memorystore.Config{
		Tokens:   processorCreationLimitPerMinute,
		Interval: time.Minute,
	})
}
