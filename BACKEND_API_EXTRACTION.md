# GRD-WEBSITE-old Backend API Extraction

> **Extracted on:** December 30, 2025  
> **Purpose:** Complete API documentation for admin dashboard integration  
> **Source:** `GRD-WEBSITE-old/greenride-api/greenride-api-clean/`

---

## 1. Folder Structure

```
GRD-WEBSITE-old/
├── greenride-api/
│   └── greenride-api-clean/         # Main Go backend
│       ├── main/main.go             # Application entry point
│       ├── internal/
│       │   ├── config/              # Configuration management
│       │   ├── handlers/            # HTTP handlers (API & Admin)
│       │   ├── middleware/          # Auth, logging middleware
│       │   ├── models/              # Database models (GORM)
│       │   ├── protocol/            # Request/Response structs
│       │   ├── services/            # Business logic
│       │   ├── i18n/                # Internationalization
│       │   ├── locales/             # Translation files (en, fr, rw)
│       │   ├── log/                 # Logging utilities
│       │   ├── queue/               # Task queue
│       │   ├── task/                # Background tasks
│       │   └── utils/               # Utility functions
│       ├── docs/                    # Swagger documentation
│       ├── config.yaml              # Base configuration
│       ├── prod.yaml                # Production config
│       ├── dev.yaml                 # Development config
│       └── go.mod                   # Go module definition
├── greenride-admin/                 # Next.js Admin Dashboard (existing)
└── greenride-frontend/              # Public website
```

---

## 2. Go Backend Location

**Path:** `GRD-WEBSITE-old/greenride-api/greenride-api-clean/`

**Tech Stack:**
- **Go Version:** 1.23.0
- **Web Framework:** Gin (github.com/gin-gonic/gin)
- **ORM:** GORM (gorm.io/gorm)
- **Database:** MySQL
- **Cache:** Redis
- **Authentication:** JWT (github.com/golang-jwt/jwt/v5)
- **API Documentation:** Swagger (swaggo)

**Key Dependencies:**
```go
github.com/gin-gonic/gin v1.9.1
github.com/golang-jwt/jwt/v5 v5.2.2
gorm.io/driver/mysql v1.6.0
gorm.io/gorm v1.30.0
github.com/go-redis/redis/v8 v8.11.5
github.com/twilio/twilio-go v1.28.0
firebase.google.com/go v3.13.0
github.com/spf13/viper v1.16.0
```

---

## 3. All API Endpoints

### 3.1 Mobile API Service (Port 8610)

#### Public Endpoints (No Auth Required)

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | /health | inline | Health check |
| POST | /register | Register | User registration |
| POST | /login | Login | User login |
| POST | /send-verify-code | SendVerifyCode | Send OTP code |
| POST | /verify-code | VerifyCode | Verify OTP code |
| POST | /reset-password | ResetPassword | Reset password |
| POST | /feedback/submit | SubmitFeedback | Submit feedback |
| POST | /checkout/status | GetCheckoutStatus | Check payment status |
| POST | /webhook/kpay/:payment_id | KPayWebhook | Payment webhook |

#### Authenticated Endpoints (JWT Required)

**User Profile:**
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | /profile | Profile | Get user profile |
| POST | /logout | Logout | User logout |
| POST | /change-password | ChangePassword | Change password |
| POST | /profile/update | UpdateProfile | Update profile |
| POST | /profile/update/avatar | UpdateAvatar | Upload avatar |
| POST | /account/delete | DeleteAccount | Delete account |

**Driver Status:**
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | /online | UserOnline | Driver go online |
| POST | /offline | UserOffline | Driver go offline |

**Orders:**
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | /order/estimate | EstimateOrder | Estimate order price |
| POST | /order/create | CreateOrder | Create new order |
| POST | /orders | GetOrders | Get order list |
| POST | /order/detail | GetOrderDetail | Get order details |
| POST | /order/accept | AcceptOrder | Driver accept order |
| POST | /order/reject | RejectOrder | Driver reject order |
| POST | /order/arrived | ArrivedOrder | Driver arrived |
| POST | /order/start | StartOrder | Start trip |
| POST | /order/finish | FinishOrder | Finish trip |
| POST | /order/cancel | CancelOrder | Cancel order |
| POST | /order/rating | CreateOrderRating | Rate order |
| POST | /order/ratings | GetOrderRatings | Get ratings |
| POST | /order/nearby | GetNearbyOrders | Get nearby orders |
| POST | /order/cash/received | OrderCashReceived | Confirm cash payment |
| POST | /order/payment | OrderPayment | Process payment |

