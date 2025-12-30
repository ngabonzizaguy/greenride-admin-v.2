'use client';

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
  Leaf
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
];

const bottomNavItems = [
  { href: '/settings', icon: Settings, label: 'Settings' },
];

export function Sidebar() {
  const pathname = usePathname();
  const { isCollapsed, toggleCollapsed } = useSidebarStore();
  const { user, logout } = useAuthStore();

  return (
    <TooltipProvider delayDuration={0}>
      <aside
        className={cn(
          'fixed left-0 top-0 z-40 h-screen bg-sidebar text-sidebar-foreground transition-all duration-300 ease-in-out',
          isCollapsed ? 'w-[72px]' : 'w-[260px]'
        )}
      >
        <div className="flex h-full flex-col">
          {/* Header */}
          <div className={cn(
            'flex h-16 items-center border-b border-sidebar-border px-4',
            isCollapsed ? 'justify-center' : 'justify-between'
          )}>
            <Link href="/" className="flex items-center gap-2">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary">
                <Leaf className="h-5 w-5 text-primary-foreground" />
              </div>
              {!isCollapsed && (
                <span className="text-lg font-semibold">GreenRide</span>
              )}
            </Link>
            {!isCollapsed && (
              <Button
                variant="ghost"
                size="icon"
                className="h-8 w-8 text-sidebar-foreground hover:bg-sidebar-accent"
                onClick={toggleCollapsed}
              >
                <ChevronLeft className="h-4 w-4" />
              </Button>
            )}
          </div>

          {/* Navigation */}
          <nav className="flex-1 space-y-1 px-3 py-4">
            {navItems.map((item) => {
              const isActive = pathname === item.href || 
                (item.href !== '/' && pathname.startsWith(item.href));
              
              const linkContent = (
                <Link
                  href={item.href}
                  className={cn(
                    'flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors',
                    isActive
                      ? 'bg-sidebar-accent text-primary border-l-3 border-primary'
                      : 'text-sidebar-foreground/70 hover:bg-sidebar-accent hover:text-sidebar-foreground',
                    isCollapsed && 'justify-center px-2'
                  )}
                >
                  <item.icon className="h-5 w-5 shrink-0" />
                  {!isCollapsed && (
                    <>
                      <span className="flex-1">{item.label}</span>
                      {item.badge && (
                        <span className="rounded-full bg-primary/20 px-2 py-0.5 text-xs text-primary">
                          {item.badge}
                        </span>
                      )}
                      {item.hasNotification && (
                        <span className="h-2 w-2 rounded-full bg-destructive" />
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
                      {item.badge && (
                        <span className="rounded-full bg-primary/20 px-2 py-0.5 text-xs text-primary">
                          {item.badge}
                        </span>
                      )}
                    </TooltipContent>
                  </Tooltip>
                );
              }

              return <div key={item.href}>{linkContent}</div>;
            })}
          </nav>

          {/* Divider */}
          <div className="mx-3 border-t border-sidebar-border" />

          {/* Bottom Navigation */}
          <nav className="space-y-1 px-3 py-4">
            {bottomNavItems.map((item) => {
              const isActive = pathname === item.href;
              
              const linkContent = (
                <Link
                  href={item.href}
                  className={cn(
                    'flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors',
                    isActive
                      ? 'bg-sidebar-accent text-primary'
                      : 'text-sidebar-foreground/70 hover:bg-sidebar-accent hover:text-sidebar-foreground',
                    isCollapsed && 'justify-center px-2'
                  )}
                >
                  <item.icon className="h-5 w-5 shrink-0" />
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
            'border-t border-sidebar-border p-3',
            isCollapsed && 'flex justify-center'
          )}>
            {isCollapsed ? (
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-10 w-10"
                    onClick={toggleCollapsed}
                  >
                    <Avatar className="h-8 w-8">
                      <AvatarImage src={user?.avatar || undefined} />
                      <AvatarFallback className="bg-primary text-primary-foreground text-xs">
                        {user?.full_name?.charAt(0) || user?.username?.charAt(0) || 'A'}
                      </AvatarFallback>
                    </Avatar>
                  </Button>
                </TooltipTrigger>
                <TooltipContent side="right">
                  <p className="font-medium">{user?.full_name || user?.username || 'Admin User'}</p>
                  <p className="text-xs text-muted-foreground">{user?.role || 'Super Admin'}</p>
                </TooltipContent>
              </Tooltip>
            ) : (
              <div className="flex items-center gap-3">
                <Avatar className="h-10 w-10">
                  <AvatarImage src={undefined} />
                  <AvatarFallback className="bg-primary text-primary-foreground">
                    {user?.full_name?.charAt(0) || user?.username?.charAt(0) || 'A'}
                  </AvatarFallback>
                </Avatar>
                <div className="flex-1 min-w-0">
                  <p className="truncate text-sm font-medium">{user?.full_name || user?.username || 'Admin User'}</p>
                  <p className="truncate text-xs text-sidebar-foreground/60">
                    {user?.role?.replace('_', ' ') || 'Super Admin'}
                  </p>
                </div>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-8 w-8 text-sidebar-foreground/70 hover:text-destructive"
                      onClick={logout}
                    >
                      <LogOut className="h-4 w-4" />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent side="right">Logout</TooltipContent>
                </Tooltip>
              </div>
            )}
          </div>
        </div>
      </aside>
    </TooltipProvider>
  );
}
