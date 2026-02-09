# Ride Assignment, Call Permissions & ETA - Implementation Tracker

## Branch: `feature/ride-assignment-call-permissions`

## Status: ALL PHASES COMPLETE - Ready for review/merge

---

## Overview

Fixing ride assignment logic, call permissions, and ETA accuracy in the GreenRide backend and admin dashboard. The mobile app is in a separate repo -- all changes here are backend-enforced.

## Identified Gaps (from analysis)

| # | Gap | Severity | Phase | Status |
|---|-----|----------|-------|--------|
| 1 | Passenger phone exposed to ALL drivers via `/order/nearby` | CRITICAL | Phase 1 | FIXED |
| 2 | No call permission enforcement | CRITICAL | Phase 1 | FIXED |
| 3 | Driver phone exposed to passengers via `/drivers/nearby` | HIGH | Phase 1 | FIXED |
| 4 | `GetOrderDetail` leaks phone to any driver for unassigned orders | HIGH | Phase 1 | FIXED |
| 5 | `AcceptOrder` WHERE clause fails for manually pre-assigned orders | HIGH | Phase 2 | FIXED |
| 6 | Dispatched orders still visible in `/order/nearby` browse | MEDIUM | Phase 2 | FIXED |
| 7 | No live ETA endpoint after driver assignment | MEDIUM | Phase 3 | FIXED |

---

## Phase 1: Phone Number Protection & Call Permission Enforcement - COMPLETE

### 1.1 Create phone masking utility
- **File:** `backend/greenride-api-clean/internal/utils/mask.go` (NEW)
- **Status:** COMPLETE
- **What:** `MaskPhone(phone string) string` -- masks middle of phone number
- **Pattern:** `+25****456` (first 3 + `****` + last 3)

### 1.2 Add `GetOrderInfoSanitized()` method
- **File:** `backend/greenride-api-clean/internal/services/order_service.go` (after line ~288)
- **Status:** COMPLETE
- **What:** Wrapper around `GetOrderInfo()` that masks phone based on requester role and order status
- **Rules:**
  - Assigned driver sees passenger phone ONLY if status in {accepted, driver_coming, driver_arrived, in_progress}
  - Passenger sees driver phone ONLY if order has provider and status in same set
  - Anyone else: ALL phones masked
  - Masks both `info.Passenger.Phone`/`info.Driver.Phone` AND `info.Details.PassengerPhone`/`info.Details.DriverPhone`

### 1.3 Sanitize `GetNearbyOrders` response
- **File:** `backend/greenride-api-clean/internal/services/order_service.go` (GetNearbyOrders func)
- **File:** `backend/greenride-api-clean/internal/protocol/request.api.go` (added `RequesterID` field)
- **File:** `backend/greenride-api-clean/internal/handlers/api.order.go` (passes `req.RequesterID = user.UserID`)
- **Status:** COMPLETE

### 1.4 Sanitize `GetOrderDetail` response
- **File:** `backend/greenride-api-clean/internal/handlers/api.order.go` (GetOrderDetail handler)
- **Status:** COMPLETE
- **What:** Uses `models.GetOrderByID()` directly then `GetOrderInfoSanitized()`. Also added `models` import.

### 1.5 Strip phone from `GetNearbyDrivers` for passengers
- **File:** `backend/greenride-api-clean/internal/handlers/api.location.go` (GetNearbyDrivers handler)
- **Status:** COMPLETE
- **What:** If requester `IsPassenger()`, strips `Phone` field from all returned drivers

### 1.6 Create `/order/contact` endpoint (call permission)
- **Files:**
  - `backend/greenride-api-clean/internal/protocol/request.api.go` -- added `OrderContactRequest`, `OrderContactResponse`
  - `backend/greenride-api-clean/internal/services/order_service.go` -- added `GetOrderContactInfo()`
  - `backend/greenride-api-clean/internal/handlers/api.order.go` -- added `GetOrderContact` handler
  - `backend/greenride-api-clean/internal/handlers/api.go` -- added route `authRequired.POST("/order/contact", ...)`
- **Status:** COMPLETE
- **Logic:** Only returns phone if requester is assigned driver/passenger AND order status is in active set

### 1.7 Sanitize `GetOrdersByUser` response
- **File:** `backend/greenride-api-clean/internal/services/order_service.go` (GetOrdersByUser func)
- **Status:** COMPLETE
- **What:** Changed `GetOrderInfo(order)` to `GetOrderInfoSanitized(order, req.UserID, userType)`

---

## Phase 2: Locking & Exclusivity Improvements - COMPLETE

### 2.1 Filter dispatched orders from `GetNearbyOrders`
- **File:** `backend/greenride-api-clean/internal/services/order_service.go` (GetNearbyOrders)
- **Status:** COMPLETE
- **What:** Added two WHERE clauses:
  - `provider_id IS NULL OR provider_id = ''` (exclude pre-assigned)
  - `order_id NOT IN (SELECT order_id FROM dispatch_records WHERE status = 'pending')` (exclude actively dispatched)