**Payment:**
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | /payment/methods | GetPaymentMethods | Get payment methods |
| POST | /payment/cancel | CancelPayment | Cancel payment |

**Vehicles:**
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | /vehicle | GetUserVehicle | Get user's vehicle |
| POST | /vehicles | GetVehicles | Get vehicle list |

**Location:**
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | /location/update | UpdateLocation | Update driver location |
| GET | /location/current | CurrentLocation | Get current location |
| GET | /drivers/nearby | GetNearbyDrivers | Get nearby drivers |

**Ratings:**
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | /rating/update | UpdateOrderRating | Update rating |
| POST | /rating/delete | DeleteOrderRating | Delete rating |
| POST | /rating/reply | ReplyToRating | Reply to rating |

**Promotions & Ads:**
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | /promotions | Promotions | Get promotions |
| POST | /ads/list | GetLocalAdvertisements | Get ads list |
| POST | /ads/detail | GetLocalAdvertisementByID | Get ad details |
| POST | /ads/stats | UpdateAdvertisementStats | Update ad stats |

---

### 3.2 Admin API Service (Port 8611)

#### Public Admin Endpoints

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | /health | inline | Health check |
| POST | /login | Login | Admin login |

#### Authenticated Admin Endpoints (JWT Required)

**Admin Auth:**
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | /logout | Logout | Admin logout |
| GET | /info | Info | Get admin info |
| POST | /change-password | ChangePassword | Change password |
| POST | /reset-password | ResetPassword | Reset password |

**Dashboard:**
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | /dashboard/stats | GetDashboardStats | Dashboard statistics |
| GET | /dashboard/revenue | GetRevenueChart | Revenue chart data |
| GET | /dashboard/user-growth | GetUserGrowthChart | User growth chart |

**User Management:**
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | /users/search | GetUserList | Search users |
| POST | /users/detail | GetUserDetail | Get user details |
| POST | /users/create | CreateUser | Create user |
| POST | /users/update | UpdateUser | Update user |
| POST | /users/status | UpdateUserStatus | Update status |
| POST | /users/verify | VerifyUser | Verify user |
| POST | /users/rides | GetUserRides | Get user rides |

**Vehicle Management:**
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | /vehicles/search | SearchVehicles | Search vehicles |
| POST | /vehicles/detail | GetVehicleDetail | Get vehicle details |
| POST | /vehicles/update | UpdateVehicle | Update vehicle |
| POST | /vehicles/status | UpdateVehicleStatus | Update status |
| POST | /vehicles/delete | DeleteVehicle | Delete vehicle |
| POST | /vehicles/create | CreateVehicle | Create vehicle |

**Order Management:**
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | /orders/search | SearchOrders | Search orders |
| POST | /orders/detail | GetOrderDetail | Get order details |
| POST | /orders/estimate | EstimateOrder | Estimate order |
| POST | /orders/create | CreateOrder | Create order |
| POST | /orders/cancel | CancelOrder | Cancel order |

---

## 4. Database Models

### 4.1 User Model (t_users)

