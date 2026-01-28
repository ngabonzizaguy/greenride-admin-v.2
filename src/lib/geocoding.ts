/**
 * Geocoding service with hybrid approach:
 * - Primary: Google Maps Geocoding API (if API key available)
 * - Fallback: OpenStreetMap Nominatim (free, no key needed)
 * - Caching: localStorage to reduce API calls
 */

export interface GeocodeResult {
  lat: number;
  lng: number;
  address: string;
  formatted_address?: string;
}

interface CachedResult {
  result: GeocodeResult;
  timestamp: number;
}

const CACHE_DURATION = 24 * 60 * 60 * 1000; // 24 hours
const CACHE_KEY_PREFIX = 'geocode_cache_';
const REVERSE_CACHE_KEY_PREFIX = 'reverse_geocode_cache_';

/**
 * Get cached geocode result if available and not expired
 */
function getCachedResult(address: string): GeocodeResult | null {
  if (typeof window === 'undefined') return null;
  
  try {
    const cacheKey = CACHE_KEY_PREFIX + address.toLowerCase().trim();
    const cached = localStorage.getItem(cacheKey);
    if (!cached) return null;

    const data: CachedResult = JSON.parse(cached);
    const now = Date.now();
    
    // Check if cache is still valid
    if (now - data.timestamp < CACHE_DURATION) {
      return data.result;
    }
    
    // Cache expired, remove it
    localStorage.removeItem(cacheKey);
    return null;
  } catch (error) {
    console.error('Error reading geocode cache:', error);
    return null;
  }
}

/**
 * Cache geocode result
 */
function cacheResult(address: string, result: GeocodeResult): void {
  if (typeof window === 'undefined') return;
  
  try {
    const cacheKey = CACHE_KEY_PREFIX + address.toLowerCase().trim();
    const data: CachedResult = {
      result,
      timestamp: Date.now(),
    };
    localStorage.setItem(cacheKey, JSON.stringify(data));
  } catch (error) {
    console.error('Error caching geocode result:', error);
  }
}

/**
 * Geocode using Google Maps Geocoding API
 */
async function geocodeWithGoogle(address: string): Promise<GeocodeResult | null> {
  const apiKey = process.env.NEXT_PUBLIC_GOOGLE_MAPS_KEY;
  if (!apiKey) return null;

  try {
    const url = `https://maps.googleapis.com/maps/api/geocode/json?address=${encodeURIComponent(address)}&key=${apiKey}`;
    const response = await fetch(url);
    
    if (!response.ok) {
      throw new Error(`Google Geocoding API error: ${response.status}`);
    }

    const data = await response.json();
    
    if (data.status === 'OK' && data.results && data.results.length > 0) {
      const result = data.results[0];
      const location = result.geometry.location;
      
      return {
        lat: location.lat,
        lng: location.lng,
        address,
        formatted_address: result.formatted_address,
      };
    }
    
    return null;
  } catch (error) {
    console.error('Google Geocoding failed:', error);
    return null;
  }
}

/**
 * Geocode using OpenStreetMap Nominatim (fallback)
 */
async function geocodeWithNominatim(address: string): Promise<GeocodeResult | null> {
  try {
    // Add Rwanda context for better results
    const searchQuery = address.includes('Rwanda') || address.includes('Kigali') 
      ? address 
      : `${address}, Kigali, Rwanda`;
    
    const url = `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(searchQuery)}&limit=1`;
    
    const response = await fetch(url, {
      headers: {
        'User-Agent': 'GreenRide Admin App', // Required by Nominatim
      },
    });
    
    if (!response.ok) {
      throw new Error(`Nominatim API error: ${response.status}`);
    }

    const data = await response.json();
    
    if (Array.isArray(data) && data.length > 0) {
      const result = data[0];
      
      return {
        lat: parseFloat(result.lat),
        lng: parseFloat(result.lon),
        address,
        formatted_address: result.display_name,
      };
    }
    
    return null;
  } catch (error) {
    console.error('Nominatim Geocoding failed:', error);
    return null;
  }
}

/**
 * Geocode an address to coordinates
 * Uses hybrid approach: Google first, then Nominatim fallback
 * Results are cached in localStorage
 */
