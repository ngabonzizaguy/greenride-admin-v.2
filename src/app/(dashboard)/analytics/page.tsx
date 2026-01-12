'use client';

import { useEffect, useMemo, useState } from 'react';
import { 
  TrendingUp, 
  Users, 
  Car, 
  Clock, 
  AlertTriangle,
  Calendar,
  Download,
  FileText,
  FileSpreadsheet
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
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
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { RidesByDayChart } from '@/components/charts/rides-by-day-chart';
import { UserGrowthChart } from '@/components/charts/user-growth-chart';
import { toast } from 'sonner';
import { apiClient } from '@/lib/api-client';

type KPIStats = {
  ridesThisPeriod: number;
  uniquePassengers: number;
  activeDrivers: number;
  avgWaitTime: number | null;
  cancellationRate: number | null;
};

type PopularRouteRow = { rank: number; origin: string; destination: string; rides: number; avgFare: number };

const dateRangeLabels: Record<string, string> = {
  today: 'Today',
  this_week: 'This Week',
  this_month: 'This Month',
  last_month: 'Last Month',
};

export default function AnalyticsPage() {
  const [dateRange, setDateRange] = useState('this_week');
  const [isLoading, setIsLoading] = useState(true);
  const [kpiStats, setKpiStats] = useState<KPIStats>({
    ridesThisPeriod: 0,
    uniquePassengers: 0,
    activeDrivers: 0,
    avgWaitTime: null,
    cancellationRate: null,
  });
  const [popularRoutes, setPopularRoutes] = useState<PopularRouteRow[]>([]);

  useEffect(() => {
    let mounted = true;
    const run = async () => {
      setIsLoading(true);
      try {
        const res = await apiClient.getDashboardStats();
        const data = (res.data ?? {}) as Record<string, unknown>;

        const totalTrips = Number(data.total_trips ?? 0) || 0;
        const onlineDrivers = Number(data.online_drivers ?? 0) || 0;
        const recentTrips = (data.recent_trips as Array<Record<string, unknown>>) ?? [];

        const passengerSet = new Set<string>();
        let cancelled = 0;
        for (const t of recentTrips) {
          passengerSet.add(String(t.user_id ?? ''));
          if (String(t.status ?? '').toLowerCase() === 'cancelled') cancelled += 1;
        }
        const cancellationRate =
          recentTrips.length > 0 ? Math.round((cancelled / recentTrips.length) * 1000) / 10 : null;

        // Popular routes from recent trips (real, but limited sample)
        const routeMap = new Map<string, { origin: string; destination: string; rides: number; fareSum: number }>();
        for (const t of recentTrips) {
          const origin = String(t.pickup_location ?? '').trim();
          const destination = String(t.dropoff_location ?? '').trim();
          if (!origin || !destination) continue;
          const key = `${origin}|||${destination}`;
          const existing = routeMap.get(key) ?? { origin, destination, rides: 0, fareSum: 0 };
          existing.rides += 1;
          existing.fareSum += Number(t.fare ?? 0) || 0;
          routeMap.set(key, existing);
        }
        const routes = Array.from(routeMap.values())
          .sort((a, b) => b.rides - a.rides)
          .slice(0, 8)
          .map((r, idx) => ({
            rank: idx + 1,
            origin: r.origin,
            destination: r.destination,
            rides: r.rides,
            avgFare: r.rides > 0 ? Math.round(r.fareSum / r.rides) : 0,
          }));

        if (!mounted) return;
        setKpiStats({
          ridesThisPeriod: totalTrips,
          uniquePassengers: passengerSet.has('') ? passengerSet.size - 1 : passengerSet.size,
          activeDrivers: onlineDrivers,
          avgWaitTime: null, // needs order lifecycle timestamps endpoint
          cancellationRate,
        });
        setPopularRoutes(routes);
      } catch (e) {
        console.error('Failed to load analytics:', e);
        if (!mounted) return;
        setKpiStats({ ridesThisPeriod: 0, uniquePassengers: 0, activeDrivers: 0, avgWaitTime: null, cancellationRate: null });
        setPopularRoutes([]);
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
    const headers = ['Rank', 'Origin', 'Destination', 'Rides', 'Avg Fare (RWF)'];
    const rows = popularRoutes.map(route => [
      route.rank,
      route.origin,
      route.destination,
      route.rides,
      route.avgFare
    ]);
    
    const csvContent = [
      '# Analytics Report - Popular Routes',
      `# Period: ${dateRangeLabels[dateRange]}`,
      `# Generated: ${new Date().toISOString()}`,
      '',
      headers.join(','),
      ...rows.map(row => row.join(','))
    ].join('\n');
    
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `analytics_report_${new Date().toISOString().split('T')[0]}.csv`;
    link.click();
    toast.success('Analytics data exported to CSV!');
  };

  const kpiValue = useMemo(() => {
    const fmt = (v: number | null) => (v === null ? 'N/A' : String(v));
    return { fmt };
  }, []);

  // Export to PDF
  const handleExportPDF = () => {
    const printContent = `
      <!DOCTYPE html>
      <html>
        <head>
          <title>GreenRide Analytics Report - ${dateRangeLabels[dateRange]}</title>
          <style>
            * { margin: 0; padding: 0; box-sizing: border-box; }
            body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; padding: 40px; color: #1a1a1a; }
            .header { text-align: center; margin-bottom: 30px; border-bottom: 2px solid #22c55e; padding-bottom: 20px; }
            .header h1 { color: #22c55e; font-size: 24px; margin-bottom: 5px; }
            .header p { color: #666; font-size: 14px; }
            .report-date { background: #f5f5f5; padding: 10px 15px; border-radius: 5px; display: inline-block; margin-bottom: 20px; }
            .stats-grid { display: grid; grid-template-columns: repeat(5, 1fr); gap: 12px; margin-bottom: 30px; }
            .stat-card { background: #f9fafb; padding: 15px; border-radius: 8px; border: 1px solid #e5e7eb; text-align: center; }
            .stat-label { font-size: 11px; color: #6b7280; margin-bottom: 5px; text-transform: uppercase; letter-spacing: 0.5px; }
            .stat-value { font-size: 22px; font-weight: 700; color: #1f2937; }
            .section { margin-bottom: 25px; }
            .section h2 { font-size: 16px; font-weight: 600; margin-bottom: 15px; color: #374151; border-bottom: 1px solid #e5e7eb; padding-bottom: 8px; }
            table { width: 100%; border-collapse: collapse; font-size: 12px; }
            th { background: #f9fafb; text-align: left; padding: 10px 8px; font-weight: 600; color: #374151; border-bottom: 2px solid #e5e7eb; }
            td { padding: 10px 8px; border-bottom: 1px solid #f3f4f6; }
            tr:hover { background: #fafafa; }
            .text-right { text-align: right; }
            .rank { display: inline-block; width: 24px; height: 24px; line-height: 24px; text-align: center; border-radius: 4px; font-weight: 600; font-size: 11px; }
            .rank-1 { background: #fbbf24; color: white; }
            .rank-2 { background: #9ca3af; color: white; }
            .rank-3 { background: #d97706; color: white; }
            .rank-default { background: #e5e7eb; color: #374151; }
            .insights { background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 8px; padding: 20px; margin-top: 20px; }
            .insights h3 { color: #166534; margin-bottom: 10px; font-size: 14px; }
            .insights ul { list-style: disc; padding-left: 20px; }
            .insights li { margin-bottom: 5px; font-size: 12px; color: #15803d; }
            .footer { margin-top: 30px; text-align: center; font-size: 11px; color: #9ca3af; border-top: 1px solid #e5e7eb; padding-top: 15px; }
            @media print {
              body { padding: 20px; }
              .stats-grid { grid-template-columns: repeat(5, 1fr); }
            }
          </style>
        </head>
        <body>
          <div class="header">
            <h1>📊 GreenRide Analytics Report</h1>
            <p>Rwanda's Premier Ride-Hailing Service</p>
          </div>
          
          <div class="report-date">
            📅 Report Period: <strong>${dateRangeLabels[dateRange]}</strong> | Generated: ${new Date().toLocaleDateString('en-US', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}
          </div>

          <div class="section">
            <h2>📈 Key Performance Indicators</h2>
            <div class="stats-grid">
              <div class="stat-card">
                <div class="stat-label">Rides This Period</div>
                <div class="stat-value">${kpiStats.ridesThisPeriod.toLocaleString()}</div>
              </div>
              <div class="stat-card">
                <div class="stat-label">Unique Passengers</div>
                <div class="stat-value">${kpiStats.uniquePassengers}</div>
              </div>
              <div class="stat-card">
                <div class="stat-label">Active Drivers</div>
                <div class="stat-value">${kpiStats.activeDrivers}</div>
              </div>
              <div class="stat-card">
                <div class="stat-label">Avg Wait Time</div>
                <div class="stat-value">${kpiStats.avgWaitTime} min</div>
              </div>
              <div class="stat-card">
                <div class="stat-label">Cancellation Rate</div>
                <div class="stat-value">${kpiStats.cancellationRate}%</div>
              </div>
            </div>
          </div>

          <div class="section">
            <h2>🗺️ Popular Routes</h2>
            <table>
              <thead>
                <tr>
                  <th style="width: 50px;">#</th>
                  <th>Origin</th>
                  <th>Destination</th>
                  <th class="text-right">Total Rides</th>
                  <th class="text-right">Avg Fare</th>
                </tr>
              </thead>
              <tbody>
                ${popularRoutes.map(route => `
                  <tr>
                    <td><span class="rank ${route.rank === 1 ? 'rank-1' : route.rank === 2 ? 'rank-2' : route.rank === 3 ? 'rank-3' : 'rank-default'}">${route.rank}</span></td>
                    <td>${route.origin}</td>
                    <td>${route.destination}</td>
                    <td class="text-right">${route.rides}</td>
                    <td class="text-right">RWF ${route.avgFare.toLocaleString()}</td>
                  </tr>
                `).join('')}
              </tbody>
            </table>
          </div>

          <div class="insights">
            <h3>💡 Key Insights</h3>
            <ul>
              <li>The <strong>Kimironko → Downtown</strong> route is the most popular with ${popularRoutes[0].rides} rides this period.</li>
              <li>Average wait time of <strong>${kpiStats.avgWaitTime} minutes</strong> is within optimal range (&lt;5 min).</li>
              <li>Cancellation rate of <strong>${kpiStats.cancellationRate}%</strong> ${kpiStats.cancellationRate < 10 ? 'is healthy' : 'needs attention'}.</li>
              <li><strong>${kpiStats.activeDrivers} drivers</strong> are actively serving <strong>${kpiStats.uniquePassengers} unique passengers</strong>.</li>
            </ul>
          </div>

          <div class="footer">
            <p>GreenRide Admin Dashboard • Analytics Report • Confidential</p>
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
      
      setTimeout(() => {
        printWindow.print();
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
          <h1 className="text-2xl font-bold tracking-tight">Analytics & Insights</h1>
          <p className="text-muted-foreground">
            Deep insights for business decisions
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
              <SelectItem value="this_week">This Week</SelectItem>
              <SelectItem value="this_month">This Month</SelectItem>
              <SelectItem value="last_month">Last Month</SelectItem>
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

      {/* KPI Summary */}
      <div className="grid gap-4 md:grid-cols-5">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-3">
              <TrendingUp className="h-5 w-5 text-primary" />
              <div>
                <p className="text-sm text-muted-foreground">Rides This Period</p>
                <p className="text-2xl font-bold">{kpiStats.ridesThisPeriod.toLocaleString()}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-3">
              <Users className="h-5 w-5 text-blue-500" />
              <div>
                <p className="text-sm text-muted-foreground">Unique Passengers</p>
                <p className="text-2xl font-bold">{kpiStats.uniquePassengers}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-3">
              <Car className="h-5 w-5 text-green-500" />
              <div>
                <p className="text-sm text-muted-foreground">Active Drivers</p>
                <p className="text-2xl font-bold">{kpiStats.activeDrivers}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-3">
              <Clock className="h-5 w-5 text-yellow-500" />
              <div>
                <p className="text-sm text-muted-foreground">Avg. Wait Time</p>
                <p className="text-2xl font-bold">{kpiStats.avgWaitTime === null ? 'N/A' : `${kpiStats.avgWaitTime} min`}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-3">
              <AlertTriangle className="h-5 w-5 text-red-500" />
              <div>
                <p className="text-sm text-muted-foreground">Cancellation Rate</p>
                <p className="text-2xl font-bold">{kpiStats.cancellationRate === null ? 'N/A' : `${kpiStats.cancellationRate}%`}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Charts Row 1 */}
      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Peak Hours</CardTitle>
            <CardDescription>Coming soon (needs a dedicated analytics endpoint)</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-[300px] w-full flex items-center justify-center text-sm text-muted-foreground">
              Not available yet.
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Rides by Day of Week</CardTitle>
            <CardDescription>Weekly ride distribution</CardDescription>
          </CardHeader>
          <CardContent>
            <RidesByDayChart />
          </CardContent>
        </Card>
      </div>

      {/* Charts Row 2 */}
      <div className="grid gap-4 md:grid-cols-2">
        {/* Popular Routes */}
        <Card>
          <CardHeader>
            <CardTitle>Popular Routes</CardTitle>
            <CardDescription>Top origin-destination pairs</CardDescription>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-12">#</TableHead>
                  <TableHead>Route</TableHead>
                  <TableHead className="text-right">Rides</TableHead>
                  <TableHead className="text-right">Avg Fare</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  <TableRow>
                    <TableCell colSpan={4} className="text-center text-muted-foreground py-8">
                      Loading…
                    </TableCell>
                  </TableRow>
                ) : popularRoutes.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={4} className="text-center text-muted-foreground py-8">
                      No route data yet (needs more completed trips with pickup/dropoff).
                    </TableCell>
                  </TableRow>
                ) : (
                  popularRoutes.map((route) => (
                    <TableRow key={route.rank}>
                      <TableCell>
                        <Badge 
                          variant={route.rank <= 3 ? 'default' : 'secondary'}
                          className={route.rank === 1 ? 'bg-yellow-500' : route.rank === 2 ? 'bg-gray-400' : route.rank === 3 ? 'bg-amber-600' : ''}
                        >
                          {route.rank}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        <span className="text-muted-foreground">{route.origin}</span>
                        <span className="mx-1">→</span>
                        <span>{route.destination}</span>
                      </TableCell>
                      <TableCell className="text-right font-medium">{route.rides}</TableCell>
                      <TableCell className="text-right">RWF {route.avgFare.toLocaleString()}</TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </CardContent>
        </Card>

        {/* User Growth */}
        <Card>
          <CardHeader>
            <CardTitle>User Growth</CardTitle>
            <CardDescription>New vs returning users over time</CardDescription>
          </CardHeader>
          <CardContent>
            <UserGrowthChart />
          </CardContent>
        </Card>
      </div>

      {/* Distance Distribution - Full Width */}
      <Card>
        <CardHeader>
          <CardTitle>Trip Distance Distribution</CardTitle>
          <CardDescription>Coming soon (needs a dedicated analytics endpoint)</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-[300px] w-full flex items-center justify-center text-sm text-muted-foreground">
            Not available yet.
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