```go
type User struct {
    ID     int64  `gorm:"column:id;primaryKey;autoIncrement"`
    UserID string `gorm:"column:user_id;type:varchar(64);uniqueIndex"`
    Salt   string `gorm:"column:salt;type:varchar(256)"`
    *UserValues
    CreatedAt int64 `gorm:"column:created_at;autoCreateTime:milli"`
}

type UserValues struct {
    UserType    *string `gorm:"column:user_type;type:varchar(32);index;default:'user'"` // passenger, driver
    Password    *string `gorm:"column:password;type:varchar(255)"`
    Email       *string `gorm:"column:email;type:varchar(255);index"`
    Phone       *string `gorm:"column:phone;type:varchar(64);index"`
    CountryCode *string `gorm:"column:country_code;type:varchar(10)"`
    Username    *string `gorm:"column:username;type:varchar(100);index"`
    DisplayName *string `gorm:"column:display_name;type:varchar(100)"`
    FirstName   *string `gorm:"column:first_name;type:varchar(100)"`
    LastName    *string `gorm:"column:last_name;type:varchar(100)"`
    Avatar      *string `gorm:"column:avatar;type:varchar(500)"`
    Gender      *string `gorm:"column:gender;type:varchar(10)"` // male, female, other
    Birthday    *int64  `gorm:"column:birthday"`
    Language    *string `gorm:"column:language;type:varchar(10);default:'en'"`
    Timezone    *string `gorm:"column:timezone;type:varchar(50);default:'UTC'"`
    
    // Address
    Address    *string `gorm:"column:address;type:text"`
    City       *string `gorm:"column:city;type:varchar(100)"`
    State      *string `gorm:"column:state;type:varchar(100)"`
    Country    *string `gorm:"column:country;type:varchar(100)"`
    PostalCode *string `gorm:"column:postal_code;type:varchar(20)"`
    
    // Location
    Latitude          *float64 `gorm:"column:latitude;type:decimal(10,8)"`
    Longitude         *float64 `gorm:"column:longitude;type:decimal(11,8)"`
    LocationUpdatedAt *int64   `gorm:"column:location_updated_at"`
    
    // Status
    Status          *string `gorm:"column:status;type:varchar(32);index;default:'active'"` // active, inactive, suspended, banned
    IsEmailVerified *bool   `gorm:"column:is_email_verified;default:false"`
    IsPhoneVerified *bool   `gorm:"column:is_phone_verified;default:false"`
    OnlineStatus    *string `gorm:"column:online_status;type:varchar(20);default:'offline'"` // online, offline, busy
    
    // Driver fields
    LicenseNumber    *string  `gorm:"column:license_number;type:varchar(50)"`
    LicenseExpiry    *int64   `gorm:"column:license_expiry"`
    QueuedOrderIds   []string `gorm:"column:queued_order_ids;type:json;serializer:json"`
    CurrentOrderId   *string  `gorm:"column:current_order_id;type:varchar(64)"`
    NextAvailableAt  *int64   `gorm:"column:next_available_at"`
    MaxQueueCapacity *int     `gorm:"column:max_queue_capacity;default:3"`
    
    // Stats
    Score      *float64 `gorm:"column:score;type:decimal(3,2);default:5.0"`
    TotalRides *int     `gorm:"column:total_rides;default:0"`
    
    // Referral
    InviteCode  *string `gorm:"column:invite_code;type:varchar(20);uniqueIndex"`
    InvitedBy   *string `gorm:"column:invited_by;type:varchar(64);index"`
    InviteCount *int    `gorm:"column:invite_count;default:0"`
    
    // Device
    FCMToken   *string `gorm:"column:fcm_token;type:varchar(500)"`
    DeviceID   *string `gorm:"column:device_id;type:varchar(255)"`
    DeviceType *string `gorm:"column:device_type;type:varchar(20)"` // ios, android, web
    
    // Timestamps
    LastLoginAt     *int64 `gorm:"column:last_login_at"`
    EmailVerifiedAt *int64 `gorm:"column:email_verified_at"`
    PhoneVerifiedAt *int64 `gorm:"column:phone_verified_at"`
    
    // Sandbox (0 = production, 1 = test user)
    Sandbox   *int   `gorm:"column:sandbox;type:tinyint(1);default:0"`
    UpdatedAt int64  `gorm:"column:updated_at;autoUpdateTime:milli"`
    DeletedAt *int64 `gorm:"column:deleted_at"` // Soft delete
}
```

### 4.2 Vehicle Model (t_vehicles)

