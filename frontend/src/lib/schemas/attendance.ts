import { z } from 'zod';

export const attendanceSchema = z.object({
  status: z.enum(['present', 'absent', 'sick', 'vacation']),
  check_in_time: z.string().optional(),
  check_out_time: z.string().optional(),
  note: z.string().max(500).optional(),
});

export type AttendanceFormData = z.infer<typeof attendanceSchema>;
