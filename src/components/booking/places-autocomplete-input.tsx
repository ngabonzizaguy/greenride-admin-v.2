/* eslint-disable @next/next/no-img-element */
'use client';

import { useEffect, useMemo, useRef, useState } from 'react';
import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';

type PlaceSuggestion = {
  placeId: string;
  description: string;
};

export type PlaceSelection = {
  placeId: string;
  address: string;
  lat: number;
  lng: number;
};

export function PlacesAutocompleteInput(props: {
  id: string;
  placeholder?: string;
  value: string;
  onChange: (value: string) => void;
  onSelect: (place: PlaceSelection) => void;
  disabled?: boolean;
  className?: string;
}) {
  const { id, placeholder, value, onChange, onSelect, disabled, className } = props;

  const containerRef = useRef<HTMLDivElement | null>(null);
  const [isOpen, setIsOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [suggestions, setSuggestions] = useState<PlaceSuggestion[]>([]);

  const canUsePlaces = useMemo(
    () => typeof window !== 'undefined' && Boolean((window as any).google?.maps?.places),
    []
  );

  useEffect(() => {
    if (!isOpen) return;
    const onDocMouseDown = (e: MouseEvent) => {
      if (!containerRef.current) return;
      if (e.target instanceof Node && !containerRef.current.contains(e.target)) {
        setIsOpen(false);
      }
    };
    document.addEventListener('mousedown', onDocMouseDown);
    return () => document.removeEventListener('mousedown', onDocMouseDown);
  }, [isOpen]);

  useEffect(() => {
    if (disabled) return;
    if (!canUsePlaces) return;

    const q = value.trim();
    if (!q) {
      setSuggestions([]);
      setError(null);
      return;
    }

    setIsLoading(true);
    setError(null);
    const handle = window.setTimeout(() => {
      try {
        const service = new google.maps.places.AutocompleteService();
        service.getPlacePredictions(
          {
            input: q,
            // Bias toward Rwanda
            componentRestrictions: { country: 'rw' },
          },
          (predictions, status) => {
            setIsLoading(false);
            if (status !== google.maps.places.PlacesServiceStatus.OK || !predictions) {
              setSuggestions([]);
              if (status === google.maps.places.PlacesServiceStatus.ZERO_RESULTS) {
                setError(null);
              } else {
                setError('Unable to load suggestions');
              }
              return;
            }
            setSuggestions(
              predictions.slice(0, 7).map((p) => ({
                placeId: p.place_id,
                description: p.description,
              }))
            );
          }
        );
      } catch (e) {
        setIsLoading(false);
        setSuggestions([]);
        setError(e instanceof Error ? e.message : 'Unable to load suggestions');
      }
    }, 220);

    return () => window.clearTimeout(handle);
  }, [value, disabled, canUsePlaces]);

  const resolvePlace = async (placeId: string, description: string) => {
    if (!canUsePlaces) return;
    setIsLoading(true);
    setError(null);

    try {
      // PlacesService requires an HTMLDivElement.
      const el = document.createElement('div');
      const service = new google.maps.places.PlacesService(el);

      service.getDetails(
        {
          placeId,
          fields: ['formatted_address', 'geometry', 'name'],
        },
        (place, status) => {
          setIsLoading(false);
          if (status !== google.maps.places.PlacesServiceStatus.OK || !place?.geometry?.location) {
            setError('Unable to resolve selected place');
            return;
          }

          const lat = place.geometry.location.lat();
          const lng = place.geometry.location.lng();
          const address = place.formatted_address || place.name || description;

          onSelect({ placeId, address, lat, lng });
          setIsOpen(false);
          setSuggestions([]);
        }
      );
    } catch (e) {
      setIsLoading(false);
      setError(e instanceof Error ? e.message : 'Unable to resolve selected place');
    }
  };

  return (
    <div ref={containerRef} className={cn('relative', className)}>
      <Input
        id={id}
        placeholder={placeholder}
        value={value}
        onFocus={() => setIsOpen(true)}
        onChange={(e) => {
          setIsOpen(true);
          onChange(e.target.value);
        }}
        disabled={disabled}
        autoComplete="off"
      />

      {isOpen && !disabled && (
        <div className="absolute z-20 mt-1 w-full rounded-md border bg-background shadow-lg">
          <div className="p-2 text-xs text-muted-foreground flex items-center justify-between">
            <span>Suggestions</span>
            {isLoading && <span>Loading…</span>}
          </div>

          {error && <div className="px-2 pb-2 text-xs text-red-600">{error}</div>}

          {!canUsePlaces && (
            <div className="px-2 pb-2 text-xs text-amber-700">
              Google Places isn’t available (check API key + “Places API” enabled). Falling back to
              basic geocoding.
            </div>
          )}

          {canUsePlaces && suggestions.length === 0 && value.trim() && !isLoading && !error && (
            <div className="px-2 pb-2 text-xs text-muted-foreground">No suggestions</div>
          )}

          {suggestions.map((s) => (
            <button
              key={s.placeId}
              type="button"
              className="w-full text-left px-2 py-2 text-sm hover:bg-accent"
              onClick={() => resolvePlace(s.placeId, s.description)}
            >
              {s.description}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}

