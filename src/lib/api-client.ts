/**
 * GreenRide Admin API Client
 * 
 * Connects to the Go backend Admin API (port 8611)
 * Based on BACKEND_API_EXTRACTION.md
 */

// API Base URL - defaults to development server
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://18.143.118.157:8611';

// Demo mode - returns mock data instead of real API calls
// Hardcoded to true for local testing until real credentials are available
const DEMO_MODE = true; // process.env.NEXT_PUBLIC_DEMO_MODE === 'true';

// ============================================
// MOCK DATA FOR DEMO MODE
// ============================================
const MOCK_DASHBOARD_STATS = {
  active_rides: 23,
  online_drivers: 156,
  today_revenue: 245000,
  today_rides: 342,
  pending_payments: 12,
  cancellation_rate: 4.2,
  total_users: 15420,
  total_drivers: 892,
};

const MOCK_DRIVERS = [
  { id: 1, user_id: 'DRV001', full_name: 'Peter Mutombo', email: 'peter.m@email.com', phone: '+250788123456', avatar: '', user_type: 'driver', status: 'active', online_status: 'online', verification_status: 'verified', rating: 4.8, total_rides: 1250, created_at: Date.now() - 90 * 24 * 60 * 60 * 1000 },
  { id: 2, user_id: 'DRV002', full_name: 'David Kagame', email: 'david.k@email.com', phone: '+250788234567', avatar: '', user_type: 'driver', status: 'active', online_status: 'busy', verification_status: 'verified', rating: 4.6, total_rides: 890, created_at: Date.now() - 60 * 24 * 60 * 60 * 1000 },
  { id: 3, user_id: 'DRV003', full_name: 'Paul Rwema', email: 'paul.r@email.com', phone: '+250788345678', avatar: '', user_type: 'driver', status: 'active', online_status: 'offline', verification_status: 'verified', rating: 4.9, total_rides: 2100, created_at: Date.now() - 120 * 24 * 60 * 60 * 1000 },
  { id: 4, user_id: 'DRV004', full_name: 'James Tuyisenge', email: 'james.t@email.com', phone: '+250788456789', avatar: '', user_type: 'driver', status: 'suspended', online_status: 'offline', verification_status: 'verified', rating: 3.9, total_rides: 450, created_at: Date.now() - 45 * 24 * 60 * 60 * 1000 },
  { id: 5, user_id: 'DRV005', full_name: 'Alex Munyaneza', email: 'alex.m@email.com', phone: '+250788567890', avatar: '', user_type: 'driver', status: 'active', online_status: 'online', verification_status: 'pending', rating: 4.5, total_rides: 320, created_at: Date.now() - 30 * 24 * 60 * 60 * 1000 },
];

const MOCK_USERS = [
  { id: 1, user_id: 'USR001', full_name: 'John Doe', email: 'john.doe@email.com', phone: '+250788111111', avatar: '', user_type: 'passenger', status: 'active', verification_status: 'verified', total_rides: 45, created_at: Date.now() - 100 * 24 * 60 * 60 * 1000 },
  { id: 2, user_id: 'USR002', full_name: 'Jane Smith', email: 'jane.smith@email.com', phone: '+250788222222', avatar: '', user_type: 'passenger', status: 'active', verification_status: 'verified', total_rides: 120, created_at: Date.now() - 80 * 24 * 60 * 60 * 1000 },
  { id: 3, user_id: 'USR003', full_name: 'Mike Johnson', email: 'mike.j@email.com', phone: '+250788333333', avatar: '', user_type: 'passenger', status: 'active', verification_status: 'pending', total_rides: 12, created_at: Date.now() - 15 * 24 * 60 * 60 * 1000 },
  { id: 4, user_id: 'USR004', full_name: 'Sarah Wilson', email: 'sarah.w@email.com', phone: '+250788444444', avatar: '', user_type: 'passenger', status: 'suspended', verification_status: 'verified', total_rides: 78, created_at: Date.now() - 150 * 24 * 60 * 60 * 1000 },
  { id: 5, user_id: 'USR005', full_name: 'Chris Brown', email: 'chris.b@email.com', phone: '+250788555555', avatar: '', user_type: 'passenger', status: 'active', verification_status: 'verified', total_rides: 200, created_at: Date.now() - 200 * 24 * 60 * 60 * 1000 },
];

