'use client';

import { useState } from 'react';
import Link from 'next/link';
import { 
  Search, 
  MoreHorizontal,
  Eye,
  XCircle,
  ChevronLeft,
  ChevronRight,
  Car,
  Filter,
  MapPin,
  Clock
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent } from '@/components/ui/card';
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
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

// Mock ride data
const mockRides = [
  {
    id: 'R001',
    passenger: 'John Doe',
    driver: 'Peter Mugisha',
    pickup: 'Kimironko Market',
    dropoff: 'Downtown Kigali',
    distance: 5.2,
    duration: 18,
    fare: 4500,
    paymentMethod: 'momo',
    paymentStatus: 'paid',
    status: 'completed',
    createdAt: '2024-12-28T14:30:00',
  },
  {
    id: 'R002',
    passenger: 'Jane Smith',
    driver: 'David Kamanzi',
    pickup: 'Remera',
    dropoff: 'Nyarutarama',
    distance: 3.8,
    duration: 12,
    fare: 3200,
    paymentMethod: 'cash',
    paymentStatus: 'paid',
    status: 'completed',
    createdAt: '2024-12-28T14:15:00',
  },
  {
    id: 'R003',
    passenger: 'Mike Johnson',
    driver: 'Claude Uwimana',
    pickup: 'Kicukiro',
    dropoff: 'Gisozi',
    distance: 8.1,
    duration: null,
    fare: 6800,
    paymentMethod: 'momo',
    paymentStatus: 'pending',
    status: 'in_progress',
    createdAt: '2024-12-28T14:00:00',
  },
  {
    id: 'R004',
    passenger: 'Sarah Wilson',
    driver: 'Emmanuel Habimana',
    pickup: 'Downtown',
    dropoff: 'Kimihurura',
    distance: 2.5,
    duration: null,
    fare: 2500,
    paymentMethod: 'cash',
    paymentStatus: 'pending',
    status: 'driver_arriving',
    createdAt: '2024-12-28T13:55:00',
  },
  {
    id: 'R005',
    passenger: 'Chris Brown',
    driver: null,
    pickup: 'Kacyiru',
    dropoff: 'Kibagabaga',
    distance: 4.0,
    duration: null,
    fare: 3600,
    paymentMethod: 'card',
    paymentStatus: 'pending',
    status: 'requesting',
    createdAt: '2024-12-28T13:50:00',
  },
  {
    id: 'R006',
    passenger: 'Emma Davis',
    driver: 'Jean Pierre',
    pickup: 'Nyamirambo',
    dropoff: 'Kigali Heights',
    distance: 6.3,
    duration: null,
    fare: 5400,
    paymentMethod: 'momo',
    paymentStatus: 'failed',
    status: 'cancelled',
    createdAt: '2024-12-28T13:30:00',
  },
];

const getStatusBadge = (status: string) => {
  switch (status) {
    case 'completed':
      return <Badge className="bg-green-100 text-green-700 hover:bg-green-100">Completed</Badge>;
    case 'in_progress':
      return <Badge className="bg-blue-100 text-blue-700 hover:bg-blue-100">In Progress</Badge>;
    case 'driver_arriving':
      return <Badge className="bg-yellow-100 text-yellow-700 hover:bg-yellow-100">Driver Arriving</Badge>;
    case 'requesting':
      return <Badge className="bg-purple-100 text-purple-700 hover:bg-purple-100">Requesting</Badge>;
    case 'cancelled':
      return <Badge className="bg-red-100 text-red-700 hover:bg-red-100">Cancelled</Badge>;
    default:
      return <Badge variant="secondary">{status}</Badge>;
  }
};

const getPaymentBadge = (method: string) => {
  switch (method) {
    case 'cash':
      return <Badge variant="outline" className="border-gray-300">Cash</Badge>;
    case 'momo':
      return <Badge variant="outline" className="border-yellow-400 text-yellow-600">MoMo</Badge>;
    case 'card':
      return <Badge variant="outline" className="border-blue-400 text-blue-600">Card</Badge>;
    default:
      return <Badge variant="outline">{method}</Badge>;
  }
};

const getPaymentStatusBadge = (status: string) => {
  switch (status) {
    case 'paid':
      return <Badge className="bg-green-100 text-green-700 hover:bg-green-100">Paid</Badge>;
    case 'pending':
      return <Badge className="bg-yellow-100 text-yellow-700 hover:bg-yellow-100">Pending</Badge>;
    case 'failed':
      return <Badge className="bg-red-100 text-red-700 hover:bg-red-100">Failed</Badge>;
    default:
      return <Badge variant="secondary">{status}</Badge>;
  }
};

