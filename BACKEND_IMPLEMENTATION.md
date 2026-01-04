# Backend Implementation Guide

Since you are the backend developer, you can apply these changes directly to your Go codebase (`greenride-api-clean`).

## 1. Database Updates (MySQL)

Run these SQL commands to prepare your database.

```sql
-- 1. Create Support Configuration Table
CREATE TABLE IF NOT EXISTS `t_support_config` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `config_key` varchar(50) NOT NULL UNIQUE, -- e.g., 'default'
  `phone` varchar(50) DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL,
  `whatsapp` varchar(50) DEFAULT NULL,
  `operating_hours` varchar(100) DEFAULT NULL,
  `faq_url` varchar(255) DEFAULT NULL,
  `updated_at` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Initialize default config if empty
INSERT IGNORE INTO `t_support_config` (`config_key`, `phone`, `email`, `whatsapp`, `operating_hours`, `faq_url`, `updated_at`)
VALUES ('default', '+250 788 000 000', 'support@greenride.com', '+250 788 000 000', 'Mon-Sun, 24/7', 'https://greenride.com/faq', UNIX_TIMESTAMP()*1000);

-- 2. Update Feedback Table (if missing columns)
-- Check if 'admin_response' exists, if not add it:
ALTER TABLE `t_feedbacks` ADD COLUMN `admin_response` TEXT DEFAULT NULL;
ALTER TABLE `t_feedbacks` ADD COLUMN `status` VARCHAR(20) DEFAULT 'pending'; -- pending, reviewing, resolved
ALTER TABLE `t_feedbacks` ADD COLUMN `category` VARCHAR(50) DEFAULT 'general';
ALTER TABLE `t_feedbacks` ADD COLUMN `severity` VARCHAR(20) DEFAULT 'low';
ALTER TABLE `t_feedbacks` ADD COLUMN `resolved_at` bigint(20) DEFAULT NULL;
```

---

## 2. Go Models (`internal/models/`)

Create or update `internal/models/support.go`:

```go
package models

type SupportConfig struct {
    ID             int64   `gorm:"column:id;primaryKey;autoIncrement"`
    ConfigKey      string  `gorm:"column:config_key;type:varchar(50);unique"`
    Phone          *string `gorm:"column:phone;type:varchar(50)"`
    Email          *string `gorm:"column:email;type:varchar(100)"`
    Whatsapp       *string `gorm:"column:whatsapp;type:varchar(50)"`
    OperatingHours *string `gorm:"column:operating_hours;type:varchar(100)"`
    FaqURL         *string `gorm:"column:faq_url;type:varchar(255)"`
    UpdatedAt      int64   `gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (SupportConfig) TableName() string {
    return "t_support_config"
}
```

Update `internal/models/feedback.go`:

```go
package models

type Feedback struct {
    ID            int64   `gorm:"column:id;primaryKey;autoIncrement"`
    FeedbackID    string  `gorm:"column:feedback_id;type:varchar(64);uniqueIndex"`
    UserID        string  `gorm:"column:user_id;type:varchar(64);index"`
    OrderID       *string `gorm:"column:order_id;type:varchar(64);index"`
    Content       string  `gorm:"column:content;type:text"`
    
    // New Fields
    Status        string  `gorm:"column:status;default:'pending'"`
    Category      string  `gorm:"column:category;default:'general'"`
    Severity      string  `gorm:"column:severity;default:'low'"`
    AdminResponse *string `gorm:"column:admin_response;type:text"`
    ResolvedAt    *int64  `gorm:"column:resolved_at"`
    
    CreatedAt     int64   `gorm:"column:created_at;autoCreateTime:milli"`
}

func (Feedback) TableName() string {
    return "t_feedbacks"
}
```

---

## 3. Protocol Structs (`internal/protocol/`)

Add these to a file like `internal/protocol/admin_extra.go`:

```go
package protocol

// Support Config Request
type UpdateSupportConfigRequest struct {
    Phone          string `json:"phone"`
    Email          string `json:"email"`
    Whatsapp       string `json:"whatsapp"`
    OperatingHours string `json:"operating_hours"`
    FaqURL         string `json:"faq_url"`
}

// Feedback Search Request
type FeedbackSearchRequest struct {
    Page      int    `json:"page"`
    Limit     int    `json:"limit"`
    Status    string `json:"status"`
    Category  string `json:"category"`
    StartDate *int64 `json:"start_date"`
    EndDate   *int64 `json:"end_date"`
}

