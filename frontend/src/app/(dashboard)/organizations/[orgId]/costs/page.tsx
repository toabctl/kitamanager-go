'use client';

import { useMemo, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery } from '@tanstack/react-query';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Button } from '@/components/ui/button';
import { apiClient } from '@/lib/api/client';
import { queryKeys } from '@/lib/api/queryKeys';
import type { Cost, CostCreateRequest, CostUpdateRequest } from '@/lib/api/types';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useCrudMutations } from '@/lib/hooks/use-crud-mutations';
import { useCrudDialogs } from '@/lib/hooks/use-crud-dialogs';
import { CrudPageHeader, ResourceTable, DeleteConfirmDialog, Column } from '@/components/crud';
import { Pagination } from '@/components/ui/pagination';
import { costSchema, type CostFormData } from '@/lib/schemas';

const defaultValues: CostFormData = {
  name: '',
};

export default function CostsPage() {
  const params = useParams();
  const router = useRouter();
  const orgId = Number(params.orgId);
  const t = useTranslations();
  const [page, setPage] = useState(1);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<CostFormData>({
    resolver: zodResolver(costSchema),
    defaultValues,
  });

  const { data: paginatedData, isLoading } = useQuery({
    queryKey: queryKeys.costs.list(orgId, page),
    queryFn: () => apiClient.getCosts(orgId, { page }),
    enabled: !!orgId,
  });

  const costs = paginatedData?.data;

  const dialogs = useCrudDialogs<Cost, CostFormData>({
    reset,
    itemToFormData: (cost) => ({ name: cost.name }),
    defaultValues,
  });

  const mutations = useCrudMutations<Cost, CostCreateRequest, CostUpdateRequest>({
    resourceName: 'costs',
    queryKey: queryKeys.costs.all(orgId),
    createFn: (data) => apiClient.createCost(orgId, data),
    updateFn: (id, data) => apiClient.updateCost(orgId, id, data),
    deleteFn: (id) => apiClient.deleteCost(orgId, id),
    onSuccess: dialogs.closeDialog,
    onDeleteSuccess: dialogs.closeDeleteDialog,
  });

  const onSubmit = (data: CostFormData) => {
    if (dialogs.editingItem) {
      mutations.updateMutation.mutate({ id: dialogs.editingItem.id, data });
    } else {
      mutations.createMutation.mutate(data);
    }
  };

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
        onNew={dialogs.handleCreate}
        newButtonText="costs.newCost"
      />

      <Card>
        <CardHeader>
          <CardTitle>{t('costs.title')}</CardTitle>
        </CardHeader>
        <CardContent>
          <ResourceTable
            items={costs}
            columns={columns}
            getItemKey={(cost) => cost.id}
            isLoading={isLoading}
            onView={handleView}
            onEdit={dialogs.handleEdit}
            onDelete={dialogs.handleDelete}
          />
          {paginatedData && (
            <Pagination
              page={paginatedData.page}
              totalPages={paginatedData.total_pages}
              total={paginatedData.total}
              limit={paginatedData.limit}
              onPageChange={setPage}
              isLoading={isLoading}
            />
          )}
        </CardContent>
      </Card>

      <Dialog open={dialogs.isDialogOpen} onOpenChange={dialogs.setIsDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{dialogs.isEditing ? t('costs.edit') : t('costs.create')}</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">{t('common.name')}</Label>
              <Input id="name" {...register('name')} />
              {errors.name && (
                <p className="text-sm text-destructive">{t('validation.nameRequired')}</p>
              )}
            </div>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => dialogs.setIsDialogOpen(false)}
              >
                {t('common.cancel')}
              </Button>
              <Button type="submit" disabled={mutations.isMutating}>
                {t('common.save')}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      <DeleteConfirmDialog
        open={dialogs.isDeleteDialogOpen}
        onOpenChange={dialogs.setIsDeleteDialogOpen}
        onConfirm={() =>
          dialogs.deletingItem && mutations.deleteMutation.mutate(dialogs.deletingItem.id)
        }
        isLoading={mutations.deleteMutation.isPending}
        resourceName="costs"
      />
    </div>
  );
}
