'use client';

import { useState } from 'react';
import Link from 'next/link';
import { 
  Search, 
  Plus, 
  Filter, 
  MoreHorizontal,
  Star,
  Phone,
  Mail,
  Eye,
  Edit,
  Ban,
  Trash2,
  ChevronLeft,
  ChevronRight,
  Car
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
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
import { Checkbox } from '@/components/ui/checkbox';

// Mock driver data
const mockDrivers = [
  {
    id: '1',
    name: 'Peter Mugisha',
    email: 'peter.m@email.com',
    phone: '+250 788 123 456',
    avatar: null,
    status: 'online',
    rating: 4.8,
    totalTrips: 1234,
    todayTrips: 8,
    todayEarnings: 45000,
    vehicle: { plate: 'RAD 123A', model: 'Toyota Corolla' },
    joinedAt: '2024-06-15',
  },
  {
    id: '2',
    name: 'David Kamanzi',
    email: 'david.k@email.com',
    phone: '+250 788 234 567',
    avatar: null,
    status: 'on_trip',
    rating: 4.6,
    totalTrips: 856,
    todayTrips: 5,
    todayEarnings: 28000,
    vehicle: { plate: 'RAB 456B', model: 'Honda Fit' },
    joinedAt: '2024-08-20',
  },
  {
    id: '3',
    name: 'Jean Pierre',
    email: 'jp@email.com',
    phone: '+250 788 345 678',
    avatar: null,
    status: 'offline',
    rating: 4.9,
    totalTrips: 2341,
    todayTrips: 0,
    todayEarnings: 0,
    vehicle: { plate: 'RAC 789C', model: 'Toyota Vitz' },
    joinedAt: '2024-01-10',
  },
  {
    id: '4',
    name: 'Emmanuel Habimana',
    email: 'emmanuel.h@email.com',
    phone: '+250 788 456 789',
    avatar: null,
    status: 'online',
    rating: 4.7,
    totalTrips: 567,
    todayTrips: 6,
    todayEarnings: 32000,
    vehicle: { plate: 'RAD 012D', model: 'Mazda Demio' },
    joinedAt: '2024-09-05',
  },
  {
    id: '5',
    name: 'Patrick Ndayisaba',
    email: 'patrick.n@email.com',
    phone: '+250 788 567 890',
    avatar: null,
    status: 'suspended',
    rating: 3.2,
    totalTrips: 234,
    todayTrips: 0,
    todayEarnings: 0,
    vehicle: { plate: 'RAE 345E', model: 'Toyota Passo' },
    joinedAt: '2024-10-01',
  },
  {
    id: '6',
    name: 'Claude Uwimana',
    email: 'claude.u@email.com',
    phone: '+250 788 678 901',
    avatar: null,
    status: 'on_trip',
    rating: 4.5,
    totalTrips: 789,
    todayTrips: 4,
    todayEarnings: 22000,
    vehicle: { plate: 'RAF 678F', model: 'Suzuki Swift' },
    joinedAt: '2024-07-22',
  },
];

const getStatusBadge = (status: string) => {
  switch (status) {
    case 'online':
      return (
        <Badge className="bg-green-100 text-green-700 hover:bg-green-100">
          <span className="mr-1.5 h-2 w-2 rounded-full bg-green-500" />
          Online
        </Badge>
      );
    case 'on_trip':
      return (
        <Badge className="bg-yellow-100 text-yellow-700 hover:bg-yellow-100">
          <span className="mr-1.5 h-2 w-2 rounded-full bg-yellow-500" />
          On Trip
        </Badge>
      );
    case 'offline':
      return (
        <Badge className="bg-gray-100 text-gray-700 hover:bg-gray-100">
          <span className="mr-1.5 h-2 w-2 rounded-full bg-gray-400" />
          Offline
        </Badge>
      );
    case 'suspended':
      return (
        <Badge className="bg-red-100 text-red-700 hover:bg-red-100">
          <span className="mr-1.5 h-2 w-2 rounded-full bg-red-500" />
          Suspended
        </Badge>
      );
    default:
      return <Badge variant="secondary">{status}</Badge>;
  }
};

export default function DriversPage() {
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [selectedDrivers, setSelectedDrivers] = useState<string[]>([]);

  const filteredDrivers = mockDrivers.filter((driver) => {
    const matchesSearch =
      driver.name.toLowerCase().includes(search.toLowerCase()) ||
      driver.email.toLowerCase().includes(search.toLowerCase()) ||
      driver.phone.includes(search) ||
      driver.vehicle.plate.toLowerCase().includes(search.toLowerCase());
    const matchesStatus = statusFilter === 'all' || driver.status === statusFilter;
    return matchesSearch && matchesStatus;
  });

  const toggleSelectAll = () => {
    if (selectedDrivers.length === filteredDrivers.length) {
      setSelectedDrivers([]);
    } else {
      setSelectedDrivers(filteredDrivers.map((d) => d.id));
    }
  };

  const toggleSelectDriver = (id: string) => {
    setSelectedDrivers((prev) =>
      prev.includes(id) ? prev.filter((i) => i !== id) : [...prev, id]
    );
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Driver Management</h1>
          <p className="text-muted-foreground">
            Manage and monitor all registered drivers
          </p>
        </div>
        <Button className="gap-2">
          <Plus className="h-4 w-4" />
          Add Driver
        </Button>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Total Drivers</p>
                <p className="text-2xl font-bold">234</p>
              </div>
              <div className="h-10 w-10 rounded-full bg-primary/10 flex items-center justify-center">
                <Car className="h-5 w-5 text-primary" />
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Online Now</p>
                <p className="text-2xl font-bold text-green-600">24</p>
              </div>
              <span className="h-3 w-3 rounded-full bg-green-500 animate-pulse" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">On Trip</p>
                <p className="text-2xl font-bold text-yellow-600">12</p>
              </div>
              <span className="h-3 w-3 rounded-full bg-yellow-500" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Suspended</p>
                <p className="text-2xl font-bold text-red-600">3</p>
              </div>
              <span className="h-3 w-3 rounded-full bg-red-500" />
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
                placeholder="Search by name, email, phone, or plate..."
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
                <SelectItem value="online">Online</SelectItem>
                <SelectItem value="on_trip">On Trip</SelectItem>
                <SelectItem value="offline">Offline</SelectItem>
                <SelectItem value="suspended">Suspended</SelectItem>
              </SelectContent>
            </Select>
            <Button variant="outline" className="gap-2">
              <Filter className="h-4 w-4" />
              More Filters
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Drivers Table */}
      <Card>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-12">
                  <Checkbox
                    checked={selectedDrivers.length === filteredDrivers.length && filteredDrivers.length > 0}
                    onCheckedChange={toggleSelectAll}
                  />
                </TableHead>
                <TableHead>Driver</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Rating</TableHead>
                <TableHead className="text-right">Total Trips</TableHead>
                <TableHead className="text-right">Today&apos;s Trips</TableHead>
                <TableHead className="text-right">Today&apos;s Earnings</TableHead>
                <TableHead>Vehicle</TableHead>
                <TableHead>Joined</TableHead>
                <TableHead className="w-12"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredDrivers.map((driver) => (
                <TableRow key={driver.id} className="group">
                  <TableCell>
                    <Checkbox
                      checked={selectedDrivers.includes(driver.id)}
                      onCheckedChange={() => toggleSelectDriver(driver.id)}
                    />
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center gap-3">
                      <Avatar className="h-10 w-10">
                        <AvatarImage src={driver.avatar || undefined} />
                        <AvatarFallback className="bg-primary/10 text-primary">
                          {driver.name.split(' ').map((n) => n[0]).join('')}
                        </AvatarFallback>
                      </Avatar>
                      <div>
                        <Link
                          href={`/drivers/${driver.id}`}
                          className="font-medium hover:text-primary hover:underline"
                        >
                          {driver.name}
                        </Link>
                        <p className="text-sm text-muted-foreground">{driver.email}</p>
                      </div>
                    </div>
                  </TableCell>
                  <TableCell>{getStatusBadge(driver.status)}</TableCell>
                  <TableCell>
                    <div className="flex items-center gap-1">
                      <Star className="h-4 w-4 fill-yellow-400 text-yellow-400" />
                      <span className="font-medium">{driver.rating}</span>
                    </div>
                  </TableCell>
                  <TableCell className="text-right font-medium">
                    {driver.totalTrips.toLocaleString()}
                  </TableCell>
                  <TableCell className="text-right">{driver.todayTrips}</TableCell>
                  <TableCell className="text-right font-medium">
                    RWF {driver.todayEarnings.toLocaleString()}
                  </TableCell>
                  <TableCell>
                    <div>
                      <p className="font-medium">{driver.vehicle.plate}</p>
                      <p className="text-sm text-muted-foreground">{driver.vehicle.model}</p>
                    </div>
                  </TableCell>
                  <TableCell className="text-muted-foreground">
                    {new Date(driver.joinedAt).toLocaleDateString('en-US', {
                      month: 'short',
                      day: 'numeric',
                      year: 'numeric',
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
                          <Link href={`/drivers/${driver.id}`} className="flex items-center gap-2">
                            <Eye className="h-4 w-4" />
                            View Details
                          </Link>
                        </DropdownMenuItem>
                        <DropdownMenuItem className="gap-2">
                          <Edit className="h-4 w-4" />
                          Edit
                        </DropdownMenuItem>
                        <DropdownMenuItem className="gap-2">
                          <Phone className="h-4 w-4" />
                          Call Driver
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        {driver.status === 'suspended' ? (
                          <DropdownMenuItem className="gap-2 text-green-600">
                            <Eye className="h-4 w-4" />
                            Activate
                          </DropdownMenuItem>
                        ) : (
                          <DropdownMenuItem className="gap-2 text-yellow-600">
                            <Ban className="h-4 w-4" />
                            Suspend
                          </DropdownMenuItem>
                        )}
                        <DropdownMenuItem className="gap-2 text-destructive">
                          <Trash2 className="h-4 w-4" />
                          Delete
                        </DropdownMenuItem>
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
              Showing 1-{filteredDrivers.length} of {mockDrivers.length} drivers
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
