// API Client for GreenRide Backend
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'https://api.greenrideafrica.com';

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
  }

  getToken(): string | null {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('admin_token');
    }
    return this.token;
  }

  private async request<T>(endpoint: string, options: ApiOptions = {}): Promise<T> {
    const { method = 'GET', body, headers = {} } = options;

    const token = this.getToken();
    const requestHeaders: Record<string, string> = {
      'Content-Type': 'application/json',
      ...headers,
    };

    if (token) {
      requestHeaders['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      method,
      headers: requestHeaders,
      body: body ? JSON.stringify(body) : undefined,
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: 'An error occurred' }));
      throw new Error(error.message || `HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  // Auth endpoints
  async login(email: string, password: string) {
    const response = await this.request<{ token: string; user: unknown }>('/admin/auth/login', {
      method: 'POST',
      body: { email, password },
    });
    if (typeof window !== 'undefined') {
      localStorage.setItem('admin_token', response.token);
    }
    this.token = response.token;
    return response;
  }

  async logout() {
    await this.request('/admin/auth/logout', { method: 'POST' });
    if (typeof window !== 'undefined') {
      localStorage.removeItem('admin_token');
    }
    this.token = null;
  }

  async getCurrentUser() {
    return this.request('/admin/auth/me');
  }

  // Dashboard endpoints
  async getDashboardStats() {
    return this.request('/admin/dashboard/stats');
  }

  async getRecentActivity() {
    return this.request('/admin/dashboard/activity');
  }

  // Driver endpoints
  async getDrivers(params?: { page?: number; limit?: number; status?: string; search?: string }) {
    const query = new URLSearchParams(params as Record<string, string>).toString();
    return this.request(`/admin/drivers${query ? `?${query}` : ''}`);
  }

  async getDriver(id: string) {
    return this.request(`/admin/drivers/${id}`);
  }

  async updateDriver(id: string, data: unknown) {
    return this.request(`/admin/drivers/${id}`, { method: 'PUT', body: data });
  }

  async suspendDriver(id: string) {
    return this.request(`/admin/drivers/${id}/suspend`, { method: 'POST' });
  }

  async activateDriver(id: string) {
    return this.request(`/admin/drivers/${id}/activate`, { method: 'POST' });
  }

  async getDriverTrips(id: string) {
    return this.request(`/admin/drivers/${id}/trips`);
  }

  async getDriverEarnings(id: string) {
    return this.request(`/admin/drivers/${id}/earnings`);
  }

  async getDriverLocations() {
    return this.request('/admin/drivers/locations');
  }

  // User endpoints
  async getUsers(params?: { page?: number; limit?: number; search?: string }) {
    const query = new URLSearchParams(params as Record<string, string>).toString();
    return this.request(`/admin/users${query ? `?${query}` : ''}`);
  }

  async getUser(id: string) {
    return this.request(`/admin/users/${id}`);
  }

  async suspendUser(id: string) {
    return this.request(`/admin/users/${id}/suspend`, { method: 'POST' });
  }

  // Ride endpoints
  async getRides(params?: { page?: number; limit?: number; status?: string; date?: string }) {
    const query = new URLSearchParams(params as Record<string, string>).toString();
    return this.request(`/admin/rides${query ? `?${query}` : ''}`);
  }

  async getRide(id: string) {
    return this.request(`/admin/rides/${id}`);
  }

  async getActiveRides() {
    return this.request('/admin/rides/active');
  }

  async cancelRide(id: string, reason: string) {
    return this.request(`/admin/rides/${id}/cancel`, { method: 'POST', body: { reason } });
  }

  // Revenue endpoints
  async getRevenueSummary(params?: { startDate?: string; endDate?: string }) {
    const query = new URLSearchParams(params as Record<string, string>).toString();
    return this.request(`/admin/revenue/summary${query ? `?${query}` : ''}`);
  }

  async getTransactions(params?: { page?: number; limit?: number }) {
    const query = new URLSearchParams(params as Record<string, string>).toString();
    return this.request(`/admin/revenue/transactions${query ? `?${query}` : ''}`);
  }

  async getRevenueByDate(params?: { startDate?: string; endDate?: string }) {
    const query = new URLSearchParams(params as Record<string, string>).toString();
    return this.request(`/admin/revenue/by-date${query ? `?${query}` : ''}`);
  }

  // Analytics endpoints
  async getPeakHours() {
    return this.request('/admin/analytics/peak-hours');
  }

  async getPopularRoutes() {
    return this.request('/admin/analytics/popular-routes');
  }

  async getDriverPerformance() {
    return this.request('/admin/analytics/driver-performance');
  }

  // Promotions endpoints
  async getPromotions() {
    return this.request('/admin/promotions');
  }

  async createPromotion(data: unknown) {
    return this.request('/admin/promotions', { method: 'POST', body: data });
  }

  async updatePromotion(id: string, data: unknown) {
    return this.request(`/admin/promotions/${id}`, { method: 'PUT', body: data });
  }

  async deletePromotion(id: string) {
    return this.request(`/admin/promotions/${id}`, { method: 'DELETE' });
  }

  // Notifications endpoints
  async sendNotification(data: { audience: string; title: string; message: string }) {
    return this.request('/admin/notifications/send', { method: 'POST', body: data });
  }

  async getNotificationHistory() {
    return this.request('/admin/notifications/history');
  }
}

export const apiClient = new ApiClient(API_BASE_URL);
export default apiClient;
