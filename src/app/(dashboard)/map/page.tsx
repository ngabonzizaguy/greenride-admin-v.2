'use client';

import { useState, useEffect, useCallback, useMemo } from 'react';
import { GoogleMap, useJsApiLoader, Marker, InfoWindow } from '@react-google-maps/api';
import { 
  Car, 
  RefreshCw,
  Eye,
  Navigation,
  Search,
  Phone,
  Star,
  X,
  Layers,
  MapPin
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

// Google Maps API Key - from environment variable
const GOOGLE_MAPS_API_KEY = process.env.NEXT_PUBLIC_GOOGLE_MAPS_KEY || 'AIzaSyDif39v3Gx4YXonS3-A8pINUMi3hxRfC3U';

// Kigali center coordinates
const KIGALI_CENTER = { lat: -1.9403, lng: 29.8739 };

// Map container style
const containerStyle = {
  width: '100%',
  height: '100%'
};

// 30 Mock drivers spread across visible Kigali downtown area
const initialDrivers = [
  { id: 'DRV001', name: 'Peter Mutombo', status: 'available', lat: -1.9080, lng: 29.8250, vehicle: 'RAD 123A', vehicleType: 'sedan', rating: 4.8, phone: '+250788123456', heading: 45, totalRides: 1250 },
  { id: 'DRV002', name: 'David Kagame', status: 'on_trip', lat: -1.9120, lng: 29.8550, vehicle: 'RAB 456B', vehicleType: 'suv', rating: 4.6, phone: '+250788234567', heading: 180, totalRides: 890 },
  { id: 'DRV003', name: 'Jean Pierre', status: 'available', lat: -1.9050, lng: 29.8850, vehicle: 'RAC 789C', vehicleType: 'sedan', rating: 4.9, phone: '+250788345678', heading: 270, totalRides: 2100 },
  { id: 'DRV004', name: 'Emmanuel Habimana', status: 'on_trip', lat: -1.9100, lng: 29.9150, vehicle: 'RAD 012D', vehicleType: 'sedan', rating: 4.5, phone: '+250788456789', heading: 90, totalRides: 650 },
  { id: 'DRV005', name: 'Claude Uwimana', status: 'available', lat: -1.9150, lng: 29.8380, vehicle: 'RAF 678F', vehicleType: 'moto', rating: 4.7, phone: '+250788567890', heading: 135, totalRides: 430 },
  { id: 'DRV006', name: 'Alice Mukamana', status: 'on_trip', lat: -1.9220, lng: 29.8680, vehicle: 'RAG 234G', vehicleType: 'sedan', rating: 4.9, phone: '+250788678901', heading: 315, totalRides: 1890 },
  { id: 'DRV007', name: 'Patrick Niyonzima', status: 'available', lat: -1.9180, lng: 29.9020, vehicle: 'RAH 567H', vehicleType: 'suv', rating: 4.4, phone: '+250788789012', heading: 225, totalRides: 560 },
  { id: 'DRV008', name: 'Grace Ingabire', status: 'offline', lat: -1.9250, lng: 29.8200, vehicle: 'RAI 890I', vehicleType: 'sedan', rating: 4.6, phone: '+250788890123', heading: 0, totalRides: 780 },
  { id: 'DRV009', name: 'Eric Ndayisaba', status: 'available', lat: -1.9200, lng: 29.8450, vehicle: 'RAJ 123J', vehicleType: 'moto', rating: 4.8, phone: '+250788901234', heading: 60, totalRides: 320 },
  { id: 'DRV010', name: 'Marie Uwase', status: 'on_trip', lat: -1.9280, lng: 29.8780, vehicle: 'RAK 456K', vehicleType: 'sedan', rating: 4.7, phone: '+250788012345', heading: 150, totalRides: 1120 },
  { id: 'DRV011', name: 'Joseph Bizimana', status: 'available', lat: -1.9350, lng: 29.8320, vehicle: 'RAL 789L', vehicleType: 'suv', rating: 4.5, phone: '+250788123567', heading: 200, totalRides: 940 },
  { id: 'DRV012', name: 'Diane Umutoni', status: 'on_trip', lat: -1.9380, lng: 29.8600, vehicle: 'RAM 012M', vehicleType: 'sedan', rating: 4.9, phone: '+250788234678', heading: 280, totalRides: 1560 },
  { id: 'DRV013', name: 'François Habiyaremye', status: 'available', lat: -1.9320, lng: 29.8920, vehicle: 'RAN 345N', vehicleType: 'sedan', rating: 4.7, phone: '+250788345789', heading: 30, totalRides: 890 },
  { id: 'DRV014', name: 'Chantal Nyiramana', status: 'on_trip', lat: -1.9400, lng: 29.9100, vehicle: 'RAO 456O', vehicleType: 'suv', rating: 4.8, phone: '+250788456890', heading: 120, totalRides: 1340 },
  { id: 'DRV015', name: 'Innocent Nshimiye', status: 'available', lat: -1.9420, lng: 29.8480, vehicle: 'RAP 567P', vehicleType: 'moto', rating: 4.6, phone: '+250788567901', heading: 240, totalRides: 560 },
  { id: 'DRV016', name: 'Yvonne Mukeshimana', status: 'on_trip', lat: -1.9480, lng: 29.8250, vehicle: 'RAQ 678Q', vehicleType: 'sedan', rating: 4.5, phone: '+250788678012', heading: 90, totalRides: 720 },
  { id: 'DRV017', name: 'Théogène Nsengimana', status: 'available', lat: -1.9450, lng: 29.8720, vehicle: 'RAR 789R', vehicleType: 'sedan', rating: 4.9, phone: '+250788789123', heading: 180, totalRides: 2340 },
  { id: 'DRV018', name: 'Vestine Uwamahoro', status: 'offline', lat: -1.9520, lng: 29.8980, vehicle: 'RAS 890S', vehicleType: 'moto', rating: 4.4, phone: '+250788890234', heading: 270, totalRides: 430 },
  { id: 'DRV019', name: 'Olivier Mugabo', status: 'available', lat: -1.9500, lng: 29.8550, vehicle: 'RAT 901T', vehicleType: 'suv', rating: 4.8, phone: '+250788901345', heading: 45, totalRides: 1670 },
  { id: 'DRV020', name: 'Clarisse Umuhoza', status: 'on_trip', lat: -1.9550, lng: 29.9180, vehicle: 'RAU 012U', vehicleType: 'sedan', rating: 4.7, phone: '+250788012456', heading: 135, totalRides: 980 },
  { id: 'DRV021', name: 'Pacifique Niyibizi', status: 'available', lat: -1.9580, lng: 29.8380, vehicle: 'RAV 123V', vehicleType: 'sedan', rating: 4.6, phone: '+250788123678', heading: 300, totalRides: 1120 },
  { id: 'DRV022', name: 'Eugène Hakizimana', status: 'on_trip', lat: -1.9620, lng: 29.8680, vehicle: 'RAW 234W', vehicleType: 'moto', rating: 4.5, phone: '+250788234789', heading: 60, totalRides: 340 },
  { id: 'DRV023', name: 'Sylvie Nibagwire', status: 'available', lat: -1.9650, lng: 29.9050, vehicle: 'RAX 345X', vehicleType: 'sedan', rating: 4.9, phone: '+250788345890', heading: 210, totalRides: 1890 },
  { id: 'DRV024', name: 'Jean-Claude Ndayisaba', status: 'on_trip', lat: -1.9600, lng: 29.8200, vehicle: 'RAY 456Y', vehicleType: 'suv', rating: 4.8, phone: '+250788456901', heading: 150, totalRides: 2100 },
  { id: 'DRV025', name: 'Josiane Mukamana', status: 'available', lat: -1.9680, lng: 29.8850, vehicle: 'RAZ 567Z', vehicleType: 'sedan', rating: 4.7, phone: '+250788567012', heading: 330, totalRides: 1450 },
  { id: 'DRV026', name: 'Faustin Niyonsenga', status: 'on_trip', lat: -1.9720, lng: 29.8450, vehicle: 'RBA 678A', vehicleType: 'sedan', rating: 4.6, phone: '+250788678123', heading: 90, totalRides: 890 },
  { id: 'DRV027', name: 'Béatrice Uwimana', status: 'available', lat: -1.9750, lng: 29.8750, vehicle: 'RBB 789B', vehicleType: 'moto', rating: 4.8, phone: '+250788789234', heading: 180, totalRides: 560 },
  { id: 'DRV028', name: 'Alphonse Nzeyimana', status: 'offline', lat: -1.9700, lng: 29.9100, vehicle: 'RBC 890C', vehicleType: 'suv', rating: 4.5, phone: '+250788890345', heading: 270, totalRides: 1230 },
  { id: 'DRV029', name: 'Pascaline Nyirahabimana', status: 'available', lat: -1.9780, lng: 29.8580, vehicle: 'RBD 901D', vehicleType: 'sedan', rating: 4.7, phone: '+250788901456', heading: 45, totalRides: 780 },
  { id: 'DRV030', name: 'Thaddée Nkundimana', status: 'on_trip', lat: -1.9680, lng: 29.8320, vehicle: 'RBE 012E', vehicleType: 'sedan', rating: 4.9, phone: '+250788012567', heading: 135, totalRides: 1670 },
];

type Driver = typeof initialDrivers[0];

const getStatusColor = (status: string) => {
  switch (status) {
    case 'available': return '#22c55e';
    case 'on_trip': return '#eab308';
    case 'offline': return '#6b7280';
    default: return '#6b7280';
  }
};

const getStatusBadgeClass = (status: string) => {
  switch (status) {
    case 'available': return 'bg-green-100 text-green-700';
    case 'on_trip': return 'bg-yellow-100 text-yellow-700';
    case 'offline': return 'bg-gray-100 text-gray-700';
    default: return 'bg-gray-100 text-gray-700';
  }
};

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'available': return 'Available';
    case 'on_trip': return 'On Trip';
    case 'offline': return 'Offline';
    default: return status;
  }
};

