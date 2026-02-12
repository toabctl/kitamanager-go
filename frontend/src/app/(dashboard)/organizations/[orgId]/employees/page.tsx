'use client';

import { useState, useCallback, useMemo } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Plus, Pencil, Trash2, FileText, History, Search } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Skeleton } from '@/components/ui/skeleton';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { useToast } from '@/lib/hooks/use-toast';
import { apiClient, getErrorMessage } from '@/lib/api/client';
import type {
  Employee,
  EmployeeContract,
  EmployeeContractCreateRequest,
  PayPlan,
} from '@/lib/api/types';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Checkbox } from '@/components/ui/checkbox';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import {
  formatDate,
  calculateAge,
  formatDateForInput,
  formatDateForApi,
} from '@/lib/utils/formatting';
import { getActiveContract, getCurrentContract, getDayBefore } from '@/lib/utils/contracts';
import { calculateMonthlySalary } from '@/lib/utils/salary';
import { calculateYearsOfService } from '@/lib/utils/step-promotions';
import { formatCurrency } from '@/lib/utils/formatting';
import { Pagination } from '@/components/ui/pagination';
import { useDebouncedValue } from '@/lib/hooks/use-debounced-value';
import { DeleteConfirmDialog } from '@/components/crud/delete-confirm-dialog';
import { PersonFormDialog } from '@/components/crud/person-form-dialog';

const employeeSchema = z.object({
  first_name: z.string().min(1),
  last_name: z.string().min(1),
  gender: z.enum(['male', 'female', 'diverse']),
  birthdate: z.string().min(1),
});

const contractSchema = z.object({
  from: z.string().min(1),
  to: z.string().optional(),
  payplan_id: z.number().min(1),
  staff_category: z.enum(['qualified', 'supplementary', 'non_pedagogical']),
  grade: z.string().min(1),
  step: z.number().min(1).max(6),
  weekly_hours: z.number().min(0).max(168),
});

type EmployeeFormData = z.infer<typeof employeeSchema>;
type ContractFormData = z.infer<typeof contractSchema>;

