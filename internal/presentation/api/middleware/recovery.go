package middleware

import (
	apierrors "ContractAnalysis/pkg/errors"
	"ContractAnalysis/pkg/utils"
	"ContractAnalysis/internal/infrastructure/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery returns a recovery middleware
func Recovery(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				log.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				// Return error response
				apiErr := apierrors.NewInternalServerError("Internal server error")
				utils.ErrorResponse(c, apiErr)

				c.Abort()
			}
		}()

		c.Next()
	}
}
