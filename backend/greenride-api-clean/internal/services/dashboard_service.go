package services

import (
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"greenride/internal/models"
	"greenride/internal/protocol"

	"gorm.io/gorm"
)

// 时区常量

// 金额格式化工具函数

// formatAmount 格式化金额为2位小数
func formatAmount(amount float64) float64 {
	return math.Round(amount*100) / 100
}

// formatPercentage 格式化百分比为2位小数，最终显示为xx.xx%格式
func formatPercentage(percentage float64) float64 {
	return math.Round(percentage*100) / 100
}

// 时区相关的辅助函数

// getRwandaTimezone 获取卢旺达时区
func getRwandaTimezone() *time.Location {
	// 卢旺达时区为 UTC+2 (Africa/Kigali)
	loc, err := time.LoadLocation("Africa/Kigali")
	if err != nil {
		// 如果无法加载时区，使用固定偏移 UTC+2
		return time.FixedZone("Rwanda", 2*3600)
	}
	return loc
}

// parseTimezone 解析时区字符串
func parseTimezone(timezone string) (*time.Location, error) {
	if timezone == "" {
		return getRwandaTimezone(), nil
	}

	// 支持的时区格式
	switch timezone {
	case "UTC", "utc":
		return time.UTC, nil
	case "Rwanda", "rwanda", "Africa/Kigali":
		return getRwandaTimezone(), nil
	default:
		// 尝试解析标准时区名称
		if loc, err := time.LoadLocation(timezone); err == nil {
			return loc, nil
		}
		// 如果解析失败，返回卢旺达时区作为默认
		return getRwandaTimezone(), nil
	}
}

// getNowInTimezone 获取指定时区的当前时间
func getNowInTimezone(timezone string) time.Time {
	loc, err := parseTimezone(timezone)
	if err != nil {
		loc = getRwandaTimezone()
	}
	return time.Now().In(loc)
}

// getDateRangeInTimezone 在指定时区计算日期范围的Unix时间戳
// 返回该日期在指定时区的0点和次日0点的Unix时间戳（毫秒）
// getDaysBackInTimezone 获取指定时区下过去N天的日期列表和查询时间戳范围

// getDaysBackInTimezone 获取指定时区下往前N天的日期列表和时间范围
func getDaysBackInTimezone(days int, timezone string) ([]time.Time, int64, int64) {
	loc, err := parseTimezone(timezone)
	if err != nil {
		loc = getRwandaTimezone()
	}

	// 获取当前时区的今天
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	// 计算开始日期（往前days-1天）
	startDate := today.AddDate(0, 0, -(days - 1))

	// 生成日期列表
	dateList := make([]time.Time, days)
	for i := 0; i < days; i++ {
		dateList[i] = startDate.AddDate(0, 0, i)
	}

	// 计算查询时间范围（UTC时间戳毫秒）
	startTimestamp := startDate.UTC().UnixMilli()
	endTimestamp := today.AddDate(0, 0, 1).UTC().UnixMilli() // 明天0点

	return dateList, startTimestamp, endTimestamp
}

// DashboardService 仪表盘服务
type DashboardService struct {
	db *gorm.DB
}

// 默认时区常量
const (
	DefaultTimezone = "Africa/Kigali" // 卢旺达时区
)

var (
	dashboardServiceInstance *DashboardService
	dashboardServiceOnce     sync.Once
)

// GetDashboardService 获取仪表盘服务单例
func GetDashboardService() *DashboardService {
	dashboardServiceOnce.Do(func() {
		SetupDashboardService()
	})
	return dashboardServiceInstance
}

// SetupDashboardService 设置仪表盘服务
func SetupDashboardService() {
	dashboardServiceInstance = &DashboardService{
		db: models.GetDB(),
	}
}

// NewDashboardService 创建仪表盘服务
func NewDashboardService(db *gorm.DB) *DashboardService {
	return &DashboardService{
		db: db,
	}
}

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	TotalUsers         int64                `json:"total_users"`
	TotalDrivers       int64                `json:"total_drivers"`
	TotalVehicles      int64                `json:"total_vehicles"`
	TotalTrips         int64                `json:"total_trips"`
	TotalRevenue       float64              `json:"total_revenue"`
	ActiveTrips        int64                `json:"active_trips"`
	MonthlyGrowth      MonthlyGrowth        `json:"monthly_growth"`
	RecentTrips        []RecentTrip         `json:"recent_trips"`
	TopDrivers         []TopDriver          `json:"top_drivers"`
	VehicleUtilization []VehicleUtilization `json:"vehicle_utilization"`
}

