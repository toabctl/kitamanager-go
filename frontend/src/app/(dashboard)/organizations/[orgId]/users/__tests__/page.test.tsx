import { screen, waitFor } from '@testing-library/react';
import UsersPage from '../page';
import { apiClient } from '@/lib/api/client';
import { renderWithProviders, createMockPaginatedResponse } from '@/test-utils';

jest.mock('@/lib/api/client', () => ({
  apiClient: {
    getUsers: jest.fn(),
    createUser: jest.fn(),
    updateUser: jest.fn(),
    deleteUser: jest.fn(),
    setSuperAdmin: jest.fn(),
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

jest.mock('@/stores/auth-store', () => ({
  useAuthStore: () => ({
    user: { id: 1, is_superadmin: true, name: 'Admin', email: 'admin@test.com' },
  }),
}));

const mockUsers = [
  {
    id: 1,
    name: 'Admin',
    email: 'admin@test.com',
    active: true,
    is_superadmin: true,
    created_at: '2024-01-01T00:00:00Z',
    created_by: 'system',
    updated_at: '2024-01-01T00:00:00Z',
  },
  {
    id: 2,
    name: 'User',
    email: 'user@test.com',
    active: false,
    is_superadmin: false,
    created_at: '2024-01-01T00:00:00Z',
    created_by: 'admin',
    updated_at: '2024-01-01T00:00:00Z',
  },
];

const mockPaginatedResponse = createMockPaginatedResponse(mockUsers);
const mockEmptyResponse = createMockPaginatedResponse([]);

describe('UsersPage', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders page title', async () => {
    (apiClient.getUsers as jest.Mock).mockResolvedValue(mockEmptyResponse);

    renderWithProviders(<UsersPage />);

    const titles = screen.getAllByText('users.title');
    expect(titles.length).toBeGreaterThanOrEqual(1);
  });

  it('renders new user button', async () => {
    (apiClient.getUsers as jest.Mock).mockResolvedValue(mockEmptyResponse);

    renderWithProviders(<UsersPage />);

    expect(screen.getByText('users.newUser')).toBeInTheDocument();
  });

  it('displays users in table', async () => {
    (apiClient.getUsers as jest.Mock).mockResolvedValue(mockPaginatedResponse);

    renderWithProviders(<UsersPage />);

    await waitFor(() => {
      expect(screen.getByText('Admin')).toBeInTheDocument();
    });

    expect(screen.getByText('User')).toBeInTheDocument();
  });

  it('shows active/inactive badges', async () => {
    (apiClient.getUsers as jest.Mock).mockResolvedValue(mockPaginatedResponse);

    renderWithProviders(<UsersPage />);

    await waitFor(() => {
      expect(screen.getByText('Admin')).toBeInTheDocument();
    });

    expect(screen.getByText('common.active')).toBeInTheDocument();
    expect(screen.getByText('common.inactive')).toBeInTheDocument();
  });

  it('shows no results when empty', async () => {
    (apiClient.getUsers as jest.Mock).mockResolvedValue(mockEmptyResponse);

    renderWithProviders(<UsersPage />);

    await waitFor(() => {
      expect(screen.getByText('common.noResults')).toBeInTheDocument();
    });
  });

  it('renders superadmin column for superadmin users', async () => {
    (apiClient.getUsers as jest.Mock).mockResolvedValue(mockPaginatedResponse);

    renderWithProviders(<UsersPage />);

    await waitFor(() => {
      expect(screen.getByText('Admin')).toBeInTheDocument();
    });

    expect(screen.getByText('users.superadmin')).toBeInTheDocument();
  });
});
