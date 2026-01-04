package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"greenride/internal/log"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggerMiddleware 记录HTTP请求日志的中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果是 OPTIONS 请求，直接处理请求不记录日志
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// 开始时间
		start := time.Now()

		// 读取请求体
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = c.GetRawData()
			// 使用 NopCloser 包装 bytes.Buffer，确保实现 io.ReadCloser 接口
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 解析请求体到map（如果需要的话）
		requestBody := map[string]any{}
		if len(bodyBytes) > 0 {
			_ = json.Unmarshal(bodyBytes, &requestBody)
		}

		// 设置响应体写入器
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 计算持续时间
		duration := time.Since(start)

		// 获取请求URL的查询参数
		query := c.Request.URL.RawQuery
		path := c.Request.URL.Path
		if query != "" {
			path = path + "?" + query
		}

		// 获取响应体
		responseBody := blw.body.String()

		// 记录完整日志
		log.Get().WithFields(logrus.Fields{
			"status":   c.Writer.Status(),
			"method":   c.Request.Method,
			"path":     path,
			"ip":       c.ClientIP(),
			"duration": duration.String(),
			"request":  requestBody,
			"response": responseBody,
		}).Info(fmt.Sprintf("[API] %s %s", c.Request.Method, path))
	}
}

// RequestIDMiddleware 为每个请求生成唯一ID
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := generateRequestID()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// ErrorHandlerMiddleware 错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 处理错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			requestID, _ := c.Get("request_id")
			log.Get().WithFields(logrus.Fields{
				"request_id": requestID,
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"ip":         c.ClientIP(),
				"error":      err.Error(),
			}).Error("Request error occurred")
		}
	}
}

// RecoveryWithLogger 带日志记录的恢复中间件，处理panic
func RecoveryWithLogger() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID, _ := c.Get("request_id")
		log.Get().WithFields(logrus.Fields{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"panic":      recovered,
		}).Error("Request panic occurred")

		c.JSON(500, gin.H{
			"error":      "Internal Server Error",
			"request_id": requestID,
		})
	})
}
