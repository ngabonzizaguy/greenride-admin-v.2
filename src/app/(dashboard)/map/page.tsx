'use client';

import { useState, useEffect, useCallback, useMemo, useRef } from 'react';
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
  MapPin,
  AlertCircle,
  Wifi,
  WifiOff,
  ChevronLeft,
  ChevronRight,
  PanelRightClose,
  PanelRightOpen
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
import { apiClient, NearbyDriverLocation } from '@/lib/api-client';

// Google Maps API Key - from environment variable
const GOOGLE_MAPS_API_KEY = process.env.NEXT_PUBLIC_GOOGLE_MAPS_KEY || '';
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || '';

// Kigali center coordinates
const KIGALI_CENTER = { lat: -1.9403, lng: 29.8739 };

// Map container style
const containerStyle = {
  width: '100%',
  height: '100%'
};

// Polling interval for real-time updates (in ms)
const POLLING_INTERVAL = 5000; // 5 seconds

// Driver interface extended for map display
interface MapDriver extends NearbyDriverLocation {
  status: 'available' | 'on_trip' | 'offline';
}

const getStatusFromDriver = (driver: NearbyDriverLocation): 'available' | 'on_trip' | 'offline' => {
  if (!driver.is_online && !driver.is_busy) return 'offline';
  if (driver.is_busy) return 'on_trip';
  return 'available';
};

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

function coerceFiniteNumber(value: unknown): number {
  const n =
    typeof value === 'number'
      ? value
      : typeof value === 'string'
        ? Number(value.trim())
        : Number(value);
  return Number.isFinite(n) ? n : Number.NaN;
}