const MOCK_RIDES = [
  { id: 1, order_id: 'ORD001', user_id: 'USR001', provider_id: 'DRV001', pickup_location: 'Kigali Convention Centre', dropoff_location: 'Kigali International Airport', status: 'completed', payment_status: 'paid', amount: 5200, distance: 12.5, duration: 25, created_at: Date.now() - 2 * 60 * 60 * 1000 },
  { id: 2, order_id: 'ORD002', user_id: 'USR002', provider_id: 'DRV002', pickup_location: 'Nyarutarama', dropoff_location: 'Downtown Kigali', status: 'in_progress', payment_status: 'pending', amount: 3800, distance: 8.2, duration: 18, created_at: Date.now() - 30 * 60 * 1000 },
  { id: 3, order_id: 'ORD003', user_id: 'USR003', provider_id: 'DRV003', pickup_location: 'Kimihurura', dropoff_location: 'Remera', status: 'completed', payment_status: 'paid', amount: 2500, distance: 5.0, duration: 12, created_at: Date.now() - 4 * 60 * 60 * 1000 },
  { id: 4, order_id: 'ORD004', user_id: 'USR004', provider_id: 'DRV001', pickup_location: 'Gisozi', dropoff_location: 'Kacyiru', status: 'cancelled', payment_status: 'refunded', amount: 0, distance: 0, duration: 0, created_at: Date.now() - 6 * 60 * 60 * 1000 },
  { id: 5, order_id: 'ORD005', user_id: 'USR005', provider_id: 'DRV005', pickup_location: 'Kabuga', dropoff_location: 'Kigali Heights', status: 'completed', payment_status: 'paid', amount: 6100, distance: 15.0, duration: 35, created_at: Date.now() - 1 * 60 * 60 * 1000 },
];

// Response codes from backend
export const API_CODES = {
  SUCCESS: '0000',
  PARAM_ERROR: '2001',
  AUTH_ERROR: '3000',
  BUSINESS_ERROR: '1003',
  SYSTEM_ERROR: '1000',
} as const;

// Standard API response format from backend
export interface ApiResponse<T = unknown> {
  code: string;
  msg: string;
  data: T;
}

// Paginated response format
export interface PageResult<T> {
  result_type: string;
  size: number;
  current: number;
  total: number;
  count: number;
  records: T[];
  attach?: Record<string, unknown>;
}

// Search/pagination request
export interface SearchRequest {
  keyword?: string;
  page?: number;
  limit?: number;
  user_type?: 'passenger' | 'driver';
  status?: string;
  online_status?: string;
  is_email_verified?: boolean;
  is_phone_verified?: boolean;
  is_active?: boolean;
}

// Order search request
export interface OrderSearchRequest {
  keyword?: string;
  page?: number;
  limit?: number;
  order_id?: string;
  order_type?: string;
  status?: string;
  payment_status?: string;
  user_id?: string;
  provider_id?: string;
  start_date?: number;
  end_date?: number;
  min_amount?: number;
  max_amount?: number;
}

// Vehicle search request
export interface VehicleSearchRequest {
  keyword?: string;
  page?: number;
  limit?: number;
  driver_id?: string;
  status?: string;
  verify_status?: string;
  category?: string;
  level?: string;
}

interface ApiOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';
  body?: unknown;
  headers?: Record<string, string>;
}

class ApiClient {
  private baseUrl: string;
  private token: string | null = null;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  setToken(token: string | null) {
    this.token = token;
    if (typeof window !== 'undefined') {
      if (token) {
        localStorage.setItem('admin_token', token);
      } else {
        localStorage.removeItem('admin_token');
      }
    }
  }

