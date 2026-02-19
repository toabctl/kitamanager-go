'use client';

import { useTranslations } from 'next-intl';
import { CheckCircle, XCircle, Thermometer, Palmtree, Pencil, Trash2 } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';
import type { ChildAttendanceResponse, ChildAttendanceStatus } from '@/lib/api/types';
import { formatTime } from '@/lib/utils/formatting';

export interface AttendanceRow {
  childId: number;
  childName: string;
  attendance: ChildAttendanceResponse | null;
}

interface AttendanceTableProps {
  rows: AttendanceRow[];
  onQuickStatus: (childId: number, status: ChildAttendanceStatus, attendanceId?: number) => void;
  onEdit: (row: AttendanceRow) => void;
  onDelete: (row: AttendanceRow) => void;
}

const STATUS_BADGE_VARIANT: Record<
  ChildAttendanceStatus,
  'success' | 'destructive' | 'warning' | 'default'
> = {
  present: 'success',
  absent: 'destructive',
  sick: 'warning',
  vacation: 'default',
};

const QUICK_BUTTONS: { status: ChildAttendanceStatus; icon: typeof CheckCircle; color: string }[] =
  [
    { status: 'present', icon: CheckCircle, color: 'text-green-600 hover:text-green-700' },
    { status: 'absent', icon: XCircle, color: 'text-red-600 hover:text-red-700' },
    { status: 'sick', icon: Thermometer, color: 'text-orange-600 hover:text-orange-700' },
    { status: 'vacation', icon: Palmtree, color: 'text-blue-600 hover:text-blue-700' },
  ];

export function AttendanceTable({ rows, onQuickStatus, onEdit, onDelete }: AttendanceTableProps) {
  const t = useTranslations('attendance');
  const tCommon = useTranslations('common');

  if (rows.length === 0) {
    return <p className="text-muted-foreground py-8 text-center">{t('noChildren')}</p>;
  }

  return (
    <TooltipProvider>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>{tCommon('name')}</TableHead>
            <TableHead>{tCommon('status')}</TableHead>
            <TableHead>{t('checkIn')}</TableHead>
            <TableHead>{t('checkOut')}</TableHead>
            <TableHead>{t('note')}</TableHead>
            <TableHead className="text-center">{t('quickMark')}</TableHead>
            <TableHead className="text-right">{tCommon('actions')}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {rows.map((row) => (
            <TableRow key={row.childId}>
              <TableCell className="font-medium">{row.childName}</TableCell>
              <TableCell>
                {row.attendance ? (
                  <Badge variant={STATUS_BADGE_VARIANT[row.attendance.status]}>
                    {t(row.attendance.status)}
                  </Badge>
                ) : (
                  <span className="text-muted-foreground text-sm">{t('notRecorded')}</span>
                )}
              </TableCell>
              <TableCell>{formatTime(row.attendance?.check_in_time)}</TableCell>
              <TableCell>{formatTime(row.attendance?.check_out_time)}</TableCell>
              <TableCell className="max-w-[200px] truncate">{row.attendance?.note || ''}</TableCell>
              <TableCell>
                <div className="flex items-center justify-center gap-1">
                  {QUICK_BUTTONS.map(({ status, icon: Icon, color }) => (
                    <Tooltip key={status}>
                      <TooltipTrigger asChild>
                        <Button
                          variant="ghost"
                          size="icon"
                          className={`h-7 w-7 ${color} ${row.attendance?.status === status ? 'bg-muted' : ''}`}
                          onClick={() => onQuickStatus(row.childId, status, row.attendance?.id)}
                        >
                          <Icon className="h-4 w-4" />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>{t(status)}</TooltipContent>
                    </Tooltip>
                  ))}
                </div>
              </TableCell>
              <TableCell className="text-right">
                {row.attendance && (
                  <div className="flex items-center justify-end gap-1">
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-7 w-7"
                      onClick={() => onEdit(row)}
                      aria-label={tCommon('edit')}
                    >
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="text-destructive h-7 w-7"
                      onClick={() => onDelete(row)}
                      aria-label={tCommon('delete')}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                )}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TooltipProvider>
  );
}
