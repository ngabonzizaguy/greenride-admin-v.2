'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { cn } from '@/lib/utils';
import { useSidebarStore } from '@/stores/sidebar-store';
import { 
  LayoutDashboard, 
  Map, 
  Users, 
  Car, 
  DollarSign, 
  BarChart3, 
  Tag, 
  Bell, 
  Settings,
  ChevronLeft,
  LogOut,
  Leaf,
  PhoneCall,
  MessageSquare
} from 'lucide-react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { useAuthStore } from '@/stores/auth-store';
import { apiClient } from '@/lib/api-client';

// Dynamic nav items - badges will be updated from API
const getNavItems = (counts: { feedback: number; drivers: number; rides: number }) => [
  { href: '/', icon: LayoutDashboard, label: 'Dashboard' },
  { href: '/quick-booking', icon: PhoneCall, label: 'Quick Booking', highlight: true },
  { href: '/feedback', icon: MessageSquare, label: 'Feedback', badge: counts.feedback > 0 ? String(counts.feedback) : undefined },
  { href: '/map', icon: Map, label: 'Live Map' },
  { href: '/drivers', icon: Car, label: 'Drivers', badge: counts.drivers > 0 ? String(counts.drivers) : undefined },
  { href: '/users', icon: Users, label: 'Users' },
  { href: '/rides', icon: Car, label: 'Rides', badge: counts.rides > 0 ? String(counts.rides) : undefined },
  { href: '/revenue', icon: DollarSign, label: 'Revenue' },
  { href: '/analytics', icon: BarChart3, label: 'Analytics' },
  { href: '/promotions', icon: Tag, label: 'Promotions' },
  { href: '/notifications', icon: Bell, label: 'Notifications', hasNotification: true },
];

const bottomNavItems = [
  { href: '/settings', icon: Settings, label: 'Settings' },
];