  getToken(): string | null {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('admin_token');
    }
    return this.token;
  }

  /**
   * Make an API request and handle the standard response format
   */
  private async request<T>(endpoint: string, options: ApiOptions = {}): Promise<ApiResponse<T>> {
    const { method = 'GET', body, headers = {} } = options;

    const token = this.getToken();
    const requestHeaders: Record<string, string> = {
      'Content-Type': 'application/json',
      'Accept-Language': 'en',
      ...headers,
    };

    if (token) {
      requestHeaders['Authorization'] = `Bearer ${token}`;
    }

    try {
      const response = await fetch(`${this.baseUrl}${endpoint}`, {
        method,
        headers: requestHeaders,
        body: body ? JSON.stringify(body) : undefined,
      });

      const data: ApiResponse<T> = await response.json();

      // Check for authentication errors
      if (data.code === API_CODES.AUTH_ERROR) {
        // Clear token and redirect to login
        this.setToken(null);
        if (typeof window !== 'undefined') {
          window.location.href = '/login';
        }
        throw new ApiError('Authentication failed', data.code, data.msg);
      }

      // Check for other errors
      if (data.code !== API_CODES.SUCCESS) {
        throw new ApiError(data.msg || 'Request failed', data.code, data.msg);
      }

      return data;
    } catch (error) {
      if (error instanceof ApiError) {
        throw error;
      }
      // Network or parsing error
      throw new ApiError(
        error instanceof Error ? error.message : 'Network error',
        API_CODES.SYSTEM_ERROR,
        'Unable to connect to server'
      );
    }
  }

  // ============================================
  // AUTHENTICATION ENDPOINTS
  // ============================================

  /**
   * Admin login
   * POST /login
   */
  async login(username: string, password: string): Promise<ApiResponse<{ token: string; user: unknown }>> {
    const response = await this.request<{ token: string; user: unknown }>('/login', {
      method: 'POST',
      body: { username, password },
    });
    
    // Store token on successful login
    if (response.data?.token) {
      this.setToken(response.data.token);
    }
    
    return response;
  }

  /**
   * Admin logout
   * POST /logout
   */
  async logout(): Promise<ApiResponse<null>> {
    try {
      const response = await this.request<null>('/logout', { method: 'POST' });
      return response;
    } finally {
      this.setToken(null);
    }
  }

  /**
   * Get current admin info
   * GET /info
   */
  async getAdminInfo(): Promise<ApiResponse<unknown>> {
    return this.request('/info');
  }

  /**
   * Change password
   * POST /change-password
   */
  async changePassword(oldPassword: string, newPassword: string): Promise<ApiResponse<null>> {
    return this.request('/change-password', {
      method: 'POST',
      body: { old_password: oldPassword, new_password: newPassword },
    });
  }

  // ============================================
  // DASHBOARD ENDPOINTS
  // ============================================

  /**
   * Get dashboard statistics
   * GET /dashboard/stats
   */
  async getDashboardStats(): Promise<ApiResponse<unknown>> {
    if (DEMO_MODE) {
      return { code: API_CODES.SUCCESS, msg: 'Success', data: MOCK_DASHBOARD_STATS };
    }
    return this.request('/dashboard/stats');
  }

  /**
   * Get revenue chart data
   * GET /dashboard/revenue
   */
  async getRevenueChart(): Promise<ApiResponse<unknown>> {
    return this.request('/dashboard/revenue');
  }

  /**
   * Get user growth chart data
   * GET /dashboard/user-growth
   */
  async getUserGrowthChart(): Promise<ApiResponse<unknown>> {
    return this.request('/dashboard/user-growth');
  }

  // ============================================
  // USER MANAGEMENT ENDPOINTS
  // ============================================

  /**
   * Search users (passengers or drivers)
   * POST /users/search
   */
  async searchUsers(params: SearchRequest = {}): Promise<ApiResponse<PageResult<unknown>>> {
    if (DEMO_MODE) {
      const isDriver = params.user_type === 'driver';
      const mockData = isDriver ? MOCK_DRIVERS : MOCK_USERS;
      let filtered = [...mockData];
      
      // Apply keyword filter
      if (params.keyword) {
        const kw = params.keyword.toLowerCase();
        filtered = filtered.filter(u => 
          u.full_name.toLowerCase().includes(kw) || 
          u.email.toLowerCase().includes(kw) ||
          u.phone.includes(kw)
        );
      }
      
      // Apply status filter
      if (params.status && params.status !== 'all') {
        filtered = filtered.filter(u => u.status === params.status);
      }
      
      // Apply online_status filter for drivers
      if (isDriver && params.online_status && params.online_status !== 'all') {
        filtered = filtered.filter((u: typeof MOCK_DRIVERS[0]) => u.online_status === params.online_status);
      }
      
      const page = params.page || 1;
      const limit = params.limit || 10;
      const start = (page - 1) * limit;
      const records = filtered.slice(start, start + limit);
      
      return {
        code: API_CODES.SUCCESS,
        msg: 'Success',
        data: {
          result_type: isDriver ? 'drivers' : 'users',
          size: limit,
          current: page,
          total: Math.ceil(filtered.length / limit),
          count: filtered.length,
          records,
          attach: isDriver ? {
            total_count: MOCK_DRIVERS.length,
            online_count: MOCK_DRIVERS.filter(d => d.online_status === 'online').length,
            busy_count: MOCK_DRIVERS.filter(d => d.online_status === 'busy').length,
            suspended_count: MOCK_DRIVERS.filter(d => d.status === 'suspended').length,
          } : undefined,
        },
      };
    }
    return this.request('/users/search', {
      method: 'POST',
      body: {
        page: params.page || 1,
        limit: params.limit || 10,
        ...params,
      },
    });
  }

  /**
   * Get drivers list
   * POST /users/search with user_type: 'driver'
   */
  async getDrivers(params: Omit<SearchRequest, 'user_type'> = {}): Promise<ApiResponse<PageResult<unknown>>> {
    return this.searchUsers({ ...params, user_type: 'driver' });
  }

  /**
   * Get passengers list
   * POST /users/search with user_type: 'passenger'
   */
  async getUsers(params: Omit<SearchRequest, 'user_type'> = {}): Promise<ApiResponse<PageResult<unknown>>> {
    return this.searchUsers({ ...params, user_type: 'passenger' });
  }

  /**
   * Get user details
   * POST /users/detail
   */
  async getUserDetail(userId: string): Promise<ApiResponse<unknown>> {
    if (DEMO_MODE) {
      const driver = MOCK_DRIVERS.find(d => d.user_id === userId || String(d.id) === userId);
      const user = MOCK_USERS.find(u => u.user_id === userId || String(u.id) === userId);
      const found = driver || user;
      if (found) {
        return { code: API_CODES.SUCCESS, msg: 'Success', data: found };
      }
      // Return first mock driver/user as fallback
      return { code: API_CODES.SUCCESS, msg: 'Success', data: MOCK_DRIVERS[0] };
    }
    return this.request('/users/detail', {
      method: 'POST',
      body: { user_id: userId },
    });
  }

  /**
   * Create a new user
   * POST /users/create
   */
  async createUser(userData: Record<string, unknown>): Promise<ApiResponse<unknown>> {
    return this.request('/users/create', {
      method: 'POST',
      body: userData,
    });
  }

  /**
   * Update user
   * POST /users/update
   */
  async updateUser(userId: string, userData: Record<string, unknown>): Promise<ApiResponse<unknown>> {
    return this.request('/users/update', {
      method: 'POST',
      body: { user_id: userId, ...userData },
    });
  }

  /**
   * Update user status (activate, suspend, ban)
   * POST /users/status
   */
  async updateUserStatus(userId: string, status: string): Promise<ApiResponse<unknown>> {
    return this.request('/users/status', {
      method: 'POST',
      body: { user_id: userId, status },
    });
  }

  /**
   * Verify user
   * POST /users/verify
   */
  async verifyUser(userId: string): Promise<ApiResponse<unknown>> {
    return this.request('/users/verify', {
      method: 'POST',
      body: { user_id: userId },
    });
  }

  /**
   * Get user's rides/orders
   * POST /users/rides
   */
  async getUserRides(userId: string, params: { page?: number; limit?: number } = {}): Promise<ApiResponse<PageResult<unknown>>> {
    return this.request('/users/rides', {
      method: 'POST',
      body: { user_id: userId, page: params.page || 1, limit: params.limit || 10 },
    });
  }

  // ============================================
  // VEHICLE MANAGEMENT ENDPOINTS
  // ============================================

  /**
   * Search vehicles
   * POST /vehicles/search
   */
  async searchVehicles(params: VehicleSearchRequest = {}): Promise<ApiResponse<PageResult<unknown>>> {
    return this.request('/vehicles/search', {
      method: 'POST',
      body: {
        page: params.page || 1,
        limit: params.limit || 10,
        ...params,
      },
    });
  }

  /**
   * Get vehicle details
   * POST /vehicles/detail
   */
  async getVehicleDetail(vehicleId: string): Promise<ApiResponse<unknown>> {
    return this.request('/vehicles/detail', {
      method: 'POST',
      body: { vehicle_id: vehicleId },
    });
  }

  /**
   * Create vehicle
   * POST /vehicles/create
   */
  async createVehicle(vehicleData: Record<string, unknown>): Promise<ApiResponse<unknown>> {
    return this.request('/vehicles/create', {
      method: 'POST',
      body: vehicleData,
    });
  }

  /**
   * Update vehicle
   * POST /vehicles/update
   */
  async updateVehicle(vehicleId: string, vehicleData: Record<string, unknown>): Promise<ApiResponse<unknown>> {
    return this.request('/vehicles/update', {
      method: 'POST',
      body: { vehicle_id: vehicleId, ...vehicleData },
    });
  }

  /**
   * Update vehicle status
   * POST /vehicles/status
   */
  async updateVehicleStatus(vehicleId: string, status: string): Promise<ApiResponse<unknown>> {
    return this.request('/vehicles/status', {
      method: 'POST',
      body: { vehicle_id: vehicleId, status },
    });
  }

  /**
   * Delete vehicle
   * POST /vehicles/delete
   */
  async deleteVehicle(vehicleId: string): Promise<ApiResponse<null>> {
    return this.request('/vehicles/delete', {
      method: 'POST',
      body: { vehicle_id: vehicleId },
    });
  }

  // ============================================
  // ORDER MANAGEMENT ENDPOINTS
  // ============================================

  /**
   * Search orders/rides
   * POST /orders/search
   */
  async searchOrders(params: OrderSearchRequest = {}): Promise<ApiResponse<PageResult<unknown>>> {
    if (DEMO_MODE) {
      let filtered = [...MOCK_RIDES];
      
      // Apply keyword filter
      if (params.keyword) {
        const kw = params.keyword.toLowerCase();
        filtered = filtered.filter(r => 
          r.order_id.toLowerCase().includes(kw) || 
          r.pickup_location.toLowerCase().includes(kw) ||
          r.dropoff_location.toLowerCase().includes(kw)
        );
      }
      
      // Apply status filter
      if (params.status && params.status !== 'all') {
        filtered = filtered.filter(r => r.status === params.status);
      }
      
      const page = params.page || 1;
      const limit = params.limit || 10;
      const start = (page - 1) * limit;
      const records = filtered.slice(start, start + limit);
      
      return {
        code: API_CODES.SUCCESS,
        msg: 'Success',
        data: {
          result_type: 'orders',
          size: limit,
          current: page,
          total: Math.ceil(filtered.length / limit),
          count: filtered.length,
          records,
        },
      };
    }
    return this.request('/orders/search', {
      method: 'POST',
      body: {
        page: params.page || 1,
        limit: params.limit || 10,
        ...params,
      },
    });
  }

  /**
   * Get order details
   * POST /orders/detail
   */
  async getOrderDetail(orderId: string): Promise<ApiResponse<unknown>> {
    if (DEMO_MODE) {
      const order = MOCK_RIDES.find(r => r.order_id === orderId || String(r.id) === orderId);
      if (order) {
        return { code: API_CODES.SUCCESS, msg: 'Success', data: order };
      }
      // Return first mock order as fallback
      return { code: API_CODES.SUCCESS, msg: 'Success', data: MOCK_RIDES[0] };
    }
    return this.request('/orders/detail', {
      method: 'POST',
      body: { order_id: orderId },
    });
  }

  /**
   * Estimate order
   * POST /orders/estimate
   */
  async estimateOrder(orderData: Record<string, unknown>): Promise<ApiResponse<unknown>> {
    return this.request('/orders/estimate', {
      method: 'POST',
      body: orderData,
    });
  }

  /**
   * Create order
   * POST /orders/create
   */
  async createOrder(orderData: Record<string, unknown>): Promise<ApiResponse<unknown>> {
    return this.request('/orders/create', {
      method: 'POST',
      body: orderData,
    });
  }

  /**
   * Cancel order
   * POST /orders/cancel
   */
  async cancelOrder(orderId: string, reason?: string): Promise<ApiResponse<null>> {
    return this.request('/orders/cancel', {
      method: 'POST',
      body: { order_id: orderId, reason },
    });
  }
}

/**
 * Custom error class for API errors
 */
export class ApiError extends Error {
  code: string;
  serverMessage: string;

  constructor(message: string, code: string, serverMessage: string) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
    this.serverMessage = serverMessage;
  }

  isAuthError(): boolean {
    return this.code === API_CODES.AUTH_ERROR;
  }

  isParamError(): boolean {
    return this.code === API_CODES.PARAM_ERROR;
  }
}

// Export singleton instance
export const apiClient = new ApiClient(API_BASE_URL);
export default apiClient;
