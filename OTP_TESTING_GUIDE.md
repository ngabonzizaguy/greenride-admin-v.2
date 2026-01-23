# 🔐 OTP Testing Guide

## 📍 **Where is OTP Handled?**

**Answer: BACKEND** ✅

The OTP flow is entirely backend-driven:

1. **App** → Calls `POST /send-verify-code` API endpoint
2. **Backend** → Generates OTP code, stores in Redis cache, sends SMS via InnoPaaS/Twilio
3. **App** → User enters code, calls `POST /verify-code` API endpoint
4. **Backend** → Validates code from Redis cache

**The app does NOT generate or validate OTP codes** - it only:
- Requests OTP via API
- Displays input field for user to enter code
- Submits code for verification via API

---

## 🐛 **Current Issue: OTP Not Received Until Timeout**

### **Root Cause:**

Your config has `bypass_otp: true` enabled, which means:

```yaml
verify_code:
  bypass_otp: true   # ⚠️ THIS IS THE PROBLEM
```

**What happens:**
1. ✅ Backend generates OTP code (or uses "1234" in sandbox mode)
2. ✅ Code is stored in Redis cache
3. ❌ **SMS is NOT sent** (bypass mode skips actual SMS sending)
4. ⏳ App waits for SMS that never arrives
5. ⏰ Eventually times out

**Code Logic:**
```go
// If bypass_otp is true, SMS is NOT sent
if s.config.BypassOTP && contactType == protocol.MsgChannelSms {
    isSandbox = true  // Forces sandbox mode
}

if isSandbox {
    code = "1234"  // Uses fixed code
}

// ⚠️ SMS is only sent if NOT sandbox
if !isSandbox {
    // Send SMS via InnoPaaS/Twilio
    s.msgService.SendMessage(msg)
}
```

---

## ✅ **Solution: Disable Bypass Mode**

### **Option 1: Disable for Testing (Recommended)**

**File:** `backend/greenride-api-clean/local.yaml` (or `dev.yaml`)

```yaml
verify_code:
  length: 4
  expiration: 5
  send_interval: 60
  max_send_times: 10
  bypass_otp: false   # ✅ Change to false
```

**Then restart backend:**
```bash
cd backend/greenride-api-clean
go run main/main.go
```

---

### **Option 2: Use Sandbox Mode with Known Code**

If you want to test without sending real SMS:

1. **Keep `bypass_otp: true`**
2. **Use code "1234" for all OTPs** (this is what sandbox mode uses)
3. **Tell users to enter "1234"** when testing

**Note:** This is only for development/testing, not production!

---

## 🧪 **Testing OTP Locally (Terminal)**

### **Step 1: Start Backend**

```bash
cd backend/greenride-api-clean
go run main/main.go
```

**Expected output:**
```
Loaded environment specific config: dev.yaml  # or local.yaml
Server starting on port 8610
```

---

### **Step 2: Send OTP (Request Verification Code)**

**Using curl:**

```bash
curl -X POST "http://localhost:8610/send-verify-code" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "register",
    "phone": "+250784928786",
    "user_type": "passenger",
    "language": "en"
  }'
```

**Expected Response (Success):**
```json
{
  "code": "0000",
  "msg": "Success",
  "data": null
}
```

**Expected Response (Error - Cooldown):**
```json
{
  "code": "E0015",
  "msg": "Verification code sent too frequently",
  "data": {
    "remaining_seconds": "45"
  }
}
```

---

### **Step 3: Check Backend Logs**

**Look for:**
```
Sent verification code 1234 to +250784928786 for purpose register and user type passenger
```

**If bypass_otp is false, also look for:**
```
SMS sent via InnoPaaS, MessageID: <message_id>
```
or
```
SMS sent via Twilio with SID: <sid>
```

**If there's an error:**
```
Failed to send verification code to +250784928786: <error>
```

---

### **Step 4: Verify OTP Code**

**Using curl:**

```bash
curl -X POST "http://localhost:8610/verify-code" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "register",
    "method": "phone",
    "phone": "+250784928786",
    "code": "1234",
    "user_type": "passenger"
  }'
```

**Expected Response (Success):**
```json
{
  "code": "0000",
  "msg": "Success",
  "data": {
    "message": "Verification successful",
    "type": "register",
    "method": "phone"
  }
}
```

**Expected Response (Invalid Code):**
```json
{
  "code": "E0016",
  "msg": "Invalid verification code"
}
```

---

## 🔍 **Debugging: Check What's Happening**

### **1. Check if Bypass is Enabled**

**Look in backend logs at startup:**
```
Loaded environment specific config: dev.yaml
```

Then check `dev.yaml` or `local.yaml`:
```yaml
verify_code:
  bypass_otp: true   # ⚠️ If true, SMS won't be sent
```

---

### **2. Check SMS Service Configuration**

**Check which SMS provider is configured:**

```bash
# In dev.yaml or local.yaml
grep -A 5 "sms:" backend/greenride-api-clean/dev.yaml
```

**Output:**
```yaml
sms:
  service_name: "twilio"  # or "innopaas"
```

---

### **3. Check Redis Cache (OTP Storage)**

**If you have Redis CLI access:**

```bash
redis-cli
> KEYS *verify_code*
> GET sms_verify_code_register_passenger_+250784928786
```

**Expected output:**
```
"1234"  # or actual OTP code
```

---

### **4. Test SMS Provider Directly**

**Test InnoPaaS API:**

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

**Test Twilio API:**

