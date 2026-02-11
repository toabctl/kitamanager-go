import { screen, waitFor } from '@testing-library/react';
import GroupsPage from '../page';
import { apiClient } from '@/lib/api/client';
import { renderWithProviders, createMockPaginatedResponse } from '@/test-utils';

jest.mock('@/lib/api/client', () => ({
  apiClient: {
    getGroups: jest.fn(),
    createGroup: jest.fn(),
    updateGroup: jest.fn(),
    deleteGroup: jest.fn(),
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

const mockGroups = [
  {
    id: 1,
    name: 'Group A',
    organization_id: 1,
    active: true,
    created_at: '2024-01-01T00:00:00Z',
    created_by: 'admin',
    updated_at: '2024-01-01T00:00:00Z',
  },
  {
    id: 2,
    name: 'Group B',
    organization_id: 1,
    active: false,
    created_at: '2024-01-01T00:00:00Z',
    created_by: 'admin',
    updated_at: '2024-01-01T00:00:00Z',
  },
];

const mockPaginatedResponse = createMockPaginatedResponse(mockGroups);
const mockEmptyResponse = createMockPaginatedResponse([]);

describe('GroupsPage', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders page title', async () => {
    (apiClient.getGroups as jest.Mock).mockResolvedValue(mockEmptyResponse);

    renderWithProviders(<GroupsPage />);

    const titles = screen.getAllByText('groups.title');
    expect(titles.length).toBeGreaterThanOrEqual(1);
  });

  it('renders new group button', async () => {
    (apiClient.getGroups as jest.Mock).mockResolvedValue(mockEmptyResponse);

    renderWithProviders(<GroupsPage />);

    expect(screen.getByText('groups.newGroup')).toBeInTheDocument();
  });

  it('shows loading skeleton while fetching', async () => {
    (apiClient.getGroups as jest.Mock).mockImplementation(() => new Promise(() => {}));

    renderWithProviders(<GroupsPage />);

    const skeletons = document.querySelectorAll('.animate-pulse');
    expect(skeletons.length).toBeGreaterThan(0);
  });

  it('displays groups in table', async () => {
    (apiClient.getGroups as jest.Mock).mockResolvedValue(mockPaginatedResponse);

    renderWithProviders(<GroupsPage />);

    await waitFor(() => {
      expect(screen.getByText('Group A')).toBeInTheDocument();
    });

    expect(screen.getByText('Group B')).toBeInTheDocument();
  });

  it('shows no results when empty', async () => {
    (apiClient.getGroups as jest.Mock).mockResolvedValue(mockEmptyResponse);

    renderWithProviders(<GroupsPage />);

    await waitFor(() => {
      expect(screen.getByText('common.noResults')).toBeInTheDocument();
    });
  });

  it('renders table headers', async () => {
    (apiClient.getGroups as jest.Mock).mockResolvedValue(mockEmptyResponse);

    renderWithProviders(<GroupsPage />);

    await waitFor(() => {
      expect(screen.getByText('common.noResults')).toBeInTheDocument();
    });

    expect(screen.getByText('common.id')).toBeInTheDocument();
    expect(screen.getByText('common.name')).toBeInTheDocument();
    expect(screen.getByText('common.status')).toBeInTheDocument();
    expect(screen.getByText('common.actions')).toBeInTheDocument();
  });
});
