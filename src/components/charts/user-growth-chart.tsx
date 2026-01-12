'use client';

import { useEffect, useMemo, useState } from 'react';
import {
  Area,
  AreaChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
  Legend,
} from 'recharts';

import { apiClient } from '@/lib/api-client';

type Point = { date: string; newUsers: number; totalUsers: number };

const labelFromDate = (dateStr: string) => {
  const d = new Date(dateStr);
  if (!Number.isFinite(d.getTime())) return dateStr;
  return d.toLocaleDateString(undefined, { month: 'short', day: '2-digit' });
};

export function UserGrowthChart() {
  const [data, setData] = useState<Point[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    let mounted = true;
    const run = async () => {
      setIsLoading(true);
      try {
        const res = await apiClient.getUserGrowthChart({ period: '30d' });
        if (res.code === '0000' && Array.isArray(res.data) && mounted) {
          const points = (res.data as Array<{ date: string; new_users: number; total_users: number }>).map((p) => ({
            date: labelFromDate(p.date),
            newUsers: Number(p.new_users) || 0,
            totalUsers: Number(p.total_users) || 0,
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

  const legendFormatter = useMemo(
    () => (value: string) => (
      <span className="text-sm text-muted-foreground">
        {value === 'newUsers' ? 'New Users' : value === 'totalUsers' ? 'Total Users' : value}
      </span>
    ),
    []
  );

  if (isLoading) {
    return <div className="h-[300px] w-full flex items-center justify-center text-sm text-muted-foreground">Loading…</div>;
  }
  if (data.length === 0) {
    return <div className="h-[300px] w-full flex items-center justify-center text-sm text-muted-foreground">No user growth data yet.</div>;
  }

  return (
    <div className="h-[300px] w-full">
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart data={data} margin={{ top: 10, right: 10, left: 0, bottom: 0 }}>
          <defs>
            <linearGradient id="colorNew" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#22C55E" stopOpacity={0.3} />
              <stop offset="95%" stopColor="#22C55E" stopOpacity={0} />
            </linearGradient>
            <linearGradient id="colorReturning" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#3B82F6" stopOpacity={0.3} />
              <stop offset="95%" stopColor="#3B82F6" stopOpacity={0} />
            </linearGradient>
          </defs>
          <XAxis
            dataKey="date"
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
                return (
                  <div className="rounded-lg border bg-background p-2 shadow-sm">
                    <p className="text-sm font-medium">{payload[0].payload.date}</p>
                    <p className="text-sm text-green-600">
                      New: {payload[0].value}
                    </p>
                    <p className="text-sm text-blue-600">
                      Total: {payload[1].value}
                    </p>
                  </div>
                );
              }
              return null;
            }}
          />
          <Legend
            verticalAlign="bottom"
            height={36}
            formatter={legendFormatter}
          />
          <Area
            type="monotone"
            dataKey="newUsers"
            stroke="#22C55E"
            strokeWidth={2}
            fill="url(#colorNew)"
          />
          <Area
            type="monotone"
            dataKey="totalUsers"
            stroke="#3B82F6"
            strokeWidth={2}
            fill="url(#colorReturning)"
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}
