'use client';

import { format } from 'date-fns';
import { de, enUS } from 'date-fns/locale';
import { useLocale, useTranslations } from 'next-intl';
import { LogIn, LogOut } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import type { ChildAttendanceResponse } from '@/lib/api/types';
import type { Child } from '@/lib/api/types';
import { formatTime } from '@/lib/utils/formatting';

const dateFnsLocales: Record<string, typeof de> = {
  de: de,
  en: enUS,
};

interface AttendanceCellProps {
  attendance: ChildAttendanceResponse | undefined;
  childId: number;
  dateStr: string;
  onCheckIn: (childId: number, dateStr: string) => void;
  onCheckOut: (childId: number, dateStr: string, attendanceId: number) => void;
}

function AttendanceCell({
  attendance,
  childId,
  dateStr,
  onCheckIn,
  onCheckOut,
}: AttendanceCellProps) {
  const t = useTranslations('attendance');

  // No record — show check-in button
  if (!attendance) {
    return (
      <Button
        variant="outline"
        size="sm"
        className="h-8 gap-1 text-green-600 hover:bg-green-50 hover:text-green-700"
        onClick={() => onCheckIn(childId, dateStr)}
        aria-label={t('checkIn')}
      >
        <LogIn className="h-4 w-4" />
        <span className="hidden sm:inline">{t('checkIn')}</span>
      </Button>
    );
  }

  const checkIn = formatTime(attendance.check_in_time);
  const checkOut = formatTime(attendance.check_out_time);

  // Checked in but not checked out — show time + check-out button
  if (checkIn && !checkOut) {
    return (
      <div className="flex items-center gap-1">
        <span className="text-sm font-medium text-green-700">{checkIn}</span>
        <Button
          variant="outline"
          size="sm"
          className="h-8 gap-1 text-orange-600 hover:bg-orange-50 hover:text-orange-700"
          onClick={() => onCheckOut(childId, dateStr, attendance.id)}
          aria-label={t('checkOut')}
        >
          <LogOut className="h-4 w-4" />
          <span className="hidden sm:inline">{t('checkOut')}</span>
        </Button>
      </div>
    );
  }

  // Checked out — show both times
  if (checkIn && checkOut) {
    return (
      <span className="text-muted-foreground text-sm">
        {checkIn} – {checkOut}
      </span>
    );
  }

  // Has record but no times (status-only, e.g. marked absent/sick)
  return <span className="text-muted-foreground text-sm italic">{t(attendance.status)}</span>;
}

interface AttendanceWeekTableProps {
  childRecords: Child[];
  attendanceByDate: Map<string, ChildAttendanceResponse[]>;
  onCheckIn: (childId: number, dateStr: string) => void;
  onCheckOut: (childId: number, dateStr: string, attendanceId: number) => void;
  days: Date[];
}

export function AttendanceWeekTable({
  childRecords,
  attendanceByDate,
  onCheckIn,
  onCheckOut,
  days,
}: AttendanceWeekTableProps) {
  const t = useTranslations('attendance');
  const tCommon = useTranslations('common');
  const locale = useLocale();
  const dfLocale = dateFnsLocales[locale] ?? enUS;

  const sortedChildren = [...childRecords].sort((a, b) => {
    const nameA = `${a.first_name} ${a.last_name}`;
    const nameB = `${b.first_name} ${b.last_name}`;
    return nameA.localeCompare(nameB);
  });

  if (sortedChildren.length === 0) {
    return <p className="text-muted-foreground py-8 text-center">{t('noChildren')}</p>;
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>{tCommon('name')}</TableHead>
          {days.map((day) => (
            <TableHead key={day.toISOString()} className="text-center">
              {format(day, 'EEE dd.MM', { locale: dfLocale })}
            </TableHead>
          ))}
        </TableRow>
      </TableHeader>
      <TableBody>
        {sortedChildren.map((child) => (
          <TableRow key={child.id}>
            <TableCell className="font-medium">
              {child.first_name} {child.last_name}
            </TableCell>
            {days.map((day) => {
              const dayStr = day.toISOString().slice(0, 10);
              const dayRecords = attendanceByDate.get(dayStr) ?? [];
              const attendance = dayRecords.find((a) => a.child_id === child.id);
              return (
                <TableCell key={dayStr} className="text-center">
                  <div className="flex justify-center">
                    <AttendanceCell
                      attendance={attendance}
                      childId={child.id}
                      dateStr={dayStr}
                      onCheckIn={onCheckIn}
                      onCheckOut={onCheckOut}
                    />
                  </div>
                </TableCell>
              );
            })}
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
