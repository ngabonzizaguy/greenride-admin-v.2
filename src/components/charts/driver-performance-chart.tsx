'use client';

import {
  Bar,
  BarChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts';

// Generate mock data for last 30 days
const generateData = () => {
  const data = [];
  for (let i = 30; i >= 1; i--) {
    const date = new Date();
    date.setDate(date.getDate() - i);
    data.push({
      date: date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
      trips: Math.floor(Math.random() * 12) + 3,
    });
  }
  return data;
};

const data = generateData();

export function DriverPerformanceChart() {
  return (
    <div className="h-[300px] w-full">
      <ResponsiveContainer width="100%" height="100%">
        <BarChart data={data} margin={{ top: 10, right: 10, left: 0, bottom: 0 }}>
          <XAxis
            dataKey="date"
            axisLine={false}
            tickLine={false}
            tick={{ fill: '#64748B', fontSize: 10 }}
            interval={4}
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
                    <p className="text-sm text-muted-foreground">
                      {payload[0].value} trips
                    </p>
                  </div>
                );
              }
              return null;
            }}
          />
          <Bar
            dataKey="trips"
            fill="#22C55E"
            radius={[4, 4, 0, 0]}
          />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
