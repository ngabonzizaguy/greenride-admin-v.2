/**
 * GreenRide Admin Dashboard Type Definitions
 * 
 * Aligned with backend models from BACKEND_API_EXTRACTION.md
 * Note: Backend uses snake_case, we convert to camelCase on the frontend
 * All timestamps are in milliseconds (Unix epoch * 1000)
 */

// ============================================
// API RESPONSE TYPES
// ============================================

export interface ApiResponse<T = unknown> {
  code: string;
  msg: string;
  data: T;
}

export interface PageResult<T> {
  result_type: string;
  size: number;
  current: number;
  total: number;
  count: number;
  records: T[];
  attach?: Record<string, unknown>;
}

// ============================================
// USER TYPES (t_users)
// ============================================

export type UserType = 'passenger' | 'driver';
export type UserStatus = 'active' | 'inactive' | 'suspended' | 'banned';
export type OnlineStatus = 'online' | 'offline' | 'busy';
export type Gender = 'male' | 'female' | 'other';

export interface User {
  // Core identifiers
  id: number;
  user_id: string;
  
  // User type
  user_type: UserType;
  
  // Auth
  email?: string;
  phone?: string;
  country_code?: string;
  
  // Profile
  username?: string;
  display_name?: string;
  full_name?: string;
  first_name?: string;
  last_name?: string;
  avatar?: string;
  gender?: Gender;
  birthday?: number; // timestamp
  language?: string;
  timezone?: string;
  
  // Address
  address?: string;
  city?: string;
  state?: string;
  country?: string;
  postal_code?: string;
  
  // Location
  latitude?: number;
  longitude?: number;
  location_updated_at?: number;
  
  // Status
  status: UserStatus;
  is_email_verified?: boolean;
  is_phone_verified?: boolean;
  online_status?: OnlineStatus;
  
  // Driver-specific fields
  license_number?: string;
  license_expiry?: number;
  queued_order_ids?: string[];
  current_order_id?: string;
  next_available_at?: number;
  max_queue_capacity?: number;
  
  // Stats
  score?: number; // Rating out of 5.0
  total_rides?: number;
  
  // Referral
  invite_code?: string;
  invited_by?: string;
  invite_count?: number;
  
  // Device
  fcm_token?: string;
  device_id?: string;
  device_type?: 'ios' | 'android' | 'web';
  
  // Timestamps
  last_login_at?: number;
  email_verified_at?: number;
  phone_verified_at?: number;
  created_at: number;
  updated_at: number;
  deleted_at?: number;
  
  // Test flag
  sandbox?: number; // 0 = production, 1 = test
}

// Driver is a User with user_type = 'driver'
export interface Driver extends User {
  user_type: 'driver';
  vehicle?: Vehicle;
}

// Passenger is a User with user_type = 'passenger'  
export interface Passenger extends User {
  user_type: 'passenger';
}

// ============================================
// VEHICLE TYPES (t_vehicles)
// ============================================

export type VehicleCategory = 'sedan' | 'suv' | 'mpv' | 'van' | 'hatchback';
export type VehicleLevel = 'economy' | 'comfort' | 'premium' | 'luxury';
export type VehicleStatus = 'active' | 'inactive' | 'maintenance' | 'retired';
export type VerifyStatus = 'pending' | 'verified' | 'rejected';
export type FuelType = 'gasoline' | 'diesel' | 'electric' | 'hybrid';
export type Transmission = 'manual' | 'automatic';

export interface Vehicle {
  // Core identifiers
  id: number;
  vehicle_id: string;
  driver_id?: string;
  driver?: User;
  
  // Basic info
  brand?: string;
  model?: string;
  year?: number;
  color?: string;
  plate_number?: string;
  vin?: string;
  
  // Type and specs
  type_id?: string;
  category?: VehicleCategory;
  level?: VehicleLevel;
  seat_capacity?: number;
  fuel_type?: FuelType;
  transmission?: Transmission;
  
  // Status
  status: VehicleStatus;
  verify_status?: VerifyStatus;
  
  // Registration/Insurance
  registration_number?: string;
  registration_expiry?: number;
  insurance_company?: string;
  insurance_policy_number?: string;
  insurance_expiry?: number;
  
  // Location
  current_latitude?: number;
  current_longitude?: number;
  location_updated_at?: number;
  
  // Media
  photos?: string[];
  documents?: string[];
  
  // Rating
  rating?: number;
  
