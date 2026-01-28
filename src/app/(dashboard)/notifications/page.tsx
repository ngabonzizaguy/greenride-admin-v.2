'use client';

import { useState, useEffect } from 'react';
import { 
  Send, 
  Bell,
  Users,
  Car,
  Clock,
  CheckCircle,
  Filter,
  Search
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Textarea } from '@/components/ui/textarea';
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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs';
import { useNotificationStore } from '@/stores/notification-store';
import { apiClient } from '@/lib/api-client';

// Add textarea component
import * as React from 'react';

// Mock notification history
const notificationHistory = [
  {
    id: '1',
    title: 'Holiday Season Bonus',
    message: 'Earn extra 20% on all rides this weekend!',
    audience: 'drivers',
    sentAt: '2024-12-28 10:00',
    deliveredCount: 234,
    openRate: 78,
    status: 'delivered',
  },
  {
    id: '2',
    title: 'New Year Discount',
    message: 'Use code NEWYEAR25 for 25% off your next ride',
    audience: 'users',
    sentAt: '2024-12-27 14:30',
    deliveredCount: 1234,
    openRate: 45,
    status: 'delivered',
  },
  {
    id: '3',
    title: 'Service Maintenance',
    message: 'App will be under maintenance from 2 AM to 4 AM',
    audience: 'all',
    sentAt: '2024-12-25 18:00',
    deliveredCount: 1468,
    openRate: 62,
    status: 'delivered',
  },
  {
    id: '4',
    title: 'Christmas Special',
    message: 'Merry Christmas! Enjoy free rides today.',
    audience: 'users',
    sentAt: '2024-12-25 00:00',
    deliveredCount: 1234,
    openRate: 89,
    status: 'delivered',
  },
];

const stats = {
  totalSent: 25,
  totalDelivered: 4200,
  avgOpenRate: 68,
  scheduledCount: 2,
};

