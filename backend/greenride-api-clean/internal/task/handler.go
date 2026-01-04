package task

import (
	"context"
	"greenride/internal/protocol"
)

// Handler 任务处理函数类型
type Handler func(ctx context.Context, params protocol.MapData) error

// handlers 任务处理器注册表
var handlers = make(map[string]Handler)

// RegisterHandler 注册任务处理器
func RegisterHandler(key string, handler Handler) {
	handlers[key] = handler
}

// GetHandler 获取任务处理器
func GetHandler(key string) (Handler, bool) {
	handler, exists := handlers[key]
	return handler, exists
}
