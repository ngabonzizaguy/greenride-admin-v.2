'use client';

import {
  Area,
  AreaChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts';

const data = [
  { day: 'Mon', revenue: 185000 },
  { day: 'Tue', revenue: 210000 },
  { day: 'Wed', revenue: 195000 },
  { day: 'Thu', revenue: 230000 },
  { day: 'Fri', revenue: 280000 },
  { day: 'Sat', revenue: 320000 },
  { day: 'Sun', revenue: 245000 },
];

export function RevenueChart() {
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
            tickFormatter={(value) => `${(value / 1000).toFixed(0)}K`}
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
