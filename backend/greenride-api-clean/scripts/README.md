# Utility Scripts

This directory contains standalone utility scripts for the GreenRide backend.

## Scripts

### `create_admin.go`
Creates or resets the production admin user.

**Usage:**
```bash
cd scripts
go run create_admin.go
```

**Creates:**
- Username: `admin`
- Password: `admin123`
- Email: `admin@greenrideafrica.com`
- Role: Super Admin

### `create_dev_admin.go`
Creates or resets a development admin user.

**Usage:**
```bash
cd scripts
go run create_dev_admin.go
```

**Creates:**
- Username: `devadmin`
- Password: `password123`
- Email: `dev@greenrideafrica.com`
- Role: Super Admin

### `test_hash.go`
Tests password hashing and verification logic.

**Usage:**
```bash
cd scripts
go run test_hash.go
```

## Why These Are Separate

These scripts are in a separate directory because they each have their own `main()` function. In Go, you can only have one `main()` function per package. By placing them in the `scripts/` directory, they can be run independently without conflicting with the main application's `main()` function in `main/main.go`.
