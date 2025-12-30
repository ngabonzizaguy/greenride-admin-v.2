'use client';

import {
  Bar,
  BarChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
  Cell,
} from 'recharts';

const data = [
  { day: 'Mon', rides: 145, isPeak: false },
  { day: 'Tue', rides: 132, isPeak: false },
  { day: 'Wed', rides: 156, isPeak: false },
  { day: 'Thu', rides: 178, isPeak: false },
  { day: 'Fri', rides: 234, isPeak: true },
  { day: 'Sat', rides: 198, isPeak: false },
  { day: 'Sun', rides: 156, isPeak: false },
];

export function RidesByDayChart() {
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
