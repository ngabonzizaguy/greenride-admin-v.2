'use client';

import { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { 
  Search, 
  Plus, 
  Filter, 
  MoreHorizontal,
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
  Upload,
  X,
  FileSpreadsheet
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
import { apiClient, ApiError } from '@/lib/api-client';
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
  vehicle_id: '',
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
  const [isBulkDeleteModalOpen, setIsBulkDeleteModalOpen] = useState(false);
  const [isBulkSuspendModalOpen, setIsBulkSuspendModalOpen] = useState(false);
  const [isUploadModalOpen, setIsUploadModalOpen] = useState(false);
  const [selectedDriver, setSelectedDriver] = useState<Driver | null>(null);
  const [formData, setFormData] = useState(emptyDriverForm);
  const [isSubmitting, setIsSubmitting] = useState(false);
  
  // CSV Upload states
  const [uploadFile, setUploadFile] = useState<File | null>(null);
  const [uploadPreview, setUploadPreview] = useState<Array<Record<string, string>>>([]);
  const [uploadError, setUploadError] = useState<string | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  
  // Vehicles list for assignment
  const [vehicles, setVehicles] = useState<Array<{ vehicle_id: string; display_name: string; plate_number: string; driver_id?: string }>>([]);

  const fetchVehicles = useCallback(async () => {
    try {
      const response = await apiClient.searchVehicles({ limit: 100 });
      const pageResult = response.data as PageResult<{ vehicle_id: string; brand?: string; model?: string; plate_number?: string; driver_id?: string }>;
      const vehicleList = (pageResult.records || []).map(v => ({
        vehicle_id: v.vehicle_id,
        display_name: `${v.brand || ''} ${v.model || ''}`.trim() || 'Unknown Vehicle',
        plate_number: v.plate_number || '',
        driver_id: v.driver_id,
      }));
      setVehicles(vehicleList);
    } catch (err) {
      console.error('Failed to fetch vehicles:', err);
    }
  }, []);

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
    fetchVehicles();
  }, [fetchDrivers, fetchVehicles]);

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

  const getDriverVehicleLabel = (driver: Driver) => {
    const v = driver.vehicle;
    if (v && (v.plate_number || v.brand || v.model)) {
      const makeModel = `${v.brand || ''} ${v.model || ''}`.trim();
      const plate = v.plate_number ? v.plate_number : '';
      return [plate, makeModel].filter(Boolean).join(' • ') || plate || makeModel || 'Assigned';
    }

    const assigned = vehicles.find(vv => vv.driver_id === driver.user_id);
    if (assigned) {
      const makeModel = assigned.display_name || '';
      const plate = assigned.plate_number || '';
      return [plate, makeModel].filter(Boolean).join(' • ') || plate || makeModel || 'Assigned';
    }

    return 'Unassigned';
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
      const response = await apiClient.createUser({
        user_type: 'driver',
        first_name: formData.first_name,
        last_name: formData.last_name,
        email: formData.email,
        phone: formData.phone,
        license_number: formData.license_number,
        status: 'active',
      });

      // If a vehicle is selected, assign it to the new driver
      if (formData.vehicle_id && response.code === '0000') {
        const driverData = response.data as { user_id?: string };
        if (driverData?.user_id) {
          // Check if vehicle is already assigned to another driver
          const vehicleToAssign = vehicles.find(v => v.vehicle_id === formData.vehicle_id);
          if (vehicleToAssign?.driver_id && vehicleToAssign.driver_id !== driverData.user_id) {
            // Unassign the old driver first (optional - could show warning instead)
            try {
              await apiClient.updateVehicle(formData.vehicle_id, {
                driver_id: driverData.user_id,
              });
            } catch (vehicleErr) {
              console.error('Failed to assign vehicle:', vehicleErr);
            }
          } else {
            // Vehicle is free, assign it
            try {
              await apiClient.updateVehicle(formData.vehicle_id, {
                driver_id: driverData.user_id,
              });
            } catch (vehicleErr) {
              console.error('Failed to assign vehicle:', vehicleErr);
            }
          }
        }
      }

      setSuccessMessage('Driver added successfully!');
      setIsAddModalOpen(false);
      setFormData(emptyDriverForm);
      fetchDrivers();
      fetchVehicles(); // Refresh vehicles list
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
      // Update driver info
      await apiClient.updateUser(selectedDriver.user_id, {
        first_name: formData.first_name,
        last_name: formData.last_name,
        email: formData.email,
        phone: formData.phone,
        license_number: formData.license_number,
      });

      // Handle vehicle assignment/unassignment
      // First, find if driver currently has a vehicle assigned
      const currentVehicle = vehicles.find(v => v.driver_id === selectedDriver.user_id);
      
      // If a new vehicle is selected
      if (formData.vehicle_id) {
        // If driver had a different vehicle, unassign it first
        if (currentVehicle && currentVehicle.vehicle_id !== formData.vehicle_id) {
          try {
            await apiClient.updateVehicle(currentVehicle.vehicle_id, {
              driver_id: '',
            });
          } catch (unassignErr) {
            console.error('Failed to unassign old vehicle:', unassignErr);
            // Continue anyway
          }
        }
        
        // Assign the new vehicle
        try {
          await apiClient.updateVehicle(formData.vehicle_id, {
            driver_id: selectedDriver.user_id,
          });
        } catch (vehicleErr) {
          console.error('Failed to assign vehicle:', vehicleErr);
          // Don't fail the whole operation, just log
        }
      } else if (currentVehicle) {
        // If no vehicle selected but driver had one, unassign it
        try {
          await apiClient.updateVehicle(currentVehicle.vehicle_id, {
            driver_id: '',
          });
        } catch (unassignErr) {
          console.error('Failed to unassign vehicle:', unassignErr);
          // Continue anyway
        }
      }

      setSuccessMessage('Driver updated successfully!');
      setIsEditModalOpen(false);
      setSelectedDriver(null);
      setFormData(emptyDriverForm);
      fetchDrivers();
      fetchVehicles(); // Refresh vehicles list
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
      await apiClient.deleteUser(selectedDriver.user_id, 'Deleted by admin');
      setSuccessMessage('Driver deleted successfully!');
      setIsDeleteModalOpen(false);
      setSelectedDriver(null);
      fetchDrivers();
    } catch (err) {
      console.error('Failed to delete driver:', err);
      const message =
        err instanceof ApiError
          ? (err.serverMessage || err.message)
          : (err instanceof Error ? err.message : 'Failed to delete driver.');
      setError(message);
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
    // Find if this driver has a vehicle assigned by checking driver_id on vehicles
    const assignedVehicle = vehicles.find(v => v.driver_id === driver.user_id);
    setFormData({
      first_name: driver.first_name || '',
      last_name: driver.last_name || '',
      email: driver.email || '',
      phone: driver.phone || '',
      license_number: driver.license_number || '',
      vehicle_id: assignedVehicle?.vehicle_id || driver.vehicle?.vehicle_id || '',
    });
    setIsEditModalOpen(true);
  };

  // Export to CSV
  const handleExportCSV = () => {
    const headers = ['Name', 'Email', 'Phone', 'Status', 'Vehicle', 'Total Rides', 'Joined'];
    const csvContent = [
      headers.join(','),
      ...drivers.map(driver => [
        `"${getDisplayName(driver)}"`,
        driver.email || '',
        driver.phone || '',
        driver.status,
        `"${getDriverVehicleLabel(driver)}"`,
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

  // Bulk Delete Drivers
  const handleBulkDelete = async () => {
    if (selectedDrivers.length === 0) return;

    setIsSubmitting(true);
    setError(null);

    try {
      // Delete each selected driver
      const results = await Promise.allSettled(
        selectedDrivers.map(userId => apiClient.deleteUser(userId, 'Bulk deleted by admin'))
      );

      const successCount = results.filter(r => r.status === 'fulfilled').length;
      const failCount = results.filter(r => r.status === 'rejected').length;

      if (failCount > 0) {
        setError(`${failCount} driver(s) failed to delete. ${successCount} deleted successfully.`);
      } else {
        setSuccessMessage(`${successCount} driver(s) deleted successfully!`);
      }

      setIsBulkDeleteModalOpen(false);
      setSelectedDrivers([]);
      fetchDrivers();
    } catch (err) {
      console.error('Failed to bulk delete:', err);
      const message =
        err instanceof ApiError
          ? (err.serverMessage || err.message)
          : (err instanceof Error ? err.message : 'Failed to delete drivers.');
      setError(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  // Bulk Suspend/Activate Drivers
  const handleBulkSuspend = async (suspend: boolean) => {
    if (selectedDrivers.length === 0) return;

    setIsSubmitting(true);
    setError(null);

    const newStatus = suspend ? 'suspended' : 'active';

    try {
      const results = await Promise.allSettled(
        selectedDrivers.map(userId => apiClient.updateUserStatus(userId, newStatus))
      );

      const successCount = results.filter(r => r.status === 'fulfilled').length;
      const failCount = results.filter(r => r.status === 'rejected').length;

      if (failCount > 0) {
        setError(`${failCount} driver(s) failed to ${suspend ? 'suspend' : 'activate'}. ${successCount} updated successfully.`);
      } else {
        setSuccessMessage(`${successCount} driver(s) ${suspend ? 'suspended' : 'activated'} successfully!`);
      }

      setIsBulkSuspendModalOpen(false);
      setSelectedDrivers([]);
      fetchDrivers();
    } catch (err) {
      console.error('Failed to bulk update status:', err);
      setError(`Failed to ${suspend ? 'suspend' : 'activate'} drivers. Please try again.`);
    } finally {
      setIsSubmitting(false);
    }
  };

  // Parse CSV file - handles various formats including the Excel format with Name, Phone, Plate
  const parseCSV = (text: string): Array<Record<string, string>> => {
    const lines = text.split('\n').filter(line => line.trim());
    if (lines.length < 2) return [];

    const headers = lines[0].split(',').map(h => h.trim().toLowerCase().replace(/["']/g, ''));
    const rows: Array<Record<string, string>> = [];

    // Check if this is the dual-driver format (Names/1st Driver, Names/2nd Driver)
    const isDualDriverFormat = headers.some(h => h.includes('1st') || h.includes('2nd'));

    for (let i = 1; i < lines.length; i++) {
      const values = lines[i].split(',').map(v => v.trim().replace(/["']/g, ''));
      
      if (isDualDriverFormat) {
        // Handle dual-driver format: Names/1st Driver | Phone | Names/2nd Driver | Phone | Plate
        const row1: Record<string, string> = {};
        const row2: Record<string, string> = {};
        
        headers.forEach((header, index) => {
          const value = values[index] || '';
          if (header.includes('1st') || (header === 'phone number' && index === 1)) {
            // First driver columns
            if (header.includes('driver') || header.includes('name')) {
              row1['name'] = value;
            } else if (header.includes('phone') || header === 'phone number') {
              row1['phone'] = value;
            }
          } else if (header.includes('2nd') || (header === 'phone number' && index === 3)) {
            // Second driver columns  
            if (header.includes('driver') || header.includes('name')) {
              row2['name'] = value;
            } else if (header.includes('phone') || header === 'phone number') {
              row2['phone'] = value;
            }
          } else if (header.includes('plate')) {
            // Plate number shared by both drivers
            row1['plate_number'] = value;
            row2['plate_number'] = value;
          }
        });
        
        // Only add rows with names
        if (row1['name']) rows.push(row1);
        if (row2['name']) rows.push(row2);
      } else {
        // Standard format
        const row: Record<string, string> = {};
        headers.forEach((header, index) => {
          row[header] = values[index] || '';
        });
        rows.push(row);
      }
    }

    return rows;
  };

  // Format phone number with +250 prefix if needed
  const formatPhoneNumber = (phone: string): string => {
    if (!phone) return '';
    // Remove spaces, dashes, and other non-digit characters except +
    let cleaned = phone.replace(/[^\d+]/g, '');
    // Add +250 prefix if it's a 9-digit Rwandan number
    if (/^\d{9}$/.test(cleaned)) {
      cleaned = '+250' + cleaned;
    } else if (/^0\d{9}$/.test(cleaned)) {
      cleaned = '+250' + cleaned.slice(1);
    } else if (/^250\d{9}$/.test(cleaned)) {
      cleaned = '+' + cleaned;
    }
    return cleaned;
  };

  // Split full name into first and last name
  const splitName = (fullName: string): { firstName: string; lastName: string } => {
    const parts = fullName.trim().split(/\s+/);
    if (parts.length === 1) {
      return { firstName: parts[0], lastName: '' };
    }
    const firstName = parts[0];
    const lastName = parts.slice(1).join(' ');
    return { firstName, lastName };
  };

  // Handle file selection
  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    setUploadError(null);
    setUploadPreview([]);

    if (!file) {
      setUploadFile(null);
      return;
    }

    // Validate file type
    const validTypes = ['text/csv', 'application/vnd.ms-excel', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'];
    if (!validTypes.includes(file.type) && !file.name.endsWith('.csv')) {
      setUploadError('Please upload a CSV file');
      return;
    }

    setUploadFile(file);

    // Read and preview the file
    const reader = new FileReader();
    reader.onload = (event) => {
      try {
        const text = event.target?.result as string;
        const parsed = parseCSV(text);
        if (parsed.length === 0) {
          setUploadError('No valid data found in file');
          return;
        }
        setUploadPreview(parsed.slice(0, 5)); // Preview first 5 rows
      } catch {
        setUploadError('Failed to parse file');
      }
    };
    reader.readAsText(file);
  };

  // Handle CSV Upload - supports various formats including Excel driver list
  const handleUploadCSV = async () => {
    if (!uploadFile) return;

    setIsUploading(true);
    setUploadError(null);

    try {
      const reader = new FileReader();
      reader.onload = async (event) => {
        try {
          const text = event.target?.result as string;
          const rows = parseCSV(text);

          if (rows.length === 0) {
            setUploadError('No valid data found in file');
            setIsUploading(false);
            return;
          }

          // Track created vehicles for plate numbers
          const plateToVehicleId: Record<string, string> = {};

          // Process each row - create drivers and optionally vehicles
          const results = await Promise.allSettled(
            rows.map(async (row) => {
              // Get name - try various column formats
              const fullName = row['name'] || row['names/1st driver'] || row['names/2nd driver'] || 
                             row['first name'] || row['first_name'] || row['firstname'] || '';
              const { firstName, lastName } = splitName(fullName);
              
              // Get phone - format with +250 prefix if needed
              const rawPhone = row['phone'] || row['phone number'] || row['phone_number'] || '';
              const phone = formatPhoneNumber(rawPhone);
              
              // Get plate number if available
              const plateNumber = row['plate_number'] || row['plate number'] || row['plate'] || '';
              
              // Skip rows without name or phone
              if (!firstName && !phone) {
                return Promise.reject(new Error('Missing name and phone'));
              }
              
              // Create the driver
              const driverResponse = await apiClient.createUser({
                user_type: 'driver',
                first_name: firstName,
                last_name: lastName,
                email: row['email'] || '',
                phone: phone,
                license_number: row['license'] || row['license number'] || row['license_number'] || '',
                status: 'active',
              });
              
              // If plate number exists, create or link vehicle
              if (plateNumber && driverResponse.code === '0000') {
                const driverData = driverResponse.data as { user_id?: string };
                const driverId = driverData?.user_id;
                
                // Check if we already created this vehicle
                if (!plateToVehicleId[plateNumber] && driverId) {
                  try {
                    // Create vehicle with the plate number
                    const vehicleResponse = await apiClient.createVehicle({
                      plate_number: plateNumber,
                      brand: 'Unknown', // Can be updated later
                      model: 'Unknown',
                      category: 'sedan',
                      level: 'economy',
                      driver_id: driverId,
                      status: 'active',
                    });
                    
                    if (vehicleResponse.code === '0000') {
                      const vehicleData = vehicleResponse.data as { vehicle_id?: string };
                      if (vehicleData?.vehicle_id) {
                        plateToVehicleId[plateNumber] = vehicleData.vehicle_id;
                      }
                    }
                  } catch (e) {
                    // Vehicle creation failed, but driver was created successfully
                    console.error('Failed to create vehicle for plate:', plateNumber, e);
                  }
                }
              }
              
              return driverResponse;
            })
          );

          const successCount = results.filter(r => r.status === 'fulfilled').length;
          const failCount = results.filter(r => r.status === 'rejected').length;
          const vehiclesCreated = Object.keys(plateToVehicleId).length;

          if (failCount > 0) {
            setError(`${failCount} driver(s) failed to import. ${successCount} imported successfully.${vehiclesCreated > 0 ? ` ${vehiclesCreated} vehicle(s) created.` : ''}`);
          } else {
            setSuccessMessage(`${successCount} driver(s) imported successfully!${vehiclesCreated > 0 ? ` ${vehiclesCreated} vehicle(s) created.` : ''}`);
          }

          setIsUploadModalOpen(false);
          setUploadFile(null);
          setUploadPreview([]);
          fetchDrivers();
        } catch {
          setUploadError('Failed to process file');
        } finally {
          setIsUploading(false);
        }
      };
      reader.readAsText(uploadFile);
    } catch {
      setUploadError('Failed to upload file');
      setIsUploading(false);
    }
  };

  // Download CSV template - supports multiple formats
  const handleDownloadTemplate = () => {
    // Template with simpler format matching user's Excel structure
    const template = `Name,Phone,Plate Number
Serge Ntwali,784871704,RAJ746C
Benimana Christiane,784149020,RAJ748C
Nkurunziza Aloys,788268767,RAJ745C
Rutayisire Bosco,785040266,RAJ783C
Nyamuvugwa Jesus,788438122,RAJ835C`;
    const blob = new Blob([template], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = 'drivers_import_template.csv';
    link.click();
  };

  // Clear selection
  const clearSelection = () => {
    setSelectedDrivers([]);
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
          <Button variant="outline" size="sm" onClick={() => setIsUploadModalOpen(true)}>
            <Upload className="h-4 w-4 mr-2" />
            Import CSV
          </Button>
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

      {/* Bulk Actions Bar */}
      {selectedDrivers.length > 0 && (
        <div className="flex items-center justify-between rounded-lg bg-primary/5 border border-primary/20 p-3">
          <div className="flex items-center gap-3">
            <Checkbox
              checked={selectedDrivers.length === drivers.length && drivers.length > 0}
              onCheckedChange={toggleSelectAll}
            />
            <span className="text-sm font-medium">
              {selectedDrivers.length} driver{selectedDrivers.length > 1 ? 's' : ''} selected
            </span>
            <Button variant="ghost" size="sm" onClick={clearSelection}>
              Clear
            </Button>
          </div>
          <div className="flex items-center gap-2">
            <Button 
              variant="outline" 
              size="sm" 
              onClick={() => setIsBulkSuspendModalOpen(true)}
              className="text-yellow-600 hover:text-yellow-700 hover:bg-yellow-50"
            >
              <Ban className="h-4 w-4 mr-2" />
              Suspend Selected
            </Button>
            <Button 
              variant="outline" 
              size="sm" 
              onClick={() => handleBulkSuspend(false)}
              className="text-green-600 hover:text-green-700 hover:bg-green-50"
            >
              <CheckCircle className="h-4 w-4 mr-2" />
              Activate Selected
            </Button>
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
                <TableHead>Vehicle</TableHead>
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
                      <span className="text-sm font-medium">{getDriverVehicleLabel(driver)}</span>
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
            <div className="space-y-2">
              <Label htmlFor="vehicle_id">Assign Vehicle</Label>
              <Select value={formData.vehicle_id || "none"} onValueChange={(v) => setFormData({ ...formData, vehicle_id: v === "none" ? "" : v })}>
                <SelectTrigger>
                  <SelectValue placeholder="Select a vehicle (optional)" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="none">No Vehicle</SelectItem>
                  {vehicles.map((vehicle) => (
                    <SelectItem key={vehicle.vehicle_id} value={vehicle.vehicle_id}>
                      {vehicle.display_name} ({vehicle.plate_number})
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
            <div className="space-y-2">
              <Label htmlFor="edit_vehicle_id">Assign Vehicle</Label>
              <Select value={formData.vehicle_id || "none"} onValueChange={(v) => setFormData({ ...formData, vehicle_id: v === "none" ? "" : v })}>
                <SelectTrigger>
                  <SelectValue placeholder="Select a vehicle (optional)" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="none">No Vehicle</SelectItem>
                  {vehicles.map((vehicle) => (
                    <SelectItem key={vehicle.vehicle_id} value={vehicle.vehicle_id}>
                      {vehicle.display_name} ({vehicle.plate_number})
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

      {/* Bulk Delete Confirmation Modal */}
      <Dialog open={isBulkDeleteModalOpen} onOpenChange={setIsBulkDeleteModalOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="text-destructive">Delete {selectedDrivers.length} Driver(s)</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete {selectedDrivers.length} selected driver(s)? 
              This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsBulkDeleteModalOpen(false)} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleBulkDelete} disabled={isSubmitting}>
              {isSubmitting ? 'Deleting...' : `Delete ${selectedDrivers.length} Driver(s)`}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Bulk Suspend Confirmation Modal */}
      <Dialog open={isBulkSuspendModalOpen} onOpenChange={setIsBulkSuspendModalOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Suspend {selectedDrivers.length} Driver(s)</DialogTitle>
            <DialogDescription>
              Are you sure you want to suspend {selectedDrivers.length} selected driver(s)? 
              They will not be able to receive ride requests.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsBulkSuspendModalOpen(false)} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button 
              variant="destructive"
              onClick={() => handleBulkSuspend(true)} 
              disabled={isSubmitting}
            >
              {isSubmitting ? 'Suspending...' : `Suspend ${selectedDrivers.length} Driver(s)`}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* CSV Upload Modal */}
      <Dialog open={isUploadModalOpen} onOpenChange={(open) => {
        setIsUploadModalOpen(open);
        if (!open) {
          setUploadFile(null);
          setUploadPreview([]);
          setUploadError(null);
        }
      }}>
        <DialogContent className="sm:max-w-xl">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <FileSpreadsheet className="h-5 w-5" />
              Import Drivers from CSV
            </DialogTitle>
            <DialogDescription>
              Upload a CSV file with driver information. Download the template to see the required format.
            </DialogDescription>
          </DialogHeader>
          
          <div className="space-y-4 py-4">
            {/* Template Download */}
            <div className="flex items-center justify-between rounded-lg border border-dashed p-4">
              <div>
                <p className="text-sm font-medium">Need a template?</p>
                <p className="text-xs text-muted-foreground">Download our CSV template with the correct headers</p>
              </div>
              <Button variant="outline" size="sm" onClick={handleDownloadTemplate}>
                <Download className="h-4 w-4 mr-2" />
                Template
              </Button>
            </div>

            {/* File Upload Area */}
            <div className="rounded-lg border-2 border-dashed p-6 text-center">
              <input
                type="file"
                accept=".csv,.xlsx,.xls"
                onChange={handleFileSelect}
                className="hidden"
                id="csv-upload"
              />
              <label htmlFor="csv-upload" className="cursor-pointer">
                <Upload className="h-10 w-10 mx-auto text-muted-foreground mb-2" />
                <p className="text-sm font-medium">
                  {uploadFile ? uploadFile.name : 'Click to upload or drag and drop'}
                </p>
                <p className="text-xs text-muted-foreground mt-1">
                  CSV, XLS, or XLSX files up to 10MB
                </p>
              </label>
            </div>

            {/* Upload Error */}
            {uploadError && (
              <div className="flex items-center gap-2 rounded-lg bg-red-50 border border-red-200 p-3 text-sm text-red-800">
                <AlertCircle className="h-4 w-4 flex-shrink-0" />
                <span>{uploadError}</span>
              </div>
            )}

            {/* Preview Table */}
            {uploadPreview.length > 0 && (
              <div className="space-y-2">
                <p className="text-sm font-medium">Preview (first 5 rows):</p>
                <div className="rounded-lg border overflow-hidden">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        {Object.keys(uploadPreview[0]).slice(0, 5).map((key) => (
                          <TableHead key={key} className="text-xs capitalize">
                            {key}
                          </TableHead>
                        ))}
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {uploadPreview.map((row, i) => (
                        <TableRow key={i}>
                          {Object.values(row).slice(0, 5).map((value, j) => (
                            <TableCell key={j} className="text-xs py-2">
                              {value || '-'}
                            </TableCell>
                          ))}
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </div>
                <p className="text-xs text-muted-foreground">
                  Ready to import {uploadPreview.length > 5 ? `${uploadPreview.length}+` : uploadPreview.length} driver(s)
                </p>
              </div>
            )}
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setIsUploadModalOpen(false)} disabled={isUploading}>
              Cancel
            </Button>
            <Button 
              onClick={handleUploadCSV} 
              disabled={!uploadFile || isUploading}
            >
              {isUploading ? 'Importing...' : 'Import Drivers'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
