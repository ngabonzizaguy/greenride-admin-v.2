# 📱 Feedback Submission Flow Guide

> **Purpose:** Understand what happens when a user submits feedback from the mobile app  
> **Flow:** Mobile App → Backend API → Database → Admin Dashboard

---

## 🎯 **Complete Feedback Flow**

### **1. User Submits Feedback (Mobile App)**

**What Happens:**
1. User fills out feedback form in mobile app
2. User provides:
   - Title (required)
   - Content/Message (required)
   - Email (required - for contact)
   - Optional: Category, Rating, Attachments, etc.

**API Call:**
```
POST http://18.143.118.157:8610/feedback/submit
Content-Type: application/json

{
  "title": "App crash when booking ride",
  "content": "The app crashes whenever I try to book a ride. It happens on Android version 12.",
  "email": "user@example.com"
}
```

---

### **2. Backend Processes Feedback**

**What Happens on Backend:**

#### **Step 1: Validation**
- ✅ Validates request format (JSON)
- ✅ Validates email format (must contain "@")
- ❌ If invalid → Returns `400 Bad Request`

#### **Step 2: Rate Limiting**
- ✅ Checks IP-based rate limit (1 submission per minute per IP)
- ✅ Uses Redis cache: `feedback:ratelimit:{IP_ADDRESS}`
- ❌ If rate limit exceeded → Returns `429 Too Many Requests`

#### **Step 3: Create Feedback Record**
- ✅ Creates feedback in database (`t_feedbacks` table)
- ✅ Sets default values:
  - `feedback_id`: Auto-generated (e.g., `FDB_xxxxx`)
  - `status`: `"pending"` (default)
  - `feedback_type`: `"suggestion"` (default)
  - `category`: `"other"` (default)
  - `severity`: `"medium"` (default)
  - `contact_email`: User's email
  - `created_at`: Current timestamp

#### **Step 4: Set Rate Limit Cache**
- ✅ Sets Redis cache for 1 minute to prevent spam

#### **Step 5: Return Response**
- ✅ Returns `200 OK` with feedback ID:
```json
{
  "code": "0000",
  "msg": "Success",
  "data": {
    "feedback_id": "FDB_abc123xyz"
  }
}
```

---

### **3. Feedback Stored in Database**

**Database Table:** `t_feedbacks`

**Record Created:**
```sql
INSERT INTO t_feedbacks (
  feedback_id,
  title,
  content,
  contact_email,
  feedback_type,
  category,
  status,
  severity,
  priority,
  created_at,
  updated_at
) VALUES (
  'FDB_abc123xyz',
  'App crash when booking ride',
  'The app crashes whenever I try to book a ride...',
  'user@example.com',
  'suggestion',
  'other',
  'pending',
  'medium',
  'medium',
  1705123456789,
  1705123456789
);
```

---

### **4. Feedback Appears in Admin Dashboard**

**When:** Immediately after submission (if admin dashboard is open and refreshing)

**Where:** Admin Dashboard → Feedback & Complaints page (`/feedback`)

**What Admin Sees:**

#### **In Feedback List:**
- ✅ New feedback appears at the top of the list
- ✅ Status badge: **"Pending"** (yellow badge)
- ✅ Category badge: **"Other"** (or user-selected category)
- ✅ User info: Email (from `contact_email` field)
- ✅ Title and preview of content
- ✅ Timestamp: "Just now" or "X minutes ago"

#### **In Feedback Stats:**
- ✅ Total count increases by 1
- ✅ Pending count increases by 1
- ✅ Stats cards update automatically

---

## 📊 **What to Expect in Admin Dashboard**

### **1. Feedback List View**

**What You'll See:**
```
┌─────────────────────────────────────────────────────────┐
│ Feedback & Complaints                                    │
├─────────────────────────────────────────────────────────┤
│ Total: 14 | Pending: 3 | Reviewing: 1 | Resolved: 10   │
├─────────────────────────────────────────────────────────┤
│ [Search] [Filter: All Categories] [Filter: All Status]  │
├─────────────────────────────────────────────────────────┤
│ 📱 App crash when booking ride                          │
│    user@example.com • Other • Pending • 2m ago         │
│    [View Details] [Actions ▼]                           │
├─────────────────────────────────────────────────────────┤
│ ... (other feedback)                                    │
└─────────────────────────────────────────────────────────┘
```

### **2. Feedback Detail View**

**When You Click "View Details":**
- ✅ Full feedback content
- ✅ User information (email, phone if available)
- ✅ Status and category
- ✅ Created timestamp
- ✅ Actions: Update Status, Add Response, Delete

---

## 🔄 **Feedback Status Lifecycle**

### **Default Status:** `pending`

**Status Flow:**
```
pending → reviewing → resolved
         ↓
      cancelled
```

**What Admin Can Do:**
1. **View Feedback** → See full details
2. **Update Status** → Change to "reviewing", "resolved", or "cancelled"
3. **Add Response** → Provide admin response/notes
4. **Assign to Admin** → Assign feedback to specific admin
5. **Delete Feedback** → Remove feedback (with confirmation)

