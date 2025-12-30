/**
 * GreenRide Admin API Client
 * 
 * Connects to the Go backend Admin API (port 8611)
 * Based on BACKEND_API_EXTRACTION.md
 */

// API Base URL - defaults to development server
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://18.143.118.157:8611';

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
