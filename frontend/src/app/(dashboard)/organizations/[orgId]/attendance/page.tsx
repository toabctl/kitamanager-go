'use client';

import { useState, useMemo, useCallback } from 'react';
import { useParams } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { DayStepper } from '@/components/ui/day-stepper';
import { AttendanceSummary } from '@/components/attendance/attendance-summary';
import { AttendanceTable, type AttendanceRow } from '@/components/attendance/attendance-table';
import { AttendanceEditDialog } from '@/components/attendance/attendance-edit-dialog';
import { DeleteConfirmDialog } from '@/components/crud/delete-confirm-dialog';
import { QueryError } from '@/components/crud/query-error';
import { apiClient } from '@/lib/api/client';
import { queryKeys } from '@/lib/api/queryKeys';
import type { ChildAttendanceStatus } from '@/lib/api/types';
import type { AttendanceFormData } from '@/lib/schemas/attendance';
import { combineDateAndTime } from '@/lib/utils/formatting';
import { useToast } from '@/lib/hooks/use-toast';

export default function AttendancePage() {
  const params = useParams();
  const orgId = Number(params.orgId);
  const t = useTranslations('attendance');
  const { toast } = useToast();
  const queryClient = useQueryClient();

  const [selectedDate, setSelectedDate] = useState(() => new Date());
  const dateStr = selectedDate.toISOString().slice(0, 10);

  // Edit dialog state
  const [editRow, setEditRow] = useState<AttendanceRow | null>(null);
  const [isEditOpen, setIsEditOpen] = useState(false);

  // Delete dialog state
  const [deleteRow, setDeleteRow] = useState<AttendanceRow | null>(null);
  const [isDeleteOpen, setIsDeleteOpen] = useState(false);

  // Fetch active children for the date
  const {
    data: children,
    isLoading: childrenLoading,
    error: childrenError,
    refetch: refetchChildren,
  } = useQuery({
    queryKey: [...queryKeys.children.allUnpaginated(orgId), dateStr],
    queryFn: () => apiClient.getChildrenAllForDate(orgId, dateStr),
    enabled: !!orgId,
  });

  // Fetch attendance records for the date
  const {
    data: attendanceRecords,
    isLoading: attendanceLoading,
    error: attendanceError,
    refetch: refetchAttendance,
  } = useQuery({
    queryKey: queryKeys.attendance.byDate(orgId, dateStr),
    queryFn: () => apiClient.getChildAttendanceByDateAll(orgId, dateStr),
    enabled: !!orgId,
  });

  const isLoading = childrenLoading || attendanceLoading;
  const queryError = childrenError || attendanceError;

  // Merge children + attendance records into rows
  const rows: AttendanceRow[] = useMemo(() => {
    if (!children) return [];
    const attendanceByChildId = new Map((attendanceRecords ?? []).map((a) => [a.child_id, a]));
    const merged = children.map((child) => ({
      childId: child.id,
      childName: `${child.first_name} ${child.last_name}`,
      attendance: attendanceByChildId.get(child.id) ?? null,
    }));
    // Sort: recorded children first, then unrecorded; alphabetical within each group
    merged.sort((a, b) => {
      const aRecorded = a.attendance ? 0 : 1;
      const bRecorded = b.attendance ? 0 : 1;
      if (aRecorded !== bRecorded) return aRecorded - bRecorded;
      return a.childName.localeCompare(b.childName);
    });
    return merged;
  }, [children, attendanceRecords]);

  const invalidateAll = useCallback(() => {
    queryClient.invalidateQueries({ queryKey: queryKeys.attendance.byDate(orgId, dateStr) });
    queryClient.invalidateQueries({ queryKey: queryKeys.attendance.summary(orgId, dateStr) });
  }, [queryClient, orgId, dateStr]);

  // Quick status mutation: create or update
  const quickStatusMutation = useMutation({
    mutationFn: async ({
      childId,
      status,
      attendanceId,
    }: {
      childId: number;
      status: ChildAttendanceStatus;
      attendanceId?: number;
    }) => {
      if (attendanceId) {
        return apiClient.updateChildAttendance(orgId, childId, attendanceId, { status });
      }
      return apiClient.createChildAttendance(orgId, childId, { date: dateStr, status });
    },
    onSuccess: () => {
      invalidateAll();
    },
    onError: () => {
      toast({ title: t('failedToSave'), variant: 'destructive' });
    },
  });

  // Update attendance mutation (from edit dialog)
  const updateMutation = useMutation({
    mutationFn: async ({ data }: { data: AttendanceFormData }) => {
      if (!editRow?.attendance) return;
      return apiClient.updateChildAttendance(orgId, editRow.childId, editRow.attendance.id, {
        status: data.status,
        check_in_time: combineDateAndTime(dateStr, data.check_in_time ?? ''),
        check_out_time: combineDateAndTime(dateStr, data.check_out_time ?? ''),
        note: data.note || undefined,
      });
    },
    onSuccess: () => {
      invalidateAll();
      setIsEditOpen(false);
      setEditRow(null);
      toast({ title: t('updateSuccess') });
    },
    onError: () => {
      toast({ title: t('failedToSave'), variant: 'destructive' });
    },
  });

  // Delete attendance mutation
  const deleteMutation = useMutation({
    mutationFn: async () => {
      if (!deleteRow?.attendance) return;
      return apiClient.deleteChildAttendance(orgId, deleteRow.childId, deleteRow.attendance.id);
    },
    onSuccess: () => {
      invalidateAll();
      setIsDeleteOpen(false);
      setDeleteRow(null);
      toast({ title: t('deleteSuccess') });
    },
    onError: () => {
      toast({ title: t('failedToDelete'), variant: 'destructive' });
    },
  });

  const handleQuickStatus = useCallback(
    (childId: number, status: ChildAttendanceStatus, attendanceId?: number) => {
      quickStatusMutation.mutate({ childId, status, attendanceId });
    },
    [quickStatusMutation]
  );

  const handleEdit = useCallback((row: AttendanceRow) => {
    setEditRow(row);
    setIsEditOpen(true);
  }, []);

  const handleDelete = useCallback((row: AttendanceRow) => {
    setDeleteRow(row);
    setIsDeleteOpen(true);
  }, []);

  const handleEditSubmit = useCallback(
    (data: AttendanceFormData) => {
      updateMutation.mutate({ data });
    },
    [updateMutation]
  );

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold tracking-tight">{t('title')}</h1>
      </div>

      <DayStepper value={selectedDate} onChange={setSelectedDate} />

      <AttendanceSummary orgId={orgId} date={dateStr} />

      <Card>
        <CardHeader>
          <CardTitle>{t('title')}</CardTitle>
        </CardHeader>
        <CardContent>
          {queryError ? (
            <QueryError
              error={queryError}
              onRetry={() => {
                refetchChildren();
                refetchAttendance();
              }}
            />
          ) : isLoading ? (
            <div className="space-y-2">
              {[...Array(5)].map((_, i) => (
                <Skeleton key={i} className="h-12 w-full" />
              ))}
            </div>
          ) : (
            <AttendanceTable
              rows={rows}
              onQuickStatus={handleQuickStatus}
              onEdit={handleEdit}
              onDelete={handleDelete}
            />
          )}
        </CardContent>
      </Card>

      <AttendanceEditDialog
        open={isEditOpen}
        onOpenChange={setIsEditOpen}
        attendance={editRow?.attendance ?? null}
        childName={editRow?.childName ?? ''}
        isSaving={updateMutation.isPending}
        onSubmit={handleEditSubmit}
      />

      <DeleteConfirmDialog
        open={isDeleteOpen}
        onOpenChange={setIsDeleteOpen}
        onConfirm={() => deleteMutation.mutate()}
        isLoading={deleteMutation.isPending}
        resourceName="attendance"
      />
    </div>
  );
}