```bash
curl -X POST "https://verify.twilio.com/v2/Services/VAd9a138c03d20b4ca69fde9c17b57d039/Verifications" \
  -u "AC67e59ab2092ad2e5d9839c16ff1e1372:46e32b8b761e924156ca87eb96e7ae4d" \
  -d "To=+250784928786" \
  -d "Channel=sms" \
  -d "CustomCode=1234"
```

---

## 📋 **Complete Test Flow**

### **Test 1: Send OTP (Bypass Mode - Code "1234")**

```bash
# 1. Send OTP request
curl -X POST "http://localhost:8610/send-verify-code" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "register",
    "phone": "+250784928786",
    "user_type": "passenger"
  }'

# 2. Check logs - should see:
# "Sent verification code 1234 to +250784928786"
# NO SMS sent (bypass mode)

# 3. Verify with code "1234"
curl -X POST "http://localhost:8610/verify-code" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "register",
    "method": "phone",
    "phone": "+250784928786",
    "code": "1234",
    "user_type": "passenger"
  }'
```

---

### **Test 2: Send OTP (Real SMS Mode)**

**First, disable bypass:**

```yaml
# local.yaml or dev.yaml
verify_code:
  bypass_otp: false  # ✅ Disable bypass
```

**Restart backend, then:**

```bash
# 1. Send OTP request
curl -X POST "http://localhost:8610/send-verify-code" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "register",
    "phone": "+250784928786",
    "user_type": "passenger"
  }'

# 2. Check logs - should see:
# "Sent verification code <random_code> to +250784928786"
# "SMS sent via InnoPaaS, MessageID: <id>"  # or Twilio

# 3. Check phone for SMS with actual code
# 4. Verify with received code
curl -X POST "http://localhost:8610/verify-code" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "register",
    "method": "phone",
    "phone": "+250784928786",
    "code": "<code_from_sms>",
    "user_type": "passenger"
  }'
```

---

## 🎯 **Quick Fix for Your Issue**

### **To Fix "No OTP Received Until Timeout":**

1. **Disable bypass mode:**
   ```yaml
   # backend/greenride-api-clean/local.yaml
   verify_code:
     bypass_otp: false  # Change from true to false
   ```

2. **Restart backend:**
   ```bash
   cd backend/greenride-api-clean
   go run main/main.go
   ```

3. **Test again:**
   - Request OTP from app
   - Check backend logs for SMS send confirmation
   - Check phone for SMS

4. **If still not working:**
   - Check SMS provider credentials (InnoPaaS/Twilio)
   - Check backend logs for errors
   - Verify phone number format

---

## 📊 **Summary**

| Setting | Behavior |
|---------|----------|
| `bypass_otp: true` | ✅ Code stored in cache<br>❌ SMS NOT sent<br>✅ Use code "1234" for testing |
| `bypass_otp: false` | ✅ Code stored in cache<br>✅ SMS sent via provider<br>✅ Use code from SMS |

**Your current issue:** `bypass_otp: true` means SMS is never sent, so users wait forever for a code that never arrives.

**Solution:** Set `bypass_otp: false` to actually send SMS, or use code "1234" if you want to test without sending SMS.

---

## 🔧 **Environment-Specific Configs**

**Production (`prod.yaml`):**
```yaml
verify_code:
  bypass_otp: true   # ⚠️ Currently enabled (SMS not sent)
```

**Development (`dev.yaml`):**
```yaml
verify_code:
  bypass_otp: true   # ⚠️ Currently enabled (SMS not sent)
```

**Local (`local.yaml`):**
```yaml
verify_code:
  bypass_otp: true   # ⚠️ Currently enabled (SMS not sent)
```

**Recommendation:** 
- **Production:** `bypass_otp: false` (send real SMS)
- **Development:** `bypass_otp: false` (test real SMS) or `true` (use "1234")
- **Local:** `bypass_otp: true` (use "1234" for quick testing)

---

## ✅ **Testing Checklist**

- [ ] Backend is running
- [ ] `bypass_otp` setting checked
- [ ] SMS provider configured (InnoPaaS or Twilio)
- [ ] OTP request sent successfully
- [ ] Backend logs show code generation
- [ ] Backend logs show SMS send (if bypass disabled)
- [ ] SMS received on phone (if bypass disabled)
- [ ] OTP verification works with correct code
- [ ] OTP verification fails with wrong code

---

## 🚨 **Common Issues**

### **Issue 1: "No SMS Received"**

**Causes:**
- `bypass_otp: true` (SMS not sent)
- SMS provider credentials invalid
- Phone number format incorrect
- SMS provider API error

**Fix:**
- Check `bypass_otp` setting
- Check backend logs for errors
- Verify SMS provider credentials
- Test SMS provider API directly

---

### **Issue 2: "Invalid Code"**

**Causes:**
- Code expired (default 5 minutes)
- Wrong code entered
- Code already used
- Cache issue

**Fix:**
- Request new code
- Check code expiration time
- Verify code format (4-6 digits)
- Check Redis cache

---

### **Issue 3: "Cooldown Error"**

**Causes:**
- Requested OTP too soon (default 60 seconds between requests)

**Fix:**
- Wait for cooldown period
- Check `send_interval` in config
- Use different phone number for testing

---

## 📞 **Need Help?**

If OTP still doesn't work:

1. **Check backend logs** for errors
2. **Verify SMS provider** credentials
3. **Test SMS API** directly
4. **Check Redis** for stored codes
5. **Verify phone number** format

The most common issue is `bypass_otp: true` preventing SMS from being sent! 🔧
