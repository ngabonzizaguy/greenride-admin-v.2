'use client';

import { useEffect, useState, useRef } from 'react';
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
  Clock,
  FileText,
  FileSpreadsheet,
  Printer
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
import { toast } from 'sonner';
import { apiClient } from '@/lib/api-client';

type RevenueStats = {
  totalRevenue: number;
  revenueChange: number; // keep 0 unless you add a historical comparison endpoint
  completedRides: number;
  averageFare: number;
  platformCommission: number;
  cashPayments: { amount: number; percentage: number };
  momoPayments: { amount: number; percentage: number };
  cardPayments: { amount: number; percentage: number };
  pendingPayments: { amount: number; count: number };
};

type TransactionRow = {
  id: string;
  rideId: string;
  date: string;
  passenger: string;
  driver: string;
  amount: number;
  method: string;
  status: string;
};

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

const dateRangeLabels: Record<string, string> = {
  today: 'Today',
  yesterday: 'Yesterday',
  this_week: 'This Week',
  this_month: 'This Month',
  last_month: 'Last Month',
  custom: 'Custom Range',
};

export default function RevenuePage() {
  const [dateRange, setDateRange] = useState('today');
  const [isLoading, setIsLoading] = useState(true);
  const [revenueStats, setRevenueStats] = useState<RevenueStats>({
    totalRevenue: 0,
    revenueChange: 0,
    completedRides: 0,
    averageFare: 0,
    platformCommission: 0,
    cashPayments: { amount: 0, percentage: 0 },
    momoPayments: { amount: 0, percentage: 0 },
    cardPayments: { amount: 0, percentage: 0 },
    pendingPayments: { amount: 0, count: 0 },
  });
  const [transactions, setTransactions] = useState<TransactionRow[]>([]);

  useEffect(() => {
    let mounted = true;
    const run = async () => {
      setIsLoading(true);
      try {
        const res = await apiClient.getDashboardStats();
        const data = (res.data ?? {}) as Record<string, unknown>;

        const totalRevenue = Number(data.total_revenue ?? 0) || 0;
        const totalTrips = Number(data.total_trips ?? 0) || 0;
        const averageFare = totalTrips > 0 ? Math.round(totalRevenue / totalTrips) : 0;
        const platformCommission = Math.round(totalRevenue * 0.15);

        const recentTrips = (data.recent_trips as Array<Record<string, unknown>>) ?? [];
        const txns: TransactionRow[] = recentTrips.slice(0, 20).map((t) => ({
          id: String(t.id ?? ''),
          rideId: String(t.id ?? ''),
          date: new Date(String(t.created_at ?? t.requested_at ?? new Date().toISOString())).toLocaleString(),
          passenger: String(t.user_id ?? ''),
          driver: String(t.driver_id ?? ''),
          amount: Number(t.fare ?? 0) || 0,
          method: String(t.payment_method ?? 'unknown'),
          status: String(t.payment_status ?? t.status ?? 'unknown'),
        }));

        const sums: Record<string, number> = { cash: 0, momo: 0, card: 0 };
        for (const txn of txns) {
          const key = txn.method.toLowerCase();
          if (key in sums) sums[key] += txn.amount;
        }
        const totalByMethod = Object.values(sums).reduce((a, b) => a + b, 0);
        const pct = (n: number) => (totalByMethod > 0 ? Math.round((n / totalByMethod) * 100) : 0);

        if (!mounted) return;
        setTransactions(txns);
        setRevenueStats({
          totalRevenue,
          revenueChange: 0, // avoid fake numbers
          completedRides: totalTrips,
          averageFare,
          platformCommission,
          cashPayments: { amount: sums.cash, percentage: pct(sums.cash) },
          momoPayments: { amount: sums.momo, percentage: pct(sums.momo) },
          cardPayments: { amount: sums.card, percentage: pct(sums.card) },
          pendingPayments: { amount: 0, count: 0 }, // needs payments endpoint
        });
      } catch (e) {
        console.error('Failed to load revenue data:', e);
        if (!mounted) return;
        setTransactions([]);
        setRevenueStats({
          totalRevenue: 0,
          revenueChange: 0,
          completedRides: 0,
          averageFare: 0,
          platformCommission: 0,
          cashPayments: { amount: 0, percentage: 0 },
          momoPayments: { amount: 0, percentage: 0 },
          cardPayments: { amount: 0, percentage: 0 },
          pendingPayments: { amount: 0, count: 0 },
        });
      } finally {
        if (mounted) setIsLoading(false);
      }
    };
    run();
    return () => {
      mounted = false;
    };
  }, [dateRange]);

  // Export to CSV
  const handleExportCSV = () => {
    const headers = ['Transaction ID', 'Ride ID', 'Date/Time', 'Passenger', 'Driver', 'Amount (RWF)', 'Method', 'Status'];
    const rows = transactions.map(txn => [
      txn.id,
      txn.rideId,
      txn.date,
      txn.passenger,
      txn.driver,
      txn.amount,
      txn.method,
      txn.status
    ]);
    
    const csvContent = [
      headers.join(','),
      ...rows.map(row => row.join(','))
    ].join('\n');
    
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `revenue_transactions_${new Date().toISOString().split('T')[0]}.csv`;
    link.click();
    toast.success('Transactions exported to CSV!');
  };

  // Export to PDF (using print dialog)
  const handleExportPDF = () => {
    const printContent = `
      <!DOCTYPE html>
      <html>
        <head>
          <title>GreenRide Revenue Report - ${dateRangeLabels[dateRange]}</title>
          <style>
            * { margin: 0; padding: 0; box-sizing: border-box; }
            body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; padding: 40px; color: #1a1a1a; }
            .header { text-align: center; margin-bottom: 30px; border-bottom: 2px solid #22c55e; padding-bottom: 20px; }
            .header h1 { color: #22c55e; font-size: 24px; margin-bottom: 5px; }
            .header p { color: #666; font-size: 14px; }
            .report-date { background: #f5f5f5; padding: 10px 15px; border-radius: 5px; display: inline-block; margin-bottom: 20px; }
            .stats-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 15px; margin-bottom: 30px; }
            .stat-card { background: #f9fafb; padding: 15px; border-radius: 8px; border: 1px solid #e5e7eb; }
            .stat-label { font-size: 12px; color: #6b7280; margin-bottom: 5px; }
            .stat-value { font-size: 20px; font-weight: 700; color: #1f2937; }
            .stat-change { font-size: 12px; color: #22c55e; }
            .section { margin-bottom: 25px; }
            .section h2 { font-size: 16px; font-weight: 600; margin-bottom: 15px; color: #374151; border-bottom: 1px solid #e5e7eb; padding-bottom: 8px; }
            table { width: 100%; border-collapse: collapse; font-size: 12px; }
            th { background: #f9fafb; text-align: left; padding: 10px 8px; font-weight: 600; color: #374151; border-bottom: 2px solid #e5e7eb; }
            td { padding: 10px 8px; border-bottom: 1px solid #f3f4f6; }
            tr:hover { background: #fafafa; }
            .amount { text-align: right; font-weight: 600; }
            .status-completed { color: #22c55e; }
            .status-pending { color: #eab308; }
            .status-failed { color: #ef4444; }
            .status-refunded { color: #6b7280; }
            .payment-summary { display: grid; grid-template-columns: repeat(3, 1fr); gap: 15px; margin-bottom: 20px; }
            .payment-card { padding: 15px; border-radius: 8px; text-align: center; }
            .payment-card.cash { background: #f5f5f5; }
            .payment-card.momo { background: #fef9c3; }
            .payment-card.card { background: #dbeafe; }
            .footer { margin-top: 30px; text-align: center; font-size: 11px; color: #9ca3af; border-top: 1px solid #e5e7eb; padding-top: 15px; }
            @media print {
              body { padding: 20px; }
              .stats-grid { grid-template-columns: repeat(2, 1fr); }
              .payment-summary { grid-template-columns: repeat(3, 1fr); }
            }
          </style>
        </head>
        <body>
          <div class="header">
            <h1>🚗 GreenRide Revenue Report</h1>
            <p>Rwanda's Premier Ride-Hailing Service</p>
          </div>
          
          <div class="report-date">
            📅 Report Period: <strong>${dateRangeLabels[dateRange]}</strong> | Generated: ${new Date().toLocaleDateString('en-US', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}
          </div>

          <div class="section">
            <h2>📊 Financial Summary</h2>
            <div class="stats-grid">
              <div class="stat-card">
                <div class="stat-label">Total Revenue</div>
                <div class="stat-value">RWF ${revenueStats.totalRevenue.toLocaleString()}</div>
                <div class="stat-change">+${revenueStats.revenueChange}% vs last period</div>
              </div>
              <div class="stat-card">
                <div class="stat-label">Completed Rides</div>
                <div class="stat-value">${revenueStats.completedRides}</div>
                <div class="stat-change">rides this period</div>
              </div>
              <div class="stat-card">
                <div class="stat-label">Average Fare</div>
                <div class="stat-value">RWF ${revenueStats.averageFare.toLocaleString()}</div>
                <div class="stat-change">per ride</div>
              </div>
              <div class="stat-card">
                <div class="stat-label">Platform Commission</div>
                <div class="stat-value">RWF ${revenueStats.platformCommission.toLocaleString()}</div>
                <div class="stat-change">15% of revenue</div>
              </div>
            </div>
          </div>

          <div class="section">
            <h2>💳 Payment Methods Breakdown</h2>
            <div class="payment-summary">
              <div class="payment-card cash">
                <div class="stat-label">Cash Payments</div>
                <div class="stat-value">RWF ${revenueStats.cashPayments.amount.toLocaleString()}</div>
                <div class="stat-change">${revenueStats.cashPayments.percentage}% of total</div>
              </div>
              <div class="payment-card momo">
                <div class="stat-label">Mobile Money</div>
                <div class="stat-value">RWF ${revenueStats.momoPayments.amount.toLocaleString()}</div>
                <div class="stat-change">${revenueStats.momoPayments.percentage}% of total</div>
              </div>
              <div class="payment-card card">
                <div class="stat-label">Card Payments</div>
                <div class="stat-value">RWF ${revenueStats.cardPayments.amount.toLocaleString()}</div>
                <div class="stat-change">${revenueStats.cardPayments.percentage}% of total</div>
              </div>
            </div>
          </div>

          <div class="section">
            <h2>📋 Recent Transactions</h2>
            <table>
              <thead>
                <tr>
                  <th>Transaction ID</th>
                  <th>Date/Time</th>
                  <th>Ride ID</th>
                  <th>Passenger</th>
                  <th>Driver</th>
                  <th class="amount">Amount</th>
                  <th>Method</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                ${transactions.map(txn => `
                  <tr>
                    <td>${txn.id}</td>
                    <td>${txn.date}</td>
                    <td>${txn.rideId}</td>
                    <td>${txn.passenger}</td>
                    <td>${txn.driver}</td>
                    <td class="amount">RWF ${txn.amount.toLocaleString()}</td>
                    <td>${txn.method.toUpperCase()}</td>
                    <td class="status-${txn.status}">${txn.status.charAt(0).toUpperCase() + txn.status.slice(1)}</td>
                  </tr>
                `).join('')}
              </tbody>
            </table>
          </div>

          <div class="footer">
            <p>GreenRide Admin Dashboard • Revenue Report • Confidential</p>
            <p>Generated on ${new Date().toISOString()}</p>
          </div>
        </body>
      </html>
    `;

    const printWindow = window.open('', '_blank');
    if (printWindow) {
      printWindow.document.write(printContent);
      printWindow.document.close();
      printWindow.focus();
      
      // Wait for content to load then trigger print
      setTimeout(() => {
        printWindow.print();
        // Close window after print dialog closes (user can cancel to keep it open)
      }, 250);
      
      toast.success('PDF report opened in new window. Use print dialog to save as PDF.');
    } else {
      toast.error('Unable to open print window. Please allow popups.');
    }
  };

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
              <DropdownMenuItem onClick={handleExportPDF} className="gap-2">
                <FileText className="h-4 w-4" />
                Export as PDF
              </DropdownMenuItem>
              <DropdownMenuItem onClick={handleExportCSV} className="gap-2">
                <FileSpreadsheet className="h-4 w-4" />
                Export as CSV
              </DropdownMenuItem>
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
          <Button variant="outline" size="sm" asChild>
            <a href="/rides">View All Transactions</a>
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
              {isLoading ? (
                <TableRow>
                  <TableCell colSpan={8} className="text-center text-muted-foreground py-8">
                    Loading…
                  </TableCell>
                </TableRow>
              ) : transactions.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={8} className="text-center text-muted-foreground py-8">
                    No transactions yet.
                  </TableCell>
                </TableRow>
              ) : (
                transactions.map((txn) => (
                  <TableRow key={txn.id || txn.rideId}>
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
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
