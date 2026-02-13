'use client';

import { useTranslations } from 'next-intl';
import { useQuery } from '@tanstack/react-query';
import { Building2, Users, Baby, UserCog } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { StepPromotionsWidget } from '@/components/dashboard/step-promotions-widget';
import { UpcomingChildrenWidget } from '@/components/dashboard/upcoming-children-widget';
import { SectionAgeAlertsWidget } from '@/components/dashboard/section-age-alerts-widget';
import { apiClient } from '@/lib/api/client';
import { queryKeys } from '@/lib/api/queryKeys';
import { useUiStore } from '@/stores/ui-store';
import { useAuthStore } from '@/stores/auth-store';

function StatCard({
  title,
  value,
  icon: Icon,
  loading,
}: {
  title: string;
  value: number | string;
  icon: React.ComponentType<{ className?: string }>;
  loading?: boolean;
}) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        <Icon className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        {loading ? (
          <Skeleton className="h-8 w-20" />
        ) : (
          <div className="text-2xl font-bold">{value}</div>
        )}
      </CardContent>
    </Card>
  );
}

export default function DashboardPage() {
  const t = useTranslations();
  const { organizations, organizationsLoading, selectedOrganizationId, getSelectedOrganization } =
    useUiStore();
  const { user } = useAuthStore();
  const selectedOrg = getSelectedOrganization();

  const { data: employeesData, isLoading: employeesLoading } = useQuery({
    queryKey: [...queryKeys.employees.list(selectedOrganizationId ?? 0, 1), 'count'],
    queryFn: () => apiClient.getEmployees(selectedOrganizationId!, { page: 1, limit: 1 }),
    enabled: !!selectedOrganizationId,
    staleTime: 2 * 60 * 1000,
  });

  const { data: childrenData, isLoading: childrenLoading } = useQuery({
    queryKey: [...queryKeys.children.list(selectedOrganizationId ?? 0, 1), 'count'],
    queryFn: () => apiClient.getChildren(selectedOrganizationId!, { page: 1, limit: 1 }),
    enabled: !!selectedOrganizationId,
    staleTime: 2 * 60 * 1000,
  });

  const { data: usersData, isLoading: usersLoading } = useQuery({
    queryKey: [...queryKeys.users.list(1), 'count'],
    queryFn: () => apiClient.getUsers({ page: 1, limit: 1 }),
    staleTime: 5 * 60 * 1000,
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">{t('dashboard.title')}</h1>
        <p className="text-muted-foreground">
          {t('dashboard.welcome')}
          {user?.name && `, ${user.name}`}
        </p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <StatCard
          title={t('dashboard.totalOrganizations')}
          value={organizations.length}
          icon={Building2}
          loading={organizationsLoading}
        />
        <StatCard
          title={t('dashboard.totalEmployees')}
          value={employeesData?.total ?? '-'}
          icon={Users}
          loading={employeesLoading}
        />
        <StatCard
          title={t('dashboard.totalChildren')}
          value={childrenData?.total ?? '-'}
          icon={Baby}
          loading={childrenLoading}
        />
        <StatCard
          title={t('dashboard.totalUsers')}
          value={usersData?.total ?? '-'}
          icon={UserCog}
          loading={usersLoading}
        />
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{t('dashboard.quickStats')}</CardTitle>
        </CardHeader>
        <CardContent>
          {selectedOrganizationId && selectedOrg ? (
            <p className="text-muted-foreground">
              {t('dashboard.statsForOrg', { name: selectedOrg.name })}
            </p>
          ) : (
            <p className="text-muted-foreground">{t('statistics.selectOrgForStats')}</p>
          )}
        </CardContent>
      </Card>

      {selectedOrganizationId && <StepPromotionsWidget orgId={selectedOrganizationId} />}
      {selectedOrganizationId && <UpcomingChildrenWidget orgId={selectedOrganizationId} />}
      {selectedOrganizationId && <SectionAgeAlertsWidget orgId={selectedOrganizationId} />}
    </div>
  );
}
