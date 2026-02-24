# MoMo and Cash Payment – Investigation and Fix Plan

## 1. MoMo (MTN Mobile Money) – Passenger Does Not Receive Prompt

### 1.1 Backend Flow (admin)

- **Endpoint:** `POST /order/payment` (handler `OrderPayment`, body `OrderPaymentRequest`).
- **Request body:** `order_id` (required), `phone` (optional, defaults to user’s phone), `payment_method` (optional; must be `"momo"` for MoMo).
- **Flow:** `OrderService.OrderPayment` → `PaymentService.OrderPayment` → `GetPaymentRouter(payment_method: momo)` → `MoMoService.Pay(payment)`.
- **MoMoService.Pay:** Builds MTN Request to Pay (Collection API), sends to MTN with `X-Callback-Url` = `{CallbackHost}/webhook/momo/{payment_id}`. On 202 Accepted, returns pending; MTN sends the USSD/prompt to the payer’s phone; MTN later calls the webhook; backend updates payment and order.
- **Webhook:** `POST /webhook/momo/:payment_id` (no auth) → `MoMoWebhook` → `MoMoService.ResolveResponse` → `CheckOrderPayment` to update order.

Backend behaviour is correct: it expects **one** endpoint for all payment methods (`/order/payment`) and supports MoMo when `payment_method` is `momo` and `phone` is set.

### 1.2 App Flow (green_ride_app)

- **MomoPaymentController** (`lib/features/payments/controller/momo_payment_controller.dart`):
  - **`useMockPayment = true`** (line 17): When true, **no backend call** is made. `_mockInitiatePayment` only logs and waits; `_mockCheckStatus` simulates success/failure after a few polls. So with the current build, **the backend and MTN are never called** – the passenger cannot receive a real MoMo prompt.
  - When `useMockPayment` is false, the app calls:
    - **Initiate:** `POST payment/momo/initiate` with `order_id`, `amount`, `currency`, `phone`.
    - **Status:** `GET payment/momo/status/{orderId}`.
- **Backend has no such routes:** There is no `payment/momo/initiate` or `payment/momo/status`. The only payment endpoint used for orders is `POST /order/payment`.

So even if mock is turned off, the app would call non‑existent endpoints and get 404/4xx, and no MoMo request would be sent to MTN.

### 1.3 Root Causes (MoMo)

1. **Mock mode on:** `useMockPayment = true` → no real API call, no MTN request, no prompt.
2. **Wrong endpoints when not mock:** App uses `payment/momo/initiate` and `payment/momo/status`; backend only has `order/payment` and order detail (for status).

### 1.4 Fix Plan (MoMo)

**App only (no backend change):**

1. **Use real payment for production:**  
   - Either set `useMockPayment = false` in code, or (preferred) make it configurable (e.g. from flavor or remote config) so production uses the real API and dev/test can keep mock.

2. **Initiate MoMo via existing order payment endpoint:**  
   - In `_realInitiatePayment`, call **`POST order/payment`** (reuse `payWithCard` endpoint string or add a named constant) with body:
     - `order_id`
     - `payment_method`: `"momo"`
     - `phone`: payer’s phone (e.g. E.164 or same format backend expects).
   - Do **not** call `payment/momo/initiate`.

3. **Status via order detail:**  
   - In `_realCheckStatus`, **poll `order/detail`** (e.g. `rideDetail` with `order_id`) and read `data.payment_status` (`success` / `failed` / `pending`).  
   - Stop polling when status is success or failed; surface errors from response.

4. **Error handling and logging:**  
   - Parse backend error code/msg from `order/payment` and status responses; show a clear message to the user (e.g. “Payment could not be started”) and log request/response (no sensitive data).

**Backend / ops (no code change required for the flow):**

- Ensure production has correct MoMo credentials and **callback_host** so MTN can reach `https://api.greenrideafrica.com/webhook/momo/{payment_id}`.
- Ensure migration/router for `momo` and the MoMo channel account are active for the environment.

---

## 2. Cash Payment and Verification Code

### 2.1 Backend Flow (admin)

- **Who can call `POST /order/payment` with cash:** Only the **driver** is allowed to use `payment_method: "cash"`. If a passenger sends `payment_method: "cash"`, the backend returns **PermissionDenied**.
- **Driver confirms cash:** `POST /order/cash/received` with body `OrderPaymentRequest` (e.g. `order_id`). Handler sets `payment_method = cash` and calls `OrderPayment`; the driver is then allowed, so the order is marked paid. There is **no verification code** in the request or in the backend logic.

### 2.2 App Flow (green_ride_app)