function resolvePhotoUrl(raw: unknown): string {
  const s = typeof raw === 'string' ? raw.trim() : '';
  if (!s) return '';
  if (/^(data:|blob:)/i.test(s)) return s;
  if (/^https?:\/\//i.test(s)) return s;
  if (s.startsWith('//')) {
    const proto = typeof window !== 'undefined' ? window.location.protocol : 'https:';
    return `${proto}${s}`;
  }

  const base =
    API_BASE_URL ||
    (typeof window !== 'undefined' ? window.location.origin : '');

  try {
    // Handles both "/uploads/..." and "uploads/..."
    return new URL(s, base.endsWith('/') ? base : `${base}/`).toString();
  } catch {
    return '';
  }
}

function escapeXmlAttr(value: string): string {
  return value
    .replace(/&/g, '&amp;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&apos;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;');
}

// Create vehicle-type specific icons with gradients and animations
const createVehicleIcon = (
  color: string, 
  heading: number, 
  vehicleCategory?: string,
  isMoving: boolean = false
) => {
  const category = vehicleCategory?.toLowerCase() || 'sedan';
  const size = 48;
  const center = size / 2;
  
  // Vehicle-specific paths
  const vehiclePaths: Record<string, string> = {
    sedan: 'M12 22 L20 10 L28 22 L25 22 L25 32 L15 32 L15 22 Z',
    suv: 'M11 24 L20 6 L29 24 L26 24 L26 34 L14 34 L14 24 Z',
    mpv: 'M10 26 L20 4 L30 26 L27 26 L27 36 L13 36 L13 26 Z',
    van: 'M8 28 L20 2 L32 28 L29 28 L29 38 L11 38 L11 28 Z',
    hatchback: 'M13 20 L20 12 L27 20 L24 20 L24 30 L16 30 L16 20 Z',
    moto: 'M18 30 L20 8 L22 30 M16 28 L24 28 L24 32 L16 32 Z',
  };
  
  const vehiclePath = vehiclePaths[category] || vehiclePaths.sedan;
  
  // Pulsing effect for moving vehicles
  const pulseAnimation = isMoving ? `
    <animate attributeName="r" values="20;22;20" dur="1.5s" repeatCount="indefinite"/>
  ` : '';
  
  // Gradient definition
  const gradientId = `grad-${color.replace('#', '')}`;
  
  const svg = `
    <svg xmlns="http://www.w3.org/2000/svg" width="${size}" height="${size}" viewBox="0 0 ${size} ${size}">
      <defs>
        <radialGradient id="${gradientId}" cx="50%" cy="50%" r="50%">
          <stop offset="0%" style="stop-color:${color};stop-opacity:1" />
          <stop offset="100%" style="stop-color:${color};stop-opacity:0.4" />
        </radialGradient>
        <filter id="glow-${gradientId}">
          <feGaussianBlur stdDeviation="3" result="coloredBlur"/>
          <feMerge>
            <feMergeNode in="coloredBlur"/>
            <feMergeNode in="SourceGraphic"/>
          </feMerge>
        </filter>
      </defs>
      <g transform="rotate(${heading}, ${center}, ${center})">
        <!-- Outer glow circle (pulsing if moving) -->
        <circle cx="${center}" cy="${center}" r="22" fill="url(#${gradientId})" opacity="0.3">
          ${pulseAnimation}
        </circle>
        <!-- Main circle with border -->
        <circle cx="${center}" cy="${center}" r="18" fill="${color}" stroke="white" stroke-width="2.5" filter="url(#glow-${gradientId})"/>
        <!-- Vehicle silhouette -->
        <path d="${vehiclePath}" fill="white" stroke="${color}" stroke-width="0.5" opacity="0.95"/>
        <!-- Status indicator dot -->
        <circle cx="${center}" cy="${center - 12}" r="3" fill="white" stroke="${color}" stroke-width="1"/>
      </g>
    </svg>
  `;
  return `data:image/svg+xml;charset=UTF-8,${encodeURIComponent(svg)}`;
};

// Create avatar marker icon (driver photo clipped in a circle)
const createAvatarIcon = (photoUrl: string, ringColor: string, heading: number, isMoving: boolean = false) => {
  const size = 48;
  const center = size / 2;
  const inner = 32;
  const pad = (size - inner) / 2;
  const pulseAnimation = isMoving ? `
    <animate attributeName="opacity" values="0.22;0.38;0.22" dur="1.5s" repeatCount="indefinite"/>
  ` : '';
  // IMPORTANT: The URL is interpolated into XML before encoding the SVG.
  // If it contains characters like `&` or `"`, it can break the SVG and make the marker disappear.
  const safePhotoUrl = escapeXmlAttr(photoUrl);
  const svg = `
    <svg xmlns="http://www.w3.org/2000/svg" width="${size}" height="${size}" viewBox="0 0 ${size} ${size}">
      <defs>
        <clipPath id="avatar-clip">
          <circle cx="${center}" cy="${center}" r="${inner / 2}" />
        </clipPath>
      </defs>
      <g transform="rotate(${heading}, ${center}, ${center})">
        <circle cx="${center}" cy="${center}" r="22" fill="${ringColor}" opacity="0.25">
          ${pulseAnimation}
        </circle>
        <circle cx="${center}" cy="${center}" r="18" fill="${ringColor}" stroke="white" stroke-width="2.5"/>
        <image href="${safePhotoUrl}" x="${pad}" y="${pad}" width="${inner}" height="${inner}" clip-path="url(#avatar-clip)" preserveAspectRatio="xMidYMid slice"/>
      </g>
    </svg>
  `;
  return `data:image/svg+xml;charset=UTF-8,${encodeURIComponent(svg)}`;
};

// Legacy function for backward compatibility
const createCarIcon = (color: string, heading: number) => {
  return createVehicleIcon(color, heading, 'sedan', false);
};

export default function MapPage() {
  const [drivers, setDrivers] = useState<MapDriver[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isLiveConnected, setIsLiveConnected] = useState(false);
  
  // Filters
  const [showAvailable, setShowAvailable] = useState(true);
  const [showOnTrip, setShowOnTrip] = useState(true);
  const [showOffline, setShowOffline] = useState(false);
  const [vehicleType, setVehicleType] = useState('all');
  const [searchQuery, setSearchQuery] = useState('');
  
  // Map controls
  const [autoRefresh, setAutoRefresh] = useState(true);
  const [selectedDriver, setSelectedDriver] = useState<MapDriver | null>(null);
  const [infoWindowDriver, setInfoWindowDriver] = useState<MapDriver | null>(null);
  const [mapType, setMapType] = useState<google.maps.MapTypeId | string>('roadmap');
  const [lastUpdated, setLastUpdated] = useState(new Date());
  const [map, setMap] = useState<google.maps.Map | null>(null);
  const [isPanelCollapsed, setIsPanelCollapsed] = useState(false);

  // Refs for cleanup
  const pollingRef = useRef<NodeJS.Timeout | null>(null);
  const isFirstLoad = useRef(true);

  // Load Google Maps
  const { isLoaded, loadError } = useJsApiLoader({
    id: 'google-map-script',
    googleMapsApiKey: GOOGLE_MAPS_API_KEY,
    libraries: ['places'],
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

  // Fetch drivers from API
  const fetchDrivers = useCallback(async (showLoadingState = false) => {
    if (showLoadingState) {
      setIsLoading(true);
    }
    
    try {
      const response = await apiClient.getDriversWithLocations({
        latitude: KIGALI_CENTER.lat,
        longitude: KIGALI_CENTER.lng,
        limit: 100,
      });

      if (response.code === '0000' && response.data) {
        const mapDrivers: MapDriver[] = response.data.drivers.map((driver) => {
          // Backend sometimes returns coords as strings; coerce to finite numbers so markers don't vanish.
          const latitude = coerceFiniteNumber((driver as any).latitude);
          const longitude = coerceFiniteNumber((driver as any).longitude);
          const normalized = {
            ...driver,
            latitude,
            longitude,
            photo_url: resolvePhotoUrl((driver as any).photo_url),
          };
          return {
            ...(normalized as any),
            status: getStatusFromDriver(normalized as any),
          } as MapDriver;
        });
        
        setDrivers(mapDrivers);
        setLastUpdated(new Date());
        setIsLiveConnected(true);
        setError(null);
        
        if (isFirstLoad.current) {
          isFirstLoad.current = false;
          console.log('[Map] Initial load complete:', mapDrivers.length, 'drivers');
        }
      } else {
        throw new Error(response.msg || 'Failed to fetch drivers');
      }
    } catch (err) {
      console.error('[Map] Failed to fetch drivers:', err);
      setError(err instanceof Error ? err.message : 'Failed to load drivers');
      setIsLiveConnected(false);
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Initial load
  useEffect(() => {
    fetchDrivers(true);
  }, [fetchDrivers]);

  // Polling for real-time updates
  useEffect(() => {
    if (!autoRefresh) {
      if (pollingRef.current) {
        clearInterval(pollingRef.current);
        pollingRef.current = null;
      }
      return;
    }

    pollingRef.current = setInterval(() => {
      fetchDrivers(false);
    }, POLLING_INTERVAL);

    return () => {
      if (pollingRef.current) {
        clearInterval(pollingRef.current);
        pollingRef.current = null;
      }
    };
  }, [autoRefresh, fetchDrivers]);

  // Filter drivers
  const filteredDrivers = useMemo(() => {
    return drivers.filter((driver) => {
      const matchesSearch = driver.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                            (driver.plate_number?.toLowerCase().includes(searchQuery.toLowerCase()) ?? false);
      const matchesStatus = 
        (showAvailable && driver.status === 'available') ||
        (showOnTrip && driver.status === 'on_trip') ||
        (showOffline && driver.status === 'offline');
      const matchesVehicle = vehicleType === 'all' || driver.vehicle_category === vehicleType;
      return matchesSearch && matchesStatus && matchesVehicle;
    });
  }, [drivers, searchQuery, showAvailable, showOnTrip, showOffline, vehicleType]);

  const hasValidCoords = (d: Pick<MapDriver, 'latitude' | 'longitude'>) =>
    Number.isFinite(d.latitude) &&
    Number.isFinite(d.longitude) &&
    Math.abs(d.latitude) <= 90 &&
    Math.abs(d.longitude) <= 180;

  const markerDrivers = useMemo(
    () => filteredDrivers.filter(hasValidCoords),
    [filteredDrivers]
  );

  // Stats
  const stats = useMemo(() => ({
    total: drivers.length,
    available: drivers.filter(d => d.status === 'available').length,
    onTrip: drivers.filter(d => d.status === 'on_trip').length,
    offline: drivers.filter(d => d.status === 'offline').length,
  }), [drivers]);

  // Focus on driver
  const focusOnDriver = useCallback((driver: MapDriver) => {
    setSelectedDriver(driver);
    if (map) {
      if (!hasValidCoords(driver)) {
        console.warn('[Map] Driver has no valid coordinates:', {
          driver_id: driver.driver_id,
          latitude: driver.latitude,
          longitude: driver.longitude,
        });
        return;
      }
      map.panTo({ lat: driver.latitude, lng: driver.longitude });
      map.setZoom(16);
    }
  }, [map]);

  // Handle marker click
  const handleMarkerClick = (driver: MapDriver) => {
    setInfoWindowDriver(driver);
    setSelectedDriver(driver);
  };

  // Manual refresh
  const handleManualRefresh = () => {
    fetchDrivers(true);
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
          {markerDrivers.map((driver) => {
            const isMoving = driver.status === 'available' || driver.status === 'on_trip';
            const photoUrl = resolvePhotoUrl(driver.photo_url);
            const iconUrl = photoUrl
              ? createAvatarIcon(photoUrl, getStatusColor(driver.status), driver.heading || 0, isMoving)
              : createVehicleIcon(
                  getStatusColor(driver.status),
                  driver.heading || 0,
                  driver.vehicle_category,
                  isMoving
                );
            return (
              <Marker
                key={driver.driver_id}
                position={{ lat: driver.latitude, lng: driver.longitude }}
                icon={{
                  url: iconUrl,
                  scaledSize: new google.maps.Size(48, 48),
                  anchor: new google.maps.Point(24, 24),
                }}
                onClick={() => handleMarkerClick(driver)}
                title={driver.name}
              />
            );
          })}

          {/* Info Window */}
          {infoWindowDriver && (
            <InfoWindow
              position={{ lat: infoWindowDriver.latitude, lng: infoWindowDriver.longitude }}
              onCloseClick={() => setInfoWindowDriver(null)}
            >
              <div className="p-2 min-w-[150px]">
                <p className="font-semibold">{infoWindowDriver.name}</p>
                <p className="text-sm text-gray-600">{infoWindowDriver.plate_number || 'No plate'}</p>
                <div className="flex items-center gap-1 mt-1">
                  <span className="text-xs px-2 py-0.5 rounded" style={{ 
                    backgroundColor: getStatusColor(infoWindowDriver.status) + '20',
                    color: getStatusColor(infoWindowDriver.status)
                  }}>
                    {getStatusLabel(infoWindowDriver.status)}
                  </span>
                  {infoWindowDriver.rating && (
                    <span className="text-xs">⭐ {infoWindowDriver.rating.toFixed(1)}</span>
                  )}
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

      {/* Live Connection Status */}
      <div className="absolute top-4 left-1/2 -translate-x-1/2 z-10">
        <Badge 
          variant={isLiveConnected ? 'default' : 'destructive'} 
          className={`shadow-lg ${isLiveConnected ? 'bg-green-600' : ''}`}
        >
          {isLiveConnected ? (
            <>
              <Wifi className="h-3 w-3 mr-1" />
              Live Data
            </>
          ) : (
            <>
              <WifiOff className="h-3 w-3 mr-1" />
              Disconnected
            </>
          )}
        </Badge>
      </div>

      {/* Error Banner */}
      {error && (
        <div className="absolute top-16 left-1/2 -translate-x-1/2 z-10">
          <div className="bg-red-100 border border-red-300 text-red-700 px-4 py-2 rounded-lg shadow-lg flex items-center gap-2">
            <AlertCircle className="h-4 w-4" />
            <span className="text-sm">{error}</span>
            <Button variant="ghost" size="sm" onClick={handleManualRefresh}>
              Retry
            </Button>
          </div>
        </div>
      )}

      {/* Toggle Panel Button (visible when collapsed) */}
      {isPanelCollapsed && (
        <Button
          className="absolute top-4 right-4 z-10 shadow-xl"
          variant="secondary"
          size="icon"
          onClick={() => setIsPanelCollapsed(false)}
        >
          <PanelRightOpen className="h-5 w-5" />
        </Button>
      )}

      {/* Control Panel */}
      <Card className={`absolute top-4 right-4 shadow-xl max-h-[calc(100vh-12rem)] overflow-y-auto z-10 transition-all duration-300 ${isPanelCollapsed ? 'translate-x-[400px] opacity-0 pointer-events-none' : 'translate-x-0 opacity-100'} w-80`}>
        <CardContent className="p-4 space-y-4">
          <div className="flex items-center justify-between">
            <div className="font-semibold flex items-center gap-2">
              <Layers className="h-4 w-4" />
              Fleet Tracker
            </div>
            <div className="flex items-center gap-1">
              <Button 
                variant="ghost" 
                size="icon" 
                className="h-8 w-8" 
                onClick={handleManualRefresh}
                disabled={isLoading}
                title="Refresh"
              >
                <RefreshCw className={`h-4 w-4 ${isLoading ? 'animate-spin' : ''}`} />
              </Button>
              <Button 
                variant="ghost" 
                size="icon" 
                className="h-8 w-8" 
                onClick={() => setIsPanelCollapsed(true)}
                title="Collapse panel"
              >
                <PanelRightClose className="h-4 w-4" />
              </Button>
            </div>
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

          {/* Loading State */}
          {isLoading && drivers.length === 0 && (
            <div className="text-center py-4">
              <RefreshCw className="h-6 w-6 animate-spin mx-auto text-primary" />
              <p className="text-sm text-muted-foreground mt-2">Loading drivers...</p>
            </div>
          )}

          {/* Quick Stats */}
          {!isLoading && (
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
          )}

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
              <Label htmlFor="auto-refresh" className="cursor-pointer text-sm">
                Live updates ({POLLING_INTERVAL / 1000}s)
              </Label>
            </div>
          </div>

          {/* Driver List */}
          <div className="space-y-2 pt-2 border-t">
            <p className="text-sm font-medium text-muted-foreground">Drivers ({filteredDrivers.length})</p>
            <div className="max-h-48 overflow-y-auto space-y-2">
              {filteredDrivers.length === 0 && !isLoading && (
                <p className="text-sm text-muted-foreground text-center py-4">
                  No drivers match your filters
                </p>
              )}
              {filteredDrivers.map((driver) => (
                <div
                  key={driver.driver_id}
                  className={`p-2 rounded-lg border cursor-pointer transition-colors hover:bg-accent ${selectedDriver?.driver_id === driver.driver_id ? 'border-primary bg-accent' : ''}`}
                  onClick={() => focusOnDriver(driver)}
                >
                  <div className="flex items-center gap-2">
                    <div className="w-3 h-3 rounded-full" style={{ backgroundColor: getStatusColor(driver.status) }} />
                    <span className="font-medium text-sm flex-1 truncate">{driver.name}</span>
                    <span className="text-xs text-muted-foreground">
                      {!hasValidCoords(driver) ? 'No GPS' : (driver.plate_number || '—')}
                    </span>
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
              <span className={`h-3 w-3 rounded-full ${isLiveConnected ? 'bg-green-500 animate-pulse' : 'bg-gray-400'}`} />
              <span className="text-sm font-medium">{stats.available + stats.onTrip} online</span>
            </div>
            <div className="flex items-center gap-2">
              <Car className="h-4 w-4 text-yellow-500" />
              <span className="text-sm font-medium">{stats.onTrip} active</span>
            </div>
            <div className="text-xs text-muted-foreground">
              Updated: {lastUpdated.toLocaleTimeString()}
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
                  <p className="text-sm text-muted-foreground">{selectedDriver.plate_number || 'No plate'}</p>
                  <div className="flex items-center gap-1 mt-1">
                    {selectedDriver.rating && (
                      <>
                        <Star className="h-3 w-3 fill-yellow-400 text-yellow-400" />
                        <span className="text-sm font-medium">{selectedDriver.rating.toFixed(1)}</span>
                      </>
                    )}
                    {selectedDriver.total_rides && (
                      <span className="text-xs text-muted-foreground">• {selectedDriver.total_rides.toLocaleString()} rides</span>
                    )}
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
                {selectedDriver.vehicle_category === 'moto' ? '🏍️' : selectedDriver.vehicle_category === 'suv' ? '🚙' : '🚗'} {selectedDriver.vehicle_category || 'Unknown'}
              </Badge>
            </div>
            
            <div className="flex gap-2 mt-4">
              <Button size="sm" className="flex-1" asChild>
                <a href={`/drivers/${selectedDriver.driver_id}`}>
                  <Eye className="h-4 w-4 mr-1" />
                  Profile
                </a>
              </Button>
              {selectedDriver.phone && (
                <Button size="sm" variant="outline" asChild>
                  <a href={`tel:${selectedDriver.phone}`}>
                    <Phone className="h-4 w-4" />
                  </a>
                </Button>
              )}
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
