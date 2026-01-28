import { create } from 'zustand';
import { apiClient } from '@/lib/api-client';

export interface Notification {
  id: number;
  notification_id: string;
  type: string;
  category: string;
  title: string;
  content: string;
  summary?: string;
  status: string;
  priority: string;
  is_read: boolean;
  is_archived: boolean;
  user_id?: string;
  user_type?: string;
  created_at: number;
  updated_at: number;
  sent_at?: number;
  read_at?: number;
  delivered_at?: number;
}

interface NotificationStore {
  notifications: Notification[];
  unreadCount: number;
  isLoading: boolean;
  error: string | null;
  
  // Actions
  fetchNotifications: (params?: { page?: number; limit?: number; keyword?: string; type?: string; status?: string }) => Promise<void>;
  markAsRead: (notificationId: string) => Promise<void>;
  markAllAsRead: () => Promise<void>;
  refreshUnreadCount: () => Promise<void>;
}

export const useNotificationStore = create<NotificationStore>((set, get) => ({
  notifications: [],
  unreadCount: 0,
  isLoading: false,
  error: null,

  fetchNotifications: async (params = {}) => {
    set({ isLoading: true, error: null });
    try {
      const response = await apiClient.getNotifications({
        page: params.page || 1,
        limit: params.limit || 20,
        keyword: params.keyword,
        type: params.type,
        status: params.status,
      });

      if (response.code === '0000' && response.data) {
        const pageResult = response.data as { records?: Notification[]; total?: number };
        set({
          notifications: pageResult.records || [],
          isLoading: false,
        });
      } else {
        set({ error: response.msg || 'Failed to fetch notifications', isLoading: false });
      }
    } catch (error) {
      console.error('Failed to fetch notifications:', error);
      set({ error: 'Failed to fetch notifications', isLoading: false });
    }
  },

  markAsRead: async (notificationId: string) => {
    try {
      const response = await apiClient.markNotificationAsRead(notificationId);
      if (response.code === '0000') {
        // Update local state
        set((state) => ({
          notifications: state.notifications.map((n) =>
            n.notification_id === notificationId
              ? { ...n, is_read: true, read_at: Date.now() }
              : n
          ),
          unreadCount: Math.max(0, state.unreadCount - 1),
        }));
      }
    } catch (error) {
      console.error('Failed to mark notification as read:', error);
    }
  },

  markAllAsRead: async () => {
    try {
      const response = await apiClient.markAllNotificationsAsRead();
      if (response.code === '0000') {
        // Update local state
        set((state) => ({
          notifications: state.notifications.map((n) => ({
            ...n,
            is_read: true,
            read_at: n.read_at || Date.now(),
          })),
          unreadCount: 0,
        }));
      }
    } catch (error) {
      console.error('Failed to mark all notifications as read:', error);
    }
  },

  refreshUnreadCount: async () => {
    try {
      const response = await apiClient.getUnreadNotificationCount();
      if (response.code === '0000' && response.data) {
        const data = response.data as { count?: number };
        set({ unreadCount: data.count || 0 });
      }
    } catch (error) {
      console.error('Failed to refresh unread count:', error);
    }
  },
}));