// MonthlyGrowth 月度增长数据
type MonthlyGrowth struct {
	Users   float64 `json:"users"`
	Trips   float64 `json:"trips"`
	Revenue float64 `json:"revenue"`
}

// RecentTrip 最近行程
type RecentTrip struct {
	ID              int64      `json:"id"`
	UserID          int64      `json:"user_id"`
	DriverID        int64      `json:"driver_id"`
	VehicleID       int64      `json:"vehicle_id"`
	PickupLocation  string     `json:"pickup_location"`
	DropoffLocation string     `json:"dropoff_location"`
	Status          string     `json:"status"`
	Fare            float64    `json:"fare"`
	Distance        float64    `json:"distance"`
	Duration        int        `json:"duration"`
	PaymentMethod   string     `json:"payment_method"`
	PaymentStatus   string     `json:"payment_status"`
	Rating          int        `json:"rating"`
	RequestedAt     time.Time  `json:"requested_at"`
	StartedAt       *time.Time `json:"started_at"`
	CompletedAt     *time.Time `json:"completed_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// TopDriver 顶级司机
type TopDriver struct {
	ID            int64   `json:"id"`
	UserID        int64   `json:"user_id"`
	DriverName    string  `json:"driver_name"`
	Rating        float64 `json:"rating"`
	TotalTrips    int64   `json:"total_trips"`
	TotalEarnings float64 `json:"total_earnings"`
	Status        string  `json:"status"`
}

// VehicleUtilization 车辆利用率
type VehicleUtilization struct {
	VehicleID       int64   `json:"vehicle_id"`
	VehicleName     string  `json:"vehicle_name"`
	UtilizationRate float64 `json:"utilization_rate"`
	TotalTrips      int64   `json:"total_trips"`
}

// RevenueData 收入数据
type RevenueData struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
	Trips   int64   `json:"trips"`
	Users   int64   `json:"users"`
}

// UserGrowthData 用户增长数据
type UserGrowthData struct {
	Date         string `json:"date"`
	NewUsers     int64  `json:"new_users"`
	TotalUsers   int64  `json:"total_users"`
	NewDrivers   int64  `json:"new_drivers"`
	TotalDrivers int64  `json:"total_drivers"`
}

// GetDashboardStats 获取仪表盘统计数据
func (s *DashboardService) GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	// 使用结构体接收所有计数查询结果
	type CountResults struct {
		TotalUsers    int64           `json:"total_users"`
		TotalDrivers  int64           `json:"total_drivers"`
		TotalVehicles int64           `json:"total_vehicles"`
		TotalTrips    int64           `json:"total_trips"`
		TotalRevenue  sql.NullFloat64 `json:"total_revenue"`
		ActiveTrips   int64           `json:"active_trips"`
	}

	// 执行事务以确保数据一致性
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var counts CountResults

		// 获取用户总数和司机总数
		var userCounts struct {
			TotalUsers   int64 `json:"total_users"`
			TotalDrivers int64 `json:"total_drivers"`
		}
		if err := tx.Raw(`
			SELECT 
				COUNT(*) as total_users,
				SUM(CASE WHEN user_type = ? THEN 1 ELSE 0 END) as total_drivers
			FROM t_users
		`, protocol.UserTypeDriver).Scan(&userCounts).Error; err != nil {
			return fmt.Errorf("failed to count users and drivers: %w", err)
		}
		counts.TotalUsers = userCounts.TotalUsers
		counts.TotalDrivers = userCounts.TotalDrivers
		fmt.Printf("DEBUG: User counts - Users: %d, Drivers: %d\n", counts.TotalUsers, counts.TotalDrivers)

		// 获取车辆总数
		if err := tx.Model(&models.Vehicle{}).Count(&counts.TotalVehicles).Error; err != nil {
			return fmt.Errorf("failed to count vehicles: %w", err)
		}
		fmt.Printf("DEBUG: Vehicle count: %d\n", counts.TotalVehicles)

		// 获取订单相关统计: 总订单数、总收入和活跃订单数
		// 总订单数
		if err := tx.Model(&models.Order{}).Count(&counts.TotalTrips).Error; err != nil {
			return fmt.Errorf("failed to count total trips: %w", err)
		}
		fmt.Printf("DEBUG: Total trips: %d\n", counts.TotalTrips)

		// 总收入（已完成订单）
		if err := tx.Model(&models.Order{}).
			Where("status = ? AND payment_amount IS NOT NULL", protocol.StatusCompleted).
			Select("COALESCE(SUM(payment_amount), 0)").
			Scan(&counts.TotalRevenue).Error; err != nil {
			return fmt.Errorf("failed to calculate total revenue: %w", err)
		}

		// 活跃订单数
		if err := tx.Model(&models.Order{}).
			Where("status IN ?", []string{protocol.StatusRequested, protocol.StatusAccepted, protocol.StatusInProgress}).
			Count(&counts.ActiveTrips).Error; err != nil {
			return fmt.Errorf("failed to count active trips: %w", err)
		}

		// 将结果赋值给 stats
		stats.TotalUsers = counts.TotalUsers
		stats.TotalDrivers = counts.TotalDrivers
		stats.TotalVehicles = counts.TotalVehicles
		stats.TotalTrips = counts.TotalTrips
		stats.TotalRevenue = formatAmount(counts.TotalRevenue.Float64)
		stats.ActiveTrips = counts.ActiveTrips

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 计算月度增长
	monthlyGrowth, err := s.calculateMonthlyGrowth("")
	if err != nil {
		return nil, fmt.Errorf("failed to calculate monthly growth: %w", err)
	}
	stats.MonthlyGrowth = *monthlyGrowth

	// 获取最近行程
	recentTrips, err := s.getRecentTrips(10)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent trips: %w", err)
	}
	stats.RecentTrips = recentTrips

	// 获取顶级司机
	topDrivers, err := s.getTopDrivers(5)
	if err != nil {
		return nil, fmt.Errorf("failed to get top drivers: %w", err)
	}
	stats.TopDrivers = topDrivers

	// 获取车辆利用率
	vehicleUtilization, err := s.getVehicleUtilization(5)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle utilization: %w", err)
	}
	stats.VehicleUtilization = vehicleUtilization

	return stats, nil
}

// calculateMonthlyGrowth 计算月度增长
func (s *DashboardService) calculateMonthlyGrowth(timezone string) (*MonthlyGrowth, error) {
	now := getNowInTimezone(timezone)
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastMonth := currentMonth.AddDate(0, -1, 0)
	// 时间范围已经足够，不需要更早的时间

	growth := &MonthlyGrowth{}

	// 定义结构体接收用户增长数据
	type UserGrowthStats struct {
		CurrentMonthUsers int64 `json:"current_month_users"`
		LastMonthUsers    int64 `json:"last_month_users"`
	}

	// 一次性查询用户增长数据
	var userStats UserGrowthStats
	userQuery := `
		SELECT
			(SELECT COUNT(*) FROM t_users WHERE created_at >= ?) AS current_month_users,
			(SELECT COUNT(*) FROM t_users WHERE created_at >= ? AND created_at < ?) AS last_month_users
	`
	if err := s.db.Raw(userQuery,
		currentMonth.UnixMilli(),
		lastMonth.UnixMilli(),
		currentMonth.UnixMilli()).Scan(&userStats).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate user growth: %w", err)
	}

	// 计算用户增长率
	if userStats.LastMonthUsers > 0 {
		growth.Users = formatPercentage(float64(userStats.CurrentMonthUsers-userStats.LastMonthUsers) / float64(userStats.LastMonthUsers) * 100)
	}

	// 定义结构体接收订单增长数据
	type OrderGrowthStats struct {
		CurrentMonthTrips   int64   `json:"current_month_trips"`
		LastMonthTrips      int64   `json:"last_month_trips"`
		CurrentMonthRevenue float64 `json:"current_month_revenue"`
		LastMonthRevenue    float64 `json:"last_month_revenue"`
	}

	// 一次性查询订单增长数据
	var orderStats OrderGrowthStats
	orderQuery := `
		SELECT
			(SELECT COUNT(*) FROM t_orders WHERE created_at >= ?) AS current_month_trips,
			(SELECT COUNT(*) FROM t_orders WHERE created_at >= ? AND created_at < ?) AS last_month_trips,
			(SELECT COALESCE(SUM(payment_amount), 0) FROM t_orders WHERE created_at >= ? AND status = ?) AS current_month_revenue,
			(SELECT COALESCE(SUM(payment_amount), 0) FROM t_orders WHERE created_at >= ? AND created_at < ? AND status = ?) AS last_month_revenue
	`
	if err := s.db.Raw(orderQuery,
		currentMonth.UnixMilli(),
		lastMonth.UnixMilli(),
		currentMonth.UnixMilli(),
		currentMonth.UnixMilli(),
		protocol.StatusCompleted,
		lastMonth.UnixMilli(),
		currentMonth.UnixMilli(),
		protocol.StatusCompleted).Scan(&orderStats).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate order growth: %w", err)
	}

	// 计算行程增长率
	if orderStats.LastMonthTrips > 0 {
		growth.Trips = formatPercentage(float64(orderStats.CurrentMonthTrips-orderStats.LastMonthTrips) / float64(orderStats.LastMonthTrips) * 100)
	}

	// 计算收入增长率
	if orderStats.LastMonthRevenue > 0 {
		growth.Revenue = formatPercentage((orderStats.CurrentMonthRevenue - orderStats.LastMonthRevenue) / orderStats.LastMonthRevenue * 100)
	}

	return growth, nil
}

// getRecentTrips 获取最近行程
func (s *DashboardService) getRecentTrips(limit int) ([]RecentTrip, error) {
	var orders []models.Order
	if err := s.db.Model(&models.Order{}).
		Order("created_at DESC").
		Limit(limit).
		Find(&orders).Error; err != nil {
		return nil, err
	}

	trips := make([]RecentTrip, 0, len(orders))
	for _, order := range orders {
		paymentAmount, _ := order.GetPaymentAmount().Float64()
		trip := RecentTrip{
			ID:            int64(order.ID),
			Status:        order.GetStatus(),
			Fare:          paymentAmount,
			PaymentMethod: order.GetPaymentMethod(),
			PaymentStatus: order.GetPaymentStatus(),
			RequestedAt:   time.Unix(order.CreatedAt/1000, 0),
			CreatedAt:     time.Unix(order.CreatedAt/1000, 0),
			UpdatedAt:     time.Unix(order.CreatedAt/1000, 0), // 使用CreatedAt作为fallback
		}

		// 解析用户ID、司机ID等
		if userIDStr := order.GetUserID(); userIDStr != "" {
			if userID, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
				trip.UserID = userID
			}
		}

		// 设置位置信息（简化处理）
		trip.PickupLocation = "Pickup Location"
		trip.DropoffLocation = "Dropoff Location"

		trips = append(trips, trip)
	}

	return trips, nil
}

// getTopDrivers 获取顶级司机
func (s *DashboardService) getTopDrivers(limit int) ([]TopDriver, error) {
	type DriverStats struct {
		DriverID      string  `json:"driver_id"`
		TotalTrips    int64   `json:"total_trips"`
		TotalEarnings float64 `json:"total_earnings"`
		AvgRating     float64 `json:"avg_rating"`
	}

	var driverStats []DriverStats
	if err := s.db.Model(&models.Order{}).
		Select("provider_id as driver_id, COUNT(*) as total_trips, SUM(payment_amount) as total_earnings, 5.0 as avg_rating").
		Where("status = ? AND provider_id IS NOT NULL AND provider_id != ''", protocol.StatusCompleted).
		Group("provider_id").
		Order("total_earnings DESC").
		Limit(limit).
		Find(&driverStats).Error; err != nil {
		return nil, err
	}

	drivers := make([]TopDriver, 0, len(driverStats))
	for _, stat := range driverStats {
		// 尝试解析driver_id为数字，如果失败则使用0
		driverIDInt := int64(0)
		if id, err := strconv.ParseInt(stat.DriverID, 10, 64); err == nil {
			driverIDInt = id
		}

		driver := TopDriver{
			ID:            driverIDInt,
			UserID:        driverIDInt,
			DriverName:    fmt.Sprintf("Driver %s", stat.DriverID),
			Rating:        stat.AvgRating,
			TotalTrips:    stat.TotalTrips,
			TotalEarnings: formatAmount(stat.TotalEarnings),
			Status:        "active",
		}
		drivers = append(drivers, driver)
	}

	return drivers, nil
}

// getVehicleUtilization 获取车辆利用率
func (s *DashboardService) getVehicleUtilization(limit int) ([]VehicleUtilization, error) {
	type VehicleStats struct {
		VehicleID   string `json:"vehicle_id"`
		VehicleName string `json:"vehicle_name"`
		TotalTrips  int64  `json:"total_trips"`
	}

	// 通过JOIN查询获取车辆利用率（基于provider_id关联车辆的driver_id）
	var vehicleStats []VehicleStats
	query := `
		SELECT 
			v.vehicle_id,
			CONCAT(COALESCE(v.brand, ''), ' ', COALESCE(v.model, '')) as vehicle_name,
			COUNT(o.id) as total_trips
		FROM t_vehicles v
		INNER JOIN t_orders o ON v.driver_id = o.provider_id
		WHERE o.status = ? AND v.driver_id IS NOT NULL AND v.driver_id != ''
		GROUP BY v.vehicle_id, v.brand, v.model
		ORDER BY total_trips DESC
		LIMIT ?
	`

	if err := s.db.Raw(query, protocol.StatusCompleted, limit).Scan(&vehicleStats).Error; err != nil {
		return nil, err
	}

	utilizations := make([]VehicleUtilization, 0, len(vehicleStats))
	for _, stat := range vehicleStats {
		// 尝试解析vehicle_id为数字，如果失败则使用0
		vehicleIDInt := int64(0)
		if id, err := strconv.ParseInt(stat.VehicleID, 10, 64); err == nil {
			vehicleIDInt = id
		}

		// 简化的利用率计算（实际应该基于工作时间）
		utilizationRate := float64(stat.TotalTrips) / 100.0 * 85 // 模拟计算
		if utilizationRate > 100 {
			utilizationRate = 95
		}

		// 清理车辆名称
		vehicleName := strings.TrimSpace(stat.VehicleName)
		if vehicleName == "" {
			vehicleName = fmt.Sprintf("Vehicle %s", stat.VehicleID)
		}

		utilization := VehicleUtilization{
			VehicleID:       vehicleIDInt,
			VehicleName:     vehicleName,
			UtilizationRate: formatPercentage(utilizationRate),
			TotalTrips:      stat.TotalTrips,
		}
		utilizations = append(utilizations, utilization)
	}

	return utilizations, nil
}

// GetRevenueChart 获取收入图表数据
func (s *DashboardService) GetRevenueChart(period string, timezone string) ([]RevenueData, error) {
	var data []RevenueData

	switch period {
	case "7d":
		data = s.getRevenueByDays(7, timezone)
	case "30d":
		data = s.getRevenueByDays(30, timezone)
	case "12m":
		data = s.getRevenueByMonths(12, timezone)
	default:
		data = s.getRevenueByMonths(12, timezone)
	}

	return data, nil
}

// getRevenueByDays 按天获取收入数据
func (s *DashboardService) getRevenueByDays(days int, timezone string) []RevenueData {
	if days <= 0 {
		days = 7
	}

	// 在应用层计算时区下的日期范围
	dateList, queryStartTimestamp, queryEndTimestamp := getDaysBackInTimezone(days, timezone)

	// 定义数据结构接收查询结果
	type DailyStats struct {
		CreatedAt     int64           `json:"created_at"`
		Status        string          `json:"status"`
		PaymentAmount sql.NullFloat64 `json:"payment_amount"`
	}

	// 查询指定时间范围内的订单数据
	var allOrders []DailyStats
	query := `
		SELECT created_at, status, payment_amount
		FROM t_orders
		WHERE created_at >= ? AND created_at < ?
		ORDER BY created_at ASC
	`

	s.db.Raw(query, queryStartTimestamp, queryEndTimestamp).Scan(&allOrders)

	// 解析时区用于数据分组
	loc, err := parseTimezone(timezone)
	if err != nil {
		loc = getRwandaTimezone()
	}

	// 在应用层按日期分组统计
	type DailyRevenue struct {
		Date    string  `json:"date"`
		Revenue float64 `json:"revenue"`
		Trips   int64   `json:"trips"`
	}

	revenueMap := make(map[string]DailyRevenue)
	for _, order := range allOrders {
		// 将Unix时间戳转换为指定时区的日期
		orderTime := time.Unix(order.CreatedAt/1000, (order.CreatedAt%1000)*1000000).In(loc)
		dateKey := orderTime.Format("2006-01-02")

		data := revenueMap[dateKey]
		data.Date = dateKey
		data.Trips++

		// 只统计已完成订单的收入
		if order.Status == protocol.StatusCompleted && order.PaymentAmount.Valid {
			data.Revenue += order.PaymentAmount.Float64
		}

		revenueMap[dateKey] = data
	}

	// 按顺序组装所有日期的数据
	result := make([]RevenueData, 0, days)
	for _, date := range dateList {
		dateStr := date.Format("2006-01-02")
		formattedDate := date.Format("Jan 02")

		// 获取当天数据
		revenue := 0.0
		trips := int64(0)
		userCount := int64(0) // 这里可以根据需要实现用户数统计

		if dailyData, ok := revenueMap[dateStr]; ok {
			revenue = dailyData.Revenue
			trips = dailyData.Trips
		}

		result = append(result, RevenueData{
			Date:    formattedDate,
			Revenue: formatAmount(revenue),
			Trips:   trips,
			Users:   userCount,
		})
	}

	return result
}

// getRevenueByMonths 按月获取收入数据
func (s *DashboardService) getRevenueByMonths(months int, timezone string) []RevenueData {
	data := make([]RevenueData, 0, months)
	now := getNowInTimezone(timezone)

	// 计算时间范围
	startDate := now.AddDate(0, -(months - 1), 0)
	startOfStartMonth := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
	endOfCurrentMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())

	// 创建月份映射，用于快速查找和填充
	monthMap := make(map[string]time.Time)
	for i := 0; i < months; i++ {
		date := now.AddDate(0, -i, 0)
		monthKey := fmt.Sprintf("%d-%02d", date.Year(), date.Month())
		monthMap[monthKey] = date
	}

	// 定义数据结构接收查询结果
	type MonthlyStats struct {
		Year    int             `json:"year"`
		Month   int             `json:"month"`
		Revenue sql.NullFloat64 `json:"revenue"`
		Trips   int64           `json:"trips"`
	}

	// 查询指定时间范围内的订单数据
	var allOrders []struct {
		CreatedAt     int64   `json:"created_at"`
		Status        string  `json:"status"`
		PaymentAmount float64 `json:"payment_amount"`
	}

	// 使用简单的时间戳范围查询
	orderQuery := `
		SELECT created_at, status, COALESCE(payment_amount, 0) as payment_amount
		FROM t_orders
		WHERE created_at >= ? AND created_at < ?
		ORDER BY created_at ASC
	`

	s.db.Raw(orderQuery, startOfStartMonth.UnixMilli(), endOfCurrentMonth.UnixMilli()).Scan(&allOrders)

	// 解析时区用于数据分组
	loc, err := parseTimezone(timezone)
	if err != nil {
		loc = getRwandaTimezone()
	}

	// 在应用层按月份分组统计
	statsMap := make(map[string]MonthlyStats)
	for _, order := range allOrders {
		// 将Unix时间戳转换为指定时区的月份
		orderTime := time.Unix(order.CreatedAt/1000, (order.CreatedAt%1000)*1000000).In(loc)
		monthKey := fmt.Sprintf("%d-%02d", orderTime.Year(), orderTime.Month())

		data := statsMap[monthKey]
		data.Year = orderTime.Year()
		data.Month = int(orderTime.Month())
		data.Trips++
		if order.Status == protocol.StatusCompleted {
			data.Revenue.Float64 += order.PaymentAmount
			data.Revenue.Valid = true
		}
		statsMap[monthKey] = data
	}

	// 按照顺序组装所有月份的数据
	for i := months - 1; i >= 0; i-- {
		date := now.AddDate(0, -i, 0)
		monthKey := fmt.Sprintf("%d-%02d", date.Year(), date.Month())
		formattedDate := date.Format("Jan")

		// 如果有查询结果则使用，否则使用零值
		var revenue float64
		var trips int64

		if stat, ok := statsMap[monthKey]; ok {
			revenue = stat.Revenue.Float64
			trips = stat.Trips
		}

		data = append(data, RevenueData{
			Date:    formattedDate,
			Revenue: formatAmount(revenue),
			Trips:   trips,
		})
	}

	return data
}

// GetUserGrowthChart 获取用户增长图表数据
func (s *DashboardService) GetUserGrowthChart(period string, timezone string) ([]UserGrowthData, error) {
	var data []UserGrowthData

	switch period {
	case "7d":
		data = s.getUserGrowthByDays(7, timezone)
	case "30d":
		data = s.getUserGrowthByDays(30, timezone)
	case "12m":
		data = s.getUserGrowthByMonths(12, timezone)
	default:
		data = s.getUserGrowthByMonths(12, timezone)
	}

	return data, nil
}

// getUserGrowthByDays 按天获取用户增长数据
func (s *DashboardService) getUserGrowthByDays(days int, timezone string) []UserGrowthData {
	if days <= 0 {
		days = 7
	}

	// 在应用层计算时区下的日期范围
	dateList, queryStartTimestamp, queryEndTimestamp := getDaysBackInTimezone(days, timezone)

	// 定义查询结果结构
	type DailyNewUsers struct {
		Date       string `json:"date"`
		NewUsers   int64  `json:"new_users"`
		NewDrivers int64  `json:"new_drivers"`
	}

	// 查询指定时间范围内的用户数据，按created_at分组
	var allUsers []struct {
		CreatedAt int64  `json:"created_at"`
		UserType  string `json:"user_type"`
	}

	// 使用简单的时间戳范围查询，避免在WHERE子句中使用函数
	userQuery := `
		SELECT created_at, user_type
		FROM t_users
		WHERE created_at >= ? AND created_at < ?
		ORDER BY created_at ASC
	`

	s.db.Raw(userQuery, queryStartTimestamp, queryEndTimestamp).Scan(&allUsers)

	// 解析时区用于数据分组
	loc, err := parseTimezone(timezone)
	if err != nil {
		loc = getRwandaTimezone()
	}

	// 在应用层按日期分组统计
	newUsersMap := make(map[string]DailyNewUsers)
	for _, user := range allUsers {
		// 将Unix时间戳转换为指定时区的日期
		userTime := time.Unix(user.CreatedAt/1000, (user.CreatedAt%1000)*1000000).In(loc)
		dateKey := userTime.Format("2006-01-02")

		data := newUsersMap[dateKey]
		data.Date = dateKey
		data.NewUsers++
		if user.UserType == protocol.UserTypeDriver {
			data.NewDrivers++
		}
		newUsersMap[dateKey] = data
	}

	// 查询在统计开始日期之前的用户总数作为基准
	var baseTotal struct {
		TotalUsers   int64 `json:"total_users"`
		TotalDrivers int64 `json:"total_drivers"`
	}

	baseQuery := `
		SELECT 
			COUNT(*) as total_users,
			SUM(CASE WHEN user_type = ? THEN 1 ELSE 0 END) as total_drivers
		FROM t_users
		WHERE created_at < ?
	`

	s.db.Raw(baseQuery, protocol.UserTypeDriver, queryStartTimestamp).Scan(&baseTotal)

	// 计算每天的累计总数（从基准开始累加）
	runningTotalUsers := baseTotal.TotalUsers
	runningTotalDrivers := baseTotal.TotalDrivers

	// 按顺序组装所有日期的数据
	result := make([]UserGrowthData, 0, days)
	for _, date := range dateList {
		dateStr := date.Format("2006-01-02")
		formattedDate := date.Format("Jan 02")

		// 获取当天新增数据
		newUsers := int64(0)
		newDrivers := int64(0)

		if dailyNew, ok := newUsersMap[dateStr]; ok {
			newUsers = dailyNew.NewUsers
			newDrivers = dailyNew.NewDrivers
		}

		// 累加到总数
		runningTotalUsers += newUsers
		runningTotalDrivers += newDrivers

		result = append(result, UserGrowthData{
			Date:         formattedDate,
			NewUsers:     newUsers,
			TotalUsers:   runningTotalUsers,
			NewDrivers:   newDrivers,
			TotalDrivers: runningTotalDrivers,
		})
	}

	return result
}

// getUserGrowthByMonths 按月获取用户增长数据
func (s *DashboardService) getUserGrowthByMonths(months int, timezone string) []UserGrowthData {
	data := make([]UserGrowthData, 0, months)
	now := getNowInTimezone(timezone)

	// 计算时间范围
	startDate := now.AddDate(0, -(months - 1), 0)
	startOfStartMonth := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
	endOfCurrentMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())

	// 创建月份映射，用于结果处理
	monthMap := make(map[string]time.Time)
	for i := 0; i < months; i++ {
		date := now.AddDate(0, -i, 0)
		monthKey := fmt.Sprintf("%d-%02d", date.Year(), date.Month())
		monthMap[monthKey] = date
	}

	// 查询每月新增用户和司机
	type MonthlyNewUsers struct {
		Year       int   `json:"year"`
		Month      int   `json:"month"`
		NewUsers   int64 `json:"new_users"`
		NewDrivers int64 `json:"new_drivers"`
	}

	// 查询指定时间范围内的用户数据
	var allUsers []struct {
		CreatedAt int64  `json:"created_at"`
		UserType  string `json:"user_type"`
	}

	// 使用简单的时间戳范围查询
	userQuery := `
		SELECT created_at, user_type
		FROM t_users
		WHERE created_at >= ? AND created_at < ?
		ORDER BY created_at ASC
	`

	s.db.Raw(userQuery, startOfStartMonth.UnixMilli(), endOfCurrentMonth.UnixMilli()).Scan(&allUsers)

	// 解析时区用于数据分组
	loc, err := parseTimezone(timezone)
	if err != nil {
		loc = getRwandaTimezone()
	}

	// 在应用层按月份分组统计
	newUsersMap := make(map[string]MonthlyNewUsers)
	for _, user := range allUsers {
		// 将Unix时间戳转换为指定时区的月份
		userTime := time.Unix(user.CreatedAt/1000, (user.CreatedAt%1000)*1000000).In(loc)
		monthKey := fmt.Sprintf("%d-%02d", userTime.Year(), userTime.Month())

		data := newUsersMap[monthKey]
		data.Year = userTime.Year()
		data.Month = int(userTime.Month())
		data.NewUsers++
		if user.UserType == protocol.UserTypeDriver {
			data.NewDrivers++
		}
		newUsersMap[monthKey] = data
	}

	// 查询最新的总用户和总司机数
	var latestTotals struct {
		TotalUsers   int64 `json:"total_users"`
		TotalDrivers int64 `json:"total_drivers"`
	}

	s.db.Raw(`
		SELECT 
			COUNT(*) as total_users,
			SUM(CASE WHEN user_type = ? THEN 1 ELSE 0 END) as total_drivers
		FROM t_users
	`, protocol.UserTypeDriver).Scan(&latestTotals)

	// 查询从开始时间到现在的所有用户数据，在应用层计算累计总数
	var allUsersAfter []struct {
		CreatedAt int64  `json:"created_at"`
		UserType  string `json:"user_type"`
	}

	usersAfterQuery := `
		SELECT created_at, user_type
		FROM t_users
		WHERE created_at >= ?
		ORDER BY created_at ASC
	`

	s.db.Raw(usersAfterQuery, startOfStartMonth.UnixMilli()).Scan(&allUsersAfter)

	// 在应用层按月份分组统计
	monthToUsersAfter := make(map[string]int64)
	monthToDriversAfter := make(map[string]int64)

	// 初始化
	for monthKey := range monthMap {
		monthToUsersAfter[monthKey] = 0
		monthToDriversAfter[monthKey] = 0
	}

	// 解析时区用于数据分组
	loc, parseErr := parseTimezone(timezone)
	if parseErr != nil {
		loc = getRwandaTimezone()
	}

	// 计算每个月之后的新增用户总数
	for _, user := range allUsersAfter {
		// 将Unix时间戳转换为指定时区的月份
		userTime := time.Unix(user.CreatedAt/1000, (user.CreatedAt%1000)*1000000).In(loc)

		for monthKey, monthDate := range monthMap {
			// 如果用户注册日期在当前月份之后，加入计数
			if userTime.After(monthDate) {
				monthToUsersAfter[monthKey]++
				if user.UserType == protocol.UserTypeDriver {
					monthToDriversAfter[monthKey]++
				}
			}
		}
	}

	// 按顺序组装所有月份的数据
	for i := months - 1; i >= 0; i-- {
		date := now.AddDate(0, -i, 0)
		monthKey := fmt.Sprintf("%d-%02d", date.Year(), date.Month())
		formattedDate := date.Format("Jan")

		// 默认值
		newUsers := int64(0)
		newDrivers := int64(0)

		// 如果有该月新增数据，则使用
		if monthlyNew, ok := newUsersMap[monthKey]; ok {
			newUsers = monthlyNew.NewUsers
			newDrivers = monthlyNew.NewDrivers
		}

		// 计算累计总数 = 最新总数 - 该月之后的新增
		totalUsers := latestTotals.TotalUsers - monthToUsersAfter[monthKey]
		totalDrivers := latestTotals.TotalDrivers - monthToDriversAfter[monthKey]

		data = append(data, UserGrowthData{
			Date:         formattedDate,
			NewUsers:     newUsers,
			TotalUsers:   totalUsers,
			NewDrivers:   newDrivers,
			TotalDrivers: totalDrivers,
		})
	}

	return data
}
