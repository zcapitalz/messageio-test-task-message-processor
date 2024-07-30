package messagecontroller

import (
	"message-processor/internal/controllers/httputils"
	"message-processor/internal/domain"

	"github.com/gin-gonic/gin"
)

type saveMessageRequestBody struct {
	Message string `json:"message" binding:"required"`
}

// @Summary Save Message
// @Description Saves and process message.
// @Tags Message
// @Accept json
// @Param orderBook body saveMessageRequestBody true "Message data"
// @Success 200
// @Failure 400 {object} httputils.HTTPError "Invalid request body"
// @Failure 500 {object} httputils.HTTPError "Internal server error"
// @Router /messages [post]
func (c *MessageController) saveMessage(ctx *gin.Context) {
	var reqBody saveMessageRequestBody
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		httputils.BindJSONBodyError(ctx, err)
		return
	}

	err = c.messageService.SaveMessage(&domain.SaveMessageDTO{Text: reqBody.Message})
	if err != nil {
		httputils.InternalError(ctx)
		return
	}
}
