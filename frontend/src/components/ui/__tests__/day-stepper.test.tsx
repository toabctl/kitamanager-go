import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { DayStepper } from '../day-stepper';
import { renderWithProviders } from '@/test-utils';

jest.mock('next-intl', () => ({
  useTranslations: () => (key: string) => key,
  useLocale: () => 'en',
}));

describe('DayStepper', () => {
  const onChange = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders navigation buttons', () => {
    renderWithProviders(<DayStepper value={new Date(2026, 0, 15)} onChange={onChange} />);

    expect(screen.getByRole('button', { name: 'previousDay' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'nextDay' })).toBeInTheDocument();
    expect(screen.getByText('today')).toBeInTheDocument();
  });

  it('calls onChange with previous day when left arrow clicked', async () => {
    const user = userEvent.setup();
    renderWithProviders(<DayStepper value={new Date(2026, 0, 15)} onChange={onChange} />);

    await user.click(screen.getByRole('button', { name: 'previousDay' }));

    expect(onChange).toHaveBeenCalledTimes(1);
    const calledDate = onChange.mock.calls[0][0] as Date;
    expect(calledDate.getDate()).toBe(14);
    expect(calledDate.getMonth()).toBe(0);
  });

  it('calls onChange with next day when right arrow clicked', async () => {
    const user = userEvent.setup();
    renderWithProviders(<DayStepper value={new Date(2026, 0, 15)} onChange={onChange} />);

    await user.click(screen.getByRole('button', { name: 'nextDay' }));

    expect(onChange).toHaveBeenCalledTimes(1);
    const calledDate = onChange.mock.calls[0][0] as Date;
    expect(calledDate.getDate()).toBe(16);
    expect(calledDate.getMonth()).toBe(0);
  });

  it('calls onChange with today when today button clicked', async () => {
    const user = userEvent.setup();
    const now = new Date();
    renderWithProviders(<DayStepper value={new Date(2020, 0, 1)} onChange={onChange} />);

    await user.click(screen.getByText('today'));

    expect(onChange).toHaveBeenCalledTimes(1);
    const calledDate = onChange.mock.calls[0][0] as Date;
    expect(calledDate.toDateString()).toBe(now.toDateString());
  });

  it('crosses month boundary going backward', async () => {
    const user = userEvent.setup();
    renderWithProviders(<DayStepper value={new Date(2026, 1, 1)} onChange={onChange} />);

    await user.click(screen.getByRole('button', { name: 'previousDay' }));

    const calledDate = onChange.mock.calls[0][0] as Date;
    expect(calledDate.getMonth()).toBe(0); // January
    expect(calledDate.getDate()).toBe(31);
  });
});
