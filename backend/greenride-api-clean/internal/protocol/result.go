package protocol

import (
	"greenride/internal/i18n"
)

// 基础结果结构
type Result struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

// Page 分页请求参数
type Pagination struct {
	Size int `json:"size" form:"size" binding:"min=1" example:"10"` // 每页记录数
	Page int `json:"page" form:"page" binding:"min=1" example:"1"`  // 当前页码，从1开始
}

// GetOffset 获取数据库查询的偏移量
func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.Size
}

// GetLimit 获取数据库查询的限制数
func (p *Pagination) GetLimit() int {
	return p.Size
}

// NewDefaultPagination 创建默认分页参数（第1页，每页10条）
func NewDefaultPagination() *Pagination {
	return &Pagination{
		Page: 1,
		Size: 10,
	}
}

// PageResult 分页数据结构
type PageResult struct {
	ResultType string         `json:"result_type"` // 数据类型
	Size       int64          `json:"size"`        // 每页记录数
	Current    int64          `json:"current"`     // 当前页码
	Total      int64          `json:"total"`       // 总页数
	Count      int64          `json:"count"`       // 总记录数
	Records    any            `json:"records"`     // 记录列表
	Attach     map[string]any `json:"attach"`      // 附加数据
}

// CountPage 计算总页数
func (p *PageResult) CountPage() int64 {
	if p.Size == 0 {
		return 0
	}
	p.Count = p.Total / p.Size
	if p.Total%p.Size > 0 {
		p.Count++
	}
	return p.Count
}

// NewPageResult 创建分页数据结构
func NewPageResult(records any, total int64, page *Pagination) *PageResult {
	if page == nil {
		page = NewDefaultPagination()
	}
	count := total / int64(page.Size)
	if total%int64(page.Size) > 0 {
		count++
	}
	return &PageResult{
		Size:    int64(page.Size),
		Current: int64(page.Page),
		Total:   total,
		Count:   count,
		Records: records,
		Attach:  make(map[string]any),
	}
}

// NewSuccessPageResult 创建分页成功响应
func NewSuccessPageResult(records interface{}, total int64, page *Pagination) *Result {
	return NewSuccessResult(NewPageResult(records, total, page))
}

// AddAttach 添加附加数据
func (p *PageResult) AddAttach(key string, value interface{}) {
	if p.Attach == nil {
		p.Attach = make(map[string]interface{})
	}
	p.Attach[key] = value
}

// 常用的响应码（保持向后兼容）
const (
	CODE_SUCCESS        = "0000" // 对应 codes.Success
	CODE_PARAM_ERROR    = "2001" // 对应 codes.InvalidParams
	CODE_AUTH_ERROR     = "3000" // 对应 codes.AuthenticationFailed
	CODE_BUSINESS_ERROR = "1003" // 对应 codes.NetworkError
	CODE_SYSTEM_ERROR   = "1000" // 对应 codes.SystemError
)

// NewResult 创建带错误码和多语言支持的响应
func NewResult(code ErrorCode, lang string, data interface{}, args ...interface{}) *Result {
	return &Result{
		Code: code.GetCode(),
		Msg:  i18n.Translate(code.GetCode(), lang, args...),
		Data: data,
	}
}

// NewSuccessResult 创建成功响应
func NewSuccessResult(data interface{}) *Result {
	return NewSuccessResultWithLang(data, i18n.DefaultLanguage)
}

// NewSuccessResultWithLang 创建多语言成功响应
func NewSuccessResultWithLang(data any, lang string) *Result {
	return &Result{
		Code: Success.GetCode(),
		Msg:  i18n.Translate(Success.GetCode(), lang),
		Data: data,
	}
}

// NewSuccessMessageResult 创建带成功消息的响应
func NewSuccessMessageResult(messageKey string, lang string, args ...interface{}) *Result {
	return &Result{
		Code: Success.GetCode(),
		Msg:  i18n.TranslateMessage(messageKey, lang, args...),
		Data: nil,
	}
}

// NewErrorResult 创建错误响应
func NewErrorResult(code ErrorCode, lang string, args ...interface{}) *Result {
	// Get translated message
	msg := i18n.Translate(code.GetCode(), lang, args...)

	// If the message is the same as the error code, attempt to use the English message as fallback
	if msg == code.GetCode() {
		msg = code.GetMessage() // Use the default English message from the ErrorCode
	}

	return &Result{
		Code: code.GetCode(),
		Msg:  msg,
		Data: nil,
	}
}

// NewErrorResultWithMessage 创建带自定义消息的错误响应
func NewErrorResultWithMessage(code ErrorCode, lang string, message string) *Result {
	return &Result{
		Code: code.GetCode(),
		Msg:  message,
		Data: nil,
	}
}

// 保持向后兼容的方法
// NewParamErrorResult 创建参数错误响应
func NewParamErrorResult(message string) *Result {
	return &Result{
		Code: CODE_PARAM_ERROR,
		Msg:  message,
	}
}

// NewParamErrorResultWithLang 创建多语言参数错误响应
func NewParamErrorResultWithLang(lang string, args ...interface{}) *Result {
	return NewErrorResult(InvalidParams, lang, args...)
}

// NewAuthErrorResult 创建认证错误响应
func NewAuthErrorResult() *Result {
	return &Result{
		Code: CODE_AUTH_ERROR,
		Msg:  "Authentication failed",
	}
}

// NewAuthErrorResultWithLang 创建多语言认证错误响应
func NewAuthErrorResultWithLang(lang string) *Result {
	return NewErrorResult(AuthenticationFailed, lang)
}

// NewBusinessErrorResult 创建业务错误响应
func NewBusinessErrorResult(message string) *Result {
	return &Result{
		Code: CODE_BUSINESS_ERROR,
		Msg:  message,
	}
}

// NewBusinessErrorResultWithCode 创建带错误码的业务错误响应
func NewBusinessErrorResultWithCode(code ErrorCode, lang string, args ...interface{}) *Result {
	return NewErrorResult(code, lang, args...)
}

// NewSystemErrorResult 创建系统错误响应
func NewSystemErrorResult(message string) *Result {
	return &Result{
		Code: CODE_SYSTEM_ERROR,
		Msg:  message,
	}
}

// NewSystemErrorResultWithLang 创建多语言系统错误响应
func NewSystemErrorResultWithLang(lang string) *Result {
	return NewErrorResult(SystemError, lang)
}
