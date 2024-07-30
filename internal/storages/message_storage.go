package storages

import (
	"fmt"
	"log/slog"
	"message-processor/internal/domain"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
)

type MessageStorage struct {
	db      *sqlx.DB
	builder sq.StatementBuilderType
}

func NewMessageStorage(db *sqlx.DB) *MessageStorage {
	return &MessageStorage{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *MessageStorage) SaveMessage(message *domain.Message) error {
	builder := s.builder.
		Insert("messages").
		Columns("id, text, status").
		Values(message.ID, message.Text, message.Status)

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "build query")
	}
	slog.Debug(fmt.Sprintf("SQL query: %s", query))

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "execute query")
	}

	return nil
}

func (s *MessageStorage) DeleteMessageByID(messageID ksuid.KSUID) error {
	builder := s.builder.
		Delete("messages").
		Where(sq.Eq{"id": messageID})

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "build query")
	}
	slog.Debug(fmt.Sprintf("SQL query: %s", query))

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "execute query")
	}

	return nil
}

func (s *MessageStorage) SetMessageStatusProcessed(messageID ksuid.KSUID, status domain.MessageStatus) error {
	builder := s.builder.
		Update("messages").
		Set("status", status).
		Set("processed_at", time.Now())

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "build query")
	}
	slog.Debug(fmt.Sprintf("SQL query: %s", query))

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "execute query")
	}

	return nil
}

func (s *MessageStorage) GetMessageProcessingStats(startTime, endTime time.Time) (*domain.MessageProcessingStats, error) {
	createdMessagesCount, err := s.getCreatedMessagesCountByPeriod(startTime, endTime)
	if err != nil {
		return nil, errors.Wrap(err, "get created messages count")
	}
	processedMessagesCount, err := s.getProcessedMessagesCountByPeriod(startTime, endTime)
	if err != nil {
		return nil, errors.Wrap(err, "get processed messages count")
	}

	return &domain.MessageProcessingStats{
		SavedTotal:     createdMessagesCount,
		ProcessedTotal: processedMessagesCount,
	}, nil
}

func (s *MessageStorage) getCreatedMessagesCountByPeriod(startTime, endTime time.Time) (int, error) {
	builder := s.builder.
		Select("COUNT(*)").
		From("messages").
		Where(sq.And{sq.GtOrEq{"created_at": startTime}, sq.LtOrEq{"created_at": endTime}})

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "build query")
	}
	slog.Debug(fmt.Sprintf("SQL query: %s", query))

	var res int
	err = s.db.QueryRow(query, args...).Scan(&res)
	if err != nil {
		return 0, errors.Wrap(err, "execute query")
	}

	return res, nil
}

func (s *MessageStorage) getProcessedMessagesCountByPeriod(startTime, endTime time.Time) (int, error) {
	builder := s.builder.
		Select("COUNT(*)").
		From("messages").
		Where(sq.And{sq.GtOrEq{"processed_at": startTime}, sq.LtOrEq{"processed_at": endTime}})

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "build query")
	}
	slog.Debug(fmt.Sprintf("SQL query: %s", query))

	var res int
	err = s.db.QueryRow(query, args...).Scan(&res)
	if err != nil {
		return 0, errors.Wrap(err, "execute query")
	}

	return res, nil
}
