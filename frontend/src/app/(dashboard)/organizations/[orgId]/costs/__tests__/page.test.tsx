import { screen, waitFor } from '@testing-library/react';
import CostsPage from '../page';
import { apiClient } from '@/lib/api/client';
import { renderWithProviders, createMockPaginatedResponse } from '@/test-utils';

jest.mock('@/lib/api/client', () => ({
  apiClient: {
    getCosts: jest.fn(),
    createCost: jest.fn(),
    updateCost: jest.fn(),
    deleteCost: jest.fn(),
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

const mockCosts = [
  {
    id: 1,
    name: 'Rent',
    organization_id: 1,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  },
  {
    id: 2,
    name: 'Insurance',
    organization_id: 1,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  },
];

const mockPaginatedResponse = createMockPaginatedResponse(mockCosts);
const mockEmptyResponse = createMockPaginatedResponse([]);

describe('CostsPage', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders page title', async () => {
    (apiClient.getCosts as jest.Mock).mockResolvedValue(mockEmptyResponse);

    renderWithProviders(<CostsPage />);

    const titles = screen.getAllByText('costs.title');
    expect(titles.length).toBeGreaterThanOrEqual(1);
  });

  it('renders new cost button', async () => {
    (apiClient.getCosts as jest.Mock).mockResolvedValue(mockEmptyResponse);

    renderWithProviders(<CostsPage />);

    expect(screen.getByText('costs.newCost')).toBeInTheDocument();
  });

  it('shows loading skeleton while fetching', async () => {
    (apiClient.getCosts as jest.Mock).mockImplementation(() => new Promise(() => {}));

    renderWithProviders(<CostsPage />);

    const skeletons = document.querySelectorAll('.animate-pulse');
    expect(skeletons.length).toBeGreaterThan(0);
  });

  it('displays costs in table', async () => {
    (apiClient.getCosts as jest.Mock).mockResolvedValue(mockPaginatedResponse);

    renderWithProviders(<CostsPage />);

    await waitFor(() => {
      expect(screen.getByText('Rent')).toBeInTheDocument();
    });

    expect(screen.getByText('Insurance')).toBeInTheDocument();
  });

  it('shows no results when empty', async () => {
    (apiClient.getCosts as jest.Mock).mockResolvedValue(mockEmptyResponse);

    renderWithProviders(<CostsPage />);

    await waitFor(() => {
      expect(screen.getByText('common.noResults')).toBeInTheDocument();
    });
  });
});
