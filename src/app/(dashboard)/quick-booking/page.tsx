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
import type { User as UserType, PageResult } from '@/types';

type BookingStep = 'passenger' | 'locations' | 'confirm';

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

  // Move to next step
  const goToStep = (step: BookingStep) => {
    setCurrentStep(step);
  };

  // Confirm booking
  const confirmBooking = async () => {
    if (!foundPassenger || !pickupLocation || !dropoffLocation) {
      setError('Please complete all booking details');
      return;
    }

    setIsBooking(true);
    setError(null);

    try {
      // Create order via API
      const orderData = {
        user_id: foundPassenger.user_id,
        order_type: 'ride',
        pickup_address: pickupLocation,
        dropoff_address: dropoffLocation,
        // Driver will be auto-assigned by the system
      };

      const response = await apiClient.createOrder(orderData);
      
      if (response.code === '0000' && response.data) {
        const orderData = response.data as Record<string, unknown>;
        const orderId = (orderData.order_id as string) || 'N/A';
        
        setBookingDetails({
          orderId,
          passenger: foundPassenger.full_name || foundPassenger.first_name || 'Passenger',
          driver: 'Auto-assigned',
          pickup: pickupLocation,
          dropoff: dropoffLocation,
          eta: 'Estimated',
        });
        
        setBookingComplete(true);
        setSuccessMessage('Booking created successfully! Driver will be auto-assigned.');
      } else {
        throw new Error(response.msg || 'Failed to create booking');
      }
    } catch (err) {
      console.error('Booking failed:', err);
      setError('Failed to create booking. Please try again.');
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
    setBookingComplete(false);
    setBookingDetails(null);
    setNewPassengerForm({ first_name: '', last_name: '', phone: '', email: '' });
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
                onClick={() => goToStep('confirm')} 
                className="flex-1"
                disabled={!pickupLocation || !dropoffLocation}
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

            {/* Driver Assignment Note */}
            <div className="p-4 rounded-lg bg-blue-50 border border-blue-200">
              <div className="flex items-start gap-2">
                <Car className="h-5 w-5 text-blue-600 mt-0.5" />
                <div>
                  <p className="font-medium text-blue-900">Driver Assignment</p>
                  <p className="text-sm text-blue-700">A driver will be automatically assigned by the system when the booking is confirmed.</p>
                </div>
              </div>
            </div>

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