// Create SVG car icon as data URL
const createCarIcon = (color: string, heading: number) => {
  const svg = `
    <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 40 40">
      <g transform="rotate(${heading}, 20, 20)">
        <circle cx="20" cy="20" r="18" fill="${color}" stroke="white" stroke-width="2"/>
        <path d="M20 8 L26 20 L24 20 L24 28 L16 28 L16 20 L14 20 Z" fill="white"/>
      </g>
    </svg>
  `;
  return `data:image/svg+xml;charset=UTF-8,${encodeURIComponent(svg)}`;
};

export default function MapPage() {
  const [drivers, setDrivers] = useState<Driver[]>(initialDrivers);
  const [simulateMovement, setSimulateMovement] = useState(true);
  const [showAvailable, setShowAvailable] = useState(true);
  const [showOnTrip, setShowOnTrip] = useState(true);
  const [showOffline, setShowOffline] = useState(false);
  const [vehicleType, setVehicleType] = useState('all');
  const [autoRefresh, setAutoRefresh] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedDriver, setSelectedDriver] = useState<Driver | null>(null);
  const [infoWindowDriver, setInfoWindowDriver] = useState<Driver | null>(null);
  const [mapType, setMapType] = useState<google.maps.MapTypeId | string>('roadmap');
  const [lastUpdated, setLastUpdated] = useState(new Date());
  const [map, setMap] = useState<google.maps.Map | null>(null);

  // Load Google Maps
  const { isLoaded, loadError } = useJsApiLoader({
    id: 'google-map-script',
    googleMapsApiKey: GOOGLE_MAPS_API_KEY,
  });

  // Map options
  const mapOptions = useMemo(() => ({
    disableDefaultUI: false,
    zoomControl: true,
    streetViewControl: false,
    mapTypeControl: false,
    fullscreenControl: true,
    mapTypeId: mapType,
  }), [mapType]);

  // On map load
  const onLoad = useCallback((map: google.maps.Map) => {
    setMap(map);
  }, []);

  // On map unmount
  const onUnmount = useCallback(() => {
    setMap(null);
  }, []);

  // Filter drivers
  const filteredDrivers = drivers.filter((driver) => {
    const matchesSearch = driver.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                          driver.vehicle.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesStatus = 
      (showAvailable && driver.status === 'available') ||
      (showOnTrip && driver.status === 'on_trip') ||
      (showOffline && driver.status === 'offline');
    const matchesVehicle = vehicleType === 'all' || driver.vehicleType === vehicleType;
    return matchesSearch && matchesStatus && matchesVehicle;
  });

  // Stats
  const stats = {
    total: drivers.length,
    available: drivers.filter(d => d.status === 'available').length,
    onTrip: drivers.filter(d => d.status === 'on_trip').length,
    offline: drivers.filter(d => d.status === 'offline').length,
  };

  // Simulated movement
  useEffect(() => {
    if (!simulateMovement) return;
    
    const moveInterval = setInterval(() => {
      setDrivers(prevDrivers => 
        prevDrivers.map(driver => {
          if (driver.status === 'offline') return driver;
          
          const latDelta = (Math.random() - 0.5) * 0.004;
          const lngDelta = (Math.random() - 0.5) * 0.004;
          const headingDelta = (Math.random() - 0.5) * 30;
          let newHeading = driver.heading + headingDelta;
          if (newHeading < 0) newHeading += 360;
          if (newHeading >= 360) newHeading -= 360;
          
          return {
            ...driver,
            lat: driver.lat + latDelta,
            lng: driver.lng + lngDelta,
            heading: Math.round(newHeading),
          };
        })
      );
      setLastUpdated(new Date());
    }, 2000);
    
    return () => clearInterval(moveInterval);
  }, [simulateMovement]);

  // Auto-refresh
  useEffect(() => {
    if (!autoRefresh) return;
    const interval = setInterval(() => setLastUpdated(new Date()), 10000);
    return () => clearInterval(interval);
  }, [autoRefresh]);

  // Focus on driver
  const focusOnDriver = useCallback((driver: Driver) => {
    setSelectedDriver(driver);
    if (map) {
      map.panTo({ lat: driver.lat, lng: driver.lng });
      map.setZoom(16);
    }
  }, [map]);

  // Handle marker click
  const handleMarkerClick = (driver: Driver) => {
    setInfoWindowDriver(driver);
    setSelectedDriver(driver);
  };

  if (loadError) {
    return (
      <div className="h-[calc(100vh-8rem)] flex items-center justify-center bg-gray-100 rounded-lg">
        <div className="text-center">
          <p className="text-red-500 font-medium">Failed to load Google Maps</p>
          <p className="text-sm text-muted-foreground mt-2">{loadError.message}</p>
        </div>
      </div>
    );
  }

  if (!isLoaded) {
    return (
      <div className="h-[calc(100vh-8rem)] flex items-center justify-center bg-gray-100 rounded-lg">
        <div className="text-center">
          <RefreshCw className="h-8 w-8 animate-spin mx-auto text-primary" />
          <p className="mt-2 text-muted-foreground">Loading map...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="h-[calc(100vh-8rem)] relative">
      {/* Map Container */}
      <div className="absolute inset-0 rounded-lg overflow-hidden">
        <GoogleMap
          mapContainerStyle={containerStyle}
          center={KIGALI_CENTER}
          zoom={13}
          onLoad={onLoad}
          onUnmount={onUnmount}
          options={mapOptions}
        >
          {/* Driver Markers */}
          {filteredDrivers.map((driver) => (
            <Marker
              key={driver.id}
              position={{ lat: driver.lat, lng: driver.lng }}
              icon={{
                url: createCarIcon(getStatusColor(driver.status), driver.heading),
                scaledSize: new google.maps.Size(40, 40),
                anchor: new google.maps.Point(20, 20),
              }}
              onClick={() => handleMarkerClick(driver)}
              title={driver.name}
            />
          ))}

          {/* Info Window */}
          {infoWindowDriver && (
            <InfoWindow
              position={{ lat: infoWindowDriver.lat, lng: infoWindowDriver.lng }}
              onCloseClick={() => setInfoWindowDriver(null)}
            >
              <div className="p-2 min-w-[150px]">
                <p className="font-semibold">{infoWindowDriver.name}</p>
                <p className="text-sm text-gray-600">{infoWindowDriver.vehicle}</p>
                <div className="flex items-center gap-1 mt-1">
                  <span className="text-xs px-2 py-0.5 rounded" style={{ 
                    backgroundColor: getStatusColor(infoWindowDriver.status) + '20',
                    color: getStatusColor(infoWindowDriver.status)
                  }}>
                    {getStatusLabel(infoWindowDriver.status)}
                  </span>
                  <span className="text-xs">⭐ {infoWindowDriver.rating}</span>
                </div>
              </div>
            </InfoWindow>
          )}
        </GoogleMap>
      </div>

      {/* Map Type Controls */}
      <div className="absolute top-4 left-4 flex gap-2 z-10">
        <Button
          variant={mapType === 'roadmap' ? 'default' : 'secondary'}
          size="sm"
          onClick={() => setMapType('roadmap')}
          className="shadow-lg"
        >
          <MapPin className="h-4 w-4 mr-1" />
          Map
        </Button>
        <Button
          variant={mapType === 'satellite' ? 'default' : 'secondary'}
          size="sm"
          onClick={() => setMapType('satellite')}
          className="shadow-lg"
        >
          Satellite
        </Button>
        <Button
          variant={mapType === 'hybrid' ? 'default' : 'secondary'}
          size="sm"
          onClick={() => setMapType('hybrid')}
          className="shadow-lg"
        >
          Hybrid
        </Button>
      </div>

      {/* Control Panel */}
      <Card className="absolute top-4 right-4 w-80 shadow-xl max-h-[calc(100vh-12rem)] overflow-y-auto z-10">
        <CardContent className="p-4 space-y-4">
          <div className="flex items-center justify-between">
            <div className="font-semibold flex items-center gap-2">
              <Layers className="h-4 w-4" />
              Fleet Tracker
            </div>
            <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => setLastUpdated(new Date())}>
              <RefreshCw className={`h-4 w-4 ${autoRefresh ? 'animate-spin' : ''}`} />
            </Button>
          </div>
          
          {/* Search */}
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Find driver or vehicle..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-9"
            />
          </div>

          {/* Quick Stats */}
          <div className="grid grid-cols-3 gap-2 text-center">
            <div className="bg-green-50 rounded-lg p-2">
              <p className="text-lg font-bold text-green-600">{stats.available}</p>
              <p className="text-xs text-muted-foreground">Available</p>
            </div>
            <div className="bg-yellow-50 rounded-lg p-2">
              <p className="text-lg font-bold text-yellow-600">{stats.onTrip}</p>
              <p className="text-xs text-muted-foreground">On Trip</p>
            </div>
            <div className="bg-gray-50 rounded-lg p-2">
              <p className="text-lg font-bold text-gray-600">{stats.offline}</p>
              <p className="text-xs text-muted-foreground">Offline</p>
            </div>
          </div>

          {/* Filters */}
          <div className="space-y-2">
            <p className="text-sm font-medium text-muted-foreground">Filter by Status</p>
            <div className="flex items-center space-x-2">
              <Checkbox id="show-available" checked={showAvailable} onCheckedChange={(c) => setShowAvailable(!!c)} />
              <Label htmlFor="show-available" className="flex items-center gap-2 cursor-pointer">
                <span className="h-3 w-3 rounded-full bg-green-500" />
                Available ({stats.available})
              </Label>
            </div>
            <div className="flex items-center space-x-2">
              <Checkbox id="show-on-trip" checked={showOnTrip} onCheckedChange={(c) => setShowOnTrip(!!c)} />
              <Label htmlFor="show-on-trip" className="flex items-center gap-2 cursor-pointer">
                <span className="h-3 w-3 rounded-full bg-yellow-500" />
                On Trip ({stats.onTrip})
              </Label>
            </div>
            <div className="flex items-center space-x-2">
              <Checkbox id="show-offline" checked={showOffline} onCheckedChange={(c) => setShowOffline(!!c)} />
              <Label htmlFor="show-offline" className="flex items-center gap-2 cursor-pointer">
                <span className="h-3 w-3 rounded-full bg-gray-400" />
                Offline ({stats.offline})
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
              <SelectItem value="sedan">🚗 Sedan</SelectItem>
              <SelectItem value="suv">🚙 SUV</SelectItem>
              <SelectItem value="moto">🏍️ Moto</SelectItem>
            </SelectContent>
          </Select>

          {/* Options */}
          <div className="flex items-center justify-between pt-2 border-t">
            <div className="flex items-center space-x-2">
              <Checkbox id="auto-refresh" checked={autoRefresh} onCheckedChange={(c) => setAutoRefresh(!!c)} />
              <Label htmlFor="auto-refresh" className="cursor-pointer">Auto-refresh</Label>
            </div>
            <div className="flex items-center space-x-2">
              <Checkbox id="simulate-movement" checked={simulateMovement} onCheckedChange={(c) => setSimulateMovement(!!c)} />
              <Label htmlFor="simulate-movement" className="cursor-pointer">🚗 Simulate</Label>
            </div>
          </div>

          {/* Driver List */}
          <div className="space-y-2 pt-2 border-t">
            <p className="text-sm font-medium text-muted-foreground">Drivers ({filteredDrivers.length})</p>
            <div className="max-h-48 overflow-y-auto space-y-2">
              {filteredDrivers.map((driver) => (
                <div
                  key={driver.id}
                  className={`p-2 rounded-lg border cursor-pointer transition-colors hover:bg-accent ${selectedDriver?.id === driver.id ? 'border-primary bg-accent' : ''}`}
                  onClick={() => focusOnDriver(driver)}
                >
                  <div className="flex items-center gap-2">
                    <div className="w-3 h-3 rounded-full" style={{ backgroundColor: getStatusColor(driver.status) }} />
                    <span className="font-medium text-sm flex-1">{driver.name}</span>
                    <span className="text-xs text-muted-foreground">{driver.vehicle}</span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Stats Panel */}
      <Card className="absolute bottom-4 left-1/2 -translate-x-1/2 shadow-xl z-10">
        <CardContent className="p-3">
          <div className="flex items-center gap-6">
            <div className="flex items-center gap-2">
              <span className="h-3 w-3 rounded-full bg-green-500 animate-pulse" />
              <span className="text-sm font-medium">{stats.available + stats.onTrip} online</span>
            </div>
            <div className="flex items-center gap-2">
              <Car className="h-4 w-4 text-yellow-500" />
              <span className="text-sm font-medium">{stats.onTrip} active</span>
            </div>
            <div className="text-xs text-muted-foreground">
              {lastUpdated.toLocaleTimeString()}
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Selected Driver Panel */}
      {selectedDriver && (
        <Card className="absolute bottom-4 left-4 w-80 shadow-xl z-10">
          <CardContent className="p-4">
            <div className="flex items-start justify-between">
              <div className="flex items-center gap-3">
                <Avatar className="h-14 w-14 border-2 border-white shadow">
                  <AvatarFallback className="bg-primary text-primary-foreground text-lg">
                    {selectedDriver.name.split(' ').map(n => n[0]).join('')}
                  </AvatarFallback>
                </Avatar>
                <div>
                  <p className="font-semibold">{selectedDriver.name}</p>
                  <p className="text-sm text-muted-foreground">{selectedDriver.vehicle}</p>
                  <div className="flex items-center gap-1 mt-1">
                    <Star className="h-3 w-3 fill-yellow-400 text-yellow-400" />
                    <span className="text-sm font-medium">{selectedDriver.rating}</span>
                    <span className="text-xs text-muted-foreground">• {selectedDriver.totalRides.toLocaleString()} rides</span>
                  </div>
                </div>
              </div>
              <Button variant="ghost" size="icon" className="h-6 w-6" onClick={() => setSelectedDriver(null)}>
                <X className="h-4 w-4" />
              </Button>
            </div>
            
            <div className="flex items-center gap-2 mt-3">
              <Badge className={getStatusBadgeClass(selectedDriver.status)}>{getStatusLabel(selectedDriver.status)}</Badge>
              <Badge variant="outline" className="capitalize">
                {selectedDriver.vehicleType === 'moto' ? '🏍️' : selectedDriver.vehicleType === 'suv' ? '🚙' : '🚗'} {selectedDriver.vehicleType}
              </Badge>
            </div>
            
            <div className="flex gap-2 mt-4">
              <Button size="sm" className="flex-1" asChild>
                <a href={`/drivers/${selectedDriver.id}`}>
                  <Eye className="h-4 w-4 mr-1" />
                  Profile
                </a>
              </Button>
              <Button size="sm" variant="outline" asChild>
                <a href={`tel:${selectedDriver.phone}`}>
                  <Phone className="h-4 w-4" />
                </a>
              </Button>
              <Button size="sm" variant="outline">
                <Navigation className="h-4 w-4" />
              </Button>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}

