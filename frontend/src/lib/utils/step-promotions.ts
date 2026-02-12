/**
 * Calculate years of service from contract history.
 * Uses the earliest contract start date.
 */
export function calculateYearsOfService(
  contracts: { from: string }[],
  asOf: Date = new Date()
): number {
  if (contracts.length === 0) return 0;

  const earliest = contracts.reduce((min, c) => {
    const d = new Date(c.from);
    return d < min ? d : min;
  }, new Date(contracts[0].from));

  const diffMs = asOf.getTime() - earliest.getTime();
  if (diffMs <= 0) return 0;

  return diffMs / (365.25 * 24 * 60 * 60 * 1000);
}

/**
 * Determine the eligible step based on years of service and pay plan entries.
 * Only considers entries that have step_min_years defined.
 * Returns 0 if no entries have step rules.
 */
export function determineEligibleStep(
  yearsOfService: number,
  entries: { step: number; grade: string; step_min_years?: number | null }[],
  grade: string
): number {
  const eligible = entries.filter(
    (e) => e.grade === grade && e.step_min_years != null && e.step_min_years <= yearsOfService
  );

  if (eligible.length === 0) return 0;

  return Math.max(...eligible.map((e) => e.step));
}
