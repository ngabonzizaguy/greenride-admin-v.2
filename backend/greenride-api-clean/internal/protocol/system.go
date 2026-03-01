package protocol

// SystemConfigResponse 系统配置公共响应（移动端/前端读取）
type SystemConfigResponse struct {
	MaintenanceMode     bool   `json:"maintenance_mode"`
	MaintenanceMessage  string `json:"maintenance_message,omitempty"`
	MaintenancePhone    string `json:"maintenance_phone,omitempty"`
	MaintenanceStartAt  int64  `json:"maintenance_started_at,omitempty"`
	UpdateNoticeEnabled bool   `json:"update_notice_enabled"`
	UpdateNoticeTitle   string `json:"update_notice_title,omitempty"`
	UpdateNoticeMessage string `json:"update_notice_message,omitempty"`
	ForceUpdateEnabled  bool   `json:"force_update_enabled"`
	MinimumAppVersion   string `json:"minimum_app_version,omitempty"`
	LatestAppVersion    string `json:"latest_app_version,omitempty"`
	AndroidStoreURL     string `json:"android_store_url,omitempty"`
	IOSStoreURL         string `json:"ios_store_url,omitempty"`
}

// SystemConfigUpdateRequest 管理员更新系统配置请求
type SystemConfigUpdateRequest struct {
	MaintenanceMode     *bool   `json:"maintenance_mode"`
	MaintenanceMessage  *string `json:"maintenance_message"`
	MaintenancePhone    *string `json:"maintenance_phone"`
	UpdateNoticeEnabled *bool   `json:"update_notice_enabled"`
	UpdateNoticeTitle   *string `json:"update_notice_title"`
	UpdateNoticeMessage *string `json:"update_notice_message"`
	ForceUpdateEnabled  *bool   `json:"force_update_enabled"`
	MinimumAppVersion   *string `json:"minimum_app_version"`
	LatestAppVersion    *string `json:"latest_app_version"`
	AndroidStoreURL     *string `json:"android_store_url"`
	IOSStoreURL         *string `json:"ios_store_url"`
}
