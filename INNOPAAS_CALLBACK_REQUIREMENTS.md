# 📱 InnoPaaS Callback Requirements - From Configuration Images

## 🔍 **Key Information Extracted from InnoPaaS Dashboard**

Based on the configuration screenshots, here are the critical requirements for completing the OTP/SMS service setup:

---

## 1. **Callback URL Configuration** ✅

### **Configured Callback Address:**
```
https://api.greenrideafrica.com/v1/sms/callback
```

### **Critical Requirements:**
- ✅ **Must return HTTP 200 status code within 3 seconds**
- ✅ **No response message body required** (just status code)
- ✅ **If callback fails validation, InnoPaaS will mark connectivity as abnormal**

**Important Note:** The callback URL will be validated when saving the configuration in InnoPaaS dashboard.

---

## 2. **Events Configuration** ✅

### **Selected Event:**
- ✅ **"message status updated"** - Checked and enabled

**What this means:**
- InnoPaaS will send POST requests to your callback URL when SMS message status changes
- You'll receive notifications for:
  - Message sent
  - Message delivered
  - Message failed
  - Message expired
  - etc.

---

## 3. **Request Mode** ✅

### **Selected Method:**
- ✅ **POST** (selected in configuration)

**What this means:**
- InnoPaaS will send POST requests (not GET)
- Request body will contain message status information
- Content-Type will likely be `application/json`

---

## 4. **Callback Authentication (Optional but Recommended)** 🔐

### **Authentication Method:**
InnoPaaS supports optional authentication using the `X-CALLBACK-ID` header.

### **Header Format:**
```
X-CALLBACK-ID: timestamp={timestamp};nonce={nonce};username={username};signature={signature}
```

### **Signature Calculation:**
```
signature = HMAC-SHA256(secret, timestamp+nonce+username)
```

**Where:**
- `timestamp` = Unix timestamp of the callback message
- `nonce` = Random number
- `username` = Username configured in InnoPaaS dashboard (optional)
- `secret` = Your authorization key (from `prod.yaml`: `authorization: "3u1K73"`)

### **Configuration Fields:**
- **Username:** Optional (if set, key is required)
- **Authorization:** Optional (authorization information)

**Note:** If you configure username/authorization in InnoPaaS dashboard, you MUST verify the signature on your backend.

---

## 5. **Authorization Header (Optional)** 🔑

### **From Configuration:**
- Authorization field is available for callback requests
- If your callback endpoint requires authentication, provide it here
- InnoPaaS will include this Authorization header when sending requests

**Current Config:**
- Your `prod.yaml` has: `authorization: "3u1K73"`
- This might be used for:
  1. Signature verification (HMAC secret)
  2. Authorization header in callback requests

---

## 📋 **Implementation Checklist**

### **Backend Requirements:**

- [ ] **Create callback endpoint:** `POST /v1/sms/callback`
  - **Note:** Current API uses `/` as base path, may need to add `/v1` prefix or use `/sms/callback`
  
- [ ] **Return HTTP 200 within 3 seconds:**
  - Process must be fast
  - Consider async processing for heavy operations
  - Return 200 immediately, process in background if needed

- [ ] **Handle POST requests:**
  - Accept `application/json` content type
  - Parse message status data from request body

- [ ] **Optional: Verify X-CALLBACK-ID signature:**
  - Extract header: `X-CALLBACK-ID`
  - Parse: `timestamp`, `nonce`, `username`, `signature`
  - Calculate: `HMAC-SHA256(secret, timestamp+nonce+username)`
  - Compare with received signature
  - Reject if signature doesn't match

- [ ] **Process message status updates:**
  - Update SMS delivery status in database
  - Log status changes
  - Handle different status types (sent, delivered, failed, etc.)

- [ ] **Error handling:**
  - Return 200 even on errors (to avoid InnoPaaS retries)
  - Log errors for debugging
  - Don't expose internal errors in response

---

## 🔧 **Implementation Guide**

### **Step 1: Create Callback Endpoint**

**File:** `backend/greenride-api-clean/internal/handlers/api.sms_callback.go`

```go
package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"greenride/internal/config"
	"greenride/internal/log"
	"greenride/internal/protocol"

	"github.com/gin-gonic/gin"
)

// InnoPaaSCallbackRequest represents the callback payload from InnoPaaS
type InnoPaaSCallbackRequest struct {
	MessageID string `json:"message_id"`
	To        string `json:"to"`
	Status    string `json:"status"`    // sent, delivered, failed, expired, etc.
	ErrorCode string `json:"error_code,omitempty"`
	ErrorMsg  string `json:"error_msg,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

