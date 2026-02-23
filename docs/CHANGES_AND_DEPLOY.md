# Summary of changes and how to deploy

This describes the **driver-request and broadcast** changes in both codebases, what to expect after deploy, and how to push to `main` and deploy.

---

## 1. What the changes are about

### Admin codebase (`greenride-admin-v.2-1`)

**Backend API (Go)**

- **Drivers receive requests**
  - **GetNearbyOrders** now returns (1) broadcast-style orders and (2) **orders dispatched to the requesting driver**, with **dispatch_id** so the app can send it on accept (first-accept-wins).
  - **FindEligibleDrivers** no longer filters by vehicle category/level: it returns **all drivers who have a vehicle**. So every new order is offered to all such drivers.
  - **EvaluateDriverForOrder** only **requires**: driver is **online** and has **no active ride**. Distance, vehicle match, and time window are optional (distance only when `MaxDistance > 0`).
- **First-accept-wins**: Only the first driver to accept can proceed (call passenger); others get “Ride already booked” and cannot get contact.
- **Order protocol**: Optional **dispatch_id** in order response for drivers.
- **Safety**: Nil checks for `order.Details` in dispatch and logging.

**Admin frontend (Next.js)**

- No code changes required for the driver-request flow. Quick booking with “Auto-assign” creates an order and the backend broadcasts to drivers.

**Docs**

- `DRIVER_REQUESTS_AND_NEARBY_FLOW.md`: flow, eligibility, passenger flows, what you need to do.
- `CHANGES_AND_DEPLOY.md`: this file.

---

### App codebase (`green_ride_app`)

**Driver app**

- **RideOrderBean**: added **dispatchId**; parsed from API and sent on accept when present.
- **FoodRequestController.driverAcceptOrder**: optional **dispatchId**; included in accept request so backend can apply first-accept-wins.
- **Available Requests screen**: loads on open and **polls every 10s** so new requests appear without pull-to-refresh; sends **driver’s real location** for `/nearby` (no longer 0,0).
- **Dashboard**: calls **refreshDriverBooking()** instead of removed **initializeDriverBooking()** when refreshing the driver tab.
- **FCM / ride_order_parts**: **goToTrackOrderDetailsScreen** accepts optional **dispatchId** and passes it to **TrackOrderDetailsScreen** so accept from notification uses first-accept-wins.

**Passenger app**

- **NearbyDriver**: added **isBusy** (from API **is_busy**).
- **DriverDetailSheet**: when **driver.isBusy**, shows a banner and **“Book Without Selecting”**; that opens **BookRideScreen** (broadcast to any available driver).
- **FindNearbyDriversScreen**: passes **onBroadcastInstead** so “Book Without Selecting” opens Book Ride.

**Other**

- **MaintenanceScreen**: new screen shown when API returns 503 + maintenance flag.
- **Feature flag**: **kShowNearbyDrivers** in `utils/feature_flags.dart` (e.g. `true` to show Nearby Drivers on home).

---

## 2. What to expect after deployment

- **Drivers**
  - See new requests in “Available Requests” (from polling and/or FCM).
  - Only need to be **online** and have **no active ride** to receive requests (no vehicle-type requirement).
  - First driver to accept gets the order and can call the passenger; others see “Ride already booked” and cannot get contact.

- **Passengers**
  - **Taxi Booking** (no driver selected): request is broadcast to all online drivers with no active ride; first to accept wins.
  - **Nearby Drivers** → select a driver who has an active ride: banner + “Book Without Selecting” → same as Taxi Booking (broadcast).

- **Admin**
  - Quick booking with **Auto-assign**: order is broadcast to all eligible drivers; first to accept wins. No change to UI flow.

---

## 3. Push to main and deploy

### Admin codebase (`greenride-admin-v.2-1`)

1. **Commit and push to `main`**
   - In the project folder (e.g. `d:\greenride-admin-v.2-1`):
   - Stage all changes, commit with a message like:  
     `Driver requests: broadcast eligibility, GetNearbyOrders dispatched-to-me, docs`
   - Push:  
     `git push origin main`

2. **Deploy**
   - **If you use GitHub Actions:** After push, the workflow deploys. Check the **Actions** tab on GitHub.
   - **If you deploy manually:** SSH to the server, then:
     - `cd /home/ubuntu/greenride-admin-v.2`   *(or your repo path)*
     - `git pull origin main`
     - `./deploy.sh`   (or `./deploy.sh backend` for API only)

3. **Verify**
   - Call the API (e.g. health/maintenance) or use admin dashboard. After deploy, drivers should see requests when they are online and have no active ride.

---

### App codebase (`green_ride_app`)

1. **Commit and push to `main`**
   - In the app project folder (e.g. `C:\Users\G\Desktop\2025\Projects\GRD\GRD-APP\green_ride_app`):
   - Stage all changes, commit with a message like:  
     `Driver requests: dispatchId, polling, real location, nearby busy → broadcast, maintenance screen`
   - Push:  
     `git push origin main`

2. **Build and release**
   - **No automatic server deploy** for the mobile app. Use your normal process:
     - **Codemagic:** Trigger the usual workflow (e.g. “Android Production (Shorebird)” / “iOS Production (Shorebird)”) and use the produced `.aab` / `.ipa` for store or internal testing.
     - **Local:** `flutter build apk` or `flutter build appbundle` / `flutter build ipa` and distribute as you do today.
   - Users get the new behavior only after installing the new build (or after an OTA update if you use Shorebird patches for compatible changes).

---

## 4. Order of operations

1. **Push and deploy admin (backend) first** so the API has GetNearbyOrders (dispatched-to-me + dispatch_id), broadcast eligibility (online + no active ride only), and first-accept-wins.
2. **Push app to `main`** and **build/release** the app so drivers and passengers use the new logic (dispatchId, polling, location, nearby busy → broadcast, maintenance screen).

After both are live, expect: drivers (online, no active ride) see requests and first to accept can call the passenger; passengers can use Taxi Booking or “Book Without Selecting” for broadcast; admin quick booking with Auto-assign broadcasts to all eligible drivers.