- **Payment method selection (trip ended):**  
  `IsTripEndedView` → “Pay Now” → `PaymentFlowCoordinator` → `EnhancedPaymentMethodSheet`. On **Cash** selected, `_handleCashPayment()` only shows a snackbar (“Please pay RWF … to the driver”) and calls `onCashSelected`; in `IsTripEndedView`, `onCashSelected` is an empty callback. So:
  - The app **does not** call `POST /order/payment` with `payment_method: cash`.
  - The app **does not** show the **CashPaymentSheet** (code + “I’ve Paid the Driver”) in this flow.

- **CashPaymentSheet** (used from **Developer Options** only):
  - Generates a **4‑digit code locally** (`_generateCode()`).
  - Shows “Pay the driver”, “Show this code to the driver”, “Driver enters code to confirm”.
  - “I’ve Paid the Driver” only sets local state (`_paymentConfirmed = true`) and shows “Waiting for Driver”. It **does not** call the backend.

- **Driver:**  
  `paymentConfirmApi(rideId)` calls `POST order/cash/received` with `order_id` only (no code). So the driver can confirm cash without entering any code; backend does not expect one.

### 2.3 Root Causes (Cash)

1. **Cash not wired in post‑trip flow:** Selecting Cash after “Pay Now” does not call the backend and does not show the cash sheet with code.
2. **Code is client‑only and not validated:** The 4‑digit code is generated and displayed only in the app; backend never sees or validates it, so “driver enters code” is currently cosmetic.

### 2.4 Fix Plan (Cash) – Minimal (Implemented)

**Goal:** Make cash payment flow visible and clear after trip end; backend allows only the driver to confirm cash.

1. **When passenger selects Cash after “Pay Now”:**
   - Show **CashPaymentSheet** with `orderId`, `amount`, 4‑digit code, and “I’ve Paid the Driver” / “Use different method”.

2. **When passenger taps “I’ve Paid the Driver”:**
   - Do **not** call `POST /order/payment` with cash (backend would return PermissionDenied for passenger).
   - Close the sheet and show a message: e.g. “Please pay RWF X to the driver. They will confirm receipt in the app.”

3. **Driver:**  
   “Confirm cash received” continues to call `POST /order/cash/received` with `order_id`. That is the only way the order is marked paid for cash; no code is required.

4. **Code in the sheet:**  
   The 4‑digit code remains a **reference** for the driver (e.g. to match the ride). Backend does not generate or validate it.

**Impact:**  
- No backend changes.  
- Passenger flow: select cash → see amount and code → “I’ve Paid” → sheet closes, message to pay driver; driver confirms via existing “Confirm cash received” in their app.

### 2.5 Optional Later: Verification Code on Backend

If you want the driver to **enter** the code in the app and backend to validate it:

- Backend: When order payment is initiated with `payment_method: cash`, create a **pending** cash payment and generate a short‑lived **verification code** (e.g. 4 digits), store it on the order or payment, return it in the response; expose an endpoint (e.g. extend `order/cash/received` or add `order/cash/confirm`) that accepts `order_id` + `verification_code` and, if valid, marks order paid and clears the code.
- App: Passenger flow gets the code from the API and shows it; driver flow has an input for the code and sends it with the confirm request. This requires backend changes and is out of scope for the minimal fix above.

---

## 3. Logs You Provided

- Request to `order/eta` and response from `order/detail`: order is `trip_ended`, `payment_status: "pending"`, `payment_method: ""`, amount 5682 RWF. So the order is in the correct state for payment; no payment method was stored yet.
- There is **no** log of a call to `order/payment` or any MoMo endpoint in the snippet. That is consistent with: (1) mock mode so no real initiate call, or (2) a call to a wrong path that fails before logging. Fixing the app to call `POST /order/payment` with `payment_method: momo` and `phone` will allow the backend to call MTN and the user to receive the prompt (subject to MTN and callback configuration).

---

## 4. Summary Table

| Item | Finding | Fix (concise) |
|------|--------|----------------|
| MoMo – no prompt | Mock on + wrong endpoints | Use `order/payment` for initiate; poll `order/detail` for status; turn off mock in prod or make configurable. |
| Cash – not functional after trip | Cash selection only shows snackbar; sheet not shown; backend allows only driver to confirm cash | On Cash selected, show CashPaymentSheet with code; on “I’ve Paid” close sheet and show message; driver confirms via existing “Confirm cash received” (order/cash/received). |
| Cash – code | Code is UI‑only; backend doesn’t use it | Keep code as reference; optional later: backend generates/validates code and driver submits it. |

No fake success, no suppressing errors, no breaking existing flows: use existing backend contracts and add proper error handling and logging in the app.
