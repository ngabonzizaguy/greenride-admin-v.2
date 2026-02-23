# What to do next (deployment)

**I didn’t push anything.** All changes are only on your machine until you commit and push.

---

## 1. Save and push the changes

In your project folder (`greenride-admin-v.2-1`):

1. **Commit** the changes (e.g. in VS Code: Source Control → stage all → message: “Deployment fixes and docs” → Commit).
2. **Push** to GitHub (e.g. Push button, or: `git push origin main`).

That’s it for “publishing” the updates. The rest depends on how you deploy.

---

## 2. Deploy the new version

You have two options.

### Option A: Automatic (if you use GitHub Actions)

- After you **push to `main`**, GitHub Actions will deploy to your server.
- You don’t need to run any commands on the server.
- Check the **Actions** tab on GitHub to see if the deploy succeeded.

### Option B: Manual (you run a script on the server)

1. **Open a terminal on the server** (e.g. SSH: `ssh ubuntu@18.143.118.157`).
2. **Go to the project folder:**  
   `cd /home/ubuntu/greenride-admin-v.2`
3. **Update the code:**  
   `git pull origin main`
4. **Run the deploy script:**  
   `./deploy.sh`
   - This updates and restarts both backend and admin.
   - To update only backend: `./deploy.sh backend`  
   - Only admin: `./deploy.sh frontend`

---

## 3. Check that it worked

- **Admin dashboard:** Open your admin URL in the browser (e.g. the one that uses `admin-api.greenrideafrica.com` or your server IP).
- **Maintenance toggle:** In Settings → System, turn maintenance mode ON then OFF again; it should save without the “Invalid parameters” error.

---

**Summary:** Commit and push your changes, then either let GitHub deploy automatically or run `git pull` and `./deploy.sh` on the server. No other steps are required for the deployment changes we made.
