'use client';

import { useState, useCallback } from 'react';
import { useDebouncedValue } from './use-debounced-value';

interface UseResourceListFiltersConfig {
  /** Debounce delay for search input in ms (default: 300) */
  debounceMs?: number;
}

/**
 * Shared hook for resource list page filter state:
 * pagination, debounced search, date filter, and a generic category filter.
 *
 * All setters automatically reset the page to 1.
 */
export function useResourceListFilters({ debounceMs = 300 }: UseResourceListFiltersConfig = {}) {
  const [page, setPage] = useState(1);
  const [searchInput, setSearchInput] = useState('');
  const search = useDebouncedValue(searchInput, debounceMs);
  const [activeOn, setActiveOn] = useState(() => new Date());

  const setSearchAndResetPage = useCallback((value: string) => {
    setSearchInput(value);
    setPage(1);
  }, []);

  const setActiveOnAndResetPage = useCallback((date: Date) => {
    setActiveOn(date);
    setPage(1);
  }, []);

  return {
    page,
    setPage,
    searchInput,
    setSearchInput: setSearchAndResetPage,
    search,
    activeOn,
    setActiveOn: setActiveOnAndResetPage,
  };
}
