import { screen, waitFor } from '@testing-library/react';
import PayPlansPage from '../page';
import { apiClient } from '@/lib/api/client';
import { renderWithProviders, createMockPaginatedResponse } from '@/test-utils';

jest.mock('@/lib/api/client', () => ({
  apiClient: {
    getPayPlans: jest.fn(),
    createPayPlan: jest.fn(),
    updatePayPlan: jest.fn(),
    deletePayPlan: jest.fn(),
  },
  getErrorMessage: jest.fn((error, fallback) => fallback),
}));

jest.mock('next/navigation', () => ({
  useParams: () => ({ orgId: '1' }),
  useRouter: () => ({ push: jest.fn() }),
}));

jest.mock('next-intl', () => ({
  useTranslations: () => (key: string, params?: Record<string, unknown>) => {
    if (params) return `${key}`;
    return key;
  },
}));

jest.mock('@/lib/hooks/use-toast', () => ({
  useToast: () => ({ toast: jest.fn() }),
}));

const mockPayPlans = [
  {
    id: 1,
    name: 'TV-L Berlin',
    organization_id: 1,
    total_periods: 3,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  },
  {
    id: 2,
    name: 'TV-L Brandenburg',
    organization_id: 1,
    total_periods: 2,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  },
];

const mockPaginatedResponse = createMockPaginatedResponse(mockPayPlans);
const mockEmptyResponse = createMockPaginatedResponse([]);

describe('PayPlansPage', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders page title', async () => {
    (apiClient.getPayPlans as jest.Mock).mockResolvedValue(mockEmptyResponse);

    renderWithProviders(<PayPlansPage />);

    const titles = screen.getAllByText('payPlans.title');
    expect(titles.length).toBeGreaterThanOrEqual(1);
  });

  it('renders new payplan button', async () => {
    (apiClient.getPayPlans as jest.Mock).mockResolvedValue(mockEmptyResponse);

    renderWithProviders(<PayPlansPage />);

    expect(screen.getByText('payPlans.newPayPlan')).toBeInTheDocument();
  });

  it('shows loading skeleton while fetching', async () => {
    (apiClient.getPayPlans as jest.Mock).mockImplementation(() => new Promise(() => {}));

    renderWithProviders(<PayPlansPage />);

    const skeletons = document.querySelectorAll('.animate-pulse');
    expect(skeletons.length).toBeGreaterThan(0);
  });

  it('displays payplans in table', async () => {
    (apiClient.getPayPlans as jest.Mock).mockResolvedValue(mockPaginatedResponse);

    renderWithProviders(<PayPlansPage />);

    await waitFor(() => {
      expect(screen.getByText('TV-L Berlin')).toBeInTheDocument();
    });

    expect(screen.getByText('TV-L Brandenburg')).toBeInTheDocument();
  });

  it('shows no results when empty', async () => {
    (apiClient.getPayPlans as jest.Mock).mockResolvedValue(mockEmptyResponse);

    renderWithProviders(<PayPlansPage />);

    await waitFor(() => {
      expect(screen.getByText('common.noResults')).toBeInTheDocument();
    });
  });
});
