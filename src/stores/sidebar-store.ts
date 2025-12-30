import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface SidebarState {
  isCollapsed: boolean;
  isMobileOpen: boolean;
  toggleCollapsed: () => void;
  setCollapsed: (collapsed: boolean) => void;
  toggleMobile: () => void;
  setMobileOpen: (open: boolean) => void;
}

export const useSidebarStore = create<SidebarState>()(
  persist(
    (set) => ({
      isCollapsed: false,
      isMobileOpen: false,
      toggleCollapsed: () => set((state) => ({ isCollapsed: !state.isCollapsed })),
      setCollapsed: (isCollapsed) => set({ isCollapsed }),
      toggleMobile: () => set((state) => ({ isMobileOpen: !state.isMobileOpen })),
      setMobileOpen: (isMobileOpen) => set({ isMobileOpen }),
    }),
    {
      name: 'sidebar-storage',
      partialize: (state) => ({ isCollapsed: state.isCollapsed }),
    }
  )
);
