package protocol

import "time"

// DirectionsResponse Google Directions API 响应结构
type DirectionsResponse struct {
	Status               string   `json:"status"`
	ErrorMessage         string   `json:"error_message,omitempty"`
	AvailableTravelModes []string `json:"available_travel_modes,omitempty"`
	GeocodedWaypoints    []struct {
		GeocoderStatus string   `json:"geocoder_status"`
		PlaceID        string   `json:"place_id"`
		Types          []string `json:"types"`
	} `json:"geocoded_waypoints,omitempty"`
	Routes []DirectionsRoute `json:"routes"`
}

// DirectionsRoute 路线信息
type DirectionsRoute struct {
	Summary          string           `json:"summary"`
	Legs             []DirectionsLeg  `json:"legs"`
	Bounds           Bounds           `json:"bounds"`
	Copyrights       string           `json:"copyrights"`
	Warnings         []string         `json:"warnings"`
	WaypointOrder    []int            `json:"waypoint_order"`
	OverviewPolyline OverviewPolyline `json:"overview_polyline"`
}

// DirectionsLeg 路线段信息
type DirectionsLeg struct {
	Distance          Distance         `json:"distance"`
	Duration          Duration         `json:"duration"`
	DurationInTraffic *Duration        `json:"duration_in_traffic,omitempty"`
	EndAddress        string           `json:"end_address"`
	EndLocation       LatLng           `json:"end_location"`
	StartAddress      string           `json:"start_address"`
	StartLocation     LatLng           `json:"start_location"`
	Steps             []DirectionsStep `json:"steps"`
}

// DirectionsStep 路线步骤信息
type DirectionsStep struct {
	Distance         Distance         `json:"distance"`
	Duration         Duration         `json:"duration"`
	EndLocation      LatLng           `json:"end_location"`
	HTMLInstructions string           `json:"html_instructions"`
	Polyline         OverviewPolyline `json:"polyline"`
	StartLocation    LatLng           `json:"start_location"`
	TravelMode       string           `json:"travel_mode"`
	Maneuver         string           `json:"maneuver,omitempty"`
}

// Distance 距离信息
type Distance struct {
	Text  string `json:"text"`
	Value int    `json:"value"` // 距离值，单位为米
}

// Duration 时长信息
type Duration struct {
	Text  string `json:"text"`
	Value int    `json:"value"` // 时长值，单位为秒
}

// LatLng 经纬度坐标
type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Bounds 地理边界
type Bounds struct {
	Northeast LatLng `json:"northeast"`
	Southwest LatLng `json:"southwest"`
}

// OverviewPolyline 路线编码
type OverviewPolyline struct {
	Points string `json:"points"`
}

// RouteRequest 路线计算请求参数
type RouteRequest struct {
	Origin      string   `json:"origin" validate:"required"`      // 起点坐标或地址
	Destination string   `json:"destination" validate:"required"` // 终点坐标或地址
	Mode        string   `json:"mode,omitempty"`                  // 出行方式: driving, walking, bicycling, transit
	Avoid       string   `json:"avoid,omitempty"`                 // 避开: tolls, highways, ferries, indoor
	Language    string   `json:"language,omitempty"`              // 语言设置
	Units       string   `json:"units,omitempty"`                 // 单位: metric, imperial
	Region      string   `json:"region,omitempty"`                // 地区设置
	Waypoints   []string `json:"waypoints,omitempty"`             // 途径点
}

// RouteResponse 路线计算响应
type RouteResponse struct {
	Success   bool             `json:"success"`
	Message   string           `json:"message,omitempty"`
	Route     *DirectionsRoute `json:"route,omitempty"`
	Distance  *Distance        `json:"distance,omitempty"`
	Duration  *Duration        `json:"duration,omitempty"`
	Error     *GoogleAPIError  `json:"error,omitempty"`
	RequestID string           `json:"request_id,omitempty"`
	Timestamp time.Time        `json:"timestamp"`
}

// GoogleAPIError Google API 错误信息
type GoogleAPIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// Validate 验证请求参数
func (r *RouteRequest) Validate() error {
	if r.Origin == "" {
		return NewValidationError("origin", "起点不能为空")
	}
	if r.Destination == "" {
		return NewValidationError("destination", "终点不能为空")
	}

	// 设置默认值
	if r.Mode == "" {
		r.Mode = "driving"
	}
	if r.Language == "" {
		r.Language = "zh-CN"
	}
	if r.Units == "" {
		r.Units = "metric"
	}

	return nil
}

// ValidationError 参数验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError 创建验证错误
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// Common Google API status codes
const (
	StatusOK                     = "OK"
	StatusNotFound               = "NOT_FOUND"
	StatusZeroResults            = "ZERO_RESULTS"
	StatusMaxWaypointsExceeded   = "MAX_WAYPOINTS_EXCEEDED"
	StatusMaxRouteLengthExceeded = "MAX_ROUTE_LENGTH_EXCEEDED"
	StatusInvalidRequest         = "INVALID_REQUEST"
	StatusOverDailyLimit         = "OVER_DAILY_LIMIT"
	StatusOverQueryLimit         = "OVER_QUERY_LIMIT"
	StatusRequestDenied          = "REQUEST_DENIED"
	StatusUnknownError           = "UNKNOWN_ERROR"
)

// IsSuccess 检查API响应是否成功
func (r *DirectionsResponse) IsSuccess() bool {
	return r.Status == StatusOK
}

// HasRoutes 检查是否有可用路线
func (r *DirectionsResponse) HasRoutes() bool {
	return len(r.Routes) > 0
}

// GetFirstRoute 获取第一条路线
func (r *DirectionsResponse) GetFirstRoute() *DirectionsRoute {
	if r.HasRoutes() {
		return &r.Routes[0]
	}
	return nil
}

// GetTotalDistance 获取总距离（米）
func (r *DirectionsResponse) GetTotalDistance() int {
	route := r.GetFirstRoute()
	if route != nil && len(route.Legs) > 0 {
		total := 0
		for _, leg := range route.Legs {
			total += leg.Distance.Value
		}
		return total
	}
	return 0
}

// GetTotalDuration 获取总时长（秒）
func (r *DirectionsResponse) GetTotalDuration() int {
	route := r.GetFirstRoute()
	if route != nil && len(route.Legs) > 0 {
		total := 0
		for _, leg := range route.Legs {
			total += leg.Duration.Value
		}
		return total
	}
	return 0
}
