# Driver offline after shift & vehicle selection — investigation and options

## Current behaviour (backend)

- **POST /offline** (driver): Sets driver `online_status` to `offline` and **unbinds the vehicle** (`vehicle.driver_id` cleared for that driver’s vehicle). So when a driver goes offline, their vehicle is freed and **other drivers can select it** when they go online (POST /online with that `vehicle_id`).
- **POST /online** (driver): Requires a `vehicle_id`. Sets driver online and **binds** that vehicle to the driver (`vehicle.driver_id = driver_id`). So only one driver per vehicle at a time.
- There is **no shift or schedule** in the codebase: no shift start/end, no cron, no admin UI for driver shifts.

So “other drivers can select vehicles” already works **if** the current driver goes offline. The problem is **drivers who forget to tap “End shift” / “Go offline”**: they stay online and keep the vehicle bound, so others cannot use it and dispatch still considers them available.

---

## Options

### A. Auto online/offline based on shift times (recommended first step)

**Idea:** Store a simple “shift end” (or “shift end time”) per driver. A scheduled job (cron or worker) runs periodically and sets any driver to **offline** (and unbinds vehicle) when current time is past their shift end. No need for a full “shifts” UI at first.

**Pros:**

- Small change: one or two fields (e.g. `shift_end_time` or `shift_end_at`), one job, reuse existing `UserOffline`.
- Handles “forgot to log out”: after shift end, driver is forced offline and vehicle is freed.
- Admin can set “shift end” when editing a driver (or in a small “Driver shift” section).

**Cons:**

- No “shift start” auto-online unless you add it (e.g. set online at shift start). Usually drivers tap “Start shift” in the app (go online + select vehicle), so auto-offline at end is enough.

**Implementation outline:**

1. **Backend**
   - Add optional `shift_end_at` (or `shift_end_time` per day) on driver/user or in a small `driver_shift` table.
   - Admin API: set/clear shift end for a driver (e.g. “today 18:00” or “2025-02-17T18:00:00Z”).
   - A **scheduled task** (e.g. every 15 min): find drivers where `online_status IN ('online','busy')` and `shift_end_at < now()`; for each call the same logic as `UserOffline` (set offline + unbind vehicle). Optionally set `shift_end_at = null` after applying so you don’t re-run every time.
2. **Admin UI**
   - On driver edit (or “Driver shift” block): time picker / dropdown for “Shift end (optional)”. Save to backend.
3. **App**
   - No change required for basic behaviour; driver can still tap “End shift” to go offline. Optional: show “Shift ends at 18:00” and a reminder to go offline.

---

### B. Full “shifts” feature

**Idea:** Shifts with start and end (e.g. per day or recurring). Admin creates/edits shifts; system can auto-set driver online at shift start and offline at shift end.

**Pros:**

- Clear model: shift start/end, history, maybe recurring.
- Can auto-online at shift start and auto-offline at shift end.

**Cons:**

- More work: shifts table, admin CRUD, UI, timezone handling, and the same cron/job to enforce end (and optionally start).
- Overkill if you only need “driver goes offline at end of day so others can use the vehicle”.

**When to choose:** If you need recurring weekly shifts, multiple shifts per day, or shift history/reports, then implement B. Otherwise A is enough.

---

## Recommendation

- **Start with Option A (auto-offline at shift end):**
  - Add `shift_end_at` (or equivalent) for drivers.
  - Admin UI: set optional “Shift end” when editing a driver.
  - Scheduled job: past shift end → call same logic as `UserOffline` (offline + unbind vehicle).
- **Keep current app flow:** Drivers continue to go **online** by tapping “Start shift” and selecting a vehicle, and **offline** by tapping “End shift”. Auto-offline is a safety net for those who forget.
- **Add Option B later** only if you need full shift scheduling, recurrence, or reporting.

---

## Summary

| Topic | Finding |
|-------|--------|
| Vehicle selection | Backend already unbinds vehicle when driver goes offline; other drivers can then select it when going online. |
| Problem | Drivers who forget to go offline stay online and keep the vehicle bound. |
| Suggested solution | Auto-offline at shift end: store shift end time per driver, run a job that sets drivers offline (and unbinds vehicle) when current time is past shift end. |
| Shifts vs auto toggler | A single “shift end” time + job is simpler; full shifts (B) only if you need full scheduling and history. |

If you want to proceed with Option A, the next steps are: (1) add `shift_end_at` (or similar) and admin API in the backend, (2) add the scheduled job that calls the same logic as `UserOffline`, and (3) add the “Shift end” field in the admin driver edit UI.
