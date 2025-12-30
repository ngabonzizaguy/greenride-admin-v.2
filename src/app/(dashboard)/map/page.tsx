'use client';

import { useState } from 'react';
import { 
  Car, 
  MapPin, 
  Filter,
  RefreshCw,
  Eye,
  Navigation,
  Info
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Checkbox } from '@/components/ui/checkbox';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';

// Mock driver locations
const mockDrivers = [
  { id: '1', name: 'Peter M.', status: 'available', lat: -1.9403, lng: 29.8739, vehicle: 'RAD 123A' },
  { id: '2', name: 'David K.', status: 'on_trip', lat: -1.9453, lng: 29.8789, vehicle: 'RAB 456B' },
  { id: '3', name: 'Jean P.', status: 'arriving', lat: -1.9353, lng: 29.8689, vehicle: 'RAC 789C' },
  { id: '4', name: 'Emmanuel H.', status: 'available', lat: -1.9503, lng: 29.8839, vehicle: 'RAD 012D' },
  { id: '5', name: 'Claude U.', status: 'on_trip', lat: -1.9303, lng: 29.8639, vehicle: 'RAF 678F' },
];

const stats = {
  onlineDrivers: 24,
  activeRides: 12,
  availableDrivers: 12,
  arrivingDrivers: 4,
};

const getStatusColor = (status: string) => {
  switch (status) {
    case 'available':
      return 'bg-green-500';
    case 'on_trip':
      return 'bg-yellow-500';
    case 'arriving':
      return 'bg-blue-500';
    default:
      return 'bg-gray-500';
  }
};

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'available':
      return 'Available';
    case 'on_trip':
      return 'On Trip';
    case 'arriving':
      return 'Arriving';
    default:
      return status;
  }
};

