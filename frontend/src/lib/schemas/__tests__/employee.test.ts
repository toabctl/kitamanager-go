import { employeeSchema, employeeContractSchema } from '../employee';

describe('employeeSchema', () => {
  const valid = {
    first_name: 'Anna',
    last_name: 'Müller',
    gender: 'female' as const,
    birthdate: '1985-03-10',
  };

  it('accepts valid employee data', () => {
    expect(employeeSchema.safeParse(valid).success).toBe(true);
  });

  it('rejects empty first_name', () => {
    expect(employeeSchema.safeParse({ ...valid, first_name: '' }).success).toBe(false);
  });

  it('rejects invalid gender', () => {
    expect(employeeSchema.safeParse({ ...valid, gender: 'other' }).success).toBe(false);
  });

  it('accepts all valid genders', () => {
    for (const gender of ['male', 'female', 'diverse']) {
      expect(employeeSchema.safeParse({ ...valid, gender }).success).toBe(true);
    }
  });
});

describe('employeeContractSchema', () => {
  const valid = {
    from: '2024-01-01',
    section_id: 1,
    payplan_id: 1,
    staff_category: 'qualified' as const,
    grade: 'S8a',
    step: 3,
    weekly_hours: 39,
  };

  it('accepts valid contract data', () => {
    expect(employeeContractSchema.safeParse(valid).success).toBe(true);
  });

  it('accepts contract with end date after start date', () => {
    expect(employeeContractSchema.safeParse({ ...valid, to: '2024-12-31' }).success).toBe(true);
  });

  it('rejects end date before start date', () => {
    const result = employeeContractSchema.safeParse({ ...valid, to: '2023-06-01' });
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].path).toContain('to');
    }
  });

  it('accepts contract without end date', () => {
    expect(employeeContractSchema.safeParse(valid).success).toBe(true);
  });

  it('rejects section_id of 0', () => {
    expect(employeeContractSchema.safeParse({ ...valid, section_id: 0 }).success).toBe(false);
  });

  it('accepts all valid staff categories', () => {
    for (const cat of ['qualified', 'supplementary', 'non_pedagogical']) {
      expect(employeeContractSchema.safeParse({ ...valid, staff_category: cat }).success).toBe(
        true
      );
    }
  });

  it('rejects invalid staff category', () => {
    expect(employeeContractSchema.safeParse({ ...valid, staff_category: 'invalid' }).success).toBe(
      false
    );
  });

  it('rejects step > 10', () => {
    expect(employeeContractSchema.safeParse({ ...valid, step: 11 }).success).toBe(false);
  });

  it('rejects weekly_hours > 168', () => {
    expect(employeeContractSchema.safeParse({ ...valid, weekly_hours: 169 }).success).toBe(false);
  });

  it('rejects negative weekly_hours', () => {
    expect(employeeContractSchema.safeParse({ ...valid, weekly_hours: -1 }).success).toBe(false);
  });
});
