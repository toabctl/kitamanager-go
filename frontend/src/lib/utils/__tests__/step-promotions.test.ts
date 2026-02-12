import { calculateYearsOfService, determineEligibleStep } from '../step-promotions';

describe('calculateYearsOfService', () => {
  beforeEach(() => {
    jest.useFakeTimers();
    jest.setSystemTime(new Date('2025-06-15'));
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it('returns 0 for empty contracts', () => {
    expect(calculateYearsOfService([])).toBe(0);
  });

  it('calculates correctly for single contract started 3 years ago', () => {
    const contracts = [{ from: '2022-06-15' }];
    const result = calculateYearsOfService(contracts);
    expect(result).toBeCloseTo(3, 0);
  });

  it('uses earliest contract from_date with multiple contracts', () => {
    const contracts = [
      { from: '2023-01-01' },
      { from: '2020-06-15' }, // earliest
      { from: '2024-01-01' },
    ];
    const result = calculateYearsOfService(contracts);
    expect(result).toBeCloseTo(5, 0);
  });

  it('returns 0 for contract started today', () => {
    const contracts = [{ from: '2025-06-15' }];
    expect(calculateYearsOfService(contracts)).toBe(0);
  });

  it('returns 0 for future contract', () => {
    const contracts = [{ from: '2026-01-01' }];
    expect(calculateYearsOfService(contracts)).toBe(0);
  });

  it('uses custom asOf date', () => {
    const contracts = [{ from: '2020-01-01' }];
    const asOf = new Date('2023-01-01');
    const result = calculateYearsOfService(contracts, asOf);
    expect(result).toBeCloseTo(3, 0);
  });
});

describe('determineEligibleStep', () => {
  const entries = [
    { step: 1, grade: 'S8a', step_min_years: 0 },
    { step: 2, grade: 'S8a', step_min_years: 1 },
    { step: 3, grade: 'S8a', step_min_years: 3 },
    { step: 4, grade: 'S8a', step_min_years: 6 },
    { step: 5, grade: 'S8a', step_min_years: 10 },
    { step: 6, grade: 'S8a', step_min_years: 15 },
  ];

  it('returns 0 when no entries have step_min_years', () => {
    const entriesNoRules = [
      { step: 1, grade: 'S8a', step_min_years: null },
      { step: 2, grade: 'S8a' },
    ];
    expect(determineEligibleStep(5, entriesNoRules, 'S8a')).toBe(0);
  });

  it('returns step 1 for 0 years of service', () => {
    expect(determineEligibleStep(0, entries, 'S8a')).toBe(1);
  });

  it('returns correct step between thresholds', () => {
    expect(determineEligibleStep(2, entries, 'S8a')).toBe(2); // 1 <= 2 < 3
    expect(determineEligibleStep(5, entries, 'S8a')).toBe(3); // 3 <= 5 < 6
    expect(determineEligibleStep(8, entries, 'S8a')).toBe(4); // 6 <= 8 < 10
  });

  it('returns step at exact threshold', () => {
    expect(determineEligibleStep(3, entries, 'S8a')).toBe(3);
    expect(determineEligibleStep(10, entries, 'S8a')).toBe(5);
  });

  it('returns max step when beyond all thresholds', () => {
    expect(determineEligibleStep(20, entries, 'S8a')).toBe(6);
  });

  it('returns 0 for wrong grade', () => {
    expect(determineEligibleStep(5, entries, 'S9')).toBe(0);
  });

  it('handles mix of entries with and without step_min_years', () => {
    const mixed = [
      { step: 1, grade: 'S8a', step_min_years: 0 },
      { step: 2, grade: 'S8a', step_min_years: null },
      { step: 3, grade: 'S8a', step_min_years: 3 },
    ];
    expect(determineEligibleStep(5, mixed, 'S8a')).toBe(3);
  });
});
