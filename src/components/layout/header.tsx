'use client';

import { useState } from 'react';
import { 
  Bell, 
  Search, 
  Menu,
  X,
  Check
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { useSidebarStore } from '@/stores/sidebar-store';
import { useAuthStore } from '@/stores/auth-store';
import { cn } from '@/lib/utils';

// Mock notifications for demo
const notifications = [
  {
    id: '1',
    type: 'ride',
    title: 'New ride completed',
    message: 'Trip #1234 completed successfully',
    time: '2 min ago',
    read: false,
  },
  {
    id: '2',
    type: 'payment',
    title: 'Payment received',
    message: 'RWF 5,200 via MoMo',
    time: '5 min ago',
    read: false,
  },
  {
    id: '3',
    type: 'driver',
    title: 'Driver went online',
    message: 'John Doe is now available',
    time: '10 min ago',
    read: true,
  },
];

export function Header() {
  const { isCollapsed, toggleCollapsed, isMobileOpen, toggleMobile } = useSidebarStore();
  const { user, logout } = useAuthStore();
  const [searchOpen, setSearchOpen] = useState(false);

  const unreadCount = notifications.filter((n) => !n.read).length;

  return (
    <header className={cn(
      'sticky top-0 z-30 flex h-16 items-center gap-4 border-b bg-background px-4 transition-all duration-300',
      isCollapsed ? 'md:pl-[88px]' : 'md:pl-[276px]'
    )}>
      {/* Mobile Menu Button */}
      <Button
        variant="ghost"
        size="icon"
        className="md:hidden"
        onClick={toggleMobile}
      >
        {isMobileOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
      </Button>

      {/* Desktop Collapse Button */}
      <Button
        variant="ghost"
        size="icon"
        className="hidden md:flex"
        onClick={toggleCollapsed}
      >
        <Menu className="h-5 w-5" />
      </Button>

      {/* Search */}
      <div className="flex flex-1 items-center gap-2">
        {searchOpen ? (
          <div className="flex flex-1 items-center gap-2 md:max-w-md">
            <Input
              placeholder="Search drivers, users, rides..."
              className="h-9"
              autoFocus
              onBlur={() => setSearchOpen(false)}
            />
            <Button
              variant="ghost"
              size="icon"
              className="h-9 w-9 shrink-0"
              onClick={() => setSearchOpen(false)}
            >
              <X className="h-4 w-4" />
            </Button>
          </div>
        ) : (
          <Button
            variant="ghost"
            size="sm"
            className="hidden gap-2 text-muted-foreground md:flex"
            onClick={() => setSearchOpen(true)}
          >
            <Search className="h-4 w-4" />
            <span>Search...</span>
            <kbd className="pointer-events-none ml-2 hidden h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100 sm:flex">
              <span className="text-xs">⌘</span>K
            </kbd>
          </Button>
        )}
        <Button
          variant="ghost"
          size="icon"
          className="h-9 w-9 md:hidden"
          onClick={() => setSearchOpen(true)}
        >
          <Search className="h-4 w-4" />
        </Button>
      </div>

      {/* Right Side */}
      <div className="flex items-center gap-2">
        {/* Notifications */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon" className="relative h-9 w-9">
              <Bell className="h-5 w-5" />
              {unreadCount > 0 && (
                <span className="absolute -right-0.5 -top-0.5 flex h-4 w-4 items-center justify-center rounded-full bg-destructive text-[10px] font-medium text-destructive-foreground">
                  {unreadCount}
                </span>
              )}
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-80">
            <DropdownMenuLabel className="flex items-center justify-between">
              <span>Notifications</span>
              <Button variant="ghost" size="sm" className="h-auto p-0 text-xs text-primary hover:bg-transparent">
                Mark all as read
              </Button>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            {notifications.map((notification) => (
              <DropdownMenuItem key={notification.id} className="flex flex-col items-start gap-1 p-3">
                <div className="flex w-full items-start justify-between gap-2">
                  <p className={cn(
                    'text-sm',
                    !notification.read && 'font-medium'
                  )}>
                    {notification.title}
                  </p>
                  {!notification.read && (
                    <span className="h-2 w-2 shrink-0 rounded-full bg-primary" />
                  )}
                </div>
                <p className="text-xs text-muted-foreground">{notification.message}</p>
                <p className="text-xs text-muted-foreground/60">{notification.time}</p>
              </DropdownMenuItem>
            ))}
            <DropdownMenuSeparator />
            <DropdownMenuItem className="justify-center text-primary">
              View all notifications
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>

        {/* User Menu */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon" className="h-9 w-9">
              <Avatar className="h-8 w-8">
                <AvatarImage src={undefined} />
                <AvatarFallback className="bg-primary text-primary-foreground text-xs">
                  {user?.name?.charAt(0) || 'A'}
                </AvatarFallback>
              </Avatar>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-56">
            <DropdownMenuLabel>
              <div className="flex flex-col space-y-1">
                <p className="text-sm font-medium">{user?.name || 'Admin User'}</p>
                <p className="text-xs text-muted-foreground">{user?.email || 'admin@greenrideafrica.com'}</p>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem>Profile</DropdownMenuItem>
            <DropdownMenuItem>Settings</DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem 
              className="text-destructive focus:text-destructive"
              onClick={logout}
            >
              Log out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  );
}
