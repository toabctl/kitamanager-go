import { screen } from '@testing-library/react';
import { Breadcrumb } from '../breadcrumb';
import { renderWithProviders } from '@/test-utils';

describe('Breadcrumb', () => {
  it('renders all items', () => {
    renderWithProviders(
      <Breadcrumb
        items={[
          { label: 'Home', href: '/' },
          { label: 'Users', href: '/users' },
          { label: 'John' },
        ]}
      />
    );

    expect(screen.getByText('Home')).toBeInTheDocument();
    expect(screen.getByText('Users')).toBeInTheDocument();
    expect(screen.getByText('John')).toBeInTheDocument();
  });

  it('renders links for non-last items with href', () => {
    renderWithProviders(
      <Breadcrumb items={[{ label: 'Home', href: '/' }, { label: 'Current' }]} />
    );

    const homeLink = screen.getByText('Home').closest('a');
    expect(homeLink).toHaveAttribute('href', '/');

    // Last item should not be a link
    const current = screen.getByText('Current');
    expect(current.closest('a')).toBeNull();
  });

  it('renders last item as span even with href', () => {
    renderWithProviders(
      <Breadcrumb
        items={[
          { label: 'Home', href: '/' },
          { label: 'Page', href: '/page' },
        ]}
      />
    );

    // Last item rendered as span, not link
    const page = screen.getByText('Page');
    expect(page.tagName).toBe('SPAN');
  });

  it('renders item without href as span', () => {
    renderWithProviders(<Breadcrumb items={[{ label: 'No Link' }, { label: 'Current' }]} />);

    expect(screen.getByText('No Link').tagName).toBe('SPAN');
  });

  it('has accessible nav landmark', () => {
    renderWithProviders(<Breadcrumb items={[{ label: 'Home' }]} />);
    expect(screen.getByLabelText('Breadcrumb')).toBeInTheDocument();
  });
});
