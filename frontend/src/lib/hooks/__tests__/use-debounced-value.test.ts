import { renderHook, act } from '@testing-library/react';
import { useDebouncedValue } from '../use-debounced-value';

describe('useDebouncedValue', () => {
  beforeEach(() => {
    jest.useFakeTimers();
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it('returns the initial value immediately', () => {
    const { result } = renderHook(() => useDebouncedValue('hello', 500));
    expect(result.current).toBe('hello');
  });

  it('updates the debounced value after the delay', () => {
    const { result, rerender } = renderHook(({ value, delay }) => useDebouncedValue(value, delay), {
      initialProps: { value: 'hello', delay: 500 },
    });

    // Change the value
    rerender({ value: 'world', delay: 500 });

    // Before the timer fires, the debounced value should still be the old one
    expect(result.current).toBe('hello');

    // Advance time past the delay
    act(() => {
      jest.advanceTimersByTime(500);
    });

    expect(result.current).toBe('world');
  });

  it('does not update the debounced value before the delay expires', () => {
    const { result, rerender } = renderHook(({ value, delay }) => useDebouncedValue(value, delay), {
      initialProps: { value: 'initial', delay: 300 },
    });

    rerender({ value: 'updated', delay: 300 });

    // Advance time to just before the delay
    act(() => {
      jest.advanceTimersByTime(299);
    });

    expect(result.current).toBe('initial');

    // Now advance past the delay
    act(() => {
      jest.advanceTimersByTime(1);
    });

    expect(result.current).toBe('updated');
  });

  it('resets the timer on rapid changes and only applies the final value', () => {
    const { result, rerender } = renderHook(({ value, delay }) => useDebouncedValue(value, delay), {
      initialProps: { value: 'a', delay: 500 },
    });

    // Rapid sequence of changes
    rerender({ value: 'b', delay: 500 });
    act(() => {
      jest.advanceTimersByTime(200);
    });

    rerender({ value: 'c', delay: 500 });
    act(() => {
      jest.advanceTimersByTime(200);
    });

    rerender({ value: 'd', delay: 500 });

    // Still showing the initial value since timer keeps resetting
    expect(result.current).toBe('a');

    // Advance past the delay from the last change
    act(() => {
      jest.advanceTimersByTime(500);
    });

    // Only the final value should be applied
    expect(result.current).toBe('d');
  });
});
