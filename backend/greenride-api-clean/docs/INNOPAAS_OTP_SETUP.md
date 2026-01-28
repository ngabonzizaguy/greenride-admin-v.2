# InnoPaaS OTP setup

## 1. Callback URL (InnoPaaS console)

When configuring API Keys / Callback in InnoPaaS, set:

**Callback address:**  
`https://api.greenrideafrica.com/webhook/innopaas`

- The backend exposes `POST /webhook/innopaas` and returns HTTP 200 within 3 seconds (required by InnoPaaS).
- Optional: enable "message status updated" (or Check All) to receive delivery status.

---

## 2. InnoPaaS “Basic Information” (Step 3) — what to put

- **Request Mode:** Leave **POST** selected (our webhook expects POST).
- **Username:** Leave **empty** (optional; we don’t use callback auth).
- **Authorization:** Leave **empty** (optional).

Then click **Submit** to save the callback configuration.

---

## 3. Where to put the keys

**Do not push real `app_key` or `app_secret` to the repo.** Keep placeholders in repo; put real values only on the server.

### Production (AWS server)

The deploy workflow **does not** overwrite `prod.yaml` on the server. **Important:** On the server, `prod.yaml` must be a **file**, not a directory. If you have an empty `prod.yaml` directory (e.g. `cd prod.yaml` then `ls` is empty), fix it once:

```bash
cd /home/ubuntu/greenride-admin-v.2
git pull origin main
rm -rf /home/ubuntu/greenride-api/prod.yaml
cp backend/greenride-api-clean/prod.yaml /home/ubuntu/greenride-api/prod.yaml
```

Then add your real InnoPaaS keys:

1. **SSH into the server** (e.g. `ssh ubuntu@18.143.118.157`).
2. **Edit** `prod.yaml` in the backend directory:
   ```bash
   nano /home/ubuntu/greenride-api/prod.yaml
   ```
   (If the file doesn’t exist, create it or copy from repo once, then edit.)
3. **Ensure an `innopaas:` block** exists and set real values:
   ```yaml
   innopaas:
     endpoint: "https://api.innopaas.com/api/otp/v3/msg/send/verify"
     app_key: "YOUR_REAL_API_KEY_FROM_INNOPAAS"
     app_secret: "YOUR_REAL_API_PASSWORD_FROM_INNOPAAS"
     sender_id: "GreenRide"
   ```
4. **Restart the API** so it reloads config:
   ```bash
   docker ps
   docker restart greenride-api
   ```
   If your API container has a different name (e.g. `api`), use that: `docker restart api`.

- **`app_key`** = value from InnoPaaS → Developer Tools → API Keys → **API Key** column (copy full value).
- **`app_secret`** = value from the same row → **API Password** column (copy via the copy icon).

### Local development

In the repo, **`backend/greenride-api-clean/local.yaml`** under `innopaas:` replace the placeholders with your real API Key and API Password. Do not commit real secrets (e.g. use a local override that is gitignored).

---

## 4. Deploy (no secrets in repo)

- **Push to `main`** only code and config **without** real InnoPaaS secrets (repo `prod.yaml` and `local.yaml` keep placeholders).
- GitHub Actions will deploy frontend and backend to the Singapore server.
- **After deploy**, SSH to the server and add real `app_key` and `app_secret` in **`/home/ubuntu/greenride-api/prod.yaml`** as in section 3, then restart the API.

---

## 5. Flow

- **OTP send:** Twilio is tried first (with 5s timeout); on failure, InnoPaaS is used (OTP v3 with `app_key` + `app_secret` signature).
- **Callback:** InnoPaaS can POST to `https://api.greenrideafrica.com/webhook/innopaas` for message status; the handler responds 200 and logs the body.
