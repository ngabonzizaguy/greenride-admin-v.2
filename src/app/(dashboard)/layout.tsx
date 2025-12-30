'use client';

import { Sidebar, Header, MobileSidebar } from '@/components/layout';
import { useSidebarStore } from '@/stores/sidebar-store';
import { cn } from '@/lib/utils';

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { isCollapsed } = useSidebarStore();

  return (
    <div className="min-h-screen bg-background">
      {/* Desktop Sidebar */}
      <div className="hidden md:block">
        <Sidebar />
      </div>

      {/* Mobile Sidebar */}
      <MobileSidebar />

      {/* Main Content */}
      <div
        className={cn(
          'flex flex-col transition-all duration-300',
          isCollapsed ? 'md:pl-[72px]' : 'md:pl-[260px]'
        )}
      >
        <Header />
        <main className="flex-1 p-4 md:p-6">{children}</main>
      </div>
    </div>
  );
}
