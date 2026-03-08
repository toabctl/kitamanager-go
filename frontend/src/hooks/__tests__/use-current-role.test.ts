import { hasMinimumRole } from '../use-current-role';

describe('hasMinimumRole', () => {
  it('returns false for null current role', () => {
    expect(hasMinimumRole(null, 'member')).toBe(false);
  });

  it('superadmin has minimum role for everything', () => {
    expect(hasMinimumRole('superadmin', 'superadmin')).toBe(true);
    expect(hasMinimumRole('superadmin', 'admin')).toBe(true);
    expect(hasMinimumRole('superadmin', 'manager')).toBe(true);
    expect(hasMinimumRole('superadmin', 'member')).toBe(true);
  });

  it('admin meets admin and below', () => {
    expect(hasMinimumRole('admin', 'superadmin')).toBe(false);
    expect(hasMinimumRole('admin', 'admin')).toBe(true);
    expect(hasMinimumRole('admin', 'manager')).toBe(true);
    expect(hasMinimumRole('admin', 'member')).toBe(true);
  });

  it('manager meets manager and below', () => {
    expect(hasMinimumRole('manager', 'superadmin')).toBe(false);
    expect(hasMinimumRole('manager', 'admin')).toBe(false);
    expect(hasMinimumRole('manager', 'manager')).toBe(true);
    expect(hasMinimumRole('manager', 'member')).toBe(true);
  });

  it('member only meets member', () => {
    expect(hasMinimumRole('member', 'superadmin')).toBe(false);
    expect(hasMinimumRole('member', 'admin')).toBe(false);
    expect(hasMinimumRole('member', 'manager')).toBe(false);
    expect(hasMinimumRole('member', 'member')).toBe(true);
  });
});