```go
type Vehicle struct {
    ID        int64  `gorm:"column:id;primaryKey;autoIncrement"`
    VehicleID string `gorm:"column:vehicle_id;type:varchar(64);uniqueIndex"`
    Salt      string `gorm:"column:salt;type:varchar(256)"`
    *VehicleValues
    CreatedAt int64 `gorm:"column:created_at;autoCreateTime:milli"`
}

type VehicleValues struct {
    DriverID *string `gorm:"column:driver_id;type:varchar(64);index"`
    
    // Basic info
    Brand       *string `gorm:"column:brand;type:varchar(100)"`
    Model       *string `gorm:"column:model;type:varchar(100)"`
    Year        *int    `gorm:"column:year;type:int"`
    Color       *string `gorm:"column:color;type:varchar(50)"`
    PlateNumber *string `gorm:"column:plate_number;type:varchar(20);index"`
    VIN         *string `gorm:"column:vin;type:varchar(50);index"`
    
    // Type and specs
    TypeID       *string `gorm:"column:type_id;type:varchar(64);index"`     // Links to VehicleType
    Category     *string `gorm:"column:category;type:varchar(50);index"`    // sedan, suv, mpv, van, hatchback
    Level        *string `gorm:"column:level;type:varchar(50);index"`       // economy, comfort, premium, luxury
    SeatCapacity *int    `gorm:"column:seat_capacity;default:4"`
    FuelType     *string `gorm:"column:fuel_type;type:varchar(20)"`         // gasoline, diesel, electric, hybrid
    Transmission *string `gorm:"column:transmission;type:varchar(20)"`      // manual, automatic
    
    // Status
    Status       *string `gorm:"column:status;type:varchar(32);index;default:'active'"`      // active, inactive, maintenance, retired
    VerifyStatus *string `gorm:"column:verify_status;type:varchar(32);index;default:'verified'"`
    
    // Registration/Insurance
    RegistrationNumber    *string `gorm:"column:registration_number;type:varchar(50)"`
    RegistrationExpiry    *int64  `gorm:"column:registration_expiry"`
    InsuranceCompany      *string `gorm:"column:insurance_company;type:varchar(100)"`
    InsurancePolicyNumber *string `gorm:"column:insurance_policy_number;type:varchar(100)"`
    InsuranceExpiry       *int64  `gorm:"column:insurance_expiry"`
    
    // Location
    CurrentLatitude   *float64 `gorm:"column:current_latitude;type:decimal(10,8)"`
    CurrentLongitude  *float64 `gorm:"column:current_longitude;type:decimal(11,8)"`
    LocationUpdatedAt *int64   `gorm:"column:location_updated_at"`
    
    // Media
    Photos    []string `gorm:"column:photos;type:json;serializer:json"`
    Documents []string `gorm:"column:documents;type:json;serializer:json"`
    
    // Rating
    Rating *float64 `gorm:"column:rating;type:decimal(3,2);default:5.0"`
    
    UpdatedAt int64 `gorm:"column:updated_at;autoUpdateTime:milli"`
}
```

### 4.3 Order Model (t_orders)

