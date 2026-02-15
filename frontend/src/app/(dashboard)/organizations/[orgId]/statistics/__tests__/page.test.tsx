import { screen, waitFor } from '@testing-library/react';
import StatisticsPage from '../page';
import { apiClient } from '@/lib/api/client';
import { renderWithProviders } from '@/test-utils';

jest.mock('@/lib/api/client', () => ({
  apiClient: {
    getAgeDistribution: jest.fn(),
    getChildrenContractCountByMonth: jest.fn(),
    getContractPropertiesDistribution: jest.fn(),
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

jest.mock('@/components/charts/age-distribution-chart', () => ({
  AgeDistributionChart: () => <div data-testid="age-chart">Age Chart</div>,
}));

jest.mock('@/components/charts/monthly-contract-chart', () => ({
  MonthlyContractChart: () => <div data-testid="contract-chart">Contract Chart</div>,
}));

jest.mock('@/components/charts/contract-properties-chart', () => ({
  ContractPropertiesChart: () => (
    <div data-testid="contract-properties-chart">Contract Properties Chart</div>
  ),
}));

const mockAgeDistribution = [
  { age: 1, count: 5 },
  { age: 2, count: 8 },
  { age: 3, count: 12 },
];

const mockContractCounts = [
  { month: '2024-01', count: 20 },
  { month: '2024-02', count: 22 },
];

const mockContractProperties = {
  date: '2024-01-01',
  total_children: 20,
  properties: [
    { key: 'care_type', value: 'ganztag', count: 15 },
    { key: 'care_type', value: 'halbtag', count: 5 },
  ],
};

describe('StatisticsPage', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders page title', async () => {
    (apiClient.getAgeDistribution as jest.Mock).mockResolvedValue(mockAgeDistribution);
    (apiClient.getChildrenContractCountByMonth as jest.Mock).mockResolvedValue(mockContractCounts);
    (apiClient.getContractPropertiesDistribution as jest.Mock).mockResolvedValue(
      mockContractProperties
    );

    renderWithProviders(<StatisticsPage />);

    expect(screen.getByText('statistics.title')).toBeInTheDocument();
  });

  it('renders card titles', async () => {
    (apiClient.getAgeDistribution as jest.Mock).mockResolvedValue(mockAgeDistribution);
    (apiClient.getChildrenContractCountByMonth as jest.Mock).mockResolvedValue(mockContractCounts);
    (apiClient.getContractPropertiesDistribution as jest.Mock).mockResolvedValue(
      mockContractProperties
    );

    renderWithProviders(<StatisticsPage />);

    expect(screen.getByText('statistics.ageDistribution')).toBeInTheDocument();
    expect(screen.getByText('statistics.childrenContractCount')).toBeInTheDocument();
    expect(screen.getByText('statistics.contractProperties')).toBeInTheDocument();
  });

  it('shows loading skeletons while fetching', async () => {
    (apiClient.getAgeDistribution as jest.Mock).mockImplementation(() => new Promise(() => {}));
    (apiClient.getChildrenContractCountByMonth as jest.Mock).mockImplementation(
      () => new Promise(() => {})
    );
    (apiClient.getContractPropertiesDistribution as jest.Mock).mockImplementation(
      () => new Promise(() => {})
    );

    renderWithProviders(<StatisticsPage />);

    const skeletons = document.querySelectorAll('.animate-pulse');
    expect(skeletons.length).toBeGreaterThan(0);
  });

  it('renders charts when data is loaded', async () => {
    (apiClient.getAgeDistribution as jest.Mock).mockResolvedValue(mockAgeDistribution);
    (apiClient.getChildrenContractCountByMonth as jest.Mock).mockResolvedValue(mockContractCounts);
    (apiClient.getContractPropertiesDistribution as jest.Mock).mockResolvedValue(
      mockContractProperties
    );

    renderWithProviders(<StatisticsPage />);

    await waitFor(() => {
      expect(screen.getByTestId('age-chart')).toBeInTheDocument();
    });

    expect(screen.getByTestId('contract-chart')).toBeInTheDocument();
  });

  it('renders contract properties chart when data is loaded', async () => {
    (apiClient.getAgeDistribution as jest.Mock).mockResolvedValue(mockAgeDistribution);
    (apiClient.getChildrenContractCountByMonth as jest.Mock).mockResolvedValue(mockContractCounts);
    (apiClient.getContractPropertiesDistribution as jest.Mock).mockResolvedValue(
      mockContractProperties
    );

    renderWithProviders(<StatisticsPage />);

    await waitFor(() => {
      expect(screen.getByTestId('contract-properties-chart')).toBeInTheDocument();
    });
  });
});