---

## ⚠️ **Important Notes**

### **1. Rate Limiting**
- **Limit:** 1 submission per minute per IP address
- **Cache Duration:** 1 minute
- **If Exceeded:** User sees `429 Too Many Requests` error

### **2. No Authentication Required**
- ✅ Feedback submission does **NOT** require login
- ✅ Anyone can submit feedback
- ✅ Email is used for contact (not for authentication)

### **3. Default Values**
- **Status:** Always starts as `"pending"`
- **Category:** Defaults to `"other"` (unless user specifies)
- **Feedback Type:** Defaults to `"suggestion"` (unless user specifies)
- **Severity:** Defaults to `"medium"`

### **4. Data Storage**
- ✅ Feedback stored in `t_feedbacks` table
- ✅ User email stored in `contact_email` field
- ✅ Timestamps: `created_at` and `updated_at` are auto-set
- ✅ Feedback ID: Auto-generated unique identifier

---

## ✅ **Expected Behavior Summary**

### **When User Submits Feedback:**

1. **Mobile App:**
   - ✅ Shows success message
   - ✅ Returns feedback ID
   - ✅ User can submit feedback (if rate limit not exceeded)

2. **Backend:**
   - ✅ Validates input
   - ✅ Rate limits (1 per minute per IP)
   - ✅ Creates feedback record
   - ✅ Returns feedback ID

3. **Database:**
   - ✅ New record created in `t_feedbacks`
   - ✅ Status: `"pending"`
   - ✅ Timestamps: Current time

4. **Admin Dashboard:**
   - ✅ New feedback appears in list
   - ✅ Status: **"Pending"** badge
   - ✅ Stats updated (Total, Pending counts)
   - ✅ Admin can view, update, respond, delete

---

## 🔍 **Testing the Flow**

### **Test 1: Submit Feedback (Mobile App)**

```bash
# Test feedback submission
curl -X POST http://18.143.118.157:8610/feedback/submit \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test feedback",
    "content": "This is a test feedback submission",
    "email": "test@example.com"
  }'
```

**Expected Response:**
```json
{
  "code": "0000",
  "msg": "Success",
  "data": {
    "feedback_id": "FDB_xxxxx"
  }
}
```

### **Test 2: Check Admin Dashboard**

1. Open Admin Dashboard: `http://localhost:3000/feedback`
2. Look for new feedback in the list
3. Status should be **"Pending"**
4. Stats should show updated counts

### **Test 3: Rate Limiting**

```bash
# Submit feedback twice quickly (within 1 minute)
curl -X POST http://18.143.118.157:8610/feedback/submit ...
# (Wait < 1 minute)
curl -X POST http://18.143.118.157:8610/feedback/submit ...
```

**Expected:**
- First request: `200 OK` ✅
- Second request: `429 Too Many Requests` ❌

---

## 📝 **What Admin Should Do**

### **After Feedback is Submitted:**

1. **View Feedback** (Recommended)
   - ✅ Check new feedback in dashboard
   - ✅ Read content and understand issue
   - ✅ Check user contact information

2. **Update Status** (Optional)
   - ✅ Change status from "pending" to "reviewing"
   - ✅ Assign to specific admin (if needed)

3. **Respond** (Optional)
   - ✅ Add admin response/notes
   - ✅ Update status to "resolved" when done

4. **Delete** (Only if spam/invalid)
   - ✅ Delete feedback if it's spam or invalid
   - ✅ Use bulk delete for multiple spam entries

---

## 🎯 **Quick Reference**

| Step | Action | Result |
|------|--------|--------|
| **1** | User submits feedback | Feedback sent to backend |
| **2** | Backend validates | Validates format and email |
| **3** | Backend rate limits | Checks 1 per minute limit |
| **4** | Backend creates record | Saves to database |
| **5** | Backend returns ID | User gets feedback ID |
| **6** | Admin views dashboard | Sees new feedback |
| **7** | Admin updates status | Changes to "reviewing" |
| **8** | Admin responds | Adds response/notes |
| **9** | Admin resolves | Changes to "resolved" |

---

## ✅ **Summary**

**When a user submits feedback from the mobile app:**

1. ✅ **Backend validates** and rate limits (1 per minute)
2. ✅ **Feedback is created** in database with status "pending"
3. ✅ **User gets feedback ID** in response
4. ✅ **Admin dashboard shows** new feedback immediately
5. ✅ **Admin can view, update, respond, or delete** feedback

**The feedback appears in the Admin Dashboard with:**
- Status: **"Pending"**
- Category: **"Other"** (or user-selected)
- User email: From `contact_email` field
- Timestamp: When feedback was submitted

**Admin can then:**
- View full details
- Update status (pending → reviewing → resolved)
- Add admin response
- Delete if needed

---

**This is the complete feedback flow!** 🚀
