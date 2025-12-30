'use client';

import { useState } from 'react';
import Link from 'next/link';
import { 
  Search, 
  MoreHorizontal,
  Phone,
  Eye,
  Ban,
  ChevronLeft,
  ChevronRight,
  Users,
  UserPlus,
  Filter
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent } from '@/components/ui/card';
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

// Mock user data
const mockUsers = [
  {
    id: '1',
    name: 'John Doe',
    email: 'john.doe@email.com',
    phone: '+250 788 111 222',
    avatar: null,
    status: 'active',
    totalTrips: 45,
    totalSpent: 185000,
    lastTrip: '2024-12-28',
    joinedAt: '2024-03-15',
  },
  {
    id: '2',
    name: 'Jane Smith',
    email: 'jane.smith@email.com',
    phone: '+250 788 222 333',
    avatar: null,
    status: 'active',
    totalTrips: 78,
    totalSpent: 312000,
    lastTrip: '2024-12-28',
    joinedAt: '2024-01-20',
  },
  {
    id: '3',
    name: 'Mike Johnson',
    email: 'mike.j@email.com',
    phone: '+250 788 333 444',
    avatar: null,
    status: 'active',
    totalTrips: 23,
    totalSpent: 92000,
    lastTrip: '2024-12-27',
    joinedAt: '2024-06-10',
  },
  {
    id: '4',
    name: 'Sarah Wilson',
    email: 'sarah.w@email.com',
    phone: '+250 788 444 555',
    avatar: null,
    status: 'suspended',
    totalTrips: 12,
    totalSpent: 48000,
    lastTrip: '2024-12-20',
    joinedAt: '2024-08-05',
  },
  {
    id: '5',
    name: 'Chris Brown',
    email: 'chris.b@email.com',
    phone: '+250 788 555 666',
    avatar: null,
    status: 'active',
    totalTrips: 156,
    totalSpent: 624000,
    lastTrip: '2024-12-28',
    joinedAt: '2023-11-12',
  },
  {
    id: '6',
    name: 'Emma Davis',
    email: 'emma.d@email.com',
    phone: '+250 788 666 777',
    avatar: null,
    status: 'active',
    totalTrips: 8,
    totalSpent: 32000,
    lastTrip: '2024-12-26',
    joinedAt: '2024-12-01',
  },
];

export default function UsersPage() {
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');

  const filteredUsers = mockUsers.filter((user) => {
    const matchesSearch =
      user.name.toLowerCase().includes(search.toLowerCase()) ||
      user.email.toLowerCase().includes(search.toLowerCase()) ||
      user.phone.includes(search);
    const matchesStatus = statusFilter === 'all' || user.status === statusFilter;
    return matchesSearch && matchesStatus;
  });

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-2xl font-bold tracking-tight">User Management</h1>
        <p className="text-muted-foreground">
          Manage and monitor all registered passengers
        </p>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Total Users</p>
                <p className="text-2xl font-bold">1,234</p>
              </div>
              <div className="h-10 w-10 rounded-full bg-primary/10 flex items-center justify-center">
                <Users className="h-5 w-5 text-primary" />
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Active Today</p>
                <p className="text-2xl font-bold text-green-600">89</p>
              </div>
              <span className="h-3 w-3 rounded-full bg-green-500 animate-pulse" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">New This Week</p>
                <p className="text-2xl font-bold text-blue-600">24</p>
              </div>
              <UserPlus className="h-5 w-5 text-blue-600" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Suspended</p>
                <p className="text-2xl font-bold text-red-600">5</p>
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
                placeholder="Search by name, email, or phone..."
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
                <SelectItem value="active">Active</SelectItem>
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

      {/* Users Table */}
      <Card>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>User</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="text-right">Total Trips</TableHead>
                <TableHead className="text-right">Total Spent</TableHead>
                <TableHead>Last Trip</TableHead>
                <TableHead>Joined</TableHead>
                <TableHead className="w-12"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredUsers.map((user) => (
                <TableRow key={user.id}>
                  <TableCell>
                    <div className="flex items-center gap-3">
                      <Avatar className="h-10 w-10">
                        <AvatarImage src={user.avatar || undefined} />
                        <AvatarFallback className="bg-blue-100 text-blue-600">
                          {user.name.split(' ').map((n) => n[0]).join('')}
                        </AvatarFallback>
                      </Avatar>
                      <div>
                        <Link
                          href={`/users/${user.id}`}
                          className="font-medium hover:text-primary hover:underline"
                        >
                          {user.name}
                        </Link>
                        <p className="text-sm text-muted-foreground">{user.email}</p>
                      </div>
                    </div>
                  </TableCell>
                  <TableCell>
                    {user.status === 'active' ? (
                      <Badge className="bg-green-100 text-green-700 hover:bg-green-100">
                        Active
                      </Badge>
                    ) : (
                      <Badge className="bg-red-100 text-red-700 hover:bg-red-100">
                        Suspended
                      </Badge>
                    )}
                  </TableCell>
                  <TableCell className="text-right font-medium">
                    {user.totalTrips}
                  </TableCell>
                  <TableCell className="text-right font-medium">
                    RWF {user.totalSpent.toLocaleString()}
                  </TableCell>
                  <TableCell className="text-muted-foreground">
                    {new Date(user.lastTrip).toLocaleDateString('en-US', {
                      month: 'short',
                      day: 'numeric',
                    })}
                  </TableCell>
                  <TableCell className="text-muted-foreground">
                    {new Date(user.joinedAt).toLocaleDateString('en-US', {
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
                          <Link href={`/users/${user.id}`} className="flex items-center gap-2">
                            <Eye className="h-4 w-4" />
                            View Details
                          </Link>
                        </DropdownMenuItem>
                        <DropdownMenuItem className="gap-2">
                          <Phone className="h-4 w-4" />
                          Call User
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        {user.status === 'suspended' ? (
                          <DropdownMenuItem className="gap-2 text-green-600">
                            <Eye className="h-4 w-4" />
                            Activate
                          </DropdownMenuItem>
                        ) : (
                          <DropdownMenuItem className="gap-2 text-red-600">
                            <Ban className="h-4 w-4" />
                            Suspend
                          </DropdownMenuItem>
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
              Showing 1-{filteredUsers.length} of {mockUsers.length} users
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
