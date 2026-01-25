'use client';

import { useState, useEffect, useCallback } from 'react';
import { 
  Search, 
  Plus, 
  Filter, 
  MoreHorizontal,
  Eye,
  Edit,
  Trash2,
  ChevronLeft,
  ChevronRight,
  Car,
  RefreshCw,
  AlertCircle,
  CheckCircle,
  Download,
  Upload,
  X,
  FileSpreadsheet,
  ImageIcon,
  User
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
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
import type { Vehicle, PageResult, VehicleStatus, VehicleCategory, VehicleLevel } from '@/types';

const getStatusBadge = (status: string) => {
  switch (status) {
    case 'active':
      return (
        <Badge className="bg-green-100 text-green-700 hover:bg-green-100">
          Active
        </Badge>
      );
    case 'inactive':
      return (
        <Badge className="bg-gray-100 text-gray-700 hover:bg-gray-100">
          Inactive
        </Badge>
      );
    case 'maintenance':
      return (
        <Badge className="bg-yellow-100 text-yellow-700 hover:bg-yellow-100">
          Maintenance
        </Badge>
      );
    case 'retired':
      return (
        <Badge className="bg-red-100 text-red-700 hover:bg-red-100">
          Retired
        </Badge>
      );
    default:
      return <Badge variant="secondary">{status}</Badge>;
  }
};

const getLevelBadge = (level: string) => {
  switch (level) {
    case 'economy':
      return <Badge variant="outline" className="text-blue-600 border-blue-300">Economy</Badge>;
    case 'comfort':
      return <Badge variant="outline" className="text-purple-600 border-purple-300">Comfort</Badge>;
    case 'premium':
      return <Badge variant="outline" className="text-amber-600 border-amber-300">Premium</Badge>;
    case 'luxury':
      return <Badge variant="outline" className="text-rose-600 border-rose-300">Luxury</Badge>;
    default:
      return <Badge variant="outline">{level}</Badge>;
  }
};

// Empty vehicle form
const emptyVehicleForm = {
  brand: '',
  model: '',
  plate_number: '',
  year: new Date().getFullYear(),
  color: '',
  seat_capacity: 4,
  category: 'sedan' as VehicleCategory,
  level: 'economy' as VehicleLevel,
  driver_id: '',
  photo_url: '',
};

interface VehicleWithDriver extends Vehicle {
  driver_name?: string;
}

export default function VehiclesPage() {
  const [vehicles, setVehicles] = useState<VehicleWithDriver[]>([]);
  const [drivers, setDrivers] = useState<Array<{ user_id: string; full_name: string }>>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [categoryFilter, setCategoryFilter] = useState('all');
  const [selectedVehicles, setSelectedVehicles] = useState<string[]>([]);
  
  // Pagination
  const [page, setPage] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const limit = 10;
  
  // Stats
  const [stats, setStats] = useState({
    total: 0,
    active: 0,
    maintenance: 0,
    inactive: 0,
  });

  // Modal states
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [isViewModalOpen, setIsViewModalOpen] = useState(false);
  const [isBulkDeleteModalOpen, setIsBulkDeleteModalOpen] = useState(false);
  
  // Photo upload mode: 'url' or 'file'
  const [photoUploadMode, setPhotoUploadMode] = useState<'url' | 'file'>('file');
  const [selectedPhotoFile, setSelectedPhotoFile] = useState<File | null>(null);
  const [photoPreviewUrl, setPhotoPreviewUrl] = useState<string>('');
  const [selectedVehicle, setSelectedVehicle] = useState<VehicleWithDriver | null>(null);
  const [formData, setFormData] = useState(emptyVehicleForm);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const fetchVehicles = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await apiClient.searchVehicles({
        page,
        limit,
        keyword: search || undefined,
        status: statusFilter !== 'all' ? statusFilter as VehicleStatus : undefined,
        category: categoryFilter !== 'all' ? categoryFilter as VehicleCategory : undefined,
      });
      
      const pageResult = response.data as PageResult<VehicleWithDriver>;
      setVehicles(pageResult.records || []);
      setTotalCount(pageResult.count || 0);
      setTotalPages(pageResult.total || 0);
      
      // Calculate stats from records
      const records = pageResult.records || [];
      setStats({
        total: pageResult.count || 0,
        active: records.filter(v => v.status === 'active').length,
        maintenance: records.filter(v => v.status === 'maintenance').length,
        inactive: records.filter(v => v.status === 'inactive' || v.status === 'retired').length,
      });
    } catch (err) {
      console.error('Failed to fetch vehicles:', err);
      setError('Failed to load vehicles. Please try again.');
    } finally {
      setIsLoading(false);
    }
  }, [page, limit, search, statusFilter, categoryFilter]);

  // Fetch drivers for assignment dropdown
  const fetchDrivers = useCallback(async () => {
    try {
      const response = await apiClient.getDrivers({ limit: 100 });
      const pageResult = response.data as PageResult<{ user_id: string; full_name?: string; first_name?: string; last_name?: string }>;
      const driverList = (pageResult.records || []).map(d => ({
        user_id: d.user_id,
        full_name: d.full_name || `${d.first_name || ''} ${d.last_name || ''}`.trim() || 'Unknown Driver',
      }));
      setDrivers(driverList);
    } catch (err) {
      console.error('Failed to fetch drivers:', err);
    }
  }, []);

  useEffect(() => {
    fetchVehicles();
    fetchDrivers();
  }, [fetchVehicles, fetchDrivers]);

  useEffect(() => {
    const timer = setTimeout(() => {
      if (page !== 1) {
        setPage(1);
      } else {
        fetchVehicles();
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
    if (selectedVehicles.length === vehicles.length) {
      setSelectedVehicles([]);
    } else {
      setSelectedVehicles(vehicles.map((v) => v.vehicle_id));
    }
  };

  const toggleSelectVehicle = (id: string) => {
    setSelectedVehicles((prev) =>
      prev.includes(id) ? prev.filter((i) => i !== id) : [...prev, id]
    );
  };

  const clearSelection = () => {
    setSelectedVehicles([]);
  };

  // Handle Add Vehicle
  const handleAddVehicle = async () => {
    if (!formData.brand || !formData.model || !formData.plate_number) {
      setError('Brand, Model, and Plate Number are required');
      return;
    }

    setIsSubmitting(true);
    setError(null);

    try {
      await apiClient.createVehicle({
        brand: formData.brand,
        model: formData.model,
        plate_number: formData.plate_number,
        year: formData.year,
        color: formData.color,
        seat_capacity: formData.seat_capacity,
        category: formData.category,
        level: formData.level,
        driver_id: formData.driver_id || undefined,
        photos: formData.photo_url ? [formData.photo_url] : undefined,
        status: 'active',
      });

      setSuccessMessage('Vehicle added successfully!');
      setIsAddModalOpen(false);
      setFormData(emptyVehicleForm);
      fetchVehicles();
    } catch (err) {
      console.error('Failed to add vehicle:', err);
      setError('Failed to add vehicle. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  // Handle Edit Vehicle
  const handleEditVehicle = async () => {
    if (!selectedVehicle) return;

    setIsSubmitting(true);
    setError(null);

    try {
      await apiClient.updateVehicle(selectedVehicle.vehicle_id, {
        brand: formData.brand,
        model: formData.model,
        plate_number: formData.plate_number,
        year: formData.year,
        color: formData.color,
        seat_capacity: formData.seat_capacity,
        category: formData.category,
        level: formData.level,
        driver_id: formData.driver_id || undefined,
        photos: formData.photo_url ? [formData.photo_url] : undefined,
      });

      setSuccessMessage('Vehicle updated successfully!');
      setIsEditModalOpen(false);
      setSelectedVehicle(null);
      setFormData(emptyVehicleForm);
      fetchVehicles();
    } catch (err) {
      console.error('Failed to update vehicle:', err);
      setError('Failed to update vehicle. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  // Handle Delete Vehicle
  const handleDeleteVehicle = async () => {
    if (!selectedVehicle) return;

    setIsSubmitting(true);
    setError(null);

    try {
      await apiClient.deleteVehicle(selectedVehicle.vehicle_id);
      setSuccessMessage('Vehicle deleted successfully!');
      setIsDeleteModalOpen(false);
      setSelectedVehicle(null);
      fetchVehicles();
    } catch (err) {
      console.error('Failed to delete vehicle:', err);
      setError('Failed to delete vehicle. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  // Handle Status Change
  const handleStatusChange = async (vehicleId: string, newStatus: VehicleStatus) => {
    try {
      await apiClient.updateVehicleStatus(vehicleId, newStatus);
      setSuccessMessage(`Vehicle status updated to ${newStatus}!`);
      fetchVehicles();
    } catch (err) {
      console.error('Failed to update vehicle status:', err);
      setError('Failed to update vehicle status. Please try again.');
    }
  };

  // Handle photo file selection
  const handlePhotoFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) {
      setSelectedPhotoFile(null);
      setPhotoPreviewUrl('');
      return;
    }

    // Validate file type
    const validTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/webp', 'image/gif'];
    if (!validTypes.includes(file.type)) {
      setError('Please upload an image file (JPEG, PNG, WebP, or GIF)');
      return;
    }

    // Validate file size (max 5MB)
    if (file.size > 5 * 1024 * 1024) {
      setError('Image file must be less than 5MB');
      return;
    }

    setSelectedPhotoFile(file);
    
    // Create preview URL
    const previewUrl = URL.createObjectURL(file);
    setPhotoPreviewUrl(previewUrl);
    setFormData({ ...formData, photo_url: '' }); // Clear URL when file is selected
  };

  // Clear photo selection
  const clearPhotoSelection = (alsoResetUrl = true) => {
    setSelectedPhotoFile(null);
    if (photoPreviewUrl) {
      URL.revokeObjectURL(photoPreviewUrl);
    }
    setPhotoPreviewUrl('');
    if (alsoResetUrl) {
      setFormData(prev => ({ ...prev, photo_url: '' }));
    }
  };

  // Get the effective photo URL (from file preview or URL input)
  const getEffectivePhotoUrl = () => {
    if (photoPreviewUrl) return photoPreviewUrl;
    if (formData.photo_url) return formData.photo_url;
    return '';
  };

  // Open Edit Modal
  const openEditModal = (vehicle: VehicleWithDriver) => {
    setSelectedVehicle(vehicle);
    setSelectedPhotoFile(null);
    if (photoPreviewUrl) URL.revokeObjectURL(photoPreviewUrl);
    setPhotoPreviewUrl('');
    setFormData({
      brand: vehicle.brand || '',
      model: vehicle.model || '',
      plate_number: vehicle.plate_number || '',
      year: vehicle.year || new Date().getFullYear(),
      color: vehicle.color || '',
      seat_capacity: vehicle.seat_capacity || 4,
      category: vehicle.category || 'sedan',
      level: vehicle.level || 'economy',
      driver_id: vehicle.driver_id || '',
      photo_url: vehicle.photos?.[0] || '',
    });
    setPhotoUploadMode(vehicle.photos?.[0] ? 'url' : 'file');
    setIsEditModalOpen(true);
  };

  // Open View Modal
  const openViewModal = (vehicle: VehicleWithDriver) => {
    setSelectedVehicle(vehicle);
    setIsViewModalOpen(true);
  };

  // Export to CSV
  const handleExportCSV = () => {
    const headers = ['Brand', 'Model', 'Plate Number', 'Category', 'Level', 'Status', 'Driver', 'Year'];
    const csvContent = [
      headers.join(','),
      ...vehicles.map(vehicle => [
        `"${vehicle.brand || ''}"`,
        `"${vehicle.model || ''}"`,
        `"${vehicle.plate_number || ''}"`,
        vehicle.category || '',
        vehicle.level || '',
        vehicle.status,
        `"${vehicle.driver_name || 'N/A'}"`,
        vehicle.year || '',
      ].join(','))
    ].join('\n');

    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `vehicles_export_${new Date().toISOString().split('T')[0]}.csv`;
    link.click();
    setSuccessMessage('Vehicles exported to CSV!');
  };

  // Bulk Delete
  const handleBulkDelete = async () => {
    if (selectedVehicles.length === 0) return;

    setIsSubmitting(true);
    setError(null);

    try {
      const results = await Promise.allSettled(
        selectedVehicles.map(vehicleId => apiClient.deleteVehicle(vehicleId))
      );

      const successCount = results.filter(r => r.status === 'fulfilled').length;
      const failCount = results.filter(r => r.status === 'rejected').length;

      if (failCount > 0) {
        setError(`${failCount} vehicle(s) failed to delete. ${successCount} deleted successfully.`);
      } else {
        setSuccessMessage(`${successCount} vehicle(s) deleted successfully!`);
      }

      setIsBulkDeleteModalOpen(false);
      setSelectedVehicles([]);
      fetchVehicles();
    } catch (err) {
      console.error('Failed to bulk delete:', err);
      setError('Failed to delete vehicles. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  // Get driver name by ID
  const getDriverName = (driverId?: string) => {
    if (!driverId) return 'N/A';
    const driver = drivers.find(d => d.user_id === driverId);
    return driver?.full_name || 'Unknown Driver';
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Vehicle Management</h1>
          <p className="text-muted-foreground">
            Manage all vehicles in the system
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" size="sm" onClick={handleExportCSV} disabled={vehicles.length === 0}>
            <Download className="h-4 w-4 mr-2" />
            Export CSV
          </Button>
          <Button variant="outline" size="sm" onClick={fetchVehicles} disabled={isLoading}>
            <RefreshCw className={`h-4 w-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
          <Button className="gap-2" onClick={() => { 
            setFormData(emptyVehicleForm); 
            setSelectedPhotoFile(null); 
            if (photoPreviewUrl) URL.revokeObjectURL(photoPreviewUrl);
            setPhotoPreviewUrl('');
            setPhotoUploadMode('file'); 
            setIsAddModalOpen(true); 
          }}>
            <Plus className="h-4 w-4" />
            Add Vehicle
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

      {/* Bulk Actions Bar */}
      {selectedVehicles.length > 0 && (
        <div className="flex items-center justify-between rounded-lg bg-primary/5 border border-primary/20 p-3">
          <div className="flex items-center gap-3">
            <Checkbox
              checked={selectedVehicles.length === vehicles.length && vehicles.length > 0}
              onCheckedChange={toggleSelectAll}
            />
            <span className="text-sm font-medium">
              {selectedVehicles.length} vehicle{selectedVehicles.length > 1 ? 's' : ''} selected
            </span>
            <Button variant="ghost" size="sm" onClick={clearSelection}>
              Clear
            </Button>
          </div>
          <div className="flex items-center gap-2">
            <Button 
              variant="destructive" 
              size="sm" 
              onClick={() => setIsBulkDeleteModalOpen(true)}
            >
              <Trash2 className="h-4 w-4 mr-2" />
              Delete Selected
            </Button>
          </div>
        </div>
      )}

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Total Vehicles</p>
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
                <p className="text-sm text-muted-foreground">Active</p>
                {isLoading ? (
                  <Skeleton className="h-8 w-16 mt-1" />
                ) : (
                  <p className="text-2xl font-bold text-green-600">{stats.active}</p>
                )}
              </div>
              <span className="h-3 w-3 rounded-full bg-green-500" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Maintenance</p>
                {isLoading ? (
                  <Skeleton className="h-8 w-16 mt-1" />
                ) : (
                  <p className="text-2xl font-bold text-yellow-600">{stats.maintenance}</p>
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
                <p className="text-sm text-muted-foreground">Inactive/Retired</p>
                {isLoading ? (
                  <Skeleton className="h-8 w-16 mt-1" />
                ) : (
                  <p className="text-2xl font-bold text-gray-600">{stats.inactive}</p>
                )}
              </div>
              <span className="h-3 w-3 rounded-full bg-gray-400" />
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
                placeholder="Search by plate number, brand, model, or driver..."
                className="pl-10"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
            </div>
            <Select value={statusFilter} onValueChange={(val) => { setStatusFilter(val); setPage(1); }}>
              <SelectTrigger className="w-full md:w-[150px]">
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Status</SelectItem>
                <SelectItem value="active">Active</SelectItem>
                <SelectItem value="inactive">Inactive</SelectItem>
                <SelectItem value="maintenance">Maintenance</SelectItem>
                <SelectItem value="retired">Retired</SelectItem>
              </SelectContent>
            </Select>
            <Select value={categoryFilter} onValueChange={(val) => { setCategoryFilter(val); setPage(1); }}>
              <SelectTrigger className="w-full md:w-[150px]">
                <SelectValue placeholder="Category" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Categories</SelectItem>
                <SelectItem value="sedan">Sedan</SelectItem>
                <SelectItem value="suv">SUV</SelectItem>
                <SelectItem value="mpv">MPV</SelectItem>
                <SelectItem value="van">Van</SelectItem>
                <SelectItem value="hatchback">Hatchback</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {/* Vehicles Table */}
      <Card>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-12">
                  <Checkbox
                    checked={selectedVehicles.length === vehicles.length && vehicles.length > 0}
                    onCheckedChange={toggleSelectAll}
                  />
                </TableHead>
                <TableHead>Photo</TableHead>
                <TableHead>Brand & Model</TableHead>
                <TableHead>Plate Number</TableHead>
                <TableHead>Category</TableHead>
                <TableHead>Level</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Driver</TableHead>
                <TableHead className="w-12"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {isLoading ? (
                Array.from({ length: 5 }).map((_, i) => (
                  <TableRow key={i}>
                    <TableCell><Skeleton className="h-4 w-4" /></TableCell>
                    <TableCell><Skeleton className="h-12 w-16 rounded" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-32" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-16" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-16" /></TableCell>
                    <TableCell><Skeleton className="h-6 w-20" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                    <TableCell><Skeleton className="h-8 w-8" /></TableCell>
                  </TableRow>
                ))
              ) : vehicles.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={9} className="h-32 text-center">
                    <div className="flex flex-col items-center gap-2 text-muted-foreground">
                      <Car className="h-8 w-8" />
                      <p>No vehicles found</p>
                      {search && <p className="text-sm">Try adjusting your search</p>}
                    </div>
                  </TableCell>
                </TableRow>
              ) : (
                vehicles.map((vehicle) => (
                  <TableRow key={vehicle.vehicle_id} className="group">
                    <TableCell>
                      <Checkbox
                        checked={selectedVehicles.includes(vehicle.vehicle_id)}
                        onCheckedChange={() => toggleSelectVehicle(vehicle.vehicle_id)}
                      />
                    </TableCell>
                    <TableCell>
                      <div className="h-12 w-16 rounded bg-gray-100 flex items-center justify-center overflow-hidden">
                        {vehicle.photos?.[0] ? (
                          <img 
                            src={vehicle.photos[0]} 
                            alt={`${vehicle.brand} ${vehicle.model}`}
                            className="h-full w-full object-cover"
                          />
                        ) : (
                          <span className="text-xs text-gray-400 font-medium">N/A</span>
                        )}
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="font-medium">
                        {vehicle.brand} {vehicle.model}
                      </div>
                      {vehicle.year && (
                        <p className="text-xs text-muted-foreground">{vehicle.year}</p>
                      )}
                    </TableCell>
                    <TableCell className="font-mono text-sm">
                      {vehicle.plate_number || '-'}
                    </TableCell>
                    <TableCell className="capitalize">
                      {vehicle.category || '-'}
                    </TableCell>
                    <TableCell>
                      {vehicle.level ? getLevelBadge(vehicle.level) : '-'}
                    </TableCell>
                    <TableCell>
                      {getStatusBadge(vehicle.status)}
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {vehicle.driver_name || getDriverName(vehicle.driver_id)}
                    </TableCell>
                    <TableCell>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="icon" className="h-8 w-8">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem className="gap-2" onClick={() => openViewModal(vehicle)}>
                            <Eye className="h-4 w-4" />
                            View Details
                          </DropdownMenuItem>
                          <DropdownMenuItem className="gap-2" onClick={() => openEditModal(vehicle)}>
                            <Edit className="h-4 w-4" />
                            Edit
                          </DropdownMenuItem>
                          <DropdownMenuSeparator />
                          <DropdownMenuItem 
                            className="gap-2"
                            onClick={() => handleStatusChange(vehicle.vehicle_id, vehicle.status === 'active' ? 'inactive' : 'active')}
                          >
                            {vehicle.status === 'active' ? 'Mark Inactive' : 'Mark Active'}
                          </DropdownMenuItem>
                          <DropdownMenuItem 
                            className="gap-2"
                            onClick={() => handleStatusChange(vehicle.vehicle_id, 'maintenance')}
                          >
                            Set to Maintenance
                          </DropdownMenuItem>
                          <DropdownMenuSeparator />
                          <DropdownMenuItem 
                            className="gap-2 text-destructive"
                            onClick={() => { setSelectedVehicle(vehicle); setIsDeleteModalOpen(true); }}
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
                <>Showing {((page - 1) * limit) + 1}-{Math.min(page * limit, totalCount)} of {totalCount} vehicles</>
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

      {/* Add Vehicle Modal */}
      <Dialog open={isAddModalOpen} onOpenChange={setIsAddModalOpen}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Add New Vehicle</DialogTitle>
            <DialogDescription>
              Enter the vehicle information to register it in the system.
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4 max-h-[60vh] overflow-y-auto">
            {/* Photo Upload */}
            <div className="space-y-3">
              <Label>Vehicle Photo</Label>
              
              {/* Upload Mode Tabs */}
              <div className="flex rounded-lg border p-1 gap-1">
                <button
                  type="button"
                  className={`flex-1 px-3 py-1.5 text-sm rounded-md transition-colors ${photoUploadMode === 'file' ? 'bg-primary text-primary-foreground' : 'hover:bg-muted'}`}
                  onClick={() => setPhotoUploadMode('file')}
                >
                  <Upload className="h-4 w-4 inline mr-1" />
                  Upload File
                </button>
                <button
                  type="button"
                  className={`flex-1 px-3 py-1.5 text-sm rounded-md transition-colors ${photoUploadMode === 'url' ? 'bg-primary text-primary-foreground' : 'hover:bg-muted'}`}
                  onClick={() => setPhotoUploadMode('url')}
                >
                  <ImageIcon className="h-4 w-4 inline mr-1" />
                  URL
                </button>
              </div>

              {/* File Upload Mode */}
              {photoUploadMode === 'file' && (
                <div className="space-y-2">
                  <div className="rounded-lg border-2 border-dashed p-4 text-center">
                    <input
                      type="file"
                      accept="image/jpeg,image/jpg,image/png,image/webp,image/gif"
                      onChange={handlePhotoFileSelect}
                      className="hidden"
                      id="photo-upload"
                    />
                    <label htmlFor="photo-upload" className="cursor-pointer">
                      <Upload className="h-8 w-8 mx-auto text-muted-foreground mb-2" />
                      <p className="text-sm font-medium">
                        {selectedPhotoFile ? selectedPhotoFile.name : 'Click to upload photo'}
                      </p>
                      <p className="text-xs text-muted-foreground mt-1">
                        JPEG, PNG, WebP, GIF up to 5MB
                      </p>
                    </label>
                  </div>
                </div>
              )}

              {/* URL Input Mode */}
              {photoUploadMode === 'url' && (
                <Input
                  id="photo_url"
                  placeholder="https://example.com/vehicle.jpg"
                  value={formData.photo_url}
                  onChange={(e) => { setFormData({ ...formData, photo_url: e.target.value }); setSelectedPhotoFile(null); setPhotoPreviewUrl(''); }}
                />
              )}

              {/* Photo Preview */}
              {getEffectivePhotoUrl() && (
                <div className="relative h-32 w-full rounded-lg border overflow-hidden">
                  <img src={getEffectivePhotoUrl()} alt="Preview" className="h-full w-full object-cover" />
                  <Button
                    variant="destructive"
                    size="icon"
                    className="absolute top-2 right-2 h-6 w-6"
                    onClick={() => clearPhotoSelection()}
                  >
                    <X className="h-4 w-4" />
                  </Button>
                </div>
              )}
            </div>
            
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="brand">Brand *</Label>
                <Input
                  id="brand"
                  placeholder="Toyota"
                  value={formData.brand}
                  onChange={(e) => setFormData({ ...formData, brand: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="model">Model *</Label>
                <Input
                  id="model"
                  placeholder="Corolla"
                  value={formData.model}
                  onChange={(e) => setFormData({ ...formData, model: e.target.value })}
                />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="plate_number">Plate Number *</Label>
                <Input
                  id="plate_number"
                  placeholder="RAJ 123C"
                  value={formData.plate_number}
                  onChange={(e) => setFormData({ ...formData, plate_number: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="seat_capacity">Seat Capacity</Label>
                <Input
                  id="seat_capacity"
                  type="number"
                  min={1}
                  max={20}
                  value={formData.seat_capacity}
                  onChange={(e) => setFormData({ ...formData, seat_capacity: parseInt(e.target.value) || 4 })}
                />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="category">Category *</Label>
                <Select value={formData.category} onValueChange={(v) => setFormData({ ...formData, category: v as VehicleCategory })}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="sedan">Sedan</SelectItem>
                    <SelectItem value="suv">SUV</SelectItem>
                    <SelectItem value="mpv">MPV</SelectItem>
                    <SelectItem value="van">Van</SelectItem>
                    <SelectItem value="hatchback">Hatchback</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label htmlFor="level">Level *</Label>
                <Select value={formData.level} onValueChange={(v) => setFormData({ ...formData, level: v as VehicleLevel })}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="economy">Economy</SelectItem>
                    <SelectItem value="comfort">Comfort</SelectItem>
                    <SelectItem value="premium">Premium</SelectItem>
                    <SelectItem value="luxury">Luxury</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="year">Year</Label>
                <Input
                  id="year"
                  type="number"
                  min={1990}
                  max={new Date().getFullYear() + 1}
                  value={formData.year}
                  onChange={(e) => setFormData({ ...formData, year: parseInt(e.target.value) || new Date().getFullYear() })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="color">Color</Label>
                <Input
                  id="color"
                  placeholder="White"
                  value={formData.color}
                  onChange={(e) => setFormData({ ...formData, color: e.target.value })}
                />
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="driver_id">Assign Driver</Label>
              <Select value={formData.driver_id || "none"} onValueChange={(v) => setFormData({ ...formData, driver_id: v === "none" ? "" : v })}>
                <SelectTrigger>
                  <SelectValue placeholder="Select a driver (optional)" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="none">No Driver</SelectItem>
                  {drivers.map((driver) => (
                    <SelectItem key={driver.user_id} value={driver.user_id}>
                      {driver.full_name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsAddModalOpen(false)} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button onClick={handleAddVehicle} disabled={isSubmitting}>
              {isSubmitting ? 'Adding...' : 'Add Vehicle'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Edit Vehicle Modal */}
      <Dialog open={isEditModalOpen} onOpenChange={setIsEditModalOpen}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Edit Vehicle</DialogTitle>
            <DialogDescription>
              Update vehicle information.
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4 max-h-[60vh] overflow-y-auto">
            {/* Photo Upload */}
            <div className="space-y-3">
              <Label>Vehicle Photo</Label>
              
              {/* Upload Mode Tabs */}
              <div className="flex rounded-lg border p-1 gap-1">
                <button
                  type="button"
                  className={`flex-1 px-3 py-1.5 text-sm rounded-md transition-colors ${photoUploadMode === 'file' ? 'bg-primary text-primary-foreground' : 'hover:bg-muted'}`}
                  onClick={() => setPhotoUploadMode('file')}
                >
                  <Upload className="h-4 w-4 inline mr-1" />
                  Upload File
                </button>
                <button
                  type="button"
                  className={`flex-1 px-3 py-1.5 text-sm rounded-md transition-colors ${photoUploadMode === 'url' ? 'bg-primary text-primary-foreground' : 'hover:bg-muted'}`}
                  onClick={() => setPhotoUploadMode('url')}
                >
                  <ImageIcon className="h-4 w-4 inline mr-1" />
                  URL
                </button>
              </div>

              {/* File Upload Mode */}
              {photoUploadMode === 'file' && (
                <div className="space-y-2">
                  <div className="rounded-lg border-2 border-dashed p-4 text-center">
                    <input
                      type="file"
                      accept="image/jpeg,image/jpg,image/png,image/webp,image/gif"
                      onChange={handlePhotoFileSelect}
                      className="hidden"
                      id="edit-photo-upload"
                    />
                    <label htmlFor="edit-photo-upload" className="cursor-pointer">
                      <Upload className="h-8 w-8 mx-auto text-muted-foreground mb-2" />
                      <p className="text-sm font-medium">
                        {selectedPhotoFile ? selectedPhotoFile.name : 'Click to upload photo'}
                      </p>
                      <p className="text-xs text-muted-foreground mt-1">
                        JPEG, PNG, WebP, GIF up to 5MB
                      </p>
                    </label>
                  </div>
                </div>
              )}

              {/* URL Input Mode */}
              {photoUploadMode === 'url' && (
                <Input
                  id="edit_photo_url"
                  placeholder="https://example.com/vehicle.jpg"
                  value={formData.photo_url}
                  onChange={(e) => { setFormData({ ...formData, photo_url: e.target.value }); setSelectedPhotoFile(null); setPhotoPreviewUrl(''); }}
                />
              )}

              {/* Photo Preview */}
              {getEffectivePhotoUrl() && (
                <div className="relative h-32 w-full rounded-lg border overflow-hidden">
                  <img src={getEffectivePhotoUrl()} alt="Preview" className="h-full w-full object-cover" />
                  <Button
                    variant="destructive"
                    size="icon"
                    className="absolute top-2 right-2 h-6 w-6"
                    onClick={() => clearPhotoSelection()}
                  >
                    <X className="h-4 w-4" />
                  </Button>
                </div>
              )}
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="edit_brand">Brand *</Label>
                <Input
                  id="edit_brand"
                  value={formData.brand}
                  onChange={(e) => setFormData({ ...formData, brand: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="edit_model">Model *</Label>
                <Input
                  id="edit_model"
                  value={formData.model}
                  onChange={(e) => setFormData({ ...formData, model: e.target.value })}
                />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="edit_plate_number">Plate Number *</Label>
                <Input
                  id="edit_plate_number"
                  value={formData.plate_number}
                  onChange={(e) => setFormData({ ...formData, plate_number: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="edit_seat_capacity">Seat Capacity</Label>
                <Input
                  id="edit_seat_capacity"
                  type="number"
                  min={1}
                  max={20}
                  value={formData.seat_capacity}
                  onChange={(e) => setFormData({ ...formData, seat_capacity: parseInt(e.target.value) || 4 })}
                />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="edit_category">Category *</Label>
                <Select value={formData.category} onValueChange={(v) => setFormData({ ...formData, category: v as VehicleCategory })}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="sedan">Sedan</SelectItem>
                    <SelectItem value="suv">SUV</SelectItem>
                    <SelectItem value="mpv">MPV</SelectItem>
                    <SelectItem value="van">Van</SelectItem>
                    <SelectItem value="hatchback">Hatchback</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label htmlFor="edit_level">Level *</Label>
                <Select value={formData.level} onValueChange={(v) => setFormData({ ...formData, level: v as VehicleLevel })}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="economy">Economy</SelectItem>
                    <SelectItem value="comfort">Comfort</SelectItem>
                    <SelectItem value="premium">Premium</SelectItem>
                    <SelectItem value="luxury">Luxury</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="edit_year">Year</Label>
                <Input
                  id="edit_year"
                  type="number"
                  min={1990}
                  max={new Date().getFullYear() + 1}
                  value={formData.year}
                  onChange={(e) => setFormData({ ...formData, year: parseInt(e.target.value) || new Date().getFullYear() })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="edit_color">Color</Label>
                <Input
                  id="edit_color"
                  value={formData.color}
                  onChange={(e) => setFormData({ ...formData, color: e.target.value })}
                />
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="edit_driver_id">Assign Driver</Label>
              <Select value={formData.driver_id || "none"} onValueChange={(v) => setFormData({ ...formData, driver_id: v === "none" ? "" : v })}>
                <SelectTrigger>
                  <SelectValue placeholder="Select a driver (optional)" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="none">No Driver</SelectItem>
                  {drivers.map((driver) => (
                    <SelectItem key={driver.user_id} value={driver.user_id}>
                      {driver.full_name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsEditModalOpen(false)} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button onClick={handleEditVehicle} disabled={isSubmitting}>
              {isSubmitting ? 'Saving...' : 'Save Changes'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* View Vehicle Modal */}
      <Dialog open={isViewModalOpen} onOpenChange={setIsViewModalOpen}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Vehicle Details</DialogTitle>
          </DialogHeader>
          {selectedVehicle && (
            <div className="space-y-4 py-4">
              {/* Photo */}
              <div className="h-40 w-full rounded-lg bg-gray-100 flex items-center justify-center overflow-hidden">
                {selectedVehicle.photos?.[0] ? (
                  <img 
                    src={selectedVehicle.photos[0]} 
                    alt={`${selectedVehicle.brand} ${selectedVehicle.model}`}
                    className="h-full w-full object-cover"
                  />
                ) : (
                  <div className="flex flex-col items-center text-gray-400">
                    <ImageIcon className="h-10 w-10 mb-2" />
                    <span className="text-sm">No Photo</span>
                  </div>
                )}
              </div>
              
              <div className="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <p className="text-muted-foreground">Brand & Model</p>
                  <p className="font-medium">{selectedVehicle.brand} {selectedVehicle.model}</p>
                </div>
                <div>
                  <p className="text-muted-foreground">Plate Number</p>
                  <p className="font-medium font-mono">{selectedVehicle.plate_number || '-'}</p>
                </div>
                <div>
                  <p className="text-muted-foreground">Category</p>
                  <p className="font-medium capitalize">{selectedVehicle.category || '-'}</p>
                </div>
                <div>
                  <p className="text-muted-foreground">Level</p>
                  <div>{selectedVehicle.level ? getLevelBadge(selectedVehicle.level) : '-'}</div>
                </div>
                <div>
                  <p className="text-muted-foreground">Status</p>
                  <div>{getStatusBadge(selectedVehicle.status)}</div>
                </div>
                <div>
                  <p className="text-muted-foreground">Seat Capacity</p>
                  <p className="font-medium">{selectedVehicle.seat_capacity || 4} seats</p>
                </div>
                <div>
                  <p className="text-muted-foreground">Year</p>
                  <p className="font-medium">{selectedVehicle.year || '-'}</p>
                </div>
                <div>
                  <p className="text-muted-foreground">Color</p>
                  <p className="font-medium">{selectedVehicle.color || '-'}</p>
                </div>
                <div className="col-span-2">
                  <p className="text-muted-foreground">Assigned Driver</p>
                  <div className="flex items-center gap-2 mt-1">
                    <User className="h-4 w-4 text-muted-foreground" />
                    <span className="font-medium">
                      {selectedVehicle.driver_name || getDriverName(selectedVehicle.driver_id)}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          )}
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsViewModalOpen(false)}>
              Close
            </Button>
            <Button onClick={() => { setIsViewModalOpen(false); if (selectedVehicle) openEditModal(selectedVehicle); }}>
              <Edit className="h-4 w-4 mr-2" />
              Edit Vehicle
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Modal */}
      <Dialog open={isDeleteModalOpen} onOpenChange={setIsDeleteModalOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="text-destructive">Delete Vehicle</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete <strong>{selectedVehicle?.brand} {selectedVehicle?.model}</strong> ({selectedVehicle?.plate_number})? 
              This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsDeleteModalOpen(false)} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDeleteVehicle} disabled={isSubmitting}>
              {isSubmitting ? 'Deleting...' : 'Delete Vehicle'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Bulk Delete Confirmation Modal */}
      <Dialog open={isBulkDeleteModalOpen} onOpenChange={setIsBulkDeleteModalOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="text-destructive">Delete {selectedVehicles.length} Vehicle(s)</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete {selectedVehicles.length} selected vehicle(s)? 
              This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsBulkDeleteModalOpen(false)} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleBulkDelete} disabled={isSubmitting}>
              {isSubmitting ? 'Deleting...' : `Delete ${selectedVehicles.length} Vehicle(s)`}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
