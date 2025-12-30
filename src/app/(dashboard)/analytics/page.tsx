'use client';

import { 
  TrendingUp, 
  Users, 
  Car, 
  Clock, 
  AlertTriangle,
  Calendar
} from 'lucide-react';
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
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { PeakHoursChart } from '@/components/charts/peak-hours-chart';
import { RidesByDayChart } from '@/components/charts/rides-by-day-chart';
import { UserGrowthChart } from '@/components/charts/user-growth-chart';
import { DistanceDistributionChart } from '@/components/charts/distance-distribution-chart';

// Mock data
const kpiStats = {
  ridesThisPeriod: 1234,
  uniquePassengers: 456,
  activeDrivers: 24,
  avgWaitTime: 4.2,
  cancellationRate: 8.5,
};

const popularRoutes = [
  { rank: 1, origin: 'Kimironko', destination: 'Downtown', rides: 234, avgFare: 4500 },
  { rank: 2, origin: 'Remera', destination: 'Nyarutarama', rides: 189, avgFare: 3800 },
  { rank: 3, origin: 'Kicukiro', destination: 'Gisozi', rides: 156, avgFare: 5200 },
  { rank: 4, origin: 'Nyamirambo', destination: 'Kigali Heights', rides: 134, avgFare: 5800 },
  { rank: 5, origin: 'Kacyiru', destination: 'Kibagabaga', rides: 121, avgFare: 3200 },
  { rank: 6, origin: 'Gikondo', destination: 'Remera', rides: 98, avgFare: 4100 },
  { rank: 7, origin: 'Kanombe', destination: 'CBD', rides: 87, avgFare: 6500 },
  { rank: 8, origin: 'Kimihurura', destination: 'Kigali Arena', rides: 76, avgFare: 2800 },
];

export default function AnalyticsPage() {
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
        <Select defaultValue="this_week">
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
                <p className="text-2xl font-bold">{kpiStats.avgWaitTime} min</p>
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
                <p className="text-2xl font-bold">{kpiStats.cancellationRate}%</p>
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
            <CardDescription>Ride demand by hour and day of week</CardDescription>
          </CardHeader>
          <CardContent>
            <PeakHoursChart />
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
                {popularRoutes.map((route) => (
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
                ))}
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
          <CardDescription>How trip distances are distributed across all rides</CardDescription>
        </CardHeader>
        <CardContent>
          <DistanceDistributionChart />
        </CardContent>
      </Card>
    </div>
  );
}
