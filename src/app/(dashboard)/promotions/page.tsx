'use client';

import { useState } from 'react';
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
  TrendingUp
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
  DialogTrigger,
} from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

// Mock promotions data
const mockPromotions = [
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

const stats = {
  activePromotions: 3,
  totalRedemptions: 969,
  revenueImpact: -485000,
  avgDiscount: 1850,
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
  const [search, setSearch] = useState('');
  const [isCreateOpen, setIsCreateOpen] = useState(false);

  const filteredPromotions = mockPromotions.filter((promo) =>
    promo.code.toLowerCase().includes(search.toLowerCase())
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
        <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
          <DialogTrigger asChild>
            <Button className="gap-2">
              <Plus className="h-4 w-4" />
              Create Promotion
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle>Create New Promotion</DialogTitle>
              <DialogDescription>
                Set up a new promotional discount code
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4 py-4">
              <div className="space-y-2">
                <Label htmlFor="code">Promo Code</Label>
                <Input id="code" placeholder="e.g., SUMMER25" />
              </div>
              <div className="space-y-2">
                <Label>Discount Type</Label>
                <Select>
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
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="value">Value</Label>
                  <Input id="value" type="number" placeholder="25" />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="max">Max Discount (RWF)</Label>
                  <Input id="max" type="number" placeholder="5000" />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="validFrom">Valid From</Label>
                  <Input id="validFrom" type="date" />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="validUntil">Valid Until</Label>
                  <Input id="validUntil" type="date" />
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="limit">Usage Limit (optional)</Label>
                <Input id="limit" type="number" placeholder="Unlimited" />
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsCreateOpen(false)}>
                Cancel
              </Button>
              <Button onClick={() => setIsCreateOpen(false)}>Create Promotion</Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
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
              {filteredPromotions.map((promo) => (
                <TableRow key={promo.id}>
                  <TableCell>
                    <div className="flex items-center gap-2">
                      <code className="font-mono font-bold text-primary">{promo.code}</code>
                      <Button variant="ghost" size="icon" className="h-6 w-6">
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
                        <DropdownMenuItem className="gap-2">
                          <Edit className="h-4 w-4" />
                          Edit
                        </DropdownMenuItem>
                        <DropdownMenuItem className="gap-2">
                          <Copy className="h-4 w-4" />
                          Duplicate
                        </DropdownMenuItem>
                        <DropdownMenuItem className="gap-2">
                          <ToggleLeft className="h-4 w-4" />
                          {promo.status === 'active' ? 'Disable' : 'Enable'}
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
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
        </CardContent>
      </Card>
    </div>
  );
}