  // Timestamps
  created_at: number;
  updated_at: number;
}

// ============================================
// ORDER TYPES (t_orders) - "Rides" in UI
// ============================================

export type OrderType = 'ride' | 'delivery' | 'shopping';
export type OrderStatus = 
  | 'requested'
  | 'accepted'
  | 'arrived'
  | 'in_progress'
  | 'trip_ended'
  | 'completed'
  | 'cancelled';
export type PaymentStatus = 'pending' | 'success' | 'failed';
export type PaymentMethod = 'cash' | 'card' | 'wallet' | 'momo';
export type ScheduleType = 'instant' | 'scheduled';
export type DispatchStatus = 'not_started' | 'in_progress' | 'completed' | 'failed';

export interface Order {
  // Core identifiers
  id: number;
  order_id: string;
  
  // Type
  order_type?: OrderType;
  
  // User relations
  user_id?: string;      // Customer/Passenger ID
  provider_id?: string;  // Driver ID
  
  // Status
  status: OrderStatus;
  payment_status?: PaymentStatus;
  schedule_type?: ScheduleType;
  
  // Amounts (as strings for decimal precision)
  currency?: string;
  original_amount?: string;
  discounted_amount?: string;
  payment_amount?: string;
  total_discount_amount?: string;
  platform_fee?: string;
  
  // Payment
  payment_method?: PaymentMethod;
  payment_id?: string;
  channel_payment_id?: string;
  payment_result?: string;
  payment_redirect_url?: string;
  
  // Promotions
  promo_codes?: string[];
  promo_discount?: string;
  user_promotion_ids?: string[];
  
  // Timestamps
  scheduled_at?: number;
  accepted_at?: number;
  started_at?: number;
  ended_at?: number;
  completed_at?: number;
  cancelled_at?: number;
  expired_at?: number;
  created_at: number;
  updated_at: number;
  
  // Cancellation
  cancelled_by?: string;
  cancel_reason?: string;
  cancellation_fee?: string;
  
  // Dispatch
  dispatch_status?: DispatchStatus;
  current_round?: number;
  max_rounds?: number;
  auto_dispatch_enabled?: boolean;
  
  // Version for optimistic locking
  version?: number;
  sandbox?: number;
  
  // Joined data (from detail endpoint)
  user?: User;
  provider?: Driver;
  details?: OrderDetail;
}

// Order details (for ride orders)
export interface OrderDetail {
  pickup_address?: string;
  pickup_latitude?: number;
  pickup_longitude?: number;
  dropoff_address?: string;
  dropoff_latitude?: number;
  dropoff_longitude?: number;
  distance?: number;   // in km
  duration?: number;   // in minutes
  route_polyline?: string;
}

// Alias for UI - "Ride" is more user-friendly than "Order"
export type Ride = Order;
export type RideStatus = OrderStatus;

// ============================================
// ADMIN TYPES (t_admins)
// ============================================

export type AdminRole = 'super_admin' | 'admin' | 'moderator' | 'support' | 'analyst';
export type AdminStatus = 'active' | 'inactive' | 'suspended';
export type AdminActiveStatus = 'online' | 'offline';

export interface AdminUser {
  // Core identifiers
  id: number;
  admin_id: string;
  
  // Auth
  username?: string;
  email?: string;
  phone?: string;
  
  // Profile
  first_name?: string;
  last_name?: string;
  full_name?: string;
  avatar?: string;
  
  // Role & Permissions
  role?: AdminRole;
  permissions?: string; // JSON array
  department?: string;
  job_title?: string;
  
  // Status
  status: AdminStatus;
  active_status?: AdminActiveStatus;
  
  // Login tracking
  last_login_at?: number;
  last_login_ip?: string;
  login_count?: number;
  failed_attempts?: number;
  locked_until?: number;
  
  // Timestamps
  created_at: number;
  updated_at: number;
}

// ============================================
// DASHBOARD & ANALYTICS TYPES
// ============================================

export interface DashboardStats {
  active_rides: number;
  online_drivers: number;
  today_revenue: number;
  today_rides: number;
  pending_payments: number;
  cancellation_rate: number;
  total_users: number;
  total_drivers: number;
}

export interface RevenueData {
  date: string;
  revenue: number;
  rides: number;
  cash: number;
  momo: number;
  card: number;
}

export interface UserGrowthData {
  date: string;
  passengers: number;
  drivers: number;
  total: number;
}

