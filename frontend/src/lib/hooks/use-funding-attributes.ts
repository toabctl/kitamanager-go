import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/lib/api/client';
import { useUiStore } from '@/stores/ui-store';

export interface AttributeSuggestion {
  name: string;
  exclusiveGroup?: string | null;
}

export interface GroupedSuggestions {
  [group: string]: string[]; // group name -> attribute names
}

/**
 * Hook to get suggested attributes for a child contract based on the
 * organization's government funding configuration.
 *
 * Returns unique property names from funding periods that overlap with
 * the given date range, along with their exclusive_group information.
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

  // Extract unique attributes with their exclusive groups
  const attributeMap = new Map<string, string | null>();

  if (fundingDetails?.periods) {
    const contractFrom = fromDate || '';
    const contractTo = toDate || '9999-12-31'; // Far future if no end date

    for (const period of fundingDetails.periods) {
      const periodFrom = period.from;
      const periodTo = period.to || '9999-12-31';

      // Check if periods overlap
      const overlaps = periodFrom <= contractTo && periodTo >= contractFrom;

      if (overlaps && period.properties) {
        for (const prop of period.properties) {
          const name = prop.name?.toLowerCase();
          if (name && !attributeMap.has(name)) {
            attributeMap.set(name, prop.exclusive_group || null);
          }
        }
      }
    }
  }

  // Build flat list for simple usage
  const suggestedAttributes = Array.from(attributeMap.keys()).sort();

  // Build grouped structure for advanced UI
  const groupedSuggestions: GroupedSuggestions = {};
  const ungrouped: string[] = [];

  Array.from(attributeMap.entries()).forEach(([name, group]) => {
    if (group) {
      if (!groupedSuggestions[group]) {
        groupedSuggestions[group] = [];
      }
      groupedSuggestions[group].push(name);
    } else {
      ungrouped.push(name);
    }
  });

  // Sort within groups
  for (const group of Object.keys(groupedSuggestions)) {
    groupedSuggestions[group].sort();
  }
  ungrouped.sort();

  // Map from attribute name to its exclusive group (for conflict detection)
  const exclusiveGroupMap: Record<string, string | null> = {};
  Array.from(attributeMap.entries()).forEach(([name, group]) => {
    exclusiveGroupMap[name] = group;
  });

  return {
    suggestedAttributes,
    groupedSuggestions,
    ungrouped,
    exclusiveGroupMap,
    isLoading: !fundingDetails && !!state,
    hasNoFunding: !!state && !funding,
  };
}
