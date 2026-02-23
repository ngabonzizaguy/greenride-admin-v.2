# Driver requests and nearby-drivers flow

## Why drivers weren’t receiving requests

1. **Backend**  
   After an order is created, **StartAutoDispatch** runs and creates **dispatch records** for selected drivers and sends FCM.  
   **GetNearbyOrders** used to return only orders that had **no** pending dispatch. So every new order was excluded from “nearby” and drivers only saw requests via FCM. If FCM wasn’t received or the app didn’t open, the driver saw nothing in the “Available Requests” list.

2. **App**  
   The driver “Available Requests” screen only called **carRideRequestApi** (POST `/nearby`) and had no polling, so the list didn’t refresh unless the user pulled to refresh.

## Changes made

### Backend (greenride-api-clean)

- **Dispatch eligibility (broadcast = minimal rules)**  
  - **FindEligibleDrivers**: returns **all drivers who have a vehicle** (no filter by vehicle category/level). So every order is offered to all such drivers; vehicle match is not required.  
  - **EvaluateDriverForOrder**: only **two mandatory** checks — (1) driver must be **online**, (2) driver must have **no active ride** (one ride at a time). Distance, queue capacity, and time window are optional (e.g. distance only applied when `MaxDistance > 0`).  
  - So: if a driver is **online** and **not on an active ride**, they receive the request. No mandatory vehicle type or location filters.

- **GetNearbyOrders**  
  - Still returns “broadcast-style” orders (requested, no provider, no pending dispatch).  
  - **Also** returns orders that have a **pending dispatch to the requesting driver**.  
  - For those, the response includes **dispatch_id** so the app can send it on accept (first-accept-wins).
- **Order** protocol: added **DispatchID** (optional) for driver responses.
- **AcceptOrder** already supported **dispatch_id** and first-accept-wins (only first accepter gets the order; others get RideAlreadyBooked).

### App (green_ride_app)

- **RideOrderBean**: added **dispatchId**; parsed from API and sent on accept when present.
- **FoodRequestController.driverAcceptOrder**: optional **dispatchId**; included in the accept request when provided.
- **Driver list**: passes **order.dispatchId** when calling accept; **TrackOrderDetailsScreen** already had **dispatchId** and passes it to **acceptApi**.
- **Driver “Available Requests” screen**: **initState** loads requests and starts a **10s polling** of **carRideRequestApi** so new requests appear without pull-to-refresh.

### Passenger – nearby drivers and active rides

- **Nearby Drivers (select a driver)**  
  - **NearbyDriver** model: added **isBusy** (from API **is_busy**).  
  - **DriverDetailSheet**: when **driver.isBusy**:
    - Shows a banner: “This driver has an active ride. Book without selecting a driver so the first available driver can accept.”
    - Replaces “Select This Driver” with **“Book Without Selecting”**; that callback closes the sheet and opens **BookRideScreen**, which creates an order **without** a pre-selected driver → backend runs **StartAutoDispatch** and **broadcasts to all online drivers** (no active ride). First driver to accept gets the order and can call the passenger; others see “request already accepted”.
- **Taxi Booking (no driver selected)**  
  - When the passenger uses **“Taxi Booking”** (book without choosing a driver), the order is created with no pre-assigned driver. Backend runs **StartAutoDispatch** and sends the request to **all online drivers with no active ride**. First to accept wins; only that driver can see the passenger’s number and proceed.

### First-accept-wins and calling

- Backend **AcceptOrder** + **HandleDriverAccept** ensure only the first driver to accept gets the order; others receive **RideAlreadyBooked** (e.g. “Ride already booked”).
- **GetOrderContact** (call passenger) is only allowed for the **assigned** driver, so only the driver who accepted can get the passenger’s phone; others cannot call.

## How to verify

1. **Driver receives requests**  
   Log in as driver, go to “Available Requests”. Create a ride as passenger (or from admin quick booking). Within ~10s the request should appear in the driver list (and/or via FCM). Driver can accept; only the first accepter gets the order.

2. **Request already accepted**  
   Two drivers on “Available Requests”; create one ride. First driver accepts. Second driver, if they tap the same request (before list refresh), gets “Ride already booked” (or similar) from the API.

3. **Passenger – driver has active ride**  
   As passenger, open “Nearby Drivers” and select a driver who is busy. Sheet shows the “Driver has active ride” banner and “Book Without Selecting”; tapping it opens Book Ride (broadcast).

## What you need to do

- **Deploy the backend** that includes: (1) **GetNearbyOrders** returning “dispatched to me” + `dispatch_id`, (2) **FindEligibleDrivers** using all drivers with a vehicle (no vehicle filter), (3) **EvaluateDriverForOrder** with only online + no active ride mandatory. Until this is deployed, drivers may not see requests in the app list and dispatch may still use the old eligibility rules.
- **Driver app**: Ensure the driver is **online** and has **no active ride**; the app now sends real location for `/nearby` and polls every 10s.
- **Optional (selected busy driver notification)**: If you want the **selected (busy) driver** to always receive a notification when a passenger chose them but then used “Book Without Selecting”, that would require an extra backend/notification path (e.g. a “passenger considered you but you were busy” message). Currently only drivers who are eligible (online, no active ride) receive the broadcast.

## Relevant files

- **Backend**: `internal/services/dispatch_service.go` (FindEligibleDrivers, EvaluateDriverForOrder), `internal/services/order_service.go` (GetNearbyOrders), `internal/protocol/order.go` (DispatchID).
- **App – driver**: `lib/ui/modules/driver_screen/controller/food_request_controller.dart`, `booking_available_request_screen.dart`, `driver_booking_order_list.dart`, `controller/model/ride_request_response.dart`.
- **App – passenger**: `lib/features/nearby_drivers/nearby_drivers_controller.dart`, `driver_detail_sheet.dart`, `find_nearby_drivers_screen.dart`.
