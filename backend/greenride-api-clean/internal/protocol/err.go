package protocol

// Error 定义错误结构
type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

// NewError 创建新的错误
func NewError(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Error 实现 error 接口
func (e *Error) Error() string {
	return e.Message
}
