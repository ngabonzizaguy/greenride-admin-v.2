import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import type { AdminUser } from '@/types';
import { apiClient } from '@/lib/api-client';

interface AuthState {
  user: AdminUser | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  setUser: (user: AdminUser | null) => void;
  setLoading: (loading: boolean) => void;
  logout: () => void;
  checkAuth: () => Promise<void>;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      isAuthenticated: false,
      isLoading: true,

      setUser: (user) => {
        set({ 
          user, 
          isAuthenticated: !!user, 
          isLoading: false 
        });
      },

      setLoading: (isLoading) => set({ isLoading }),

      logout: async () => {
        try {
          // Call logout API
          await apiClient.logout();
        } catch (error) {
          // Ignore errors on logout - we're clearing local state anyway
          console.warn('Logout API call failed:', error);
        }
        
        // Clear local state
        if (typeof window !== 'undefined') {
          localStorage.removeItem('admin_token');
        }
        set({ user: null, isAuthenticated: false, isLoading: false });
      },

      // Check if current token is valid by calling /info endpoint
      checkAuth: async () => {
        const token = typeof window !== 'undefined' 
          ? localStorage.getItem('admin_token') 
          : null;

        if (!token) {
          set({ user: null, isAuthenticated: false, isLoading: false });
          return;
        }

        try {
          // Verify token by fetching admin info
          const response = await apiClient.getAdminInfo();
          const userData = response.data as AdminUser;
          
          set({ 
            user: userData, 
            isAuthenticated: true, 
            isLoading: false 
          });
        } catch (error) {
          // Token is invalid or expired
          console.warn('Auth check failed:', error);
          if (typeof window !== 'undefined') {
            localStorage.removeItem('admin_token');
          }
          set({ user: null, isAuthenticated: false, isLoading: false });
        }
      },
    }),
    {
      name: 'greenride-auth',
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({ 
        user: state.user, 
        isAuthenticated: state.isAuthenticated 
      }),
      onRehydrateStorage: () => (state) => {
        // After rehydration from localStorage, set loading to false
        if (state) {
          state.isLoading = false;
        }
      },
    }
  )
);
