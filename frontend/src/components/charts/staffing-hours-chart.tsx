'use client';

import { useTranslations } from 'next-intl';
import { ResponsiveLine } from '@nivo/line';
import type { StaffingHoursResponse } from '@/lib/api/types';

interface StaffingHoursChartProps {
  data: StaffingHoursResponse;
}

export function StaffingHoursChart({ data }: StaffingHoursChartProps) {
  const t = useTranslations();

  // Format dates as "Jan 25", "Feb 25", etc.
  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr + 'T00:00:00');
    return date.toLocaleDateString('en-US', { month: 'short', year: '2-digit' });
  };

  // Find today marker position
  const today = new Date();
  const todayStr = today.toISOString().slice(0, 10);

  const chartData = [
    {
      id: t('statistics.requiredHours'),
      color: '#f59e0b',
      data: data.data_points.map((dp) => ({
        x: formatDate(dp.date),
        y: Math.round(dp.required_hours * 100) / 100,
      })),
    },
    {
      id: t('statistics.availableHours'),
      color: '#3b82f6',
      data: data.data_points.map((dp) => ({
        x: formatDate(dp.date),
        y: dp.available_hours,
      })),
    },
  ];

  // Find the x label for today's month
  const todayLabel = formatDate(todayStr);

  return (
    <div className="h-[300px]">
      <ResponsiveLine
        data={chartData}
        margin={{ top: 20, right: 110, bottom: 50, left: 60 }}
        xScale={{ type: 'point' }}
        yScale={{ type: 'linear', min: 'auto', max: 'auto', stacked: false }}
        curve="monotoneX"
        axisTop={null}
        axisRight={null}
        axisBottom={{
          tickSize: 5,
          tickPadding: 5,
          tickRotation: -45,
        }}
        axisLeft={{
          tickSize: 5,
          tickPadding: 5,
          tickRotation: 0,
        }}
        colors={['#f59e0b', '#3b82f6']}
        pointSize={8}
        pointColor={{ theme: 'background' }}
        pointBorderWidth={2}
        pointBorderColor={{ from: 'serieColor' }}
        pointLabelYOffset={-12}
        useMesh={true}
        enableSlices="x"
        markers={[
          {
            axis: 'x',
            value: todayLabel,
            lineStyle: {
              stroke: 'hsl(var(--foreground))',
              strokeWidth: 1,
              strokeDasharray: '4 4',
            },
            legend: t('common.today'),
            legendPosition: 'top',
            textStyle: {
              fill: 'hsl(var(--muted-foreground))',
              fontSize: 11,
            },
          },
        ]}
        legends={[
          {
            anchor: 'bottom-right',
            direction: 'column',
            justify: false,
            translateX: 100,
            translateY: 0,
            itemsSpacing: 0,
            itemDirection: 'left-to-right',
            itemWidth: 80,
            itemHeight: 20,
            itemOpacity: 0.85,
            symbolSize: 12,
            symbolShape: 'circle',
          },
        ]}
        theme={{
          axis: {
            ticks: {
              text: {
                fill: 'hsl(var(--muted-foreground))',
              },
            },
          },
          grid: {
            line: {
              stroke: 'hsl(var(--border))',
            },
          },
          legends: {
            text: {
              fill: 'hsl(var(--foreground))',
            },
          },
          crosshair: {
            line: {
              stroke: 'hsl(var(--foreground))',
              strokeWidth: 1,
              strokeOpacity: 0.35,
            },
          },
          tooltip: {
            container: {
              background: 'hsl(var(--background))',
              color: 'hsl(var(--foreground))',
              border: '1px solid hsl(var(--border))',
              borderRadius: '6px',
            },
          },
        }}
      />
    </div>
  );
}