export default function EmployeesPage() {
  const params = useParams();
  const router = useRouter();
  const orgId = Number(params.orgId);
  const t = useTranslations();
  const { toast } = useToast();
  const queryClient = useQueryClient();

  const [isEmployeeDialogOpen, setIsEmployeeDialogOpen] = useState(false);
  const [isContractDialogOpen, setIsContractDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [editingEmployee, setEditingEmployee] = useState<Employee | null>(null);
  const [deletingEmployee, setDeletingEmployee] = useState<Employee | null>(null);
  const [contractEmployee, setContractEmployee] = useState<Employee | null>(null);
  const [endCurrentContract, setEndCurrentContract] = useState(true);
  const [page, setPage] = useState(1);
  const [searchInput, setSearchInput] = useState('');
  const search = useDebouncedValue(searchInput, 300);
  const [staffCategoryFilter, setStaffCategoryFilter] = useState<string>('');

  const { data: paginatedData, isLoading } = useQuery({
    queryKey: ['employees', orgId, page, search, staffCategoryFilter],
    queryFn: () =>
      apiClient.getEmployees(orgId, {
        page,
        search: search || undefined,
        staff_category: staffCategoryFilter || undefined,
      }),
    enabled: !!orgId,
  });

  const employees = paginatedData?.data;

  const { data: payPlansData } = useQuery({
    queryKey: ['payplans', orgId],
    queryFn: () => apiClient.getPayPlans(orgId, { limit: 100 }),
    enabled: !!orgId,
  });
  const payPlans = useMemo(() => payPlansData?.data ?? [], [payPlansData?.data]);

  // Collect unique payplan IDs from employee contracts to fetch details for salary calc
  const payPlanIds = Array.from(
    new Set(
      (employees ?? [])
        .flatMap((e) => e.contracts ?? [])
        .map((c) => c.payplan_id)
        .filter((id) => id > 0)
    )
  );

  const payPlanDetailsQueries = useQuery({
    queryKey: ['payplanDetails', orgId, payPlanIds],
    queryFn: async () => {
      const results = await Promise.all(payPlanIds.map((id) => apiClient.getPayPlan(orgId, id)));
      const map = new Map<number, PayPlan>();
      for (const pp of results) {
        map.set(pp.id, pp);
      }
      return map;
    },
    enabled: payPlanIds.length > 0,
  });
  const payPlanMap = payPlanDetailsQueries.data ?? new Map<number, PayPlan>();

  const createMutation = useMutation({
    mutationFn: (data: Omit<EmployeeFormData, 'organization_id'>) =>
      apiClient.createEmployee(orgId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['employees', orgId] });
      toast({ title: t('employees.createSuccess') });
      setIsEmployeeDialogOpen(false);
      resetEmployee();
    },
    onError: (error) => {
      toast({
        title: t('common.error'),
        description: getErrorMessage(error, t('common.failedToCreate', { resource: 'employee' })),
        variant: 'destructive',
      });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: Partial<EmployeeFormData> }) =>
      apiClient.updateEmployee(orgId, id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['employees', orgId] });
      toast({ title: t('employees.updateSuccess') });
      setIsEmployeeDialogOpen(false);
      setEditingEmployee(null);
      resetEmployee();
    },
    onError: (error) => {
      toast({
        title: t('common.error'),
        description: getErrorMessage(error, t('common.failedToSave', { resource: 'employee' })),
        variant: 'destructive',
      });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => apiClient.deleteEmployee(orgId, id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['employees', orgId] });
      toast({ title: t('employees.deleteSuccess') });
      setIsDeleteDialogOpen(false);
      setDeletingEmployee(null);
    },
    onError: (error) => {
      toast({
        title: t('common.error'),
        description: getErrorMessage(error, t('common.failedToDelete', { resource: 'employee' })),
        variant: 'destructive',
      });
    },
  });

  const createContractMutation = useMutation({
    mutationFn: async ({
      employeeId,
      data,
    }: {
      employeeId: number;
      data: EmployeeContractCreateRequest;
    }) => {
      // If we need to end the current contract first
      if (contractEmployee && endCurrentContract) {
        const active = getActiveContract(contractEmployee.contracts);
        if (active && data.from) {
          const endDate = getDayBefore(data.from);
          await apiClient.updateEmployeeContract(orgId, employeeId, active.id, {
            to: formatDateForApi(endDate),
          });
        }
      }
      return apiClient.createEmployeeContract(orgId, employeeId, data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['employees', orgId] });
      toast({
        title: endCurrentContract
          ? t('contracts.previousContractEnded')
          : t('contracts.createSuccess'),
      });
      setIsContractDialogOpen(false);
      setContractEmployee(null);
      setEndCurrentContract(true);
      resetContract();
    },
    onError: (error) => {
      toast({
        title: t('common.error'),
        description: getErrorMessage(error, t('common.failedToCreate', { resource: 'contract' })),
        variant: 'destructive',
      });
    },
  });

  const {
    register: registerEmployee,
    handleSubmit: handleSubmitEmployee,
    reset: resetEmployee,
    setValue: setValueEmployee,
    watch: watchEmployee,
    formState: { errors: errorsEmployee },
  } = useForm<EmployeeFormData>({
    resolver: zodResolver(employeeSchema),
    defaultValues: {
      first_name: '',
      last_name: '',
      gender: 'male',
      birthdate: '',
    },
  });

  const {
    register: registerContract,
    handleSubmit: handleSubmitContract,
    reset: resetContract,
    watch: watchContract,
    setValue: setValueContract,
    formState: { errors: errorsContract },
  } = useForm<ContractFormData>({
    resolver: zodResolver(contractSchema),
    defaultValues: {
      from: '',
      to: '',
      payplan_id: 0,
      staff_category: 'qualified',
      grade: '',
      step: 1,
      weekly_hours: 39,
    },
  });

  const contractFromDate = watchContract('from');

  // Calculate end date preview based on contract from date
  const activeContract = contractEmployee ? getActiveContract(contractEmployee.contracts) : null;
  const endDatePreview = contractFromDate ? getDayBefore(contractFromDate) : null;

  const handleCreateEmployee = useCallback(() => {
    setEditingEmployee(null);
    resetEmployee({ first_name: '', last_name: '', gender: 'male', birthdate: '' });
    setIsEmployeeDialogOpen(true);
  }, [resetEmployee]);

  const handleEditEmployee = useCallback(
    (employee: Employee) => {
      setEditingEmployee(employee);
      resetEmployee({
        first_name: employee.first_name,
        last_name: employee.last_name,
        gender: employee.gender,
        birthdate: formatDateForInput(employee.birthdate),
      });
      setIsEmployeeDialogOpen(true);
    },
    [resetEmployee]
  );

  const handleDeleteEmployee = useCallback((employee: Employee) => {
    setDeletingEmployee(employee);
    setIsDeleteDialogOpen(true);
  }, []);

  const handleAddContract = useCallback(
    (employee: Employee) => {
      setContractEmployee(employee);
      setEndCurrentContract(true);

      const defaultPayPlanId = payPlans.length === 1 ? payPlans[0].id : 0;

      // Prefill from active contract if exists
      const active = getActiveContract(employee.contracts);
      if (active) {
        // Suggest start date as tomorrow
        const tomorrow = new Date();
        tomorrow.setDate(tomorrow.getDate() + 1);
        const tomorrowStr = tomorrow.toISOString().split('T')[0];

        resetContract({
          from: tomorrowStr,
          to: '',
          payplan_id: active.payplan_id || defaultPayPlanId,
          staff_category: active.staff_category as
            | 'qualified'
            | 'supplementary'
            | 'non_pedagogical',
          grade: active.grade,
          step: active.step,
          weekly_hours: active.weekly_hours,
        });
      } else {
        resetContract({
          from: '',
          to: '',
          payplan_id: defaultPayPlanId,
          staff_category: 'qualified',
          grade: '',
          step: 1,
          weekly_hours: 39,
        });
      }
      setIsContractDialogOpen(true);
    },
    [resetContract, payPlans]
  );

  const handleViewContractHistory = useCallback(
    (employee: Employee) => {
      router.push(`/organizations/${orgId}/employees/${employee.id}/contracts`);
    },
    [router, orgId]
  );

  const onSubmitEmployee = useCallback(
    (data: EmployeeFormData) => {
      if (editingEmployee) {
        updateMutation.mutate({ id: editingEmployee.id, data });
      } else {
        createMutation.mutate(data);
      }
    },
    [editingEmployee, updateMutation, createMutation]
  );

  // Helper to check if contract details have changed
  const contractDetailsChanged = (
    newData: {
      staff_category: string;
      grade: string;
      step: number;
      weekly_hours: number;
      payplan_id: number;
    },
    oldContract: EmployeeContract
  ): boolean => {
    return (
      newData.staff_category !== oldContract.staff_category ||
      newData.grade !== oldContract.grade ||
      newData.step !== oldContract.step ||
      newData.weekly_hours !== oldContract.weekly_hours ||
      newData.payplan_id !== oldContract.payplan_id
    );
  };

  const onSubmitContract = useCallback(
    (data: ContractFormData) => {
      if (contractEmployee) {
        // If there's an active contract and we're ending it, check if something actually changed
        if (activeContract && endCurrentContract) {
          const hasChanges = contractDetailsChanged(
            {
              staff_category: data.staff_category,
              grade: data.grade,
              step: data.step,
              weekly_hours: data.weekly_hours,
              payplan_id: data.payplan_id,
            },
            activeContract
          );
          if (!hasChanges) {
            toast({
              title: t('contracts.noChangesDetected'),
              description: t('contracts.noChangesDescription'),
              variant: 'destructive',
            });
            return;
          }
        }

        createContractMutation.mutate({
          employeeId: contractEmployee.id,
          data: {
            from: formatDateForApi(data.from) || data.from,
            to: formatDateForApi(data.to),
            staff_category: data.staff_category,
            grade: data.grade,
            step: data.step,
            weekly_hours: data.weekly_hours,
            payplan_id: data.payplan_id,
          },
        });
      }
    },
    [contractEmployee, activeContract, endCurrentContract, createContractMutation, toast, t]
  );

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">{t('employees.title')}</h1>
        </div>
        <Button onClick={handleCreateEmployee}>
          <Plus className="mr-2 h-4 w-4" />
          {t('employees.newEmployee')}
        </Button>
      </div>

      <div className="flex items-center gap-4">
        <div className="relative max-w-sm">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <label htmlFor="search-employees" className="sr-only">
            {t('common.search')}
          </label>
          <Input
            id="search-employees"
            placeholder={t('common.search')}
            value={searchInput}
            onChange={(e) => {
              setSearchInput(e.target.value);
              setPage(1);
            }}
            className="pl-9"
          />
        </div>
        <Select
          value={staffCategoryFilter}
          onValueChange={(value) => {
            setStaffCategoryFilter(value === 'all' ? '' : value);
            setPage(1);
          }}
        >
          <SelectTrigger className="w-[200px]">
            <SelectValue placeholder={t('employees.filterByCategory')} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">{t('common.all')}</SelectItem>
            <SelectItem value="qualified">{t('employees.staffCategory.qualified')}</SelectItem>
            <SelectItem value="supplementary">
              {t('employees.staffCategory.supplementary')}
            </SelectItem>
            <SelectItem value="non_pedagogical">
              {t('employees.staffCategory.non_pedagogical')}
            </SelectItem>
          </SelectContent>
        </Select>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{t('employees.title')}</CardTitle>
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
                  <TableHead>{t('common.name')}</TableHead>
                  <TableHead>{t('gender.label')}</TableHead>
                  <TableHead>{t('employees.birthdate')}</TableHead>
                  <TableHead>{t('employees.age')}</TableHead>
                  <TableHead>{t('employees.staffCategory.label')}</TableHead>
                  <TableHead>{t('employees.grade')}</TableHead>
                  <TableHead>{t('employees.weeklyHours')}</TableHead>
                  <TableHead>{t('employees.salary')}</TableHead>
                  <TableHead>{t('employees.yearsOfService')}</TableHead>
                  <TableHead className="text-right">{t('common.actions')}</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {employees?.map((employee) => {
                  const currentContract = getCurrentContract(employee.contracts);
                  const payPlanForSalary = currentContract?.payplan_id
                    ? payPlanMap.get(currentContract.payplan_id)
                    : undefined;
                  const salary =
                    currentContract && payPlanForSalary
                      ? calculateMonthlySalary(currentContract, payPlanForSalary)
                      : null;
                  const yearsOfService = employee.contracts?.length
                    ? calculateYearsOfService(employee.contracts)
                    : null;
                  return (
                    <TableRow key={employee.id}>
                      <TableCell className="font-medium">
                        {employee.first_name} {employee.last_name}
                      </TableCell>
                      <TableCell>{t(`gender.${employee.gender}`)}</TableCell>
                      <TableCell>{formatDate(employee.birthdate)}</TableCell>
                      <TableCell>{calculateAge(employee.birthdate)}</TableCell>
                      <TableCell>
                        {currentContract ? (
                          t(`employees.staffCategory.${currentContract.staff_category}`)
                        ) : (
                          <span className="text-muted-foreground">{t('employees.noContract')}</span>
                        )}
                      </TableCell>
                      <TableCell>
                        {currentContract
                          ? `${currentContract.grade} / ${currentContract.step}`
                          : '-'}
                      </TableCell>
                      <TableCell>{currentContract?.weekly_hours || '-'}</TableCell>
                      <TableCell>{salary !== null ? formatCurrency(salary) : '-'}</TableCell>
                      <TableCell>
                        {yearsOfService !== null ? yearsOfService.toFixed(1) : '-'}
                      </TableCell>
                      <TableCell className="text-right">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleViewContractHistory(employee)}
                          title={t('employees.contractHistory')}
                          aria-label={t('employees.contractHistory')}
                        >
                          <History className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleAddContract(employee)}
                          title={t('employees.addContract')}
                          aria-label={t('employees.addContract')}
                        >
                          <FileText className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleEditEmployee(employee)}
                          aria-label={t('common.edit')}
                        >
                          <Pencil className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleDeleteEmployee(employee)}
                          aria-label={t('common.delete')}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </TableCell>
                    </TableRow>
                  );
                })}
                {employees?.length === 0 && (
                  <TableRow>
                    <TableCell colSpan={10} className="text-center text-muted-foreground">
                      {t('common.noResults')}
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          )}
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

      {/* Employee Create/Edit Dialog */}
      <PersonFormDialog
        open={isEmployeeDialogOpen}
        onOpenChange={setIsEmployeeDialogOpen}
        isEditing={!!editingEmployee}
        register={registerEmployee}
        onSubmit={handleSubmitEmployee(onSubmitEmployee)}
        errors={errorsEmployee}
        watch={watchEmployee}
        setValue={setValueEmployee}
        isSaving={createMutation.isPending || updateMutation.isPending}
        translationPrefix="employees"
      />

      {/* Contract Create Dialog */}
      <Dialog open={isContractDialogOpen} onOpenChange={setIsContractDialogOpen}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>
              {t('contracts.newContractFor', {
                name: contractEmployee
                  ? `${contractEmployee.first_name} ${contractEmployee.last_name}`
                  : '',
              })}
            </DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSubmitContract(onSubmitContract)} className="space-y-4">
            {/* Show active contract info if exists */}
            {activeContract && (
              <Alert>
                <AlertDescription className="space-y-3">
                  <p className="font-medium">{t('contracts.hasActiveContractEmployee')}</p>
                  <p className="text-sm text-muted-foreground">
                    {t('contracts.activeSinceEmployee', {
                      date: formatDate(activeContract.from),
                      staffCategory: t(`employees.staffCategory.${activeContract.staff_category}`),
                      grade: activeContract.grade,
                      step: activeContract.step,
                    })}
                  </p>
                  <div className="flex items-center space-x-2">
                    <Checkbox
                      id="endCurrentContract"
                      checked={endCurrentContract}
                      onCheckedChange={(checked) => setEndCurrentContract(checked === true)}
                    />
                    <label
                      htmlFor="endCurrentContract"
                      className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                    >
                      {endDatePreview
                        ? t('contracts.endCurrentContract', { date: formatDate(endDatePreview) })
                        : t('contracts.endCurrentContract', { date: '...' })}
                    </label>
                  </div>
                </AlertDescription>
              </Alert>
            )}

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="from">{t('contracts.startDate')}</Label>
                <Input id="from" type="date" {...registerContract('from')} />
                {errorsContract.from && (
                  <p className="text-sm text-destructive">{t('contracts.startDateRequired')}</p>
                )}
              </div>
              <div className="space-y-2">
                <Label htmlFor="to">{t('contracts.endDateOptional')}</Label>
                <Input id="to" type="date" {...registerContract('to')} />
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="payplan_id">{t('employees.payPlan')}</Label>
              <Select
                value={String(watchContract('payplan_id') || '')}
                onValueChange={(val) => setValueContract('payplan_id', Number(val))}
              >
                <SelectTrigger>
                  <SelectValue placeholder={t('employees.selectPayPlan')} />
                </SelectTrigger>
                <SelectContent>
                  {payPlans.map((pp) => (
                    <SelectItem key={pp.id} value={String(pp.id)}>
                      {pp.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              {errorsContract.payplan_id && (
                <p className="text-sm text-destructive">{t('employees.selectPayPlan')}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="staff_category">{t('employees.staffCategory.label')}</Label>
              <select
                id="staff_category"
                {...registerContract('staff_category')}
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
              >
                <option value="qualified">{t('employees.staffCategory.qualified')}</option>
                <option value="supplementary">{t('employees.staffCategory.supplementary')}</option>
                <option value="non_pedagogical">
                  {t('employees.staffCategory.non_pedagogical')}
                </option>
              </select>
              {errorsContract.staff_category && (
                <p className="text-sm text-destructive">{t('validation.staffCategoryRequired')}</p>
              )}
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="grade">{t('employees.grade')}</Label>
                <Input id="grade" {...registerContract('grade')} placeholder="S8a" />
                {errorsContract.grade && (
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
                  {...registerContract('step', { valueAsNumber: true })}
                />
                {errorsContract.step && (
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
                {...registerContract('weekly_hours', { valueAsNumber: true })}
              />
              {errorsContract.weekly_hours && (
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
              <Button type="submit" disabled={createContractMutation.isPending}>
                {t('common.save')}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <DeleteConfirmDialog
        open={isDeleteDialogOpen}
        onOpenChange={setIsDeleteDialogOpen}
        onConfirm={() => deletingEmployee && deleteMutation.mutate(deletingEmployee.id)}
        isLoading={deleteMutation.isPending}
        resourceName="employees"
        description={t('employees.confirmDeleteMessage', {
          name: deletingEmployee
            ? `${deletingEmployee.first_name} ${deletingEmployee.last_name}`
            : '',
        })}
      />
    </div>
  );
}