export default function NotificationsPage() {
  const [audience, setAudience] = useState<'all' | 'drivers' | 'users'>('all');
  const [title, setTitle] = useState('');
  const [message, setMessage] = useState('');
  const [schedule, setSchedule] = useState('now');
  const [scheduledAt, setScheduledAt] = useState<Date | null>(null);
  const [isSending, setIsSending] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  
  const { notifications, fetchNotifications, isLoading } = useNotificationStore();
  const [page, setPage] = useState(1);
  const [totalNotifications, setTotalNotifications] = useState(0);

  // Fetch notification history
  useEffect(() => {
    fetchNotifications({ page, limit: 20 });
  }, [page, fetchNotifications]);

  const handleSend = async () => {
    if (!title || !message) {
      setError('Title and message are required');
      return;
    }

    setIsSending(true);
    setError(null);
    setSuccessMessage(null);

    try {
      const scheduledTimestamp = schedule === 'now' 
        ? undefined 
        : scheduledAt?.getTime();

      const response = await apiClient.sendNotification({
        audience,
        type: 'system',
        category: 'marketing',
        title,
        content: message,
        summary: message.substring(0, 100),
        scheduled_at: scheduledTimestamp,
      });

      if (response.code === '0000') {
        setSuccessMessage('Notification sent successfully!');
        setTitle('');
        setMessage('');
        setSchedule('now');
        setScheduledAt(null);
        // Refresh notification history
        fetchNotifications({ page: 1, limit: 20 });
      } else {
        setError(response.msg || 'Failed to send notification');
      }
    } catch (err) {
      console.error('Failed to send notification:', err);
      setError('Failed to send notification. Please try again.');
    } finally {
      setIsSending(false);
    }
  };

  // Format notification history from store
  const notificationHistory = notifications.map(n => ({
    id: n.notification_id,
    title: n.title,
    message: n.content,
    audience: n.user_type || 'all',
    sentAt: new Date(n.created_at).toLocaleString(),
    deliveredCount: n.status === 'delivered' ? 1 : 0,
    openRate: n.is_read ? 100 : 0,
    status: n.status,
  }));

  const stats = {
    totalSent: notifications.length,
    totalDelivered: notifications.filter(n => n.status === 'delivered' || n.status === 'sent').length,
    avgOpenRate: notifications.length > 0 
      ? Math.round((notifications.filter(n => n.is_read).length / notifications.length) * 100)
      : 0,
    scheduledCount: notifications.filter(n => n.status === 'pending').length,
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Notifications</h1>
        <p className="text-muted-foreground">
          Send messages to drivers and users
        </p>
      </div>

      {/* Success/Error Messages */}
      {successMessage && (
        <div className="flex items-center gap-2 rounded-lg bg-green-50 border border-green-200 p-3 text-sm text-green-800">
          <CheckCircle className="h-4 w-4 flex-shrink-0" />
          <span>{successMessage}</span>
        </div>
      )}

      {error && (
        <div className="flex items-center gap-2 rounded-lg bg-red-50 border border-red-200 p-3 text-sm text-red-800">
          <span>{error}</span>
        </div>
      )}

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Sent This Month</p>
                <p className="text-2xl font-bold">{stats.totalSent}</p>
              </div>
              <Send className="h-5 w-5 text-primary" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Total Delivered</p>
                <p className="text-2xl font-bold">{stats.totalDelivered.toLocaleString()}</p>
              </div>
              <CheckCircle className="h-5 w-5 text-green-500" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Avg. Open Rate</p>
                <p className="text-2xl font-bold">{stats.avgOpenRate}%</p>
              </div>
              <Bell className="h-5 w-5 text-blue-500" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Scheduled</p>
                <p className="text-2xl font-bold">{stats.scheduledCount}</p>
              </div>
              <Clock className="h-5 w-5 text-yellow-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      <Tabs defaultValue="compose" className="space-y-4">
        <TabsList>
          <TabsTrigger value="compose">Compose</TabsTrigger>
          <TabsTrigger value="history">History</TabsTrigger>
        </TabsList>

        {/* Compose Tab */}
        <TabsContent value="compose">
          <Card>
            <CardHeader>
              <CardTitle>Send Notification</CardTitle>
              <CardDescription>
                Broadcast a message to your users or drivers
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              {/* Audience Selection */}
              <div className="space-y-2">
                <Label>Select Audience</Label>
                <div className="flex gap-2">
                  <Button
                    variant={audience === 'all' ? 'default' : 'outline'}
                    className="flex-1"
                    onClick={() => setAudience('all')}
                  >
                    <Users className="h-4 w-4 mr-2" />
                    All (1,468)
                  </Button>
                  <Button
                    variant={audience === 'drivers' ? 'default' : 'outline'}
                    className="flex-1"
                    onClick={() => setAudience('drivers')}
                  >
                    <Car className="h-4 w-4 mr-2" />
                    Drivers (234)
                  </Button>
                  <Button
                    variant={audience === 'users' ? 'default' : 'outline'}
                    className="flex-1"
                    onClick={() => setAudience('users')}
                  >
                    <Users className="h-4 w-4 mr-2" />
                    Users (1,234)
                  </Button>
                </div>
              </div>

              {/* Title */}
              <div className="space-y-2">
                <Label htmlFor="title">Notification Title</Label>
                <Input
                  id="title"
                  placeholder="Enter notification title..."
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                />
              </div>

              {/* Message */}
              <div className="space-y-2">
                <Label htmlFor="message">Message</Label>
                <Textarea
                  id="message"
                  placeholder="Enter your message..."
                  rows={4}
                  value={message}
                  onChange={(e) => setMessage(e.target.value)}
                />
                <p className="text-xs text-muted-foreground">
                  {message.length}/160 characters
                </p>
              </div>

              {/* Schedule */}
              <div className="space-y-2">
                <Label>Schedule</Label>
                <Select value={schedule} onValueChange={setSchedule}>
                  <SelectTrigger className="w-full md:w-[300px]">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="now">Send Immediately</SelectItem>
                    <SelectItem value="schedule">Schedule for Later</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* Send Button */}
              <div className="flex justify-end gap-2">
                <Button variant="outline">Save as Draft</Button>
                <Button 
                  onClick={handleSend}
                  disabled={!title || !message || isSending}
                  className="gap-2"
                >
                  {isSending ? (
                    <>
                      <Clock className="h-4 w-4 animate-spin" />
                      Sending...
                    </>
                  ) : (
                    <>
                      <Send className="h-4 w-4" />
                      Send Notification
                    </>
                  )}
                </Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* History Tab */}
        <TabsContent value="history">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <div>
                <CardTitle>Notification History</CardTitle>
                <CardDescription>
                  Previously sent notifications
                </CardDescription>
              </div>
              <div className="relative w-64">
                <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                <Input placeholder="Search notifications..." className="pl-10" />
              </div>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Title</TableHead>
                    <TableHead>Audience</TableHead>
                    <TableHead>Sent At</TableHead>
                    <TableHead className="text-right">Delivered</TableHead>
                    <TableHead className="text-right">Open Rate</TableHead>
                    <TableHead>Status</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {isLoading ? (
                    <TableRow>
                      <TableCell colSpan={6} className="text-center py-8 text-muted-foreground">
                        Loading notifications...
                      </TableCell>
                    </TableRow>
                  ) : notificationHistory.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={6} className="text-center py-8 text-muted-foreground">
                        No notifications sent yet
                      </TableCell>
                    </TableRow>
                  ) : (
                    notificationHistory.map((notification) => (
                    <TableRow key={notification.id}>
                      <TableCell>
                        <div>
                          <p className="font-medium">{notification.title}</p>
                          <p className="text-sm text-muted-foreground truncate max-w-[300px]">
                            {notification.message}
                          </p>
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge variant="outline" className="capitalize">
                          {notification.audience === 'all' ? (
                            <Users className="h-3 w-3 mr-1" />
                          ) : notification.audience === 'drivers' ? (
                            <Car className="h-3 w-3 mr-1" />
                          ) : (
                            <Users className="h-3 w-3 mr-1" />
                          )}
                          {notification.audience}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-muted-foreground">
                        {notification.sentAt}
                      </TableCell>
                      <TableCell className="text-right font-medium">
                        {notification.deliveredCount.toLocaleString()}
                      </TableCell>
                      <TableCell className="text-right">
                        <span className={notification.openRate >= 70 ? 'text-green-600' : notification.openRate >= 50 ? 'text-yellow-600' : 'text-red-600'}>
                          {notification.openRate}%
                        </span>
                      </TableCell>
                      <TableCell>
                        <Badge className="bg-green-100 text-green-700 hover:bg-green-100">
                          <CheckCircle className="h-3 w-3 mr-1" />
                          Delivered
                        </Badge>
                      </TableCell>
                    </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
