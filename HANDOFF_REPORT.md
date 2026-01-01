# 🚗 GreenRide Admin Dashboard - Developer Handoff Report

> **Last Updated:** January 1, 2026  
> **Project Status:** Demo Mode Fully Functional  
> **Branch:** `feature/api-integration`

---

## 📋 Executive Summary

The GreenRide Admin Dashboard is a **Next.js 16** web application for managing a ride-hailing platform operating in **Kigali, Rwanda**. The dashboard provides comprehensive tools for managing drivers, passengers, rides, revenue analytics, and real-time vehicle tracking.

The frontend UI is **100% complete** and operates in **Demo Mode** with mock data. Real backend API integration is pending AWS access for database credentials.

---

## 🏗️ Tech Stack

| Category | Technology | Version |
|----------|------------|---------|
| **Framework** | Next.js (App Router, Turbopack) | 16.1.1 |
| **Language** | TypeScript | 5.x |
| **Styling** | Tailwind CSS | 4.x |
| **UI Components** | shadcn/ui (Radix primitives) | Latest |
| **State Management** | Zustand | 5.x |
| **Data Fetching** | React Query (TanStack) | 5.x |
| **Charts** | Recharts | 3.x |
| **Maps** | Google Maps JavaScript API | Latest |
| **Forms** | React Hook Form + Zod | Latest |
| **Icons** | Lucide React | 0.562.0 |

---

## 📁 Project Structure

```
greenride-admin-v.2/
├── src/
│   ├── app/
│   │   ├── (auth)/
│   │   │   └── login/page.tsx          # Login page with demo mode
│   │   ├── (dashboard)/
│   │   │   ├── layout.tsx              # Dashboard layout with auth guard
│   │   │   ├── page.tsx                # Home dashboard with stats
│   │   │   ├── quick-booking/page.tsx  # CRM phone booking flow ⭐ NEW
│   │   │   ├── drivers/
│   │   │   │   ├── page.tsx            # Driver list with full CRUD
│   │   │   │   └── [id]/page.tsx       # Driver detail page
│   │   │   ├── users/
│   │   │   │   ├── page.tsx            # User list with full CRUD
│   │   │   │   └── [id]/page.tsx       # User detail page
│   │   │   ├── rides/
│   │   │   │   ├── page.tsx            # Rides/Orders list
│   │   │   │   └── [id]/page.tsx       # Ride detail page
│   │   │   ├── map/page.tsx            # Live map with Google Maps JS API
│   │   │   ├── revenue/page.tsx        # Revenue dashboard
│   │   │   ├── analytics/page.tsx      # Analytics charts
│   │   │   ├── promotions/page.tsx     # Promotions management
│   │   │   ├── notifications/page.tsx  # Notification center
│   │   │   └── settings/page.tsx       # Settings page
│   │   ├── globals.css                 # Global styles
│   │   └── layout.tsx                  # Root layout
│   ├── components/
│   │   ├── charts/                     # 8 chart components (Area, Bar, Pie, etc.)
│   │   ├── layout/
│   │   │   ├── sidebar.tsx             # Main navigation sidebar
│   │   │   ├── header.tsx              # Top header bar
│   │   │   └── mobile-sidebar.tsx      # Mobile navigation
│   │   └── ui/                         # shadcn/ui components
│   ├── lib/
│   │   ├── api-client.ts               # API client with demo mode ⭐ KEY FILE
│   │   └── utils.ts                    # Utility functions
│   ├── stores/
│   │   ├── auth-store.ts               # Authentication state
│   │   └── sidebar-store.ts            # Sidebar UI state
│   └── types/
│       └── index.ts                    # TypeScript type definitions
├── .env.local                          # Environment variables
├── package.json                        # Dependencies & scripts
├── HANDOFF_REPORT.md                   # This document
└── INTEGRATION_PROGRESS.md             # Integration phase tracking
```

---

## ✅ Feature Implementation Status

### **Core Features**

