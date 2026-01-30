import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/lib/api/client';
import { useUiStore } from '@/stores/ui-store';

/**
 * Hook to get suggested attributes for a child contract based on the
 * organization's government funding configuration.
 *
 * Returns unique property names from funding periods that overlap with
 * the given date range.
 */
export function useFundingAttributes(orgId: number, fromDate?: string, toDate?: string) {
  const organizations = useUiStore((state) => state.organizations);
  const org = organizations.find((o) => o.id === orgId);
  const state = org?.state;

  // Fetch all government fundings (cached by React Query)
  const { data: fundingsResponse } = useQuery({
    queryKey: ['governmentFundings', 'all'],
    queryFn: () => apiClient.getGovernmentFundings({ page: 1, limit: 100 }),
    staleTime: 5 * 60 * 1000, // 5 minutes - funding configs rarely change
    enabled: !!state,
  });

  // Find funding for this org's state
  const funding = fundingsResponse?.data?.find((f) => f.state === state);

  // Fetch funding details with all periods if we found a matching funding
  const { data: fundingDetails } = useQuery({
    queryKey: ['governmentFunding', funding?.id, 'details'],
    queryFn: () => apiClient.getGovernmentFunding(funding!.id, 0), // 0 = all periods
    staleTime: 5 * 60 * 1000,
    enabled: !!funding?.id,
  });

  // Extract unique attribute names from periods that overlap with the contract date range
  const suggestedAttributes: string[] = [];

  if (fundingDetails?.periods) {
    const contractFrom = fromDate || '';
    const contractTo = toDate || '9999-12-31'; // Far future if no end date

    for (const period of fundingDetails.periods) {
      const periodFrom = period.from;
      const periodTo = period.to || '9999-12-31';

      // Check if periods overlap: period starts before contract ends AND period ends after contract starts
      const overlaps = periodFrom <= contractTo && periodTo >= contractFrom;

      if (overlaps && period.properties) {
        for (const prop of period.properties) {
          if (prop.name && !suggestedAttributes.includes(prop.name.toLowerCase())) {
            suggestedAttributes.push(prop.name.toLowerCase());
          }
        }
      }
    }
  }

  return {
    suggestedAttributes: suggestedAttributes.sort(),
    isLoading: !fundingDetails && !!state,
    hasNoFunding: !!state && !funding,
  };
}
