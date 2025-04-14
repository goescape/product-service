package middlewares

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func Logger(router *gin.Engine) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] | %d | %v | %s | %s %s | %s | %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method, param.Path,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

func UnaryLoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("[Success] = %s - [%s]",
		info.FullMethod, start.Format(time.RFC1123))
	return resp, err
}
