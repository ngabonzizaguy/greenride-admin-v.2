'use client';

import { useEffect } from 'react';
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
  LogOut,
  Leaf,
  X
} from 'lucide-react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { useAuthStore } from '@/stores/auth-store';

const navItems = [
  { href: '/', icon: LayoutDashboard, label: 'Dashboard' },
  { href: '/map', icon: Map, label: 'Live Map' },
  { href: '/drivers', icon: Car, label: 'Drivers', badge: '24' },
  { href: '/users', icon: Users, label: 'Users' },
  { href: '/rides', icon: Car, label: 'Rides', badge: '12' },
  { href: '/revenue', icon: DollarSign, label: 'Revenue' },
  { href: '/analytics', icon: BarChart3, label: 'Analytics' },
  { href: '/promotions', icon: Tag, label: 'Promotions' },
  { href: '/notifications', icon: Bell, label: 'Notifications', hasNotification: true },
  { href: '/settings', icon: Settings, label: 'Settings' },
];

export function MobileSidebar() {
  const pathname = usePathname();
  const { isMobileOpen, setMobileOpen } = useSidebarStore();
  const { user, logout } = useAuthStore();

  // Close sidebar on route change
  useEffect(() => {
    setMobileOpen(false);
  }, [pathname, setMobileOpen]);

  // Prevent body scroll when sidebar is open
  useEffect(() => {
    if (isMobileOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = 'unset';
    }
    return () => {
      document.body.style.overflow = 'unset';
    };
  }, [isMobileOpen]);

  if (!isMobileOpen) return null;

  return (
    <>
      {/* Backdrop */}
      <div
        className="fixed inset-0 z-40 bg-black/50 backdrop-blur-sm md:hidden"
        onClick={() => setMobileOpen(false)}
      />

      {/* Sidebar */}
      <aside className="fixed left-0 top-0 z-50 h-screen w-[280px] bg-sidebar text-sidebar-foreground md:hidden animate-in slide-in-from-left duration-300">
        <div className="flex h-full flex-col">
          {/* Header */}
          <div className="flex h-16 items-center justify-between border-b border-sidebar-border px-4">
            <Link href="/" className="flex items-center gap-2">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary">
                <Leaf className="h-5 w-5 text-primary-foreground" />
              </div>
              <span className="text-lg font-semibold">GreenRide</span>
            </Link>
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8 text-sidebar-foreground hover:bg-sidebar-accent"
              onClick={() => setMobileOpen(false)}
            >
              <X className="h-4 w-4" />
            </Button>
          </div>

          {/* Navigation */}
          <nav className="flex-1 space-y-1 overflow-y-auto px-3 py-4">
            {navItems.map((item) => {
              const isActive = pathname === item.href || 
                (item.href !== '/' && pathname.startsWith(item.href));
              
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={cn(
                    'flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors',
                    isActive
                      ? 'bg-sidebar-accent text-primary border-l-3 border-primary'
                      : 'text-sidebar-foreground/70 hover:bg-sidebar-accent hover:text-sidebar-foreground'
                  )}
                >
                  <item.icon className="h-5 w-5 shrink-0" />
                  <span className="flex-1">{item.label}</span>
                  {item.badge && (
                    <span className="rounded-full bg-primary/20 px-2 py-0.5 text-xs text-primary">
                      {item.badge}
                    </span>
                  )}
                  {item.hasNotification && (
                    <span className="h-2 w-2 rounded-full bg-destructive" />
                  )}
                </Link>
              );
            })}
          </nav>

          {/* User Section */}
          <div className="border-t border-sidebar-border p-3">
            <div className="flex items-center gap-3">
              <Avatar className="h-10 w-10">
                <AvatarImage src={undefined} />
                <AvatarFallback className="bg-primary text-primary-foreground">
                  {user?.name?.charAt(0) || 'A'}
                </AvatarFallback>
              </Avatar>
              <div className="flex-1 min-w-0">
                <p className="truncate text-sm font-medium">{user?.name || 'Admin User'}</p>
                <p className="truncate text-xs text-sidebar-foreground/60">
                  {user?.role?.replace('_', ' ') || 'Super Admin'}
                </p>
              </div>
              <Button
                variant="ghost"
                size="icon"
                className="h-8 w-8 text-sidebar-foreground/70 hover:text-destructive"
                onClick={logout}
              >
                <LogOut className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>
      </aside>
    </>
  );
}
