# Backend Requirements for Admin Dashboard v2

> **To:** Backend Team / Backend Agent  
> **From:** Admin Dashboard Team  
> **Subject:** Missing API Endpoints for New Admin Features

We have completed the UI for the **Feedback System** and **Support Configuration**. To make these features live (disabling "Demo Mode"), we need the following endpoints implemented in the Go Backend (Admin Service, Port 8611).

---

## 1. Support Configuration (New Feature)

We need a way to dynamically update the support contact info shown in the mobile app.

### Database Requirement
We need a storage mechanism (e.g., a key-value table `t_sys_config` or a specific table `t_support_config`) to store:
- `support_phone` (string)
- `support_email` (string)
- `support_whatsapp` (string)
- `operating_hours` (string)
- `faq_url` (string)

### API Endpoints

#### A. Get Configuration
- **Method:** `GET`
- **Path:** `/config/support`
- **Auth:** Admin JWT
- **Response:**
  ```json
  {
    "code": "0000",
    "msg": "Success",
    "data": {
      "phone": "+250 788 000 000",
      "email": "support@greenride.com",
      "whatsapp": "+250 788 000 001",
      "hours": "24/7",
      "faq_url": "https://greenride.com/faq"
    }
  }
  ```

#### B. Update Configuration
- **Method:** `POST`
- **Path:** `/config/support`
- **Auth:** Admin JWT
- **Request Body:** (Same JSON structure as Response)
- **Response:** Success confirmation.

---

## 2. Feedback Management (Enhancement)

The backend currently has `t_feedbacks` and `POST /feedback/submit` (Mobile). We need **Admin Management** endpoints to view and reply to these.

### Database Updates (If not already present)
Ensure `t_feedbacks` has:
- `status`: ('pending', 'reviewing', 'resolved', 'closed')
- `admin_response`: Text field for the reply
- `category`: ('driver', 'vehicle', 'pricing', 'safety', etc.)
- `severity`: ('low', 'medium', 'high', 'critical')

### API Endpoints

#### A. Search Feedback (Admin)
- **Method:** `POST`
- **Path:** `/feedback/search` (Admin Service)
- **Request Body:**
  ```json
  {
    "page": 1,
    "limit": 10,
    "keyword": "search term",
    "status": "pending",    // optional filter
    "category": "safety",   // optional filter
    "user_id": "USR...",    // optional filter
    "start_date": 12345678, // optional timestamp
    "end_date": 12345678    // optional timestamp
  }
  ```
- **Response:** Standard `PageResult` structure.
 
#### B. Feedback Detail (Admin)
- **Method:** `POST`
- **Path:** `/feedback/detail`
- **Request Body:** `{ "feedback_id": "FB..." }`
- **Response:** Full feedback object including `admin_response`.

#### C. Update Feedback (Admin)
- **Method:** `POST`
- **Path:** `/feedback/update`
- **Request Body:**
  ```json
  {
    "feedback_id": "FB...",
    "status": "resolved",
    "admin_response": "We have processed your refund.",
    "severity": "medium" // optional update
  }
  ```

---

## 3. Implementation Checklist for Backend

1.  [ ] Create/Update MySQL tables (`t_support_config`, `t_feedbacks`).
2.  [ ] Create Go Structs matching the JSON layouts above.
3.  [ ] Implement Gin Handlers in `greenride-api-clean`.
4.  [ ] Register routes in `main.go` or `router.go` under the Admin Group.
5.  [ ] Deploy to Dev/Prod environment.

Once these are live, the Admin Dashboard will automatically start working with real data (we just toggle `DEMO_MODE = false`).