```go
type Order struct {
    ID      int64  `gorm:"column:id;primaryKey;autoIncrement"`
    OrderID string `gorm:"column:order_id;type:varchar(64);uniqueIndex"`
    Salt    string `gorm:"column:salt;type:varchar(256)"`
    *OrderValues
    Details   *OrderDetail `gorm:"-"` // Not saved to DB
    CreatedAt int64        `gorm:"column:created_at;autoCreateTime:milli"`
}

type OrderValues struct {
    OrderType *string `gorm:"column:order_type;type:varchar(32);index"` // ride, delivery, shopping
    
    // User relations
    UserID     *string `gorm:"column:user_id;type:varchar(64);index"`     // Customer ID
    ProviderID *string `gorm:"column:provider_id;type:varchar(64);index"` // Driver ID
    
    // Status
    Status        *string `gorm:"column:status;type:varchar(32);index;default:'requested'"`
    PaymentStatus *string `gorm:"column:payment_status;type:varchar(32);index;default:''"`
    ScheduleType  *string `gorm:"column:schedule_type;type:varchar(32);default:'instant'"` // instant, scheduled
    
    // Amounts (decimal for precision)
    Currency            *string          `gorm:"column:currency;type:varchar(3);default:'USD'"`
    OriginalAmount      *decimal.Decimal `gorm:"column:original_amount;type:decimal(20,6)"`
    DiscountedAmount    *decimal.Decimal `gorm:"column:discounted_amount;type:decimal(20,6)"`
    PaymentAmount       *decimal.Decimal `gorm:"column:payment_amount;type:decimal(20,6)"`
    TotalDiscountAmount *decimal.Decimal `gorm:"column:total_discount_amount;type:decimal(20,6)"`
    PlatformFee         *decimal.Decimal `gorm:"column:platform_fee;type:decimal(20,6)"`
    
    // Payment
    PaymentMethod      *string `gorm:"column:payment_method;type:varchar(32)"` // card, cash, wallet
    PaymentID          *string `gorm:"column:payment_id;type:varchar(64)"`
    ChannelPaymentID   *string `gorm:"column:channel_payment_id;type:varchar(128)"`
    PaymentResult      *string `gorm:"column:payment_result;type:text"`
    PaymentRedirectURL *string `gorm:"column:payment_redirect_url;type:varchar(512)"`
    
    // Promotions
    PromoCodes        []string         `gorm:"column:promo_codes;type:json;serializer:json"`
    PromoDiscount     *decimal.Decimal `gorm:"column:promo_discount;type:decimal(20,6)"`
    UserPromotionIDs  []string         `gorm:"column:user_promotion_ids;type:json;serializer:json"`
    
    // Timestamps
    ScheduledAt *int64 `gorm:"column:scheduled_at"`
    AcceptedAt  *int64 `gorm:"column:accepted_at"`
    StartedAt   *int64 `gorm:"column:started_at"`
    EndedAt     *int64 `gorm:"column:ended_at"`
    CompletedAt *int64 `gorm:"column:completed_at"`
    CancelledAt *int64 `gorm:"column:cancelled_at"`
    ExpiredAt   *int64 `gorm:"column:expired_at"`
    
    // Cancellation
    CancelledBy     *string          `gorm:"column:cancelled_by;type:varchar(64)"`
    CancelReason    *string          `gorm:"column:cancel_reason;type:varchar(255)"`
    CancellationFee *decimal.Decimal `gorm:"column:cancellation_fee;type:decimal(20,6)"`
    
    // Dispatch
    DispatchStatus      *string `gorm:"column:dispatch_status;type:varchar(32);default:'not_started'"`
    CurrentRound        *int    `gorm:"column:current_round;default:0"`
    MaxRounds           *int    `gorm:"column:max_rounds;default:4"`
    AutoDispatchEnabled *bool   `gorm:"column:auto_dispatch_enabled;default:true"`
    
    Sandbox   *int   `gorm:"column:sandbox;type:tinyint(1);default:0"`
    Version   *int64 `gorm:"column:version;default:1"` // Optimistic locking
    UpdatedAt int64  `gorm:"column:updated_at;autoUpdateTime:milli"`
}
```

### 4.4 Admin Model (t_admins)

```go
type Admin struct {
    ID      int64  `gorm:"column:id;primaryKey;autoIncrement"`
    AdminID string `gorm:"column:admin_id;type:varchar(64);uniqueIndex"`
    Salt    string `gorm:"column:salt;type:varchar(256)"`
    *AdminValues
    CreatedAt int64 `gorm:"column:created_at;autoCreateTime:milli"`
}

type AdminValues struct {
    Username *string `gorm:"column:username;type:varchar(50);uniqueIndex"`
    Email    *string `gorm:"column:email;type:varchar(255);uniqueIndex"`
    Phone    *string `gorm:"column:phone;type:varchar(20);index"`
    
    FirstName *string `gorm:"column:first_name;type:varchar(100)"`
    LastName  *string `gorm:"column:last_name;type:varchar(100)"`
    FullName  *string `gorm:"column:full_name;type:varchar(200)"`
    Avatar    *string `gorm:"column:avatar;type:varchar(500)"`
    
    // Auth
    PasswordHash     *string `gorm:"column:password_hash;type:varchar(255)"`
    PasswordSalt     *string `gorm:"column:password_salt;type:varchar(255)"`
    TwoFactorEnabled *bool   `gorm:"column:two_factor_enabled;default:false"`
    
    // Role & Permissions
    Role        *string `gorm:"column:role;type:varchar(50);index"`  // super_admin, admin, moderator, support, analyst
    Permissions *string `gorm:"column:permissions;type:json"`        // JSON array
    Department  *string `gorm:"column:department;type:varchar(100)"`
    JobTitle    *string `gorm:"column:job_title;type:varchar(100)"`
    
    // Status
    Status       *string `gorm:"column:status;type:varchar(32);index;default:'active'"`
    ActiveStatus *string `gorm:"column:active_status;type:varchar(32);default:'offline'"`
    
    // Login tracking
    LastLoginAt    *int64  `gorm:"column:last_login_at"`
    LastLoginIP    *string `gorm:"column:last_login_ip;type:varchar(45)"`
    LoginCount     *int    `gorm:"column:login_count;default:0"`
    FailedAttempts *int    `gorm:"column:failed_attempts;default:0"`
    LockedUntil    *int64  `gorm:"column:locked_until"`
    
    UpdatedAt int64 `gorm:"column:updated_at;autoUpdateTime:milli"`
}

// Role constants
const (
    AdminRoleSuperAdmin = "super_admin"
    AdminRoleAdmin      = "admin"
    AdminRoleModerator  = "moderator"
    AdminRoleSupport    = "support"
    AdminRoleAnalyst    = "analyst"
)
```

