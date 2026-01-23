# 📱 InnoPaaS OTP/SMS Integration Setup Guide

## ✅ **Current Status**

The InnoPaaS integration is **already implemented** in the backend! Here's what's done and what needs verification:

---

## 🔍 **What's Already Done**

### **1. Backend Implementation** ✅
- **Service:** `backend/greenride-api-clean/internal/services/innopaas.go`
- **Integration:** Fully implemented with proper error handling
- **Code Extraction:** Fixed to correctly extract OTP code from message params

### **2. Configuration** ✅
- **Production:** Already configured in `prod.yaml`
- **Service Selection:** SMS service switches between Twilio and InnoPaaS based on config

### **3. Current Production Config** (`prod.yaml`):
```yaml
sms:
  service_name: "innopaas"  # ✅ Already set to InnoPaaS

innopaas:
  endpoint: "https://api.innopaas.com/api/otp/v3/msg/send/verify"
  app_key: "KaunIJesBYOVrM29"
  authorization: "3u1K73"
  sender_id: "GreenRide"
```

---

## 🔧 **Setup & Verification Steps**

### **Step 1: Verify InnoPaaS Credentials**

Check if the credentials in `prod.yaml` match what InnoPaaS provided:

1. **Endpoint:** Should be `https://api.innopaas.com/api/otp/v3/msg/send/verify`
2. **App Key:** Verify `app_key` matches InnoPaaS dashboard
3. **Authorization:** Verify `authorization` (secret) matches InnoPaaS dashboard
4. **Sender ID:** Optional, but should match approved sender ID

**Location:** `backend/greenride-api-clean/prod.yaml` (lines 154-158)

---

### **Step 2: Update Other Environments (if needed)**

**Development/Staging Configs:**

Currently, `dev.yaml` and `local.yaml` use Twilio. If you want to test InnoPaaS in dev:

**File:** `backend/greenride-api-clean/dev.yaml`
```yaml
# Change from:
sms:
  service_name: "twilio"

# To:
sms:
  service_name: "innopaas"

# Add InnoPaaS config:
innopaas:
  endpoint: "https://api.innopaas.com/api/otp/v3/msg/send/verify"
  app_key: "YOUR_DEV_APP_KEY"  # Use dev credentials if available
  authorization: "YOUR_DEV_AUTHORIZATION"
  sender_id: "GreenRide"
```

---

### **Step 3: Test the Integration**

#### **A. Test via Backend API**

1. **Start the backend:**
   ```bash
   cd backend/greenride-api-clean
   go run main/main.go
   ```

2. **Send a test OTP:**
   ```bash
   curl -X POST "http://localhost:8610/send-verify-code" \
     -H "Content-Type: application/json" \
     -d '{
       "contact_type": "sms",
       "contact": "+250784928786",
       "user_type": "passenger",
       "purpose": "register",
       "language": "en"
     }'
   ```

3. **Expected Response:**
   ```json
   {
     "code": "0000",
     "msg": "Success",
     "data": null
   }
   ```

4. **Check Backend Logs:**
   Look for: `SMS sent via InnoPaaS, MessageID: <message_id>`

#### **B. Test via Mobile App**

1. **Register/Login Flow:**
   - Enter phone number (e.g., `+250784928786`)
   - Request OTP
   - Verify SMS is received

2. **Check Logs:**
   - Backend should log successful InnoPaaS API calls
   - Check for any error messages

---

### **Step 4: Verify API Request Format**

The backend sends requests in this format:

**Request:**
```json
POST https://api.innopaas.com/api/otp/v3/msg/send/verify
Headers:
  Content-Type: application/json
  Authorization: 3u1K73
  appKey: KaunIJesBYOVrM29

Body:
{
  "to": "250784928786",  // Phone number (digits only, no +)
  "type": "3",            // OTP type
  "language": "en",       // Language code
  "code": "1234"          // OTP code
}
```

**Expected Response:**
```json
{
  "code": "0000",
  "success": true,
  "message": "Success",
  "data": "<message_id>"
}
```

---

## 🐛 **Troubleshooting**

### **Issue 1: OTP Not Received**

**Check:**
1. ✅ Backend logs for InnoPaaS API errors
2. ✅ Phone number format (should be digits only: `250784928786`)
3. ✅ InnoPaaS dashboard for delivery status
4. ✅ Credentials match InnoPaaS dashboard

