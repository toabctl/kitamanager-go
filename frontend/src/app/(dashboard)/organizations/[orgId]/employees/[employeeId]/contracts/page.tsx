'use client';

import { useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { ArrowLeft, Pencil, Trash2, Plus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useToast } from '@/lib/hooks/use-toast';
import { apiClient, getErrorMessage } from '@/lib/api/client';
import type {
  EmployeeContract,
  EmployeeContractCreateRequest,
  EmployeeContractUpdateRequest,
} from '@/lib/api/types';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { formatDate, formatDateForInput, formatDateForApi } from '@/lib/utils/formatting';

const contractSchema = z.object({
  from: z.string().min(1),
  to: z.string().optional(),
  position: z.string().min(1),
  grade: z.string().min(1),
  step: z.number().min(1).max(6),
  weekly_hours: z.number().min(0).max(168),
});

type ContractFormData = z.infer<typeof contractSchema>;

export default function EmployeeContractsPage() {
  const params = useParams();
  const router = useRouter();
  const orgId = Number(params.orgId);
  const employeeId = Number(params.employeeId);
  const t = useTranslations();
  const { toast } = useToast();
  const queryClient = useQueryClient();

  const [isContractDialogOpen, setIsContractDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [editingContract, setEditingContract] = useState<EmployeeContract | null>(null);
  const [deletingContract, setDeletingContract] = useState<EmployeeContract | null>(null);

  // Fetch employee data
  const { data: employee, isLoading: employeeLoading } = useQuery({
    queryKey: ['employee', orgId, employeeId],
    queryFn: () => apiClient.getEmployee(orgId, employeeId),
    enabled: !!orgId && !!employeeId,
  });

  // Fetch contracts
  const { data: contracts, isLoading: contractsLoading } = useQuery({
    queryKey: ['employeeContracts', orgId, employeeId],
    queryFn: () => apiClient.getEmployeeContracts(orgId, employeeId),
    enabled: !!orgId && !!employeeId,
  });

  const createMutation = useMutation({
    mutationFn: (data: EmployeeContractCreateRequest) =>
      apiClient.createEmployeeContract(orgId, employeeId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['employeeContracts', orgId, employeeId] });
      queryClient.invalidateQueries({ queryKey: ['employee', orgId, employeeId] });
      toast({ title: t('contracts.createSuccess') });
      setIsContractDialogOpen(false);
      setEditingContract(null);
      reset();
    },
    onError: (error) => {
      toast({
        title: t('common.error'),
        description: getErrorMessage(error, t('common.failedToCreate', { resource: 'contract' })),
        variant: 'destructive',
      });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({
      contractId,
      data,
    }: {
      contractId: number;
      data: EmployeeContractUpdateRequest;
    }) => apiClient.updateEmployeeContract(orgId, employeeId, contractId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['employeeContracts', orgId, employeeId] });
      queryClient.invalidateQueries({ queryKey: ['employee', orgId, employeeId] });
      toast({ title: t('contracts.updateSuccess') });
      setIsContractDialogOpen(false);
      setEditingContract(null);
      reset();
    },
    onError: (error) => {
      toast({
        title: t('common.error'),
        description: getErrorMessage(error, t('common.failedToSave', { resource: 'contract' })),
        variant: 'destructive',
      });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (contractId: number) =>
      apiClient.deleteEmployeeContract(orgId, employeeId, contractId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['employeeContracts', orgId, employeeId] });
      queryClient.invalidateQueries({ queryKey: ['employee', orgId, employeeId] });
      toast({ title: t('contracts.deleteSuccess') });
      setIsDeleteDialogOpen(false);
      setDeletingContract(null);
    },
    onError: (error) => {
      toast({
        title: t('common.error'),
        description: getErrorMessage(error, t('common.failedToDelete', { resource: 'contract' })),
        variant: 'destructive',
      });
    },
  });

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<ContractFormData>({
    resolver: zodResolver(contractSchema),
    defaultValues: {
      from: '',
      to: '',
      position: '',
      grade: '',
      step: 1,
      weekly_hours: 39,
    },
  });

  const handleCreate = () => {
    setEditingContract(null);
    reset({ from: '', to: '', position: '', grade: '', step: 1, weekly_hours: 39 });
    setIsContractDialogOpen(true);
  };

  const handleEdit = (contract: EmployeeContract) => {
    setEditingContract(contract);
    reset({
      from: formatDateForInput(contract.from),
      to: contract.to ? formatDateForInput(contract.to) : '',
      position: contract.position,
      grade: contract.grade,
      step: contract.step,
      weekly_hours: contract.weekly_hours,
    });
    setIsContractDialogOpen(true);
  };

  const handleDelete = (contract: EmployeeContract) => {
    setDeletingContract(contract);
    setIsDeleteDialogOpen(true);
  };

  const onSubmit = (data: ContractFormData) => {
    if (editingContract) {
      updateMutation.mutate({
        contractId: editingContract.id,
        data: {
          from: formatDateForApi(data.from) || undefined,
          to: formatDateForApi(data.to) || undefined,
          position: data.position,
          grade: data.grade,
          step: data.step,
          weekly_hours: data.weekly_hours,
        },
      });
    } else {
      createMutation.mutate({
        from: formatDateForApi(data.from) || data.from,
        to: formatDateForApi(data.to),
        position: data.position,
        grade: data.grade,
        step: data.step,
        weekly_hours: data.weekly_hours,
      });
    }
  };

  const getContractStatus = (contract: EmployeeContract): 'active' | 'upcoming' | 'ended' => {
    const today = new Date().toISOString().split('T')[0];
    if (contract.from > today) return 'upcoming';
    if (contract.to && contract.to < today) return 'ended';
    return 'active';
  };

  const isLoading = employeeLoading || contractsLoading;

  // Sort contracts by start date descending (most recent first)
  const sortedContracts = contracts
    ? [...contracts].sort((a, b) => b.from.localeCompare(a.from))
    : [];

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="icon"
          onClick={() => router.push(`/organizations/${orgId}/employees`)}
        >
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <div>
          <h1 className="text-3xl font-bold tracking-tight">{t('employees.contractHistory')}</h1>
          {employee && (
            <p className="text-muted-foreground">
              {employee.first_name} {employee.last_name}
            </p>
          )}
        </div>
        <div className="ml-auto">
          <Button onClick={handleCreate}>
            <Plus className="mr-2 h-4 w-4" />
            {t('contracts.newContract')}
          </Button>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{t('contracts.title')}</CardTitle>
          <CardDescription>
            {sortedContracts.length > 0
              ? t('employees.contractHistory')
              : t('employees.noContractsFound')}
          </CardDescription>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="space-y-2">
              {[...Array(3)].map((_, i) => (
                <Skeleton key={i} className="h-12 w-full" />
              ))}
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>{t('common.status')}</TableHead>
                  <TableHead>{t('contracts.from')}</TableHead>
                  <TableHead>{t('contracts.to')}</TableHead>
                  <TableHead>{t('employees.position')}</TableHead>
                  <TableHead>{t('employees.grade')}</TableHead>
                  <TableHead>{t('employees.weeklyHours')}</TableHead>
                  <TableHead className="text-right">{t('common.actions')}</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {sortedContracts.map((contract) => {
                  const status = getContractStatus(contract);
                  return (
                    <TableRow key={contract.id}>
                      <TableCell>
                        <Badge
                          variant={
                            status === 'active'
                              ? 'success'
                              : status === 'upcoming'
                                ? 'warning'
                                : 'secondary'
                          }
                        >
                          {status === 'active'
                            ? t('common.active')
                            : status === 'upcoming'
                              ? t('common.upcoming')
                              : t('common.ended')}
                        </Badge>
                      </TableCell>
                      <TableCell>{formatDate(contract.from)}</TableCell>
                      <TableCell>
                        {contract.to ? formatDate(contract.to) : t('common.ongoing')}
                      </TableCell>
                      <TableCell>{contract.position}</TableCell>
                      <TableCell>
                        {contract.grade} / {contract.step}
                      </TableCell>
                      <TableCell>{contract.weekly_hours}h</TableCell>
                      <TableCell className="text-right">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleEdit(contract)}
                          aria-label={t('common.edit')}
                        >
                          <Pencil className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleDelete(contract)}
                          aria-label={t('common.delete')}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </TableCell>
                    </TableRow>
                  );
                })}
                {sortedContracts.length === 0 && (
                  <TableRow>
                    <TableCell colSpan={7} className="text-center text-muted-foreground">
                      {t('employees.noContractsFound')}
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      {/* Contract Create/Edit Dialog */}
      <Dialog open={isContractDialogOpen} onOpenChange={setIsContractDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {editingContract ? t('contracts.edit') : t('contracts.create')}
            </DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="from">{t('contracts.startDate')}</Label>
                <Input id="from" type="date" {...register('from')} />
                {errors.from && (
                  <p className="text-sm text-destructive">{t('contracts.startDateRequired')}</p>
                )}
              </div>
              <div className="space-y-2">
                <Label htmlFor="to">{t('contracts.endDateOptional')}</Label>
                <Input id="to" type="date" {...register('to')} />
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="position">{t('employees.position')}</Label>
              <Input id="position" {...register('position')} />
              {errors.position && (
                <p className="text-sm text-destructive">{t('validation.positionRequired')}</p>
              )}
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="grade">{t('employees.grade')}</Label>
                <Input id="grade" {...register('grade')} placeholder="S8a" />
                {errors.grade && (
                  <p className="text-sm text-destructive">{t('payPlans.gradeRequired')}</p>
                )}
              </div>
              <div className="space-y-2">
                <Label htmlFor="step">{t('employees.step')}</Label>
                <Input
                  id="step"
                  type="number"
                  min={1}
                  max={6}
                  {...register('step', { valueAsNumber: true })}
                />
                {errors.step && (
                  <p className="text-sm text-destructive">{t('payPlans.stepRequired')}</p>
                )}
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="weekly_hours">{t('employees.weeklyHours')}</Label>
              <Input
                id="weekly_hours"
                type="number"
                min={0}
                max={168}
                step={0.5}
                {...register('weekly_hours', { valueAsNumber: true })}
              />
              {errors.weekly_hours && (
                <p className="text-sm text-destructive">{t('validation.weeklyHoursRequired')}</p>
              )}
            </div>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setIsContractDialogOpen(false)}
              >
                {t('common.cancel')}
              </Button>
              <Button type="submit" disabled={createMutation.isPending || updateMutation.isPending}>
                {t('common.save')}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{t('common.confirmDelete')}</AlertDialogTitle>
            <AlertDialogDescription>{t('contracts.deleteConfirm')}</AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{t('common.cancel')}</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => deletingContract && deleteMutation.mutate(deletingContract.id)}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              {t('common.delete')}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
