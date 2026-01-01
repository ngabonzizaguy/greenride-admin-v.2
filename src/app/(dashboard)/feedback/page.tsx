'use client';

import { useState, useEffect, useCallback } from 'react';
import {
  MessageSquare,
  Search,
  Filter,
  MoreHorizontal,
  Eye,
  CheckCircle,
  Clock,
  AlertTriangle,
  AlertCircle,
  X,
  Send,
  User,
  Car,
  MapPin,
  Calendar,
  Star,
  ChevronLeft,
  ChevronRight,
  Download,
  RefreshCw
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
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
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';
import { toast } from 'sonner';

// ============================================
// TYPES
// ============================================

type FeedbackCategory = 'driver' | 'vehicle' | 'pricing' | 'safety' | 'app' | 'payment' | 'other';
type FeedbackStatus = 'pending' | 'reviewing' | 'resolved' | 'closed';
type FeedbackSeverity = 'low' | 'medium' | 'high' | 'critical';

interface Feedback {
  id: string;
  feedback_id: string;
  order_id?: string;
  user_id: string;
  user_name: string;
  user_phone: string;
  driver_id?: string;
  driver_name?: string;
  category: FeedbackCategory;
  severity: FeedbackSeverity;
  title: string;
  content: string;
  rating?: number;
  attachments?: string[];
  status: FeedbackStatus;
  admin_response?: string;
  assigned_to?: string;
  created_at: number;
  updated_at: number;
  resolved_at?: number;
}

// ============================================
// MOCK DATA
// ============================================

const MOCK_FEEDBACK: Feedback[] = [
  {
    id: '1',
    feedback_id: 'FB001',
    order_id: 'ORD001',
    user_id: 'USR001',
    user_name: 'John Doe',
    user_phone: '+250788111111',
    driver_id: 'DRV003',
    driver_name: 'Paul Rwema',
    category: 'driver',
    severity: 'high',
    title: 'Rude driver behavior',
    content: 'The driver was very rude and refused to help with luggage. He was also talking loudly on the phone the entire trip which made me very uncomfortable.',
    rating: 2,
    status: 'pending',
    created_at: Date.now() - 2 * 60 * 60 * 1000,
    updated_at: Date.now() - 2 * 60 * 60 * 1000,
  },
  {
    id: '2',
    feedback_id: 'FB002',
    order_id: 'ORD002',
    user_id: 'USR002',
    user_name: 'Jane Smith',
    user_phone: '+250788222222',
    driver_id: 'DRV001',
    driver_name: 'Peter Mutombo',
    category: 'vehicle',
    severity: 'medium',
    title: 'Car AC not working',
    content: 'The air conditioning was not working and it was very hot during the ride. The driver said he would fix it but it never worked.',
    rating: 3,
    status: 'reviewing',
    admin_response: 'We are contacting the driver to verify the AC issue.',
    created_at: Date.now() - 5 * 60 * 60 * 1000,
    updated_at: Date.now() - 1 * 60 * 60 * 1000,
  },
  {
    id: '3',
    feedback_id: 'FB003',
    user_id: 'USR003',
    user_name: 'Mike Johnson',
    user_phone: '+250788333333',
    category: 'pricing',
    severity: 'low',
    title: 'Fare was higher than estimate',
    content: 'The final fare was RWF 500 more than the initial estimate. I understand traffic can affect this but it seems too much.',
    status: 'resolved',
    admin_response: 'We reviewed the trip and found traffic conditions caused the delay. A RWF 300 credit has been added to your account as a goodwill gesture.',
    created_at: Date.now() - 24 * 60 * 60 * 1000,
    updated_at: Date.now() - 12 * 60 * 60 * 1000,
    resolved_at: Date.now() - 12 * 60 * 60 * 1000,
  },
  {
    id: '4',
    feedback_id: 'FB004',
    order_id: 'ORD004',
    user_id: 'USR004',
    user_name: 'Sarah Wilson',
    user_phone: '+250788444444',
    driver_id: 'DRV002',
    driver_name: 'David Kagame',
    category: 'safety',
    severity: 'critical',
    title: 'Dangerous driving',
    content: 'Driver was speeding and running red lights. I was very scared and asked him to slow down but he ignored me. This is unacceptable!',
    rating: 1,
    status: 'pending',
    created_at: Date.now() - 30 * 60 * 1000,
    updated_at: Date.now() - 30 * 60 * 1000,
  },
  {
    id: '5',
    feedback_id: 'FB005',
    user_id: 'USR005',
    user_name: 'Chris Brown',
    user_phone: '+250788555555',
    category: 'app',
    severity: 'low',
    title: 'App crashes when booking',
    content: 'The app keeps crashing when I try to book a ride. I have to restart it multiple times before it works.',
    status: 'closed',
    admin_response: 'This issue was fixed in app version 2.1.0. Please update your app from the store.',
    created_at: Date.now() - 48 * 60 * 60 * 1000,
    updated_at: Date.now() - 24 * 60 * 60 * 1000,
    resolved_at: Date.now() - 24 * 60 * 60 * 1000,
  },
  {
    id: '6',
    feedback_id: 'FB006',
    order_id: 'ORD006',
    user_id: 'USR001',
    user_name: 'John Doe',
    user_phone: '+250788111111',
    category: 'payment',
    severity: 'high',
    title: 'Double charged for ride',
    content: 'I was charged twice for the same ride. My bank shows two transactions of RWF 4,500 each. Please refund one.',
    status: 'reviewing',
    created_at: Date.now() - 4 * 60 * 60 * 1000,
    updated_at: Date.now() - 2 * 60 * 60 * 1000,
  },
  {
    id: '7',
    feedback_id: 'FB007',
    order_id: 'ORD007',
    user_id: 'USR002',
    user_name: 'Jane Smith',
    user_phone: '+250788222222',
    driver_id: 'DRV005',
    driver_name: 'Alex Munyaneza',
    category: 'other',
    severity: 'medium',
    title: 'Driver took wrong route',
    content: 'The driver took a much longer route than necessary. Google Maps showed 10 minutes but we drove for 25 minutes.',
    rating: 2,
    status: 'pending',
    created_at: Date.now() - 6 * 60 * 60 * 1000,
    updated_at: Date.now() - 6 * 60 * 60 * 1000,
  },
];

// ============================================
// HELPER FUNCTIONS
// ============================================

const getCategoryIcon = (category: FeedbackCategory) => {
  switch (category) {
    case 'driver': return <User className="h-4 w-4" />;
    case 'vehicle': return <Car className="h-4 w-4" />;
    case 'pricing': return <span className="text-sm font-bold">RWF</span>;
    case 'safety': return <AlertTriangle className="h-4 w-4" />;
    case 'app': return <span className="text-sm">📱</span>;
    case 'payment': return <span className="text-sm">💳</span>;
    default: return <MessageSquare className="h-4 w-4" />;
  }
};

const getCategoryBadge = (category: FeedbackCategory) => {
  const styles: Record<FeedbackCategory, string> = {
    driver: 'bg-blue-100 text-blue-700',
    vehicle: 'bg-purple-100 text-purple-700',
    pricing: 'bg-yellow-100 text-yellow-700',
    safety: 'bg-red-100 text-red-700',
    app: 'bg-gray-100 text-gray-700',
    payment: 'bg-green-100 text-green-700',
    other: 'bg-slate-100 text-slate-700',
  };
  return (
    <Badge className={`${styles[category]} gap-1`}>
      {getCategoryIcon(category)}
      <span className="capitalize">{category}</span>
    </Badge>
  );
};

const getStatusBadge = (status: FeedbackStatus) => {
  switch (status) {
    case 'pending':
      return <Badge className="bg-yellow-100 text-yellow-700 gap-1"><Clock className="h-3 w-3" />Pending</Badge>;
    case 'reviewing':
      return <Badge className="bg-blue-100 text-blue-700 gap-1"><Eye className="h-3 w-3" />Reviewing</Badge>;
    case 'resolved':
      return <Badge className="bg-green-100 text-green-700 gap-1"><CheckCircle className="h-3 w-3" />Resolved</Badge>;
    case 'closed':
      return <Badge className="bg-gray-100 text-gray-700 gap-1"><X className="h-3 w-3" />Closed</Badge>;
    default:
      return <Badge variant="secondary">{status}</Badge>;
  }
};

const getSeverityBadge = (severity: FeedbackSeverity) => {
  switch (severity) {
    case 'low':
      return <Badge variant="outline" className="text-green-600 border-green-300">Low</Badge>;
    case 'medium':
      return <Badge variant="outline" className="text-yellow-600 border-yellow-300">Medium</Badge>;
    case 'high':
      return <Badge variant="outline" className="text-orange-600 border-orange-300">High</Badge>;
    case 'critical':
      return <Badge className="bg-red-500 text-white">Critical</Badge>;
    default:
      return <Badge variant="outline">{severity}</Badge>;
  }
};

const formatTimeAgo = (timestamp: number) => {
  const seconds = Math.floor((Date.now() - timestamp) / 1000);
  if (seconds < 60) return 'Just now';
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  return `${days}d ago`;
};

// ============================================
// COMPONENT
// ============================================

export default function FeedbackPage() {
  const [feedback, setFeedback] = useState<Feedback[]>(MOCK_FEEDBACK);
  const [search, setSearch] = useState('');
  const [categoryFilter, setCategoryFilter] = useState<string>('all');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [severityFilter, setSeverityFilter] = useState<string>('all');
  const [page, setPage] = useState(1);
  const [selectedFeedback, setSelectedFeedback] = useState<Feedback | null>(null);
  const [isDetailOpen, setIsDetailOpen] = useState(false);
  const [isResponseDialogOpen, setIsResponseDialogOpen] = useState(false);
  const [responseText, setResponseText] = useState('');
  const [newStatus, setNewStatus] = useState<FeedbackStatus>('reviewing');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const limit = 10;

  // Filter feedback
  const filteredFeedback = feedback.filter(f => {
    const matchesSearch = search === '' || 
      f.title.toLowerCase().includes(search.toLowerCase()) ||
      f.content.toLowerCase().includes(search.toLowerCase()) ||
      f.user_name.toLowerCase().includes(search.toLowerCase()) ||
      f.feedback_id.toLowerCase().includes(search.toLowerCase());
    
    const matchesCategory = categoryFilter === 'all' || f.category === categoryFilter;
    const matchesStatus = statusFilter === 'all' || f.status === statusFilter;
    const matchesSeverity = severityFilter === 'all' || f.severity === severityFilter;
    
    return matchesSearch && matchesCategory && matchesStatus && matchesSeverity;
  });

  // Pagination
  const totalPages = Math.ceil(filteredFeedback.length / limit);
  const paginatedFeedback = filteredFeedback.slice((page - 1) * limit, page * limit);

  // Stats
  const stats = {
    total: feedback.length,
    pending: feedback.filter(f => f.status === 'pending').length,
    reviewing: feedback.filter(f => f.status === 'reviewing').length,
    resolved: feedback.filter(f => f.status === 'resolved').length,
    critical: feedback.filter(f => f.severity === 'critical' && f.status !== 'resolved' && f.status !== 'closed').length,
  };

  // Open detail view
  const openDetail = (item: Feedback) => {
    setSelectedFeedback(item);
    setIsDetailOpen(true);
  };

  // Open response dialog
  const openResponseDialog = (item: Feedback) => {
    setSelectedFeedback(item);
    setResponseText(item.admin_response || '');
    setNewStatus(item.status === 'pending' ? 'reviewing' : item.status);
    setIsResponseDialogOpen(true);
  };

  // Submit response
  const handleSubmitResponse = () => {
    if (!selectedFeedback) return;
    
    setIsSubmitting(true);
    
    // Simulate API call
    setTimeout(() => {
      setFeedback(prev => prev.map(f => {
        if (f.id === selectedFeedback.id) {
          return {
            ...f,
            status: newStatus,
            admin_response: responseText,
            updated_at: Date.now(),
            resolved_at: (newStatus === 'resolved' || newStatus === 'closed') ? Date.now() : f.resolved_at,
          };
        }
        return f;
      }));
      
      setIsResponseDialogOpen(false);
      setIsSubmitting(false);
      toast.success(`Feedback ${selectedFeedback.feedback_id} updated successfully!`);
      
      // Update selected feedback if detail is open
      if (isDetailOpen) {
        setSelectedFeedback(prev => prev ? {
          ...prev,
          status: newStatus,
          admin_response: responseText,
          updated_at: Date.now(),
        } : null);
      }
    }, 500);
  };

  // Quick status update
  const quickStatusUpdate = (item: Feedback, status: FeedbackStatus) => {
    setFeedback(prev => prev.map(f => {
      if (f.id === item.id) {
        return { ...f, status, updated_at: Date.now() };
      }
      return f;
    }));
    toast.success(`Status updated to ${status}`);
  };

  // Export to CSV
  const exportToCSV = () => {
    const headers = ['ID', 'Category', 'Severity', 'Status', 'Title', 'User', 'Phone', 'Driver', 'Created', 'Response'];
    const rows = filteredFeedback.map(f => [
      f.feedback_id,
      f.category,
      f.severity,
      f.status,
      `"${f.title.replace(/"/g, '""')}"`,
      f.user_name,
      f.user_phone,
      f.driver_name || '-',
      new Date(f.created_at).toISOString(),
      f.admin_response ? `"${f.admin_response.replace(/"/g, '""')}"` : '-'
    ]);
    
    const csv = [headers.join(','), ...rows.map(r => r.join(','))].join('\n');
    const blob = new Blob([csv], { type: 'text/csv' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `feedback_${new Date().toISOString().split('T')[0]}.csv`;
    a.click();
    toast.success('Feedback exported to CSV');
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Feedback & Complaints</h1>
          <p className="text-muted-foreground">
            Manage customer feedback and resolve issues
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" size="sm" onClick={exportToCSV}>
            <Download className="h-4 w-4 mr-2" />
            Export
          </Button>
          <Button variant="outline" size="sm" onClick={() => setFeedback([...MOCK_FEEDBACK])}>
            <RefreshCw className="h-4 w-4 mr-2" />
            Refresh
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-5">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Total</p>
                <p className="text-2xl font-bold">{stats.total}</p>
              </div>
              <MessageSquare className="h-5 w-5 text-muted-foreground" />
            </div>
          </CardContent>
        </Card>
        <Card className="border-yellow-200 bg-yellow-50">
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-yellow-700">Pending</p>
                <p className="text-2xl font-bold text-yellow-900">{stats.pending}</p>
              </div>
              <Clock className="h-5 w-5 text-yellow-600" />
            </div>
          </CardContent>
        </Card>
        <Card className="border-blue-200 bg-blue-50">
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-blue-700">Reviewing</p>
                <p className="text-2xl font-bold text-blue-900">{stats.reviewing}</p>
              </div>
              <Eye className="h-5 w-5 text-blue-600" />
            </div>
          </CardContent>
        </Card>
        <Card className="border-green-200 bg-green-50">
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-green-700">Resolved</p>
                <p className="text-2xl font-bold text-green-900">{stats.resolved}</p>
              </div>
              <CheckCircle className="h-5 w-5 text-green-600" />
            </div>
          </CardContent>
        </Card>
        {stats.critical > 0 && (
          <Card className="border-red-300 bg-red-50">
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-red-700">Critical</p>
                  <p className="text-2xl font-bold text-red-900">{stats.critical}</p>
                </div>
                <AlertCircle className="h-5 w-5 text-red-600" />
              </div>
            </CardContent>
          </Card>
        )}
      </div>

      {/* Filters */}
      <Card>
        <CardContent className="p-4">
          <div className="flex flex-col gap-4 md:flex-row md:items-center">
            <div className="relative flex-1 max-w-md">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Search by title, content, user..."
                className="pl-10"
                value={search}
                onChange={(e) => { setSearch(e.target.value); setPage(1); }}
              />
            </div>
            <div className="flex gap-2 flex-wrap">
              <Select value={categoryFilter} onValueChange={(v) => { setCategoryFilter(v); setPage(1); }}>
                <SelectTrigger className="w-[140px]">
                  <SelectValue placeholder="Category" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Categories</SelectItem>
                  <SelectItem value="driver">Driver</SelectItem>
                  <SelectItem value="vehicle">Vehicle</SelectItem>
                  <SelectItem value="pricing">Pricing</SelectItem>
                  <SelectItem value="safety">Safety</SelectItem>
                  <SelectItem value="app">App</SelectItem>
                  <SelectItem value="payment">Payment</SelectItem>
                  <SelectItem value="other">Other</SelectItem>
                </SelectContent>
              </Select>
              <Select value={statusFilter} onValueChange={(v) => { setStatusFilter(v); setPage(1); }}>
                <SelectTrigger className="w-[130px]">
                  <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Status</SelectItem>
                  <SelectItem value="pending">Pending</SelectItem>
                  <SelectItem value="reviewing">Reviewing</SelectItem>
                  <SelectItem value="resolved">Resolved</SelectItem>
                  <SelectItem value="closed">Closed</SelectItem>
                </SelectContent>
              </Select>
              <Select value={severityFilter} onValueChange={(v) => { setSeverityFilter(v); setPage(1); }}>
                <SelectTrigger className="w-[130px]">
                  <SelectValue placeholder="Severity" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Severity</SelectItem>
                  <SelectItem value="low">Low</SelectItem>
                  <SelectItem value="medium">Medium</SelectItem>
                  <SelectItem value="high">High</SelectItem>
                  <SelectItem value="critical">Critical</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Feedback Table */}
      <Card>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-[100px]">ID</TableHead>
                <TableHead>Issue</TableHead>
                <TableHead>User</TableHead>
                <TableHead>Category</TableHead>
                <TableHead>Severity</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Created</TableHead>
                <TableHead className="w-[50px]"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {paginatedFeedback.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={8} className="text-center py-12 text-muted-foreground">
                    No feedback found matching your filters.
                  </TableCell>
                </TableRow>
              ) : (
                paginatedFeedback.map((item) => (
                  <TableRow 
                    key={item.id} 
                    className={`cursor-pointer hover:bg-muted/50 ${item.severity === 'critical' && item.status === 'pending' ? 'bg-red-50' : ''}`}
                    onClick={() => openDetail(item)}
                  >
                    <TableCell className="font-mono text-sm">{item.feedback_id}</TableCell>
                    <TableCell>
                      <div className="max-w-[300px]">
                        <p className="font-medium truncate">{item.title}</p>
                        <p className="text-sm text-muted-foreground truncate">{item.content}</p>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <Avatar className="h-6 w-6">
                          <AvatarFallback className="text-xs bg-primary/10 text-primary">
                            {item.user_name.split(' ').map(n => n[0]).join('')}
                          </AvatarFallback>
                        </Avatar>
                        <span className="text-sm">{item.user_name}</span>
                      </div>
                    </TableCell>
                    <TableCell>{getCategoryBadge(item.category)}</TableCell>
                    <TableCell>{getSeverityBadge(item.severity)}</TableCell>
                    <TableCell>{getStatusBadge(item.status)}</TableCell>
                    <TableCell className="text-sm text-muted-foreground">{formatTimeAgo(item.created_at)}</TableCell>
                    <TableCell onClick={(e) => e.stopPropagation()}>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="icon" className="h-8 w-8">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem onClick={() => openDetail(item)}>
                            <Eye className="h-4 w-4 mr-2" />
                            View Details
                          </DropdownMenuItem>
                          <DropdownMenuItem onClick={() => openResponseDialog(item)}>
                            <Send className="h-4 w-4 mr-2" />
                            Respond
                          </DropdownMenuItem>
                          <DropdownMenuSeparator />
                          {item.status !== 'reviewing' && (
                            <DropdownMenuItem onClick={() => quickStatusUpdate(item, 'reviewing')}>
                              <Eye className="h-4 w-4 mr-2" />
                              Mark as Reviewing
                            </DropdownMenuItem>
                          )}
                          {item.status !== 'resolved' && (
                            <DropdownMenuItem onClick={() => quickStatusUpdate(item, 'resolved')}>
                              <CheckCircle className="h-4 w-4 mr-2" />
                              Mark as Resolved
                            </DropdownMenuItem>
                          )}
                          {item.status !== 'closed' && (
                            <DropdownMenuItem onClick={() => quickStatusUpdate(item, 'closed')}>
                              <X className="h-4 w-4 mr-2" />
                              Close
                            </DropdownMenuItem>
                          )}
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-between border-t px-4 py-3">
              <p className="text-sm text-muted-foreground">
                Showing {((page - 1) * limit) + 1}-{Math.min(page * limit, filteredFeedback.length)} of {filteredFeedback.length}
              </p>
              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setPage(p => Math.max(1, p - 1))}
                  disabled={page === 1}
                >
                  <ChevronLeft className="h-4 w-4" />
                </Button>
                <span className="text-sm">Page {page} of {totalPages}</span>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                  disabled={page === totalPages}
                >
                  <ChevronRight className="h-4 w-4" />
                </Button>
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Detail Sheet */}
      <Sheet open={isDetailOpen} onOpenChange={setIsDetailOpen}>
        <SheetContent className="w-full sm:max-w-lg overflow-y-auto">
          {selectedFeedback && (
            <>
              <SheetHeader>
                <div className="flex items-center gap-2">
                  <SheetTitle>{selectedFeedback.feedback_id}</SheetTitle>
                  {getSeverityBadge(selectedFeedback.severity)}
                </div>
                <SheetDescription>{selectedFeedback.title}</SheetDescription>
              </SheetHeader>
              
              <div className="mt-6 space-y-6">
                {/* Status */}
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Status</span>
                  {getStatusBadge(selectedFeedback.status)}
                </div>

                {/* Category */}
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Category</span>
                  {getCategoryBadge(selectedFeedback.category)}
                </div>

                <Separator />

                {/* User Info */}
                <div>
                  <h4 className="font-medium mb-3 flex items-center gap-2">
                    <User className="h-4 w-4" />
                    Submitted By
                  </h4>
                  <div className="space-y-2 text-sm">
                    <p><span className="text-muted-foreground">Name:</span> {selectedFeedback.user_name}</p>
                    <p><span className="text-muted-foreground">Phone:</span> {selectedFeedback.user_phone}</p>
                    {selectedFeedback.order_id && (
                      <p><span className="text-muted-foreground">Order:</span> {selectedFeedback.order_id}</p>
                    )}
                  </div>
                </div>

                {/* Driver Info */}
                {selectedFeedback.driver_name && (
                  <>
                    <Separator />
                    <div>
                      <h4 className="font-medium mb-3 flex items-center gap-2">
                        <Car className="h-4 w-4" />
                        Related Driver
                      </h4>
                      <p className="text-sm">{selectedFeedback.driver_name}</p>
                    </div>
                  </>
                )}

                <Separator />

                {/* Content */}
                <div>
                  <h4 className="font-medium mb-3">Description</h4>
                  <p className="text-sm text-muted-foreground whitespace-pre-wrap">
                    {selectedFeedback.content}
                  </p>
                </div>

                {/* Rating */}
                {selectedFeedback.rating && (
                  <div className="flex items-center gap-2">
                    <Star className="h-4 w-4 text-yellow-500 fill-yellow-500" />
                    <span className="text-sm">{selectedFeedback.rating}/5 rating given</span>
                  </div>
                )}

                {/* Admin Response */}
                {selectedFeedback.admin_response && (
                  <>
                    <Separator />
                    <div>
                      <h4 className="font-medium mb-3 text-primary">Admin Response</h4>
                      <p className="text-sm bg-primary/5 p-3 rounded-lg">
                        {selectedFeedback.admin_response}
                      </p>
                    </div>
                  </>
                )}

                {/* Timestamps */}
                <Separator />
                <div className="text-xs text-muted-foreground space-y-1">
                  <p>Created: {new Date(selectedFeedback.created_at).toLocaleString()}</p>
                  <p>Updated: {new Date(selectedFeedback.updated_at).toLocaleString()}</p>
                  {selectedFeedback.resolved_at && (
                    <p>Resolved: {new Date(selectedFeedback.resolved_at).toLocaleString()}</p>
                  )}
                </div>

                {/* Actions */}
                <div className="flex gap-2 pt-4">
                  <Button className="flex-1" onClick={() => { setIsDetailOpen(false); openResponseDialog(selectedFeedback); }}>
                    <Send className="h-4 w-4 mr-2" />
                    Respond
                  </Button>
                  <Button variant="outline" onClick={() => setIsDetailOpen(false)}>
                    Close
                  </Button>
                </div>
              </div>
            </>
          )}
        </SheetContent>
      </Sheet>

      {/* Response Dialog */}
      <Dialog open={isResponseDialogOpen} onOpenChange={setIsResponseDialogOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Respond to Feedback</DialogTitle>
            <DialogDescription>
              {selectedFeedback?.feedback_id} - {selectedFeedback?.title}
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <Label>Update Status</Label>
              <Select value={newStatus} onValueChange={(v) => setNewStatus(v as FeedbackStatus)}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="pending">Pending</SelectItem>
                  <SelectItem value="reviewing">Reviewing</SelectItem>
                  <SelectItem value="resolved">Resolved</SelectItem>
                  <SelectItem value="closed">Closed</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label>Admin Response</Label>
              <Textarea
                placeholder="Write your response to the customer..."
                value={responseText}
                onChange={(e) => setResponseText(e.target.value)}
                rows={4}
              />
              <p className="text-xs text-muted-foreground">
                This response may be shared with the customer.
              </p>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsResponseDialogOpen(false)} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button onClick={handleSubmitResponse} disabled={isSubmitting}>
              {isSubmitting ? 'Saving...' : 'Save Response'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}

