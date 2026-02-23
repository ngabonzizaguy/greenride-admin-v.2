# Critical Booking & Dispatch — Investigation Report

## A. Why “Book Ride” Can Appear to Do Nothing

### Flow traced

1. **Button** (`book_ride_screen.dart`): `BaseButton` with `onPressed`: if `isEstimatingFare` → **disabled** (`onPressed: null`). Else if `estimatedFare != null` → label **"Book Ride"**, calls `controller.bookRide()`. Else → label **"Confirm Location"**, calls `controller.estimateFare()`.
2. **bookRide()** (`book_ride_controller.dart`): Builds payload (including `price_id` from `estimatedFare.value?.priceId`), POSTs to `order/create`, then checks `code != '0000'` and shows snackbar, or parses and navigates to `TrackOrderDetailsScreen`.

### Exact failure points that can look like “nothing happens”

| # | Failure point | File / function | Cause | Fix applied |
|---|----------------|------------------|--------|-------------|
| 1 | **No feedback when fare missing** | `book_ride_controller.dart` → `bookRide()` | `estimatedFare` or `priceId` null (e.g. estimate not run or expired). Request was still sent with `price_id: null`; backend may reject. User got no message. | **Pre-validation added**: if no valid fare/price_id, show snackbar *"Please wait for the fare to load, or tap Confirm Location again"* and return. No request is sent. |
| 2 | **Backend error 6018/6017 not mapped** | Same | Codes 6018 (driver has active order), 6017 (driver offline) fell into `default` and showed raw `msg`. In some builds/locale, message could be empty or unclear. | **Explicit handling**: 6018 → *"The selected driver is busy. Please choose another driver or book without selecting one."* 6017 → *"The selected driver is offline. Please choose another driver."* |
| 3 | **Generic catch could hide server message** | Same → `catch (e, s)` | Any exception (e.g. `ResponseError` from parsing) was turned into a generic *"Failed to book ride. Please try again."* | **ResponseError handled first**: show `e.msg` so backend message is shown when available. |
| 4 | **Button disabled with no hint** | `book_ride_screen.dart` | When `isEstimatingFare == true`, button is disabled and title is empty. If estimate never completes (e.g. network hang), user sees a dead-looking button. | Logic unchanged; loader and timeout already in place. Pre-validation in `bookRide()` now avoids sending invalid requests and gives a clear message when fare is missing. |

**Conclusion:** The main fix is **pre-validation of fare/price_id** and **explicit handling of 6018/6017** plus **ResponseError.msg** in catch, so the user **always** gets a clear message and the button never appears to “do nothing” without explanation.

---

## B. Whether the Ride Record Is Created

- **When the API returns success (`code == '0000'`):**  
  **Yes.** The backend creates the order (main table + ride detail), then either:
  - **Manual assign** (`provider_id` in request): creates one dispatch record for that driver and calls `SendDispatchNotifications(record)` (FCM).
  - **Auto-dispatch** (no `provider_id`): runs `StartAutoDispatch(order)` → finds drivers by vehicle, filters by runtime/distance/availability, creates dispatch records, then sends FCM per record.

- **When the API returns an error (e.g. 6018, 6017, 6007, 6008):**  
  **No.** The handler returns before creating an order; no ride record and no dispatch.

- **6007 "Ride in progress":** The backend blocks the **passenger** from creating a new order if they have any order with status in `requested`, `pending`, `accepted`, or `in_progress` (`CountActiveRideOrdersByUser`). A single stuck order blocks new bookings until it is cancelled or completed (e.g. admin Rides → Cancel, or DB).

- **Evidence:** Backend `handlers/api.order.go` → `CreateOrder`; `services/order_service.go` → `CreateOrder` (DB create + dispatch branch); success response returns `protocol.Order` with `order_id`.

---

## C. Why Drivers Aren’t Receiving Orders

Dispatch chain (backend → driver app):

