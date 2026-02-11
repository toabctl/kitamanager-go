import { screen, waitFor } from '@testing-library/react';
import EmployeesPage from '../page';
import { apiClient } from '@/lib/api/client';
import { renderWithProviders, createMockPaginatedResponse } from '@/test-utils';

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

jest.mock('@/lib/api/client', () => ({
  apiClient: {
    getEmployees: jest.fn(),
    createEmployee: jest.fn(),
    updateEmployee: jest.fn(),
    deleteEmployee: jest.fn(),
    createEmployeeContract: jest.fn(),
    updateEmployeeContract: jest.fn(),
  },
  getErrorMessage: jest.fn((e: unknown, f: string) => f),
}));

const mockEmployees = [
  {
    id: 1,
    first_name: 'John',
    last_name: 'Doe',
    gender: 'male',
    birthdate: '1990-01-15T00:00:00Z',
    organization_id: 1,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
    contracts: [
      {
        id: 1,
        employee_id: 1,
        from: '2020-01-01',
        staff_category: 'qualified',
        grade: 'S8a',
        step: 3,
        weekly_hours: 39,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      },
    ],
  },
  {
    id: 2,
    first_name: 'Jane',
    last_name: 'Smith',
    gender: 'female',
    birthdate: '1985-06-20T00:00:00Z',
    organization_id: 1,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
    contracts: [],
  },
];

const mockPaginatedResponse = createMockPaginatedResponse(mockEmployees);
const mockEmptyResponse = createMockPaginatedResponse([]);

describe('EmployeesPage', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders page title', async () => {
    (apiClient.getEmployees as jest.Mock).mockResolvedValue(mockEmptyResponse);

    renderWithProviders(<EmployeesPage />);

    const titles = screen.getAllByText('employees.title');
    expect(titles.length).toBeGreaterThanOrEqual(1);
  });

  it('renders new employee button', async () => {
    (apiClient.getEmployees as jest.Mock).mockResolvedValue(mockEmptyResponse);

    renderWithProviders(<EmployeesPage />);

    expect(screen.getByText('employees.newEmployee')).toBeInTheDocument();
  });

  it('shows loading skeleton while fetching', async () => {
    (apiClient.getEmployees as jest.Mock).mockImplementation(
      () => new Promise(() => {}) // Never resolves
    );

    renderWithProviders(<EmployeesPage />);

    const skeletons = document.querySelectorAll('.animate-pulse');
    expect(skeletons.length).toBeGreaterThan(0);
  });

  it('displays employees in table', async () => {
    (apiClient.getEmployees as jest.Mock).mockResolvedValue(mockPaginatedResponse);

    renderWithProviders(<EmployeesPage />);

    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument();
    });

    expect(screen.getByText('Jane Smith')).toBeInTheDocument();
  });

  it('shows no results when empty', async () => {
    (apiClient.getEmployees as jest.Mock).mockResolvedValue(mockEmptyResponse);

    renderWithProviders(<EmployeesPage />);

    await waitFor(() => {
      expect(screen.getByText('common.noResults')).toBeInTheDocument();
    });
  });

  it('renders search input', async () => {
    (apiClient.getEmployees as jest.Mock).mockResolvedValue(mockEmptyResponse);

    renderWithProviders(<EmployeesPage />);

    expect(screen.getByPlaceholderText('common.search')).toBeInTheDocument();
  });

  it('renders table headers', async () => {
    (apiClient.getEmployees as jest.Mock).mockResolvedValue(mockEmptyResponse);

    renderWithProviders(<EmployeesPage />);

    await waitFor(() => {
      expect(screen.getByText('common.noResults')).toBeInTheDocument();
    });

    expect(screen.getByText('common.name')).toBeInTheDocument();
    expect(screen.getByText('gender.label')).toBeInTheDocument();
    expect(screen.getByText('employees.birthdate')).toBeInTheDocument();
    expect(screen.getByText('employees.age')).toBeInTheDocument();
    expect(screen.getByText('employees.staffCategory.label')).toBeInTheDocument();
    expect(screen.getByText('employees.grade')).toBeInTheDocument();
    expect(screen.getByText('employees.weeklyHours')).toBeInTheDocument();
    expect(screen.getByText('common.actions')).toBeInTheDocument();
  });
});
