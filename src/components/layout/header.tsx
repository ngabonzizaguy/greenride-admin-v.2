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
      'h-16 glass-panel border-b-0 sticky top-0 z-30 flex items-center justify-between px-6 mx-6 mt-4 rounded-xl transition-all duration-300',
      isCollapsed ? 'md:ml-[88px]' : 'md:ml-[276px]' // Adjusted margin-left instead of padding-left because it's floating
    )}>
      {/* Mobile Menu Button */}
      <Button
        variant="ghost"
        size="icon"
        className="md:hidden mr-4 text-gray-500 hover:text-gray-700"
        onClick={toggleMobile}
      >
        {isMobileOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
      </Button>

      {/* Desktop Collapse Button - Removed as per design snippet, but kept if needed for functionality? 
          Snippet doesn't have it explicitly besides mobile menu. 
          I'll keep it but style it minimally or hide if desired. 
          Actually, snippet has a menu button only for md:hidden.
          I'll keep the collapse button logic but maybe style it simpler.
      */}
      <Button
        variant="ghost"
        size="icon"
        className="hidden md:flex mr-4 text-gray-400 hover:text-gray-600"
        onClick={toggleCollapsed}
      >
        <Menu className="h-5 w-5" />
      </Button>

      {/* Search */}
      <div className="relative hidden md:block w-full max-w-md">
        <span className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <Search className="h-4 w-4 text-gray-400" />
        </span>
        <Input
            className="block w-full pl-10 pr-3 py-2 border border-gray-200/50 dark:border-gray-600/30 rounded-lg leading-5 bg-white/50 dark:bg-gray-800/50 text-gray-900 dark:text-gray-100 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary/50 focus:bg-white/80 sm:text-sm transition duration-200 ease-in-out backdrop-blur-sm shadow-none"
            placeholder="Search anything..."
        />
        <div className="absolute inset-y-0 right-0 pr-2 flex items-center">
            <kbd className="inline-flex items-center border border-gray-200/50 dark:border-gray-600/50 rounded px-2 text-xs font-sans font-medium text-gray-400 bg-white/30">⌘K</kbd>
        </div>
      </div>

      {/* Right Side */}
      <div className="flex items-center space-x-4 ml-4">
        {/* Notifications */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button className="p-2 rounded-full text-gray-500 hover:bg-white/50 dark:hover:bg-gray-700/50 transition-colors relative hover:shadow-lg outline-none">
              <Bell className="h-5 w-5" />
              {unreadCount > 0 && (
                <span className="absolute top-2 right-2 h-2 w-2 rounded-full bg-red-500 ring-2 ring-white/80 dark:ring-gray-800/80 shadow-[0_0_8px_rgba(239,68,68,0.6)]"></span>
              )}
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-80 glass-panel border-0">
             {/* ... content ... */}
            <DropdownMenuLabel className="flex items-center justify-between">
              <span>Notifications</span>
              <Button variant="ghost" size="sm" className="h-auto p-0 text-xs text-primary hover:bg-transparent">
                Mark all as read
              </Button>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            {notifications.map((notification) => (
              <DropdownMenuItem key={notification.id} className="flex flex-col items-start gap-1 p-3 focus:bg-white/40">
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
            <DropdownMenuItem className="justify-center text-primary focus:bg-white/40" asChild>
              <a href="/notifications">View all notifications</a>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>

        {/* User Menu - Keeping simple logout button or avatar if needed, but sidebar has profile. 
           Snippet shows "dark mode" button here. 
           I'll keep the avatar menu for consistency but style it minimally.
        */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
             <button className="p-2 rounded-full text-gray-500 hover:bg-white/50 dark:hover:bg-gray-700/50 transition-colors hover:shadow-lg outline-none">
                <Avatar className="h-8 w-8">
                    <AvatarImage src={undefined} />
                    <AvatarFallback className="bg-primary text-primary-foreground text-xs">
                    {user?.full_name?.charAt(0) || user?.username?.charAt(0) || 'A'}
                    </AvatarFallback>
                </Avatar>
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-56 glass-panel border-0">
            <DropdownMenuLabel>
              <div className="flex flex-col space-y-1">
                <p className="text-sm font-medium">{user?.full_name || user?.username || 'Admin User'}</p>
                <p className="text-xs text-muted-foreground">{user?.email || 'admin@greenrideafrica.com'}</p>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem className="focus:bg-white/40">Profile</DropdownMenuItem>
            <DropdownMenuItem className="focus:bg-white/40">Settings</DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem 
              className="text-destructive focus:text-destructive focus:bg-white/40"
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