1. **Order created** → `CreateOrder` finishes and, for auto-dispatch, runs `go s.DispatchOrder(orderInfo)`.
2. **Dispatch** → `StartAutoDispatch` → `FindEligibleDrivers` (by vehicle category/level) → `GetDriversRuntime` (Redis/DB) → `EvaluateDriverForOrder` (online, distance ≤ max_distance, queue) → `ExecuteDispatch` (create records) → `SendDispatchNotifications(record)` per driver.
3. **FCM** → `MessageService.SendFcmMessage` uses `params["to"]` = driver user ID → `GetUserFCMTokens(userID)`. If no tokens, returns error and logs *"no FCM tokens found for user"*.
4. **Driver app** → `FirebaseMessaging.onMessage` checks `data['notification_type'] == 'new_order_available'` and shows local notification; tap opens order screen.

**Where the chain can break (no order received by driver):**

| Break point | Cause |
|-------------|--------|
| No eligible drivers | `FindDriversByVehicle(category, level)` returns empty (no vehicles/drivers for that category/level). |
| No runtime | `GetDriversRuntime` returns empty (drivers not “online” or no Redis/DB runtime). |
| All filtered | Every driver rejected by `EvaluateDriverForOrder`: offline, or distance &gt; max_distance (e.g. 15 km), or queue full. |
| No FCM token | Driver has no active token in `t_fcm_tokens` → `SendDispatchNotifications` fails with *"no FCM tokens found for user"*. |
| Driver not “online” | Driver must have gone online (toggle) and selected a vehicle; otherwise they are offline and filtered out. |
| Manual assign (6018) | Passenger selected a driver who has an active order → backend returns 6018, no order created, so no notification. |

So “drivers not receiving” is either **no order created** (e.g. 6018 when selecting a busy driver) or **order created but dispatch/FCM failing** (no eligible drivers, no tokens, or driver not online).

---

## D. Concrete Fix Plan (Minimal Regression Risk)

### Done in this pass

1. **Passenger app (Flutter)**
   - **book_ride_controller.dart**
     - Require valid `estimatedFare` and `priceId` before calling the API; otherwise show: *"Please wait for the fare to load, or tap Confirm Location again."*
     - Map backend codes **6018** and **6017** to clear user-facing messages (driver busy / driver offline).
     - In catch: handle **ResponseError** first and show `e.msg`; keep generic message for other exceptions.
   - **No fake success:** Only navigate to track screen when `code == '0000'` and parsing succeeds.
   - **No suppression:** All error paths show a snackbar; no silent return without feedback.

### Recommended next steps (no code change in this doc)

2. **Backend**
   - Ensure **stuck-order cleanup** is running (scheduled task that cancels `driver_arrived` / `in_progress` &gt; 2 hours) so 6018 is less frequent.
   - Log and monitor: *"No eligible drivers found"*, *"No online drivers found"*, *"0/N drivers eligible after evaluation"*, *"no FCM tokens found for user"*, *"Dispatch notification sent to driver X"*.

3. **Driver app**
   - Confirm FCM token is sent at login and on refresh; confirm **going online** and **location updates** so backend has runtime and marks driver available.
   - Optionally: after login, re-send FCM token if it was null at first attempt.
   - **FCM → Accept alignment (done):** Backend sends `OrderID` and `dispatch_id` in new-order FCM data. The new app now passes `dispatch_id` when opening `TrackOrderDetailsScreen` from FCM (getInitialMessage, onMessageOpenedApp, and notification tap payload). The screen already forwards `dispatchId` to `acceptApi`, so the backend’s `HandleDriverAccept` runs correctly when the driver accepts from a push.

4. **Operational**
   - Use admin **Rides → Cancel** for stuck orders to free drivers immediately.
   - Check DB: `t_fcm_tokens` (active tokens per driver), `t_vehicles` (category/level), `t_orders` (stuck in `driver_arrived`/`in_progress`).

This keeps the booking flow **deterministic**: the button always produces a **visible result** (success navigation or an explicit error message), and the dispatch chain is documented so you can trace why a driver did not receive a request.
