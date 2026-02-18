import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import BudgetPage from '../page';
import { apiClient } from '@/lib/api/client';
import { renderWithProviders } from '@/test-utils';

jest.mock('@/lib/api/client', () => ({
  apiClient: {
    getFinancials: jest.fn(),
  },
  getErrorMessage: jest.fn((error, fallback) => fallback),
}));

jest.mock('next/navigation', () => ({
  useParams: () => ({ orgId: '1' }),
  useRouter: () => ({ push: jest.fn() }),
}));

jest.mock('next-intl', () => ({
  useTranslations: () => (key: string) => key,
}));

jest.mock('@/lib/hooks/use-toast', () => ({
  useToast: () => ({ toast: jest.fn() }),
}));

const mockFinancials = {
  data_points: [
    {
      date: '2026-01-01',
      funding_income: 500000,
      gross_salary: 200000,
      employer_costs: 50000,
      budget_income: 10000,
      budget_expenses: 5000,
      total_income: 510000,
      total_expenses: 255000,
      balance: 255000,
      child_count: 10,
      staff_count: 5,
    },
    {
      date: '2026-02-01',
      funding_income: 500000,
      gross_salary: 200000,
      employer_costs: 50000,
      budget_income: 10000,
      budget_expenses: 5000,
      total_income: 510000,
      total_expenses: 255000,
      balance: 255000,
      child_count: 10,
      staff_count: 5,
    },
  ],
};

describe('BudgetPage', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders page title', async () => {
    (apiClient.getFinancials as jest.Mock).mockResolvedValue(mockFinancials);

    renderWithProviders(<BudgetPage />);

    expect(screen.getByText('nav.statisticsBudget')).toBeInTheDocument();
  });

  it('shows loading skeleton while fetching', () => {
    (apiClient.getFinancials as jest.Mock).mockImplementation(() => new Promise(() => {}));

    renderWithProviders(<BudgetPage />);

    const skeletons = document.querySelectorAll('.animate-pulse');
    expect(skeletons.length).toBeGreaterThan(0);
  });

  it('renders table when data loads', async () => {
    (apiClient.getFinancials as jest.Mock).mockResolvedValue(mockFinancials);

    renderWithProviders(<BudgetPage />);

    await waitFor(() => {
      expect(screen.getByText('Jan. 26')).toBeInTheDocument();
    });

    expect(screen.getByText('Feb. 26')).toBeInTheDocument();
  });

  it('renders year stepper with navigation buttons', () => {
    (apiClient.getFinancials as jest.Mock).mockResolvedValue(mockFinancials);

    renderWithProviders(<BudgetPage />);

    expect(screen.getByLabelText('previousYear')).toBeInTheDocument();
    expect(screen.getByLabelText('nextYear')).toBeInTheDocument();
  });

  it('changes query when year stepper is used', async () => {
    (apiClient.getFinancials as jest.Mock).mockResolvedValue(mockFinancials);

    renderWithProviders(<BudgetPage />);

    await userEvent.click(screen.getByLabelText('nextYear'));

    await waitFor(() => {
      const calls = (apiClient.getFinancials as jest.Mock).mock.calls;
      const lastCall = calls[calls.length - 1];
      expect(lastCall[1]).toEqual({
        from: `${new Date().getFullYear() + 1}-01-01`,
        to: `${new Date().getFullYear() + 1}-12-01`,
      });
    });
  });

  it('shows error fallback when API fails', async () => {
    (apiClient.getFinancials as jest.Mock).mockRejectedValue(new Error('API error'));

    renderWithProviders(<BudgetPage />);

    await waitFor(() => {
      expect(screen.getByText('statistics.chartError')).toBeInTheDocument();
    });
  });

  it('handles empty data points', async () => {
    (apiClient.getFinancials as jest.Mock).mockResolvedValue({ data_points: [] });

    renderWithProviders(<BudgetPage />);

    await waitFor(() => {
      // BudgetTable shows chartError for empty data
      expect(screen.getByText('chartError')).toBeInTheDocument();
    });
  });
});