export default function MapPage() {
  const [showAvailable, setShowAvailable] = useState(true);
  const [showOnTrip, setShowOnTrip] = useState(true);
  const [showRoutes, setShowRoutes] = useState(true);
  const [vehicleType, setVehicleType] = useState('all');
  const [autoRefresh, setAutoRefresh] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedDriver, setSelectedDriver] = useState<typeof mockDrivers[0] | null>(null);

  const filteredDrivers = mockDrivers.filter((driver) => {
    const matchesSearch = driver.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                          driver.vehicle.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesStatus = 
      (showAvailable && driver.status === 'available') ||
      (showOnTrip && driver.status === 'on_trip') ||
      driver.status === 'arriving';
    return matchesSearch && matchesStatus;
  });

  return (
    <div className="h-[calc(100vh-8rem)] relative">
      {/* Map Container - Placeholder for Google Maps */}
      <div className="absolute inset-0 bg-gray-200 rounded-lg overflow-hidden">
        {/* Map placeholder with gradient */}
        <div 
          className="w-full h-full"
          style={{
            background: 'linear-gradient(135deg, #e0f2e9 0%, #c6e2d5 50%, #a8d4be 100%)',
          }}
        >
          {/* Grid overlay to simulate map */}
          <div 
            className="w-full h-full opacity-10"
            style={{
              backgroundImage: `
                linear-gradient(rgba(0,0,0,0.1) 1px, transparent 1px),
                linear-gradient(90deg, rgba(0,0,0,0.1) 1px, transparent 1px)
              `,
              backgroundSize: '50px 50px',
            }}
          />
          
          {/* Driver markers */}
          <div className="absolute inset-0">
            {filteredDrivers.map((driver, index) => (
              <div
                key={driver.id}
                className="absolute cursor-pointer transform -translate-x-1/2 -translate-y-1/2 transition-transform hover:scale-110"
                style={{
                  left: `${20 + index * 15}%`,
                  top: `${30 + (index % 3) * 20}%`,
                }}
                onClick={() => setSelectedDriver(driver)}
              >
                <div className={`w-10 h-10 rounded-full ${getStatusColor(driver.status)} flex items-center justify-center shadow-lg border-2 border-white`}>
                  <Car className="h-5 w-5 text-white" />
                </div>
                <div className="absolute -bottom-6 left-1/2 -translate-x-1/2 whitespace-nowrap">
                  <span className="text-xs font-medium bg-white px-2 py-1 rounded shadow">
                    {driver.name}
                  </span>
                </div>
              </div>
            ))}
          </div>

          {/* Map center marker */}
          <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2">
            <div className="text-center">
              <MapPin className="h-8 w-8 text-primary mx-auto" />
              <p className="text-xs font-medium mt-1 bg-white px-2 py-1 rounded shadow">
                Kigali City Center
              </p>
            </div>
          </div>
        </div>

        {/* Map Controls Hint */}
        <div className="absolute bottom-4 left-4 bg-white/90 backdrop-blur rounded-lg p-3 shadow-lg">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Info className="h-4 w-4" />
            <span>Google Maps integration ready. Add API key to enable.</span>
          </div>
        </div>
      </div>

      {/* Control Panel - Top Right */}
      <Card className="absolute top-4 right-4 w-72 shadow-xl">
        <CardContent className="p-4 space-y-4">
          <div className="font-semibold">Map Controls</div>
          
          {/* Search */}
          <Input
            placeholder="Find driver..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />

          {/* Filters */}
          <div className="space-y-2">
            <div className="flex items-center space-x-2">
              <Checkbox 
                id="show-available" 
                checked={showAvailable}
                onCheckedChange={(checked) => setShowAvailable(!!checked)}
              />
              <Label htmlFor="show-available" className="flex items-center gap-2">
                <span className="h-2 w-2 rounded-full bg-green-500" />
                Show Available
              </Label>
            </div>
            <div className="flex items-center space-x-2">
              <Checkbox 
                id="show-on-trip" 
                checked={showOnTrip}
                onCheckedChange={(checked) => setShowOnTrip(!!checked)}
              />
              <Label htmlFor="show-on-trip" className="flex items-center gap-2">
                <span className="h-2 w-2 rounded-full bg-yellow-500" />
                Show On Trip
              </Label>
            </div>
            <div className="flex items-center space-x-2">
              <Checkbox 
                id="show-routes" 
                checked={showRoutes}
                onCheckedChange={(checked) => setShowRoutes(!!checked)}
              />
              <Label htmlFor="show-routes" className="flex items-center gap-2">
                <Navigation className="h-3 w-3 text-primary" />
                Show Routes
              </Label>
            </div>
          </div>

          {/* Vehicle Type */}
          <Select value={vehicleType} onValueChange={setVehicleType}>
            <SelectTrigger>
              <SelectValue placeholder="Vehicle Type" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Types</SelectItem>
              <SelectItem value="sedan">Sedan</SelectItem>
              <SelectItem value="suv">SUV</SelectItem>
              <SelectItem value="moto">Moto</SelectItem>
            </SelectContent>
          </Select>

          {/* Auto Refresh */}
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <Checkbox 
                id="auto-refresh" 
                checked={autoRefresh}
                onCheckedChange={(checked) => setAutoRefresh(!!checked)}
              />
              <Label htmlFor="auto-refresh">Auto-refresh</Label>
            </div>
            <Button variant="ghost" size="icon" className="h-8 w-8">
              <RefreshCw className="h-4 w-4" />
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Stats Panel - Bottom Left */}
      <Card className="absolute bottom-4 left-1/2 -translate-x-1/2 shadow-xl">
        <CardContent className="p-4">
          <div className="flex items-center gap-6">
            <div className="flex items-center gap-2">
              <span className="h-3 w-3 rounded-full bg-green-500 animate-pulse" />
              <span className="text-sm font-medium">{stats.onlineDrivers} drivers online</span>
            </div>
            <div className="flex items-center gap-2">
              <span className="h-3 w-3 rounded-full bg-blue-500" />
              <span className="text-sm font-medium">{stats.activeRides} active rides</span>
            </div>
            <div className="text-xs text-muted-foreground">
              Updated just now
            </div>
          </div>
          {/* Legend */}
          <div className="flex items-center gap-4 mt-3 pt-3 border-t">
            <div className="flex items-center gap-1.5">
              <span className="h-3 w-3 rounded-full bg-green-500" />
              <span className="text-xs text-muted-foreground">Available</span>
            </div>
            <div className="flex items-center gap-1.5">
              <span className="h-3 w-3 rounded-full bg-yellow-500" />
              <span className="text-xs text-muted-foreground">On Trip</span>
            </div>
            <div className="flex items-center gap-1.5">
              <span className="h-3 w-3 rounded-full bg-blue-500" />
              <span className="text-xs text-muted-foreground">Arriving</span>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Selected Driver Panel */}
      {selectedDriver && (
        <Card className="absolute top-4 left-4 w-72 shadow-xl">
          <CardContent className="p-4">
            <div className="flex items-start justify-between">
              <div className="flex items-center gap-3">
                <Avatar className="h-12 w-12">
                  <AvatarFallback className="bg-primary text-primary-foreground">
                    {selectedDriver.name.split(' ').map(n => n[0]).join('')}
                  </AvatarFallback>
                </Avatar>
                <div>
                  <p className="font-semibold">{selectedDriver.name}</p>
                  <p className="text-sm text-muted-foreground">{selectedDriver.vehicle}</p>
                </div>
              </div>
              <Button 
                variant="ghost" 
                size="icon" 
                className="h-6 w-6"
                onClick={() => setSelectedDriver(null)}
              >
                ×
              </Button>
            </div>
            <div className="mt-3">
              <Badge className={`${getStatusColor(selectedDriver.status)} text-white`}>
                {getStatusLabel(selectedDriver.status)}
              </Badge>
            </div>
            <div className="flex gap-2 mt-4">
              <Button size="sm" className="flex-1">
                <Eye className="h-4 w-4 mr-1" />
                View Details
              </Button>
              <Button size="sm" variant="outline" className="flex-1">
                <Navigation className="h-4 w-4 mr-1" />
                Track
              </Button>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
