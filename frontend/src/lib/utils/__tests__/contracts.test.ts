import {
  getActiveContract,
  getCurrentContract,
  getDayBefore,
  getContractStatus,
} from '../contracts';

// ---------------------------------------------------------------------------
// getActiveContract
// ---------------------------------------------------------------------------
describe('getActiveContract', () => {
  beforeEach(() => {
    jest.useFakeTimers();
    jest.setSystemTime(new Date('2025-06-15'));
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it('returns null for undefined', () => {
    expect(getActiveContract(undefined)).toBeNull();
  });

  it('returns null for empty array', () => {
    expect(getActiveContract([])).toBeNull();
  });

  it('returns active contract (no end date)', () => {
    const contracts = [{ from: '2025-01-01' }];
    expect(getActiveContract(contracts)).toBe(contracts[0]);
  });

  it('returns active contract (end date in future)', () => {
    const contracts = [{ from: '2025-01-01', to: '2025-12-31' }];
    expect(getActiveContract(contracts)).toBe(contracts[0]);
  });

  it('returns null when contract has not started yet', () => {
    const contracts = [{ from: '2025-09-01', to: '2025-12-31' }];
    expect(getActiveContract(contracts)).toBeNull();
  });

  it('returns null when contract has ended', () => {
    const contracts = [{ from: '2024-01-01', to: '2025-01-01' }];
    expect(getActiveContract(contracts)).toBeNull();
  });

  it('returns active contract among multiple', () => {
    const contracts = [
      { from: '2024-01-01', to: '2024-12-31' },
      { from: '2025-01-01', to: '2025-12-31' },
      { from: '2026-01-01', to: '2026-12-31' },
    ];
    expect(getActiveContract(contracts)).toBe(contracts[1]);
  });

  it('handles null to value', () => {
    const contracts = [{ from: '2025-01-01', to: null }];
    expect(getActiveContract(contracts)).toBe(contracts[0]);
  });
});

// ---------------------------------------------------------------------------
// getCurrentContract
// ---------------------------------------------------------------------------
describe('getCurrentContract', () => {
  beforeEach(() => {
    jest.useFakeTimers();
    jest.setSystemTime(new Date('2025-06-15'));
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it('returns null for undefined', () => {
    expect(getCurrentContract(undefined)).toBeNull();
  });

  it('returns null for empty array', () => {
    expect(getCurrentContract([])).toBeNull();
  });

  it('returns active contract when one exists', () => {
    const contracts = [{ from: '2025-01-01', to: '2025-12-31' }];
    expect(getCurrentContract(contracts)).toBe(contracts[0]);
  });

  it('falls back to contract with latest start date', () => {
    const contracts = [
      { from: '2023-01-01', to: '2023-12-31' },
      { from: '2024-06-01', to: '2024-12-31' },
      { from: '2024-01-01', to: '2024-06-30' },
    ];
    expect(getCurrentContract(contracts)).toEqual(contracts[1]);
  });

  it('prefers active contract over later-starting ended contract', () => {
    const contracts = [{ from: '2024-01-01', to: '2024-12-31' }, { from: '2025-01-01' }];
    expect(getCurrentContract(contracts)).toBe(contracts[1]);
  });

  it('does not mutate the original array when sorting', () => {
    const contracts = [
      { from: '2023-01-01', to: '2023-12-31' },
      { from: '2024-06-01', to: '2024-12-31' },
    ];
    const copy = [...contracts];
    getCurrentContract(contracts);
    expect(contracts).toEqual(copy);
  });
});

// ---------------------------------------------------------------------------
// getDayBefore
// ---------------------------------------------------------------------------
describe('getDayBefore', () => {
  it('returns day before a date', () => {
    expect(getDayBefore('2025-06-15')).toBe('2025-06-14');
  });

  it('crosses month boundary', () => {
    expect(getDayBefore('2025-03-01')).toBe('2025-02-28');
  });

  it('crosses year boundary', () => {
    expect(getDayBefore('2025-01-01')).toBe('2024-12-31');
  });

  it('handles leap year', () => {
    expect(getDayBefore('2024-03-01')).toBe('2024-02-29');
  });
});

// ---------------------------------------------------------------------------
// getContractStatus
// ---------------------------------------------------------------------------
describe('getContractStatus', () => {
  beforeEach(() => {
    jest.useFakeTimers();
    jest.setSystemTime(new Date('2025-06-15'));
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it('returns null for null contract', () => {
    expect(getContractStatus(null)).toBeNull();
  });

  it('returns active for current contract (no end date)', () => {
    expect(getContractStatus({ from: '2025-01-01' })).toBe('active');
  });

  it('returns active for current contract (end date in future)', () => {
    expect(getContractStatus({ from: '2025-01-01', to: '2025-12-31' })).toBe('active');
  });

  it('returns upcoming for future contract', () => {
    expect(getContractStatus({ from: '2025-09-01' })).toBe('upcoming');
  });

  it('returns ended for past contract', () => {
    expect(getContractStatus({ from: '2024-01-01', to: '2025-01-01' })).toBe('ended');
  });

  it('returns active for contract ending today', () => {
    expect(getContractStatus({ from: '2025-01-01', to: '2025-06-15' })).toBe('active');
  });

  it('returns active for contract starting today', () => {
    expect(getContractStatus({ from: '2025-06-15' })).toBe('active');
  });

  it('handles null to value', () => {
    expect(getContractStatus({ from: '2025-01-01', to: null })).toBe('active');
  });
});
