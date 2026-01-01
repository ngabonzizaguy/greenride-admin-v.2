'use client';

import { useState, useEffect } from 'react';
import { 
  Plus, 
  Search, 
  MoreHorizontal,
  Tag,
  Edit,
  Trash2,
  Copy,
  ToggleLeft,
  Calendar,
  Users,
  TrendingUp,
  AlertCircle,
  X,
  Download
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
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
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { toast } from 'sonner';

// Promotion type definition
interface Promotion {
  id: string;
  code: string;
  type: 'percentage' | 'fixed' | 'free_ride';
  value: number;
  usageLimit: number | null;
  usedCount: number;
  validFrom: string;
  validUntil: string;
  status: 'active' | 'expired' | 'disabled';
  minOrder: number;
  maxDiscount: number;
}

// Initial mock promotions data
const initialPromotions: Promotion[] = [
  {
    id: '1',
    code: 'WELCOME50',
    type: 'percentage',
    value: 50,
    usageLimit: 1000,
    usedCount: 234,
    validFrom: '2024-12-01',
    validUntil: '2025-01-31',
    status: 'active',
    minOrder: 5000,
    maxDiscount: 3000,
  },
  {
    id: '2',
    code: 'NEWYEAR25',
    type: 'percentage',
    value: 25,
    usageLimit: 500,
    usedCount: 123,
    validFrom: '2024-12-25',
    validUntil: '2025-01-05',
    status: 'active',
    minOrder: 3000,
    maxDiscount: 2000,
  },
  {
    id: '3',
    code: 'FREERIDE',
    type: 'free_ride',
    value: 100,
    usageLimit: 100,
    usedCount: 45,
    validFrom: '2024-12-15',
    validUntil: '2024-12-31',
    status: 'active',
    minOrder: 0,
    maxDiscount: 10000,
  },
  {
    id: '4',
    code: 'FLAT2000',
    type: 'fixed',
    value: 2000,
    usageLimit: 200,
    usedCount: 200,
    validFrom: '2024-11-01',
    validUntil: '2024-11-30',
    status: 'expired',
    minOrder: 5000,
    maxDiscount: 2000,
  },
  {
    id: '5',
    code: 'VIP10',
    type: 'percentage',
    value: 10,
    usageLimit: null,
    usedCount: 567,
    validFrom: '2024-01-01',
    validUntil: '2025-12-31',
    status: 'disabled',
    minOrder: 0,
    maxDiscount: 5000,
  },
];

// Empty form state
const emptyFormData = {
  code: '',
  type: 'percentage' as Promotion['type'],
  value: 0,
  maxDiscount: 0,
  minOrder: 0,
  validFrom: '',
  validUntil: '',
  usageLimit: null as number | null,
};

const getStatusBadge = (status: string) => {
  switch (status) {
    case 'active':
      return <Badge className="bg-green-100 text-green-700 hover:bg-green-100">Active</Badge>;
    case 'expired':
      return <Badge className="bg-gray-100 text-gray-700 hover:bg-gray-100">Expired</Badge>;
    case 'disabled':
      return <Badge className="bg-red-100 text-red-700 hover:bg-red-100">Disabled</Badge>;
    default:
      return <Badge variant="secondary">{status}</Badge>;
  }
};

const getTypeLabel = (type: string, value: number) => {
  switch (type) {
    case 'percentage':
      return `${value}% off`;
    case 'fixed':
      return `RWF ${value.toLocaleString()} off`;
    case 'free_ride':
      return 'Free Ride';
    default:
      return type;
  }
};

export default function PromotionsPage() {
  const [promotions, setPromotions] = useState<Promotion[]>(initialPromotions);
  const [search, setSearch] = useState('');
  
  // Modal states
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [isDeleteOpen, setIsDeleteOpen] = useState(false);
  
  // Form state
  const [formData, setFormData] = useState(emptyFormData);
  const [editingPromotion, setEditingPromotion] = useState<Promotion | null>(null);
  const [deletingPromotion, setDeletingPromotion] = useState<Promotion | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [formErrors, setFormErrors] = useState<Record<string, string>>({});

  // Calculate stats from promotions
  const stats = {
    activePromotions: promotions.filter(p => p.status === 'active').length,
    totalRedemptions: promotions.reduce((acc, p) => acc + p.usedCount, 0),
    revenueImpact: promotions.reduce((acc, p) => acc - (p.usedCount * (p.type === 'fixed' ? p.value : p.maxDiscount / 2)), 0),
    avgDiscount: Math.round(promotions.reduce((acc, p) => acc + (p.type === 'fixed' ? p.value : p.maxDiscount), 0) / promotions.length),
  };

  const filteredPromotions = promotions.filter((promo) =>
    promo.code.toLowerCase().includes(search.toLowerCase())
  );

  // Validate form
  const validateForm = (): boolean => {
    const errors: Record<string, string> = {};
    
    if (!formData.code.trim()) {
      errors.code = 'Promo code is required';
    } else if (formData.code.length < 3) {
      errors.code = 'Promo code must be at least 3 characters';
    } else if (!/^[A-Z0-9]+$/.test(formData.code.toUpperCase())) {
      errors.code = 'Promo code can only contain letters and numbers';
    }
    
    if (formData.type !== 'free_ride' && formData.value <= 0) {
      errors.value = 'Discount value must be greater than 0';
    }
    
    if (formData.type === 'percentage' && formData.value > 100) {
      errors.value = 'Percentage cannot exceed 100%';
    }
    
    if (!formData.validFrom) {
      errors.validFrom = 'Start date is required';
    }
    
    if (!formData.validUntil) {
      errors.validUntil = 'End date is required';
    }
    
    if (formData.validFrom && formData.validUntil && formData.validFrom > formData.validUntil) {
      errors.validUntil = 'End date must be after start date';
    }
    
    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  // Handle Create
  const handleCreate = () => {
    if (!validateForm()) return;
    
    setIsSubmitting(true);
    
    // Simulate API delay
    setTimeout(() => {
      const newPromotion: Promotion = {
        id: String(Date.now()),
        code: formData.code.toUpperCase(),
        type: formData.type,
        value: formData.type === 'free_ride' ? 100 : formData.value,
        usageLimit: formData.usageLimit,
        usedCount: 0,
        validFrom: formData.validFrom,
        validUntil: formData.validUntil,
        status: 'active',
        minOrder: formData.minOrder,
        maxDiscount: formData.maxDiscount,
      };
      
      setPromotions(prev => [newPromotion, ...prev]);
      setIsCreateOpen(false);
      setFormData(emptyFormData);
      setFormErrors({});
      setIsSubmitting(false);
      toast.success(`Promotion "${newPromotion.code}" created successfully!`);
    }, 500);
  };

  // Handle Edit
  const openEditModal = (promo: Promotion) => {
    setEditingPromotion(promo);
    setFormData({
      code: promo.code,
      type: promo.type,
      value: promo.value,
      maxDiscount: promo.maxDiscount,
      minOrder: promo.minOrder,
      validFrom: promo.validFrom,
      validUntil: promo.validUntil,
      usageLimit: promo.usageLimit,
    });
    setFormErrors({});
    setIsEditOpen(true);
  };

  const handleEdit = () => {
    if (!validateForm() || !editingPromotion) return;
    
    setIsSubmitting(true);
    
    setTimeout(() => {
      setPromotions(prev => prev.map(p => {
        if (p.id === editingPromotion.id) {
          return {
            ...p,
            code: formData.code.toUpperCase(),
            type: formData.type,
            value: formData.type === 'free_ride' ? 100 : formData.value,
            maxDiscount: formData.maxDiscount,
            minOrder: formData.minOrder,
            validFrom: formData.validFrom,
            validUntil: formData.validUntil,
            usageLimit: formData.usageLimit,
          };
        }
        return p;
      }));
      
      setIsEditOpen(false);
      setEditingPromotion(null);
      setFormData(emptyFormData);
      setFormErrors({});
      setIsSubmitting(false);
      toast.success(`Promotion "${formData.code.toUpperCase()}" updated successfully!`);
    }, 500);
  };

  // Handle Delete
  const openDeleteModal = (promo: Promotion) => {
    setDeletingPromotion(promo);
    setIsDeleteOpen(true);
  };

  const handleDelete = () => {
    if (!deletingPromotion) return;
    
    setPromotions(prev => prev.filter(p => p.id !== deletingPromotion.id));
    setIsDeleteOpen(false);
    toast.success(`Promotion "${deletingPromotion.code}" deleted successfully!`);
    setDeletingPromotion(null);
  };

  // Handle Toggle Status
  const handleToggleStatus = (promo: Promotion) => {
    const newStatus = promo.status === 'active' ? 'disabled' : 'active';
    setPromotions(prev => prev.map(p => {
      if (p.id === promo.id) {
        return { ...p, status: newStatus };
      }
      return p;
    }));
    toast.success(`Promotion "${promo.code}" ${newStatus === 'active' ? 'enabled' : 'disabled'}!`);
  };

  // Handle Duplicate
  const handleDuplicate = (promo: Promotion) => {
    const duplicatedPromo: Promotion = {
      ...promo,
      id: String(Date.now()),
      code: `${promo.code}_COPY`,
      usedCount: 0,
      status: 'disabled',
    };
    setPromotions(prev => [duplicatedPromo, ...prev]);
    toast.success(`Promotion duplicated as "${duplicatedPromo.code}"!`);
  };

  // Handle Copy Code
  const handleCopyCode = (code: string) => {
    navigator.clipboard.writeText(code);
    toast.success(`Copied "${code}" to clipboard!`);
  };

  // Handle CSV Export
  const handleExportCSV = () => {
    const headers = ['Code', 'Type', 'Value', 'Usage Limit', 'Used Count', 'Valid From', 'Valid Until', 'Status', 'Min Order', 'Max Discount'];
    const rows = filteredPromotions.map(p => [
      p.code,
      p.type,
      p.value,
      p.usageLimit ?? 'Unlimited',
      p.usedCount,
      p.validFrom,
      p.validUntil,
      p.status,
      p.minOrder,
      p.maxDiscount
    ]);
    
    const csvContent = [
      headers.join(','),
      ...rows.map(row => row.join(','))
    ].join('\n');
    
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `promotions_${new Date().toISOString().split('T')[0]}.csv`;
    link.click();
    toast.success('Promotions exported to CSV!');
  };

  // Form input handler
  const updateFormField = (field: string, value: string | number | null) => {
    setFormData(prev => ({ ...prev, [field]: value }));
    if (formErrors[field]) {
      setFormErrors(prev => {
        const newErrors = { ...prev };
        delete newErrors[field];
        return newErrors;
      });
    }
  };

  // Render form fields (shared between Create and Edit dialogs)
  const renderFormFields = () => (
    <div className="space-y-4 py-4">
      <div className="space-y-2">
        <Label htmlFor="code">Promo Code</Label>
        <Input 
          id="code" 
          placeholder="e.g., SUMMER25" 
          value={formData.code}
          onChange={(e) => updateFormField('code', e.target.value.toUpperCase())}
          className={formErrors.code ? 'border-red-500' : ''}
        />
        {formErrors.code && <p className="text-xs text-red-500">{formErrors.code}</p>}
      </div>
      <div className="space-y-2">
        <Label>Discount Type</Label>
        <Select 
          value={formData.type}
          onValueChange={(value) => updateFormField('type', value as Promotion['type'])}
        >
          <SelectTrigger>
            <SelectValue placeholder="Select type" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="percentage">Percentage</SelectItem>
            <SelectItem value="fixed">Fixed Amount</SelectItem>
            <SelectItem value="free_ride">Free Ride</SelectItem>
          </SelectContent>
        </Select>
      </div>
      {formData.type !== 'free_ride' && (
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label htmlFor="value">
              {formData.type === 'percentage' ? 'Percentage (%)' : 'Amount (RWF)'}
            </Label>
            <Input 
              id="value" 
              type="number" 
              placeholder={formData.type === 'percentage' ? '25' : '2000'} 
              value={formData.value || ''}
              onChange={(e) => updateFormField('value', parseInt(e.target.value) || 0)}
              className={formErrors.value ? 'border-red-500' : ''}
            />
            {formErrors.value && <p className="text-xs text-red-500">{formErrors.value}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="max">Max Discount (RWF)</Label>
            <Input 
              id="max" 
              type="number" 
              placeholder="5000" 
              value={formData.maxDiscount || ''}
              onChange={(e) => updateFormField('maxDiscount', parseInt(e.target.value) || 0)}
            />
          </div>
        </div>
      )}
      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="validFrom">Valid From</Label>
          <Input 
            id="validFrom" 
            type="date" 
            value={formData.validFrom}
            onChange={(e) => updateFormField('validFrom', e.target.value)}
            className={formErrors.validFrom ? 'border-red-500' : ''}
          />
          {formErrors.validFrom && <p className="text-xs text-red-500">{formErrors.validFrom}</p>}
        </div>
        <div className="space-y-2">
          <Label htmlFor="validUntil">Valid Until</Label>
          <Input 
            id="validUntil" 
            type="date" 
            value={formData.validUntil}
            onChange={(e) => updateFormField('validUntil', e.target.value)}
            className={formErrors.validUntil ? 'border-red-500' : ''}
          />
          {formErrors.validUntil && <p className="text-xs text-red-500">{formErrors.validUntil}</p>}
        </div>
      </div>
      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="minOrder">Min Order (RWF)</Label>
          <Input 
            id="minOrder" 
            type="number" 
            placeholder="0" 
            value={formData.minOrder || ''}
            onChange={(e) => updateFormField('minOrder', parseInt(e.target.value) || 0)}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="limit">Usage Limit (optional)</Label>
          <Input 
            id="limit" 
            type="number" 
            placeholder="Unlimited" 
            value={formData.usageLimit ?? ''}
            onChange={(e) => updateFormField('usageLimit', e.target.value ? parseInt(e.target.value) : null)}
          />
        </div>
      </div>
    </div>
  );

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Promotions & Discounts</h1>
          <p className="text-muted-foreground">
            Create and manage promotional offers
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" size="sm" onClick={handleExportCSV}>
            <Download className="h-4 w-4 mr-2" />
            Export CSV
          </Button>
          <Button className="gap-2" onClick={() => {
            setFormData(emptyFormData);
            setFormErrors({});
            setIsCreateOpen(true);
          }}>
            <Plus className="h-4 w-4" />
            Create Promotion
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Active Promotions</p>
                <p className="text-2xl font-bold">{stats.activePromotions}</p>
              </div>
              <Tag className="h-5 w-5 text-primary" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Total Redemptions</p>
                <p className="text-2xl font-bold">{stats.totalRedemptions.toLocaleString()}</p>
              </div>
              <Users className="h-5 w-5 text-blue-500" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Revenue Impact</p>
                <p className="text-2xl font-bold text-red-600">
                  RWF {Math.abs(stats.revenueImpact).toLocaleString()}
                </p>
              </div>
              <TrendingUp className="h-5 w-5 text-red-500 rotate-180" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Avg. Discount</p>
                <p className="text-2xl font-bold">RWF {stats.avgDiscount.toLocaleString()}</p>
              </div>
              <Tag className="h-5 w-5 text-yellow-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Search */}
      <Card>
        <CardContent className="p-4">
          <div className="relative max-w-md">
            <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              placeholder="Search by promo code..."
              className="pl-10"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
          </div>
        </CardContent>
      </Card>

      {/* Promotions Table */}
      <Card>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Code</TableHead>
                <TableHead>Discount</TableHead>
                <TableHead>Usage</TableHead>
                <TableHead>Valid Period</TableHead>
                <TableHead>Min Order</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="w-12"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredPromotions.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} className="text-center py-8 text-muted-foreground">
                    {search ? 'No promotions match your search.' : 'No promotions yet. Create your first one!'}
                  </TableCell>
                </TableRow>
              ) : (
                filteredPromotions.map((promo) => (
                  <TableRow key={promo.id}>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <code className="font-mono font-bold text-primary">{promo.code}</code>
                        <Button 
                          variant="ghost" 
                          size="icon" 
                          className="h-6 w-6"
                          onClick={() => handleCopyCode(promo.code)}
                        >
                          <Copy className="h-3 w-3" />
                        </Button>
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge variant="outline">{getTypeLabel(promo.type, promo.value)}</Badge>
                    </TableCell>
                    <TableCell>
                      <div className="flex flex-col">
                        <span className="font-medium">{promo.usedCount.toLocaleString()}</span>
                        <span className="text-xs text-muted-foreground">
                          of {promo.usageLimit ? promo.usageLimit.toLocaleString() : '∞'}
                        </span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-1 text-sm">
                        <Calendar className="h-3 w-3 text-muted-foreground" />
                        <span>
                          {new Date(promo.validFrom).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}
                          {' - '}
                          {new Date(promo.validUntil).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}
                        </span>
                      </div>
                    </TableCell>
                    <TableCell>
                      {promo.minOrder > 0 ? `RWF ${promo.minOrder.toLocaleString()}` : '-'}
                    </TableCell>
                    <TableCell>{getStatusBadge(promo.status)}</TableCell>
                    <TableCell>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="icon" className="h-8 w-8">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem className="gap-2" onClick={() => openEditModal(promo)}>
                            <Edit className="h-4 w-4" />
                            Edit
                          </DropdownMenuItem>
                          <DropdownMenuItem className="gap-2" onClick={() => handleDuplicate(promo)}>
                            <Copy className="h-4 w-4" />
                            Duplicate
                          </DropdownMenuItem>
                          <DropdownMenuItem 
                            className="gap-2" 
                            onClick={() => handleToggleStatus(promo)}
                            disabled={promo.status === 'expired'}
                          >
                            <ToggleLeft className="h-4 w-4" />
                            {promo.status === 'active' ? 'Disable' : 'Enable'}
                          </DropdownMenuItem>
                          <DropdownMenuSeparator />
                          <DropdownMenuItem 
                            className="gap-2 text-destructive"
                            onClick={() => openDeleteModal(promo)}
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
        </CardContent>
      </Card>

      {/* Create Promotion Dialog */}
      <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Create New Promotion</DialogTitle>
            <DialogDescription>
              Set up a new promotional discount code
            </DialogDescription>
          </DialogHeader>
          {renderFormFields()}
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsCreateOpen(false)} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button onClick={handleCreate} disabled={isSubmitting}>
              {isSubmitting ? 'Creating...' : 'Create Promotion'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Edit Promotion Dialog */}
      <Dialog open={isEditOpen} onOpenChange={setIsEditOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Edit Promotion</DialogTitle>
            <DialogDescription>
              Update the promotion details
            </DialogDescription>
          </DialogHeader>
          {renderFormFields()}
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsEditOpen(false)} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button onClick={handleEdit} disabled={isSubmitting}>
              {isSubmitting ? 'Saving...' : 'Save Changes'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={isDeleteOpen} onOpenChange={setIsDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Promotion</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete the promotion <strong>"{deletingPromotion?.code}"</strong>? 
              This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={handleDelete} className="bg-destructive text-destructive-foreground hover:bg-destructive/90">
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
