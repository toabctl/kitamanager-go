import { payPlanSchema, payPlanPeriodSchema, payPlanEntrySchema } from '../payplan';

describe('payPlanSchema', () => {
  it('accepts valid name', () => {
    expect(payPlanSchema.safeParse({ name: 'TVöD-SuE' }).success).toBe(true);
  });

  it('rejects empty name', () => {
    expect(payPlanSchema.safeParse({ name: '' }).success).toBe(false);
  });

  it('rejects name over 255 chars', () => {
    expect(payPlanSchema.safeParse({ name: 'A'.repeat(256) }).success).toBe(false);
  });
});

describe('payPlanPeriodSchema', () => {
  const valid = {
    from: '2024-01-01',
    weekly_hours: 39,
    employer_contribution_rate: 22,
  };

  it('accepts valid period', () => {
    expect(payPlanPeriodSchema.safeParse(valid).success).toBe(true);
  });

  it('accepts period with end date after start', () => {
    expect(payPlanPeriodSchema.safeParse({ ...valid, to: '2024-12-31' }).success).toBe(true);
  });

  it('rejects end date before start date', () => {
    const result = payPlanPeriodSchema.safeParse({ ...valid, to: '2023-06-01' });
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].path).toContain('to');
    }
  });

  it('rejects zero weekly hours', () => {
    expect(payPlanPeriodSchema.safeParse({ ...valid, weekly_hours: 0 }).success).toBe(false);
  });

  it('rejects weekly hours above 168', () => {
    expect(payPlanPeriodSchema.safeParse({ ...valid, weekly_hours: 169 }).success).toBe(false);
  });

  it('rejects negative contribution rate', () => {
    expect(
      payPlanPeriodSchema.safeParse({ ...valid, employer_contribution_rate: -1 }).success
    ).toBe(false);
  });

  it('rejects contribution rate above 100', () => {
    expect(
      payPlanPeriodSchema.safeParse({ ...valid, employer_contribution_rate: 101 }).success
    ).toBe(false);
  });

  it('accepts 0% contribution rate', () => {
    expect(payPlanPeriodSchema.safeParse({ ...valid, employer_contribution_rate: 0 }).success).toBe(
      true
    );
  });

  it('accepts 100% contribution rate', () => {
    expect(
      payPlanPeriodSchema.safeParse({ ...valid, employer_contribution_rate: 100 }).success
    ).toBe(true);
  });
});

describe('payPlanEntrySchema', () => {
  const valid = {
    grade: 'S8a',
    step: 3,
    monthly_amount_euros: 3500.89,
  };

  it('accepts valid entry', () => {
    expect(payPlanEntrySchema.safeParse(valid).success).toBe(true);
  });

  it('accepts entry with step_min_years', () => {
    expect(payPlanEntrySchema.safeParse({ ...valid, step_min_years: 2 }).success).toBe(true);
  });

  it('rejects step below 1', () => {
    expect(payPlanEntrySchema.safeParse({ ...valid, step: 0 }).success).toBe(false);
  });

  it('rejects step above 10', () => {
    expect(payPlanEntrySchema.safeParse({ ...valid, step: 11 }).success).toBe(false);
  });

  it('rejects empty grade', () => {
    expect(payPlanEntrySchema.safeParse({ ...valid, grade: '' }).success).toBe(false);
  });

  it('rejects negative monthly amount', () => {
    expect(payPlanEntrySchema.safeParse({ ...valid, monthly_amount_euros: -1 }).success).toBe(
      false
    );
  });
});
