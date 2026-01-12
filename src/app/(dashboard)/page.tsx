'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { 
  Car, 
  Users, 
  DollarSign, 
  TrendingUp,
  Activity,
  Clock,
  Eye,
  MessageSquare,
  Download,
  AlertCircle,
  RefreshCw
} from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Skeleton } from '@/components/ui/skeleton';
import { StatsCard } from '@/components/charts/stats-card';
import { RevenueChart } from '@/components/charts/revenue-chart';
import { PaymentMethodsChart } from '@/components/charts/payment-methods-chart';
import { apiClient } from '@/lib/api-client';
import type { DashboardStats } from '@/types';

// Recent activity type
interface RecentActivity {
  id: string | number;
  time: string;
  event: string;
  user: string;
  driver: string;
  status: string;
  amount: number;
}

// Helper to format time ago
const formatTimeAgo = (timestamp: number): string => {
  const seconds = Math.floor((Date.now() - timestamp) / 1000);
  if (seconds < 60) return `${seconds} sec ago`;
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes} min ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours} hr ago`;
  const days = Math.floor(hours / 24);
  return `${days} day${days > 1 ? 's' : ''} ago`;
};

const getStatusBadge = (status: string) => {
  switch (status) {
    case 'completed':
      return <Badge className="bg-green-100 text-green-700 hover:bg-green-100">Completed</Badge>;
    case 'in_progress':
      return <Badge className="bg-yellow-100 text-yellow-700 hover:bg-yellow-100">In Progress</Badge>;
    case 'cancelled':
      return <Badge className="bg-red-100 text-red-700 hover:bg-red-100">Cancelled</Badge>;
    case 'paid':
      return <Badge className="bg-blue-100 text-blue-700 hover:bg-blue-100">Paid</Badge>;
    case 'online':
      return <Badge className="bg-green-100 text-green-700 hover:bg-green-100">Online</Badge>;
    case 'new_user':
      return <Badge className="bg-purple-100 text-purple-700 hover:bg-purple-100">New User</Badge>;
    default:
      return <Badge variant="secondary">{status}</Badge>;
  }
};

// Default stats for when API fails
const defaultStats: DashboardStats = {
  active_rides: 0,
  online_drivers: 0,
  today_revenue: 0,
  today_rides: 0,
  pending_payments: 0,
  cancellation_rate: 0,
  total_users: 0,
  total_drivers: 0,
};

export default function DashboardPage() {
  const [stats, setStats] = useState<DashboardStats>(defaultStats);
  const [recentActivity, setRecentActivity] = useState<RecentActivity[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchDashboardData = async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      // Fetch dashboard stats
      const response = await apiClient.getDashboardStats();
      // Map backend response to frontend format with fallbacks
      const backendData = response.data as Record<string, unknown>;
      setStats({
        active_rides: (backendData.active_trips as number) ?? 0,
        online_drivers: (backendData.online_drivers as number) ?? 0,
        today_revenue: (backendData.total_revenue as number) ?? (backendData.today_revenue as number) ?? 0,
        today_rides: (backendData.total_trips as number) ?? (backendData.today_rides as number) ?? 0,
        pending_payments: (backendData.pending_payments as number) ?? 0,
        cancellation_rate: (backendData.cancellation_rate as number) ?? 0,
        total_users: (backendData.total_users as number) ?? 0,
        total_drivers: (backendData.total_drivers as number) ?? 0,
      });

      // Recent Activity: Prefer backend-provided recent trips (real data), fallback to orders search.
      const toMillis = (value: unknown): number => {
        if (typeof value === 'number') {
          // Heuristic: seconds vs milliseconds
          return value < 10_000_000_000 ? value * 1000 : value;
        }
        if (typeof value === 'string') {
          const parsed = Date.parse(value);
          return Number.isFinite(parsed) ? parsed : Date.now();
        }
        return Date.now();
      };

      const recentTrips = (backendData.recent_trips as Array<Record<string, unknown>>) ?? [];
      if (Array.isArray(recentTrips) && recentTrips.length > 0) {
        const activities: RecentActivity[] = recentTrips.slice(0, 5).map((trip) => {
          const status = (trip.status as string) || 'pending';
          const createdAt = toMillis(trip.created_at ?? trip.requested_at ?? trip.updated_at);
          return {
            id: (trip.id as number) ?? Math.random(),
            time: formatTimeAgo(createdAt),
            event:
              status === 'completed' ? 'Ride completed' :
              status === 'cancelled' ? 'Ride cancelled' :
              status === 'in_progress' ? 'Ride in progress' :
              status === 'pending' ? 'New ride request' : 'Ride update',
            user: String(trip.user_id ?? 'Customer'),
            driver: String(trip.driver_id ?? 'Driver'),
            status,
            amount: (trip.fare as number) ?? 0,
          };
        });
        setRecentActivity(activities);
      } else {
        try {
          const ordersResponse = await apiClient.searchOrders({ page: 1, limit: 5 });
          if (ordersResponse.code === '0000' && ordersResponse.data?.records) {
            const orders = ordersResponse.data.records as Array<Record<string, unknown>>;
            const activities: RecentActivity[] = orders.map((order, idx) => {
              const status = (order.status as string) || 'pending';
              const createdAt = toMillis(order.created_at ?? order.updated_at);
              return {
                id: (order.order_id as string) ?? (order.id as string) ?? idx,
                time: formatTimeAgo(createdAt),
                event:
                  status === 'completed' ? 'Ride completed' :
                  status === 'cancelled' ? 'Ride cancelled' :
                  status === 'in_progress' ? 'Ride in progress' :
                  status === 'pending' ? 'New ride request' : 'Ride update',
                user: String(order.user_id ?? order.customer_id ?? 'Customer'),
                driver: String(order.provider_id ?? order.driver_id ?? 'Driver'),
                status,
                amount: (order.payment_amount as number) ?? (order.fare as number) ?? 0,
              };
            });
            setRecentActivity(activities);
          }
        } catch (orderErr) {
          console.error('Failed to fetch recent activity:', orderErr);
        }
      }
    } catch (err) {
      console.error('Failed to fetch dashboard stats:', err);
      setError('Failed to load dashboard data. Using cached data.');
      // Keep existing stats or use defaults
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchDashboardData();
    
    // Refresh stats every 30 seconds
    const interval = setInterval(fetchDashboardData, 30000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Dashboard</h1>
          <p className="text-muted-foreground">
            Welcome back! Here&apos;s what&apos;s happening with GreenRide today.
          </p>
        </div>
        <Button 
          variant="outline" 
          size="sm" 
          onClick={fetchDashboardData}
          disabled={isLoading}
        >
          <RefreshCw className={`h-4 w-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      {/* Error Banner */}
      {error && (
        <div className="flex items-center gap-2 rounded-lg bg-yellow-50 border border-yellow-200 p-3 text-sm text-yellow-800">
          <AlertCircle className="h-4 w-4 flex-shrink-0" />
          <span>{error}</span>
        </div>
      )}

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {isLoading ? (
          <>
            {[1, 2, 3, 4].map((i) => (
              <Card key={i}>
                <CardContent className="p-6">
                  <Skeleton className="h-4 w-24 mb-2" />
                  <Skeleton className="h-8 w-16" />
                </CardContent>
              </Card>
            ))}
          </>
        ) : (
          <>
            <StatsCard
              title="Active Rides"
              value={stats.active_rides ?? 0}
              icon={Car}
              changeLabel="in progress"
            />
            <StatsCard
              title="Online Drivers"
              value={stats.online_drivers ?? 0}
              icon={Users}
              subtext={`${stats.total_drivers ?? 0} total drivers`}
            />
            <StatsCard
              title="Today's Revenue"
              value={`RWF ${(stats.today_revenue ?? 0).toLocaleString()}`}
              icon={DollarSign}
              changeLabel="today"
            />
            <StatsCard
              title="Rides Today"
              value={stats.today_rides ?? 0}
              icon={Activity}
              subtext={`${stats.total_users ?? 0} total users`}
            />
          </>
        )}
      </div>

      {/* Secondary Stats */}
      <div className="grid gap-4 md:grid-cols-2">
        {isLoading ? (
          <>
            {[1, 2].map((i) => (
              <Card key={i}>
                <CardContent className="p-6">
                  <Skeleton className="h-4 w-24 mb-2" />
                  <Skeleton className="h-8 w-16" />
                </CardContent>
              </Card>
            ))}
          </>
        ) : (
          <>
            <StatsCard
              title="Pending Payments"
              value={`RWF ${(stats.pending_payments ?? 0).toLocaleString()}`}
              icon={Clock}
              subtext="awaiting confirmation"
              variant="warning"
            />
            <StatsCard
              title="Cancellation Rate"
              value={`${(stats.cancellation_rate ?? 0).toFixed(1)}%`}
              icon={TrendingUp}
              changeLabel="overall"
            />
          </>
        )}
      </div>

      {/* Charts Row */}
      <div className="grid gap-4 md:grid-cols-7">
        <Card className="md:col-span-4">
          <CardHeader>
            <CardTitle>Revenue This Week</CardTitle>
            <CardDescription>Daily revenue for the current week</CardDescription>
          </CardHeader>
          <CardContent>
            <RevenueChart />
          </CardContent>
        </Card>
        <Card className="md:col-span-3">
          <CardHeader>
            <CardTitle>Payment Methods</CardTitle>
            <CardDescription>Distribution of payment methods</CardDescription>
          </CardHeader>
          <CardContent>
            <PaymentMethodsChart />
          </CardContent>
        </Card>
      </div>

      {/* Quick Actions & Recent Activity */}
      <div className="grid gap-4 lg:grid-cols-4">
        {/* Quick Actions */}
        <Card className="lg:col-span-1">
          <CardHeader>
            <CardTitle>Quick Actions</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            <Button className="w-full justify-start gap-2" variant="default" asChild>
              <a href="/rides">
                <Eye className="h-4 w-4" />
                View Active Rides
              </a>
            </Button>
            <Button className="w-full justify-start gap-2" variant="outline" asChild>
              <a href="/notifications">
                <MessageSquare className="h-4 w-4" />
                Broadcast Message
              </a>
            </Button>
            <Button className="w-full justify-start gap-2" variant="outline" asChild>
              <a href="/revenue">
                <Download className="h-4 w-4" />
                Export Report
              </a>
            </Button>
            <Button className="w-full justify-start gap-2" variant="outline" asChild>
              <a href="/drivers">
                <AlertCircle className="h-4 w-4" />
                Manage Drivers
              </a>
            </Button>
          </CardContent>
        </Card>

        {/* Recent Activity */}
        <Card className="lg:col-span-3">
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle>Recent Activity</CardTitle>
              <CardDescription>Latest events from your platform</CardDescription>
            </div>
            <Button variant="outline" size="sm" asChild>
              <Link href="/rides">View All</Link>
            </Button>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Time</TableHead>
                  <TableHead>Event</TableHead>
                  <TableHead>User</TableHead>
                  <TableHead>Driver</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Amount</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {recentActivity.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center text-muted-foreground py-8">
                      {isLoading ? 'Loading recent activity...' : 'No recent activity'}
                    </TableCell>
                  </TableRow>
                ) : recentActivity.map((activity) => (
                  <TableRow key={activity.id}>
                    <TableCell className="text-muted-foreground text-sm">
                      {activity.time}
                    </TableCell>
                    <TableCell className="font-medium">{activity.event}</TableCell>
                    <TableCell>{activity.user}</TableCell>
                    <TableCell>{activity.driver}</TableCell>
                    <TableCell>{getStatusBadge(activity.status)}</TableCell>
                    <TableCell className="text-right">
                      {activity.amount > 0 ? `RWF ${activity.amount.toLocaleString()}` : '-'}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
            <p className="text-xs text-muted-foreground mt-4 text-center">
              Showing last 5 orders • Auto-refreshes every 30 seconds
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
