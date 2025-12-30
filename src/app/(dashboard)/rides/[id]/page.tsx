'use client';

import { use } from 'react';
import Link from 'next/link';
import { 
  ArrowLeft, 
  MapPin, 
  Clock,
  DollarSign,
  Car,
  User,
  Star,
  Phone,
  Navigation,
  XCircle
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Separator } from '@/components/ui/separator';

// Mock ride data
const mockRide = {
  id: 'R001',
  status: 'completed',
  passenger: {
    id: '1',
    name: 'John Doe',
    phone: '+250 788 111 222',
    rating: 4.8,
  },
  driver: {
    id: '1',
    name: 'Peter Mugisha',
    phone: '+250 788 123 456',
    rating: 4.9,
    vehicle: 'Toyota Corolla • RAD 123A',
  },
  pickup: {
    address: 'Kimironko Market, Kigali',
    lat: -1.9403,
    lng: 29.8739,
    time: '2024-12-28 14:30',
  },
  dropoff: {
    address: 'Downtown Kigali, City Center',
    lat: -1.9453,
    lng: 29.8789,
    time: '2024-12-28 14:48',
  },
  distance: 5.2,
  duration: 18,
  fare: {
    baseFare: 1000,
    distanceFare: 2600,
    timeFare: 360,
    total: 4500,
    discount: 0,
    final: 4500,
  },
  paymentMethod: 'momo',
  paymentStatus: 'paid',
  rating: 5,
  review: 'Great driver, very professional!',
  createdAt: '2024-12-28 14:25',
  completedAt: '2024-12-28 14:48',
};

const timeline = [
  { time: '14:25', event: 'Ride requested', status: 'completed' },
  { time: '14:26', event: 'Driver assigned - Peter Mugisha', status: 'completed' },
  { time: '14:28', event: 'Driver arriving', status: 'completed' },
  { time: '14:30', event: 'Driver arrived at pickup', status: 'completed' },
  { time: '14:30', event: 'Trip started', status: 'completed' },
  { time: '14:48', event: 'Trip ended', status: 'completed' },
  { time: '14:48', event: 'Payment received via MoMo', status: 'completed' },
];

const getStatusBadge = (status: string) => {
  switch (status) {
    case 'completed':
      return <Badge className="bg-green-100 text-green-700 hover:bg-green-100">Completed</Badge>;
    case 'in_progress':
      return <Badge className="bg-blue-100 text-blue-700 hover:bg-blue-100">In Progress</Badge>;
    case 'cancelled':
      return <Badge className="bg-red-100 text-red-700 hover:bg-red-100">Cancelled</Badge>;
    default:
      return <Badge variant="secondary">{status}</Badge>;
  }
};

