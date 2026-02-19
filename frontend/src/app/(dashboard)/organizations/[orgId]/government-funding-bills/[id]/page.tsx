'use client';

import Link from 'next/link';
import { useParams } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery } from '@tanstack/react-query';
import { ArrowLeft, CheckCircle2, XCircle } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
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
import { queryKeys } from '@/lib/api/queryKeys';
import { formatCurrency } from '@/lib/utils/formatting';

export default function GovernmentFundingBillDetailPage() {
  const params = useParams();
  const orgId = Number(params.orgId);
  const id = Number(params.id);
  const t = useTranslations('governmentFundingBills');
  const tCommon = useTranslations('common');

  const { data: result, isLoading } = useQuery({
    queryKey: queryKeys.governmentFundingBillPeriods.detail(orgId, id),
    queryFn: () => apiClient.getGovernmentFundingBillPeriod(orgId, id),
  });

  if (isLoading) {
    return <p className="text-muted-foreground py-8 text-center">{tCommon('loading')}</p>;
  }

  if (!result) {
    return <p className="text-muted-foreground py-8 text-center">{tCommon('notFound')}</p>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" asChild>
          <Link href={`/organizations/${orgId}/government-funding-bills`}>
            <ArrowLeft className="h-4 w-4" />
          </Link>
        </Button>
        <div>
          <h1 className="text-2xl font-bold">{result.facility_name}</h1>
          <p className="text-muted-foreground text-sm">
            {new Date(result.from).toLocaleDateString('de-DE', {
              month: 'long',
              year: 'numeric',
            })}
            {' \u2014 '}
            {result.file_name}
          </p>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-muted-foreground text-sm font-medium">
              {t('facilityTotal')}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold">{formatCurrency(result.facility_total)}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-muted-foreground text-sm font-medium">
              {t('contractBooking')} / {t('correctionBooking')}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold">
              {formatCurrency(result.contract_booking)} /{' '}
              {formatCurrency(result.correction_booking)}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-muted-foreground text-sm font-medium">
              {t('matchedChildren')} / {t('unmatchedChildren')}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold">
              <span className="text-green-600">{result.matched_count}</span>
              {' / '}
              <span className="text-red-600">{result.unmatched_count}</span>
              <span className="text-muted-foreground ml-2 text-sm">
                ({result.children_count} {t('children')})
              </span>
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Surcharges */}
      {result.surcharges && result.surcharges.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>{t('surcharges')}</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 gap-2 sm:grid-cols-2 md:grid-cols-3">
              {result.surcharges.map((s) => (
                <div
                  key={`${s.key}-${s.value}`}
                  className="flex justify-between rounded-md border p-3"
                >
                  <span className="text-muted-foreground text-sm">
                    {s.key}: {s.value}
                  </span>
                  <span className="font-medium">{formatCurrency(s.amount)}</span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Children Table */}
      <Card>
        <CardHeader>
          <CardTitle>
            {t('children')} ({result.children_count})
          </CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t('voucherNumber')}</TableHead>
                <TableHead>{t('childName')}</TableHead>
                <TableHead className="hidden md:table-cell">{t('birthDate')}</TableHead>
                <TableHead className="hidden md:table-cell">{t('district')}</TableHead>
                <TableHead className="text-right">{t('totalAmount')}</TableHead>
                <TableHead>{t('matched')}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {result.children.map((child) => (
                <TableRow key={child.voucher_number}>
                  <TableCell className="font-mono text-sm">{child.voucher_number}</TableCell>
                  <TableCell>
                    {child.matched && child.child_id ? (
                      <Link
                        href={`/organizations/${orgId}/children/${child.child_id}`}
                        className="text-primary hover:underline"
                      >
                        {child.child_name}
                      </Link>
                    ) : (
                      child.child_name
                    )}
                  </TableCell>
                  <TableCell className="hidden md:table-cell">{child.birth_date}</TableCell>
                  <TableCell className="hidden md:table-cell">{child.district}</TableCell>
                  <TableCell className="text-right">{formatCurrency(child.total_amount)}</TableCell>
                  <TableCell>
                    {child.matched ? (
                      <Badge variant="success">
                        <CheckCircle2 className="mr-1 h-3 w-3" />
                      </Badge>
                    ) : (
                      <Badge variant="destructive">
                        <XCircle className="mr-1 h-3 w-3" />
                      </Badge>
                    )}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