| Feature | Status | Notes |
|---------|--------|-------|
| **Authentication** | ✅ Complete | Demo mode bypass, JWT structure ready |
| **Dashboard Home** | ✅ Complete | Stats cards, charts, recent activity |
| **Driver Management** | ✅ Complete | List, Add, Edit, Delete, Suspend/Activate |
| **User Management** | ✅ Complete | List, Add, Edit, Delete, Suspend/Activate |
| **Ride Management** | ✅ Complete | List, View details, status display |
| **Revenue Dashboard** | ✅ Complete | Charts, trends, payment breakdown |
| **Analytics** | ✅ Complete | Multiple chart types, time filters |
| **Live Map** | ✅ Complete | Google Maps JS API, 30 mock drivers |
| **Promotions** | ✅ Complete | Full CRUD, Duplicate, Toggle, Export CSV |
| **Notifications** | ✅ Complete | List view with filters |
| **Settings** | ✅ Complete | Profile, preferences, security tabs |

### **New Features (This Session)**

| Feature | Status | Notes |
|---------|--------|-------|
| **Quick Booking (CRM)** | ✅ Complete | 4-step flow for phone bookings |
| **Driver CRUD Modals** | ✅ Complete | Add/Edit/Delete with validation |
| **User CRUD Modals** | ✅ Complete | Add/Edit/Delete with validation |
| **CSV Export** | ✅ Complete | Client-side download for tables |
| **PDF Export** | ✅ Complete | Revenue & Analytics printable reports |
| **View All Links** | ✅ Complete | All navigation links functional |
| **Auto Lock Cleanup** | ✅ Complete | `npm run dev` auto-cleans stale locks |

### **Pending Features**

| Feature | Priority | Blocker |
|---------|----------|---------|
| Real API Integration | High | AWS SSH access needed |

---

## 🔧 Environment Setup

### **Environment Variables (`.env.local`)**

```bash
# API Configuration
NEXT_PUBLIC_API_URL=http://13.247.190.134:8611

# Demo Mode (set to 'false' for real API)
NEXT_PUBLIC_DEMO_MODE=true

# Google Maps API Key
NEXT_PUBLIC_GOOGLE_MAPS_KEY=AIzaSyDif39v3Gx4YXonS3-A8pINUMi3hxRfC3U
```

### **NPM Scripts**

```bash
npm run dev        # Start dev server (auto-cleans lock file)
npm run dev:fresh  # Clean build + start (clears .next folder)
npm run build      # Production build
npm run lint       # Run ESLint
```

### **Known Issues & Solutions**

| Issue | Solution |
|-------|----------|
| `Unable to acquire lock` error | Run `npm run dev` (auto-cleans), or manually delete `.next/dev/lock` |
| Port 3000 in use | Next.js auto-selects next available port (3001, 3002, etc.) |
| Google Maps API errors | Ensure billing is enabled on Google Cloud Console |

---

## 🔌 Backend API Integration

### **API Base URL**
- **Development:** `http://13.247.190.134:8611` (Admin API)
- **User API:** Port 8610 (for mobile app)

### **Authentication**
- JWT token stored in `localStorage` as `admin_token`
- Header format: `Authorization: Bearer {token}`

### **Response Format**
```json
{
  "code": "0000",    // "0000" = success
  "msg": "Success",
  "data": { ... }
}
```