### 2.2 Fix `AcceptOrder` for manually pre-assigned orders
- **File:** `backend/greenride-api-clean/internal/services/order_service.go` (AcceptOrder, ~line 1001)
- **Status:** COMPLETE
- **What:** Changed WHERE from `provider_id IS NULL or provider_id =''` to `provider_id IS NULL OR provider_id = '' OR provider_id = ?` (req.UserID)

### 2.3 Cancel stale dispatches on order cancellation
- **File:** `backend/greenride-api-clean/internal/services/order_service.go` (CancelOrder func)
- **Status:** COMPLETE
- **What:** Added inside transaction: cancels all pending dispatch records when order is cancelled

---

## Phase 3: ETA Accuracy Improvements - COMPLETE

### 3.1 Create `/order/eta` endpoint
- **Files:**
  - `backend/greenride-api-clean/internal/protocol/request.api.go` -- added `OrderETARequest`, `OrderETAResponse`
  - `backend/greenride-api-clean/internal/services/order_service.go` -- added `GetOrderETA()`
  - `backend/greenride-api-clean/internal/handlers/api.order.go` -- added `GetOrderETA` handler
  - `backend/greenride-api-clean/internal/handlers/api.go` -- added route `authRequired.POST("/order/eta", ...)`
- **Status:** COMPLETE
- **Logic:**
  - Permission: only passenger or assigned driver
  - If status accepted/driver_coming/driver_arrived: ETA to **pickup**
  - If status in_progress: ETA to **dropoff**
  - Uses Haversine rough estimate (2 min/km). Google Directions can be added later.

### 3.2 Tighten location freshness window
- **File:** `backend/greenride-api-clean/internal/services/user_service.go`
- **Status:** COMPLETE
- **What:** Changed from 5-minute (300000ms) to 2-minute (120000ms) window in GetNearbyDrivers SQL

### 3.3 Include ETA in acceptance notification
- **File:** `backend/greenride-api-clean/internal/services/order_service.go` (NotifyOrderAccepted)
- **Status:** COMPLETE
- **What:** Calculates driver-to-pickup ETA/distance, stores in order metadata, passed to FCM via `DriverToPickupETA` and `DriverToPickupDistance` params

### 3.4 Add ETA to dispatch notification
- **File:** `backend/greenride-api-clean/internal/services/dispatch_service.go` (SendDispatchNotifications)
- **Status:** COMPLETE
- **What:** Adds `DriverToPickupETA` (int minutes) and `DriverToPickupDistance` (string km) to dispatch notification params

---

## Phase 4: Admin Dashboard Updates - COMPLETE

### 4.1 Add `getOrderETA` and `getOrderContact` to API client
- **File:** `src/lib/api-client.ts`
- **Status:** COMPLETE
- **What:** Two new typed methods: `getOrderETA(orderId)` and `getOrderContact(orderId)`

### 4.2 Add ETA polling in Quick Booking
- **File:** `src/app/(dashboard)/quick-booking/page.tsx`
- **Status:** COMPLETE
- **What:** `useEffect` polls `/order/eta` every 10 seconds when booking is complete, updates displayed ETA

---

## New API Endpoints Created

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/order/contact` | POST | Get counterpart's phone for calling (permission-gated) |
| `/order/eta` | POST | Get live ETA from driver location to pickup/dropoff |

## Files Modified (Summary)

### Backend (Go)
| File | Changes |
|------|---------|
| `internal/utils/mask.go` | **NEW** - Phone masking utility |
| `internal/services/order_service.go` | Added `GetOrderInfoSanitized()`, `GetOrderContactInfo()`, `GetOrderETA()`. Updated `GetNearbyOrders`, `GetOrdersByUser`, `CancelOrder`, `AcceptOrder`, `NotifyOrderAccepted`, `NotifyPassenger` |
| `internal/services/dispatch_service.go` | Added ETA to dispatch notifications |
| `internal/services/user_service.go` | Tightened location freshness to 2 min |
| `internal/handlers/api.order.go` | Added `GetOrderContact`, `GetOrderETA` handlers. Updated `GetOrderDetail` to use sanitized output. Added `models` import |
| `internal/handlers/api.location.go` | Strip driver phone for passengers in GetNearbyDrivers |
| `internal/handlers/api.go` | Added routes for `/order/contact` and `/order/eta` |
| `internal/protocol/request.api.go` | Added `RequesterID` to `GetNearbyOrdersRequest`. Added `OrderContactRequest/Response`, `OrderETARequest/Response` |

### Frontend (TypeScript)
| File | Changes |
|------|---------|
| `src/lib/api-client.ts` | Added `getOrderETA()` and `getOrderContact()` methods |
| `src/app/(dashboard)/quick-booking/page.tsx` | Added ETA polling effect after booking complete |

## Build Status
- Go backend: `go build ./...` -- PASS
- TypeScript: `tsc --noEmit` -- PASS

## What the Mobile App Team Needs to Know
See the prompt/handoff document that was drafted separately. Key points:
1. Phone numbers in order detail responses are now masked unless authorized
2. New `/order/contact` endpoint is the only way to get phone for calling
3. New `/order/eta` endpoint provides live ETA polling
4. FCM acceptance notifications now include `DriverToPickupETA` and `DriverToPickupDistance`