// InnoPaaSCallback handles InnoPaaS SMS status callback
// @Summary InnoPaaS SMS Callback
// @Description Receives SMS delivery status updates from InnoPaaS
// @Tags Api,Webhook,SMS
// @Accept json
// @Produce json
// @Param callback_data body InnoPaaSCallbackRequest true "Callback data"
// @Success 200 "OK"
// @Router /sms/callback [post]
func (a *Api) InnoPaaSCallback(c *gin.Context) {
	// Parse request body
	var callbackData InnoPaaSCallbackRequest
	if err := c.ShouldBindJSON(&callbackData); err != nil {
		log.Get().Errorf("InnoPaaS Callback: invalid JSON: %v", err)
		c.Status(http.StatusOK) // Return 200 even on error
		return
	}

	// Optional: Verify X-CALLBACK-ID signature
	cfg := config.Get().InnoPaaS
	if cfg != nil && cfg.Authorization != "" {
		if !verifyCallbackSignature(c, cfg.Authorization) {
			log.Get().Warnf("InnoPaaS Callback: invalid signature for message_id=%s", callbackData.MessageID)
			// Still return 200 to avoid retries, but log the issue
		}
	}

	// Log callback received
	log.Get().Infof("InnoPaaS Callback received: message_id=%s, to=%s, status=%s",
		callbackData.MessageID, callbackData.To, callbackData.Status)

	// Process status update asynchronously (to ensure fast response)
	go func() {
		processSMSStatusUpdate(callbackData)
	}()

	// Return 200 immediately (within 3 seconds requirement)
	c.Status(http.StatusOK)
}

// verifyCallbackSignature verifies the X-CALLBACK-ID header signature
func verifyCallbackSignature(c *gin.Context, secret string) bool {
	callbackID := c.GetHeader("X-CALLBACK-ID")
	if callbackID == "" {
		return false
	}

	// Parse header: timestamp={timestamp};nonce={nonce};username={username};signature={signature}
	parts := strings.Split(callbackID, ";")
	var timestamp, nonce, username, receivedSignature string

	for _, part := range parts {
		if strings.HasPrefix(part, "timestamp=") {
			timestamp = strings.TrimPrefix(part, "timestamp=")
		} else if strings.HasPrefix(part, "nonce=") {
			nonce = strings.TrimPrefix(part, "nonce=")
		} else if strings.HasPrefix(part, "username=") {
			username = strings.TrimPrefix(part, "username=")
		} else if strings.HasPrefix(part, "signature=") {
			receivedSignature = strings.TrimPrefix(part, "signature=")
		}
	}

	if timestamp == "" || nonce == "" || receivedSignature == "" {
		return false
	}

	// Calculate expected signature: HMAC-SHA256(secret, timestamp+nonce+username)
	message := timestamp + nonce
	if username != "" {
		message += username
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	// Compare signatures (constant-time comparison)
	return hmac.Equal([]byte(receivedSignature), []byte(expectedSignature))
}

// processSMSStatusUpdate processes the SMS status update
func processSMSStatusUpdate(data InnoPaaSCallbackRequest) {
	// TODO: Update SMS status in database
	// TODO: Log status change
	// TODO: Handle different status types
	
	log.Get().Infof("Processing SMS status update: message_id=%s, status=%s, to=%s",
		data.MessageID, data.Status, data.To)

	// Example: Update verification code send status
	// This would depend on your database schema
	// You might want to:
	// 1. Find the verification code record by phone number and timestamp
	// 2. Update its delivery status
	// 3. Log the status change
}
```

---

### **Step 2: Register Endpoint**

**File:** `backend/greenride-api-clean/internal/handlers/api.go`

Add to the public routes section:

```go
api := router.Group("/")
{
    // ... existing routes ...
    
    // InnoPaaS SMS Callback - No authentication required
    api.POST("/sms/callback", a.InnoPaaSCallback)
    
    // OR if you need /v1 prefix:
    // v1 := api.Group("/v1")
    // v1.POST("/sms/callback", a.InnoPaaSCallback)
}
```

---

### **Step 3: Update Nginx Configuration**

If using `/v1/sms/callback`, ensure Nginx routes it correctly:

```nginx
location /v1/sms/callback {
    proxy_pass http://greenride_api_backend;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    
    # Preserve X-CALLBACK-ID header for signature verification
    proxy_set_header X-CALLBACK-ID $http_x_callback_id;
    
    # Fast timeout to meet 3-second requirement
    proxy_read_timeout 3s;
    proxy_connect_timeout 2s;
}
```

---

## 🧪 **Testing the Callback**

### **Test 1: Basic Callback (No Signature)**

```bash
curl -X POST "http://localhost:8610/sms/callback" \
  -H "Content-Type: application/json" \
  -d '{
    "message_id": "test_msg_123",
    "to": "250784928786",
    "status": "delivered",
    "timestamp": 1234567890
  }'
