'use client';

import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip, Legend } from 'recharts';

const data = [
  { name: 'Cash', value: 45, color: '#64748B' },
  { name: 'MoMo', value: 40, color: '#F59E0B' },
  { name: 'Card', value: 15, color: '#3B82F6' },
];

export function PaymentMethodsChart() {
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
            label={({ name, percent }) => `${name} ${((percent ?? 0) * 100).toFixed(0)}%`}
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
