package middleware

import (
	"strings"

	"greenride/internal/i18n"

	"github.com/gin-gonic/gin"
)

// LanguageContextKey 语言上下文键
const LanguageContextKey = "language"

// LanguageMiddleware 语言检测中间件
func LanguageMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		lang := detectLanguage(c)
		c.Set(LanguageContextKey, lang)
		c.Next()
	})
}

// detectLanguage 检测用户语言偏好
func detectLanguage(c *gin.Context) string {
	// 1. 优先使用查询参数中的语言设置
	if langParam := c.Query("lang"); langParam != "" {
		if isValidLanguage(langParam) {
			return normalizeLanguage(langParam)
		}
	}

	// 2. 检查请求头中的语言设置
	if langHeader := c.GetHeader("X-Language"); langHeader != "" {
		if isValidLanguage(langHeader) {
			return normalizeLanguage(langHeader)
		}
	}

	// 3. 检查Accept-Language头
	acceptLang := c.GetHeader("Accept-Language")
	if detectedLang := i18n.GetLanguageFromAcceptLanguage(acceptLang); detectedLang != "" {
		if isValidLanguage(detectedLang) {
			return normalizeLanguage(detectedLang)
		}
	}

	// 4. 检查用户偏好设置（从JWT或session中获取）
	if userClaim := getUserClaimFromContext(c); userClaim != nil {
		// 这里可以从用户信息中获取语言偏好
		// 暂时跳过，因为需要定义用户结构
	}

	// 5. 默认语言
	return i18n.DefaultLanguage
}

// isValidLanguage 检查语言是否受支持
func isValidLanguage(lang string) bool {
	supportedLangs := []string{
		i18n.LanguageEnglish,
		i18n.LanguageChinese,
		i18n.LanguageChineseTraditional,
		i18n.LanguageSpanish,
		i18n.LanguageFrench,
		i18n.LanguageGerman,
		i18n.LanguageJapanese,
		i18n.LanguageKorean,
		i18n.LanguageArabic,
		i18n.LanguageRussian,
	}

	normalizedLang := normalizeLanguage(lang)
	for _, supportedLang := range supportedLangs {
		if normalizedLang == supportedLang {
			return true
		}
	}
	return false
}

// normalizeLanguage 规范化语言代码
func normalizeLanguage(lang string) string {
	if lang == "" {
		return i18n.DefaultLanguage
	}

	// 转换为小写
	lang = strings.ToLower(lang)

	// 处理带地区的语言代码
	if strings.Contains(lang, "-") {
		parts := strings.SplitN(lang, "-", 2)
		primaryLang := parts[0]

		// 特殊处理中文
		if primaryLang == "zh" {
			if strings.Contains(lang, "tw") || strings.Contains(lang, "hk") || strings.Contains(lang, "mo") {
				return i18n.LanguageChineseTraditional
			}
			return i18n.LanguageChinese
		}

		return primaryLang
	}

	return lang
}

// getUserClaimFromContext 从上下文获取用户信息
func getUserClaimFromContext(c *gin.Context) interface{} {
	// 这里可以获取JWT中的用户信息
	// 暂时返回nil，等待用户结构定义
	if value, exists := c.Get("user_claim"); exists {
		return value
	}
	return nil
}

// GetLanguageFromContext 从上下文获取语言设置
func GetLanguageFromContext(c *gin.Context) string {
	if lang, exists := c.Get(LanguageContextKey); exists {
		if langStr, ok := lang.(string); ok {
			return langStr
		}
	}
	return i18n.DefaultLanguage
}

// SetLanguageInContext 在上下文中设置语言
func SetLanguageInContext(c *gin.Context, lang string) {
	c.Set(LanguageContextKey, normalizeLanguage(lang))
}
