# GreenRide Deployment

How to deploy the **backend API** (Go) and **admin frontend** (Next.js) for GreenRide.

**→ For a short, step-by-step “what do I do next?” see [DEPLOY_WHAT_TO_DO.md](DEPLOY_WHAT_TO_DO.md).**

---

## Quick reference

| Method | When to use | Command / trigger |
|--------|--------------|--------------------|
| **Manual (EC2)** | Deploy from your machine to the server | SSH to server, run `./deploy.sh` (or `./deploy.sh backend` / `./deploy.sh frontend`) |
| **CI/CD** | Deploy on every push to `main` | Push to `main` → GitHub Actions deploys to Singapore server |

**Server (current):** `18.143.118.157` (ubuntu)  
**Paths on server:**

- Repo: `/home/ubuntu/greenride-admin-v.2`
- Backend config/logs: `/home/ubuntu/greenride-api`
- API port: `8610`, Admin API port: `8611`, Frontend port: `3601`

---

## 1. Prerequisites on server

- Docker installed and running
- Git
- For manual deploy: repo cloned at `/home/ubuntu/greenride-admin-v.2`
- Backend directory created: `mkdir -p /home/ubuntu/greenride-api`
- `prod.yaml` and `config.yaml` in `/home/ubuntu/greenride-api` (created once; `deploy.sh` does not overwrite `prod.yaml`)

---

## 2. Manual deployment (`deploy.sh`)

Run **on the server** (e.g. after SSH):

```bash
cd /home/ubuntu/greenride-admin-v.2
chmod +x deploy.sh
./deploy.sh           # backend + frontend
./deploy.sh backend   # backend only
./deploy.sh frontend  # frontend only
```

What it does:

- **Backend:** Pulls latest `main`, builds the Go binary (or uses Docker to build), copies binary + config + locales to `/home/ubuntu/greenride-api`, builds the Docker image **from the repo** (full source), runs the container with volumes for `config.yaml`, `prod.yaml`, and `logs`.
- **Frontend:** Builds the Next.js Docker image from the repo root and runs it on port 3601.

Backend image is built from source so the Dockerfile in the repo works without a separate “runtime” image.

---

## 3. CI/CD (GitHub Actions)

- **Workflow:** `.github/workflows/deploy.yml`
- **Trigger:** Push to `main`
- **Secrets:** `SSH_PRIVATE_KEY` (private key for `ubuntu@18.143.118.157`)

Steps:

1. **Frontend:** SSH → pull, build and run the admin container.
2. **Backend:** SSH → pull, copy configs, build image from repo, restart container.

Backend is built from source on the server (same as manual deploy).

---

## 4. Backend image build

- **Dockerfile:** `backend/greenride-api-clean/Dockerfile` (multi-stage: build Go binary, then minimal runtime).
- **Config:** Container expects `config.yaml` and `prod.yaml` (and optionally `internal/locales`) mounted or copied. Production secrets live in `prod.yaml` on the server; do not commit real secrets.
- **Ports:** 8610 (user API), 8611 (admin API), 8612.

To build locally (for testing):

```bash
cd /path/to/greenride-admin-v.2
docker build -t greenride-api -f backend/greenride-api-clean/Dockerfile backend/greenride-api-clean
```

---

## 5. Frontend (admin dashboard) image

- **Dockerfile:** Repo root `Dockerfile`.
- **Build-time env:** `NEXT_PUBLIC_API_URL`, `NEXT_PUBLIC_DEMO_MODE`, `NEXT_PUBLIC_GOOGLE_MAPS_KEY` (set in Dockerfile or override with build args).
- **Port:** Container listens on 3000; host maps to 3601.

---

## 6. Environment and config

| Component | Config file | Notes |
|-----------|-------------|--------|
| Backend | `config.yaml` | Base config; committed (no secrets). |
| Backend | `prod.yaml` | **Production only**; DB, Redis, JWT, Firebase, etc. Not overwritten by deploy script. |
| Frontend | Dockerfile / build args | `NEXT_PUBLIC_*` baked at build time. |

For production admin URL and API, see `PRODUCTION_DEPLOYMENT_GUIDE.md`.

---

## 7. Health checks

After deploy:

```bash
# User API
curl -s http://localhost:8610/health

# Admin API
curl -s http://localhost:8611/health

# Frontend (optional)
curl -s -o /dev/null -w "%{http_code}" http://localhost:3601
```

---

## 8. Troubleshooting

| Issue | What to check |
|-------|----------------|
| Backend container exits | `docker logs greenride-api`; ensure `prod.yaml` and `config.yaml` exist and DB/Redis are reachable. |
| Frontend 502 / blank | `docker logs greenride-admin-v2`; confirm `NEXT_PUBLIC_API_URL` matches the admin API (e.g. `https://admin-api.greenrideafrica.com`). |
| “Invalid parameters” on system config | Admin frontend was sending double-stringified JSON; fixed in `src/lib/api-client.ts` (use `body: config`, not `JSON.stringify(config)`). |
| CORS errors | Nginx or API CORS allowlist must include the admin frontend origin (e.g. `https://admin.greenrideafrica.com`). |

---

## 9. Related docs

- **Production checklist and env:** `PRODUCTION_DEPLOYMENT_GUIDE.md`
- **Backend scripts:** `backend/greenride-api-clean/scripts/README.md`
- **Deployment test script:** `test-deployment.sh` (run against a live server to validate health and maintenance mode).
