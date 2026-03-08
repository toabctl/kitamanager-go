import { cn, getCookie } from '../utils';

describe('cn', () => {
  it('merges class names', () => {
    expect(cn('foo', 'bar')).toBe('foo bar');
  });

  it('handles conditional classes', () => {
    expect(cn('foo', false && 'bar', 'baz')).toBe('foo baz');
  });

  it('resolves Tailwind conflicts (last wins)', () => {
    expect(cn('p-4', 'p-2')).toBe('p-2');
  });

  it('handles undefined and null', () => {
    expect(cn('foo', undefined, null, 'bar')).toBe('foo bar');
  });

  it('returns empty string for no args', () => {
    expect(cn()).toBe('');
  });

  it('handles arrays', () => {
    expect(cn(['foo', 'bar'])).toBe('foo bar');
  });
});

describe('getCookie', () => {
  it('returns cookie value', () => {
    document.cookie = 'token=abc123';
    expect(getCookie('token')).toBe('abc123');
  });

  it('returns null for missing cookie', () => {
    expect(getCookie('nonexistent_cookie_xyz')).toBeNull();
  });

  it('handles cookie with = in value', () => {
    document.cookie = 'data=key=value';
    expect(getCookie('data')).toBe('key=value');
  });
});
