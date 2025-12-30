'use client';

import {
  Area,
  AreaChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
  Legend,
} from 'recharts';

const data = [
  { date: 'Week 1', newUsers: 45, returningUsers: 120 },
  { date: 'Week 2', newUsers: 52, returningUsers: 135 },
  { date: 'Week 3', newUsers: 38, returningUsers: 142 },
  { date: 'Week 4', newUsers: 67, returningUsers: 158 },
  { date: 'Week 5', newUsers: 54, returningUsers: 175 },
  { date: 'Week 6', newUsers: 72, returningUsers: 189 },
  { date: 'Week 7', newUsers: 48, returningUsers: 195 },
  { date: 'Week 8', newUsers: 63, returningUsers: 210 },
];

export function UserGrowthChart() {
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
                      Returning: {payload[1].value}
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
            formatter={(value) => (
              <span className="text-sm text-muted-foreground">
                {value === 'newUsers' ? 'New Users' : 'Returning Users'}
              </span>
            )}
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
            dataKey="returningUsers"
            stroke="#3B82F6"
            strokeWidth={2}
            fill="url(#colorReturning)"
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}
