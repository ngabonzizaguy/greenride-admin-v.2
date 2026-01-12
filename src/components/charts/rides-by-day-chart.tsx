'use client';

import { useEffect, useState } from 'react';
import {
  Bar,
  BarChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
  Cell,
} from 'recharts';

import { apiClient } from '@/lib/api-client';

type Point = { day: string; rides: number; isPeak: boolean };

const labelFromDate = (dateStr: string) => {
  const d = new Date(dateStr);
  if (!Number.isFinite(d.getTime())) return dateStr;
  return d.toLocaleDateString(undefined, { weekday: 'short' });
};

export function RidesByDayChart() {
  const [data, setData] = useState<Point[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    let mounted = true;
    const run = async () => {
      setIsLoading(true);
      try {
        const res = await apiClient.getRevenueChart({ period: '7d' });
        if (res.code === '0000' && Array.isArray(res.data) && mounted) {
          const pointsRaw = (res.data as Array<{ date: string; trips: number }>).map((p) => ({
            day: labelFromDate(p.date),
            rides: Number(p.trips) || 0,
          }));
          const max = pointsRaw.reduce((m, p) => Math.max(m, p.rides), 0);
          const points: Point[] = pointsRaw.map((p) => ({ ...p, isPeak: max > 0 && p.rides === max }));
          setData(points);
        } else if (mounted) {
          setData([]);
        }
      } catch {
        if (mounted) setData([]);
      } finally {
        if (mounted) setIsLoading(false);
      }
    };
    run();
    return () => {
      mounted = false;
    };
  }, []);

  if (isLoading) {
    return <div className="h-[300px] w-full flex items-center justify-center text-sm text-muted-foreground">Loading…</div>;
  }
  if (data.length === 0) {
    return <div className="h-[300px] w-full flex items-center justify-center text-sm text-muted-foreground">No rides-by-day data yet.</div>;
  }

  return (
    <div className="h-[300px] w-full">
      <ResponsiveContainer width="100%" height="100%">
        <BarChart data={data} margin={{ top: 10, right: 10, left: 0, bottom: 0 }}>
          <XAxis
            dataKey="day"
            axisLine={false}
            tickLine={false}
            tick={{ fill: '#64748B', fontSize: 12 }}
          />
          <YAxis
            axisLine={false}
            tickLine={false}
            tick={{ fill: '#64748B', fontSize: 12 }}
          />
          <Tooltip
            content={({ active, payload }) => {
              if (active && payload && payload.length) {
                const data = payload[0].payload;
                return (
                  <div className="rounded-lg border bg-background p-2 shadow-sm">
                    <p className="text-sm font-medium">{data.day}</p>
                    <p className="text-sm text-muted-foreground">
                      {data.rides} rides
                      {data.isPeak && (
                        <span className="ml-2 text-primary font-medium">Peak Day!</span>
                      )}
                    </p>
                  </div>
                );
              }
              return null;
            }}
          />
          <Bar dataKey="rides" radius={[4, 4, 0, 0]}>
            {data.map((entry, index) => (
              <Cell
                key={`cell-${index}`}
                fill={entry.isPeak ? '#22C55E' : '#86EFAC'}
              />
            ))}
          </Bar>
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
