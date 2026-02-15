'use client';

import { useMemo } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { apiClient } from '@/lib/api/client';
import { queryKeys } from '@/lib/api/queryKeys';
import type { Cost, CostCreateRequest, CostUpdateRequest } from '@/lib/api/types';
import { useCrudPage } from '@/lib/hooks/use-crud-page';
import {
  CrudPageHeader,
  ResourceTable,
  DeleteConfirmDialog,
  CrudFormDialog,
  Column,
} from '@/components/crud';
import { Pagination } from '@/components/ui/pagination';
import { costSchema, type CostFormData } from '@/lib/schemas';

const defaultValues: CostFormData = {
  name: '',
};

export default function CostsPage() {
  const router = useRouter();
  const params = useParams();
  const orgId = Number(params.orgId);
  const t = useTranslations();

  const crud = useCrudPage<Cost, CostFormData, CostCreateRequest, CostUpdateRequest>({
    resourceName: 'costs',
    schema: costSchema,
    defaultValues,
    itemToFormData: (cost) => ({ name: cost.name }),
    listFn: (orgId, params) => apiClient.getCosts(orgId, params),
    createFn: (orgId, data) => apiClient.createCost(orgId, data),
    updateFn: (orgId, id, data) => apiClient.updateCost(orgId, id, data),
    deleteFn: (orgId, id) => apiClient.deleteCost(orgId, id),
    queryKeys: {
      list: (orgId, page) => queryKeys.costs.list(orgId, page),
      invalidate: (orgId) => queryKeys.costs.all(orgId),
    },
  });

  const handleView = (cost: Cost) => {
    router.push(`/organizations/${orgId}/costs/${cost.id}`);
  };

  const columns = useMemo<Column<Cost>[]>(
    () => [
      { key: 'id', header: 'common.id', render: (cost) => cost.id },
      {
        key: 'name',
        header: 'common.name',
        render: (cost) => cost.name,
        className: 'font-medium',
      },
    ],
    []
  );

  return (
    <div className="space-y-6">
      <CrudPageHeader
        title="costs.title"
        onNew={crud.dialogs.handleCreate}
        newButtonText="costs.newCost"
      />

      <Card>
        <CardHeader>
          <CardTitle>{t('costs.title')}</CardTitle>
        </CardHeader>
        <CardContent>
          <ResourceTable
            items={crud.items}
            columns={columns}
            getItemKey={(cost) => cost.id}
            isLoading={crud.isLoading}
            onView={handleView}
            onEdit={crud.dialogs.handleEdit}
            onDelete={crud.dialogs.handleDelete}
          />
          {crud.paginatedData && (
            <Pagination
              page={crud.paginatedData.page}
              totalPages={crud.paginatedData.total_pages}
              total={crud.paginatedData.total}
              limit={crud.paginatedData.limit}
              onPageChange={crud.setPage}
              isLoading={crud.isLoading}
            />
          )}
        </CardContent>
      </Card>

      <CrudFormDialog
        open={crud.dialogs.isDialogOpen}
        onOpenChange={crud.dialogs.setIsDialogOpen}
        isEditing={crud.dialogs.isEditing}
        translationPrefix="costs"
        onSubmit={crud.handleSubmit(crud.onSubmit)}
        isSaving={crud.mutations.isMutating}
      >
        <div className="space-y-2">
          <Label htmlFor="name">{t('common.name')}</Label>
          <Input id="name" {...crud.register('name')} />
          {crud.errors.name && (
            <p className="text-sm text-destructive">{t('validation.nameRequired')}</p>
          )}
        </div>
      </CrudFormDialog>

      <DeleteConfirmDialog
        open={crud.dialogs.isDeleteDialogOpen}
        onOpenChange={crud.dialogs.setIsDeleteDialogOpen}
        onConfirm={() =>
          crud.dialogs.deletingItem &&
          crud.mutations.deleteMutation.mutate(crud.dialogs.deletingItem.id)
        }
        isLoading={crud.mutations.deleteMutation.isPending}
        resourceName="costs"
      />
    </div>
  );
}
