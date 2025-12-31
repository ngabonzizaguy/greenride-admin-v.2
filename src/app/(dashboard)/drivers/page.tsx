'use client';

import { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { 
  Search, 
  Plus, 
  Filter, 
  MoreHorizontal,
  Star,
  Phone,
  Eye,
  Edit,
  Ban,
  Trash2,
  ChevronLeft,
  ChevronRight,
  Car,
  RefreshCw,
  AlertCircle,
  CheckCircle,
  Download,
  X
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Label } from '@/components/ui/label';
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
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Checkbox } from '@/components/ui/checkbox';
import { Skeleton } from '@/components/ui/skeleton';
import { apiClient } from '@/lib/api-client';
import type { Driver, PageResult, UserStatus } from '@/types';

const getStatusBadge = (status: string, onlineStatus?: string) => {
  if (onlineStatus === 'online') {
    return (
      <Badge className="bg-green-100 text-green-700 hover:bg-green-100">
        <span className="mr-1.5 h-2 w-2 rounded-full bg-green-500" />
        Online
      </Badge>
    );
  }
  if (onlineStatus === 'busy') {
    return (
      <Badge className="bg-yellow-100 text-yellow-700 hover:bg-yellow-100">
        <span className="mr-1.5 h-2 w-2 rounded-full bg-yellow-500" />
        On Trip
      </Badge>
    );
  }
  
  switch (status) {
    case 'active':
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
    case 'banned':
      return (
        <Badge className="bg-red-100 text-red-700 hover:bg-red-100">
          <span className="mr-1.5 h-2 w-2 rounded-full bg-red-700" />
          Banned
        </Badge>
      );
    default:
      return <Badge variant="secondary">{status}</Badge>;
  }
};

// Empty driver form
const emptyDriverForm = {
  first_name: '',
  last_name: '',
  email: '',
  phone: '',
  license_number: '',
  vehicle_plate: '',
  vehicle_type: 'sedan',
};

