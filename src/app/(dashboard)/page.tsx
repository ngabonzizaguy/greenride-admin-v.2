'use client';

import { 
  Car, 
  Users, 
  DollarSign, 
  TrendingUp,
  Activity,
  Clock,
  ArrowUpRight,
  ArrowDownRight,
  MoreHorizontal,
  Eye,
  MessageSquare,
  Download,
  AlertCircle
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
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { StatsCard } from '@/components/charts/stats-card';
import { RevenueChart } from '@/components/charts/revenue-chart';
import { PaymentMethodsChart } from '@/components/charts/payment-methods-chart';

// Mock data for the dashboard
const stats = {
  activeRides: { value: 12, change: 3, changeType: 'increase' as const },
  onlineDrivers: { value: 8, subtext: '2 on trip' },
  todayRevenue: { value: 245000, change: 15, changeType: 'increase' as const },
  todayRides: { value: 67, subtext: '5 cancelled' },
  pendingPayments: { value: 24500, count: 8 },
  cancellationRate: { value: 7.5, change: 2.1, changeType: 'decrease' as const },
};

const recentActivity = [
  { id: 1, time: '2 min ago', event: 'Ride completed', user: 'John Doe', driver: 'Peter M.', status: 'completed', amount: 5200 },
  { id: 2, time: '5 min ago', event: 'Payment received', user: 'Jane Smith', driver: 'David K.', status: 'paid', amount: 3800 },
  { id: 3, time: '8 min ago', event: 'Ride started', user: 'Mike Johnson', driver: 'Paul R.', status: 'in_progress', amount: 4500 },
  { id: 4, time: '12 min ago', event: 'Ride cancelled', user: 'Sarah Wilson', driver: 'James T.', status: 'cancelled', amount: 0 },
  { id: 5, time: '15 min ago', event: 'Ride completed', user: 'Chris Brown', driver: 'Alex M.', status: 'completed', amount: 6100 },
  { id: 6, time: '20 min ago', event: 'Driver went online', user: '-', driver: 'Peter M.', status: 'online', amount: 0 },
  { id: 7, time: '25 min ago', event: 'New user registered', user: 'Emma Davis', driver: '-', status: 'new_user', amount: 0 },
  { id: 8, time: '30 min ago', event: 'Ride completed', user: 'Tom Harris', driver: 'David K.', status: 'completed', amount: 4200 },
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

export default function DashboardPage() {
  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground">
          Welcome back! Here&apos;s what&apos;s happening with GreenRide today.
        </p>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <StatsCard
          title="Active Rides"
          value={stats.activeRides.value}
          icon={Car}
          change={stats.activeRides.change}
          changeType={stats.activeRides.changeType}
          changeLabel="vs yesterday"
        />
        <StatsCard
          title="Online Drivers"
          value={stats.onlineDrivers.value}
          icon={Users}
          subtext={stats.onlineDrivers.subtext}
        />
        <StatsCard
          title="Today's Revenue"
          value={`RWF ${stats.todayRevenue.value.toLocaleString()}`}
          icon={DollarSign}
          change={stats.todayRevenue.change}
          changeType={stats.todayRevenue.changeType}
          changeLabel="vs yesterday"
        />
        <StatsCard
          title="Rides Today"
          value={stats.todayRides.value}
          icon={Activity}
          subtext={stats.todayRides.subtext}
        />
      </div>

      {/* Secondary Stats */}
      <div className="grid gap-4 md:grid-cols-2">
        <StatsCard
          title="Pending Payments"
          value={`RWF ${stats.pendingPayments.value.toLocaleString()}`}
          icon={Clock}
          subtext={`${stats.pendingPayments.count} pending`}
          variant="warning"
        />
        <StatsCard
          title="Cancellation Rate"
          value={`${stats.cancellationRate.value}%`}
          icon={TrendingUp}
          change={stats.cancellationRate.change}
          changeType={stats.cancellationRate.changeType}
          changeLabel="vs yesterday"
        />
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
            <Button className="w-full justify-start gap-2" variant="default">
              <Eye className="h-4 w-4" />
              View Active Rides
            </Button>
            <Button className="w-full justify-start gap-2" variant="outline">
              <MessageSquare className="h-4 w-4" />
              Broadcast Message
            </Button>
            <Button className="w-full justify-start gap-2" variant="outline">
              <Download className="h-4 w-4" />
              Export Report
            </Button>
            <Button className="w-full justify-start gap-2" variant="outline">
              <AlertCircle className="h-4 w-4" />
              View Issues
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
            <Button variant="outline" size="sm">
              View All
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
                {recentActivity.map((activity) => (
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
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
