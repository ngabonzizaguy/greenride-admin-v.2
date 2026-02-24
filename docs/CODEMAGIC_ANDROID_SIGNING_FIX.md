# Codemagic: Android build "Keystore file not found" — clear steps

## What changed / why it broke

Your **codemagic.yaml** says:

```yaml
android_signing:
  - android_signing
```

In Codemagic, that means: **"Use the Android keystore whose reference name is `android_signing`"** from **Code signing identities**. Codemagic then sets `CM_KEYSTORE_PATH`, `CM_KEYSTORE_PASSWORD`, `CM_KEY_ALIAS`, `CM_KEY_PASSWORD` on the build machine.

What you showed is the **Environment variables** screen: you have a **variable group** named `android_signing` with only `ANDROID_SIGNING_PLACEHOLDER = not_used`. That is **not** the same as uploading a keystore file.

- **Environment variables** = key/value pairs (e.g. SHOREBIRD_TOKEN, placeholders). They do **not** provide a keystore file.
- **Code signing identities** = separate section where you **upload** the actual `.jks` file. When you reference it in yaml, Codemagic injects the keystore path and passwords.

So either the keystore was never uploaded in Code signing identities, or it was under a different reference name. The app's **build.gradle** was updated to use **CM_KEYSTORE_PATH** so once the keystore is correctly set up in Code signing, the build will find it.

---

## What you need to do (step by step)

### 1. Get your keystore file

Use the same `.jks` (or `.keystore`) you use for Play Store releases. If you don't have it, you cannot sign the same app for updates.

### 2. In Codemagic: add the keystore in Code signing identities

- Open Codemagic, your app **green-ride-app**.
- Go to **Code signing identities** (not the "Environment variables" page). It may be under Team settings or in the app settings.
- Under **Android keystores**, click **Add keystore** / **Upload keystore**.
- Set **Reference name:** e.g. `android_signing` (must match the name in yaml).
- Set **Keystore password**, **Key alias**, **Key password**.
- Upload the `.jks` file and save.

### 3. Keep codemagic.yaml as is

Your yaml has `android_signing: - android_signing`. The reference name in Code signing identities should be **android_signing**. If you used a different name (e.g. `release_keystore`), use that in yaml instead.

### 4. Optional: Environment variables

You can delete `ANDROID_SIGNING_PLACEHOLDER` from the `android_signing` variable group; it doesn't affect signing. Real signing comes from Code signing identities.

### 5. Re-run the Android build

Start a new build. Codemagic will set `CM_KEYSTORE_PATH` and Gradle will find the keystore.

---

## Summary

| Where | What to do |
|-------|------------|
| Code signing identities | Add Android keystore: upload `.jks`, reference name `android_signing`, passwords and alias. |
| codemagic.yaml | Keep `android_signing: - android_signing`. |
| Environment variables | No keystore needed; optional to remove placeholder. |
| build.gradle | Already uses `CM_KEYSTORE_PATH`. |

After the keystore is in Code signing identities with reference name `android_signing`, the "Keystore file not found" error should go away.
