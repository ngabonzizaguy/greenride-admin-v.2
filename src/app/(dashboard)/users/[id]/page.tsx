'use client';

import { use } from 'react';
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
  Star
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

// Mock user data
const mockUser = {
  id: '1',
  name: 'John Doe',
  email: 'john.doe@email.com',
  phone: '+250 788 111 222',
  avatar: null,
  status: 'active',
  totalTrips: 45,
  totalSpent: 185000,
  avgRating: 4.8,
  lastTrip: '2024-12-28',
  joinedAt: '2024-03-15',
  address: 'Kimironko, Kigali',
};

const recentTrips = [
  { id: 'T001', date: '2024-12-28', driver: 'Peter M.', pickup: 'Kimironko', dropoff: 'Downtown', fare: 4500, rating: 5, status: 'completed' },
  { id: 'T002', date: '2024-12-27', driver: 'David K.', pickup: 'Remera', dropoff: 'Nyarutarama', fare: 3200, rating: 4, status: 'completed' },
  { id: 'T003', date: '2024-12-26', driver: 'Jean P.', pickup: 'Kicukiro', dropoff: 'Gisozi', fare: 6800, rating: 5, status: 'completed' },
  { id: 'T004', date: '2024-12-25', driver: 'Claude U.', pickup: 'Downtown', dropoff: 'Kimihurura', fare: 2500, rating: null, status: 'cancelled' },
];

const paymentMethods = [
  { type: 'MoMo', number: '**** 1234', isDefault: true },
  { type: 'Card', number: '**** 5678', isDefault: false },
];

export default function UserDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const user = mockUser;

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
                  {new Date(user.lastTrip).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}
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