// Feedback Update Request
type FeedbackUpdateRequest struct {
    FeedbackID    string `json:"feedback_id" binding:"required"`
    Status        string `json:"status"`
    AdminResponse string `json:"admin_response"`
    Severity      string `json:"severity"`
}
```

---

## 4. Handlers (`internal/handlers/admin_extra.go`)

Create this file to handle the logic.

```go
package handlers

import (
    "net/http"
    "time"
    "greenride-api/internal/models"
    "greenride-api/internal/protocol"
    "github.com/gin-gonic/gin"
)

// --- Support Config Handlers ---

func GetSupportConfig(c *gin.Context) {
    var config models.SupportConfig
    // Assuming 'default' key
    if err := DB.Where("config_key = ?", "default").First(&config).Error; err != nil {
        c.JSON(http.StatusOK, protocol.Result{Code: "1000", Msg: "Config not found"})
        return
    }
    c.JSON(http.StatusOK, protocol.Result{Code: "0000", Msg: "Success", Data: config})
}

func UpdateSupportConfig(c *gin.Context) {
    var req protocol.UpdateSupportConfigRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusOK, protocol.Result{Code: "2001", Msg: "Invalid params"})
        return
    }

    var config models.SupportConfig
    if err := DB.Where("config_key = ?", "default").First(&config).Error; err != nil {
        // Create if not exists
        config.ConfigKey = "default"
    }

    // Update fields
    config.Phone = &req.Phone
    config.Email = &req.Email
    config.Whatsapp = &req.Whatsapp
    config.OperatingHours = &req.OperatingHours
    config.FaqURL = &req.FaqURL
    config.UpdatedAt = time.Now().UnixMilli()

    DB.Save(&config)
    c.JSON(http.StatusOK, protocol.Result{Code: "0000", Msg: "Updated successfully"})
}

// --- Feedback Handlers ---

func SearchFeedback(c *gin.Context) {
    var req protocol.FeedbackSearchRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        req.Page = 1
        req.Limit = 10
    }
    if req.Page < 1 { req.Page = 1 }
    if req.Limit < 1 { req.Limit = 10 }

    query := DB.Model(&models.Feedback{})
    if req.Status != "" && req.Status != "all" {
        query = query.Where("status = ?", req.Status)
    }
    if req.Category != "" && req.Category != "all" {
        query = query.Where("category = ?", req.Category)
    }

    var total int64
    query.Count(&total)

    var feedbacks []models.Feedback
    offset := (req.Page - 1) * req.Limit
    query.Order("created_at DESC").Offset(offset).Limit(req.Limit).Find(&feedbacks)

    c.JSON(http.StatusOK, protocol.PageResult{
        ResultType: "feedback",
        Current:    int64(req.Page),
        Size:       int64(req.Limit),
        Total:      (total + int64(req.Limit) - 1) / int64(req.Limit),
        Count:      total,
        Records:    feedbacks,
    })
}

func UpdateFeedback(c *gin.Context) {
    var req protocol.FeedbackUpdateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusOK, protocol.Result{Code: "2001", Msg: "Invalid params"})
        return
    }

    var feedback models.Feedback
    if err := DB.Where("feedback_id = ?", req.FeedbackID).First(&feedback).Error; err != nil {
        c.JSON(http.StatusOK, protocol.Result{Code: "1000", Msg: "Feedback not found"})
        return
    }

    if req.Status != "" { feedback.Status = req.Status }
    if req.AdminResponse != "" { feedback.AdminResponse = &req.AdminResponse }
    if req.Severity != "" { feedback.Severity = req.Severity }
    
    if req.Status == "resolved" && feedback.ResolvedAt == nil {
        now := time.Now().UnixMilli()
        feedback.ResolvedAt = &now
    }

    DB.Save(&feedback)
    c.JSON(http.StatusOK, protocol.Result{Code: "0000", Msg: "Updated successfully"})
}
```

---

## 5. Register Routes (`main.go` or `router.go`)

Add these lines to your Admin API Group registration:

```go
adminGroup := r.Group("/admin/v1")
adminGroup.Use(AuthMiddleware()) // Ensure protected
{
    // ... existing routes ...

    // Support Config
    adminGroup.GET("/config/support", handlers.GetSupportConfig)
    adminGroup.POST("/config/support", handlers.UpdateSupportConfig)

    // Feedback
    adminGroup.POST("/feedback/search", handlers.SearchFeedback)
    adminGroup.POST("/feedback/update", handlers.UpdateFeedback)
}
```

## 6. Deployment

1.  Stop the backend service.
2.  Run the SQL commands.
3.  Rebuild and start the Go service.
4.  In the Admin Dashboard (`src/lib/api-client.ts`), set `DEMO_MODE = false`.



