# GreenRide Africa - Admin Dashboard

A modern, responsive admin dashboard for the GreenRide Africa ride-hailing platform. Built with Next.js 14+, Tailwind CSS, and shadcn/ui.

## 🚀 Features

### Dashboard Home
- Real-time stats cards (Active Rides, Online Drivers, Revenue, etc.)
- Revenue trend charts
- Payment methods distribution chart
- Recent activity feed
- Quick actions panel

### Live Operations Map
- Interactive map view (Google Maps ready)
- Driver location markers with status colors
- Ride route visualization
- Real-time filtering controls

### Driver Management
- Comprehensive driver list with search and filters
- Driver detail pages with tabs:
  - Overview (stats & performance chart)
  - Trip history
  - Earnings breakdown
  - Reviews
  - Vehicle & documents
  - Activity log
- Actions: Edit, Suspend, Activate, Delete

### User Management
- User list with search and filters
- User detail pages with trip history
- Payment methods view
- Account actions

### Ride Management
- Ride list with status filters
- Ride detail pages with:
  - Route map visualization
  - Fare breakdown
  - Timeline
  - Passenger & driver info
  - Rating & review

### Revenue Dashboard
- Revenue stats cards
- Payment method breakdown
- Revenue trend charts
- Transaction history table
- Export functionality (PDF/Excel/CSV)

### Analytics
- Peak hours heatmap
- Rides by day of week
- Popular routes
- User growth chart
- Distance distribution

### Promotions
- Create and manage promo codes
- Usage tracking
- Status management (Active/Disabled/Expired)

### Notifications
- Broadcast messages to drivers/users
- Notification history
- Open rate tracking

### Settings
- Company information
- Pricing configuration
- Admin user management
- Security settings
- Notification preferences

## 🛠 Tech Stack

- **Framework**: Next.js 14+ (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **UI Components**: shadcn/ui
- **Charts**: Recharts
- **Icons**: Lucide Icons
- **State Management**: Zustand
- **Data Fetching**: React Query ready
- **Forms**: React Hook Form + Zod

## 📁 Project Structure

```
greenride-admin/
├── src/
│   ├── app/
│   │   ├── (auth)/
│   │   │   └── login/
│   │   ├── (dashboard)/
│   │   │   ├── drivers/
│   │   │   ├── users/
│   │   │   ├── rides/
│   │   │   ├── revenue/
│   │   │   ├── analytics/
│   │   │   ├── map/
│   │   │   ├── promotions/
│   │   │   ├── notifications/
│   │   │   └── settings/
│   │   └── layout.tsx
│   ├── components/
│   │   ├── ui/          # shadcn components
│   │   ├── charts/      # Recharts components
│   │   └── layout/      # Sidebar, Header
│   ├── lib/
│   │   ├── api-client.ts
│   │   └── utils.ts
│   ├── stores/
│   │   ├── auth-store.ts
│   │   └── sidebar-store.ts
│   └── types/
│       └── index.ts
├── public/
├── tailwind.config.ts
└── package.json
```

## 🚀 Getting Started

### Prerequisites
- Node.js 18+
- npm or yarn

### Installation

```bash
# Navigate to the admin dashboard directory
cd greenride-admin

# Install dependencies
npm install

# Run development server
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

### Build for Production

```bash
npm run build
npm run start
```

### Deployment

**What do I do next?** → **[docs/DEPLOY_WHAT_TO_DO.md](docs/DEPLOY_WHAT_TO_DO.md)** (short, step-by-step).  
More detail: [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md) and [PRODUCTION_DEPLOYMENT_GUIDE.md](PRODUCTION_DEPLOYMENT_GUIDE.md).

## 🔑 Authentication

For demo purposes, the login accepts any email and password. In production, connect to the GreenRide API:

```
Production API: https://api.greenrideafrica.com
Development API: http://18.143.118.157:8610/
```

Set the API URL in your environment:

```bash
NEXT_PUBLIC_API_URL=https://api.greenrideafrica.com
```

## 🎨 Theming

The dashboard uses GreenRide's brand colors:

- **Primary**: `#22C55E` (Green)
- **Background**: `#F8FAFC` (Light gray)
- **Sidebar**: `#0F172A` (Dark navy)

Colors can be customized in `src/app/globals.css`.

## 📱 Responsive Design

The dashboard is fully responsive:
- **Desktop**: Full sidebar, multi-column layouts
- **Tablet**: Collapsible sidebar, adapted grids
- **Mobile**: Bottom navigation, single-column layouts

## 🔗 API Integration

The `src/lib/api-client.ts` provides methods for all API endpoints:

- Authentication (login, logout, refresh)
- Dashboard stats
- Driver CRUD operations
- User management
- Ride operations
- Revenue & analytics
- Promotions
- Notifications
- Settings

## 🚧 TODO

- [ ] Connect to live GreenRide API
- [ ] Implement WebSocket for real-time updates
- [ ] Add Google Maps integration
- [ ] Implement export functionality
- [ ] Add dark mode toggle
- [ ] Set up unit tests
- [ ] Add E2E tests with Playwright

## 📄 License

Proprietary - GreenRide Africa

---

Built for GreenRide Africa 🌿
