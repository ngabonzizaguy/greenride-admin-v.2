# Mobile App (Flutter) - Backend Integration Handoff

## Overview

The backend has been updated with security, privacy, and ETA improvements. The mobile app needs to integrate with these changes. This document covers all required mobile app changes.

---

## 1. CRITICAL: Fix Driver "Go Online" Flow

### Current Bug
The **Confirm button** in vehicle selection does nothing after the driver taps it. Drivers cannot go online.

### Required Backend Call
When the driver toggles Online and selects a vehicle, the app must call:

```
POST /online
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "vehicle_id": "<selected_vehicle_id>",
  "latitude": <current_gps_latitude>,
  "longitude": <current_gps_longitude>
}
```

**All three fields are required.** Missing any will return error code `2001` (Invalid Params).

### Success Response
```json
{
  "code": "0000",
  "msg": "Success",
  "data": ""
}
```

### Error Codes
| Code | Meaning | Action |
|------|---------|--------|
| `2001` | Missing vehicle_id, latitude, or longitude | Ensure GPS is enabled and vehicle is selected |
| `3005` | User is not a driver | Check user type |
| `6700` | Vehicle not found | Vehicle ID is invalid |
| `1000` | System error | Retry or show error message |

### Going Offline
```
POST /offline
Authorization: Bearer <jwt_token>
Content-Type: application/json

{}
```

### Important Notes
- After successful `/online`, the app **must** start sending periodic location updates via `POST /location/update`
- Location updates should be sent every 30-60 seconds while driver is online
- The driver will not appear in nearby searches without recent location updates (2-minute freshness window)

---

## 2. Location Updates While Online

Once online, the driver app must periodically send location:

```
POST /location/update
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "latitude": <float64>,
  "longitude": <float64>,
  "online_status": "online",
  "heading": <float64>,
  "speed": <float64>,
  "accuracy": <float64>,
  "altitude": <float64>
}
```

- `latitude`, `longitude`: Required
- `online_status`: Optional, defaults to "online"
- `heading`, `speed`, `accuracy`, `altitude`: Optional but recommended for map display

---

## 3. Phone Number Privacy Changes

### What Changed
Phone numbers are now **masked** in order responses unless the requester is authorized.

- **Drivers** see passenger phone ONLY after accepting the ride (status: `accepted`, `driver_coming`, `driver_arrived`, `in_progress`)
- **Passengers** see driver phone ONLY after a driver is assigned and ride is active
- **Everyone else** sees masked phones: `+25****456`

### New Endpoint: `/order/contact`
To get the actual phone number for calling, use:

```
POST /order/contact
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "order_id": "<order_id>"
}
```

### Response (Allowed)
```json
{
  "code": "0000",
  "data": {
    "allowed": true,
    "phone": "+250788123456",
    "name": "John Doe"
  }
}
```

### Response (Not Allowed)
```json
{
  "code": "0000",
  "data": {
    "allowed": false,
    "phone": "",
    "name": ""
  }
}
```

### Rules
| Requester | Order Status | Result |
|-----------|-------------|--------|
| Assigned driver | `accepted` / `driver_coming` / `driver_arrived` / `in_progress` | Passenger phone revealed |
| Passenger | Same active statuses (when driver assigned) | Driver phone revealed |
| Any other case | Any | `allowed: false` |

### Migration Steps
1. Replace any direct phone number usage from order details with `/order/contact`
2. Before showing "Call" button, check `allowed` field
3. Use `phone` from response to initiate `tel:` intent

---

## 4. Live ETA Endpoint

### New Endpoint: `/order/eta`
Get real-time ETA from driver's current location to pickup or dropoff:

```
POST /order/eta
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "order_id": "<order_id>"
}
```

### Response
```json
{
  "code": "0000",
  "data": {
    "order_id": "R20260206123456",
    "eta_minutes": 12,
    "distance_km": 5.3,
    "driver_latitude": -1.9403,
    "driver_longitude": 29.8739,
    "pickup_latitude": -1.9500,
    "pickup_longitude": 29.8800,
    "mode": "accurate",
    "updated_at": 1738847123456
  }
}
```

### Fields
| Field | Description |
|-------|-------------|
| `eta_minutes` | Estimated minutes until arrival |
| `distance_km` | Distance in kilometers |
| `mode` | `"accurate"` (Google Directions API) or `"rough"` (Haversine estimate) |
| `driver_latitude/longitude` | Driver's current position |
| `updated_at` | Timestamp of this calculation |

### ETA Logic
| Order Status | ETA Target |
|-------------|-----------|
| `accepted` / `driver_coming` / `driver_arrived` | ETA to **pickup** |
| `in_progress` | ETA to **dropoff** |
| Other statuses | No ETA returned |

### Usage
- Poll every 10-15 seconds during active rides
- Display ETA to passenger while waiting for driver
- Display ETA during trip for remaining time

### Permission
Only the assigned driver or the passenger can call this endpoint.

---

## 5. Nearby Drivers ETA Improvement

### Current Behavior
The `GET /drivers/nearby` endpoint defaults to "rough" ETA mode (2 min/km estimate).

### Recommended Change
Add `eta_mode=accurate` to the request for better ETA:

```
GET /drivers/nearby?latitude=-1.9403&longitude=29.8739&radius_km=15&limit=50&eta_mode=accurate
```

ETA modes:
- `rough`: Haversine distance * 2 min/km (fast but inaccurate)
- `accurate`: Google Directions API (slower but realistic, applied to first 12 drivers)
- `none`: No ETA calculation

---

## 6. Acceptance Notifications - New Fields

When a driver accepts a ride, the FCM notification now includes:

| Field | Type | Description |
|-------|------|-------------|
| `DriverToPickupETA` | int | Estimated minutes until driver arrives at pickup |
| `DriverToPickupDistance` | string | Distance in km (e.g., "5.3") |

Use these to display ETA immediately after acceptance, before the first `/order/eta` poll.

---

## 7. Booking Flow Summary

### A. "Find Drivers" (Manual Selection)
1. Passenger calls `GET /drivers/nearby` to see available drivers
2. Passenger selects a driver and creates order with `provider_id` set
3. **Only the selected driver** receives FCM notification
4. Only that driver can accept (backend enforced)
5. After acceptance: phone revealed via `/order/contact`, ETA via `/order/eta`

### B. "Taxi Booking" (Broadcast)
1. Passenger creates order WITHOUT `provider_id`
2. Order is broadcast to all eligible nearby drivers via FCM
3. **First driver to accept wins** (atomic DB transaction)
4. All other drivers are immediately locked out (dispatch records cancelled)
5. After acceptance: same phone/ETA access as above

### Security Guarantees (Backend Enforced)
- Unassigned drivers CANNOT see passenger phone numbers
- Unassigned drivers CANNOT accept manually-assigned orders
- Dispatched orders are hidden from `/order/nearby` browse
- Phone numbers are masked in all order responses until authorized
- `/order/contact` is the ONLY way to get the actual phone number

---

## API Base URL
```
https://api.greenrideafrica.com
```

## Authentication
All endpoints require `Authorization: Bearer <jwt_token>` header.
