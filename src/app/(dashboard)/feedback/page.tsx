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
import { apiClient } from '@/lib/api-client';
import { FeedbackCategory, FeedbackStatus, FeedbackSeverity } from '@/types';

// ============================================
// TYPES
// ============================================

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
// HELPER FUNCTIONS
// ============================================

const getCategoryIcon = (category: FeedbackCategory) => {
  switch (category) {
    case 'driver': return <User className="h-4 w-4" />;
    case 'vehicle': return <Car className="h-4 w-4" />;
    // case 'pricing': return <span className="text-sm font-bold">RWF</span>;
    case 'safety': return <AlertTriangle className="h-4 w-4" />;
    // case 'app': return <span className="text-sm">📱</span>;
    // case 'payment': return <span className="text-sm">💳</span>;
    default: return <MessageSquare className="h-4 w-4" />;
  }
};

const getCategoryBadge = (category: FeedbackCategory) => {
  const styles: Record<string, string> = {
    driver: 'bg-blue-100 text-blue-700',
    vehicle: 'bg-purple-100 text-purple-700',
    pricing: 'bg-yellow-100 text-yellow-700',
    safety: 'bg-red-100 text-red-700',
    app: 'bg-gray-100 text-gray-700',
    payment: 'bg-green-100 text-green-700',
    other: 'bg-slate-100 text-slate-700',
  };
  return (
    <Badge className={`${styles[category] || styles.other} gap-1`}>
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
  const [feedback, setFeedback] = useState<Feedback[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [search, setSearch] = useState('');
  const [categoryFilter, setCategoryFilter] = useState<string>('all');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [severityFilter, setSeverityFilter] = useState<string>('all');
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(10);
  const [totalCount, setTotalCount] = useState(0);
  
  const [selectedFeedback, setSelectedFeedback] = useState<Feedback | null>(null);
  const [isDetailOpen, setIsDetailOpen] = useState(false);
  const [isResponseDialogOpen, setIsResponseDialogOpen] = useState(false);
  const [responseText, setResponseText] = useState('');
  const [newStatus, setNewStatus] = useState<FeedbackStatus>('reviewing');
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Fetch feedback data
  const fetchFeedback = useCallback(async () => {
    setIsLoading(true);
    try {
      const response = await apiClient.searchFeedback({
        page,
        limit,
        keyword: search,
        category: categoryFilter !== 'all' ? categoryFilter : undefined,
        status: statusFilter !== 'all' ? statusFilter : undefined,
        severity: severityFilter !== 'all' ? severityFilter : undefined,
      });
      
      if (response.code === '0000') {
        // Cast the unknown[] to Feedback[]
        setFeedback(response.data.records as unknown as Feedback[]);
        setTotalCount(response.data.total);
      } else {
        toast.error('Failed to fetch feedback');
      }
    } catch (error) {
      console.error('Error fetching feedback:', error);
      toast.error('Error loading feedback');
    } finally {
      setIsLoading(false);
    }
  }, [page, limit, search, categoryFilter, statusFilter, severityFilter]);

  // Initial fetch and on filter change
  useEffect(() => {
    fetchFeedback();
  }, [fetchFeedback]);

  // Stats (in a real app, these should come from a separate API endpoint)
  // For demo, we'll calculate from the current filtered list or fetch separately if needed
  // Using a simplified stats object for now based on loaded data or a separate fetch
  const stats = {
    total: totalCount || 0, // This might be just the filtered total
    pending: 0, // Hard to get accurate global stats without dedicated endpoint
    reviewing: 0,
    resolved: 0,
    critical: 0,
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
  const handleSubmitResponse = async () => {
    if (!selectedFeedback) return;
    
    setIsSubmitting(true);
    try {
      await apiClient.updateFeedback(selectedFeedback.feedback_id, {
        status: newStatus,
        admin_response: responseText
      });
      
      toast.success(`Feedback ${selectedFeedback.feedback_id} updated successfully!`);
      setIsResponseDialogOpen(false);
      
      // Refresh list
      fetchFeedback();
      
      // Update selected item if detail is open
      if (isDetailOpen) {
        setSelectedFeedback(prev => prev ? {
          ...prev,
          status: newStatus,
          admin_response: responseText,
          updated_at: Date.now(),
        } : null);
      }
    } catch (error) {
      console.error('Failed to update feedback:', error);
      toast.error('Failed to update feedback');
    } finally {
      setIsSubmitting(false);
    }
  };

  // Quick status update
  const quickStatusUpdate = async (item: Feedback, status: FeedbackStatus) => {
    try {
      await apiClient.updateFeedback(item.feedback_id, { status });
      toast.success(`Status updated to ${status}`);
      fetchFeedback();
    } catch (error) {
      console.error('Failed to update status:', error);
      toast.error('Failed to update status');
    }
  };

  // Export to CSV
  const exportToCSV = () => {
    const headers = ['ID', 'Category', 'Severity', 'Status', 'Title', 'User', 'Phone', 'Driver', 'Created', 'Response'];
    const rows = feedback.map(f => [
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
    <div className="relative h-[calc(100vh-8rem)]">
      {/* Main Content Area */}
      <div className="h-full overflow-y-auto pb-6">
        <div className="space-y-6">
          {/* Header */}
          <div>
            <h1 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-white drop-shadow-sm">Feedback & Complaints</h1>
            <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
              Manage customer feedback, resolve issues, and ensure high service quality.
            </p>
          </div>

          {/* Stats Cards */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="glass-card p-5 rounded-xl flex items-center justify-between group">
              <div>
                <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Total Feedback</p>
                <p className="text-3xl font-bold text-gray-900 dark:text-white mt-1 group-hover:scale-105 transition-transform origin-left">{stats.total}</p>
              </div>
              <div className="h-12 w-12 rounded-xl bg-gradient-to-br from-gray-100 to-gray-200 dark:from-gray-700 dark:to-gray-800 flex items-center justify-center text-gray-600 dark:text-gray-300 shadow-inner">
                <MessageSquare className="h-6 w-6" />
              </div>
            </div>
            <div className="glass-card p-5 rounded-xl border-l-4 border-l-yellow-400 flex items-center justify-between group">
              <div>
                <p className="text-sm font-medium text-yellow-600 dark:text-yellow-500">Pending</p>
                <p className="text-3xl font-bold text-gray-900 dark:text-white mt-1 group-hover:scale-105 transition-transform origin-left">{stats.pending}</p>
              </div>
              <div className="h-12 w-12 rounded-xl bg-yellow-50/50 dark:bg-yellow-900/20 flex items-center justify-center text-yellow-600 dark:text-yellow-500 shadow-inner">
                <Clock className="h-6 w-6" />
              </div>
            </div>
            <div className="glass-card p-5 rounded-xl border-l-4 border-l-blue-500 flex items-center justify-between group">
              <div>
                <p className="text-sm font-medium text-blue-600 dark:text-blue-500">Reviewing</p>
                <p className="text-3xl font-bold text-gray-900 dark:text-white mt-1 group-hover:scale-105 transition-transform origin-left">{stats.reviewing}</p>
              </div>
              <div className="h-12 w-12 rounded-xl bg-blue-50/50 dark:bg-blue-900/20 flex items-center justify-center text-blue-600 dark:text-blue-500 shadow-inner">
                <Eye className="h-6 w-6" />
              </div>
            </div>
          </div>

          {/* Filters */}
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 glass-panel p-4 rounded-xl">
            <div className="relative flex-1 max-w-lg">
              <span className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <Search className="h-5 w-5 text-gray-400" />
              </span>
              <input 
                className="block w-full pl-10 pr-3 py-2 border border-gray-200/50 dark:border-gray-600/30 rounded-lg leading-5 bg-white/50 dark:bg-gray-800/50 text-gray-900 dark:text-gray-100 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary/50 focus:bg-white/80 sm:text-sm backdrop-blur-sm transition-all shadow-sm" 
                placeholder="Search by ID, title, user..." 
                type="text"
                value={search}
                onChange={(e) => { setSearch(e.target.value); setPage(1); }}
              />
            </div>
            <div className="flex gap-2">
              <select 
                className="block w-full pl-3 pr-10 py-2 text-base border border-gray-200/50 dark:border-gray-600/30 focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary/50 sm:text-sm rounded-lg bg-white/50 dark:bg-gray-800/50 text-gray-700 dark:text-gray-300 backdrop-blur-sm shadow-sm"
                value={categoryFilter}
                onChange={(e) => { setCategoryFilter(e.target.value); setPage(1); }}
              >
                <option value="all">All Categories</option>
                <option value="driver">Driver Behavior</option>
                <option value="vehicle">Vehicle Condition</option>
                <option value="pricing">Fare Dispute</option>
                <option value="safety">Safety Issue</option>
                <option value="app">App Issue</option>
                <option value="payment">Payment Issue</option>
                <option value="other">Other</option>
              </select>
              <select 
                className="block w-full pl-3 pr-10 py-2 text-base border border-gray-200/50 dark:border-gray-600/30 focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary/50 sm:text-sm rounded-lg bg-white/50 dark:bg-gray-800/50 text-gray-700 dark:text-gray-300 backdrop-blur-sm shadow-sm"
                value={statusFilter}
                onChange={(e) => { setStatusFilter(e.target.value); setPage(1); }}
              >
                <option value="all">All Status</option>
                <option value="pending">Pending</option>
                <option value="reviewing">Reviewing</option>
                <option value="resolved">Resolved</option>
                <option value="closed">Closed</option>
              </select>
            </div>
          </div>

          {/* Table */}
          <div className="glass-panel rounded-xl overflow-hidden shadow-lg mb-6">
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200/50 dark:divide-gray-700/50">
                <thead className="bg-gray-50/50 dark:bg-gray-800/30">
                  <tr>
                    <th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">ID</th>
                    <th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Issue</th>
                    <th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">User</th>
                    <th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Category</th>
                    <th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Status</th>
                    <th className="px-6 py-4 text-right text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Action</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200/50 dark:divide-gray-700/50">
                  {paginatedFeedback.length === 0 ? (
                    <tr>
                      <td colSpan={6} className="px-6 py-12 text-center text-muted-foreground">
                        No feedback found matching your filters.
                      </td>
                    </tr>
                  ) : (
                    paginatedFeedback.map((item) => (
                      <tr 
                        key={item.id} 
                        className={`hover:bg-white/40 dark:hover:bg-white/5 transition-colors group cursor-pointer ${item.severity === 'critical' ? 'bg-red-50/30 dark:bg-red-900/10 border-l-4 border-l-red-500' : ''}`}
                        onClick={() => setSelectedFeedback(item)}
                      >
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400 font-mono">{item.feedback_id}</td>
                        <td className="px-6 py-4">
                          <div className="flex flex-col">
                            <span className="text-sm font-medium text-gray-900 dark:text-white group-hover:text-primary transition-colors">{item.title}</span>
                            <span className="text-xs text-gray-500 dark:text-gray-400 truncate max-w-xs">{item.content}</span>
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="flex items-center">
                            <div className="h-8 w-8 rounded-full bg-gradient-to-br from-green-100 to-green-200 dark:from-green-900 dark:to-green-800 flex items-center justify-center text-green-700 dark:text-green-300 text-xs font-bold mr-2 shadow-sm">
                              {item.user_name.split(' ').map(n => n[0]).join('')}
                            </div>
                            <div className="text-sm text-gray-900 dark:text-white">{item.user_name}</div>
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${
                            item.category === 'safety' ? 'bg-red-100/80 text-red-800 border-red-200' :
                            item.category === 'driver' ? 'bg-blue-100/80 text-blue-800 border-blue-200' :
                            item.category === 'vehicle' ? 'bg-purple-100/80 text-purple-800 border-purple-200' :
                            'bg-gray-100/80 text-gray-800 border-gray-200'
                          }`}>
                            {item.category.charAt(0).toUpperCase() + item.category.slice(1)}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                           <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${
                            item.status === 'resolved' ? 'bg-green-100/80 text-green-800 border-green-200' :
                            item.status === 'reviewing' ? 'bg-blue-100/80 text-blue-800 border-blue-200' :
                            item.status === 'pending' ? 'bg-yellow-100/80 text-yellow-800 border-yellow-200' :
                            'bg-gray-100/80 text-gray-800 border-gray-200'
                          }`}>
                            {item.status.charAt(0).toUpperCase() + item.status.slice(1)}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                          <button 
                            className="text-gray-400 hover:text-primary transition-colors hover:scale-110 transform"
                            onClick={(e) => { e.stopPropagation(); setSelectedFeedback(item); }}
                          >
                            <ChevronRight className="h-5 w-5" />
                          </button>
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      {/* Detail Overlay Panel */}
      {selectedFeedback && (
        <div className="absolute inset-y-0 right-0 w-full md:w-[480px] glass-panel m-0 md:m-4 md:mb-20 rounded-xl shadow-[0_20px_50px_rgba(0,0,0,0.15)] transform transition-transform duration-300 ease-in-out z-40 flex flex-col border border-white/40 dark:border-white/10 h-[calc(100vh-10rem)]">
          <div className="flex items-start justify-between p-6 border-b border-gray-100/30 dark:border-gray-700/30 bg-white/80 dark:bg-gray-800/80 backdrop-blur-md rounded-t-xl shrink-0">
            <div>
              <div className="flex items-center gap-2 mb-2">
                <h2 className="text-lg font-bold text-gray-900 dark:text-white font-mono">{selectedFeedback.feedback_id}</h2>
                {selectedFeedback.severity === 'critical' && (
                  <span className="bg-red-100/80 text-red-700 dark:bg-red-900/40 dark:text-red-300 text-xs font-bold px-2 py-0.5 rounded border border-red-200/50 dark:border-red-800/50 shadow-sm">Critical</span>
                )}
              </div>
              <h3 className="text-xl font-semibold text-gray-800 dark:text-gray-100 leading-tight">{selectedFeedback.title}</h3>
            </div>
            <button 
              className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 transition-colors p-1 hover:bg-white/20 rounded-full"
              onClick={() => setSelectedFeedback(null)}
            >
              <X className="h-6 w-6" />
            </button>
          </div>
          
          <div className="flex-1 overflow-y-auto p-6 space-y-8">
            <div className="grid grid-cols-2 gap-4">
              <div className="bg-white/80 dark:bg-gray-800/80 p-4 rounded-xl border border-white/50 dark:border-white/10 shadow-sm backdrop-blur-sm">
                <span className="text-xs text-gray-500 uppercase tracking-wide">Status</span>
                <div className="mt-2 flex items-center">
                  <span className={`h-2.5 w-2.5 rounded-full mr-2 shadow-[0_0_8px_rgba(0,0,0,0.2)] ${
                    selectedFeedback.status === 'pending' ? 'bg-yellow-500' :
                    selectedFeedback.status === 'reviewing' ? 'bg-blue-500' :
                    selectedFeedback.status === 'resolved' ? 'bg-green-500' : 'bg-gray-500'
                  }`}></span>
                  <span className="text-sm font-medium text-gray-900 dark:text-white capitalize">{selectedFeedback.status}</span>
                </div>
              </div>
              <div className="bg-white/80 dark:bg-gray-800/80 p-4 rounded-xl border border-white/50 dark:border-white/10 shadow-sm backdrop-blur-sm">
                <span className="text-xs text-gray-500 uppercase tracking-wide">Category</span>
                <div className="mt-2 flex items-center">
                  <span className="text-sm font-medium text-gray-900 dark:text-white capitalize">{selectedFeedback.category}</span>
                </div>
              </div>
            </div>

            <div className="relative pl-5 border-l-2 border-primary/30">
              <h4 className="text-xs font-bold text-gray-500 uppercase mb-3 flex items-center gap-2">
                <User className="h-4 w-4" /> Submitted By
              </h4>
              <div className="bg-white/70 dark:bg-gray-800/70 rounded-xl p-4 space-y-2 border border-white/20 dark:border-white/5 backdrop-blur-sm">
                <div className="flex justify-between py-1 border-b border-gray-100/20 dark:border-gray-700/30 pb-2">
                  <span className="text-sm text-gray-500 dark:text-gray-400">Name:</span>
                  <span className="text-sm font-medium text-gray-900 dark:text-white">{selectedFeedback.user_name}</span>
                </div>
                <div className="flex justify-between py-1 border-b border-gray-100/20 dark:border-gray-700/30 pb-2">
                  <span className="text-sm text-gray-500 dark:text-gray-400">Phone:</span>
                  <span className="text-sm font-medium text-gray-900 dark:text-white font-mono">{selectedFeedback.user_phone}</span>
                </div>
                {selectedFeedback.order_id && (
                  <div className="flex justify-between py-1 pt-2">
                    <span className="text-sm text-gray-500 dark:text-gray-400">Order ID:</span>
                    <span className="text-sm font-medium text-primary hover:underline cursor-pointer">{selectedFeedback.order_id}</span>
                  </div>
                )}
              </div>
            </div>

            {selectedFeedback.driver_name && (
              <div className="relative pl-5 border-l-2 border-gray-200/50 dark:border-gray-700/50">
                <h4 className="text-xs font-bold text-gray-500 uppercase mb-3 flex items-center gap-2">
                  <Car className="h-4 w-4" /> Related Driver
                </h4>
                <div className="flex items-center justify-between p-4 bg-white/70 dark:bg-gray-800/70 rounded-xl border border-white/20 dark:border-white/5 backdrop-blur-sm hover:bg-white/80 transition-colors cursor-pointer group">
                  <div className="flex items-center">
                    <div className="h-10 w-10 bg-gradient-to-br from-gray-200 to-gray-300 dark:from-gray-600 dark:to-gray-700 rounded-full flex items-center justify-center overflow-hidden shadow-sm">
                      <User className="h-5 w-5 text-gray-500 dark:text-gray-300" />
                    </div>
                    <div className="ml-3">
                      <p className="text-sm font-medium text-gray-900 dark:text-white group-hover:text-primary transition-colors">{selectedFeedback.driver_name}</p>
                      <div className="flex text-xs text-yellow-500 items-center">
                        <Star className="h-3 w-3 fill-yellow-500" />
                        <span className="ml-1 text-gray-500">4.5</span>
                      </div>
                    </div>
                  </div>
                  <button className="text-gray-400 hover:text-primary group-hover:translate-x-1 transition-transform">
                    <ChevronRight className="h-5 w-5" />
                  </button>
                </div>
              </div>
            )}

            <div>
              <h4 className="text-xs font-bold text-gray-500 uppercase mb-3 flex items-center gap-2">
                Description
              </h4>
              <div className="bg-red-50/70 dark:bg-red-900/30 p-5 rounded-xl border border-red-100/50 dark:border-red-900/30 text-gray-800 dark:text-gray-200 text-sm leading-relaxed backdrop-blur-sm shadow-inner italic">
                "{selectedFeedback.content}"
              </div>
              {selectedFeedback.rating && (
                <div className="mt-4 flex items-center gap-3 bg-white/60 dark:bg-gray-800/60 p-3 rounded-lg w-fit border border-white/20">
                  <span className="text-xs font-semibold text-gray-500">Rating given:</span>
                  <div className="flex text-yellow-500">
                    {[1,2,3,4,5].map(star => (
                       <Star key={star} className={`h-3 w-3 ${star <= selectedFeedback.rating! ? 'fill-yellow-500 text-yellow-500' : 'text-gray-300/50'}`} />
                    ))}
                  </div>
                  <span className="text-xs text-gray-500 font-medium bg-white/80 dark:bg-black/40 px-1.5 rounded">{selectedFeedback.rating}/5</span>
                </div>
              )}
            </div>
            
            <div className="flex items-center justify-between text-xs text-gray-400 dark:text-gray-500 border-t border-gray-100/30 dark:border-gray-800 pt-4">
              <span>Created: {new Date(selectedFeedback.created_at).toLocaleString()}</span>
              <span>Updated: {new Date(selectedFeedback.updated_at).toLocaleString()}</span>
            </div>
          </div>

          <div className="p-6 border-t border-gray-100/30 dark:border-gray-700/30 bg-white/80 dark:bg-gray-800/80 backdrop-blur-md rounded-b-xl flex items-center justify-end gap-3 shrink-0">
            <button 
              className="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-white/50 dark:bg-gray-700/50 border border-gray-300/50 dark:border-gray-600/50 rounded-lg hover:bg-white dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary transition-all shadow-sm"
              onClick={() => setSelectedFeedback(null)}
            >
              Close
            </button>
            <button 
              className="flex items-center px-4 py-2 text-sm font-medium text-white bg-primary/90 hover:bg-primary rounded-lg focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary shadow-lg shadow-primary/30 transition-all transform hover:-translate-y-0.5 border border-white/20 backdrop-blur-sm"
              onClick={() => openResponseDialog(selectedFeedback)}
            >
              <Send className="h-4 w-4 mr-2" />
              Respond
            </button>
          </div>
        </div>
      )}
      
      {/* Response Dialog (Keep existing logic) */}
      <Dialog open={isResponseDialogOpen} onOpenChange={setIsResponseDialogOpen}>
        {/* ... existing dialog content ... */}
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

