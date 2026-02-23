'use client';

import { useState, useEffect, useCallback, useMemo } from 'react';
import { GoogleMap, Marker, useJsApiLoader } from '@react-google-maps/api';
import { 
  Phone, 
  User, 
  MapPin, 
  Car, 
  Search, 
  CheckCircle, 
  AlertCircle,
  Navigation,
  Star,
  Clock,
  X,
  UserPlus,
  RefreshCw
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import { PlacesAutocompleteInput, PlaceSelection } from '@/components/booking/places-autocomplete-input';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Skeleton } from '@/components/ui/skeleton';
import { apiClient, ApiError, NearbyDriverLocation } from '@/lib/api-client';
import { geocodeAddress, POPULAR_LOCATIONS } from '@/lib/geocoding';
import type { User as UserType, PageResult } from '@/types';

type BookingStep = 'passenger' | 'locations' | 'confirm';

export default function QuickBookingPage() {
  const { isLoaded: isGoogleLoaded } = useJsApiLoader({
    id: 'google-map-script',
    googleMapsApiKey: process.env.NEXT_PUBLIC_GOOGLE_MAPS_KEY || '',
    libraries: ['places'],
  });

  // Step tracking
  const [currentStep, setCurrentStep] = useState<BookingStep>('passenger');
  
  // Passenger info
  const [passengerPhone, setPassengerPhone] = useState('');
  const [isSearchingPassenger, setIsSearchingPassenger] = useState(false);
  const [foundPassenger, setFoundPassenger] = useState<UserType | null>(null);
  const [isNewPassenger, setIsNewPassenger] = useState(false);
  const [newPassengerForm, setNewPassengerForm] = useState({
    first_name: '',
    last_name: '',
    phone: '',
    email: '',
  });

  // Location info
  const [pickupLocation, setPickupLocation] = useState('');
  const [dropoffLocation, setDropoffLocation] = useState('');
  const [pickupCoords, setPickupCoords] = useState<{lat: number, lng: number} | null>(null);
  const [dropoffCoords, setDropoffCoords] = useState<{lat: number, lng: number} | null>(null);
  const [isGeocodingPickup, setIsGeocodingPickup] = useState(false);
  const [isGeocodingDropoff, setIsGeocodingDropoff] = useState(false);
  const [pickupFromPlaces, setPickupFromPlaces] = useState(false);
  const [dropoffFromPlaces, setDropoffFromPlaces] = useState(false);
  
  // Booking state
  const [isBooking, setIsBooking] = useState(false);
  const [bookingComplete, setBookingComplete] = useState(false);
  const [bookingDetails, setBookingDetails] = useState<{
    orderId: string;
    passenger: string;
    driver: string;
    vehicle?: string;
    plate?: string;
    driverStatus?: string;
    driverLocation?: string;
    pickup: string;
    dropoff: string;
    eta: string;
  } | null>(null);
  
  // Messages
  const [error, setError] = useState<string | null>(null);
  const [bookingErrorCode, setBookingErrorCode] = useState<string | null>(null); // e.g. 6007 = ride in progress
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [estimatedFare, setEstimatedFare] = useState<number | undefined>(undefined);
  
  // Create new passenger modal
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isCreatingPassenger, setIsCreatingPassenger] = useState(false);

  // Driver selection (Phase 5)
  const [availableDrivers, setAvailableDrivers] = useState<NearbyDriverLocation[]>([]);
  const [selectedDriverId, setSelectedDriverId] = useState<string>('auto');
  const [isLoadingDrivers, setIsLoadingDrivers] = useState(false);
  const [driversError, setDriversError] = useState<string | null>(null);

  const selectedDriver = useMemo(
    () => availableDrivers.find(d => d.driver_id === selectedDriverId) || null,
    [availableDrivers, selectedDriverId]
  );

  // Auto-hide messages
  useEffect(() => {
    if (successMessage) {
      const timer = setTimeout(() => setSuccessMessage(null), 8000);
      return () => clearTimeout(timer);
    }
  }, [successMessage]);

  useEffect(() => {
    if (error) {
      const timer = setTimeout(() => setError(null), 5000);
      return () => clearTimeout(timer);
    }
  }, [error]);

  // Geocode pickup location (debounced)
  useEffect(() => {
    if (pickupFromPlaces) return;
    if (!pickupLocation || pickupLocation.trim().length === 0) {
      setPickupCoords(null);
      return;
    }

    // Check if it's a popular location first
    const popularLocation = Object.values(POPULAR_LOCATIONS).find(
      loc => loc.address.toLowerCase().includes(pickupLocation.toLowerCase()) ||
             pickupLocation.toLowerCase().includes(loc.address.toLowerCase().split(',')[0].toLowerCase())
    );
    
    if (popularLocation) {
      setPickupCoords({ lat: popularLocation.lat, lng: popularLocation.lng });
      return;
    }

    // Debounce geocoding
    const timer = setTimeout(async () => {
      setIsGeocodingPickup(true);
      try {
        const result = await geocodeAddress(pickupLocation);
        if (result) {
          setPickupCoords({ lat: result.lat, lng: result.lng });
        } else {
          setPickupCoords(null);
        }
      } catch (err) {
        console.error('Geocoding failed for pickup:', err);
        setPickupCoords(null);
      } finally {
        setIsGeocodingPickup(false);
      }
    }, 500);

    return () => clearTimeout(timer);
  }, [pickupLocation]);

  // Geocode dropoff location (debounced)
  useEffect(() => {
    if (dropoffFromPlaces) return;
    if (!dropoffLocation || dropoffLocation.trim().length === 0) {
      setDropoffCoords(null);
      return;
    }

    // Check if it's a popular location first
    const popularLocation = Object.values(POPULAR_LOCATIONS).find(
      loc => loc.address.toLowerCase().includes(dropoffLocation.toLowerCase()) ||
             dropoffLocation.toLowerCase().includes(loc.address.toLowerCase().split(',')[0].toLowerCase())
    );
    
    if (popularLocation) {
      setDropoffCoords({ lat: popularLocation.lat, lng: popularLocation.lng });
      return;
    }

    // Debounce geocoding
    const timer = setTimeout(async () => {
      setIsGeocodingDropoff(true);
      try {
        const result = await geocodeAddress(dropoffLocation);
        if (result) {
          setDropoffCoords({ lat: result.lat, lng: result.lng });
        } else {
          setDropoffCoords(null);
        }
      } catch (err) {
        console.error('Geocoding failed for dropoff:', err);
        setDropoffCoords(null);
      } finally {
        setIsGeocodingDropoff(false);
      }
    }, 500);

    return () => clearTimeout(timer);
  }, [dropoffLocation]);

  // Search for existing passenger by phone
  const searchPassenger = async () => {
    if (!passengerPhone || passengerPhone.length < 9) {
      setError('Please enter a valid phone number');
      return;
    }

    setIsSearchingPassenger(true);
    setError(null);
    setFoundPassenger(null);
    setIsNewPassenger(false);

    try {
      const response = await apiClient.getUsers({ keyword: passengerPhone, limit: 5 });
      const pageResult = response.data as PageResult<UserType>;
      
      if (pageResult.records && pageResult.records.length > 0) {
        // Found existing passenger
        const passenger = pageResult.records[0];
        setFoundPassenger(passenger);
        setSuccessMessage(`Found passenger: ${passenger.full_name || passenger.first_name || 'Unknown'}`);
      } else {
        // No passenger found - prompt to create new
        setIsNewPassenger(true);
        setNewPassengerForm(prev => ({ ...prev, phone: passengerPhone }));
      }
    } catch (err) {
      console.error('Search failed:', err);
      setError('Failed to search for passenger. Please try again.');
    } finally {
      setIsSearchingPassenger(false);
    }
  };

  // Create new passenger
  const handleCreatePassenger = async () => {
    if (!newPassengerForm.first_name || !newPassengerForm.phone) {
      setError('First name and phone are required');
      return;
    }

    setIsCreatingPassenger(true);
    setError(null);

    try {
      const response = await apiClient.createUser({
        user_type: 'passenger',
        first_name: newPassengerForm.first_name,
        last_name: newPassengerForm.last_name,
        phone: newPassengerForm.phone,
        email: newPassengerForm.email,
        status: 'active',
      });

      const newPassenger = response.data as UserType;
      setFoundPassenger(newPassenger);
      setIsNewPassenger(false);
      setIsCreateModalOpen(false);
      setSuccessMessage('New passenger created successfully!');
    } catch (err) {
      console.error('Failed to create passenger:', err);
      setError('Failed to create passenger. Please try again.');
    } finally {
      setIsCreatingPassenger(false);
    }
  };

  // Move to next step
  const goToStep = (step: BookingStep) => {
    setCurrentStep(step);
  };

  // Poll live ETA after booking is complete
  useEffect(() => {
    if (!bookingComplete || !bookingDetails?.orderId) return;
    let cancelled = false;

    const pollETA = async () => {
      try {
        const res = await apiClient.getOrderETA(bookingDetails.orderId);
        if (!cancelled && res.code === '0000' && res.data) {
          const { eta_minutes, distance_km } = res.data;
          if (eta_minutes > 0) {
            setBookingDetails(prev => prev ? {
              ...prev,
              eta: `${eta_minutes} min (${distance_km.toFixed(1)} km away)`,
            } : prev);
          }
        }
      } catch {
        // Silently ignore ETA polling errors
      }
    };

    pollETA(); // Initial fetch
    const interval = setInterval(pollETA, 10000); // Poll every 10s

    return () => {
      cancelled = true;
      clearInterval(interval);
    };
  }, [bookingComplete, bookingDetails?.orderId]);

  // Load nearby available drivers with accurate ETA (polling while on confirm step)
  useEffect(() => {
    if (!FEATURE_NEARBY_DRIVERS) return;
    if (currentStep !== 'confirm') return;
    if (!pickupCoords) return;

    let cancelled = false;
    let interval: number | undefined;

    const load = async (showLoading: boolean) => {
      if (showLoading) setIsLoadingDrivers(true);
      setDriversError(null);
      try {
        const res = await apiClient.getNearbyDrivers({
          latitude: pickupCoords.lat,
          longitude: pickupCoords.lng,
          radius_km: 8,
          limit: 25,
          eta_mode: 'accurate',
        });
        if (cancelled) return;
        const drivers = (res.data?.drivers || [])
          .filter(d => d.is_online && !d.is_busy)
          .sort((a, b) => a.eta_minutes - b.eta_minutes);
        setAvailableDrivers(drivers);
      } catch (e) {
        if (cancelled) return;
        setDriversError(e instanceof Error ? e.message : 'Failed to load nearby drivers');
        setAvailableDrivers([]);
      } finally {
        if (!cancelled) setIsLoadingDrivers(false);
      }
    };

    load(true);
    interval = window.setInterval(() => load(false), 10000); // live updates

    return () => {
      cancelled = true;
      if (interval) window.clearInterval(interval);
    };
  }, [currentStep, pickupCoords]);

  // Confirm booking
  const confirmBooking = async () => {
    if (!foundPassenger || !pickupLocation || !dropoffLocation) {
      setError('Please complete all booking details');
      return;
    }

    setIsBooking(true);
    setError(null);
    setBookingErrorCode(null);

    // Validate coordinates are available
    if (!pickupCoords || !dropoffCoords) {
      setError('Please wait for location coordinates to be determined, or enter a valid address.');
      setIsBooking(false);
      return;
    }

    try {
      // Step 1: Get price estimate first (required to obtain price_id)
      const estimateData = {
        user_id: foundPassenger.user_id,
        pickup_latitude: pickupCoords.lat,
        pickup_longitude: pickupCoords.lng,
        pickup_address: pickupLocation,
        dropoff_latitude: dropoffCoords.lat,
        dropoff_longitude: dropoffCoords.lng,
        dropoff_address: dropoffLocation,
        order_type: 'ride',
        vehicle_category: 'sedan',
        vehicle_level: 'economy', // Match app defaults for same pricing
      };

      console.log('[QuickBooking] Getting price estimate...');
      const estimateResponse = await apiClient.estimateOrder(estimateData);
      
      let priceId: string | undefined;
      let fare: number | undefined;

      if (estimateResponse.code === '0000' && estimateResponse.data) {
        const estimate = estimateResponse.data as Record<string, unknown>;
        priceId = estimate.price_id as string;
        fare = (estimate.discounted_fare as number) || (estimate.original_fare as number);
        setEstimatedFare(fare);
        console.log('[QuickBooking] Got price_id:', priceId, 'fare:', fare);
      }

      if (!priceId) {
        setError(`Price estimation failed: ${estimateResponse.msg || 'No price ID returned'}. Please verify the locations and try again.`);
        setIsBooking(false);
        return;
      }

      // Step 2: Create order with price_id
      const orderData: Record<string, unknown> = {
        user_id: foundPassenger.user_id,
        order_type: 'ride',
        pickup_address: pickupLocation,
        dropoff_address: dropoffLocation,
        pickup_latitude: pickupCoords.lat,
        pickup_longitude: pickupCoords.lng,
        dropoff_latitude: dropoffCoords.lat,
        dropoff_longitude: dropoffCoords.lng,
        vehicle_category: 'sedan',
      };

      // Optional: manually select an available driver (provider_id)
      if (selectedDriverId && selectedDriverId !== 'auto') {
        orderData.provider_id = selectedDriverId;
      }

      // Add price_id if we got one from estimate
      if (priceId) {
        orderData.price_id = priceId;
      }
      if (fare) {
        orderData.estimated_fare = fare;
      }

      console.log('[QuickBooking] Creating order with data:', orderData);
      const response = await apiClient.createOrder(orderData);
      
      if (response.code === '0000' && response.data) {
        const resultData = response.data as Record<string, unknown>;
        const orderId = (resultData.order_id as string) || 'N/A';

        // Try to fetch assigned driver/vehicle details (if dispatch happened quickly)
        let driverName = 'Auto-assigned';
        let vehicleLabel: string | undefined;
        let plate: string | undefined;
        let driverStatus: string | undefined;
        let driverLocation: string | undefined;

        try {
          const detailRes = await apiClient.getOrderDetail(orderId);
          if (detailRes.code === '0000' && detailRes.data) {
            const o = detailRes.data as any;
            const d = o.driver || o.Driver;
            const v = o.vehicle || o.Vehicle;
            const details = o.details || o.Details;

            driverName =
              details?.driver_name ||
              d?.display_name ||
              d?.full_name ||
              d?.username ||
              `${d?.first_name || ''} ${d?.last_name || ''}`.trim() ||
              d?.phone ||
              driverName;

            plate = details?.license_plate || v?.plate_number || v?.PlateNumber || plate;
            const makeModel = `${v?.brand || v?.Brand || ''} ${v?.model || v?.Model || ''}`.trim();
            vehicleLabel = details?.vehicle_model || makeModel || vehicleLabel;

            driverStatus = d?.online_status || d?.status || o?.dispatch_status || o?.status || driverStatus;

            const lat = v?.current_latitude ?? v?.CurrentLatitude ?? d?.latitude ?? d?.Latitude;
            const lng = v?.current_longitude ?? v?.CurrentLongitude ?? d?.longitude ?? d?.Longitude;
            if (typeof lat === 'number' && typeof lng === 'number') {
              driverLocation = `${lat.toFixed(5)}, ${lng.toFixed(5)}`;
            }
          }
        } catch {
          // ignore; we can still show "Auto-assigned"
        }
        
        setBookingDetails({
          orderId,
          passenger: foundPassenger.full_name || foundPassenger.first_name || 'Passenger',
          driver: driverName,
          vehicle: vehicleLabel,
          plate,
          driverStatus,
          driverLocation,
          pickup: pickupLocation,
          dropoff: dropoffLocation,
          eta: fare ? `~${Math.round(fare).toLocaleString()} RWF` : 'Calculating...',
        });
        
        setBookingComplete(true);
        setSuccessMessage(
          selectedDriverId !== 'auto'
            ? 'Booking created successfully! Driver has been requested.'
            : 'Booking created successfully! Driver will be auto-assigned.'
        );
      } else {
        throw new Error(response.msg || 'Failed to create booking');
      }
    } catch (err) {
      console.error('Booking failed:', err);
      const isApiError = err instanceof ApiError;
      const code = isApiError ? err.code : null;
      const errorMsg = err instanceof Error ? err.message : 'Failed to create booking';
      setBookingErrorCode(code);
      setError(isApiError && err.serverMessage ? err.serverMessage : `Booking failed: ${errorMsg}. Please try again.`);
      setSuccessMessage(null);
    } finally {
      setIsBooking(false);
    }
  };

  // Reset for new booking
  const startNewBooking = () => {
    setCurrentStep('passenger');
    setPassengerPhone('');
    setFoundPassenger(null);
    setIsNewPassenger(false);
    setPickupLocation('');
    setDropoffLocation('');
    setPickupCoords(null);
    setDropoffCoords(null);
    setAvailableDrivers([]);
    setSelectedDriverId('auto');
    setDriversError(null);
    setBookingComplete(false);
    setBookingDetails(null);
    setEstimatedFare(undefined);
    setBookingErrorCode(null);
    setNewPassengerForm({ first_name: '', last_name: '', phone: '', email: '' });
  };

  const dismissError = () => {
    setError(null);
    setBookingErrorCode(null);
  };

  // Render step indicator
  const renderStepIndicator = () => {
    const steps: BookingStep[] = ['passenger', 'locations', 'confirm'];
    return (
      <div className="flex items-center justify-center gap-2 mb-8">
        {steps.map((step, index) => (
          <div key={step} className="flex items-center">
            <div
              className={`w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium transition-colors ${
                currentStep === step
                  ? 'bg-primary text-primary-foreground'
                  : index < steps.indexOf(currentStep)
                  ? 'bg-green-500 text-white'
                  : 'bg-muted text-muted-foreground'
              }`}
            >
              {index < steps.indexOf(currentStep) ? (
                <CheckCircle className="h-4 w-4" />
              ) : (
                index + 1
              )}
            </div>
            {index < steps.length - 1 && (
              <div
                className={`w-12 h-0.5 ${
                  index < steps.indexOf(currentStep)
                    ? 'bg-green-500'
                    : 'bg-muted'
                }`}
              />
            )}
          </div>
        ))}
      </div>
    );
  };

  // Booking complete view
  if (bookingComplete && bookingDetails) {
    return (
      <div className="max-w-2xl mx-auto space-y-6">
        <Card className="border-green-200 bg-green-50">
          <CardContent className="pt-6">
            <div className="text-center space-y-4">
              <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto">
                <CheckCircle className="h-8 w-8 text-green-600" />
              </div>
              <div>
                <h2 className="text-2xl font-bold text-green-900">Booking Confirmed!</h2>
                <p className="text-green-700">Order ID: {bookingDetails.orderId}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Booking Details</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-muted-foreground">Passenger</p>
                <p className="font-medium">{bookingDetails.passenger}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Driver</p>
                <p className="font-medium">{bookingDetails.driver}</p>
                {(bookingDetails.vehicle || bookingDetails.plate) && (
                  <p className="text-sm text-muted-foreground mt-1">
                    {[bookingDetails.plate, bookingDetails.vehicle].filter(Boolean).join(' • ')}
                  </p>
                )}
                {(bookingDetails.driverStatus || bookingDetails.driverLocation) && (
                  <p className="text-sm text-muted-foreground">
                    {[bookingDetails.driverStatus ? `Status: ${bookingDetails.driverStatus}` : null, bookingDetails.driverLocation ? `Location: ${bookingDetails.driverLocation}` : null]
                      .filter(Boolean)
                      .join(' • ')}
                  </p>
                )}
              </div>
            </div>
            <Separator />
            <div className="space-y-2">
              <div className="flex items-start gap-3">
                <div className="w-2 h-2 bg-green-500 rounded-full mt-2" />
                <div>
                  <p className="text-sm text-muted-foreground">Pickup</p>
                  <p className="font-medium">{bookingDetails.pickup}</p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <div className="w-2 h-2 bg-red-500 rounded-full mt-2" />
                <div>
                  <p className="text-sm text-muted-foreground">Drop-off</p>
                  <p className="font-medium">{bookingDetails.dropoff}</p>
                </div>
              </div>
            </div>
            <Separator />
            <div className="flex items-center gap-2 text-sm">
              <Clock className="h-4 w-4 text-muted-foreground" />
              <span>Estimated fare: <strong>{bookingDetails.eta}</strong></span>
            </div>
          </CardContent>
        </Card>

        <Button onClick={startNewBooking} className="w-full" size="lg">
          <Phone className="h-4 w-4 mr-2" />
          Start New Booking
        </Button>
      </div>
    );
  }

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      {/* Page Header */}
      <div className="text-center">
        <h1 className="text-2xl font-bold tracking-tight">Quick Booking</h1>
        <p className="text-muted-foreground">
          Create ride bookings for phone-in customers
        </p>
      </div>

      {/* Step Indicator */}
      {renderStepIndicator()}

      {/* Success/Error Messages */}
      {successMessage && (
        <div className="flex items-center gap-2 rounded-lg bg-green-50 border border-green-200 p-3 text-sm text-green-800">
          <CheckCircle className="h-4 w-4 flex-shrink-0" />
          <span>{successMessage}</span>
        </div>
      )}

      {error && (
        <div
          className={
            bookingErrorCode === '6007'
              ? 'rounded-lg border border-amber-200 bg-amber-50 p-4 text-amber-900'
              : 'flex items-center gap-2 rounded-lg bg-red-50 border border-red-200 p-3 text-sm text-red-800'
          }
        >
          <div className="flex items-start gap-2">
            <AlertCircle className="h-4 w-4 flex-shrink-0 mt-0.5" />
            <div className="flex-1 min-w-0">
              {bookingErrorCode === '6007' ? (
                <>
                  <p className="font-medium text-amber-900">Ride in progress</p>
                  <p className="text-sm text-amber-800 mt-1">{error}</p>
                  <p className="text-sm text-amber-700 mt-2">
                    This passenger already has an active ride. Cancel it from <strong>Rides</strong> (find the order → Cancel), then try again.
                  </p>
                </>
              ) : (
                <span>{error}</span>
              )}
            </div>
            <Button variant="ghost" size="sm" className="flex-shrink-0 h-6 w-6 p-0" onClick={dismissError}>
              <X className="h-4 w-4" />
            </Button>
          </div>
        </div>
      )}

      {/* Step 1: Passenger Search */}
      {currentStep === 'passenger' && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <User className="h-5 w-5" />
              Step 1: Find or Create Passenger
            </CardTitle>
            <CardDescription>
              Search by phone number to find existing passenger or create a new one
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex gap-2">
              <div className="flex-1">
                <Input
                  placeholder="Enter phone number (e.g., +250788123456)"
                  value={passengerPhone}
                  onChange={(e) => setPassengerPhone(e.target.value)}
                  onKeyDown={(e) => e.key === 'Enter' && searchPassenger()}
                />
              </div>
              <Button onClick={searchPassenger} disabled={isSearchingPassenger}>
                {isSearchingPassenger ? (
                  <RefreshCw className="h-4 w-4 animate-spin" />
                ) : (
                  <Search className="h-4 w-4" />
                )}
              </Button>
            </div>

            {/* Found Passenger */}
            {foundPassenger && (
              <div className="p-4 rounded-lg border bg-green-50 border-green-200">
                <div className="flex items-center gap-3">
                  <Avatar className="h-12 w-12">
                    <AvatarFallback className="bg-green-100 text-green-700">
                      {(foundPassenger.full_name || foundPassenger.first_name || 'P')[0].toUpperCase()}
                    </AvatarFallback>
                  </Avatar>
                  <div className="flex-1">
                    <p className="font-medium">{foundPassenger.full_name || `${foundPassenger.first_name || ''} ${foundPassenger.last_name || ''}`.trim() || 'Passenger'}</p>
                    <p className="text-sm text-muted-foreground">{foundPassenger.phone}</p>
                    <p className="text-xs text-muted-foreground">
                      {foundPassenger.total_rides || 0} previous rides
                    </p>
                  </div>
                  <Badge variant="outline" className="text-green-700 border-green-300">
                    <CheckCircle className="h-3 w-3 mr-1" />
                    Found
                  </Badge>
                </div>
              </div>
            )}

            {/* New Passenger Prompt */}
            {isNewPassenger && (
              <div className="p-4 rounded-lg border bg-amber-50 border-amber-200">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium text-amber-800">No passenger found</p>
                    <p className="text-sm text-amber-700">Create a new passenger account for {passengerPhone}</p>
                  </div>
                  <Button onClick={() => setIsCreateModalOpen(true)}>
                    <UserPlus className="h-4 w-4 mr-2" />
                    Create New
                  </Button>
                </div>
              </div>
            )}

            {/* Next Button */}
            {foundPassenger && (
              <Button onClick={() => goToStep('locations')} className="w-full" size="lg">
                Continue to Locations
                <Navigation className="h-4 w-4 ml-2" />
              </Button>
            )}
          </CardContent>
        </Card>
      )}

      {/* Step 2: Locations */}
      {currentStep === 'locations' && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <MapPin className="h-5 w-5" />
              Step 2: Enter Locations
            </CardTitle>
            <CardDescription>
              Enter the pickup and drop-off locations
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="pickup">Pickup Location</Label>
              <div className="relative">
                <MapPin className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-green-600" />
                <div className="pl-10">
                  {isGoogleLoaded ? (
                    <PlacesAutocompleteInput
                      id="pickup"
                      placeholder="Search pickup location..."
                      value={pickupLocation}
                      onChange={(v) => {
                        setPickupFromPlaces(false);
                        setPickupCoords(null);
                        setPickupLocation(v);
                      }}
                      onSelect={(place: PlaceSelection) => {
                        setPickupFromPlaces(true);
                        setPickupLocation(place.address);
                        setPickupCoords({ lat: place.lat, lng: place.lng });
                      }}
                      className="w-full"
                    />
                  ) : (
                    <Input
                      id="pickup"
                      placeholder="Enter pickup address..."
                      value={pickupLocation}
                      onChange={(e) => {
                        setPickupFromPlaces(false);
                        setPickupCoords(null);
                        setPickupLocation(e.target.value);
                      }}
                      className="pl-10"
                    />
                  )}
                </div>
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="dropoff">Drop-off Location</Label>
              <div className="relative">
                <MapPin className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-red-600" />
                <div className="pl-10">
                  {isGoogleLoaded ? (
                    <PlacesAutocompleteInput
                      id="dropoff"
                      placeholder="Search drop-off location..."
                      value={dropoffLocation}
                      onChange={(v) => {
                        setDropoffFromPlaces(false);
                        setDropoffCoords(null);
                        setDropoffLocation(v);
                      }}
                      onSelect={(place: PlaceSelection) => {
                        setDropoffFromPlaces(true);
                        setDropoffLocation(place.address);
                        setDropoffCoords({ lat: place.lat, lng: place.lng });
                      }}
                      className="w-full"
                    />
                  ) : (
                    <Input
                      id="dropoff"
                      placeholder="Enter drop-off address..."
                      value={dropoffLocation}
                      onChange={(e) => {
                        setDropoffFromPlaces(false);
                        setDropoffCoords(null);
                        setDropoffLocation(e.target.value);
                      }}
                      className="pl-10"
                    />
                  )}
                </div>
              </div>
            </div>

            {/* Quick Location Suggestions */}
            <div className="space-y-2">
              <p className="text-sm text-muted-foreground">Popular locations:</p>
              <div className="flex flex-wrap gap-2">
                {Object.entries(POPULAR_LOCATIONS).slice(0, 5).map(([name, location]) => (
                  <Button
                    key={name}
                    variant="outline"
                    size="sm"
                    onClick={() => {
                      if (!pickupLocation) {
                        setPickupLocation(location.address);
                        setPickupCoords({ lat: location.lat, lng: location.lng });
                      } else {
                        setDropoffLocation(location.address);
                        setDropoffCoords({ lat: location.lat, lng: location.lng });
                      }
                    }}
                  >
                    {name}
                  </Button>
                ))}
              </div>
            </div>

            {/* Coordinate Status Indicators */}
            {(isGeocodingPickup || isGeocodingDropoff) && (
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <RefreshCw className="h-4 w-4 animate-spin" />
                <span>Determining location coordinates...</span>
              </div>
            )}
            
            {pickupLocation && pickupCoords && (
              <div className="text-xs text-green-600 flex items-center gap-1">
                <CheckCircle className="h-3 w-3" />
                Pickup location verified
              </div>
            )}
            
            {dropoffLocation && dropoffCoords && (
              <div className="text-xs text-green-600 flex items-center gap-1">
                <CheckCircle className="h-3 w-3" />
                Dropoff location verified
              </div>
            )}
            
            {((pickupLocation && !pickupCoords && !isGeocodingPickup) || (dropoffLocation && !dropoffCoords && !isGeocodingDropoff)) && (
              <div className="text-xs text-amber-600 flex items-center gap-1">
                <AlertCircle className="h-3 w-3" />
                Could not determine coordinates. Please try a more specific address.
              </div>
            )}

            <div className="flex gap-2">
              <Button variant="outline" onClick={() => setCurrentStep('passenger')} className="flex-1">
                Back
              </Button>
              <Button 
                onClick={() => goToStep('confirm')} 
                className="flex-1"
                disabled={!pickupLocation || !dropoffLocation || !pickupCoords || !dropoffCoords || isGeocodingPickup || isGeocodingDropoff}
              >
                Continue to Review
                <CheckCircle className="h-4 w-4 ml-2" />
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Step 3: Confirm Booking */}
      {currentStep === 'confirm' && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <CheckCircle className="h-5 w-5" />
              Step 3: Confirm Booking
            </CardTitle>
            <CardDescription>
              Review the booking details before confirmation
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {/* Passenger */}
            <div className="p-4 rounded-lg bg-muted/50">
              <p className="text-sm text-muted-foreground mb-1">Passenger</p>
              <p className="font-medium">{foundPassenger?.full_name || foundPassenger?.first_name || 'Passenger'}</p>
              <p className="text-sm text-muted-foreground">{foundPassenger?.phone}</p>
            </div>

            {/* Locations */}
            <div className="space-y-2">
              <div className="flex items-start gap-3">
                <div className="w-3 h-3 bg-green-500 rounded-full mt-1.5" />
                <div>
                  <p className="text-sm text-muted-foreground">Pickup</p>
                  <p className="font-medium">{pickupLocation}</p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <div className="w-3 h-3 bg-red-500 rounded-full mt-1.5" />
                <div>
                  <p className="text-sm text-muted-foreground">Drop-off</p>
                  <p className="font-medium">{dropoffLocation}</p>
                </div>
              </div>
            </div>

            {/* Estimated Fare */}
            {estimatedFare && (
              <div className="p-4 rounded-lg bg-green-50 border border-green-200">
                <p className="text-sm text-muted-foreground mb-1">Estimated Fare</p>
                <p className="text-2xl font-bold text-green-800">
                  {Math.round(estimatedFare).toLocaleString()} RWF
                </p>
              </div>
            )}

            {/* Driver Assignment — off when FEATURE_NEARBY_DRIVERS is false */}
            {FEATURE_NEARBY_DRIVERS ? (
              <div className="p-4 rounded-lg bg-blue-50 border border-blue-200 space-y-3">
                <div className="flex items-start gap-2">
                  <Car className="h-5 w-5 text-blue-600 mt-0.5" />
                  <div className="flex-1">
                    <p className="font-medium text-blue-900">Driver Assignment</p>
                    <p className="text-sm text-blue-700">
                      Choose an available driver (with live location + ETA), or leave it on auto-assign.
                    </p>
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="driver-select">Select Driver</Label>
                  <Select value={selectedDriverId} onValueChange={setSelectedDriverId}>
                    <SelectTrigger id="driver-select">
                      <SelectValue placeholder="Select a driver" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="auto">Auto-assign (system)</SelectItem>
                      {availableDrivers.map((d) => (
                        <SelectItem key={d.driver_id} value={d.driver_id}>
                          {d.name} • {Math.round(d.distance_km * 10) / 10} km • {d.eta_minutes} min
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <div className="text-xs text-muted-foreground">
                    {isLoadingDrivers ? 'Updating drivers (10s)…' : 'Live updates every 10s.'}
                    {driversError ? ` • ${driversError}` : ''}
                  </div>
                </div>

                {selectedDriverId !== 'auto' && selectedDriver && (
                  <div className="rounded-md bg-background border p-3 space-y-2">
                    <div className="flex items-center justify-between gap-2">
                      <div>
                        <p className="text-sm font-medium">{selectedDriver.name}</p>
                        <p className="text-xs text-muted-foreground">
                          ETA to pickup: <strong>{selectedDriver.eta_minutes} min</strong> • Distance:{' '}
                          <strong>{Math.round(selectedDriver.distance_km * 10) / 10} km</strong>
                        </p>
                      </div>
                      <Badge className="bg-green-100 text-green-700 hover:bg-green-100">Available</Badge>
                    </div>

                    {isGoogleLoaded && pickupCoords && (
                      <div className="h-48 w-full rounded-md overflow-hidden">
                        <GoogleMap
                          mapContainerStyle={{ width: '100%', height: '100%' }}
                          center={{ lat: pickupCoords.lat, lng: pickupCoords.lng }}
                          zoom={14}
                          options={{ disableDefaultUI: true, zoomControl: true }}
                        >
                          <Marker
                            position={{ lat: pickupCoords.lat, lng: pickupCoords.lng }}
                            title="Pickup"
                            label={{ text: 'P', color: 'white' }}
                          />
                          <Marker
                            position={{ lat: selectedDriver.latitude, lng: selectedDriver.longitude }}
                            title={selectedDriver.name}
                            label={{ text: 'D', color: 'white' }}
                          />
                        </GoogleMap>
                      </div>
                    )}
                  </div>
                )}
              </div>
            ) : (
              <div className="p-4 rounded-lg bg-muted/50 border">
                <p className="text-sm text-muted-foreground">
                  Driver will be auto-assigned by the system when the booking is created.
                </p>
              </div>
            )}

            <Separator />

            <div className="flex gap-2">
              <Button variant="outline" onClick={() => setCurrentStep('locations')} className="flex-1">
                Back
              </Button>
              <Button 
                onClick={confirmBooking} 
                className="flex-1"
                disabled={isBooking}
              >
                {isBooking ? (
                  <>
                    <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                    Creating Booking...
                  </>
                ) : (
                  <>
                    <CheckCircle className="h-4 w-4 mr-2" />
                    Confirm Booking
                  </>
                )}
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Create Passenger Modal */}
      <Dialog open={isCreateModalOpen} onOpenChange={setIsCreateModalOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Create New Passenger</DialogTitle>
            <DialogDescription>
              Enter passenger details to create their account
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="new_first_name">First Name *</Label>
                <Input
                  id="new_first_name"
                  placeholder="John"
                  value={newPassengerForm.first_name}
                  onChange={(e) => setNewPassengerForm({ ...newPassengerForm, first_name: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="new_last_name">Last Name</Label>
                <Input
                  id="new_last_name"
                  placeholder="Doe"
                  value={newPassengerForm.last_name}
                  onChange={(e) => setNewPassengerForm({ ...newPassengerForm, last_name: e.target.value })}
                />
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="new_phone">Phone Number *</Label>
              <Input
                id="new_phone"
                placeholder="+250 788 123 456"
                value={newPassengerForm.phone}
                onChange={(e) => setNewPassengerForm({ ...newPassengerForm, phone: e.target.value })}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="new_email">Email (optional)</Label>
              <Input
                id="new_email"
                type="email"
                placeholder="john@example.com"
                value={newPassengerForm.email}
                onChange={(e) => setNewPassengerForm({ ...newPassengerForm, email: e.target.value })}
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsCreateModalOpen(false)} disabled={isCreatingPassenger}>
              Cancel
            </Button>
            <Button onClick={handleCreatePassenger} disabled={isCreatingPassenger}>
              {isCreatingPassenger ? 'Creating...' : 'Create Passenger'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}

