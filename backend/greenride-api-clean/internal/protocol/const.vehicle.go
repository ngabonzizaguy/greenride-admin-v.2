package protocol

// =============================================================================
// 车辆相关常量定义
// =============================================================================

// 车辆分类常量
const (
	VehicleCategorySedan      = "sedan"      // 轿车
	VehicleCategorySUV        = "suv"        // 运动型多用途车
	VehicleCategoryHatchback  = "hatchback"  // 掀背车
	VehicleCategoryTruck      = "truck"      // 卡车
	VehicleCategoryMotorcycle = "motorcycle" // 摩托车
	VehicleCategoryVan        = "van"        // 面包车
	VehicleCategoryPickup     = "pickup"     // 皮卡
	VehicleCategoryConvert    = "convert"    // 敞篷车
	VehicleCategoryCoupe      = "coupe"      // 双门轿跑车
	VehicleCategoryWagon      = "wagon"      // 旅行车
)

// 车型级别常量
const (
	VehicleLevelEconomy = "economy" // 经济型
	VehicleLevelComfort = "comfort" // 舒适型
	VehicleLevelPremium = "premium" // 高级型
	VehicleLevelLuxury  = "luxury"  // 豪华型
)

// 燃料类型常量
const (
	FuelTypeGasoline = "gasoline" // 汽油
	FuelTypeDiesel   = "diesel"   // 柴油
	FuelTypeElectric = "electric" // 电动
	FuelTypeHybrid   = "hybrid"   // 混合动力
	FuelTypeCNG      = "cng"      // 压缩天然气
	FuelTypeLPG      = "lpg"      // 液化石油气
)

// 变速箱类型常量
const (
	TransmissionManual    = "manual"    // 手动挡
	TransmissionAutomatic = "automatic" // 自动挡
	TransmissionCVT       = "cvt"       // 无级变速
	TransmissionDual      = "dual"      // 双离合
)

// 需求/供给水平常量
const (
	LevelLow    = "low"    // 低水平
	LevelMedium = "medium" // 中等水平
	LevelHigh   = "high"   // 高水平
)

// 服务类型和车辆类型的映射关系
var ServiceTypeToVehicleCategorys = map[string][]string{
	ServiceTypeStandard: {VehicleCategorySedan, VehicleCategoryHatchback},
	ServiceTypePremium:  {VehicleCategorySUV, VehicleCategoryVan},
	ServiceTypeLuxury:   {VehicleCategorySedan, VehicleCategorySUV, VehicleCategoryCoupe},
}

// 获取服务类型对应的车辆类型列表
func GetVehicleCategorysForService(serviceType string) []string {
	if types, exists := ServiceTypeToVehicleCategorys[serviceType]; exists {
		return types
	}
	return []string{VehicleCategorySedan} // 默认返回轿车
}

// 验证车辆类型是否有效
func IsValidVehicleCategory(VehicleCategory string) bool {
	validTypes := []string{
		VehicleCategorySedan,
		VehicleCategorySUV,
		VehicleCategoryHatchback,
		VehicleCategoryTruck,
		VehicleCategoryMotorcycle,
		VehicleCategoryVan,
		VehicleCategoryPickup,
		VehicleCategoryConvert,
		VehicleCategoryCoupe,
		VehicleCategoryWagon,
	}

	for _, valid := range validTypes {
		if VehicleCategory == valid {
			return true
		}
	}
	return false
}

// 验证燃料类型是否有效
func IsValidFuelType(fuelType string) bool {
	validTypes := []string{
		FuelTypeGasoline,
		FuelTypeDiesel,
		FuelTypeElectric,
		FuelTypeHybrid,
		FuelTypeCNG,
		FuelTypeLPG,
	}

	for _, valid := range validTypes {
		if fuelType == valid {
			return true
		}
	}
	return false
}

// 验证变速箱类型是否有效
func IsValidTransmissionType(transmissionType string) bool {
	validTypes := []string{
		TransmissionManual,
		TransmissionAutomatic,
		TransmissionCVT,
		TransmissionDual,
	}

	for _, valid := range validTypes {
		if transmissionType == valid {
			return true
		}
	}
	return false
}

// 验证车辆状态是否有效
func IsValidVehicleStatus(status string) bool {
	validStatuses := []string{
		StatusActive,
		StatusInactive,
		StatusMaintenance,
		StatusRetired,
		StatusUnverified,
		StatusSuspended,
	}

	for _, valid := range validStatuses {
		if status == valid {
			return true
		}
	}
	return false
}

// 获取车辆类型的显示名称
func GetVehicleCategoryDisplayName(VehicleCategory string) string {
	displayNames := map[string]string{
		VehicleCategorySedan:      "轿车",
		VehicleCategorySUV:        "SUV",
		VehicleCategoryHatchback:  "掀背车",
		VehicleCategoryTruck:      "卡车",
		VehicleCategoryMotorcycle: "摩托车",
		VehicleCategoryVan:        "面包车",
		VehicleCategoryPickup:     "皮卡",
		VehicleCategoryConvert:    "敞篷车",
		VehicleCategoryCoupe:      "轿跑车",
		VehicleCategoryWagon:      "旅行车",
	}

	if name, exists := displayNames[VehicleCategory]; exists {
		return name
	}
	return VehicleCategory
}

// 获取服务类型的显示名称
func GetServiceTypeDisplayName(serviceType string) string {
	displayNames := map[string]string{
		ServiceTypeStandard: "标准",
		ServiceTypePremium:  "高级",
		ServiceTypeLuxury:   "豪华",
	}

	if name, exists := displayNames[serviceType]; exists {
		return name
	}
	return serviceType
}