export default function DriversPage() {
  const [drivers, setDrivers] = useState<Driver[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [selectedDrivers, setSelectedDrivers] = useState<string[]>([]);
  
  // Pagination
  const [page, setPage] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const limit = 10;
  
  // Stats
  const [stats, setStats] = useState({
    total: 0,
    online: 0,
    busy: 0,
    suspended: 0,
  });

  // Modal states
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [isSuspendModalOpen, setIsSuspendModalOpen] = useState(false);
  const [selectedDriver, setSelectedDriver] = useState<Driver | null>(null);
  const [formData, setFormData] = useState(emptyDriverForm);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const fetchDrivers = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await apiClient.getDrivers({
        page,
        limit,
        keyword: search || undefined,
        status: statusFilter !== 'all' ? statusFilter as UserStatus : undefined,
        online_status: statusFilter === 'online' ? 'online' : statusFilter === 'busy' ? 'busy' : undefined,
      });
      
      const pageResult = response.data as PageResult<Driver>;
      setDrivers(pageResult.records || []);
      setTotalCount(pageResult.count || 0);
      setTotalPages(pageResult.total || 0);
      
      if (pageResult.attach) {
        setStats({
          total: (pageResult.attach.total_count as number) || pageResult.count || 0,
          online: (pageResult.attach.online_count as number) || 0,
          busy: (pageResult.attach.busy_count as number) || 0,
          suspended: (pageResult.attach.suspended_count as number) || 0,
        });
      }
    } catch (err) {
      console.error('Failed to fetch drivers:', err);
      setError('Failed to load drivers. Please try again.');
    } finally {
      setIsLoading(false);
    }
  }, [page, limit, search, statusFilter]);

  useEffect(() => {
    fetchDrivers();
  }, [fetchDrivers]);

  useEffect(() => {
    const timer = setTimeout(() => {
      if (page !== 1) {
        setPage(1);
      } else {
        fetchDrivers();
      }
    }, 300);
    return () => clearTimeout(timer);
  }, [search]);

  // Auto-hide success message
  useEffect(() => {
    if (successMessage) {
      const timer = setTimeout(() => setSuccessMessage(null), 3000);
      return () => clearTimeout(timer);
    }
  }, [successMessage]);

  const toggleSelectAll = () => {
    if (selectedDrivers.length === drivers.length) {
      setSelectedDrivers([]);
    } else {
      setSelectedDrivers(drivers.map((d) => d.user_id));
    }
  };

  const toggleSelectDriver = (id: string) => {
    setSelectedDrivers((prev) =>
      prev.includes(id) ? prev.filter((i) => i !== id) : [...prev, id]
    );
  };

  const getDisplayName = (driver: Driver) => {
    return driver.display_name || 
      (driver.first_name && driver.last_name 
        ? `${driver.first_name} ${driver.last_name}` 
        : driver.username || 'Unknown Driver');
  };

  const getInitials = (driver: Driver) => {
    const name = getDisplayName(driver);
    return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2);
  };

  // Handle Add Driver
  const handleAddDriver = async () => {
    if (!formData.first_name || !formData.phone) {
      setError('First name and phone are required');
      return;
    }

    setIsSubmitting(true);
    setError(null);

    try {
      await apiClient.createUser({
        user_type: 'driver',
        first_name: formData.first_name,
        last_name: formData.last_name,
        email: formData.email,
        phone: formData.phone,
        license_number: formData.license_number,
        status: 'active',
      });

      setSuccessMessage('Driver added successfully!');
      setIsAddModalOpen(false);
      setFormData(emptyDriverForm);
      fetchDrivers();
    } catch (err) {
      console.error('Failed to add driver:', err);
      setError('Failed to add driver. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  // Handle Edit Driver
  const handleEditDriver = async () => {
    if (!selectedDriver) return;

    setIsSubmitting(true);
    setError(null);

    try {
      await apiClient.updateUser(selectedDriver.user_id, {
        first_name: formData.first_name,
        last_name: formData.last_name,
        email: formData.email,
        phone: formData.phone,
        license_number: formData.license_number,
      });

      setSuccessMessage('Driver updated successfully!');
      setIsEditModalOpen(false);
      setSelectedDriver(null);
      setFormData(emptyDriverForm);
      fetchDrivers();
    } catch (err) {
      console.error('Failed to update driver:', err);
      setError('Failed to update driver. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  // Handle Delete Driver
  const handleDeleteDriver = async () => {
    if (!selectedDriver) return;

    setIsSubmitting(true);
    setError(null);

    try {
      await apiClient.updateUserStatus(selectedDriver.user_id, 'banned');
      setSuccessMessage('Driver removed successfully!');
      setIsDeleteModalOpen(false);
      setSelectedDriver(null);
      fetchDrivers();
    } catch (err) {
      console.error('Failed to delete driver:', err);
      setError('Failed to remove driver. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  // Handle Suspend/Activate Driver
  const handleToggleSuspend = async () => {
    if (!selectedDriver) return;

    setIsSubmitting(true);
    setError(null);

    const newStatus = selectedDriver.status === 'suspended' ? 'active' : 'suspended';

    try {
      await apiClient.updateUserStatus(selectedDriver.user_id, newStatus);
      setSuccessMessage(`Driver ${newStatus === 'suspended' ? 'suspended' : 'activated'} successfully!`);
      setIsSuspendModalOpen(false);
      setSelectedDriver(null);
      fetchDrivers();
    } catch (err) {
      console.error('Failed to update driver status:', err);
      setError('Failed to update driver status. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  // Open Edit Modal
  const openEditModal = (driver: Driver) => {
    setSelectedDriver(driver);
    setFormData({
      first_name: driver.first_name || '',
      last_name: driver.last_name || '',
      email: driver.email || '',
      phone: driver.phone || '',
      license_number: driver.license_number || '',
      vehicle_plate: '',
      vehicle_type: 'sedan',
    });
    setIsEditModalOpen(true);
  };

  // Export to CSV
  const handleExportCSV = () => {
    const headers = ['Name', 'Email', 'Phone', 'Status', 'Rating', 'Total Rides', 'Joined'];
    const csvContent = [
      headers.join(','),
      ...drivers.map(driver => [
        `"${getDisplayName(driver)}"`,
        driver.email || '',
        driver.phone || '',
        driver.status,
        driver.score?.toFixed(1) || '5.0',
        driver.total_rides || 0,
        driver.created_at ? new Date(driver.created_at).toLocaleDateString() : '',
      ].join(','))
    ].join('\n');

    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `drivers_export_${new Date().toISOString().split('T')[0]}.csv`;
    link.click();
    setSuccessMessage('Drivers exported to CSV!');
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
        <div className="flex gap-2">
          <Button variant="outline" size="sm" onClick={handleExportCSV} disabled={drivers.length === 0}>
            <Download className="h-4 w-4 mr-2" />
            Export CSV
          </Button>
          <Button variant="outline" size="sm" onClick={fetchDrivers} disabled={isLoading}>
            <RefreshCw className={`h-4 w-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
          <Button className="gap-2" onClick={() => { setFormData(emptyDriverForm); setIsAddModalOpen(true); }}>
            <Plus className="h-4 w-4" />
            Add Driver
          </Button>
        </div>
      </div>

      {/* Success Banner */}
      {successMessage && (
        <div className="flex items-center gap-2 rounded-lg bg-green-50 border border-green-200 p-3 text-sm text-green-800">
          <CheckCircle className="h-4 w-4 flex-shrink-0" />
          <span>{successMessage}</span>
        </div>
      )}

      {/* Error Banner */}
      {error && (
        <div className="flex items-center gap-2 rounded-lg bg-red-50 border border-red-200 p-3 text-sm text-red-800">
          <AlertCircle className="h-4 w-4 flex-shrink-0" />
          <span>{error}</span>
          <Button variant="ghost" size="sm" className="ml-auto h-6 w-6 p-0" onClick={() => setError(null)}>
            <X className="h-4 w-4" />
          </Button>
        </div>
      )}

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Total Drivers</p>
                {isLoading ? (
                  <Skeleton className="h-8 w-16 mt-1" />
                ) : (
                  <p className="text-2xl font-bold">{stats.total || totalCount}</p>
                )}
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
                {isLoading ? (
                  <Skeleton className="h-8 w-16 mt-1" />
                ) : (
                  <p className="text-2xl font-bold text-green-600">{stats.online}</p>
                )}
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
                {isLoading ? (
                  <Skeleton className="h-8 w-16 mt-1" />
                ) : (
                  <p className="text-2xl font-bold text-yellow-600">{stats.busy}</p>
                )}
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
                {isLoading ? (
                  <Skeleton className="h-8 w-16 mt-1" />
                ) : (
                  <p className="text-2xl font-bold text-red-600">{stats.suspended}</p>
                )}
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
            <Select value={statusFilter} onValueChange={(val) => { setStatusFilter(val); setPage(1); }}>
              <SelectTrigger className="w-full md:w-[180px]">
                <SelectValue placeholder="Filter by status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Status</SelectItem>
                <SelectItem value="active">Active</SelectItem>
                <SelectItem value="online">Online</SelectItem>
                <SelectItem value="busy">On Trip</SelectItem>
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
                    checked={selectedDrivers.length === drivers.length && drivers.length > 0}
                    onCheckedChange={toggleSelectAll}
                  />
                </TableHead>
                <TableHead>Driver</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Rating</TableHead>
                <TableHead className="text-right">Total Rides</TableHead>
                <TableHead>Phone</TableHead>
                <TableHead>Joined</TableHead>
                <TableHead className="w-12"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {isLoading ? (
                Array.from({ length: 5 }).map((_, i) => (
                  <TableRow key={i}>
                    <TableCell><Skeleton className="h-4 w-4" /></TableCell>
                    <TableCell>
                      <div className="flex items-center gap-3">
                        <Skeleton className="h-10 w-10 rounded-full" />
                        <div>
                          <Skeleton className="h-4 w-32 mb-1" />
                          <Skeleton className="h-3 w-24" />
                        </div>
                      </div>
                    </TableCell>
                    <TableCell><Skeleton className="h-6 w-20" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-12" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-16" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                    <TableCell><Skeleton className="h-8 w-8" /></TableCell>
                  </TableRow>
                ))
              ) : drivers.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={8} className="h-32 text-center">
                    <div className="flex flex-col items-center gap-2 text-muted-foreground">
                      <Car className="h-8 w-8" />
                      <p>No drivers found</p>
                      {search && <p className="text-sm">Try adjusting your search</p>}
                    </div>
                  </TableCell>
                </TableRow>
              ) : (
                drivers.map((driver) => (
                  <TableRow key={driver.user_id} className="group">
                    <TableCell>
                      <Checkbox
                        checked={selectedDrivers.includes(driver.user_id)}
                        onCheckedChange={() => toggleSelectDriver(driver.user_id)}
                      />
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-3">
                        <Avatar className="h-10 w-10">
                          <AvatarImage src={driver.avatar || undefined} />
                          <AvatarFallback className="bg-primary/10 text-primary">
                            {getInitials(driver)}
                          </AvatarFallback>
                        </Avatar>
                        <div>
                          <Link
                            href={`/drivers/${driver.user_id}`}
                            className="font-medium hover:text-primary hover:underline"
                          >
                            {getDisplayName(driver)}
                          </Link>
                          <p className="text-sm text-muted-foreground">{driver.email || 'No email'}</p>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell>{getStatusBadge(driver.status, driver.online_status)}</TableCell>
                    <TableCell>
                      <div className="flex items-center gap-1">
                        <Star className="h-4 w-4 fill-yellow-400 text-yellow-400" />
                        <span className="font-medium">{driver.score?.toFixed(1) || '5.0'}</span>
                      </div>
                    </TableCell>
                    <TableCell className="text-right font-medium">
                      {(driver.total_rides || 0).toLocaleString()}
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {driver.phone || 'No phone'}
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {driver.created_at 
                        ? new Date(driver.created_at).toLocaleDateString('en-US', {
                            month: 'short',
                            day: 'numeric',
                            year: 'numeric',
                          })
                        : '-'}
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
                            <Link href={`/drivers/${driver.user_id}`} className="flex items-center gap-2">
                              <Eye className="h-4 w-4" />
                              View Details
                            </Link>
                          </DropdownMenuItem>
                          <DropdownMenuItem className="gap-2" onClick={() => openEditModal(driver)}>
                            <Edit className="h-4 w-4" />
                            Edit
                          </DropdownMenuItem>
                          {driver.phone && (
                            <DropdownMenuItem className="gap-2" asChild>
                              <a href={`tel:${driver.phone}`}>
                                <Phone className="h-4 w-4" />
                                Call Driver
                              </a>
                            </DropdownMenuItem>
                          )}
                          <DropdownMenuSeparator />
                          <DropdownMenuItem 
                            className={driver.status === 'suspended' ? 'gap-2 text-green-600' : 'gap-2 text-yellow-600'}
                            onClick={() => { setSelectedDriver(driver); setIsSuspendModalOpen(true); }}
                          >
                            {driver.status === 'suspended' ? (
                              <>
                                <CheckCircle className="h-4 w-4" />
                                Activate
                              </>
                            ) : (
                              <>
                                <Ban className="h-4 w-4" />
                                Suspend
                              </>
                            )}
                          </DropdownMenuItem>
                          <DropdownMenuItem 
                            className="gap-2 text-destructive"
                            onClick={() => { setSelectedDriver(driver); setIsDeleteModalOpen(true); }}
                          >
                            <Trash2 className="h-4 w-4" />
                            Delete
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>

          {/* Pagination */}
          <div className="flex items-center justify-between border-t px-4 py-4">
            <div className="text-sm text-muted-foreground">
              {isLoading ? (
                <Skeleton className="h-4 w-32" />
              ) : (
                <>Showing {((page - 1) * limit) + 1}-{Math.min(page * limit, totalCount)} of {totalCount} drivers</>
              )}
            </div>
            <div className="flex items-center gap-2">
              <Button 
                variant="outline" 
                size="sm" 
                disabled={page === 1 || isLoading}
                onClick={() => setPage(p => p - 1)}
              >
                <ChevronLeft className="h-4 w-4" />
                Previous
              </Button>
              <span className="text-sm text-muted-foreground">
                Page {page} of {totalPages || 1}
              </span>
              <Button 
                variant="outline" 
                size="sm" 
                disabled={page >= totalPages || isLoading}
                onClick={() => setPage(p => p + 1)}
              >
                Next
                <ChevronRight className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Add Driver Modal */}
      <Dialog open={isAddModalOpen} onOpenChange={setIsAddModalOpen}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Add New Driver</DialogTitle>
            <DialogDescription>
              Enter the driver&apos;s information to register them in the system.
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="first_name">First Name *</Label>
                <Input
                  id="first_name"
                  placeholder="John"
                  value={formData.first_name}
                  onChange={(e) => setFormData({ ...formData, first_name: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="last_name">Last Name</Label>
                <Input
                  id="last_name"
                  placeholder="Doe"
                  value={formData.last_name}
                  onChange={(e) => setFormData({ ...formData, last_name: e.target.value })}
                />
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="phone">Phone Number *</Label>
              <Input
                id="phone"
                placeholder="+250 788 123 456"
                value={formData.phone}
                onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="email">Email Address</Label>
              <Input
                id="email"
                type="email"
                placeholder="john@example.com"
                value={formData.email}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="license">License Number</Label>
              <Input
                id="license"
                placeholder="DL-123456"
                value={formData.license_number}
                onChange={(e) => setFormData({ ...formData, license_number: e.target.value })}
              />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="vehicle_plate">Vehicle Plate</Label>
                <Input
                  id="vehicle_plate"
                  placeholder="RAD 123A"
                  value={formData.vehicle_plate}
                  onChange={(e) => setFormData({ ...formData, vehicle_plate: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="vehicle_type">Vehicle Type</Label>
                <Select value={formData.vehicle_type} onValueChange={(v) => setFormData({ ...formData, vehicle_type: v })}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="sedan">Sedan</SelectItem>
                    <SelectItem value="suv">SUV</SelectItem>
                    <SelectItem value="moto">Moto</SelectItem>
                    <SelectItem value="van">Van</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsAddModalOpen(false)} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button onClick={handleAddDriver} disabled={isSubmitting}>
              {isSubmitting ? 'Adding...' : 'Add Driver'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Edit Driver Modal */}
      <Dialog open={isEditModalOpen} onOpenChange={setIsEditModalOpen}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Edit Driver</DialogTitle>
            <DialogDescription>
              Update driver information.
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="edit_first_name">First Name</Label>
                <Input
                  id="edit_first_name"
                  value={formData.first_name}
                  onChange={(e) => setFormData({ ...formData, first_name: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="edit_last_name">Last Name</Label>
                <Input
                  id="edit_last_name"
                  value={formData.last_name}
                  onChange={(e) => setFormData({ ...formData, last_name: e.target.value })}
                />
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="edit_phone">Phone Number</Label>
              <Input
                id="edit_phone"
                value={formData.phone}
                onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="edit_email">Email Address</Label>
              <Input
                id="edit_email"
                type="email"
                value={formData.email}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="edit_license">License Number</Label>
              <Input
                id="edit_license"
                value={formData.license_number}
                onChange={(e) => setFormData({ ...formData, license_number: e.target.value })}
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsEditModalOpen(false)} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button onClick={handleEditDriver} disabled={isSubmitting}>
              {isSubmitting ? 'Saving...' : 'Save Changes'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Modal */}
      <Dialog open={isDeleteModalOpen} onOpenChange={setIsDeleteModalOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="text-destructive">Delete Driver</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete <strong>{selectedDriver && getDisplayName(selectedDriver)}</strong>? 
              This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsDeleteModalOpen(false)} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDeleteDriver} disabled={isSubmitting}>
              {isSubmitting ? 'Deleting...' : 'Delete Driver'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Suspend Confirmation Modal */}
      <Dialog open={isSuspendModalOpen} onOpenChange={setIsSuspendModalOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>
              {selectedDriver?.status === 'suspended' ? 'Activate Driver' : 'Suspend Driver'}
            </DialogTitle>
            <DialogDescription>
              {selectedDriver?.status === 'suspended' 
                ? `Are you sure you want to activate ${selectedDriver && getDisplayName(selectedDriver)}? They will be able to receive ride requests again.`
                : `Are you sure you want to suspend ${selectedDriver && getDisplayName(selectedDriver)}? They will not be able to receive ride requests.`
              }
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsSuspendModalOpen(false)} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button 
              variant={selectedDriver?.status === 'suspended' ? 'default' : 'destructive'}
              onClick={handleToggleSuspend} 
              disabled={isSubmitting}
            >
              {isSubmitting ? 'Processing...' : selectedDriver?.status === 'suspended' ? 'Activate' : 'Suspend'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