export default function RidesPage() {
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');

  const filteredRides = mockRides.filter((ride) => {
    const matchesSearch =
      ride.id.toLowerCase().includes(search.toLowerCase()) ||
      ride.passenger.toLowerCase().includes(search.toLowerCase()) ||
      (ride.driver && ride.driver.toLowerCase().includes(search.toLowerCase()));
    const matchesStatus = statusFilter === 'all' || ride.status === statusFilter;
    return matchesSearch && matchesStatus;
  });

  const activeRides = mockRides.filter(r => ['in_progress', 'driver_arriving', 'requesting'].includes(r.status)).length;
  const completedToday = mockRides.filter(r => r.status === 'completed').length;
  const cancelledToday = mockRides.filter(r => r.status === 'cancelled').length;

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Ride Management</h1>
        <p className="text-muted-foreground">
          View and manage all rides on the platform
        </p>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Active Now</p>
                <p className="text-2xl font-bold text-blue-600">{activeRides}</p>
              </div>
              <div className="h-10 w-10 rounded-full bg-blue-100 flex items-center justify-center">
                <Car className="h-5 w-5 text-blue-600" />
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Completed Today</p>
                <p className="text-2xl font-bold text-green-600">{completedToday}</p>
              </div>
              <span className="h-3 w-3 rounded-full bg-green-500" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Cancelled Today</p>
                <p className="text-2xl font-bold text-red-600">{cancelledToday}</p>
              </div>
              <span className="h-3 w-3 rounded-full bg-red-500" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Avg Wait Time</p>
                <p className="text-2xl font-bold">4.2 min</p>
              </div>
              <Clock className="h-5 w-5 text-muted-foreground" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <Card>
        <CardContent className="p-4">
          <div className="flex flex-col gap-4 md:flex-row md:items-center">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Search by ride ID, passenger, or driver..."
                className="pl-10"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
            </div>
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="w-full md:w-[180px]">
                <SelectValue placeholder="Filter by status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Status</SelectItem>
                <SelectItem value="requesting">Requesting</SelectItem>
                <SelectItem value="driver_arriving">Driver Arriving</SelectItem>
                <SelectItem value="in_progress">In Progress</SelectItem>
                <SelectItem value="completed">Completed</SelectItem>
                <SelectItem value="cancelled">Cancelled</SelectItem>
              </SelectContent>
            </Select>
            <Button variant="outline" className="gap-2">
              <Filter className="h-4 w-4" />
              More Filters
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Rides Table */}
      <Card>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Ride ID</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Passenger</TableHead>
                <TableHead>Driver</TableHead>
                <TableHead>Route</TableHead>
                <TableHead className="text-right">Distance</TableHead>
                <TableHead className="text-right">Fare</TableHead>
                <TableHead>Payment</TableHead>
                <TableHead>Time</TableHead>
                <TableHead className="w-12"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredRides.map((ride) => (
                <TableRow key={ride.id}>
                  <TableCell className="font-mono font-medium">{ride.id}</TableCell>
                  <TableCell>{getStatusBadge(ride.status)}</TableCell>
                  <TableCell>{ride.passenger}</TableCell>
                  <TableCell>{ride.driver || <span className="text-muted-foreground">-</span>}</TableCell>
                  <TableCell>
                    <div className="flex items-center gap-1 text-sm">
                      <MapPin className="h-3 w-3 text-green-500" />
                      <span className="text-muted-foreground max-w-[100px] truncate">{ride.pickup}</span>
                      <span className="mx-1">→</span>
                      <MapPin className="h-3 w-3 text-red-500" />
                      <span className="max-w-[100px] truncate">{ride.dropoff}</span>
                    </div>
                  </TableCell>
                  <TableCell className="text-right">{ride.distance} km</TableCell>
                  <TableCell className="text-right font-medium">
                    RWF {ride.fare.toLocaleString()}
                  </TableCell>
                  <TableCell>
                    <div className="flex flex-col gap-1">
                      {getPaymentBadge(ride.paymentMethod)}
                      {getPaymentStatusBadge(ride.paymentStatus)}
                    </div>
                  </TableCell>
                  <TableCell className="text-muted-foreground text-sm">
                    {new Date(ride.createdAt).toLocaleTimeString('en-US', {
                      hour: '2-digit',
                      minute: '2-digit',
                    })}
                  </TableCell>
                  <TableCell>
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon" className="h-8 w-8">
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem asChild>
                          <Link href={`/rides/${ride.id}`} className="flex items-center gap-2">
                            <Eye className="h-4 w-4" />
                            View Details
                          </Link>
                        </DropdownMenuItem>
                        {ride.status !== 'completed' && ride.status !== 'cancelled' && (
                          <>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem className="gap-2 text-red-600">
                              <XCircle className="h-4 w-4" />
                              Cancel Ride
                            </DropdownMenuItem>
                          </>
                        )}
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>

          {/* Pagination */}
          <div className="flex items-center justify-between border-t px-4 py-4">
            <p className="text-sm text-muted-foreground">
              Showing 1-{filteredRides.length} of {mockRides.length} rides
            </p>
            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm" disabled>
                <ChevronLeft className="h-4 w-4" />
                Previous
              </Button>
              <Button variant="outline" size="sm" disabled>
                Next
                <ChevronRight className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
