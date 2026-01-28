'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import { Search, X, Car, Users, Truck, Navigation, Loader2 } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import { apiClient } from '@/lib/api-client';
import { useRouter } from 'next/navigation';

interface SearchResult {
  id: string;
  type: 'driver' | 'vehicle' | 'user' | 'ride';
  title: string;
  subtitle: string;
  href: string;
}

interface GlobalSearchResults {
  drivers: SearchResult[];
  vehicles: SearchResult[];
  users: SearchResult[];
  rides: SearchResult[];
}

interface CommandPaletteProps {
  open: boolean;
  onClose: () => void;
}

export function CommandPalette({ open, onClose }: CommandPaletteProps) {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<GlobalSearchResults>({
    drivers: [],
    vehicles: [],
    users: [],
    rides: [],
  });
  const [isSearching, setIsSearching] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(0);
  const inputRef = useRef<HTMLInputElement>(null);
  const router = useRouter();

  // Get all results as flat array for navigation
  const allResults: SearchResult[] = [
    ...results.drivers,
    ...results.vehicles,
    ...results.users,
    ...results.rides,
  ];

  // Focus input when opened
  useEffect(() => {
    if (open && inputRef.current) {
      inputRef.current.focus();
      setQuery('');
      setSelectedIndex(0);
    }
  }, [open]);

  // Search with debounce
  useEffect(() => {
    if (!open || query.trim().length < 2) {
      setResults({ drivers: [], vehicles: [], users: [], rides: [] });
      return;
    }

    const timer = setTimeout(async () => {
      setIsSearching(true);
      try {
        const [driversRes, vehiclesRes, usersRes, ridesRes] = await Promise.all([
          apiClient.getDrivers({ keyword: query, limit: 5 }),
          apiClient.searchVehicles({ keyword: query, limit: 5 }),
          apiClient.getUsers({ keyword: query, limit: 5 }),
          apiClient.searchOrders({ keyword: query, limit: 5 }),
        ]);

        const formatDriver = (d: any): SearchResult => ({
          id: d.user_id,
          type: 'driver',
          title:
            d.display_name ||
            d.full_name ||
            d.username ||
            `${d.first_name || ''} ${d.last_name || ''}`.trim() ||
            d.phone ||
            'Unknown Driver',
          subtitle: d.phone || d.email || '',
          href: `/drivers/${d.user_id}`,
        });

        const formatVehicle = (v: any): SearchResult => ({
          id: v.vehicle_id,
          type: 'vehicle',
          title:
            `${v.brand || ''} ${v.model || ''}`.trim() ||
            v.plate_number ||
            v.vehicle_id ||
            'Unknown Vehicle',
          subtitle: v.plate_number || v.vehicle_id || '',
          href: `/vehicles`,
        });

        const formatUser = (u: any): SearchResult => ({
          id: u.user_id,
          type: 'user',
          title:
            u.display_name ||
            u.full_name ||
            u.username ||
            `${u.first_name || ''} ${u.last_name || ''}`.trim() ||
            u.phone ||
            'Unknown User',
          subtitle: u.phone || u.email || '',
          href: `/users/${u.user_id}`,
        });

        const formatRide = (r: any): SearchResult => ({
          id: r.order_id || r.id,
          type: 'ride',
          title: `Order #${r.order_id || r.id}`,
          subtitle: `${r.pickup_location || 'Unknown'} → ${r.dropoff_location || 'Unknown'}`,
          href: `/rides/${r.order_id || r.id}`,
        });

        setResults({
          drivers: (driversRes.data?.records || []).map(formatDriver),
          vehicles: (vehiclesRes.data?.records || []).map(formatVehicle),
          users: (usersRes.data?.records || []).map(formatUser),
          rides: (ridesRes.data?.records || []).map(formatRide),
        });
        setSelectedIndex(0);
      } catch (error) {
        console.error('Search failed:', error);
        setResults({ drivers: [], vehicles: [], users: [], rides: [] });
      } finally {
        setIsSearching(false);
      }
    }, 300);

    return () => clearTimeout(timer);
  }, [query, open]);

  // Keyboard navigation
  const handleKeyDown = useCallback((e: React.KeyboardEvent) => {
    if (e.key === 'Escape') {
      onClose();
      return;
    }

    if (e.key === 'ArrowDown') {
      e.preventDefault();
      setSelectedIndex((prev) => (prev < allResults.length - 1 ? prev + 1 : 0));
      return;
    }

    if (e.key === 'ArrowUp') {
      e.preventDefault();
      setSelectedIndex((prev) => (prev > 0 ? prev - 1 : allResults.length - 1));
      return;
    }

    if (e.key === 'Enter' && allResults[selectedIndex]) {
      e.preventDefault();
      router.push(allResults[selectedIndex].href);
      onClose();
      return;
    }
  }, [allResults, selectedIndex, router, onClose]);

  const handleSelect = (result: SearchResult) => {
    router.push(result.href);
    onClose();
  };

  const getTypeIcon = (type: SearchResult['type']) => {
    switch (type) {
      case 'driver':
        return <Car className="h-4 w-4" />;
      case 'vehicle':
        return <Truck className="h-4 w-4" />;
      case 'user':
        return <Users className="h-4 w-4" />;
      case 'ride':
        return <Navigation className="h-4 w-4" />;
    }
  };

  const getTypeLabel = (type: SearchResult['type']) => {
    switch (type) {
      case 'driver':
        return 'Driver';
      case 'vehicle':
        return 'Vehicle';
      case 'user':
        return 'User';
      case 'ride':
        return 'Ride';
    }
  };

  if (!open) return null;

  const hasResults = allResults.length > 0;
  const totalCount = allResults.length;

  return (
    <div className="fixed inset-0 z-50 flex items-start justify-center pt-[20vh]">
      {/* Backdrop */}
      <div
        className="fixed inset-0 bg-black/50 backdrop-blur-sm"
        onClick={onClose}
      />

      {/* Command Palette */}
      <div className="relative w-full max-w-2xl mx-4 bg-white dark:bg-gray-800 rounded-lg shadow-2xl border border-gray-200 dark:border-gray-700">
        {/* Search Input */}
        <div className="flex items-center gap-2 p-4 border-b border-gray-200 dark:border-gray-700">
          <Search className="h-5 w-5 text-gray-400" />
          <Input
            ref={inputRef}
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Search drivers, vehicles, users, rides..."
            className="flex-1 border-0 focus-visible:ring-0 focus-visible:ring-offset-0 text-lg"
          />
          {isSearching && <Loader2 className="h-5 w-5 animate-spin text-gray-400" />}
          <Button
            variant="ghost"
            size="icon"
            onClick={onClose}
            className="h-8 w-8"
          >
            <X className="h-4 w-4" />
          </Button>
        </div>

        {/* Results */}
        <div className="max-h-[60vh] overflow-y-auto">
          {query.trim().length < 2 ? (
            <div className="p-8 text-center text-gray-500">
              <Search className="h-12 w-12 mx-auto mb-4 opacity-50" />
              <p className="text-sm">Type at least 2 characters to search</p>
            </div>
          ) : isSearching ? (
            <div className="p-8 text-center text-gray-500">
              <Loader2 className="h-8 w-8 mx-auto mb-4 animate-spin" />
              <p className="text-sm">Searching...</p>
            </div>
          ) : !hasResults ? (
            <div className="p-8 text-center text-gray-500">
              <Search className="h-12 w-12 mx-auto mb-4 opacity-50" />
              <p className="text-sm">No results found for &quot;{query}&quot;</p>
            </div>
          ) : (
            <>
              {/* Results Summary */}
              <div className="px-4 py-2 text-xs text-gray-500 border-b border-gray-200 dark:border-gray-700">
                {totalCount} result{totalCount !== 1 ? 's' : ''} found
              </div>

              {/* Results by Type */}
              {results.drivers.length > 0 && (
                <div className="py-2">
                  <div className="px-4 py-1 text-xs font-semibold text-gray-500 uppercase">
                    Drivers ({results.drivers.length})
                  </div>
                  {results.drivers.map((result, idx) => {
                    const flatIndex = results.drivers.slice(0, idx).length;
                    return (
                      <button
                        key={result.id}
                        onClick={() => handleSelect(result)}
                        className={cn(
                          'w-full px-4 py-3 text-left hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors flex items-center gap-3',
                          flatIndex === selectedIndex && 'bg-gray-100 dark:bg-gray-700'
                        )}
                      >
                        <div className="flex-shrink-0 text-gray-400">
                          {getTypeIcon(result.type)}
                        </div>
                        <div className="flex-1 min-w-0">
                          <div className="font-medium truncate">{result.title}</div>
                          <div className="text-sm text-gray-500 truncate">{result.subtitle}</div>
                        </div>
                        <Badge variant="outline" className="text-xs">
                          {getTypeLabel(result.type)}
                        </Badge>
                      </button>
                    );
                  })}
                </div>
              )}

              {results.vehicles.length > 0 && (
                <div className="py-2">
                  <div className="px-4 py-1 text-xs font-semibold text-gray-500 uppercase">
                    Vehicles ({results.vehicles.length})
                  </div>
                  {results.vehicles.map((result, idx) => {
                    const flatIndex = results.drivers.length + idx;
                    return (
                      <button
                        key={result.id}
                        onClick={() => handleSelect(result)}
                        className={cn(
                          'w-full px-4 py-3 text-left hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors flex items-center gap-3',
                          flatIndex === selectedIndex && 'bg-gray-100 dark:bg-gray-700'
                        )}
                      >
                        <div className="flex-shrink-0 text-gray-400">
                          {getTypeIcon(result.type)}
                        </div>
                        <div className="flex-1 min-w-0">
                          <div className="font-medium truncate">{result.title}</div>
                          <div className="text-sm text-gray-500 truncate">{result.subtitle}</div>
                        </div>
                        <Badge variant="outline" className="text-xs">
                          {getTypeLabel(result.type)}
                        </Badge>
                      </button>
                    );
                  })}
                </div>
              )}

              {results.users.length > 0 && (
                <div className="py-2">
                  <div className="px-4 py-1 text-xs font-semibold text-gray-500 uppercase">
                    Users ({results.users.length})
                  </div>
                  {results.users.map((result, idx) => {
                    const flatIndex = results.drivers.length + results.vehicles.length + idx;
                    return (
                      <button
                        key={result.id}
                        onClick={() => handleSelect(result)}
                        className={cn(
                          'w-full px-4 py-3 text-left hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors flex items-center gap-3',
                          flatIndex === selectedIndex && 'bg-gray-100 dark:bg-gray-700'
                        )}
                      >
                        <div className="flex-shrink-0 text-gray-400">
                          {getTypeIcon(result.type)}
                        </div>
                        <div className="flex-1 min-w-0">
                          <div className="font-medium truncate">{result.title}</div>
                          <div className="text-sm text-gray-500 truncate">{result.subtitle}</div>
                        </div>
                        <Badge variant="outline" className="text-xs">
                          {getTypeLabel(result.type)}
                        </Badge>
                      </button>
                    );
                  })}
                </div>
              )}

              {results.rides.length > 0 && (
                <div className="py-2">
                  <div className="px-4 py-1 text-xs font-semibold text-gray-500 uppercase">
                    Rides ({results.rides.length})
                  </div>
                  {results.rides.map((result, idx) => {
                    const flatIndex = results.drivers.length + results.vehicles.length + results.users.length + idx;
                    return (
                      <button
                        key={result.id}
                        onClick={() => handleSelect(result)}
                        className={cn(
                          'w-full px-4 py-3 text-left hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors flex items-center gap-3',
                          flatIndex === selectedIndex && 'bg-gray-100 dark:bg-gray-700'
                        )}
                      >
                        <div className="flex-shrink-0 text-gray-400">
                          {getTypeIcon(result.type)}
                        </div>
                        <div className="flex-1 min-w-0">
                          <div className="font-medium truncate">{result.title}</div>
                          <div className="text-sm text-gray-500 truncate">{result.subtitle}</div>
                        </div>
                        <Badge variant="outline" className="text-xs">
                          {getTypeLabel(result.type)}
                        </Badge>
                      </button>
                    );
                  })}
                </div>
              )}
            </>
          )}
        </div>

        {/* Footer */}
        <div className="px-4 py-2 border-t border-gray-200 dark:border-gray-700 flex items-center justify-between text-xs text-gray-500">
          <div className="flex items-center gap-4">
            <span className="flex items-center gap-1">
              <kbd className="px-1.5 py-0.5 bg-gray-100 dark:bg-gray-700 rounded">↑↓</kbd>
              Navigate
            </span>
            <span className="flex items-center gap-1">
              <kbd className="px-1.5 py-0.5 bg-gray-100 dark:bg-gray-700 rounded">Enter</kbd>
              Select
            </span>
            <span className="flex items-center gap-1">
              <kbd className="px-1.5 py-0.5 bg-gray-100 dark:bg-gray-700 rounded">Esc</kbd>
              Close
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
