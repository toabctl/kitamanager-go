import { screen, waitFor } from '@testing-library/react';
import CostDetailPage from '../page';
import { apiClient } from '@/lib/api/client';
import { renderWithProviders } from '@/test-utils';

jest.mock('@/lib/api/client', () => ({
  apiClient: {
    getCost: jest.fn(),
    createCostEntry: jest.fn(),
    updateCostEntry: jest.fn(),
    deleteCostEntry: jest.fn(),
  },
  getErrorMessage: jest.fn((error, fallback) => fallback),
}));

jest.mock('next/navigation', () => ({
  useParams: () => ({ orgId: '1', id: '1' }),
  useRouter: () => ({ push: jest.fn(), back: jest.fn(), refresh: jest.fn() }),
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

const mockCostWithEntries = {
  id: 1,
  organization_id: 1,
  name: 'Rent',
  entries: [
    {
      id: 10,
      cost_id: 1,
      from: '2024-01-01T00:00:00Z',
      to: '2024-12-31T00:00:00Z',
      amount_cents: 150000,
      notes: 'Monthly office rent',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    },
  ],
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
};

const mockCostNoEntries = {
  id: 1,
  organization_id: 1,
  name: 'Insurance',
  entries: [],
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
};

describe('CostDetailPage', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders loading state', async () => {
    (apiClient.getCost as jest.Mock).mockImplementation(() => new Promise(() => {}));

    renderWithProviders(<CostDetailPage />);

    const skeletons = document.querySelectorAll('.animate-pulse');
    expect(skeletons.length).toBeGreaterThan(0);
  });

  it('renders cost name', async () => {
    (apiClient.getCost as jest.Mock).mockResolvedValue(mockCostWithEntries);

    renderWithProviders(<CostDetailPage />);

    await waitFor(() => {
      expect(screen.getByText('Rent')).toBeInTheDocument();
    });
  });

  it('shows entries with amount and notes', async () => {
    (apiClient.getCost as jest.Mock).mockResolvedValue(mockCostWithEntries);

    renderWithProviders(<CostDetailPage />);

    await waitFor(() => {
      expect(screen.getByText('Rent')).toBeInTheDocument();
    });

    expect(screen.getByText('Monthly office rent')).toBeInTheDocument();
  });

  it('shows add entry button', async () => {
    (apiClient.getCost as jest.Mock).mockResolvedValue(mockCostWithEntries);

    renderWithProviders(<CostDetailPage />);

    await waitFor(() => {
      expect(screen.getByText('Rent')).toBeInTheDocument();
    });

    expect(screen.getByText('costs.addEntry')).toBeInTheDocument();
  });

  it('shows no entries message when empty', async () => {
    (apiClient.getCost as jest.Mock).mockResolvedValue(mockCostNoEntries);

    renderWithProviders(<CostDetailPage />);

    await waitFor(() => {
      expect(screen.getByText('Insurance')).toBeInTheDocument();
    });

    expect(screen.getByText('costs.noEntriesDefined')).toBeInTheDocument();
  });

  it('renders back button', async () => {
    (apiClient.getCost as jest.Mock).mockResolvedValue(mockCostWithEntries);

    renderWithProviders(<CostDetailPage />);

    await waitFor(() => {
      expect(screen.getByText('Rent')).toBeInTheDocument();
    });

    const buttons = screen.getAllByRole('button');
    expect(buttons.length).toBeGreaterThanOrEqual(1);
  });

  it('shows error when cost fails to load', async () => {
    (apiClient.getCost as jest.Mock).mockResolvedValue(null);

    renderWithProviders(<CostDetailPage />);

    await waitFor(() => {
      expect(screen.getByText('costs.failedToLoadCost')).toBeInTheDocument();
    });
  });
});
