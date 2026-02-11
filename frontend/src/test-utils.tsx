import React from 'react';
import { render } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

export function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });
}

export function renderWithProviders(ui: React.ReactNode) {
  const queryClient = createTestQueryClient();
  return render(<QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>);
}

export function createHookWrapper(queryClient?: QueryClient) {
  const qc = queryClient ?? createTestQueryClient();
  function TestWrapper({ children }: { children: React.ReactNode }) {
    return React.createElement(QueryClientProvider, { client: qc }, children);
  }
  return TestWrapper;
}

export function createMockPaginatedResponse<T>(data: T[], overrides = {}) {
  return { data, total: data.length, page: 1, limit: 30, total_pages: 1, ...overrides };
}