---

## 5. Authentication System

### 5.1 JWT Configuration

```go
type JWTConfig struct {
    Secret     string        `mapstructure:"secret"`
    Expiration string        `mapstructure:"expiration"` // Default: "336h" (2 weeks)
    Issuer     string        `mapstructure:"issuer"`     // Default: "Greenride"
    ExpiresIn  time.Duration `mapstructure:"expires_in"`
}
```

### 5.2 JWT Claims Structure

```go
type JWTClaims struct {
    UserID   string `json:"user_id"`
    UserType string `json:"user_type"` // passenger, driver
    Username string `json:"username"`
    Role     string `json:"role"`
    Email    string `json:"email"`
    Phone    string `json:"phone"`
    jwt.RegisteredClaims
}
```

### 5.3 Token Extraction

```go
// Token can be provided via:
// 1. Authorization header: "Bearer <token>"
// 2. Query parameter: ?token=<token>

func GetTokenFromRequest(c *gin.Context) string {
    authHeader := c.GetHeader("Authorization")
    if authHeader != "" {
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString != authHeader {
            return tokenString
        }
    }
    return c.Query("token")
}
```

### 5.4 Auth Middleware (Mobile API)

```go
func (a *Api) AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := middleware.ValidToken(c, []byte(a.Jwt.Secret))
        if token == nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, protocol.NewAuthErrorResult())
            c.Abort()
            return
        }
        
        if claims, ok := token.Claims.(*middleware.JWTClaims); ok {
            user := services.GetUserService().GetUserByID(claims.UserID)
            if user == nil || user.GetStatus() != protocol.StatusActive {
                c.JSON(http.StatusUnauthorized, protocol.NewAuthErrorResult())
                c.Abort()
                return
            }
            
            c.Set("user", user)
            c.Set("user_id", claims.UserID)
            c.Set("user_type", claims.UserType)
        }
        c.Next()
    }
}
```

### 5.5 Admin Auth Middleware

```go
func (t *Admin) AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := middleware.ValidToken(c, []byte(t.Jwt.Secret))
        if token == nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, protocol.NewAuthErrorResult())
            c.Abort()
            return
        }
        
        if claims, ok := token.Claims.(*middleware.JWTClaims); ok {
            admin := services.GetAdminAdminService().GetAdminByID(claims.UserID)
            if admin == nil || admin.GetStatus() != protocol.StatusActive {
                c.JSON(http.StatusUnauthorized, protocol.NewAuthErrorResult())
                c.Abort()
                return
            }
            
            c.Set("user", admin)
            c.Set("user_id", claims.UserID)
        }
        c.Next()
    }
}
```

---

## 6. Request/Response Structures

### 6.1 Standard Response Format

```go
type Result struct {
    Code string `json:"code"`          // Error code (e.g., "0000" for success)
    Msg  string `json:"msg"`           // Human-readable message
    Data any    `json:"data,omitempty"` // Response payload
}

// Response codes
const (
    CODE_SUCCESS        = "0000"
    CODE_PARAM_ERROR    = "2001"
    CODE_AUTH_ERROR     = "3000"
    CODE_BUSINESS_ERROR = "1003"
    CODE_SYSTEM_ERROR   = "1000"
)
```

### 6.2 Pagination

```go
type Pagination struct {
    Size int `json:"size"` // Items per page
    Page int `json:"page"` // Current page (1-based)
}

type PageResult struct {
    ResultType string         `json:"result_type"`
    Size       int64          `json:"size"`
    Current    int64          `json:"current"`
    Total      int64          `json:"total"`   // Total pages
    Count      int64          `json:"count"`   // Total records
    Records    any            `json:"records"`
    Attach     map[string]any `json:"attach"`
}
```

### 6.3 Admin Request Examples

