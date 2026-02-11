'use client';

import { useTranslations } from 'next-intl';
import { AlertTriangle, RefreshCw } from 'lucide-react';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';

interface QueryErrorProps {
  error: Error | null;
  onRetry: () => void;
}

export function QueryError({ error, onRetry }: QueryErrorProps) {
  const t = useTranslations();

  if (!error) return null;

  return (
    <Alert variant="destructive">
      <AlertTriangle className="h-4 w-4" />
      <AlertTitle>{t('common.error')}</AlertTitle>
      <AlertDescription className="flex items-center justify-between">
        <span>{error.message}</span>
        <Button variant="outline" size="sm" onClick={onRetry}>
          <RefreshCw className="mr-2 h-4 w-4" />
          {t('common.retry')}
        </Button>
      </AlertDescription>
    </Alert>
  );
}
