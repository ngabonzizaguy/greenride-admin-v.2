'use client';

import {
  Bar,
  BarChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
  ReferenceLine,
  Cell,
} from 'recharts';

const data = [
  { range: '0-2 km', count: 145, color: '#22C55E' },
  { range: '2-5 km', count: 312, color: '#22C55E' },
  { range: '5-10 km', count: 234, color: '#22C55E' },
  { range: '10-20 km', count: 98, color: '#22C55E' },
  { range: '20+ km', count: 45, color: '#22C55E' },
];

const avgDistance = 6.8; // Average distance in km

export function DistanceDistributionChart() {
  return (
    <div className="h-[250px] w-full">
      <ResponsiveContainer width="100%" height="100%">
        <BarChart data={data} margin={{ top: 20, right: 20, left: 20, bottom: 5 }}>
          <XAxis
            dataKey="range"
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
                    <p className="text-sm font-medium">{payload[0].payload.range}</p>
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
            dataKey="count" 
            radius={[4, 4, 0, 0]}
            label={{ 
              position: 'top', 
              fill: '#64748B', 
              fontSize: 12 
            }}
          >
            {data.map((entry, index) => (
              <Cell key={`cell-${index}`} fill={entry.color} />
            ))}
          </Bar>
        </BarChart>
      </ResponsiveContainer>
      <div className="text-center text-sm text-muted-foreground">
        Average trip distance: <span className="font-medium text-foreground">{avgDistance} km</span>
      </div>
    </div>
  );
}
