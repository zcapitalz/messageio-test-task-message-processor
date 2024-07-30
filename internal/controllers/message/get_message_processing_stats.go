package messagecontroller

import (
	"message-processor/internal/controllers/httputils"
	"message-processor/internal/domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type getMessageProcessingStatsRequestQuery struct {
	StartTime time.Time `form:"startTime" binding:"required"`
	EndTime   time.Time `form:"endTime" binding:"required"`
}

type getMessageProcessingStatsResponseBody struct {
	Stats domain.MessageProcessingStats `json:"stats"`
}

// @Summary Get message processing stats
// @Description Get message processing statistics within a specified time range
// @Tags Message
// @Produce json
// @Param startTime query string true "Start time for the stats query"
// @Param endTime query string true "End time for the stats query"
// @Success 200 {object} getMessageProcessingStatsResponseBody
// @Failure 400 {object} httputils.HTTPError "Bad request"
// @Failure 500 {object} httputils.HTTPError "Internal server error"
// @Router /messages/processing/stats [get]
func (c *MessageController) getMessageProcessingStats(ctx *gin.Context) {
	var reqQuery getMessageProcessingStatsRequestQuery
	err := ctx.BindQuery(&reqQuery)
	if err != nil {
		httputils.BindQueryError(ctx, err)
		return
	}

	stats, err := c.messageService.
		GetMessageProcessingStats(reqQuery.StartTime, reqQuery.EndTime)
	if err != nil {
		httputils.InternalError(ctx)
		return
	}

	ctx.JSON(
		http.StatusOK,
		getMessageProcessingStatsResponseBody{
			Stats: *stats})
}
