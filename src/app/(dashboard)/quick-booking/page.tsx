'use client';

import { useState, useEffect, useCallback } from 'react';
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
import { apiClient } from '@/lib/api-client';
import type { Driver, User as UserType, PageResult } from '@/types';

// Mock nearby drivers for demo
const MOCK_NEARBY_DRIVERS = [
  { id: 'DRV001', name: 'Peter Mutombo', vehicle: 'Toyota Corolla - RAD 123A', rating: 4.8, distance: '0.5 km', eta: '3 min', status: 'online' },
  { id: 'DRV002', name: 'David Kagame', vehicle: 'Honda Fit - RAC 456B', rating: 4.6, distance: '1.2 km', eta: '5 min', status: 'online' },
  { id: 'DRV005', name: 'Alex Munyaneza', vehicle: 'Suzuki Swift - RAB 789C', rating: 4.5, distance: '2.0 km', eta: '8 min', status: 'online' },
];

type BookingStep = 'passenger' | 'locations' | 'driver' | 'confirm';

export default function QuickBookingPage() {
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
  
  // Driver selection
  const [nearbyDrivers, setNearbyDrivers] = useState<typeof MOCK_NEARBY_DRIVERS>([]);
  const [selectedDriver, setSelectedDriver] = useState<typeof MOCK_NEARBY_DRIVERS[0] | null>(null);
  const [isLoadingDrivers, setIsLoadingDrivers] = useState(false);
  
  // Booking state
  const [isBooking, setIsBooking] = useState(false);
  const [bookingComplete, setBookingComplete] = useState(false);
  const [bookingDetails, setBookingDetails] = useState<{
    orderId: string;
    passenger: string;
    driver: string;
    pickup: string;
    dropoff: string;
    eta: string;
  } | null>(null);
  
  // Messages
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  
  // Create new passenger modal
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isCreatingPassenger, setIsCreatingPassenger] = useState(false);

  // Auto-hide messages
  useEffect(() => {
    if (successMessage) {
      const timer = setTimeout(() => setSuccessMessage(null), 5000);
      return () => clearTimeout(timer);
    }
  }, [successMessage]);

  useEffect(() => {
    if (error) {
      const timer = setTimeout(() => setError(null), 5000);
      return () => clearTimeout(timer);
    }
  }, [error]);

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

  // Load nearby drivers
  const loadNearbyDrivers = useCallback(async () => {
    setIsLoadingDrivers(true);
    try {
      // In demo mode, use mock nearby drivers
      await new Promise(r => setTimeout(r, 500));
      setNearbyDrivers(MOCK_NEARBY_DRIVERS);
    } catch (err) {
      setError('Failed to load nearby drivers');
    } finally {
      setIsLoadingDrivers(false);
    }
  }, []);

  // Move to next step
  const goToStep = (step: BookingStep) => {
    if (step === 'driver') {
      loadNearbyDrivers();
    }
    setCurrentStep(step);
  };

  // Confirm booking
  const confirmBooking = async () => {
    if (!foundPassenger || !selectedDriver || !pickupLocation || !dropoffLocation) {
      setError('Please complete all booking details');
      return;
    }

    setIsBooking(true);
    setError(null);

    try {
      // Simulate booking creation
      await new Promise(r => setTimeout(r, 1000));
      
      const orderId = `ORD${Date.now().toString().slice(-6)}`;
      
      setBookingDetails({
        orderId,
        passenger: foundPassenger.full_name || foundPassenger.first_name || 'Passenger',
        driver: selectedDriver.name,
        pickup: pickupLocation,
        dropoff: dropoffLocation,
        eta: selectedDriver.eta,
      });
      
      setBookingComplete(true);
      setSuccessMessage('Booking created successfully!');
    } catch (err) {
      console.error('Booking failed:', err);
      setError('Failed to create booking. Please try again.');
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
    setSelectedDriver(null);
    setNearbyDrivers([]);
    setBookingComplete(false);
    setBookingDetails(null);
    setNewPassengerForm({ first_name: '', last_name: '', phone: '', email: '' });
  };

  // Render step indicator
  const renderStepIndicator = () => (
    <div className="flex items-center justify-center gap-2 mb-8">
      {(['passenger', 'locations', 'driver', 'confirm'] as BookingStep[]).map((step, index) => (
        <div key={step} className="flex items-center">
          <div
            className={`w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium transition-colors ${
              currentStep === step
                ? 'bg-primary text-primary-foreground'
                : index < ['passenger', 'locations', 'driver', 'confirm'].indexOf(currentStep)
                ? 'bg-green-500 text-white'
                : 'bg-muted text-muted-foreground'
            }`}
          >
            {index < ['passenger', 'locations', 'driver', 'confirm'].indexOf(currentStep) ? (
              <CheckCircle className="h-4 w-4" />
            ) : (
              index + 1
            )}
          </div>
          {index < 3 && (
            <div
              className={`w-12 h-0.5 ${
                index < ['passenger', 'locations', 'driver', 'confirm'].indexOf(currentStep)
                  ? 'bg-green-500'
                  : 'bg-muted'
              }`}
            />
          )}
        </div>
      ))}
    </div>
  );

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
              <span>Estimated arrival: <strong>{bookingDetails.eta}</strong></span>
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
        <div className="flex items-center gap-2 rounded-lg bg-red-50 border border-red-200 p-3 text-sm text-red-800">
          <AlertCircle className="h-4 w-4 flex-shrink-0" />
          <span>{error}</span>
          <Button variant="ghost" size="sm" className="ml-auto h-6 w-6 p-0" onClick={() => setError(null)}>
            <X className="h-4 w-4" />
          </Button>
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
                <Input
                  id="pickup"
                  placeholder="Enter pickup address..."
                  value={pickupLocation}
                  onChange={(e) => setPickupLocation(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="dropoff">Drop-off Location</Label>
              <div className="relative">
                <MapPin className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-red-600" />
                <Input
                  id="dropoff"
                  placeholder="Enter drop-off address..."
                  value={dropoffLocation}
                  onChange={(e) => setDropoffLocation(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>

            {/* Quick Location Suggestions */}
            <div className="space-y-2">
              <p className="text-sm text-muted-foreground">Popular locations:</p>
              <div className="flex flex-wrap gap-2">
                {['Kigali Convention Centre', 'Kigali International Airport', 'Nyarutarama', 'Downtown Kigali', 'Remera'].map(loc => (
                  <Button
                    key={loc}
                    variant="outline"
                    size="sm"
                    onClick={() => !pickupLocation ? setPickupLocation(loc) : setDropoffLocation(loc)}
                  >
                    {loc}
                  </Button>
                ))}
              </div>
            </div>

            <div className="flex gap-2">
              <Button variant="outline" onClick={() => setCurrentStep('passenger')} className="flex-1">
                Back
              </Button>
              <Button 
                onClick={() => goToStep('driver')} 
                className="flex-1"
                disabled={!pickupLocation || !dropoffLocation}
              >
                Find Nearby Drivers
                <Car className="h-4 w-4 ml-2" />
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Step 3: Select Driver */}
      {currentStep === 'driver' && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Car className="h-5 w-5" />
              Step 3: Assign Driver
            </CardTitle>
            <CardDescription>
              Select an available driver near the pickup location
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {isLoadingDrivers ? (
              <div className="space-y-3">
                {[1, 2, 3].map(i => (
                  <div key={i} className="flex items-center gap-3 p-3 border rounded-lg">
                    <Skeleton className="h-10 w-10 rounded-full" />
                    <div className="flex-1 space-y-2">
                      <Skeleton className="h-4 w-32" />
                      <Skeleton className="h-3 w-48" />
                    </div>
                    <Skeleton className="h-8 w-20" />
                  </div>
                ))}
              </div>
            ) : (
              <div className="space-y-3">
                {nearbyDrivers.map(driver => (
                  <div
                    key={driver.id}
                    className={`flex items-center gap-3 p-3 border rounded-lg cursor-pointer transition-colors ${
                      selectedDriver?.id === driver.id
                        ? 'border-primary bg-primary/5'
                        : 'hover:bg-muted/50'
                    }`}
                    onClick={() => setSelectedDriver(driver)}
                  >
                    <Avatar className="h-10 w-10">
                      <AvatarFallback className="bg-primary/10 text-primary">
                        {driver.name[0]}
                      </AvatarFallback>
                    </Avatar>
                    <div className="flex-1">
                      <div className="flex items-center gap-2">
                        <p className="font-medium">{driver.name}</p>
                        <div className="flex items-center gap-1 text-sm text-yellow-600">
                          <Star className="h-3 w-3 fill-current" />
                          {driver.rating}
                        </div>
                      </div>
                      <p className="text-sm text-muted-foreground">{driver.vehicle}</p>
                    </div>
                    <div className="text-right">
                      <p className="text-sm font-medium text-green-600">{driver.eta}</p>
                      <p className="text-xs text-muted-foreground">{driver.distance}</p>
                    </div>
                    {selectedDriver?.id === driver.id && (
                      <CheckCircle className="h-5 w-5 text-primary" />
                    )}
                  </div>
                ))}
              </div>
            )}

            <div className="flex gap-2">
              <Button variant="outline" onClick={() => setCurrentStep('locations')} className="flex-1">
                Back
              </Button>
              <Button 
                onClick={() => goToStep('confirm')} 
                className="flex-1"
                disabled={!selectedDriver}
              >
                Review Booking
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Step 4: Confirm Booking */}
      {currentStep === 'confirm' && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <CheckCircle className="h-5 w-5" />
              Step 4: Confirm Booking
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

            {/* Driver */}
            {selectedDriver && (
              <div className="p-4 rounded-lg bg-muted/50">
                <p className="text-sm text-muted-foreground mb-2">Assigned Driver</p>
                <div className="flex items-center gap-3">
                  <Avatar className="h-10 w-10">
                    <AvatarFallback className="bg-primary/10 text-primary">
                      {selectedDriver.name[0]}
                    </AvatarFallback>
                  </Avatar>
                  <div>
                    <p className="font-medium">{selectedDriver.name}</p>
                    <p className="text-sm text-muted-foreground">{selectedDriver.vehicle}</p>
                  </div>
                  <Badge variant="outline" className="ml-auto">
                    <Clock className="h-3 w-3 mr-1" />
                    ETA: {selectedDriver.eta}
                  </Badge>
                </div>
              </div>
            )}

            <Separator />

            <div className="flex gap-2">
              <Button variant="outline" onClick={() => setCurrentStep('driver')} className="flex-1">
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

