package messagecontroller

import (
	"message-processor/internal/domain"
	"time"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	messageService MessageService
}

//go:generate mockery --name MessageService --filename message_service.go
type MessageService interface {
	SaveMessage(messageDTO *domain.SaveMessageDTO) error
	GetMessageProcessingStats(startTime, endTime time.Time) (*domain.MessageProcessingStats, error)
}

func NewMessageController(messageService MessageService) *MessageController {
	return &MessageController{
		messageService: messageService,
	}
}

func (c *MessageController) RegisterRoutes(engine *gin.Engine) {
	engine.POST("/messages", c.saveMessage)
	engine.GET("/messages/processing/stats", c.getMessageProcessingStats)
}
