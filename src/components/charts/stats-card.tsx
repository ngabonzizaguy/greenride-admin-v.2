'use client';

import { LucideIcon, ArrowUpRight, ArrowDownRight } from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';
import { cn } from '@/lib/utils';

interface StatsCardProps {
  title: string;
  value: string | number;
  icon: LucideIcon;
  change?: number;
  changeType?: 'increase' | 'decrease';
  changeLabel?: string;
  subtext?: string;
  variant?: 'default' | 'warning' | 'success' | 'danger';
}

export function StatsCard({
  title,
  value,
  icon: Icon,
  change,
  changeType,
  changeLabel,
  subtext,
  variant = 'default',
}: StatsCardProps) {
  const variantStyles = {
    default: 'bg-primary/10 text-primary',
    warning: 'bg-yellow-100 text-yellow-600',
    success: 'bg-green-100 text-green-600',
    danger: 'bg-red-100 text-red-600',
  };

  return (
    <Card>
      <CardContent className="p-6">
        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <p className="text-sm font-medium text-muted-foreground">{title}</p>
            <p className="text-2xl font-bold">{value}</p>
            {change !== undefined && changeType && (
              <div className="flex items-center gap-1 text-sm">
                {changeType === 'increase' ? (
                  <ArrowUpRight className="h-4 w-4 text-green-500" />
                ) : (
                  <ArrowDownRight className="h-4 w-4 text-red-500" />
                )}
                <span
                  className={cn(
                    'font-medium',
                    changeType === 'increase' ? 'text-green-500' : 'text-red-500'
                  )}
                >
                  {changeType === 'increase' ? '+' : '-'}{change}%
                </span>
                {changeLabel && (
                  <span className="text-muted-foreground">{changeLabel}</span>
                )}
              </div>
            )}
            {subtext && (
              <p className="text-sm text-muted-foreground">{subtext}</p>
            )}
          </div>
          <div
            className={cn(
              'flex h-12 w-12 items-center justify-center rounded-full',
              variantStyles[variant]
            )}
          >
            <Icon className="h-6 w-6" />
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