```

**Expected:** HTTP 200 (no body)

---

### **Test 2: Callback with Signature Verification**

```bash
# Calculate signature first
TIMESTAMP=$(date +%s)
NONCE="random123"
USERNAME="testuser"
SECRET="3u1K73"
MESSAGE="${TIMESTAMP}${NONCE}${USERNAME}"

# Calculate HMAC-SHA256 (using openssl)
SIGNATURE=$(echo -n "$MESSAGE" | openssl dgst -sha256 -hmac "$SECRET" | cut -d' ' -f2)

curl -X POST "http://localhost:8610/sms/callback" \
  -H "Content-Type: application/json" \
  -H "X-CALLBACK-ID: timestamp=${TIMESTAMP};nonce=${NONCE};username=${USERNAME};signature=${SIGNATURE}" \
  -d '{
    "message_id": "test_msg_123",
    "to": "250784928786",
    "status": "delivered"
  }'
```

---

## 📊 **Status Types to Handle**

Based on typical SMS provider statuses, you should handle:

- `sent` - Message sent to provider
- `delivered` - Message delivered to recipient
- `failed` - Message failed to send
- `expired` - Message expired
- `rejected` - Message rejected by provider
- `pending` - Message pending

---

## ⚠️ **Important Notes**

1. **3-Second Requirement:**
   - Your endpoint MUST return 200 within 3 seconds
   - Use async processing (goroutines) for heavy operations
   - Don't do database writes synchronously if they're slow

2. **No Response Body:**
   - InnoPaaS doesn't require a response body
   - Just return HTTP 200 status code

3. **Error Handling:**
   - Always return 200, even on errors
   - Log errors for debugging
   - Don't let errors cause timeouts

4. **Signature Verification:**
   - Optional but recommended for security
   - If configured in InnoPaaS, you MUST verify it
   - Use constant-time comparison to prevent timing attacks

5. **Callback URL:**
   - Currently configured as: `https://api.greenrideafrica.com/v1/sms/callback`
   - Ensure this matches your actual endpoint path
   - Update Nginx routing if needed

---

## ✅ **Complete Setup Checklist**

- [ ] Create callback endpoint handler
- [ ] Register endpoint in router
- [ ] Implement signature verification (optional)
- [ ] Add async processing for status updates
- [ ] Update Nginx configuration (if needed)
- [ ] Test callback with sample data
- [ ] Verify 200 response within 3 seconds
- [ ] Test signature verification (if enabled)
- [ ] Update database schema for SMS status tracking (if needed)
- [ ] Deploy and verify in production

---

## 🔗 **Related Files**

- `backend/greenride-api-clean/internal/handlers/api.go` - Router setup
- `backend/greenride-api-clean/internal/services/innopaas.go` - InnoPaaS service
- `backend/greenride-api-clean/prod.yaml` - InnoPaaS configuration
- `nginx.conf` - Nginx routing (if using reverse proxy)

---

## 📝 **Summary**

From the InnoPaaS configuration images, we learned:

1. ✅ **Callback URL:** `https://api.greenrideafrica.com/v1/sms/callback`
2. ✅ **Method:** POST
3. ✅ **Response:** HTTP 200 within 3 seconds (no body required)
4. ✅ **Events:** "message status updated"
5. ✅ **Authentication:** Optional X-CALLBACK-ID header with HMAC-SHA256 signature
6. ✅ **Authorization:** Optional authorization header support

**Next Steps:**
1. Implement the callback endpoint
2. Test locally
3. Deploy to production
4. Verify InnoPaaS can reach the endpoint

This completes the OTP/SMS service setup! 🚀
