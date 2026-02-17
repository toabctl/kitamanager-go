'use client';

import { useMemo } from 'react';
import { useTranslations } from 'next-intl';
import { ResponsivePie } from '@nivo/pie';
import type { FinancialDataPoint } from '@/lib/api/types';
import { chartTheme } from './chart-utils';

interface FundingBreakdownChartProps {
  data: FinancialDataPoint;
}

function formatEur(cents: number): string {
  return (cents / 100).toLocaleString('de-DE', { style: 'currency', currency: 'EUR' });
}

const COLORS = ['#22c55e', '#14b8a6', '#06b6d4', '#8b5cf6', '#f59e0b', '#ec4899'];

export function FundingBreakdownChart({ data }: FundingBreakdownChartProps) {
  const t = useTranslations();

  const pieData = useMemo(() => {
    if (!data.funding_details?.length) return [];
    return data.funding_details.map((fd, i) => ({
      id: fd.value,
      label: fd.value,
      value: fd.amount_cents / 100,
      color: COLORS[i % COLORS.length],
    }));
  }, [data.funding_details]);

  if (pieData.length === 0) {
    return <p className="text-muted-foreground">{t('statistics.chartError')}</p>;
  }

  return (
    <div className="h-[350px]">
      <ResponsivePie
        data={pieData}
        margin={{ top: 30, right: 120, bottom: 30, left: 120 }}
        innerRadius={0.5}
        padAngle={1}
        cornerRadius={3}
        activeOuterRadiusOffset={6}
        colors={{ datum: 'data.color' }}
        arcLinkLabelsSkipAngle={10}
        arcLinkLabelsTextColor="hsl(var(--foreground))"
        arcLinkLabelsThickness={2}
        arcLinkLabelsColor={{ from: 'color' }}
        arcLabelsSkipAngle={10}
        arcLabelsTextColor="white"
        valueFormat={(v) => formatEur(v * 100)}
        tooltip={({ datum }) => (
          <div
            style={{
              background: 'hsl(var(--background))',
              color: 'hsl(var(--foreground))',
              border: '1px solid hsl(var(--border))',
              borderRadius: '6px',
              padding: '9px 12px',
              fontSize: 13,
            }}
          >
            <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
              <span
                style={{
                  width: 10,
                  height: 10,
                  borderRadius: '50%',
                  background: datum.color,
                  display: 'inline-block',
                }}
              />
              <strong>{datum.label}</strong>
            </div>
            <div style={{ marginTop: 4 }}>{formatEur(datum.value * 100)}</div>
          </div>
        )}
        theme={chartTheme}
      />
    </div>
  );
}
