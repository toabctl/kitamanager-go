'use client';

import { useQuery } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { apiClient } from '@/lib/api/client';
import { formatCurrency } from '@/lib/utils/formatting';

function formatServiceStart(dateStr: string): string {
  const d = new Date(dateStr);
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`;
}

interface StepPromotionsWidgetProps {
  orgId: number;
}

export function StepPromotionsWidget({ orgId }: StepPromotionsWidgetProps) {
  const t = useTranslations('stepPromotions');

  const { data } = useQuery({
    queryKey: ['step-promotions', orgId],
    queryFn: () => apiClient.getStepPromotions(orgId),
    enabled: !!orgId,
  });

  if (!data || data.promotions.length === 0) {
    return null;
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-base font-medium">{t('title')}</CardTitle>
        <Badge variant="secondary">
          {t('totalCost', { amount: formatCurrency(data.total_monthly_cost_delta) })}
        </Badge>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>{t('employee')}</TableHead>
              <TableHead>{t('grade')}</TableHead>
              <TableHead className="text-center">{t('currentStep')}</TableHead>
              <TableHead className="text-center">{t('eligibleStep')}</TableHead>
              <TableHead>{t('serviceStart')}</TableHead>
              <TableHead className="text-right">{t('yearsOfService')}</TableHead>
              <TableHead className="text-right">{t('monthlyCostDelta')}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {data.promotions.map((p) => (
              <TableRow key={p.employee_id}>
                <TableCell className="font-medium">{p.employee_name}</TableCell>
                <TableCell>{p.grade}</TableCell>
                <TableCell className="text-center">{p.current_step}</TableCell>
                <TableCell className="text-center">{p.eligible_step}</TableCell>
                <TableCell>{formatServiceStart(p.service_start)}</TableCell>
                <TableCell className="text-right">{p.years_of_service.toFixed(1)}</TableCell>
                <TableCell className="text-right">
                  +{formatCurrency(p.monthly_cost_delta)}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