### **Key Endpoints**

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/login` | POST | Admin login |
| `/logout` | POST | Admin logout |
| `/info` | GET | Get admin info |
| `/users/search` | POST | Search users/drivers |
| `/users/create` | POST | Create user/driver |
| `/users/update` | POST | Update user/driver |
| `/users/status` | POST | Change user status |
| `/orders/search` | POST | Search rides |
| `/orders/detail` | POST | Get ride details |
| `/dashboard/stats` | GET | Dashboard statistics |

### **Demo Mode**

When `NEXT_PUBLIC_DEMO_MODE=true`:
- All API calls return mock data from `src/lib/api-client.ts`
- CRUD operations update in-memory mock arrays
- No real backend connection needed

To switch to real API:
1. Set `NEXT_PUBLIC_DEMO_MODE=false` in `.env.local`
2. Set `DEMO_MODE = false` in `src/lib/api-client.ts` (line 13)
3. Ensure valid admin credentials exist in database

---

## 🗄️ Database Access (Pending)

### **Required for Full Integration**

Access to the MySQL database is needed to:
1. Create admin accounts with valid credentials
2. View/modify user and driver data
3. Access ride history and payment records

### **Database Credentials (from CEO)**
```
Host: 13.247.190.134
Port: 3306
Database: green_ride
Username: admin
Password: !Greenride#@2024$
```

### **SSH Access Required**
AWS Console access or SSH key is needed to connect to the database server. Contact the CEO/previous developer for:
- AWS IAM credentials, OR
- SSH private key for the EC2 instance

---

## 🗺️ Live Map Configuration

### **Google Maps APIs Required**
- ✅ Maps JavaScript API
- ✅ Maps Embed API (fallback)
- ✅ Geocoding API
- ✅ Directions API
- ✅ Roads API
- ✅ Places API

### **Current Implementation**
- Uses `@react-google-maps/api` library
- 30 mock drivers with simulated movement
- Custom markers with status colors (green/yellow/gray)
- InfoWindow popups for driver details
- Map type toggle (Roadmap/Satellite/Hybrid)

### **For Real Integration**
- Replace mock drivers with real driver locations from WebSocket/polling
- Integrate with driver tracking backend endpoint
- Add real-time ride tracking

---

## 🎯 Quick Booking Feature (CRM)

### **Purpose**
Allows CRM agents to create ride bookings for customers who call the toll-free number (no app installed).

### **4-Step Flow**
1. **Find/Create Passenger** - Search by phone, create if new
2. **Enter Locations** - Pickup and drop-off addresses
3. **Assign Driver** - Select from nearby available drivers
4. **Confirm Booking** - Review and confirm

### **Access**
- Sidebar: "Quick Booking" link (with phone icon)
- URL: `/quick-booking`

---

## 📊 Mock Data Structure

### **Mock Drivers** (`api-client.ts`)
```typescript
{
  id: number,
  user_id: string,      // "DRV001", "DRV002", etc.
  full_name: string,
  display_name: string,
  first_name: string,
  last_name: string,
  email: string,
  phone: string,
  status: 'active' | 'suspended' | 'banned',
  online_status: 'online' | 'busy' | 'offline',
  rating: number,
  score: number,
  total_rides: number,
  created_at: number    // Unix timestamp
}
```

### **Mock Users** (`api-client.ts`)
```typescript
{
  id: number,
  user_id: string,      // "USR001", "USR002", etc.
  full_name: string,
  first_name: string,
  last_name: string,
  email: string,
  phone: string,
  status: 'active' | 'suspended' | 'banned',
  is_phone_verified: boolean,
  is_email_verified: boolean,
  total_rides: number,
  created_at: number
}
```

---

## 🚀 Getting Started (For New Developer)

### **1. Clone & Install**
```bash
git clone https://github.com/ngabonzizaguy/greenride-admin-v.2.git
cd greenride-admin-v.2
npm install --legacy-peer-deps
```

### **2. Configure Environment**
```bash
cp .env.example .env.local  # Or create manually
# Edit .env.local with your values
```

### **3. Start Development**
```bash
npm run dev
# Open http://localhost:3000 (or next available port)
```

### **4. Login**
- Demo Mode: Any username/password works
- Real Mode: Use valid admin credentials from database

---

## 📝 Next Steps for Continuation

### **Priority 1: AWS Access**
- [ ] Get AWS Console access or SSH key from CEO
- [ ] Connect to MySQL database
- [ ] Create/verify admin credentials
- [ ] Test real API endpoints

### **Priority 2: Real API Integration**
- [ ] Set `DEMO_MODE = false` in `api-client.ts`
- [ ] Test authentication flow
- [ ] Test CRUD operations with real data
- [ ] Verify dashboard stats from real backend

### **Priority 3: Live Features**
- [ ] Implement real driver location tracking
- [ ] Add WebSocket for real-time updates
- [ ] Connect Quick Booking to real order creation

### **Priority 4: Polish**
- [ ] Make Promotions CRUD functional
- [ ] Add PDF export for reports
- [ ] Fix any remaining "View All" links
- [ ] Add loading states for slow connections

---

## 📞 Contacts & Resources

| Resource | Details |
|----------|---------|
| **Repository** | `https://github.com/ngabonzizaguy/greenride-admin-v.2` |
| **Branch** | `feature/api-integration` |
| **Backend API** | `http://13.247.190.134:8611` |
| **Google Cloud** | Enable billing for Maps API |
| **CEO Contact** | For AWS access and credentials |

---

## 📄 Related Documents

- `INTEGRATION_PROGRESS.md` - Phase-by-phase integration tracking
- `BACKEND_API_EXTRACTION.md` - API documentation from Go backend
- `README.md` - Basic project readme

---

*This handoff report was generated on December 31, 2025. For questions or clarifications, refer to the commit history and inline code comments.*

