package utils

import "math"

// CalculateDistance 计算两点间距离（简化版本）
// lat1, lng1: 起点的纬度和经度
// lat2, lng2: 终点的纬度和经度
// 返回距离，单位：公里
func CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	// 简化的直线距离计算，实际应使用更精确的地理计算
	deltaLat := lat2 - lat1
	deltaLng := lng2 - lng1
	return math.Sqrt(deltaLat*deltaLat+deltaLng*deltaLng) * 111.0 // 近似转换为公里
}

// CalculateDistanceHaversine 使用 Haversine 公式计算两点间的球面距离
// lat1, lng1: 起点的纬度和经度
// lat2, lng2: 终点的纬度和经度
// 返回距离，单位：公里
func CalculateDistanceHaversine(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusKm = 6371.0

	// 将度数转换为弧度
	lat1Rad := lat1 * math.Pi / 180
	lng1Rad := lng1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lng2Rad := lng2 * math.Pi / 180

	// 计算差值
	deltaLat := lat2Rad - lat1Rad
	deltaLng := lng2Rad - lng1Rad

	// Haversine 公式
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

// CalculateCoordinateRange 根据中心坐标和半径计算坐标范围
// lat: 纬度, lng: 经度, radiusKm: 半径(公里)
// 返回: minLat, maxLat, minLng, maxLng
func CalculateCoordinateRange(lat, lng, radiusKm float64) (float64, float64, float64, float64) {
	// 地球半径(公里)
	const earthRadiusKm = 6371.0

	// 将半径转换为弧度
	radiusRad := radiusKm / earthRadiusKm

	// 将坐标转换为弧度
	latRad := lat * math.Pi / 180
	lngRad := lng * math.Pi / 180

	// 计算纬度范围
	minLatRad := latRad - radiusRad
	maxLatRad := latRad + radiusRad

	// 计算经度范围(考虑纬度对经度距离的影响)
	deltaLngRad := math.Asin(math.Sin(radiusRad) / math.Cos(latRad))
	minLngRad := lngRad - deltaLngRad
	maxLngRad := lngRad + deltaLngRad

	// 转换回度数
	minLat := minLatRad * 180 / math.Pi
	maxLat := maxLatRad * 180 / math.Pi
	minLng := minLngRad * 180 / math.Pi
	maxLng := maxLngRad * 180 / math.Pi

	return minLat, maxLat, minLng, maxLng
}
