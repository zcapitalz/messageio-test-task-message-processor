package domain

import "github.com/segmentio/ksuid"

type Message struct {
	ID     ksuid.KSUID `json:"id"`
	Text   string      `json:"text"`
	Status MessageStatus
}

type MessageStatus string

const (
	MessageStatusNotProcessed MessageStatus = "not-processed"
	MessageStatusProcessed    MessageStatus = "processed"
)
