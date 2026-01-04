package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// 支持的语言列表
const (
	LanguageEnglish            = "en"
	LanguageChinese            = "zh"
	LanguageChineseTraditional = "zh-TW"
	LanguageSpanish            = "es"
	LanguageFrench             = "fr"
	LanguageKinyarwanda        = "rw"
	LanguageGerman             = "de"
	LanguageJapanese           = "ja"
	LanguageKorean             = "ko"
	LanguageArabic             = "ar"
	LanguageRussian            = "ru"
)

// DefaultLanguage 默认语言
const DefaultLanguage = LanguageEnglish

// ProjectSupportedLanguages 项目实际支持的语言列表
var ProjectSupportedLanguages = []string{
	LanguageEnglish,     // en - 英语
	LanguageChinese,     // zh - 中文
	LanguageFrench,      // fr - 法语
	LanguageKinyarwanda, // rw - 卢旺达语
}

// IsValidLanguage 检查语言代码是否被项目支持
func IsValidLanguage(lang string) bool {
	for _, supportedLang := range ProjectSupportedLanguages {
		if lang == supportedLang {
			return true
		}
	}
	return false
}

// GetValidLanguage 获取有效的语言代码，如果不支持则返回默认语言
func GetValidLanguage(lang string) string {
	if IsValidLanguage(lang) {
		return lang
	}
	return DefaultLanguage
}

// Translator 翻译器接口
type Translator interface {
	Translate(code string, lang string, args ...interface{}) string
	TranslateMessage(key string, lang string, args ...interface{}) string
	GetSupportedLanguages() []string
	SetDefaultLanguage(lang string)
	LoadTranslations(localesDir string) error
}

// FileTranslator 基于文件的翻译器
type FileTranslator struct {
	mu                sync.RWMutex
	defaultLanguage   string
	supportedLangs    []string
	translations      map[string]map[string]string // [language][key]message
	errorTranslations map[string]map[string]string // [language][code]message
}

// NewFileTranslator 创建基于文件的翻译器
func NewFileTranslator() *FileTranslator {
	return &FileTranslator{
		defaultLanguage:   DefaultLanguage,
		supportedLangs:    []string{LanguageEnglish, LanguageChinese, LanguageFrench, LanguageKinyarwanda},
		translations:      make(map[string]map[string]string),
		errorTranslations: make(map[string]map[string]string),
	}
}

// SetDefaultLanguage 设置默认语言
func (ft *FileTranslator) SetDefaultLanguage(lang string) {
	ft.mu.Lock()
	defer ft.mu.Unlock()
	ft.defaultLanguage = lang
}

// GetSupportedLanguages 获取支持的语言列表
func (ft *FileTranslator) GetSupportedLanguages() []string {
	ft.mu.RLock()
	defer ft.mu.RUnlock()
	return append([]string{}, ft.supportedLangs...)
}

// LoadTranslations 从指定目录加载翻译文件
func (ft *FileTranslator) LoadTranslations(localesDir string) error {
	ft.mu.Lock()
	defer ft.mu.Unlock()

	// 检查目录是否存在
	if _, err := os.Stat(localesDir); os.IsNotExist(err) {
		return fmt.Errorf("locales directory does not exist: %s", localesDir)
	}

	// 遍历语言目录
	entries, err := os.ReadDir(localesDir)
	if err != nil {
		return fmt.Errorf("failed to read locales directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		lang := entry.Name()
		langDir := filepath.Join(localesDir, lang)

		// 加载通用翻译
		if err := ft.loadGeneralTranslations(lang, langDir); err != nil {
			return fmt.Errorf("failed to load general translations for %s: %w", lang, err)
		}

		// 加载错误码翻译
		if err := ft.loadErrorTranslations(lang, langDir); err != nil {
			return fmt.Errorf("failed to load error translations for %s: %w", lang, err)
		}
	}

	return nil
}

// loadGeneralTranslations 加载通用翻译
func (ft *FileTranslator) loadGeneralTranslations(lang, langDir string) error {
	filePath := filepath.Join(langDir, "common.json")

	// 如果文件不存在，跳过
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var translations map[string]string
	if err := json.Unmarshal(data, &translations); err != nil {
		return err
	}

	if ft.translations[lang] == nil {
		ft.translations[lang] = make(map[string]string)
	}

	for key, value := range translations {
		ft.translations[lang][key] = value
	}

	return nil
}

// loadErrorTranslations 加载错误码翻译
func (ft *FileTranslator) loadErrorTranslations(lang, langDir string) error {
	filePath := filepath.Join(langDir, "errors.json")

	// 如果文件不存在，跳过
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var errorMessages map[string]string
	if err := json.Unmarshal(data, &errorMessages); err != nil {
		return err
	}

	if ft.errorTranslations[lang] == nil {
		ft.errorTranslations[lang] = make(map[string]string)
	}

	// 直接存储code字符串到消息的映射
	for codeStr, message := range errorMessages {
		ft.errorTranslations[lang][codeStr] = message
	}

	return nil
}

// Translate 翻译错误码
func (ft *FileTranslator) Translate(code string, lang string, args ...interface{}) string {
	ft.mu.RLock()
	defer ft.mu.RUnlock()

	// 规范化语言代码
	normalizedLang := ft.normalizeLanguage(lang)

	// 查找翻译
	if errorTranslations, exists := ft.errorTranslations[normalizedLang]; exists {
		if message, found := errorTranslations[code]; found {
			if len(args) > 0 {
				return ft.safeFormat(message, args...)
			}
			return message
		}
	}

	// 降级到默认语言
	if normalizedLang != ft.defaultLanguage {
		if errorTranslations, exists := ft.errorTranslations[ft.defaultLanguage]; exists {
			if message, found := errorTranslations[code]; found {
				if len(args) > 0 {
					return ft.safeFormat(message, args...)
				}
				return message
			}
		}
	}

	// Log if translation is missing - helps diagnose issues like 7102 showing as message
	fmt.Printf("Translation missing for code: %s, language: %s\n", code, lang)

	// 如果找不到翻译，返回原始错误码作为后备
	if len(args) > 0 {
		return ft.safeFormat(code, args...)
	}
	return code
}

// TranslateMessage 翻译通用消息
func (ft *FileTranslator) TranslateMessage(key string, lang string, args ...interface{}) string {
	ft.mu.RLock()
	defer ft.mu.RUnlock()

	// 规范化语言代码
	normalizedLang := ft.normalizeLanguage(lang)

	// 查找翻译
	if translations, exists := ft.translations[normalizedLang]; exists {
		if message, found := translations[key]; found {
			if len(args) > 0 {
				return ft.safeFormat(message, args...)
			}
			return message
		}
	}

	// 降级到默认语言
	if normalizedLang != ft.defaultLanguage {
		if translations, exists := ft.translations[ft.defaultLanguage]; exists {
			if message, found := translations[key]; found {
				if len(args) > 0 {
					return ft.safeFormat(message, args...)
				}
				return message
			}
		}
	}

	// 返回原始key
	if len(args) > 0 {
		return ft.safeFormat(key, args...)
	}
	return key
}

// normalizeLanguage 规范化语言代码
func (ft *FileTranslator) normalizeLanguage(lang string) string {
	if lang == "" {
		return ft.defaultLanguage
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
				return LanguageChineseTraditional
			}
			return LanguageChinese
		}

		return primaryLang
	}

	return lang
}

