'use client';

import { use, useState, useEffect } from 'react';
import Link from 'next/link';
import { 
  ArrowLeft, 
  Star, 
  Phone, 
  Mail, 
  MapPin,
  Edit,
  Ban,
  Trash2,
  Car,
  Calendar,
  TrendingUp,
  CheckCircle,
  XCircle,
  Clock,
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
import { DriverPerformanceChart } from '@/components/charts/driver-performance-chart';
import { apiClient } from '@/lib/api-client';
import { reverseGeocode } from '@/lib/geocoding';

interface DriverData {
  id: string;
  name: string;
  email: string;
  phone: string;
  avatar: string | null;
  status: string;
  rating: number;
  totalTrips: number;
  acceptanceRate: number;
  completionRate: number;
  cancellationRate: number;
  todayTrips: number;
  todayEarnings: number;
  weekEarnings: number;
  monthEarnings: number;
  lifetimeEarnings: number;
  vehicle: {
    plate: string;
    model: string;
    make: string;
    year: number;
    color: string;
    type: string;
  };
  documents: {
    license: { status: string; expiry: string };
    insurance: { status: string; expiry: string };
    registration: { status: string; expiry: string };
  };
  joinedAt: string;
  address: string;
}

interface Trip {
  id: string;
  date: string;
  passenger: string;
  pickup: string;
  dropoff: string;
  distance: number;
  fare: number;
  rating: number | null;
  status: string;
}

interface Review {
  id: number;
  passenger: string;
  rating: number;
  comment: string;
  date: string;
}

interface ActivityLogItem {
  id: number;
  event: string;
  time: string;
}

const getStatusBadge = (status: string) => {
  switch (status) {
    case 'online':
      return (
        <Badge className="bg-green-100 text-green-700 hover:bg-green-100">
          <span className="mr-1.5 h-2 w-2 rounded-full bg-green-500" />
          Online
        </Badge>
      );
    case 'on_trip':
      return (
        <Badge className="bg-yellow-100 text-yellow-700 hover:bg-yellow-100">
          <span className="mr-1.5 h-2 w-2 rounded-full bg-yellow-500" />
          On Trip
        </Badge>
      );
    case 'offline':
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
    default:
      return <Badge variant="secondary">{status}</Badge>;
  }
};

export default function DriverDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const [driver, setDriver] = useState<DriverData | null>(null);
  const [currentCoords, setCurrentCoords] = useState<{ lat: number; lng: number } | null>(null);
  const [currentLocationText, setCurrentLocationText] = useState<string>('N/A');
  const [recentTrips, setRecentTrips] = useState<Trip[]>([]);
  const [reviews, setReviews] = useState<Review[]>([]);
  const [activityLog, setActivityLog] = useState<ActivityLogItem[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchDriverData = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const response = await apiClient.getUserDetail(id);
        if (response.code === '0000' && response.data) {
          const data = response.data as Record<string, unknown>;
          const vehicleData = data.vehicle as Record<string, unknown> | undefined;

          const lat = data.latitude as number | undefined;
          const lng = data.longitude as number | undefined;
          if (typeof lat === 'number' && typeof lng === 'number') {
            setCurrentCoords({ lat, lng });
            setCurrentLocationText(`${lat.toFixed(5)}, ${lng.toFixed(5)}`);
          } else {
            setCurrentCoords(null);
            setCurrentLocationText((data.address as string) || 'N/A');
          }
          
          setDriver({
            id: (data.user_id as string) || id,
            name: (data.display_name as string) || 
                  `${data.first_name || ''} ${data.last_name || ''}`.trim() || 
                  'Unknown Driver',
            email: (data.email as string) || '',
            phone: (data.phone as string) || '',
            avatar: data.avatar as string | null,
            status: (data.online_status as string) || (data.status as string) || 'offline',
            rating: (data.score as number) || 0,
            totalTrips: (data.total_trips as number) || 0,
            acceptanceRate: (data.acceptance_rate as number) || 0,
            completionRate: (data.completion_rate as number) || 100,
            cancellationRate: (data.cancellation_rate as number) || 0,
            todayTrips: (data.today_trips as number) || 0,
            todayEarnings: (data.today_earnings as number) || 0,
            weekEarnings: (data.week_earnings as number) || 0,
            monthEarnings: (data.month_earnings as number) || 0,
            lifetimeEarnings: (data.total_earnings as number) || 0,
            vehicle: {
              plate: (vehicleData?.plate_number as string) || 'N/A',
              model: (vehicleData?.model as string) || 'N/A',
              make: (vehicleData?.brand as string) || 'N/A',
              year: (vehicleData?.year as number) || 0,
              color: (vehicleData?.color as string) || 'N/A',
              type: (vehicleData?.category as string) || (vehicleData?.type as string) || 'car',
            },
            documents: {
              license: { status: 'valid', expiry: 'N/A' },
              insurance: { status: 'valid', expiry: 'N/A' },
              registration: { status: 'valid', expiry: 'N/A' },
            },
            joinedAt: data.created_at 
              ? new Date(data.created_at as number).toISOString().split('T')[0]
              : 'N/A',
            address: (data.address as string) || 'N/A',
          });

          // Fetch driver's trips
          try {
            const ordersResponse = await apiClient.searchOrders({ provider_id: id, page: 1, limit: 5 });
            if (ordersResponse.code === '0000' && ordersResponse.data?.records) {
              const orders = ordersResponse.data.records as Array<Record<string, unknown>>;
              setRecentTrips(orders.map((order) => ({
                id: (order.order_id as string) || String(order.id),
                date: order.created_at 
                  ? new Date(order.created_at as number).toISOString().split('T')[0]
                  : 'N/A',
                passenger: (order.customer_name as string) || 'Customer',
                pickup: (order.pickup_address as string) || 'N/A',
                dropoff: (order.dropoff_address as string) || 'N/A',
                distance: (order.distance as number) || 0,
                fare: (order.payment_amount as number) || 0,
                rating: order.rating as number | null,
                status: (order.status as string) || 'pending',
              })));
            }
          } catch {
            // Silently fail
          }
        } else {
          setError('Driver not found');
        }
      } catch (err) {
        console.error('Failed to fetch driver:', err);
        setError('Failed to load driver details');
      } finally {
        setIsLoading(false);
      }
    };

    fetchDriverData();
  }, [id]);

  // Resolve a human-readable location from coordinates
  useEffect(() => {
    let cancelled = false;
    const run = async () => {
      if (!currentCoords) return;
      const result = await reverseGeocode(currentCoords.lat, currentCoords.lng);
      if (!cancelled && result) setCurrentLocationText(result);
    };
    run();
    return () => {
      cancelled = true;
    };
  }, [currentCoords]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  if (error || !driver) {
    return (
      <div className="space-y-6">
        <Link href="/drivers" className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground">
          <ArrowLeft className="h-4 w-4" />
          Back to Drivers
        </Link>
        <Card>
          <CardContent className="p-6 text-center text-muted-foreground">
            {error || 'Driver not found'}
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Back Button */}
      <Link
        href="/drivers"
        className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground"
      >
        <ArrowLeft className="h-4 w-4" />
        Back to Drivers
      </Link>

      <div className="grid gap-6 lg:grid-cols-3">
        {/* Profile Card */}
        <Card className="lg:col-span-1">
          <CardContent className="p-6">
            <div className="flex flex-col items-center text-center">
              <Avatar className="h-24 w-24">
                <AvatarImage src={driver.avatar || undefined} />
                <AvatarFallback className="bg-primary text-primary-foreground text-2xl">
                  {driver.name.split(' ').map((n) => n[0]).join('')}
                </AvatarFallback>
              </Avatar>
              <h2 className="mt-4 text-xl font-bold">{driver.name}</h2>
              <div className="mt-2">{getStatusBadge(driver.status)}</div>
              
              <div className="mt-4 flex items-center gap-1">
                <Star className="h-5 w-5 fill-yellow-400 text-yellow-400" />
                <span className="text-lg font-bold">{driver.rating}</span>
                <span className="text-muted-foreground">({driver.totalTrips} reviews)</span>
              </div>

              <Separator className="my-4" />

              <div className="w-full space-y-3 text-left">
                <div className="flex items-center gap-3">
                  <Phone className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm">{driver.phone}</span>
                </div>
                <div className="flex items-center gap-3">
                  <Mail className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm">{driver.email}</span>
                </div>
                <div className="flex items-center gap-3">
                  <MapPin className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm">{currentLocationText || driver.address}</span>
                </div>
                <div className="flex items-center gap-3">
                  <Calendar className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm">
                    Joined {new Date(driver.joinedAt).toLocaleDateString('en-US', {
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
                <Button variant="outline" className="w-full gap-2 text-yellow-600 hover:text-yellow-600">
                  <Ban className="h-4 w-4" />
                  Suspend Driver
                </Button>
                <Button variant="ghost" className="w-full gap-2 text-destructive hover:text-destructive">
                  <Trash2 className="h-4 w-4" />
                  Delete Driver
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Main Content */}
        <div className="lg:col-span-2 space-y-6">
          <Tabs defaultValue="overview" className="w-full">
            <TabsList className="grid w-full grid-cols-6">
              <TabsTrigger value="overview">Overview</TabsTrigger>
              <TabsTrigger value="trips">Trips</TabsTrigger>
              <TabsTrigger value="earnings">Earnings</TabsTrigger>
              <TabsTrigger value="reviews">Reviews</TabsTrigger>
              <TabsTrigger value="vehicle">Vehicle</TabsTrigger>
              <TabsTrigger value="activity">Activity</TabsTrigger>
            </TabsList>

            {/* Overview Tab */}
            <TabsContent value="overview" className="space-y-4">
              <div className="grid gap-4 md:grid-cols-4">
                <Card>
                  <CardContent className="p-4">
                    <p className="text-sm text-muted-foreground">Total Trips</p>
                    <p className="text-2xl font-bold">{driver.totalTrips.toLocaleString()}</p>
                  </CardContent>
                </Card>
                <Card>
                  <CardContent className="p-4">
                    <p className="text-sm text-muted-foreground">Acceptance Rate</p>
                    <p className="text-2xl font-bold text-green-600">{driver.acceptanceRate}%</p>
                  </CardContent>
                </Card>
                <Card>
                  <CardContent className="p-4">
                    <p className="text-sm text-muted-foreground">Completion Rate</p>
                    <p className="text-2xl font-bold text-green-600">{driver.completionRate}%</p>
                  </CardContent>
                </Card>
                <Card>
                  <CardContent className="p-4">
                    <p className="text-sm text-muted-foreground">Cancellation Rate</p>
                    <p className="text-2xl font-bold text-yellow-600">{driver.cancellationRate}%</p>
                  </CardContent>
                </Card>
              </div>

              <Card>
                <CardHeader>
                  <CardTitle>Performance (Last 30 Days)</CardTitle>
                  <CardDescription>Daily trips completed</CardDescription>
                </CardHeader>
                <CardContent>
                  <DriverPerformanceChart />
                </CardContent>
              </Card>
            </TabsContent>

            {/* Trips Tab */}
            <TabsContent value="trips">
              <Card>
                <CardHeader>
                  <CardTitle>Trip History</CardTitle>
                </CardHeader>
                <CardContent>
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Date</TableHead>
                        <TableHead>Passenger</TableHead>
                        <TableHead>Route</TableHead>
                        <TableHead className="text-right">Distance</TableHead>
                        <TableHead className="text-right">Fare</TableHead>
                        <TableHead>Rating</TableHead>
                        <TableHead>Status</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {recentTrips.map((trip) => (
                        <TableRow key={trip.id}>
                          <TableCell>{trip.date}</TableCell>
                          <TableCell>{trip.passenger}</TableCell>
                          <TableCell>
                            <span className="text-muted-foreground">{trip.pickup}</span>
                            <span className="mx-1">→</span>
                            <span>{trip.dropoff}</span>
                          </TableCell>
                          <TableCell className="text-right">{trip.distance} km</TableCell>
                          <TableCell className="text-right font-medium">
                            RWF {trip.fare.toLocaleString()}
                          </TableCell>
                          <TableCell>
                            {trip.rating ? (
                              <div className="flex items-center gap-1">
                                <Star className="h-4 w-4 fill-yellow-400 text-yellow-400" />
                                {trip.rating}
                              </div>
                            ) : (
                              '-'
                            )}
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

            {/* Earnings Tab */}
            <TabsContent value="earnings" className="space-y-4">
              <div className="grid gap-4 md:grid-cols-4">
                <Card>
                  <CardContent className="p-4">
                    <p className="text-sm text-muted-foreground">Today</p>
                    <p className="text-2xl font-bold">RWF {driver.todayEarnings.toLocaleString()}</p>
                  </CardContent>
                </Card>
                <Card>
                  <CardContent className="p-4">
                    <p className="text-sm text-muted-foreground">This Week</p>
                    <p className="text-2xl font-bold">RWF {driver.weekEarnings.toLocaleString()}</p>
                  </CardContent>
                </Card>
                <Card>
                  <CardContent className="p-4">
                    <p className="text-sm text-muted-foreground">This Month</p>
                    <p className="text-2xl font-bold">RWF {driver.monthEarnings.toLocaleString()}</p>
                  </CardContent>
                </Card>
                <Card>
                  <CardContent className="p-4">
                    <p className="text-sm text-muted-foreground">Lifetime</p>
                    <p className="text-2xl font-bold">RWF {(driver.lifetimeEarnings / 1000000).toFixed(1)}M</p>
                  </CardContent>
                </Card>
              </div>
            </TabsContent>

            {/* Reviews Tab */}
            <TabsContent value="reviews">
              <Card>
                <CardHeader>
                  <CardTitle>Passenger Reviews</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  {reviews.map((review) => (
                    <div key={review.id} className="flex gap-4 border-b pb-4 last:border-0">
                      <Avatar className="h-10 w-10">
                        <AvatarFallback>{review.passenger[0]}</AvatarFallback>
                      </Avatar>
                      <div className="flex-1">
                        <div className="flex items-center justify-between">
                          <p className="font-medium">{review.passenger}</p>
                          <div className="flex items-center gap-1">
                            {[...Array(5)].map((_, i) => (
                              <Star
                                key={i}
                                className={`h-4 w-4 ${
                                  i < review.rating
                                    ? 'fill-yellow-400 text-yellow-400'
                                    : 'text-gray-200'
                                }`}
                              />
                            ))}
                          </div>
                        </div>
                        <p className="mt-1 text-sm text-muted-foreground">{review.comment}</p>
                        <p className="mt-1 text-xs text-muted-foreground">{review.date}</p>
                      </div>
                    </div>
                  ))}
                </CardContent>
              </Card>
            </TabsContent>

            {/* Vehicle Tab */}
            <TabsContent value="vehicle">
              <Card>
                <CardHeader>
                  <CardTitle>Vehicle Information</CardTitle>
                </CardHeader>
                <CardContent className="space-y-6">
                  <div className="grid gap-4 md:grid-cols-3">
                    <div>
                      <p className="text-sm text-muted-foreground">Plate Number</p>
                      <p className="text-lg font-medium">{driver.vehicle.plate}</p>
                    </div>
                    <div>
                      <p className="text-sm text-muted-foreground">Model</p>
                      <p className="text-lg font-medium">{driver.vehicle.model}</p>
                    </div>
                    <div>
                      <p className="text-sm text-muted-foreground">Make</p>
                      <p className="text-lg font-medium">{driver.vehicle.make}</p>
                    </div>
                    <div>
                      <p className="text-sm text-muted-foreground">Year</p>
                      <p className="text-lg font-medium">{driver.vehicle.year}</p>
                    </div>
                    <div>
                      <p className="text-sm text-muted-foreground">Color</p>
                      <p className="text-lg font-medium">{driver.vehicle.color}</p>
                    </div>
                    <div>
                      <p className="text-sm text-muted-foreground">Type</p>
                      <p className="text-lg font-medium">{driver.vehicle.type}</p>
                    </div>
                  </div>

                  <Separator />

                  <div>
                    <h4 className="font-medium mb-4">Documents</h4>
                    <div className="grid gap-4 md:grid-cols-3">
                      <div className="flex items-center justify-between p-3 border rounded-lg">
                        <div>
                          <p className="font-medium">License</p>
                          <p className="text-sm text-muted-foreground">
                            Expires: {driver.documents.license.expiry}
                          </p>
                        </div>
                        <CheckCircle className="h-5 w-5 text-green-500" />
                      </div>
                      <div className="flex items-center justify-between p-3 border rounded-lg">
                        <div>
                          <p className="font-medium">Insurance</p>
                          <p className="text-sm text-muted-foreground">
                            Expires: {driver.documents.insurance.expiry}
                          </p>
                        </div>
                        <CheckCircle className="h-5 w-5 text-green-500" />
                      </div>
                      <div className="flex items-center justify-between p-3 border rounded-lg">
                        <div>
                          <p className="font-medium">Registration</p>
                          <p className="text-sm text-muted-foreground">
                            Expires: {driver.documents.registration.expiry}
                          </p>
                        </div>
                        <CheckCircle className="h-5 w-5 text-green-500" />
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            {/* Activity Tab */}
            <TabsContent value="activity">
              <Card>
                <CardHeader>
                  <CardTitle>Activity Log</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {activityLog.map((activity) => (
                      <div key={activity.id} className="flex items-center gap-4">
                        <div className="h-2 w-2 rounded-full bg-primary" />
                        <div className="flex-1">
                          <p className="text-sm font-medium">{activity.event}</p>
                          <p className="text-xs text-muted-foreground">{activity.time}</p>
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </div>
  );
}
