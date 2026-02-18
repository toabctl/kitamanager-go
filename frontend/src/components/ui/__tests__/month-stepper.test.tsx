import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MonthStepper } from '../month-stepper';
import { renderWithProviders } from '@/test-utils';

jest.mock('next-intl', () => ({
  useTranslations: () => (key: string) => key,
  useLocale: () => 'en',
}));

describe('MonthStepper', () => {
  const onChange = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders the current month and year', () => {
    renderWithProviders(<MonthStepper value={new Date(2026, 1, 1)} onChange={onChange} />);

    expect(screen.getByText('1. February 2026')).toBeInTheDocument();
  });

  it('renders navigation buttons', () => {
    renderWithProviders(<MonthStepper value={new Date(2026, 1, 1)} onChange={onChange} />);

    expect(screen.getByRole('button', { name: 'previousMonth' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'nextMonth' })).toBeInTheDocument();
    expect(screen.getByText('today')).toBeInTheDocument();
  });

  it('calls onChange with previous month when left arrow clicked', async () => {
    const user = userEvent.setup();
    renderWithProviders(<MonthStepper value={new Date(2026, 1, 1)} onChange={onChange} />);

    await user.click(screen.getByRole('button', { name: 'previousMonth' }));

    expect(onChange).toHaveBeenCalledTimes(1);
    const calledDate = onChange.mock.calls[0][0] as Date;
    expect(calledDate.getFullYear()).toBe(2026);
    expect(calledDate.getMonth()).toBe(0); // January
    expect(calledDate.getDate()).toBe(1);
  });

  it('calls onChange with next month when right arrow clicked', async () => {
    const user = userEvent.setup();
    renderWithProviders(<MonthStepper value={new Date(2026, 1, 1)} onChange={onChange} />);

    await user.click(screen.getByRole('button', { name: 'nextMonth' }));

    expect(onChange).toHaveBeenCalledTimes(1);
    const calledDate = onChange.mock.calls[0][0] as Date;
    expect(calledDate.getFullYear()).toBe(2026);
    expect(calledDate.getMonth()).toBe(2); // March
    expect(calledDate.getDate()).toBe(1);
  });

  it('calls onChange with today when Today button clicked', async () => {
    const user = userEvent.setup();
    renderWithProviders(<MonthStepper value={new Date(2024, 5, 1)} onChange={onChange} />);

    await user.click(screen.getByText('today'));

    expect(onChange).toHaveBeenCalledTimes(1);
    const calledDate = onChange.mock.calls[0][0] as Date;
    const now = new Date();
    expect(calledDate.getFullYear()).toBe(now.getFullYear());
    expect(calledDate.getMonth()).toBe(now.getMonth());
    expect(calledDate.getDate()).toBe(now.getDate());
  });

  it('handles year boundary correctly (January to December)', async () => {
    const user = userEvent.setup();
    renderWithProviders(<MonthStepper value={new Date(2026, 0, 1)} onChange={onChange} />);

    await user.click(screen.getByRole('button', { name: 'previousMonth' }));

    const calledDate = onChange.mock.calls[0][0] as Date;
    expect(calledDate.getFullYear()).toBe(2025);
    expect(calledDate.getMonth()).toBe(11); // December
  });

  it('handles year boundary correctly (December to January)', async () => {
    const user = userEvent.setup();
    renderWithProviders(<MonthStepper value={new Date(2025, 11, 1)} onChange={onChange} />);

    await user.click(screen.getByRole('button', { name: 'nextMonth' }));

    const calledDate = onChange.mock.calls[0][0] as Date;
    expect(calledDate.getFullYear()).toBe(2026);
    expect(calledDate.getMonth()).toBe(0); // January
  });
});
