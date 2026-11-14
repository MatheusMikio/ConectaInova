package response

import (
	"net/http"

	"conectainova/internal/domain/models"
	"github.com/gin-gonic/gin"
)

func SendSuccess(ctx *gin.Context, statusCode int, data any) {
	if statusCode == http.StatusNoContent {
		ctx.Status(statusCode)
		return
	}
	ctx.JSON(statusCode, gin.H{"data": data})
}

func SendError(ctx *gin.Context, statusCode int, err *models.ErrorMessage) {
	ctx.JSON(statusCode, gin.H{"error": err})
}

func SendErrors(ctx *gin.Context, statusCode int, err []*models.ErrorMessage) {
	ctx.JSON(statusCode, gin.H{"error": err})
}
