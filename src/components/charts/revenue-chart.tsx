'use client';

import { useEffect, useMemo, useState } from 'react';
import {
  Area,
  AreaChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts';

import { apiClient } from '@/lib/api-client';

type RevenuePoint = { day: string; revenue: number };

const labelFromDate = (dateStr: string) => {
  const d = new Date(dateStr);
  if (!Number.isFinite(d.getTime())) return dateStr;
  return d.toLocaleDateString(undefined, { weekday: 'short' });
};

export function RevenueChart() {
  const [data, setData] = useState<RevenuePoint[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    let mounted = true;
    const run = async () => {
      setIsLoading(true);
      try {
        const res = await apiClient.getRevenueChart({ period: '7d' });
        if (res.code === '0000' && Array.isArray(res.data) && mounted) {
          const points = (res.data as Array<{ date: string; revenue: number }>).map((p) => ({
            day: labelFromDate(p.date),
            revenue: Number(p.revenue) || 0,
          }));
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

  const yTick = useMemo(() => (value: number) => `${(value / 1000).toFixed(0)}K`, []);

  if (isLoading) {
    return <div className="h-[300px] w-full flex items-center justify-center text-sm text-muted-foreground">Loading…</div>;
  }

  if (data.length === 0) {
    return <div className="h-[300px] w-full flex items-center justify-center text-sm text-muted-foreground">No revenue data yet.</div>;
  }

  return (
    <div className="h-[300px] w-full">
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart data={data} margin={{ top: 10, right: 10, left: 0, bottom: 0 }}>
          <defs>
            <linearGradient id="colorRevenue" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#22C55E" stopOpacity={0.3} />
              <stop offset="95%" stopColor="#22C55E" stopOpacity={0} />
            </linearGradient>
          </defs>
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
            tickFormatter={yTick}
          />
          <Tooltip
            content={({ active, payload }) => {
              if (active && payload && payload.length) {
                return (
                  <div className="rounded-lg border bg-background p-2 shadow-sm">
                    <p className="text-sm font-medium">
                      RWF {payload[0].value?.toLocaleString()}
                    </p>
                  </div>
                );
              }
              return null;
            }}
          />
          <Area
            type="monotone"
            dataKey="revenue"
            stroke="#22C55E"
            strokeWidth={2}
            fill="url(#colorRevenue)"
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}
