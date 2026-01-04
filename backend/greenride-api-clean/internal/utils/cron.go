package utils

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"greenride/internal/log"

	"github.com/robfig/cron/v3"
)

// CalculateNextTime 计算下次执行时间
// 支持标准cron表达式和简化的时间表达式
// 简化表达式格式: "every 1h", "every 30m", "every 1d" 等
func CalculateNextTime(cronExpr string) time.Time {
	// 首先尝试解析简化的时间表达式
	if duration := parseSimpleTimeExpr(cronExpr); duration > 0 {
		return time.Now().Add(duration)
	}

	// 如果不是简化表达式，尝试解析标准cron表达式
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(cronExpr)
	if err != nil {
		log.Get().Errorf("Parse cron expression error: %v, expression: %s", err, cronExpr)
		// 解析失败时，默认1分钟后重试
		return time.Now().Add(time.Minute)
	}

	// 计算下次执行时间
	next := schedule.Next(time.Now())

	// 如果计算失败（返回零值），默认1分钟后重试
	if next.IsZero() {
		log.Get().Errorf("Calculate next time failed, expression: %s", cronExpr)
		return time.Now().Add(time.Minute)
	}

	return next
}

// parseSimpleTimeExpr 解析简化的时间表达式
// 支持格式: "every 1h", "every 30m", "every 1d", "every 2w" 等
func parseSimpleTimeExpr(expr string) time.Duration {
	// 匹配 "every 数字单位" 格式
	re := regexp.MustCompile(`^every\s+(\d+)([smhdw])$`)
	matches := re.FindStringSubmatch(strings.ToLower(strings.TrimSpace(expr)))

	if len(matches) != 3 {
		return 0
	}

	num, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0
	}

	unit := matches[2]
	switch unit {
	case "s": // 秒
		return time.Duration(num) * time.Second
	case "m": // 分钟
		return time.Duration(num) * time.Minute
	case "h": // 小时
		return time.Duration(num) * time.Hour
	case "d": // 天
		return time.Duration(num) * 24 * time.Hour
	case "w": // 周
		return time.Duration(num) * 7 * 24 * time.Hour
	default:
		return 0
	}
}
