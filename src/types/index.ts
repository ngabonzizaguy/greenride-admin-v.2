// User/Driver types
export interface User {
  id: string;
  name: string;
  email: string;
  phone: string;
  avatar?: string;
  role: 'passenger' | 'driver' | 'admin';
  status: 'active' | 'suspended' | 'deleted';
  createdAt: Date;
  updatedAt: Date;
}

export interface Driver extends User {
  role: 'driver';
  rating: number;
  totalTrips: number;
  acceptanceRate: number;
  completionRate: number;
  isOnline: boolean;
  currentLocation?: {
    lat: number;
    lng: number;
    updatedAt: Date;
  };
  vehicle?: Vehicle;
}

export interface Vehicle {
  id: string;
  driverId: string;
  plateNumber: string;
  model: string;
  make: string;
  year: number;
  color: string;
  type: 'sedan' | 'suv' | 'moto' | 'premium';
  photos: string[];
  documents: VehicleDocument[];
}

export interface VehicleDocument {
  type: 'license' | 'insurance' | 'registration';
  url: string;
  expiryDate: Date;
  status: 'valid' | 'expired' | 'pending';
}

// Ride types
export interface Ride {
  id: string;
  passengerId: string;
  passenger: User;
  driverId?: string;
  driver?: Driver;
  status: RideStatus;
  pickup: Location;
  dropoff: Location;
  distance: number; // in km
  duration: number; // in minutes
  fare: number;
  paymentMethod: 'cash' | 'momo' | 'card';
  paymentStatus: 'pending' | 'paid' | 'failed' | 'refunded';
  rating?: number;
  review?: string;
  createdAt: Date;
  startedAt?: Date;
  completedAt?: Date;
}

export type RideStatus = 
  | 'requesting'
  | 'driver_assigned'
  | 'driver_arriving'
  | 'driver_arrived'
  | 'trip_started'
  | 'trip_ended'
  | 'completed'
  | 'cancelled';

export interface Location {
  address: string;
  lat: number;
  lng: number;
}

// Financial types
export interface Transaction {
  id: string;
  rideId: string;
  amount: number;
  currency: 'RWF';
  method: 'cash' | 'momo' | 'card';
  status: 'pending' | 'completed' | 'failed' | 'refunded';
  reference?: string;
  createdAt: Date;
}

// Analytics types
export interface DashboardStats {
  activeRides: number;
  onlineDrivers: number;
  todayRevenue: number;
  todayRides: number;
  pendingPayments: number;
  cancellationRate: number;
}

export interface RevenueData {
  date: string;
  revenue: number;
  rides: number;
  cash: number;
  momo: number;
  card: number;
}

// Promotion types
export interface Promotion {
  id: string;
  code: string;
  type: 'percentage' | 'fixed' | 'free_ride';
  value: number;
  minOrderAmount?: number;
  maxDiscount?: number;
  usageLimit?: number;
  usedCount: number;
  validFrom: Date;
  validUntil: Date;
  status: 'active' | 'expired' | 'disabled';
}

// Admin types
export interface AdminUser {
  id: string;
  name: string;
  email: string;
  role: AdminRole;
  permissions: Permission[];
  lastLogin?: Date;
  createdAt: Date;
}

export type AdminRole = 'super_admin' | 'operations' | 'finance' | 'support';
export type Permission = 'read' | 'write' | 'delete' | 'export' | 'settings';

// Activity types
export interface Activity {
  id: string;
  type: 'ride_completed' | 'ride_cancelled' | 'driver_online' | 'driver_offline' | 'payment_received' | 'new_user';
  description: string;
  userId?: string;
  driverId?: string;
  rideId?: string;
  timestamp: Date;
}

// Notification types
export interface Notification {
  id: string;
  type: 'ride' | 'payment' | 'driver' | 'user' | 'system';
  title: string;
  message: string;
  read: boolean;
  createdAt: Date;
}