// safeFormat 安全地格式化字符串，防止panic
func (ft *FileTranslator) safeFormat(format string, args ...interface{}) (result string) {
	if len(args) == 0 {
		return format
	}

	// 使用defer recover来捕获fmt.Sprintf可能的panic
	defer func() {
		if r := recover(); r != nil {
			// 如果格式化失败，返回原始字符串
			result = format
		}
	}()

	// 尝试格式化字符串
	result = fmt.Sprintf(format, args...)
	return result
}

// GetLanguageFromAcceptLanguage 从Accept-Language头解析语言
func GetLanguageFromAcceptLanguage(acceptLang string) string {
	if acceptLang == "" {
		return DefaultLanguage
	}

	// 解析Accept-Language头
	// 例: "zh-CN,zh;q=0.9,en;q=0.8"
	languages := strings.Split(acceptLang, ",")

	for _, lang := range languages {
		// 移除权重
		lang = strings.TrimSpace(strings.Split(lang, ";")[0])
		if lang != "" {
			// 转换为标准格式
			if strings.HasPrefix(lang, "zh") {
				if strings.Contains(lang, "TW") || strings.Contains(lang, "HK") || strings.Contains(lang, "MO") {
					return LanguageChineseTraditional
				}
				return LanguageChinese
			}

			// 提取主要语言代码
			primaryLang := strings.Split(lang, "-")[0]
			return strings.ToLower(primaryLang)
		}
	}

	return DefaultLanguage
}

// 全局翻译器实例
var globalTranslator Translator

// SetGlobalTranslator 设置全局翻译器
func SetGlobalTranslator(translator Translator) {
	globalTranslator = translator
}

// GetGlobalTranslator 获取全局翻译器
func GetGlobalTranslator() Translator {
	if globalTranslator == nil {
		globalTranslator = NewFileTranslator()
	}
	return globalTranslator
}

// Translate 使用全局翻译器翻译错误码
func Translate(code string, lang string, args ...interface{}) string {
	return GetGlobalTranslator().Translate(code, lang, args...)
}

// TranslateMessage 使用全局翻译器翻译通用消息
func TranslateMessage(key string, lang string, args ...interface{}) string {
	return GetGlobalTranslator().TranslateMessage(key, lang, args...)
}

// T 从gin context中获取语言信息并翻译消息
func T(c interface{}, key string, args ...interface{}) string {
	// 这里可以从context中获取Accept-Language头或其他语言设置
	// 暂时使用默认语言
	lang := DefaultLanguage

	// 如果传入的是gin.Context，可以从header中获取语言设置
	if ginCtx, ok := c.(interface {
		GetHeader(string) string
	}); ok {
		acceptLang := ginCtx.GetHeader("Accept-Language")
		if acceptLang != "" {
			if strings.HasPrefix(acceptLang, "zh") {
				lang = LanguageChinese
			}
		}
	}

	return GetGlobalTranslator().TranslateMessage(key, lang, args...)
}
