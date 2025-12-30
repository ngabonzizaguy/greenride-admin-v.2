'use client';

import { ResponsiveContainer, Tooltip } from 'recharts';

// Generate mock heatmap data
const days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
const hours = ['6 AM', '8 AM', '10 AM', '12 PM', '2 PM', '4 PM', '6 PM', '8 PM', '10 PM'];

const generateHeatmapData = () => {
  return days.map((day) => ({
    day,
    ...hours.reduce((acc, hour, idx) => {
      // Simulate higher demand during rush hours
      let value = Math.floor(Math.random() * 30) + 10;
      if (idx === 1 || idx === 6) value += 30; // 8 AM and 6 PM peak
      if (day === 'Fri' || day === 'Sat') value += 15;
      return { ...acc, [hour]: value };
    }, {}),
  }));
};

const data = generateHeatmapData();

const getColor = (value: number) => {
  if (value >= 60) return 'bg-green-600';
  if (value >= 45) return 'bg-green-500';
  if (value >= 30) return 'bg-green-400';
  if (value >= 20) return 'bg-green-300';
  return 'bg-green-200';
};

export function PeakHoursChart() {
  return (
    <div className="space-y-4">
      {/* Hours Header */}
      <div className="flex">
        <div className="w-12" />
        {hours.map((hour) => (
          <div key={hour} className="flex-1 text-center text-xs text-muted-foreground">
            {hour}
          </div>
        ))}
      </div>

      {/* Heatmap Grid */}
      <div className="space-y-1">
        {data.map((row) => (
          <div key={row.day} className="flex items-center">
            <div className="w-12 text-xs font-medium text-muted-foreground">{row.day}</div>
            <div className="flex flex-1 gap-1">
              {hours.map((hour) => {
                const value = (row as Record<string, number | string>)[hour] as number;
                return (
                  <div
                    key={hour}
                    className={`flex-1 h-8 rounded ${getColor(value)} cursor-pointer transition-opacity hover:opacity-80`}
                    title={`${row.day} ${hour}: ${value} rides`}
                  />
                );
              })}
            </div>
          </div>
        ))}
      </div>

      {/* Legend */}
      <div className="flex items-center justify-center gap-4 pt-2">
        <div className="flex items-center gap-1.5">
          <div className="h-3 w-3 rounded bg-green-200" />
          <span className="text-xs text-muted-foreground">Low</span>
        </div>
        <div className="flex items-center gap-1.5">
          <div className="h-3 w-3 rounded bg-green-400" />
          <span className="text-xs text-muted-foreground">Medium</span>
        </div>
        <div className="flex items-center gap-1.5">
          <div className="h-3 w-3 rounded bg-green-600" />
          <span className="text-xs text-muted-foreground">High</span>
        </div>
      </div>
    </div>
  );
}
