'use client';

import { useState } from 'react';
import { 
  DollarSign, 
  TrendingUp, 
  ArrowUpRight, 
  ArrowDownRight,
  Download,
  Calendar,
  CreditCard,
  Smartphone,
  Banknote,
  Clock
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { RevenueChart } from '@/components/charts/revenue-chart';
import { PaymentMethodsChart } from '@/components/charts/payment-methods-chart';

// Mock revenue data
const revenueStats = {
  totalRevenue: 1245000,
  revenueChange: 12,
  completedRides: 342,
  averageFare: 3640,
  platformCommission: 186750,
  cashPayments: { amount: 560250, percentage: 45 },
  momoPayments: { amount: 498000, percentage: 40 },
  cardPayments: { amount: 186750, percentage: 15 },
  pendingPayments: { amount: 24500, count: 8 },
};

const transactions = [
  { id: 'TXN-2024122812345', rideId: 'R001', date: '2024-12-28 14:30', passenger: 'John Doe', driver: 'Peter M.', amount: 4500, method: 'momo', status: 'completed' },
  { id: 'TXN-2024122812344', rideId: 'R002', date: '2024-12-28 14:15', passenger: 'Jane Smith', driver: 'David K.', amount: 3200, method: 'cash', status: 'completed' },
  { id: 'TXN-2024122812343', rideId: 'R003', date: '2024-12-28 14:00', passenger: 'Mike Johnson', driver: 'Claude U.', amount: 6800, method: 'momo', status: 'pending' },
  { id: 'TXN-2024122812342', rideId: 'R004', date: '2024-12-28 13:55', passenger: 'Sarah Wilson', driver: 'Emmanuel H.', amount: 2500, method: 'cash', status: 'completed' },
  { id: 'TXN-2024122812341', rideId: 'R005', date: '2024-12-28 13:50', passenger: 'Chris Brown', driver: 'Jean P.', amount: 3600, method: 'card', status: 'completed' },
  { id: 'TXN-2024122812340', rideId: 'R006', date: '2024-12-28 13:30', passenger: 'Emma Davis', driver: 'Patrick N.', amount: 5400, method: 'momo', status: 'failed' },
  { id: 'TXN-2024122812339', rideId: 'R007', date: '2024-12-28 13:15', passenger: 'Tom Harris', driver: 'David K.', amount: 4200, method: 'cash', status: 'completed' },
  { id: 'TXN-2024122812338', rideId: 'R008', date: '2024-12-28 13:00', passenger: 'Lisa Brown', driver: 'Peter M.', amount: 3800, method: 'momo', status: 'refunded' },
];

const getMethodIcon = (method: string) => {
  switch (method) {
    case 'cash':
      return <Banknote className="h-4 w-4 text-gray-500" />;
    case 'momo':
      return <Smartphone className="h-4 w-4 text-yellow-500" />;
    case 'card':
      return <CreditCard className="h-4 w-4 text-blue-500" />;
    default:
      return null;
  }
};

const getStatusBadge = (status: string) => {
  switch (status) {
    case 'completed':
      return <Badge className="bg-green-100 text-green-700 hover:bg-green-100">Completed</Badge>;
    case 'pending':
      return <Badge className="bg-yellow-100 text-yellow-700 hover:bg-yellow-100">Pending</Badge>;
    case 'failed':
      return <Badge className="bg-red-100 text-red-700 hover:bg-red-100">Failed</Badge>;
    case 'refunded':
      return <Badge className="bg-gray-100 text-gray-700 hover:bg-gray-100">Refunded</Badge>;
    default:
      return <Badge variant="secondary">{status}</Badge>;
  }
};

export default function RevenuePage() {
  const [dateRange, setDateRange] = useState('today');

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Revenue & Financials</h1>
          <p className="text-muted-foreground">
            Track revenue, payments, and financial metrics
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Select value={dateRange} onValueChange={setDateRange}>
            <SelectTrigger className="w-[180px]">
              <Calendar className="mr-2 h-4 w-4" />
              <SelectValue placeholder="Select period" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="today">Today</SelectItem>
              <SelectItem value="yesterday">Yesterday</SelectItem>
              <SelectItem value="this_week">This Week</SelectItem>
              <SelectItem value="this_month">This Month</SelectItem>
              <SelectItem value="last_month">Last Month</SelectItem>
              <SelectItem value="custom">Custom Range</SelectItem>
            </SelectContent>
          </Select>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" className="gap-2">
                <Download className="h-4 w-4" />
                Export
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem>Export as PDF</DropdownMenuItem>
              <DropdownMenuItem>Export as Excel</DropdownMenuItem>
              <DropdownMenuItem>Export as CSV</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      {/* Primary Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Revenue</p>
                <p className="text-2xl font-bold">RWF {revenueStats.totalRevenue.toLocaleString()}</p>
                <div className="flex items-center gap-1 text-sm">
                  <ArrowUpRight className="h-4 w-4 text-green-500" />
                  <span className="font-medium text-green-500">+{revenueStats.revenueChange}%</span>
                  <span className="text-muted-foreground">vs last period</span>
                </div>
              </div>
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-green-100">
                <DollarSign className="h-6 w-6 text-green-600" />
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Completed Rides</p>
                <p className="text-2xl font-bold">{revenueStats.completedRides}</p>
                <p className="text-sm text-muted-foreground">rides this period</p>
              </div>
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-blue-100">
                <TrendingUp className="h-6 w-6 text-blue-600" />
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Average Fare</p>
                <p className="text-2xl font-bold">RWF {revenueStats.averageFare.toLocaleString()}</p>
                <p className="text-sm text-muted-foreground">per ride</p>
              </div>
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-purple-100">
                <DollarSign className="h-6 w-6 text-purple-600" />
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Platform Commission</p>
                <p className="text-2xl font-bold">RWF {revenueStats.platformCommission.toLocaleString()}</p>
                <p className="text-sm text-muted-foreground">15% of revenue</p>
              </div>
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
                <DollarSign className="h-6 w-6 text-primary" />
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Payment Method Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Cash Payments</p>
                <p className="text-2xl font-bold">RWF {revenueStats.cashPayments.amount.toLocaleString()}</p>
              </div>
              <Badge variant="secondary">{revenueStats.cashPayments.percentage}%</Badge>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">MoMo Payments</p>
                <p className="text-2xl font-bold">RWF {revenueStats.momoPayments.amount.toLocaleString()}</p>
              </div>
              <Badge className="bg-yellow-100 text-yellow-700">{revenueStats.momoPayments.percentage}%</Badge>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Card Payments</p>
                <p className="text-2xl font-bold">RWF {revenueStats.cardPayments.amount.toLocaleString()}</p>
              </div>
              <Badge className="bg-blue-100 text-blue-700">{revenueStats.cardPayments.percentage}%</Badge>
            </div>
          </CardContent>
        </Card>
        <Card className="border-yellow-200 bg-yellow-50">
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-yellow-800">Pending Payments</p>
                <p className="text-2xl font-bold text-yellow-900">RWF {revenueStats.pendingPayments.amount.toLocaleString()}</p>
              </div>
              <div className="flex items-center gap-2">
                <Clock className="h-5 w-5 text-yellow-600" />
                <Badge className="bg-yellow-200 text-yellow-800">{revenueStats.pendingPayments.count} pending</Badge>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Charts Row */}
      <div className="grid gap-4 md:grid-cols-7">
        <Card className="md:col-span-4">
          <CardHeader>
            <CardTitle>Revenue Trend</CardTitle>
            <CardDescription>Daily revenue over the selected period</CardDescription>
          </CardHeader>
          <CardContent>
            <RevenueChart />
          </CardContent>
        </Card>
        <Card className="md:col-span-3">
          <CardHeader>
            <CardTitle>Payment Methods</CardTitle>
            <CardDescription>Distribution by payment type</CardDescription>
          </CardHeader>
          <CardContent>
            <PaymentMethodsChart />
          </CardContent>
        </Card>
      </div>

      {/* Transactions Table */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle>Recent Transactions</CardTitle>
            <CardDescription>Latest payment transactions</CardDescription>
          </div>
          <Button variant="outline" size="sm">
            View All
          </Button>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Transaction ID</TableHead>
                <TableHead>Date/Time</TableHead>
                <TableHead>Ride ID</TableHead>
                <TableHead>Passenger</TableHead>
                <TableHead>Driver</TableHead>
                <TableHead className="text-right">Amount</TableHead>
                <TableHead>Method</TableHead>
                <TableHead>Status</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {transactions.map((txn) => (
                <TableRow key={txn.id}>
                  <TableCell className="font-mono text-sm">{txn.id}</TableCell>
                  <TableCell className="text-muted-foreground">{txn.date}</TableCell>
                  <TableCell className="font-medium">{txn.rideId}</TableCell>
                  <TableCell>{txn.passenger}</TableCell>
                  <TableCell>{txn.driver}</TableCell>
                  <TableCell className="text-right font-medium">
                    RWF {txn.amount.toLocaleString()}
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center gap-2">
                      {getMethodIcon(txn.method)}
                      <span className="capitalize">{txn.method}</span>
                    </div>
                  </TableCell>
                  <TableCell>{getStatusBadge(txn.status)}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
