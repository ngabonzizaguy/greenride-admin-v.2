'use client';

import { use, useState, useEffect } from 'react';
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
  XCircle,
  Loader2
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Separator } from '@/components/ui/separator';
import { apiClient } from '@/lib/api-client';
import { toast } from 'sonner';

interface RideData {
  order_id: string;
  status: string;
  passenger?: {
    user_id?: string;
    display_name?: string;
    first_name?: string;
    last_name?: string;
    phone?: string;
    email?: string;
    avatar?: string;
    score?: number;
  };
  driver?: {
    user_id?: string;
    display_name?: string;
    first_name?: string;
    last_name?: string;
    phone?: string;
    email?: string;
    avatar?: string;
    score?: number;
  };
  vehicle?: {
    plate_number?: string;
    brand?: string;
    model?: string;
    color?: string;
  };
  pickup_address?: string;
  pickup_latitude?: number;
  pickup_longitude?: number;
  dropoff_address?: string;
  dropoff_latitude?: number;
  dropoff_longitude?: number;
  estimated_distance?: number;
  actual_distance?: number;
  estimated_duration?: number;
  actual_duration?: number;
  payment_amount?: number;
  payment_method?: string;
  payment_status?: string;
  base_fare?: number;
  distance_fare?: number;
  time_fare?: number;
  total_fare?: number;
  created_at?: number;
  started_at?: number;
  completed_at?: number;
  cancelled_at?: number;
  cancelled_by?: string;
  cancel_reason?: string;
  driver_ratings?: Array<{ rating?: number; comment?: string }>;
}

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
  const [ride, setRide] = useState<RideData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchRideData = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const response = await apiClient.getOrderDetail(id);
        if (response.code === '0000' && response.data) {
          const data = response.data as Record<string, unknown>;
          setRide({
            order_id: (data.order_id as string) || id,
            status: (data.status as string) || 'pending',
            passenger: data.passenger as RideData['passenger'],
            driver: data.driver as RideData['driver'],
            vehicle: data.vehicle as RideData['vehicle'],
            pickup_address: data.pickup_address as string,
            pickup_latitude: data.pickup_latitude as number,
            pickup_longitude: data.pickup_longitude as number,
            dropoff_address: data.dropoff_address as string,
            dropoff_latitude: data.dropoff_latitude as number,
            dropoff_longitude: data.dropoff_longitude as number,
            estimated_distance: data.estimated_distance as number,
            actual_distance: data.actual_distance as number,
            estimated_duration: data.estimated_duration as number,
            actual_duration: data.actual_duration as number,
            payment_amount: data.payment_amount as number,
            payment_method: data.payment_method as string,
            payment_status: data.payment_status as string,
            base_fare: data.base_fare as number,
            distance_fare: data.distance_fare as number,
            time_fare: data.time_fare as number,
            total_fare: data.total_fare as number,
            created_at: data.created_at as number,
            started_at: data.started_at as number,
            completed_at: data.completed_at as number,
            cancelled_at: data.cancelled_at as number,
            cancelled_by: data.cancelled_by as string,
            cancel_reason: data.cancel_reason as string,
            driver_ratings: data.driver_ratings as Array<{ rating?: number; comment?: string }>,
          });
        } else {
          setError('Ride not found');
        }
      } catch (err) {
        console.error('Failed to fetch ride:', err);
        setError('Failed to load ride details');
        toast.error('Failed to load ride details');
      } finally {
        setIsLoading(false);
      }
    };

    fetchRideData();
  }, [id]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  if (error || !ride) {
    return (
      <div className="space-y-6">
        <Link
          href="/rides"
          className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground"
        >
          <ArrowLeft className="h-4 w-4" />
          Back to Rides
        </Link>
        <Card>
          <CardContent className="p-12 text-center">
            <p className="text-muted-foreground">{error || 'Ride not found'}</p>
          </CardContent>
        </Card>
      </div>
    );
  }

  // Build timeline from order timestamps
  const timeline: Array<{ time: string; event: string; status: string }> = [];
  if (ride.created_at) {
    timeline.push({
      time: new Date(ride.created_at).toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' }),
      event: 'Ride requested',
      status: 'completed',
    });
  }
  if (ride.started_at) {
    timeline.push({
      time: new Date(ride.started_at).toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' }),
      event: 'Trip started',
      status: 'completed',
    });
  }
  if (ride.completed_at) {
    timeline.push({
      time: new Date(ride.completed_at).toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' }),
      event: 'Trip ended',
      status: 'completed',
    });
  }
  if (ride.payment_status === 'paid' && ride.completed_at) {
    timeline.push({
      time: new Date(ride.completed_at).toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' }),
      event: `Payment received via ${ride.payment_method || 'payment'}`,
      status: 'completed',
    });
  }
  if (ride.cancelled_at) {
    timeline.push({
      time: new Date(ride.cancelled_at).toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' }),
      event: 'Ride cancelled',
      status: 'cancelled',
    });
  }

  const passengerName = ride.passenger?.display_name || 
    `${ride.passenger?.first_name || ''} ${ride.passenger?.last_name || ''}`.trim() || 
    'Passenger';
  const driverName = ride.driver?.display_name || 
    `${ride.driver?.first_name || ''} ${ride.driver?.last_name || ''}`.trim() || 
    'Driver';
  const vehicleInfo = ride.vehicle 
    ? `${ride.vehicle.brand || ''} ${ride.vehicle.model || ''} • ${ride.vehicle.plate_number || 'N/A'}`.trim()
    : 'N/A';
  const distance = ride.actual_distance || ride.estimated_distance || 0;
  const duration = ride.actual_duration || ride.estimated_duration || 0;
  const rating = ride.driver_ratings?.[0]?.rating || 0;
  const review = ride.driver_ratings?.[0]?.comment || '';

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
            <h1 className="text-2xl font-bold tracking-tight">Ride {ride.order_id}</h1>
            {getStatusBadge(ride.status)}
          </div>
          <p className="text-muted-foreground">
            {ride.created_at ? new Date(ride.created_at).toLocaleString() : 'N/A'}
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
                    <p className="text-xs text-muted-foreground">{distance.toFixed(1)} km • {duration} min</p>
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
                      <p className="font-medium">{ride.pickup_address || 'N/A'}</p>
                      {ride.created_at && (
                        <p className="text-sm text-muted-foreground">
                          {new Date(ride.created_at).toLocaleString()}
                        </p>
                      )}
                    </div>
                    <div>
                      <p className="font-medium">{ride.dropoff_address || 'N/A'}</p>
                      {ride.completed_at && (
                        <p className="text-sm text-muted-foreground">
                          {new Date(ride.completed_at).toLocaleString()}
                        </p>
                      )}
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
              {ride.base_fare !== undefined && (
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Base Fare</span>
                  <span>RWF {(ride.base_fare || 0).toLocaleString()}</span>
                </div>
              )}
              {ride.distance_fare !== undefined && (
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Distance ({distance.toFixed(1)} km)</span>
                  <span>RWF {(ride.distance_fare || 0).toLocaleString()}</span>
                </div>
              )}
              {ride.time_fare !== undefined && (
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Time ({duration} min)</span>
                  <span>RWF {(ride.time_fare || 0).toLocaleString()}</span>
                </div>
              )}
              <Separator />
              <div className="flex justify-between font-bold text-lg">
                <span>Total</span>
                <span>RWF {(ride.payment_amount || ride.total_fare || 0).toLocaleString()}</span>
              </div>
              <div className="flex items-center gap-2 mt-4">
                {ride.payment_method && (
                  <Badge variant="outline" className="capitalize">
                    {ride.payment_method}
                  </Badge>
                )}
                {ride.payment_status && (
                  <Badge className={ride.payment_status === 'paid' ? 'bg-green-100 text-green-700' : 'bg-yellow-100 text-yellow-700'}>
                    {ride.payment_status}
                  </Badge>
                )}
              </div>
            </CardContent>
          </Card>

          {/* Rating & Review */}
          {rating > 0 && (
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
                        i < rating
                          ? 'fill-yellow-400 text-yellow-400'
                          : 'text-gray-200'
                      }`}
                    />
                  ))}
                </div>
                {review && (
                  <p className="text-muted-foreground italic">&quot;{review}&quot;</p>
                )}
              </CardContent>
            </Card>
          )}

          {/* Cancellation Info */}
          {ride.status === 'cancelled' && (
            <Card className="border-red-200">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-red-700">
                  <XCircle className="h-4 w-4" />
                  Cancellation Details
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                {ride.cancel_reason && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Reason</span>
                    <span className="font-medium">{ride.cancel_reason}</span>
                  </div>
                )}
                {ride.cancelled_by && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Cancelled by</span>
                    <span className="font-medium">
                      {ride.cancelled_by === ride.passenger?.user_id
                        ? 'Passenger'
                        : ride.cancelled_by === ride.driver?.user_id
                        ? 'Driver'
                        : 'System'}
                    </span>
                  </div>
                )}
                {ride.cancelled_at && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Cancelled at</span>
                    <span>{new Date(ride.cancelled_at).toLocaleString()}</span>
                  </div>
                )}
              </CardContent>
            </Card>
          )}
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Passenger Info */}
          {ride.passenger && (
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
                    {ride.passenger.avatar ? (
                      <img src={ride.passenger.avatar} alt={passengerName} />
                    ) : (
                      <AvatarFallback className="bg-blue-100 text-blue-600">
                        {passengerName.split(' ').map(n => n[0]).join('').slice(0, 2).toUpperCase()}
                      </AvatarFallback>
                    )}
                  </Avatar>
                  <div>
                    <p className="font-medium">{passengerName}</p>
                    {ride.passenger.score !== undefined && (
                      <div className="flex items-center gap-1 text-sm text-muted-foreground">
                        <Star className="h-3 w-3 fill-yellow-400 text-yellow-400" />
                        {ride.passenger.score.toFixed(1)}
                      </div>
                    )}
                  </div>
                </div>
                {ride.passenger.user_id && (
                  <Link href={`/users/${ride.passenger.user_id}`}>
                    <Button variant="outline" size="sm" className="w-full mt-4 gap-2">
                      <User className="h-4 w-4" />
                      View Profile
                    </Button>
                  </Link>
                )}
                {ride.passenger.phone && (
                  <Button variant="outline" size="sm" className="w-full mt-2 gap-2">
                    <Phone className="h-4 w-4" />
                    {ride.passenger.phone}
                  </Button>
                )}
              </CardContent>
            </Card>
          )}

          {/* Driver Info */}
          {ride.driver && (
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
                    {ride.driver.avatar ? (
                      <img src={ride.driver.avatar} alt={driverName} />
                    ) : (
                      <AvatarFallback className="bg-primary/10 text-primary">
                        {driverName.split(' ').map(n => n[0]).join('').slice(0, 2).toUpperCase()}
                      </AvatarFallback>
                    )}
                  </Avatar>
                  <div>
                    <p className="font-medium">{driverName}</p>
                    {ride.driver.score !== undefined && (
                      <div className="flex items-center gap-1 text-sm text-muted-foreground">
                        <Star className="h-3 w-3 fill-yellow-400 text-yellow-400" />
                        {ride.driver.score.toFixed(1)}
                      </div>
                    )}
                  </div>
                </div>
                {vehicleInfo && vehicleInfo !== 'N/A' && (
                  <p className="text-sm text-muted-foreground mt-2">{vehicleInfo}</p>
                )}
                {ride.driver.user_id && (
                  <Link href={`/drivers/${ride.driver.user_id}`}>
                    <Button variant="outline" size="sm" className="w-full mt-4 gap-2">
                      <User className="h-4 w-4" />
                      View Profile
                    </Button>
                  </Link>
                )}
                {ride.driver.phone && (
                  <Button variant="outline" size="sm" className="w-full mt-2 gap-2">
                    <Phone className="h-4 w-4" />
                    {ride.driver.phone}
                  </Button>
                )}
              </CardContent>
            </Card>
          )}

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
