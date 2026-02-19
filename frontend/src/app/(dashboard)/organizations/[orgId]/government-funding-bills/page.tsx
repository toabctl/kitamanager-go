'use client';

import { useState, useRef } from 'react';
import Link from 'next/link';
import { useParams } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useMutation } from '@tanstack/react-query';
import { Upload, CheckCircle2, XCircle, RotateCcw } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { apiClient, getErrorMessage } from '@/lib/api/client';
import type { GovernmentFundingBillResponse } from '@/lib/api/types';
import { useToast } from '@/lib/hooks/use-toast';
import { formatCurrency } from '@/lib/utils/formatting';

export default function GovernmentFundingBillsPage() {
  const params = useParams();
  const orgId = Number(params.orgId);
  const t = useTranslations('governmentFundingBills');
  const { toast } = useToast();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  const uploadMutation = useMutation({
    mutationFn: async (file: File) => {
      const formData = new FormData();
      formData.append('file', file);
      return apiClient.uploadGovernmentFundingBill(orgId, formData);
    },
    onError: (error) => {
      toast({
        title: t('uploadError'),
        description: getErrorMessage(error, t('uploadError')),
        variant: 'destructive',
      });
    },
  });

  const handleUpload = () => {
    if (selectedFile) {
      uploadMutation.mutate(selectedFile);
    }
  };

  const handleReset = () => {
    uploadMutation.reset();
    setSelectedFile(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const result: GovernmentFundingBillResponse | undefined = uploadMutation.data;

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">{t('title')}</h1>

      {!result ? (
        <Card>
          <CardHeader>
            <CardTitle>{t('selectFile')}</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col gap-4 sm:flex-row sm:items-end">
              <div className="flex-1">
                <Input
                  ref={fileInputRef}
                  type="file"
                  accept=".xlsx"
                  onChange={(e) => setSelectedFile(e.target.files?.[0] ?? null)}
                />
              </div>
              <Button onClick={handleUpload} disabled={!selectedFile || uploadMutation.isPending}>
                <Upload className="mr-2 h-4 w-4" />
                {uploadMutation.isPending ? t('uploading') : t('upload')}
              </Button>
            </div>
          </CardContent>
        </Card>
      ) : (
        <>
          {/* Summary Cards */}
          <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-muted-foreground text-sm font-medium">
                  {t('facilityName')}
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-lg font-semibold">{result.facility_name}</p>
              </CardContent>
            </Card>
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
          </div>

          <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-muted-foreground text-sm font-medium">
                  {t('matchedChildren')}
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-lg font-semibold text-green-600">
                  {result.matched_count} / {result.children_count}
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-muted-foreground text-sm font-medium">
                  {t('unmatchedChildren')}
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-lg font-semibold text-red-600">{result.unmatched_count}</p>
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
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle>
                {t('children')} ({result.children_count})
              </CardTitle>
              <Button variant="outline" size="sm" onClick={handleReset}>
                <RotateCcw className="mr-2 h-4 w-4" />
                {t('reset')}
              </Button>
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
                      <TableCell className="text-right">
                        {formatCurrency(child.total_amount)}
                      </TableCell>
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
        </>
      )}
    </div>
  );
}
