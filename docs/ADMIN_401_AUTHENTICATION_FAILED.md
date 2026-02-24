# Admin 401 "Authentication failed" – Investigation

## What you’re seeing

- **Browser:** Requests to `https://admin-api.greenrideafrica.com` return **401 Unauthorized**.
- **Response body:** `{"code":"3000","msg":"Authentication failed"}`.
- **Console:** "API Client" logs show `hasToken: true` but the server still rejects the request.

So the frontend **is** sending a token; the backend is **rejecting** it.

---

## Did recent code changes break the admin?

**No.** In this project we did **not** change:

- Admin frontend auth (login, token storage, API client auth headers).
- Backend admin routes or JWT validation.
- Any file under `src/` that affects how the token is sent or validated.

Only docs and the **app** (green_ride_app) payment flow were changed. The admin codebase (Next.js, api-client, auth) was not modified.

So the 401 is **not** caused by the recent payment or doc changes. It is due to one of the causes below.

---

## Likely causes (in order)

1. **Expired or invalid token**
   - The stored admin token may have expired or been invalidated.
   - **Fix:** Log out in the admin UI and log in again to get a new token. If you use a long‑lived token (e.g. API key), regenerate it and update the place it’s stored.

2. **Backend JWT secret or issuer changed**
   - If the backend was redeployed or config changed (e.g. new JWT secret, different issuer), existing tokens will no longer validate.
   - **Fix:** Ensure the backend’s JWT config (secret, issuer, audience) matches what was used when the token was issued. Then log in again from the admin to get a new token.

3. **Wrong base URL or path**
   - If the admin was built with a different `NEXT_PUBLIC_API_URL` (e.g. with or without `/admin/api`), it might be calling a different server or path that expects different auth.
   - **Fix:** Confirm the admin is calling the same base URL you use in production (e.g. `https://admin-api.greenrideafrica.com` or `https://admin-api.greenrideafrica.com/admin/api`). Check browser Network tab: request URL and `Authorization` header are present.

4. **CORS or preflight**
   - Less common, but if the browser sends a preflight and the server responds with 401 on the preflight, the console can show 401.
   - **Fix:** Ensure the backend allows the admin origin and sends proper CORS headers for the admin domain.

---

## What to do next

1. **Log out and log in again** in the admin panel to refresh the token.
2. **Inspect one failing request** in DevTools → Network: confirm the request URL and that the `Authorization: Bearer <token>` header is present and unchanged after login.
3. **Check backend logs** for the admin API when you trigger a request: look for auth errors (e.g. "token expired", "invalid signature", "invalid issuer").
4. **Confirm backend config** for admin JWT (secret, issuer, audience) has not changed since the last successful login.

Adding more logging (e.g. in the API client) for 401 responses (without logging the raw token) can help confirm that a token is sent and that the server returns 401; it does not change the fact that the server is rejecting the token for one of the reasons above.
