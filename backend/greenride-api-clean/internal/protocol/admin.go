package protocol

// Admin 管理员信息V2
type Admin struct {
	AdminID      string `json:"admin_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	Role         string `json:"role"`
	Department   string `json:"department"`
	Status       string `json:"status"`
	ActiveStatus string `json:"active_status"`
	CreatedAt    int64  `json:"created_at"`
	LastLoginAt  *int64 `json:"last_login_at,omitempty"`
}
