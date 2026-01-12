'use client';

import { useEffect, useMemo, useState } from 'react';
import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip, Legend } from 'recharts';

import { apiClient } from '@/lib/api-client';

type Slice = { name: string; value: number; color: string };

const COLORS: Record<string, string> = {
  cash: '#64748B',
  momo: '#F59E0B',
  card: '#3B82F6',
  other: '#A855F7',
};

export function PaymentMethodsChart() {
  const [data, setData] = useState<Slice[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    let mounted = true;
    const run = async () => {
      setIsLoading(true);
      try {
        const res = await apiClient.getDashboardStats();
        const raw = res.data as Record<string, unknown>;
        const trips = (raw?.recent_trips as Array<Record<string, unknown>>) ?? [];

        const counts: Record<string, number> = {};
        for (const t of trips) {
          const pm = String(t.payment_method ?? 'other').toLowerCase();
          counts[pm] = (counts[pm] || 0) + 1;
        }

        const total = Object.values(counts).reduce((a, b) => a + b, 0);
        const slices: Slice[] =
          total === 0
            ? []
            : Object.entries(counts).map(([k, v]) => ({
                name: k === 'momo' ? 'MoMo' : k.charAt(0).toUpperCase() + k.slice(1),
                value: Math.round((v / total) * 100),
                color: COLORS[k] || COLORS.other,
              }));

        if (mounted) setData(slices);
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

  const hasData = data.length > 0 && data.some((d) => d.value > 0);
  const label = useMemo(
    () => ({ name, percent }: { name?: string; percent?: number }) => `${name ?? ''} ${(((percent ?? 0) * 100) || 0).toFixed(0)}%`,
    []
  );

  if (isLoading) {
    return <div className="h-[300px] w-full flex items-center justify-center text-sm text-muted-foreground">Loading…</div>;
  }

  if (!hasData) {
    return <div className="h-[300px] w-full flex items-center justify-center text-sm text-muted-foreground">No payment method data yet.</div>;
  }

  return (
    <div className="h-[300px] w-full">
      <ResponsiveContainer width="100%" height="100%">
        <PieChart>
          <Pie
            data={data}
            cx="50%"
            cy="50%"
            innerRadius={60}
            outerRadius={100}
            paddingAngle={2}
            dataKey="value"
            label={label}
            labelLine={false}
          >
            {data.map((entry, index) => (
              <Cell key={`cell-${index}`} fill={entry.color} />
            ))}
          </Pie>
          <Tooltip
            content={({ active, payload }) => {
              if (active && payload && payload.length) {
                return (
                  <div className="rounded-lg border bg-background p-2 shadow-sm">
                    <p className="text-sm font-medium">
                      {payload[0].name}: {payload[0].value}%
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
            formatter={(value, entry) => (
              <span className="text-sm text-muted-foreground">{value}</span>
            )}
          />
        </PieChart>
      </ResponsiveContainer>
    </div>
  );
}
