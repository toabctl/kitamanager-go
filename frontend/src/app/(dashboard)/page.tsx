'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { Skeleton } from '@/components/ui/skeleton';
import { useUiStore } from '@/stores/ui-store';

export default function RootRedirectPage() {
  const router = useRouter();
  const t = useTranslations();
  const { organizations, organizationsLoading, selectedOrganizationId } = useUiStore();

  useEffect(() => {
    if (organizationsLoading) return;

    const orgId = selectedOrganizationId ?? organizations[0]?.id;
    if (orgId) {
      router.replace(`/organizations/${orgId}/dashboard`);
    }
  }, [organizationsLoading, selectedOrganizationId, organizations, router]);

  if (organizationsLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-10 w-48" />
        <Skeleton className="h-32 w-full" />
      </div>
    );
  }

  if (organizations.length === 0) {
    return (
      <div className="space-y-6">
        <h1 className="text-3xl font-bold tracking-tight">{t('dashboard.title')}</h1>
        <p className="text-muted-foreground">{t('statistics.selectOrgForStats')}</p>
      </div>
    );
  }

  return null;
}