**Login Request:**
```go
type AdminLoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

// Response
{
    "code": "0000",
    "msg": "Success",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "user": {
            "admin_id": "ADM...",
            "username": "admin",
            "email": "admin@example.com",
            "role": "super_admin",
            ...
        }
    }
}
```

**Search Users Request:**
```go
type SearchRequest struct {
    Keyword         string   `json:"keyword,omitempty"`
    Page            int      `json:"page,omitempty"`      // Default: 1
    Limit           int      `json:"limit,omitempty"`     // Default: 10, Max: 100
    UserType        string   `json:"user_type,omitempty"` // passenger, driver
    Status          string   `json:"status,omitempty"`
    OnlineStatus    string   `json:"online_status,omitempty"`
    IsEmailVerified *bool    `json:"is_email_verified,omitempty"`
    IsPhoneVerified *bool    `json:"is_phone_verified,omitempty"`
    IsActive        *bool    `json:"is_active,omitempty"`
}
```

**Search Orders Request:**
```go
type OrderSearchRequest struct {
    Keyword       string   `json:"keyword,omitempty"`
    Page          int      `json:"page,omitempty"`
    Limit         int      `json:"limit,omitempty"`
    OrderID       string   `json:"order_id,omitempty"`
    OrderType     string   `json:"order_type,omitempty"`
    Status        string   `json:"status,omitempty"`
    PaymentStatus string   `json:"payment_status,omitempty"`
    UserID        string   `json:"user_id,omitempty"`
    ProviderID    string   `json:"provider_id,omitempty"`
    StartDate     *int64   `json:"start_date,omitempty"`  // Timestamp in milliseconds
    EndDate       *int64   `json:"end_date,omitempty"`
    MinAmount     *float64 `json:"min_amount,omitempty"`
    MaxAmount     *float64 `json:"max_amount,omitempty"`
}
```

**Vehicle Create Request:**
```go
type VehicleCreateRequest struct {
    DriverID     string `json:"driver_id,omitempty"`
    Brand        string `json:"brand" binding:"required"`
    Model        string `json:"model" binding:"required"`
    PlateNumber  string `json:"plate_number" binding:"required"`
    Year         *int   `json:"year,omitempty"`
    Color        *string `json:"color,omitempty"`
    Category     *string `json:"category,omitempty"`     // sedan, suv, mpv, van
    Level        *string `json:"level,omitempty"`        // economy, comfort, premium, luxury
    SeatCapacity *int    `json:"seat_capacity,omitempty"`
}
```

---

## 7. Configuration

### 7.1 config.yaml Structure

```yaml
env: prod  # or dev
debug: false

server:
  api:
    name: "greenride-api"
    port: "8610"
    version: "1.0"
    jwt:
      secret: "your-jwt-secret"
      expiration: "336h"
  admin:
    name: "greenride-admin"
    port: "8611"
    version: "1.0"
    jwt:
      secret: "your-jwt-secret"
      expiration: "336h"

database:
  dsn: "user:password@tcp(host:3306)/greenride?charset=utf8mb4&parseTime=True&loc=Local"
  max_idle_conns: 5
  max_open_conns: 25
  conn_max_lifetime: 300s

redis:
  addr: "localhost:6379"
  password: ""
  db: 0

log:
  level: "info"
  path: "/app/logs"
  format: "json"

i18n:
  locales_dir: "/app/internal/locales"
  default_language: "en"
  supported_langs: ["en", "fr", "rw"]
```

### 7.2 Production Config (prod.yaml)

```yaml
server:
  api_port: 8610
  admin_port: 8611
  ws_port: 8612
  jwt_secret: "grd_prod_secret_2025_9a64"

database:
  dsn: "greenride:GreenRide2024!@tcp(db:3306)/greenride?charset=utf8mb4&parseTime=True&loc=Local"

redis:
  addr: "redis:6379"
  password: ""
  db: 0

log:
  level: "info"
  path: "/app/logs"
  format: "json"

env: production
debug: false
```

---

## 8. Database Setup

### 8.1 Database Type
- **MySQL** (using `gorm.io/driver/mysql`)

### 8.2 Connection Setup

```go
func InitDB(cfg *config.Config) error {
    dsn := cfg.Database.DSN
    if dsn == "" {
        return fmt.Errorf("database DSN is empty")
    }
    
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logLevel),
        DisableForeignKeyConstraintWhenMigrating: true,
    })
    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }
    
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
    sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
    sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
    
    DB = db
    return nil
}
```

