'use client';

import { use } from 'react';
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
  Clock
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

// Mock driver data
const mockDriver = {
  id: '1',
  name: 'Peter Mugisha',
  email: 'peter.m@email.com',
  phone: '+250 788 123 456',
  avatar: null,
  status: 'online',
  rating: 4.8,
  totalTrips: 1234,
  acceptanceRate: 92,
  completionRate: 98,
  cancellationRate: 2,
  todayTrips: 8,
  todayEarnings: 45000,
  weekEarnings: 285000,
  monthEarnings: 1150000,
  lifetimeEarnings: 4800000,
  vehicle: {
    plate: 'RAD 123A',
    model: 'Toyota Corolla',
    make: 'Toyota',
    year: 2019,
    color: 'Silver',
    type: 'Sedan',
  },
  documents: {
    license: { status: 'valid', expiry: '2025-06-15' },
    insurance: { status: 'valid', expiry: '2025-03-20' },
    registration: { status: 'valid', expiry: '2025-12-01' },
  },
  joinedAt: '2024-06-15',
  address: 'Kigali, Gasabo District',
};

const recentTrips = [
  { id: 'T001', date: '2024-12-28', passenger: 'John Doe', pickup: 'Kimironko', dropoff: 'Downtown', distance: 5.2, fare: 4500, rating: 5, status: 'completed' },
  { id: 'T002', date: '2024-12-28', passenger: 'Jane Smith', pickup: 'Remera', dropoff: 'Nyarutarama', distance: 3.8, fare: 3200, rating: 4, status: 'completed' },
  { id: 'T003', date: '2024-12-28', passenger: 'Mike Johnson', pickup: 'Kicukiro', dropoff: 'Gisozi', distance: 8.1, fare: 6800, rating: 5, status: 'completed' },
  { id: 'T004', date: '2024-12-27', passenger: 'Sarah Wilson', pickup: 'Downtown', dropoff: 'Kimihurura', distance: 2.5, fare: 2500, rating: null, status: 'cancelled' },
  { id: 'T005', date: '2024-12-27', passenger: 'Chris Brown', pickup: 'Kacyiru', dropoff: 'Kibagabaga', distance: 4.0, fare: 3600, rating: 5, status: 'completed' },
];

const reviews = [
  { id: 1, passenger: 'John Doe', rating: 5, comment: 'Excellent driver! Very professional and punctual.', date: '2024-12-28' },
  { id: 2, passenger: 'Jane Smith', rating: 4, comment: 'Good service, car was clean.', date: '2024-12-28' },
  { id: 3, passenger: 'Mike Johnson', rating: 5, comment: 'Best driver I have had so far. Highly recommended!', date: '2024-12-28' },
  { id: 4, passenger: 'Chris Brown', rating: 5, comment: 'Very friendly and knows the city well.', date: '2024-12-27' },
];

const activityLog = [
  { id: 1, event: 'Completed trip #T003', time: '2 hours ago' },
  { id: 2, event: 'Started trip #T003', time: '2 hours ago' },
  { id: 3, event: 'Completed trip #T002', time: '3 hours ago' },
  { id: 4, event: 'Started trip #T002', time: '3 hours ago' },
  { id: 5, event: 'Went online', time: '4 hours ago' },
  { id: 6, event: 'Completed trip #T001', time: 'Yesterday' },
];

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
  const driver = mockDriver; // In real app, fetch by id

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
                  <span className="text-sm">{driver.address}</span>
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
