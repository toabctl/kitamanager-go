import { attendanceSchema } from '../attendance';

describe('attendanceSchema', () => {
  it('accepts valid present status', () => {
    expect(attendanceSchema.safeParse({ status: 'present' }).success).toBe(true);
  });

  it('accepts all valid statuses', () => {
    for (const status of ['present', 'absent', 'sick', 'vacation']) {
      expect(attendanceSchema.safeParse({ status }).success).toBe(true);
    }
  });

  it('rejects invalid status', () => {
    expect(attendanceSchema.safeParse({ status: 'late' }).success).toBe(false);
  });

  it('rejects empty status', () => {
    expect(attendanceSchema.safeParse({ status: '' }).success).toBe(false);
  });

  it('accepts optional check_in_time', () => {
    expect(attendanceSchema.safeParse({ status: 'present', check_in_time: '08:00' }).success).toBe(
      true
    );
  });

  it('accepts optional note', () => {
    expect(attendanceSchema.safeParse({ status: 'absent', note: 'Doctor visit' }).success).toBe(
      true
    );
  });

  it('rejects note over 500 chars', () => {
    expect(attendanceSchema.safeParse({ status: 'present', note: 'A'.repeat(501) }).success).toBe(
      false
    );
  });
});
