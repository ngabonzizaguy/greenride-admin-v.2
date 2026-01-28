'use client';

import { useState, useEffect, useCallback } from 'react';
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
  Filter,
  RefreshCw,
  AlertCircle,
  CheckCircle,
  Mail,
  Download,
  Upload,
  Edit,
  Trash2,
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
import { Skeleton } from '@/components/ui/skeleton';
import { Checkbox } from '@/components/ui/checkbox';
import { apiClient, ApiError } from '@/lib/api-client';
import type { User, PageResult, UserStatus } from '@/types';

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
    case 'suspended':
      return (
        <Badge className="bg-red-100 text-red-700 hover:bg-red-100">
          Suspended
        </Badge>
      );
    case 'banned':
      return (
        <Badge className="bg-red-100 text-red-700 hover:bg-red-100">
          Banned
        </Badge>
      );
    default:
      return <Badge variant="secondary">{status}</Badge>;
  }
};

// Empty user form
const emptyUserForm = {
  first_name: '',
  last_name: '',
  email: '',
  phone: '',
};

export default function UsersPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [selectedUsers, setSelectedUsers] = useState<string[]>([]);
  
  // Pagination
  const [page, setPage] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const limit = 10;
  
  // Stats
  const [stats, setStats] = useState({
    total: 0,
    active: 0,
    new_this_month: 0,
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
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [formData, setFormData] = useState(emptyUserForm);
  const [isSubmitting, setIsSubmitting] = useState(false);
  
  // CSV Upload states
  const [uploadFile, setUploadFile] = useState<File | null>(null);
  const [uploadPreview, setUploadPreview] = useState<Array<Record<string, string>>>([]);
  const [uploadError, setUploadError] = useState<string | null>(null);
  const [isUploading, setIsUploading] = useState(false);

  const fetchUsers = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await apiClient.getUsers({
        page,
        limit,
        keyword: search || undefined,
        status: statusFilter !== 'all' ? statusFilter as UserStatus : undefined,
      });
      
      const pageResult = response.data as PageResult<User>;
      setUsers(pageResult.records || []);
      setTotalCount(pageResult.count || 0);
      setTotalPages(pageResult.total || 0);
      
      if (pageResult.attach) {
        setStats({
          total: (pageResult.attach.total_count as number) || pageResult.count || 0,
          active: (pageResult.attach.active_count as number) || 0,
          new_this_month: (pageResult.attach.new_this_month as number) || 0,
          suspended: (pageResult.attach.suspended_count as number) || 0,
        });
      }
    } catch (err) {
      console.error('Failed to fetch users:', err);
      setError('Failed to load users. Please try again.');
    } finally {
      setIsLoading(false);
    }
  }, [page, limit, search, statusFilter]);

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  useEffect(() => {
    const timer = setTimeout(() => {
      if (page !== 1) {
        setPage(1);
      } else {
        fetchUsers();
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

  const getDisplayName = (user: User) => {
    return user.display_name || 
      (user.first_name && user.last_name 
        ? `${user.first_name} ${user.last_name}` 
        : user.username || 'Unknown User');
  };

  const getInitials = (user: User) => {
    const name = getDisplayName(user);
    return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2);
  };

  // Handle Add User
  const handleAddUser = async () => {
    if (!formData.first_name || !formData.phone) {
      setError('First name and phone are required');
      return;
    }

    setIsSubmitting(true);
    setError(null);

    try {
      await apiClient.createUser({
        user_type: 'passenger',
        first_name: formData.first_name,
        last_name: formData.last_name,
        email: formData.email,
        phone: formData.phone,
        status: 'active',
      });

      setSuccessMessage('User added successfully!');
      setIsAddModalOpen(false);
      setFormData(emptyUserForm);
      fetchUsers();
    } catch (err) {
      console.error('Failed to add user:', err);
      setError('Failed to add user. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  // Handle Edit User
  const handleEditUser = async () => {
    if (!selectedUser) return;

    setIsSubmitting(true);
    setError(null);

    try {
      await apiClient.updateUser(selectedUser.user_id, {
        first_name: formData.first_name,
        last_name: formData.last_name,
        email: formData.email,
        phone: formData.phone,
      });

      setSuccessMessage('User updated successfully!');
      setIsEditModalOpen(false);
      setSelectedUser(null);
      setFormData(emptyUserForm);
      fetchUsers();
    } catch (err) {
      console.error('Failed to update user:', err);
      setError('Failed to update user. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  // Handle Delete User
  const handleDeleteUser = async () => {
    if (!selectedUser) return;

    setIsSubmitting(true);
    setError(null);

    try {
      await apiClient.deleteUser(selectedUser.user_id, 'Deleted by admin');
      setSuccessMessage('User deleted successfully!');
      setIsDeleteModalOpen(false);
      setSelectedUser(null);
      fetchUsers();
    } catch (err) {
      console.error('Failed to delete user:', err);
      const message =
        err instanceof ApiError
          ? (err.serverMessage || err.message)
          : (err instanceof Error ? err.message : 'Failed to delete user.');
      setError(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  // Handle Suspend/Activate
  const handleToggleSuspend = async () => {
    if (!selectedUser) return;

    setIsSubmitting(true);
    setError(null);

    const newStatus = selectedUser.status === 'suspended' ? 'active' : 'suspended';

    try {
      await apiClient.updateUserStatus(selectedUser.user_id, newStatus);
      setSuccessMessage(`User ${newStatus === 'suspended' ? 'suspended' : 'activated'} successfully!`);
      setIsSuspendModalOpen(false);
      setSelectedUser(null);
      fetchUsers();
    } catch (err) {
      console.error('Failed to update user status:', err);
      setError('Failed to update user status. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  // Open Edit Modal
  const openEditModal = (user: User) => {
    setSelectedUser(user);
    setFormData({
      first_name: user.first_name || '',
      last_name: user.last_name || '',
      email: user.email || '',
      phone: user.phone || '',
    });
    setIsEditModalOpen(true);
  };

  // Export to CSV
  const handleExportCSV = () => {
    const headers = ['Name', 'Email', 'Phone', 'Status', 'Total Rides', 'Verified', 'Joined'];
    const csvContent = [
      headers.join(','),
      ...users.map(user => [
        `"${getDisplayName(user)}"`,
        user.email || '',
        user.phone || '',
        user.status,
        user.total_rides || 0,
        user.is_phone_verified ? 'Yes' : 'No',
        user.created_at ? new Date(user.created_at).toLocaleDateString() : '',
      ].join(','))
    ].join('\n');

    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `users_export_${new Date().toISOString().split('T')[0]}.csv`;
    link.click();
    setSuccessMessage('Users exported to CSV!');
  };

  // Selection handlers
  const toggleSelectAll = () => {
    if (selectedUsers.length === users.length) {
      setSelectedUsers([]);
    } else {
      setSelectedUsers(users.map((u) => u.user_id));
    }
  };

  const toggleSelectUser = (id: string) => {
    setSelectedUsers((prev) =>
      prev.includes(id) ? prev.filter((i) => i !== id) : [...prev, id]
    );
  };

  const clearSelection = () => {
    setSelectedUsers([]);
  };

  // Bulk Delete Users
  const handleBulkDelete = async () => {
    if (selectedUsers.length === 0) return;

    setIsSubmitting(true);
    setError(null);

    try {
      const results = await Promise.allSettled(
        selectedUsers.map(userId => apiClient.deleteUser(userId, 'Bulk deleted by admin'))
      );

      const successCount = results.filter(r => r.status === 'fulfilled').length;
      const failCount = results.filter(r => r.status === 'rejected').length;

      if (failCount > 0) {
        setError(`${failCount} user(s) failed to delete. ${successCount} deleted successfully.`);
      } else {
        setSuccessMessage(`${successCount} user(s) deleted successfully!`);
      }

      setIsBulkDeleteModalOpen(false);
      setSelectedUsers([]);
      fetchUsers();
    } catch (err) {
      console.error('Failed to bulk delete:', err);
      const message =
        err instanceof ApiError
          ? (err.serverMessage || err.message)
          : (err instanceof Error ? err.message : 'Failed to delete users.');
      setError(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  // Bulk Suspend/Activate Users
  const handleBulkSuspend = async (suspend: boolean) => {
    if (selectedUsers.length === 0) return;

    setIsSubmitting(true);
    setError(null);

    const newStatus = suspend ? 'suspended' : 'active';

    try {
      const results = await Promise.allSettled(
        selectedUsers.map(userId => apiClient.updateUserStatus(userId, newStatus))
      );

      const successCount = results.filter(r => r.status === 'fulfilled').length;
      const failCount = results.filter(r => r.status === 'rejected').length;

      if (failCount > 0) {
        setError(`${failCount} user(s) failed to ${suspend ? 'suspend' : 'activate'}. ${successCount} updated successfully.`);
      } else {
        setSuccessMessage(`${successCount} user(s) ${suspend ? 'suspended' : 'activated'} successfully!`);
      }

      setIsBulkSuspendModalOpen(false);
      setSelectedUsers([]);
      fetchUsers();
    } catch (err) {
      console.error('Failed to bulk update status:', err);
      setError(`Failed to ${suspend ? 'suspend' : 'activate'} users. Please try again.`);
    } finally {
      setIsSubmitting(false);
    }
  };

  // Parse CSV file
  const parseCSV = (text: string): Array<Record<string, string>> => {
    const lines = text.split('\n').filter(line => line.trim());
    if (lines.length < 2) return [];

    const headers = lines[0].split(',').map(h => h.trim().toLowerCase().replace(/["']/g, ''));
    const rows: Array<Record<string, string>> = [];

    for (let i = 1; i < lines.length; i++) {
      const values = lines[i].split(',').map(v => v.trim().replace(/["']/g, ''));
      const row: Record<string, string> = {};
      headers.forEach((header, index) => {
        row[header] = values[index] || '';
      });
      rows.push(row);
    }

    return rows;
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

    const validTypes = ['text/csv', 'application/vnd.ms-excel', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'];
    if (!validTypes.includes(file.type) && !file.name.endsWith('.csv')) {
      setUploadError('Please upload a CSV file');
      return;
    }

    setUploadFile(file);

    const reader = new FileReader();
    reader.onload = (event) => {
      try {
        const text = event.target?.result as string;
        const parsed = parseCSV(text);
        if (parsed.length === 0) {
          setUploadError('No valid data found in file');
          return;
        }
        setUploadPreview(parsed.slice(0, 5));
      } catch {
        setUploadError('Failed to parse file');
      }
    };
    reader.readAsText(file);
  };

  // Handle CSV Upload
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

          const results = await Promise.allSettled(
            rows.map(row => apiClient.createUser({
              user_type: 'passenger',
              first_name: row['first name'] || row['first_name'] || row['firstname'] || '',
              last_name: row['last name'] || row['last_name'] || row['lastname'] || '',
              email: row['email'] || '',
              phone: row['phone'] || row['phone number'] || row['phone_number'] || '',
              status: 'active',
            }))
          );

          const successCount = results.filter(r => r.status === 'fulfilled').length;
          const failCount = results.filter(r => r.status === 'rejected').length;

          if (failCount > 0) {
            setError(`${failCount} user(s) failed to import. ${successCount} imported successfully.`);
          } else {
            setSuccessMessage(`${successCount} user(s) imported successfully!`);
          }

          setIsUploadModalOpen(false);
          setUploadFile(null);
          setUploadPreview([]);
          fetchUsers();
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

  // Download CSV template
  const handleDownloadTemplate = () => {
    const template = 'First Name,Last Name,Email,Phone\nJohn,Doe,john@example.com,+250788123456\nJane,Smith,jane@example.com,+250788234567';
    const blob = new Blob([template], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = 'users_import_template.csv';
    link.click();
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">User Management</h1>
          <p className="text-muted-foreground">
            View and manage passenger accounts
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" size="sm" onClick={() => setIsUploadModalOpen(true)}>
            <Upload className="h-4 w-4 mr-2" />
            Import CSV
          </Button>
          <Button variant="outline" size="sm" onClick={handleExportCSV} disabled={users.length === 0}>
            <Download className="h-4 w-4 mr-2" />
            Export CSV
          </Button>
          <Button variant="outline" size="sm" onClick={fetchUsers} disabled={isLoading}>
            <RefreshCw className={`h-4 w-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
          <Button className="gap-2" onClick={() => { setFormData(emptyUserForm); setIsAddModalOpen(true); }}>
            <UserPlus className="h-4 w-4" />
            Add User
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
      {selectedUsers.length > 0 && (
        <div className="flex items-center justify-between rounded-lg bg-primary/5 border border-primary/20 p-3">
          <div className="flex items-center gap-3">
            <Checkbox
              checked={selectedUsers.length === users.length && users.length > 0}
              onCheckedChange={toggleSelectAll}
            />
            <span className="text-sm font-medium">
              {selectedUsers.length} user{selectedUsers.length > 1 ? 's' : ''} selected
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
                <p className="text-sm text-muted-foreground">Total Users</p>
                {isLoading ? (
                  <Skeleton className="h-8 w-16 mt-1" />
                ) : (
                  <p className="text-2xl font-bold">{stats.total || totalCount}</p>
                )}
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
                <p className="text-sm text-muted-foreground">Active Users</p>
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
                <p className="text-sm text-muted-foreground">New This Month</p>
                {isLoading ? (
                  <Skeleton className="h-8 w-16 mt-1" />
                ) : (
                  <p className="text-2xl font-bold text-blue-600">{stats.new_this_month}</p>
                )}
              </div>
              <UserPlus className="h-5 w-5 text-blue-500" />
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
                <SelectItem value="inactive">Inactive</SelectItem>
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
                <TableHead className="w-12">
                  <Checkbox
                    checked={selectedUsers.length === users.length && users.length > 0}
                    onCheckedChange={toggleSelectAll}
                  />
                </TableHead>
                <TableHead>User</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Phone</TableHead>
                <TableHead className="text-right">Total Rides</TableHead>
                <TableHead>Verified</TableHead>
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
                    <TableCell><Skeleton className="h-6 w-16" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-12" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-16" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                    <TableCell><Skeleton className="h-8 w-8" /></TableCell>
                  </TableRow>
                ))
              ) : users.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={8} className="h-32 text-center">
                    <div className="flex flex-col items-center gap-2 text-muted-foreground">
                      <Users className="h-8 w-8" />
                      <p>No users found</p>
                      {search && <p className="text-sm">Try adjusting your search</p>}
                    </div>
                  </TableCell>
                </TableRow>
              ) : (
                users.map((user) => (
                  <TableRow key={user.user_id} className="group">
                    <TableCell>
                      <Checkbox
                        checked={selectedUsers.includes(user.user_id)}
                        onCheckedChange={() => toggleSelectUser(user.user_id)}
                      />
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-3">
                        <Avatar className="h-10 w-10">
                          <AvatarImage src={user.avatar || undefined} />
                          <AvatarFallback className="bg-primary/10 text-primary">
                            {getInitials(user)}
                          </AvatarFallback>
                        </Avatar>
                        <div>
                          <Link
                            href={`/users/${user.user_id}`}
                            className="font-medium hover:text-primary hover:underline"
                          >
                            {getDisplayName(user)}
                          </Link>
                          <p className="text-sm text-muted-foreground">{user.email || 'No email'}</p>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell>{getStatusBadge(user.status)}</TableCell>
                    <TableCell className="text-muted-foreground">
                      {user.phone || 'No phone'}
                    </TableCell>
                    <TableCell className="text-right font-medium">
                      {(user.total_rides || 0).toLocaleString()}
                    </TableCell>
                    <TableCell>
                      <div className="flex gap-1">
                        {user.is_email_verified && (
                          <Badge variant="outline" className="text-xs">
                            <Mail className="h-3 w-3 mr-1" />
                            Email
                          </Badge>
                        )}
                        {user.is_phone_verified && (
                          <Badge variant="outline" className="text-xs">
                            <Phone className="h-3 w-3 mr-1" />
                            Phone
                          </Badge>
                        )}
                        {!user.is_email_verified && !user.is_phone_verified && (
                          <span className="text-muted-foreground text-sm">Not verified</span>
                        )}
                      </div>
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {user.created_at 
                        ? new Date(user.created_at).toLocaleDateString('en-US', {
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
                            <Link href={`/users/${user.user_id}`} className="flex items-center gap-2">
                              <Eye className="h-4 w-4" />
                              View Details
                            </Link>
                          </DropdownMenuItem>
                          <DropdownMenuItem className="gap-2" onClick={() => openEditModal(user)}>
                            <Edit className="h-4 w-4" />
                            Edit
                          </DropdownMenuItem>
                          {user.phone && (
                            <DropdownMenuItem className="gap-2" asChild>
                              <a href={`tel:${user.phone}`}>
                                <Phone className="h-4 w-4" />
                                Call User
                              </a>
                            </DropdownMenuItem>
                          )}
                          <DropdownMenuSeparator />
                          <DropdownMenuItem 
                            className={user.status === 'suspended' ? 'gap-2 text-green-600' : 'gap-2 text-yellow-600'}
                            onClick={() => { setSelectedUser(user); setIsSuspendModalOpen(true); }}
                          >
                            {user.status === 'suspended' ? (
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
                            onClick={() => { setSelectedUser(user); setIsDeleteModalOpen(true); }}
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
                <>Showing {((page - 1) * limit) + 1}-{Math.min(page * limit, totalCount)} of {totalCount} users</>
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

      {/* Add User Modal */}
      <Dialog open={isAddModalOpen} onOpenChange={setIsAddModalOpen}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Add New User</DialogTitle>
            <DialogDescription>
              Enter the user&apos;s information to create their account.
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
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsAddModalOpen(false)} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button onClick={handleAddUser} disabled={isSubmitting}>
              {isSubmitting ? 'Adding...' : 'Add User'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Edit User Modal */}
      <Dialog open={isEditModalOpen} onOpenChange={setIsEditModalOpen}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Edit User</DialogTitle>
            <DialogDescription>
              Update user information.
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
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsEditModalOpen(false)} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button onClick={handleEditUser} disabled={isSubmitting}>
              {isSubmitting ? 'Saving...' : 'Save Changes'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Modal */}
      <Dialog open={isDeleteModalOpen} onOpenChange={setIsDeleteModalOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="text-destructive">Delete User</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete <strong>{selectedUser && getDisplayName(selectedUser)}</strong>? 
              This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsDeleteModalOpen(false)} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDeleteUser} disabled={isSubmitting}>
              {isSubmitting ? 'Deleting...' : 'Delete User'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Suspend Confirmation Modal */}
      <Dialog open={isSuspendModalOpen} onOpenChange={setIsSuspendModalOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>
              {selectedUser?.status === 'suspended' ? 'Activate User' : 'Suspend User'}
            </DialogTitle>
            <DialogDescription>
              {selectedUser?.status === 'suspended' 
                ? `Are you sure you want to activate ${selectedUser && getDisplayName(selectedUser)}?`
                : `Are you sure you want to suspend ${selectedUser && getDisplayName(selectedUser)}? They will not be able to book rides.`
              }
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsSuspendModalOpen(false)} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button 
              variant={selectedUser?.status === 'suspended' ? 'default' : 'destructive'}
              onClick={handleToggleSuspend} 
              disabled={isSubmitting}
            >
              {isSubmitting ? 'Processing...' : selectedUser?.status === 'suspended' ? 'Activate' : 'Suspend'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Bulk Delete Confirmation Modal */}
      <Dialog open={isBulkDeleteModalOpen} onOpenChange={setIsBulkDeleteModalOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="text-destructive">Delete {selectedUsers.length} User(s)</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete {selectedUsers.length} selected user(s)? 
              This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsBulkDeleteModalOpen(false)} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleBulkDelete} disabled={isSubmitting}>
              {isSubmitting ? 'Deleting...' : `Delete ${selectedUsers.length} User(s)`}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Bulk Suspend Confirmation Modal */}
      <Dialog open={isBulkSuspendModalOpen} onOpenChange={setIsBulkSuspendModalOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Suspend {selectedUsers.length} User(s)</DialogTitle>
            <DialogDescription>
              Are you sure you want to suspend {selectedUsers.length} selected user(s)? 
              They will not be able to book rides.
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
              {isSubmitting ? 'Suspending...' : `Suspend ${selectedUsers.length} User(s)`}
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
              Import Users from CSV
            </DialogTitle>
            <DialogDescription>
              Upload a CSV file with user information. Download the template to see the required format.
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
                        {Object.keys(uploadPreview[0]).slice(0, 4).map((key) => (
                          <TableHead key={key} className="text-xs capitalize">
                            {key}
                          </TableHead>
                        ))}
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {uploadPreview.map((row, i) => (
                        <TableRow key={i}>
                          {Object.values(row).slice(0, 4).map((value, j) => (
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
                  Ready to import {uploadPreview.length > 5 ? `${uploadPreview.length}+` : uploadPreview.length} user(s)
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
              {isUploading ? 'Importing...' : 'Import Users'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