### 8.3 Auto-Migration (Tables)

The system auto-migrates these tables:
- t_users
- t_user_accounts
- t_user_addresses
- t_user_payment_methods
- t_user_promotions
- t_user_location_history
- t_admins
- t_identities
- t_orders
- t_order_ratings
- t_order_history_logs
- t_ride_orders
- t_dispatch_records
- t_vehicles
- t_vehicle_types
- t_price_rules
- t_price_snapshots
- t_payments
- t_payment_methods
- t_payment_channels
- t_wallets
- t_wallet_transactions
- t_withdrawals
- t_promotions
- t_messages
- t_message_templates
- t_notifications
- t_fcm_tokens
- t_fcm_message_logs
- t_service_areas
- t_announcements
- t_feedbacks
- t_tasks

---

## 9. Existing Admin Dashboard

**Location:** `GRD-WEBSITE-old/greenride-admin/greenride-admin/`

**Tech Stack:** Next.js with TypeScript, Tailwind CSS

**Key Service Files:**
- `src/services/api.ts` - Base API client
- `src/services/auth.ts` - Authentication
- `src/services/users.ts` - User management
- `src/services/drivers.ts` - Driver management
- `src/services/vehicles.ts` - Vehicle management
- `src/services/orders.ts` - Order management
- `src/services/dashboard.ts` - Dashboard stats

---

## 10. Constants & Status Values

### Order Status
```go
const (
    StatusRequested   = "requested"    // Order created, waiting for driver
    StatusAccepted    = "accepted"     // Driver accepted
    StatusArrived     = "arrived"      // Driver arrived at pickup
    StatusInProgress  = "in_progress"  // Trip started
    StatusTripEnded   = "trip_ended"   // Trip finished, awaiting payment
    StatusCompleted   = "completed"    // Fully completed
    StatusCancelled   = "cancelled"    // Cancelled
)
```

### Payment Status
```go
const (
    StatusPending = "pending"
    StatusSuccess = "success"
    StatusFailed  = "failed"
)
```

### User Status
```go
const (
    StatusActive   = "active"
    StatusInactive = "inactive"
    StatusSuspended = "suspended"
    StatusBanned   = "banned"
)
```

### Online Status
```go
const (
    StatusOnline  = "online"
    StatusOffline = "offline"
    StatusBusy    = "busy"
)
```

### User Types
```go
const (
    UserTypePassenger = "passenger"
    UserTypeDriver    = "driver"
)
```

### Payment Methods
```go
const (
    PaymentMethodCash   = "cash"
    PaymentMethodCard   = "card"
    PaymentMethodWallet = "wallet"
    PaymentMethodMomo   = "momo"  // Mobile money
)
```

---

## 11. API Usage Examples

### Admin Login
```bash
curl -X POST http://localhost:8611/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password123"}'
```

### Get Dashboard Stats (with auth)
```bash
curl -X GET http://localhost:8611/dashboard/stats \
  -H "Authorization: Bearer <token>"
```

### Search Users
```bash
curl -X POST http://localhost:8611/users/search \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "keyword": "john",
    "page": 1,
    "limit": 20,
    "user_type": "driver",
    "status": "active"
  }'
```

### Create Vehicle
```bash
curl -X POST http://localhost:8611/vehicles/create \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "driver_id": "USR_xxxx",
    "brand": "Toyota",
    "model": "Corolla",
    "plate_number": "RAB 123 A",
    "year": 2022,
    "color": "White",
    "category": "sedan",
    "level": "economy",
    "seat_capacity": 4
  }'
```

---

## 12. Important Notes

1. **All timestamps are in milliseconds** (Unix epoch * 1000)
2. **Soft deletes** - Users have a `deleted_at` field, not physically deleted
3. **Optimistic locking** - Orders use a `version` field for concurrency
4. **Sandbox mode** - Users/orders can be marked as sandbox for testing
5. **Two separate services** - API (8610) for mobile, Admin (8611) for dashboard
6. **Same JWT secret** can be shared between services if needed
7. **CORS** - Handled by nginx in production, dev mode has permissive CORS

---

*Document generated for admin dashboard integration purposes.*

