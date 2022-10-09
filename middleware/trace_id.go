package middleware

import (
	"github.com/gin-gonic/gin"
	"main/internal/uuid"
)

func AddTraceId() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := uuid.CreateUUID()
		c.Set("trace_id", traceId)
		c.Next()
	}
}
