'use client';

import { use, useState, useEffect } from 'react';
import Link from 'next/link';
import { 
  ArrowLeft, 
  Phone, 
  Mail, 
  MapPin,
  Edit,
  Ban,
  Calendar,
  CreditCard,
  Star,
  Loader2
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Separator } from '@/components/ui/separator';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { apiClient } from '@/lib/api-client';

interface UserData {
  id: string;
  user_id?: string;
  name: string;
  display_name?: string;
  first_name?: string;
  last_name?: string;
  email: string;
  phone: string;
  avatar?: string | null;
  status: string;
  totalTrips: number;
  total_trips?: number;
  totalSpent: number;
  total_spent?: number;
  avgRating: number;
  score?: number;
  lastTrip?: string;
  joinedAt: string;
  created_at?: number;
  address?: string;
}

interface Trip {
  id: string;
  order_id?: string;
  date: string;
  created_at?: number;
  driver: string;
  provider_name?: string;
  pickup: string;
  pickup_address?: string;
  dropoff: string;
  dropoff_address?: string;
  fare: number;
  payment_amount?: number;
  rating?: number | null;
  status: string;
}

export default function UserDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const [user, setUser] = useState<UserData | null>(null);
  const [recentTrips, setRecentTrips] = useState<Trip[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchUserData = async () => {
      setIsLoading(true);
      setError(null);
      try {
        // Fetch user details
        const response = await apiClient.getUserDetail(id);
        if (response.code === '0000' && response.data) {
          const userData = response.data as Record<string, unknown>;
          setUser({
            id: (userData.user_id as string) || id,
            name: (userData.display_name as string) || 
                  `${userData.first_name || ''} ${userData.last_name || ''}`.trim() || 
                  'Unknown User',
            email: (userData.email as string) || '',
            phone: (userData.phone as string) || '',
            avatar: userData.avatar as string | null,
            status: (userData.status as string) || 'active',
            totalTrips: (userData.total_trips as number) || 0,
            totalSpent: (userData.total_spent as number) || 0,
            avgRating: (userData.score as number) || 0,
            joinedAt: userData.created_at 
              ? new Date(userData.created_at as number).toISOString().split('T')[0]
              : 'N/A',
            address: (userData.address as string) || '',
          });
        } else {
          setError('User not found');
        }

        // Fetch user's trips
        try {
          const ordersResponse = await apiClient.searchOrders({ user_id: id, page: 1, limit: 10 });
          if (ordersResponse.code === '0000' && ordersResponse.data?.records) {
            const orders = ordersResponse.data.records as Array<Record<string, unknown>>;
            setRecentTrips(orders.map((order) => ({
              id: (order.order_id as string) || String(order.id),
              date: order.created_at 
                ? new Date(order.created_at as number).toISOString().split('T')[0]
                : 'N/A',
              driver: (order.provider_name as string) || 'Driver',
              pickup: (order.pickup_address as string) || 'N/A',
              dropoff: (order.dropoff_address as string) || 'N/A',
              fare: (order.payment_amount as number) || 0,
              rating: order.rating as number | null,
              status: (order.status as string) || 'pending',
            })));
          }
        } catch {
          // Silently fail - trips are optional
        }
      } catch (err) {
        console.error('Failed to fetch user:', err);
        setError('Failed to load user details');
      } finally {
        setIsLoading(false);
      }
    };

    fetchUserData();
  }, [id]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  if (error || !user) {
    return (
      <div className="space-y-6">
        <Link href="/users" className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground">
          <ArrowLeft className="h-4 w-4" />
          Back to Users
        </Link>
        <Card>
          <CardContent className="p-6 text-center text-muted-foreground">
            {error || 'User not found'}
          </CardContent>
        </Card>
      </div>
    );
  }

  const paymentMethods = [
    { type: 'MoMo', number: '**** ****', isDefault: true },
  ];

  return (
    <div className="space-y-6">
      {/* Back Button */}
      <Link
        href="/users"
        className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground"
      >
        <ArrowLeft className="h-4 w-4" />
        Back to Users
      </Link>

      <div className="grid gap-6 lg:grid-cols-3">
        {/* Profile Card */}
        <Card className="lg:col-span-1">
          <CardContent className="p-6">
            <div className="flex flex-col items-center text-center">
              <Avatar className="h-24 w-24">
                <AvatarImage src={user.avatar || undefined} />
                <AvatarFallback className="bg-blue-100 text-blue-600 text-2xl">
                  {user.name.split(' ').map((n) => n[0]).join('')}
                </AvatarFallback>
              </Avatar>
              <h2 className="mt-4 text-xl font-bold">{user.name}</h2>
              <Badge className={user.status === 'active' ? 'bg-green-100 text-green-700 mt-2' : 'bg-red-100 text-red-700 mt-2'}>
                {user.status === 'active' ? 'Active' : 'Suspended'}
              </Badge>

              <Separator className="my-4" />

              <div className="w-full space-y-3 text-left">
                <div className="flex items-center gap-3">
                  <Phone className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm">{user.phone}</span>
                </div>
                <div className="flex items-center gap-3">
                  <Mail className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm">{user.email}</span>
                </div>
                <div className="flex items-center gap-3">
                  <MapPin className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm">{user.address}</span>
                </div>
                <div className="flex items-center gap-3">
                  <Calendar className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm">
                    Joined {new Date(user.joinedAt).toLocaleDateString('en-US', {
                      month: 'long',
                      year: 'numeric',
                    })}
                  </span>
                </div>
              </div>

              <Separator className="my-4" />

              <div className="w-full space-y-2">
                <Button variant="outline" className="w-full gap-2">
                  <Edit className="h-4 w-4" />
                  Edit Profile
                </Button>
                <Button variant="outline" className="w-full gap-2 text-red-600 hover:text-red-600">
                  <Ban className="h-4 w-4" />
                  Suspend User
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Main Content */}
        <div className="lg:col-span-2 space-y-6">
          {/* Stats */}
          <div className="grid gap-4 md:grid-cols-4">
            <Card>
              <CardContent className="p-4">
                <p className="text-sm text-muted-foreground">Total Trips</p>
                <p className="text-2xl font-bold">{user.totalTrips}</p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4">
                <p className="text-sm text-muted-foreground">Total Spent</p>
                <p className="text-2xl font-bold">RWF {user.totalSpent.toLocaleString()}</p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4">
                <p className="text-sm text-muted-foreground">Avg Rating Given</p>
                <div className="flex items-center gap-1">
                  <Star className="h-5 w-5 fill-yellow-400 text-yellow-400" />
                  <span className="text-2xl font-bold">{user.avgRating}</span>
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4">
                <p className="text-sm text-muted-foreground">Last Trip</p>
                <p className="text-2xl font-bold">
                  {user.lastTrip ? new Date(user.lastTrip).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }) : 'N/A'}
                </p>
              </CardContent>
            </Card>
          </div>

          <Tabs defaultValue="trips" className="w-full">
            <TabsList>
              <TabsTrigger value="trips">Trip History</TabsTrigger>
              <TabsTrigger value="payments">Payment Methods</TabsTrigger>
            </TabsList>

            {/* Trips Tab */}
            <TabsContent value="trips">
              <Card>
                <CardHeader>
                  <CardTitle>Recent Trips</CardTitle>
                </CardHeader>
                <CardContent>
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Date</TableHead>
                        <TableHead>Driver</TableHead>
                        <TableHead>Route</TableHead>
                        <TableHead className="text-right">Fare</TableHead>
                        <TableHead>Rating</TableHead>
                        <TableHead>Status</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {recentTrips.map((trip) => (
                        <TableRow key={trip.id}>
                          <TableCell>{trip.date}</TableCell>
                          <TableCell>{trip.driver}</TableCell>
                          <TableCell>
                            <span className="text-muted-foreground">{trip.pickup}</span>
                            <span className="mx-1">→</span>
                            <span>{trip.dropoff}</span>
                          </TableCell>
                          <TableCell className="text-right font-medium">
                            RWF {trip.fare.toLocaleString()}
                          </TableCell>
                          <TableCell>
                            {trip.rating ? (
                              <div className="flex items-center gap-1">
                                <Star className="h-4 w-4 fill-yellow-400 text-yellow-400" />
                                {trip.rating}
                              </div>
                            ) : '-'}
                          </TableCell>
                          <TableCell>
                            {trip.status === 'completed' ? (
                              <Badge className="bg-green-100 text-green-700">Completed</Badge>
                            ) : (
                              <Badge className="bg-red-100 text-red-700">Cancelled</Badge>
                            )}
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </CardContent>
              </Card>
            </TabsContent>

            {/* Payments Tab */}
            <TabsContent value="payments">
              <Card>
                <CardHeader>
                  <CardTitle>Payment Methods</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  {paymentMethods.map((method, idx) => (
                    <div key={idx} className="flex items-center justify-between p-4 border rounded-lg">
                      <div className="flex items-center gap-4">
                        <CreditCard className="h-6 w-6 text-muted-foreground" />
                        <div>
                          <p className="font-medium">{method.type}</p>
                          <p className="text-sm text-muted-foreground">{method.number}</p>
                        </div>
                      </div>
                      {method.isDefault && (
                        <Badge variant="outline">Default</Badge>
                      )}
                    </div>
                  ))}
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </div>
  );
}