export function Sidebar() {
  const pathname = usePathname();
  const { isCollapsed, toggleCollapsed } = useSidebarStore();
  const { user, logout } = useAuthStore();
  const [counts, setCounts] = useState({ feedback: 0, drivers: 0, rides: 0 });

  // Fetch counts from API
  useEffect(() => {
    const fetchCounts = async () => {
      try {
        const statsResponse = await apiClient.getDashboardStats();
        if (statsResponse.code === '0000' && statsResponse.data) {
          const data = statsResponse.data as Record<string, unknown>;
          setCounts({
            feedback: (data.pending_feedback as number) ?? 0,
            drivers: (data.online_drivers as number) ?? 0,
            rides: (data.active_trips as number) ?? 0,
          });
        }
      } catch (error) {
        // Silently fail - badges will just not show
        // Timeout errors are expected if backend is slow - don't spam console
        if (error instanceof Error && !error.message.includes('timeout')) {
          console.error('Failed to fetch sidebar counts:', error);
        }
      }
    };

    fetchCounts();
    // Refresh every 60 seconds
    const interval = setInterval(fetchCounts, 60000);
    return () => clearInterval(interval);
  }, []);

  const navItems = getNavItems(counts);

  return (
    <TooltipProvider delayDuration={0}>
      <aside
        className={cn(
          'fixed left-0 top-0 z-40 h-screen glass-sidebar text-gray-300 transition-all duration-300 ease-in-out border-r border-white/10',
          isCollapsed ? 'w-[72px]' : 'w-[260px]'
        )}
      >
        <div className="flex h-full flex-col">
          {/* Header */}
          <div className={cn(
            'flex h-16 items-center border-b border-white/10 px-6',
            isCollapsed ? 'justify-center' : 'justify-start'
          )}>
            <Link href="/" className="flex items-center gap-2">
              <span className="text-primary filter drop-shadow-[0_0_8px_rgba(16,185,129,0.5)]">
                <Leaf className="h-6 w-6" />
              </span>
              {!isCollapsed && (
                <span className="text-xl font-bold text-white tracking-tight drop-shadow-sm">GreenRide</span>
              )}
            </Link>
          </div>

          {/* Navigation */}
          <nav className="flex-1 space-y-1 px-4 py-6">
            {!isCollapsed && (
                <p className="px-4 text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Menu</p>
            )}
            {navItems.map((item) => {
              const isActive = pathname === item.href || 
                (item.href !== '/' && pathname.startsWith(item.href));
              
              const linkContent = (
                <Link
                  href={item.href}
                  className={cn(
                    'flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-all duration-200 group relative',
                    isActive
                      ? 'bg-primary/20 text-white border border-primary/30 shadow-[0_0_15px_rgba(16,185,129,0.15)]'
                      : 'hover:bg-white/10 hover:text-white',
                    isCollapsed && 'justify-center px-2'
                  )}
                >
                  <item.icon className={cn(
                    "mr-3 h-5 w-5 transition-opacity", 
                    isActive ? "text-primary drop-shadow-[0_0_5px_rgba(16,185,129,0.8)] opacity-100" : "opacity-70 group-hover:opacity-100",
                    isCollapsed && "mr-0"
                  )} />
                  {!isCollapsed && (
                    <>
                      <span className="flex-1">{item.label}</span>
                      {item.badge && (
                        <span className={cn(
                            "ml-auto text-xs font-bold px-2 py-0.5 rounded-full shadow-sm",
                            isActive ? "bg-primary text-white" : "bg-white/10 text-gray-300"
                        )}>
                          {item.badge}
                        </span>
                      )}
                      {item.hasNotification && (
                        <span className="h-2 w-2 rounded-full bg-red-500 absolute top-2 right-2 ring-2 ring-white/10" />
                      )}
                    </>
                  )}
                </Link>
              );

              if (isCollapsed) {
                return (
                  <Tooltip key={item.href}>
                    <TooltipTrigger asChild>{linkContent}</TooltipTrigger>
                    <TooltipContent side="right" className="flex items-center gap-2">
                      {item.label}
                    </TooltipContent>
                  </Tooltip>
                );
              }

              return <div key={item.href}>{linkContent}</div>;
            })}
          </nav>

          {/* Divider */}
          <div className="mx-4 border-t border-white/10" />

          {/* Bottom Navigation */}
          <nav className="space-y-1 px-4 py-4">
            {bottomNavItems.map((item) => {
              const isActive = pathname === item.href;
              
              const linkContent = (
                <Link
                  href={item.href}
                  className={cn(
                    'flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-all duration-200 mb-1 group',
                    isActive
                      ? 'bg-primary/20 text-white border border-primary/30'
                      : 'hover:bg-white/10 hover:text-white',
                    isCollapsed && 'justify-center px-2'
                  )}
                >
                  <item.icon className={cn(
                    "mr-3 h-5 w-5 transition-opacity", 
                    isActive ? "text-primary opacity-100" : "opacity-70 group-hover:opacity-100",
                    isCollapsed && "mr-0"
                  )} />
                  {!isCollapsed && <span>{item.label}</span>}
                </Link>
              );

              if (isCollapsed) {
                return (
                  <Tooltip key={item.href}>
                    <TooltipTrigger asChild>{linkContent}</TooltipTrigger>
                    <TooltipContent side="right">{item.label}</TooltipContent>
                  </Tooltip>
                );
              }

              return <div key={item.href}>{linkContent}</div>;
            })}
          </nav>

          {/* User Section */}
          <div className={cn(
            'px-4 py-4 border-t border-white/10',
            isCollapsed && 'flex justify-center'
          )}>
            {isCollapsed ? (
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-10 w-10 hover:bg-white/10"
                    onClick={toggleCollapsed}
                  >
                    <div className="h-8 w-8 rounded-full bg-gradient-to-br from-primary to-emerald-600 flex items-center justify-center text-white font-bold text-xs shadow-md">
                        {user?.full_name?.charAt(0) || user?.username?.charAt(0) || 'A'}
                    </div>
                  </Button>
                </TooltipTrigger>
                <TooltipContent side="right">
                  <p className="font-medium">{user?.full_name || user?.username || 'Admin User'}</p>
                </TooltipContent>
              </Tooltip>
            ) : (
              <div className="flex items-center px-4 py-3 mt-2 rounded-lg bg-white/5 border border-white/5 hover:bg-white/10 transition-colors cursor-pointer group">
                <div className="h-8 w-8 rounded-full bg-gradient-to-br from-primary to-emerald-600 flex items-center justify-center text-white font-bold text-xs mr-3 shadow-md">
                    {user?.full_name?.charAt(0) || user?.username?.charAt(0) || 'A'}
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium text-white truncate">{user?.full_name || user?.username || 'Admin User'}</p>
                  <p className="text-xs text-gray-400 truncate">{user?.role?.replace('_', ' ') || 'Super Admin'}</p>
                </div>
                <Tooltip>
                    <TooltipTrigger asChild>
                        <Button
                            variant="ghost"
                            size="icon"
                            className="h-8 w-8 text-gray-400 hover:text-white"
                            onClick={logout}
                        >
                            <LogOut className="h-4 w-4" />
                        </Button>
                    </TooltipTrigger>
                    <TooltipContent>Logout</TooltipContent>
                </Tooltip>
              </div>
            )}
          </div>
        </div>
      </aside>
    </TooltipProvider>
  );
}
