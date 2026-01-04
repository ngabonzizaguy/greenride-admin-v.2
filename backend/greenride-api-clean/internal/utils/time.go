package utils

import (
	"fmt"
	"time"
)

// timeNowMilli 返回当前时间的毫秒时间戳
func TimeNowMilli() int64 {
	return time.Now().UnixMilli()
}

// TimeNowPtr 返回当前时间的毫秒时间戳指针
func TimeNowPtr() *int64 {
	now := time.Now().UnixMilli()
	return &now
}

// TimeNowMilliPtr 返回当前时间的毫秒时间戳指针
func TimeNowMilliPtr() *int64 {
	now := time.Now().UnixMilli()
	return &now
}

// TimeNowSeconds 返回当前时间的秒时间戳
func TimeNowSeconds() int64 {
	return time.Now().Unix()
}

// TimeNowNano 返回当前时间的纳秒时间戳
func TimeNowNano() int64 {
	return time.Now().UnixNano()
}

// MilliToTime 将毫秒时间戳转换为time.Time
func MilliToTime(millis int64) time.Time {
	return time.Unix(0, millis*int64(time.Millisecond))
}

// TimeToMilli 将time.Time转换为毫秒时间戳
func TimeToMilli(t time.Time) int64 {
	return t.UnixMilli()
}

// IsExpired 检查给定的毫秒时间戳是否已过期
func IsExpired(expiry int64) bool {
	return expiry < TimeNowMilli()
}

// AddMinutes 给毫秒时间戳添加分钟
func AddMinutes(millis int64, minutes int) int64 {
	return millis + int64(minutes)*60*1000
}

// AddHours 给毫秒时间戳添加小时
func AddHours(millis int64, hours int) int64 {
	return millis + int64(hours)*60*60*1000
}

// AddDays 给毫秒时间戳添加天数
func AddDays(millis int64, days int) int64 {
	return millis + int64(days)*24*60*60*1000
}

// GetHourlyTimeRange 获取指定时间所在小时的时间范围（开始和结束时间戳）
// 参数:
//   - t: 任意时间点
//   - timezone: 时区字符串，如 "UTC", "Africa/Kigali" 等
//
// 返回:
//   - startAt: 小时开始的毫秒时间戳
//   - endAt: 小时结束的毫秒时间戳（下一小时开始）
//   - dateStr: 日期字符串，格式为 "2006-01-02"
//   - hour: 小时数 (0-23)
//   - error: 如果时区无效则返回错误
func GetHourlyTimeRange(t time.Time, timezone string) (startAt int64, endAt int64, dateStr string, hour int, err error) {
	// 加载时区
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return 0, 0, "", 0, fmt.Errorf("无效的时区: %s", timezone)
	}

	// 转换为指定时区的时间
	localTime := t.In(loc)

	// 获取小时开始时间
	hourStart := time.Date(
		localTime.Year(),
		localTime.Month(),
		localTime.Day(),
		localTime.Hour(),
		0, 0, 0,
		loc,
	)

	// 获取小时结束时间（下一小时的开始）
	hourEnd := hourStart.Add(time.Hour)

	// 转换为毫秒时间戳
	startAt = hourStart.UnixMilli()
	endAt = hourEnd.UnixMilli()

	// 日期字符串和小时数
	dateStr = hourStart.Format("2006-01-02")
	hour = hourStart.Hour()

	return startAt, endAt, dateStr, hour, nil
}

// GetHourlyTimeRangeByParams 根据指定的日期和小时参数获取时间范围
// 参数:
//   - date: 日期，如 "2025-09-10"
//   - hour: 小时数 (0-23)
//   - timezone: 时区字符串
//
// 返回:
//   - startAt: 小时开始的毫秒时间戳
//   - endAt: 小时结束的毫秒时间戳（下一小时开始）
//   - error: 如果日期格式或时区无效则返回错误
func GetHourlyTimeRangeByParams(date string, hour int, timezone string) (startAt int64, endAt int64, err error) {
	// 验证小时范围
	if hour < 0 || hour > 23 {
		return 0, 0, fmt.Errorf("无效的小时数: %d，应为 0-23", hour)
	}

	// 解析日期字符串
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0, 0, fmt.Errorf("无效的日期格式: %s, 应为 YYYY-MM-DD", date)
	}

	// 加载时区
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return 0, 0, fmt.Errorf("无效的时区: %s", timezone)
	}

	// 创建指定日期和小时的时间
	hourStart := time.Date(
		parsedDate.Year(),
		parsedDate.Month(),
		parsedDate.Day(),
		hour,
		0, 0, 0,
		loc,
	)

	// 获取小时结束时间（下一小时开始）
	hourEnd := hourStart.Add(time.Hour)

	// 转换为毫秒时间戳
	startAt = hourStart.UnixMilli()
	endAt = hourEnd.UnixMilli()

	return startAt, endAt, nil
}
