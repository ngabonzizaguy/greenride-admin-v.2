package protocol

// SystemConfigResponse 系统配置公共响应（移动端/前端读取）
type SystemConfigResponse struct {
	MaintenanceMode    bool   `json:"maintenance_mode"`
	MaintenanceMessage string `json:"maintenance_message,omitempty"`
	MaintenancePhone   string `json:"maintenance_phone,omitempty"`
	MaintenanceStartAt int64  `json:"maintenance_started_at,omitempty"`
}

// SystemConfigUpdateRequest 管理员更新系统配置请求
type SystemConfigUpdateRequest struct {
	MaintenanceMode    *bool   `json:"maintenance_mode"`
	MaintenanceMessage *string `json:"maintenance_message"`
	MaintenancePhone   *string `json:"maintenance_phone"`
}
