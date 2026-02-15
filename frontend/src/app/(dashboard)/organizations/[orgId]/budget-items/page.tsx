'use client';

import { useMemo } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { Check } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Checkbox } from '@/components/ui/checkbox';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { apiClient } from '@/lib/api/client';
import { queryKeys } from '@/lib/api/queryKeys';
import type { BudgetItem, BudgetItemCreateRequest, BudgetItemUpdateRequest } from '@/lib/api/types';
import { useCrudPage } from '@/lib/hooks/use-crud-page';
import {
  CrudPageHeader,
  ResourceTable,
  DeleteConfirmDialog,
  CrudFormDialog,
  Column,
} from '@/components/crud';
import { Pagination } from '@/components/ui/pagination';
import { budgetItemSchema, type BudgetItemFormData } from '@/lib/schemas';

const defaultValues: BudgetItemFormData = {
  name: '',
  category: 'expense',
  per_child: false,
};

export default function BudgetItemsPage() {
  const router = useRouter();
  const params = useParams();
  const orgId = Number(params.orgId);
  const t = useTranslations();

  const crud = useCrudPage<
    BudgetItem,
    BudgetItemFormData,
    BudgetItemCreateRequest,
    BudgetItemUpdateRequest
  >({
    resourceName: 'budgetItems',
    schema: budgetItemSchema,
    defaultValues,
    itemToFormData: (item) => ({
      name: item.name,
      category: item.category,
      per_child: item.per_child,
    }),
    listFn: (orgId, params) => apiClient.getBudgetItems(orgId, params),
    createFn: (orgId, data) => apiClient.createBudgetItem(orgId, data),
    updateFn: (orgId, id, data) => apiClient.updateBudgetItem(orgId, id, data),
    deleteFn: (orgId, id) => apiClient.deleteBudgetItem(orgId, id),
    queryKeys: {
      list: (orgId, page) => queryKeys.budgetItems.list(orgId, page),
      invalidate: (orgId) => queryKeys.budgetItems.all(orgId),
    },
  });

  const handleView = (item: BudgetItem) => {
    router.push(`/organizations/${orgId}/budget-items/${item.id}`);
  };

  const columns = useMemo<Column<BudgetItem>[]>(
    () => [
      { key: 'id', header: 'common.id', render: (item) => item.id },
      {
        key: 'name',
        header: 'common.name',
        render: (item) => item.name,
        className: 'font-medium',
      },
      {
        key: 'category',
        header: 'budgetItems.category',
        render: (item) => (
          <Badge variant={item.category === 'income' ? 'default' : 'secondary'}>
            {t(`budgetItems.category${item.category === 'income' ? 'Income' : 'Expense'}`)}
          </Badge>
        ),
      },
      {
        key: 'per_child',
        header: 'budgetItems.perChild',
        render: (item) =>
          item.per_child ? <Check className="h-4 w-4 text-muted-foreground" /> : null,
      },
    ],
    [t]
  );

  return (
    <div className="space-y-6">
      <CrudPageHeader
        title="budgetItems.title"
        onNew={crud.dialogs.handleCreate}
        newButtonText="budgetItems.newBudgetItem"
      />

      <Card>
        <CardHeader>
          <CardTitle>{t('budgetItems.title')}</CardTitle>
        </CardHeader>
        <CardContent>
          <ResourceTable
            items={crud.items}
            columns={columns}
            getItemKey={(item) => item.id}
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
        translationPrefix="budgetItems"
        onSubmit={crud.handleSubmit(crud.onSubmit)}
        isSaving={crud.mutations.isMutating}
      >
        <div className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">{t('common.name')}</Label>
            <Input id="name" {...crud.register('name')} />
            {crud.errors.name && (
              <p className="text-sm text-destructive">{t('validation.nameRequired')}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="category">{t('budgetItems.category')}</Label>
            <Select
              value={crud.watch('category')}
              onValueChange={(value) =>
                crud.setValue('category', value as 'income' | 'expense', { shouldValidate: true })
              }
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="income">{t('budgetItems.categoryIncome')}</SelectItem>
                <SelectItem value="expense">{t('budgetItems.categoryExpense')}</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="flex items-center space-x-2">
            <Checkbox
              id="per_child"
              checked={crud.watch('per_child')}
              onCheckedChange={(checked) =>
                crud.setValue('per_child', checked === true, { shouldValidate: true })
              }
            />
            <Label htmlFor="per_child">{t('budgetItems.perChild')}</Label>
          </div>
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
        resourceName="budgetItems"
      />
    </div>
  );
}
