'use client';

import { useEffect, useState } from 'react';
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

// Mock recent activity (will be replaced when backend provides this endpoint)
const mockRecentActivity = [
  { id: 1, time: '2 min ago', event: 'Ride completed', user: 'John Doe', driver: 'Peter M.', status: 'completed', amount: 5200 },
  { id: 2, time: '5 min ago', event: 'Payment received', user: 'Jane Smith', driver: 'David K.', status: 'paid', amount: 3800 },
  { id: 3, time: '8 min ago', event: 'Ride started', user: 'Mike Johnson', driver: 'Paul R.', status: 'in_progress', amount: 4500 },
  { id: 4, time: '12 min ago', event: 'Ride cancelled', user: 'Sarah Wilson', driver: 'James T.', status: 'cancelled', amount: 0 },
  { id: 5, time: '15 min ago', event: 'Ride completed', user: 'Chris Brown', driver: 'Alex M.', status: 'completed', amount: 6100 },
];

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
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchDashboardData = async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await apiClient.getDashboardStats();
      setStats(response.data as DashboardStats);
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
              value={stats.active_rides}
              icon={Car}
              changeLabel="in progress"
            />
            <StatsCard
              title="Online Drivers"
              value={stats.online_drivers}
              icon={Users}
              subtext={`${stats.total_drivers} total drivers`}
            />
            <StatsCard
              title="Today's Revenue"
              value={`RWF ${stats.today_revenue.toLocaleString()}`}
              icon={DollarSign}
              changeLabel="today"
            />
            <StatsCard
              title="Rides Today"
              value={stats.today_rides}
              icon={Activity}
              subtext={`${stats.total_users} total users`}
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
              value={`RWF ${stats.pending_payments.toLocaleString()}`}
              icon={Clock}
              subtext="awaiting confirmation"
              variant="warning"
            />
            <StatsCard
              title="Cancellation Rate"
              value={`${stats.cancellation_rate.toFixed(1)}%`}
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
              <a href="/rides">View All</a>
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
                {mockRecentActivity.map((activity) => (
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
              Real-time activity feed coming soon
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
