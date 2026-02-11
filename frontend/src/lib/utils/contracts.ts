/**
 * Get the currently active contract (from <= today, no end date or end date >= today)
 */
export function getActiveContract<T extends { from: string; to?: string | null }>(
  contracts?: T[]
): T | null {
  if (!contracts || contracts.length === 0) return null;
  const today = new Date().toISOString().split('T')[0];
  return contracts.find((c) => c.from <= today && (!c.to || c.to >= today)) || null;
}

/**
 * Get the current or most recent contract.
 * Falls back to the contract with the latest start date.
 */
export function getCurrentContract<T extends { from: string; to?: string | null }>(
  contracts?: T[]
): T | null {
  if (!contracts || contracts.length === 0) return null;
  const today = new Date().toISOString().split('T')[0];
  return (
    contracts.find((c) => c.from <= today && (!c.to || c.to >= today)) ||
    [...contracts].sort((a, b) => b.from.localeCompare(a.from))[0]
  );
}

/**
 * Get the day before a given date string (YYYY-MM-DD format)
 */
export function getDayBefore(dateStr: string): string {
  const date = new Date(dateStr);
  date.setDate(date.getDate() - 1);
  return date.toISOString().split('T')[0];
}

/**
 * Get the status of a contract relative to today
 */
export function getContractStatus(
  contract: { from: string; to?: string | null } | null
): 'active' | 'upcoming' | 'ended' | null {
  if (!contract) return null;
  const today = new Date().toISOString().split('T')[0];
  if (contract.from > today) return 'upcoming';
  if (contract.to && contract.to < today) return 'ended';
  return 'active';
}