**Common Errors:**
- `InnoPaaS configuration is missing` → Check `innopaas` section in config
- `InnoPaaS API returned error` → Check credentials and API endpoint
- `failed to extract OTP code` → Should be fixed in latest code

---

### **Issue 2: Wrong Phone Number Format**

**Problem:** Phone number includes `+` or other characters

**Solution:** The code automatically cleans phone numbers:
```go
// Removes + and non-digits
cleanTo := cleanPhoneNumber(to)  // "+250784928786" → "250784928786"
```

---

### **Issue 3: Code Extraction Issues**

**Problem:** OTP code not extracted correctly

**Solution:** Fixed in latest code - now extracts from `Params["code"]` directly:
```go
// Uses code from params (preferred)
code := message.Params["code"]

// Falls back to extracting from content if needed
code = extractCodeFromContent(content)
```

---

## 📋 **Code Changes Made**

### **Fixed: OTP Code Extraction**

**File:** `backend/greenride-api-clean/internal/services/innopaas.go`

**Changes:**
1. ✅ Now extracts OTP code from `Params["code"]` (direct from VerifyCodeService)
2. ✅ Falls back to extracting from formatted content if needed
3. ✅ Added `extractCodeFromContent()` helper function
4. ✅ Better error messages

**Before:**
```go
// Assumed content was the code
Code: content,
```

**After:**
```go
// Extracts code from params or content
var code string
if codeVal, ok := message.Params["code"]; ok && codeVal != nil {
    code = fmt.Sprintf("%v", codeVal)
} else {
    code = extractCodeFromContent(content)
}
```

---

## 🔐 **Security Notes**

1. **Credentials:** Keep `app_key` and `authorization` secure
2. **Environment Variables:** Consider moving sensitive data to env vars
3. **Rate Limiting:** InnoPaaS may have rate limits - check their docs
4. **Error Handling:** Backend logs errors but doesn't expose sensitive info

---

## 📊 **Monitoring**

### **What to Monitor:**

1. **Success Rate:**
   - Check backend logs for `SMS sent via InnoPaaS`
   - Monitor InnoPaaS dashboard for delivery stats

2. **Error Rate:**
   - Watch for `InnoPaaS API returned error` in logs
   - Track failed OTP sends

3. **Response Times:**
   - InnoPaaS API timeout is set to 10 seconds
   - Monitor for slow responses

---

## ✅ **Verification Checklist**

- [ ] InnoPaaS credentials verified in `prod.yaml`
- [ ] Backend restarted with new config
- [ ] Test OTP sent successfully
- [ ] SMS received on test phone (+250 number)
- [ ] Backend logs show successful InnoPaaS calls
- [ ] Error handling works (test with invalid credentials)
- [ ] Phone number formatting correct (digits only)
- [ ] Admin dashboard tweaks completed

---

## 🚀 **Next Steps**

1. **Test with Rwanda Operators:**
   - Test with +250 numbers
   - Verify delivery to different operators
   - Check delivery times

2. **Admin Dashboard:**
   - Complete minor tweaks
   - Add SMS provider status indicator (optional)
   - Add OTP send logs view (optional)

3. **Production Deployment:**
   - Deploy backend with InnoPaaS config
   - Monitor for 24-48 hours
   - Verify no issues

---

## 📞 **Support**

If you encounter issues:

1. **Check Backend Logs:**
   ```bash
   # Look for InnoPaaS-related errors
   grep -i "innopaas" backend/logs/*.log
   ```

2. **Test API Directly:**
   ```bash
   curl -X POST "https://api.innopaas.com/api/otp/v3/msg/send/verify" \
     -H "Content-Type: application/json" \
     -H "Authorization: 3u1K73" \
     -H "appKey: KaunIJesBYOVrM29" \
     -d '{
       "to": "250784928786",
       "type": "3",
       "language": "en",
       "code": "1234"
     }'
   ```

3. **Contact InnoPaaS Support:**
   - Check their API documentation
   - Verify account status
   - Check for API changes

---

## 📝 **Summary**

✅ **Integration Status:** Fully implemented and ready for testing  
✅ **Code Quality:** Fixed OTP extraction, proper error handling  
✅ **Configuration:** Production config ready  
⏳ **Pending:** Testing and verification with Rwanda operators  

The integration should work out of the box once credentials are verified and backend is restarted! 🚀