// ============================================
// PROMOTION TYPES (t_promotions)
// ============================================

export type PromotionType = 'percentage' | 'fixed' | 'free_ride';
export type PromotionStatus = 'active' | 'expired' | 'disabled';

export interface Promotion {
  id: number;
  promotion_id: string;
  code: string;
  type: PromotionType;
  value: number;
  min_order_amount?: number;
  max_discount?: number;
  usage_limit?: number;
  used_count: number;
  valid_from: number;
  valid_until: number;
  status: PromotionStatus;
  created_at: number;
  updated_at: number;
}

// ============================================
// NOTIFICATION TYPES (t_notifications)
// ============================================

export type NotificationType = 'ride' | 'payment' | 'driver' | 'user' | 'system';

export interface Notification {
  id: number;
  notification_id: string;
  type: NotificationType;
  title: string;
  message: string;
  read: boolean;
  user_id?: string;
  created_at: number;
}

// ============================================
// ACTIVITY TYPES
// ============================================

export type ActivityType = 
  | 'ride_completed' 
  | 'ride_cancelled' 
  | 'driver_online' 
  | 'driver_offline' 
  | 'payment_received' 
  | 'new_user'
  | 'new_driver';

export interface Activity {
  id: string;
  type: ActivityType;
  description: string;
  user_id?: string;
  driver_id?: string;
  order_id?: string;
  timestamp: number;
}

// ============================================
// TRANSACTION TYPES (t_payments)
// ============================================

export interface Transaction {
  id: number;
  payment_id: string;
  order_id: string;
  amount: string;
  currency: string;
  method: PaymentMethod;
  status: PaymentStatus;
  reference?: string;
  created_at: number;
}

// ============================================
// LOCATION TYPES
// ============================================

export interface Location {
  address: string;
  latitude: number;
  longitude: number;
}

export interface DriverLocation {
  driver_id: string;
  latitude: number;
  longitude: number;
  online_status: OnlineStatus;
  updated_at: number;
}

// ============================================
// FEEDBACK TYPES
// ============================================

export type FeedbackCategory = 'driver' | 'vehicle' | 'pricing' | 'safety' | 'app' | 'payment' | 'other';
export type FeedbackStatus = 'pending' | 'reviewing' | 'resolved' | 'closed';
export type FeedbackSeverity = 'low' | 'medium' | 'high' | 'critical';

export interface Feedback {
  id: string;
  feedback_id: string;
  order_id?: string;
  user_id: string;
  user_name: string;
  user_phone: string;
  driver_id?: string;
  driver_name?: string;
  category: FeedbackCategory;
  severity: FeedbackSeverity;
  title: string;
  content: string;
  rating?: number;
  attachments?: string[];
  status: FeedbackStatus;
  admin_response?: string;
  assigned_to?: string;
  created_at: number;
  updated_at: number;
  resolved_at?: number;
}

// ============================================
// CONFIG TYPES
// ============================================

export interface SupportConfig {
  phone: string;
  email: string;
  whatsapp?: string;
  hours: string;
  faq_url?: string;
}

// ============================================
// UTILITY TYPES
// ============================================

// For form inputs - allows partial updates
export type UserUpdateInput = Partial<Omit<User, 'id' | 'user_id' | 'created_at' | 'updated_at'>>;
export type VehicleUpdateInput = Partial<Omit<Vehicle, 'id' | 'vehicle_id' | 'created_at' | 'updated_at'>>;
export type OrderUpdateInput = Partial<Omit<Order, 'id' | 'order_id' | 'created_at' | 'updated_at'>>;

// Search params
export interface UserSearchParams {
  keyword?: string;
  page?: number;
  limit?: number;
  user_type?: UserType;
  status?: UserStatus;
  online_status?: OnlineStatus;
  is_email_verified?: boolean;
  is_phone_verified?: boolean;
}

export interface OrderSearchParams {
  keyword?: string;
  page?: number;
  limit?: number;
  order_id?: string;
  order_type?: OrderType;
  status?: OrderStatus;
  payment_status?: PaymentStatus;
  user_id?: string;
  provider_id?: string;
  start_date?: number;
  end_date?: number;
  min_amount?: number;
  max_amount?: number;
}

export interface VehicleSearchParams {
  keyword?: string;
  page?: number;
  limit?: number;
  driver_id?: string;
  status?: VehicleStatus;
  verify_status?: VerifyStatus;
  category?: VehicleCategory;
  level?: VehicleLevel;
}
