import {
  governmentFundingSchema,
  governmentFundingPeriodSchema,
  governmentFundingPropertySchema,
} from '../government-funding';

describe('governmentFundingSchema', () => {
  it('accepts valid funding', () => {
    expect(
      governmentFundingSchema.safeParse({ name: 'Berlin Funding', state: 'berlin' }).success
    ).toBe(true);
  });

  it('rejects empty name', () => {
    expect(governmentFundingSchema.safeParse({ name: '', state: 'berlin' }).success).toBe(false);
  });

  it('rejects empty state', () => {
    expect(governmentFundingSchema.safeParse({ name: 'Test', state: '' }).success).toBe(false);
  });
});

describe('governmentFundingPeriodSchema', () => {
  const valid = {
    from: '2024-01-01',
    full_time_weekly_hours: 39,
  };

  it('accepts valid period without end date', () => {
    expect(governmentFundingPeriodSchema.safeParse(valid).success).toBe(true);
  });

  it('accepts valid period with end date after start', () => {
    expect(governmentFundingPeriodSchema.safeParse({ ...valid, to: '2024-12-31' }).success).toBe(
      true
    );
  });

  it('rejects end date before start date', () => {
    const result = governmentFundingPeriodSchema.safeParse({ ...valid, to: '2023-06-01' });
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].path).toContain('to');
    }
  });

  it('rejects weekly hours below 0.1', () => {
    expect(
      governmentFundingPeriodSchema.safeParse({ ...valid, full_time_weekly_hours: 0 }).success
    ).toBe(false);
  });

  it('rejects weekly hours above 80', () => {
    expect(
      governmentFundingPeriodSchema.safeParse({ ...valid, full_time_weekly_hours: 81 }).success
    ).toBe(false);
  });

  it('accepts optional comment', () => {
    const result = governmentFundingPeriodSchema.safeParse({
      ...valid,
      comment: 'Test comment',
    });
    expect(result.success).toBe(true);
  });
});

describe('governmentFundingPropertySchema', () => {
  const valid = {
    key: 'care_type',
    value: 'ganztag',
    label: 'Ganztag',
    payment_euros: 1668.47,
    requirement: 0.261,
  };

  it('accepts valid property', () => {
    expect(governmentFundingPropertySchema.safeParse(valid).success).toBe(true);
  });

  it('accepts property with age range', () => {
    expect(
      governmentFundingPropertySchema.safeParse({ ...valid, min_age: 0, max_age: 2 }).success
    ).toBe(true);
  });

  it('accepts null age fields', () => {
    expect(
      governmentFundingPropertySchema.safeParse({ ...valid, min_age: null, max_age: null }).success
    ).toBe(true);
  });

  it('rejects negative payment', () => {
    expect(governmentFundingPropertySchema.safeParse({ ...valid, payment_euros: -1 }).success).toBe(
      false
    );
  });

  it('rejects negative requirement', () => {
    expect(governmentFundingPropertySchema.safeParse({ ...valid, requirement: -0.1 }).success).toBe(
      false
    );
  });

  it('rejects empty key', () => {
    expect(governmentFundingPropertySchema.safeParse({ ...valid, key: '' }).success).toBe(false);
  });

  it('rejects empty value', () => {
    expect(governmentFundingPropertySchema.safeParse({ ...valid, value: '' }).success).toBe(false);
  });
});