export default function RideDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const ride = mockRide;

  return (
    <div className="space-y-6">
      {/* Back Button */}
      <Link
        href="/rides"
        className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground"
      >
        <ArrowLeft className="h-4 w-4" />
        Back to Rides
      </Link>

      {/* Header */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <div className="flex items-center gap-3">
            <h1 className="text-2xl font-bold tracking-tight">Ride {ride.id}</h1>
            {getStatusBadge(ride.status)}
          </div>
          <p className="text-muted-foreground">
            {new Date(ride.createdAt).toLocaleString()}
          </p>
        </div>
        {ride.status !== 'completed' && ride.status !== 'cancelled' && (
          <Button variant="destructive" className="gap-2">
            <XCircle className="h-4 w-4" />
            Cancel Ride
          </Button>
        )}
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-6">
          {/* Map Placeholder */}
          <Card>
            <CardContent className="p-0">
              <div 
                className="h-64 rounded-t-lg"
                style={{
                  background: 'linear-gradient(135deg, #e0f2e9 0%, #c6e2d5 50%, #a8d4be 100%)',
                }}
              >
                <div className="h-full w-full flex items-center justify-center">
                  <div className="text-center">
                    <Navigation className="h-12 w-12 text-primary mx-auto mb-2" />
                    <p className="text-sm text-muted-foreground">Route Map</p>
                    <p className="text-xs text-muted-foreground">{ride.distance} km • {ride.duration} min</p>
                  </div>
                </div>
              </div>
              <div className="p-4 space-y-4">
                <div className="flex items-start gap-3">
                  <div className="flex flex-col items-center">
                    <div className="h-3 w-3 rounded-full bg-green-500" />
                    <div className="w-0.5 h-8 bg-gray-200" />
                    <div className="h-3 w-3 rounded-full bg-red-500" />
                  </div>
                  <div className="flex-1 space-y-4">
                    <div>
                      <p className="font-medium">{ride.pickup.address}</p>
                      <p className="text-sm text-muted-foreground">{ride.pickup.time}</p>
                    </div>
                    <div>
                      <p className="font-medium">{ride.dropoff.address}</p>
                      <p className="text-sm text-muted-foreground">{ride.dropoff.time}</p>
                    </div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Fare Breakdown */}
          <Card>
            <CardHeader>
              <CardTitle>Fare Breakdown</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Base Fare</span>
                <span>RWF {ride.fare.baseFare.toLocaleString()}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Distance ({ride.distance} km)</span>
                <span>RWF {ride.fare.distanceFare.toLocaleString()}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Time ({ride.duration} min)</span>
                <span>RWF {ride.fare.timeFare.toLocaleString()}</span>
              </div>
              {ride.fare.discount > 0 && (
                <div className="flex justify-between text-green-600">
                  <span>Discount</span>
                  <span>-RWF {ride.fare.discount.toLocaleString()}</span>
                </div>
              )}
              <Separator />
              <div className="flex justify-between font-bold text-lg">
                <span>Total</span>
                <span>RWF {ride.fare.final.toLocaleString()}</span>
              </div>
              <div className="flex items-center gap-2 mt-4">
                <Badge variant="outline" className="capitalize">
                  {ride.paymentMethod}
                </Badge>
                <Badge className="bg-green-100 text-green-700">
                  {ride.paymentStatus}
                </Badge>
              </div>
            </CardContent>
          </Card>

          {/* Rating & Review */}
          {ride.rating && (
            <Card>
              <CardHeader>
                <CardTitle>Rating & Review</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex items-center gap-2 mb-2">
                  {[...Array(5)].map((_, i) => (
                    <Star
                      key={i}
                      className={`h-5 w-5 ${
                        i < ride.rating!
                          ? 'fill-yellow-400 text-yellow-400'
                          : 'text-gray-200'
                      }`}
                    />
                  ))}
                </div>
                {ride.review && (
                  <p className="text-muted-foreground italic">&quot;{ride.review}&quot;</p>
                )}
              </CardContent>
            </Card>
          )}
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Passenger Info */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <User className="h-4 w-4" />
                Passenger
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex items-center gap-3">
                <Avatar className="h-12 w-12">
                  <AvatarFallback className="bg-blue-100 text-blue-600">
                    {ride.passenger.name.split(' ').map(n => n[0]).join('')}
                  </AvatarFallback>
                </Avatar>
                <div>
                  <p className="font-medium">{ride.passenger.name}</p>
                  <div className="flex items-center gap-1 text-sm text-muted-foreground">
                    <Star className="h-3 w-3 fill-yellow-400 text-yellow-400" />
                    {ride.passenger.rating}
                  </div>
                </div>
              </div>
              <Button variant="outline" size="sm" className="w-full mt-4 gap-2">
                <Phone className="h-4 w-4" />
                {ride.passenger.phone}
              </Button>
            </CardContent>
          </Card>

          {/* Driver Info */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Car className="h-4 w-4" />
                Driver
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex items-center gap-3">
                <Avatar className="h-12 w-12">
                  <AvatarFallback className="bg-primary/10 text-primary">
                    {ride.driver.name.split(' ').map(n => n[0]).join('')}
                  </AvatarFallback>
                </Avatar>
                <div>
                  <p className="font-medium">{ride.driver.name}</p>
                  <div className="flex items-center gap-1 text-sm text-muted-foreground">
                    <Star className="h-3 w-3 fill-yellow-400 text-yellow-400" />
                    {ride.driver.rating}
                  </div>
                </div>
              </div>
              <p className="text-sm text-muted-foreground mt-2">{ride.driver.vehicle}</p>
              <Button variant="outline" size="sm" className="w-full mt-4 gap-2">
                <Phone className="h-4 w-4" />
                {ride.driver.phone}
              </Button>
            </CardContent>
          </Card>

          {/* Timeline */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Clock className="h-4 w-4" />
                Timeline
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {timeline.map((item, idx) => (
                  <div key={idx} className="flex gap-3">
                    <div className="flex flex-col items-center">
                      <div className="h-2 w-2 rounded-full bg-green-500" />
                      {idx < timeline.length - 1 && (
                        <div className="w-0.5 flex-1 bg-gray-200" />
                      )}
                    </div>
                    <div className="pb-4">
                      <p className="text-sm font-medium">{item.event}</p>
                      <p className="text-xs text-muted-foreground">{item.time}</p>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