export async function geocodeAddress(address: string): Promise<GeocodeResult | null> {
  if (!address || address.trim().length === 0) {
    return null;
  }

  const normalizedAddress = address.trim();
  
  // Check cache first
  const cached = getCachedResult(normalizedAddress);
  if (cached) {
    return cached;
  }

  // Try Google Maps first (if API key available)
  let result = await geocodeWithGoogle(normalizedAddress);
  
  // Fallback to Nominatim if Google fails or no API key
  if (!result) {
    result = await geocodeWithNominatim(normalizedAddress);
  }

  // Cache the result if we got one
  if (result) {
    cacheResult(normalizedAddress, result);
  }

  return result;
}

/**
 * Reverse geocode coordinates to a human-readable address.
 * Uses Google first (if key available), then Nominatim fallback.
 * Result is cached in localStorage.
 */
export async function reverseGeocode(lat: number, lng: number): Promise<string | null> {
  if (!Number.isFinite(lat) || !Number.isFinite(lng)) return null;
  if (typeof window === 'undefined') return null;

  const cacheKey = `${REVERSE_CACHE_KEY_PREFIX}${lat.toFixed(5)},${lng.toFixed(5)}`;
  try {
    const cached = localStorage.getItem(cacheKey);
    if (cached) {
      const data: { result: string; timestamp: number } = JSON.parse(cached);
      if (Date.now() - data.timestamp < CACHE_DURATION) return data.result;
      localStorage.removeItem(cacheKey);
    }
  } catch {
    // ignore cache read errors
  }

  const apiKey = process.env.NEXT_PUBLIC_GOOGLE_MAPS_KEY;
  try {
    if (apiKey) {
      const url = `https://maps.googleapis.com/maps/api/geocode/json?latlng=${lat},${lng}&key=${apiKey}`;
      const response = await fetch(url);
      if (response.ok) {
        const data = await response.json();
        if (data.status === 'OK' && Array.isArray(data.results) && data.results.length > 0) {
          const formatted = data.results[0].formatted_address as string;
          try {
            localStorage.setItem(cacheKey, JSON.stringify({ result: formatted, timestamp: Date.now() }));
          } catch {
            // ignore cache write errors
          }
          return formatted;
        }
      }
    }
  } catch {
    // ignore and fall back
  }

  try {
    const url = `https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lng}`;
    const response = await fetch(url, {
      headers: { 'User-Agent': 'GreenRide Admin App' },
    });
    if (!response.ok) return null;
    const data = await response.json();
    const formatted = (data?.display_name as string | undefined) || null;
    if (formatted) {
      try {
        localStorage.setItem(cacheKey, JSON.stringify({ result: formatted, timestamp: Date.now() }));
      } catch {
        // ignore cache write errors
      }
    }
    return formatted;
  } catch {
    return null;
  }
}

/**
 * Popular locations in Kigali with pre-defined coordinates
 */
export const POPULAR_LOCATIONS: Record<string, { address: string; lat: number; lng: number }> = {
  'Kigali Convention Centre': {
    address: 'Kigali Convention Centre, Kigali, Rwanda',
    lat: -1.9536,
    lng: 30.0935,
  },
  'Kigali International Airport': {
    address: 'Kigali International Airport, Kigali, Rwanda',
    lat: -1.9686,
    lng: 30.1395,
  },
  'Nyarutarama': {
    address: 'Nyarutarama, Kigali, Rwanda',
    lat: -1.9167,
    lng: 30.1167,
  },
  'Downtown Kigali': {
    address: 'Downtown Kigali, Kigali, Rwanda',
    lat: -1.9441,
    lng: 30.0619,
  },
  'Remera': {
    address: 'Remera, Kigali, Rwanda',
    lat: -1.9500,
    lng: 30.1000,
  },
  'Kimisagara': {
    address: 'Kimisagara, Kigali, Rwanda',
    lat: -1.9667,
    lng: 30.0833,
  },
  'Nyabugogo': {
    address: 'Nyabugogo, Kigali, Rwanda',
    lat: -1.9333,
    lng: 30.0500,
  },
};
