package handlers

import (
	"log"

	"greenride/internal/config"
	"greenride/internal/i18n"
)

// 初始化多语言支持
func InitI18n() {
	// 创建翻译器
	translator := i18n.NewFileTranslator()

	// 设置全局翻译器
	i18n.SetGlobalTranslator(translator)

	// 获取配置
	cfg := config.Get()
	if cfg == nil || cfg.I18n == nil || cfg.I18n.LocalesDir == "" {
		log.Printf("I18n: locales_dir not configured, skipping translation loading")
		return
	}

	// 使用配置中的路径
	localesDir := cfg.I18n.LocalesDir

	// 加载翻译资源
	if err := translator.LoadTranslations(localesDir); err != nil {
		log.Printf("I18n: Failed to load translations from '%s': %v", localesDir, err)
	} else {
		log.Printf("I18n: Translations loaded successfully from '%s'", localesDir)
	}
}
