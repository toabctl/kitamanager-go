// We store mock functions on a shared object so that jest.mock factory
// (which is hoisted above all const declarations) can reference them.
// Using 'var' avoids the temporal dead zone that 'const'/'let' would cause.

var __mockGet: jest.Mock;
var __mockPost: jest.Mock;
var __mockPut: jest.Mock;
var __mockDel: jest.Mock;

jest.mock('axios', () => {
  __mockGet = jest.fn();
  __mockPost = jest.fn();
  __mockPut = jest.fn();
  __mockDel = jest.fn();

  const instance = {
    get: __mockGet,
    post: __mockPost,
    put: __mockPut,
    delete: __mockDel,
    interceptors: {
      request: { use: jest.fn() },
      response: { use: jest.fn() },
    },
  };
  return {
    __esModule: true,
    default: {
      create: jest.fn(() => instance),
    },
  };
});

import { getErrorMessage, apiClient } from '../client';

describe('getErrorMessage', () => {
  it('extracts message from axios error response', () => {
    const error = {
      response: {
        data: {
          message: 'Invalid credentials',
        },
      },
    };

    expect(getErrorMessage(error, 'Fallback message')).toBe('Invalid credentials');
  });

  it('returns fallback for error without response', () => {
    const error = new Error('Network error');

    expect(getErrorMessage(error, 'Fallback message')).toBe('Fallback message');
  });

  it('returns fallback for error without message in response', () => {
    const error = {
      response: {
        data: {},
      },
    };

    expect(getErrorMessage(error, 'Fallback message')).toBe('Fallback message');
  });

  it('returns fallback for null error', () => {
    expect(getErrorMessage(null, 'Fallback message')).toBe('Fallback message');
  });

  it('returns fallback for undefined error', () => {
    expect(getErrorMessage(undefined, 'Fallback message')).toBe('Fallback message');
  });

  it('returns fallback for non-object error', () => {
    expect(getErrorMessage('string error', 'Fallback message')).toBe('Fallback message');
    expect(getErrorMessage(123, 'Fallback message')).toBe('Fallback message');
  });
});

describe('export URL builders', () => {
  describe('getEmployeesExportUrl', () => {
    it('builds URL without filters', () => {
      const url = apiClient.getEmployeesExportUrl(1);
      expect(url).toBe('/api/v1/organizations/1/employees/export/excel');
    });

    it('builds URL with filters', () => {
      const url = apiClient.getEmployeesExportUrl(1, {
        search: 'John',
        staff_category: 'qualified',
        active_on: '2026-02-01',
      });
      expect(url).toContain('/api/v1/organizations/1/employees/export/excel?');
      expect(url).toContain('search=John');
      expect(url).toContain('staff_category=qualified');
      expect(url).toContain('active_on=2026-02-01');
    });

    it('omits undefined and empty filters', () => {
      const url = apiClient.getEmployeesExportUrl(1, {
        search: undefined,
        staff_category: '',
        active_on: '2026-02-01',
      });
      expect(url).toContain('active_on=2026-02-01');
      expect(url).not.toContain('search');
      expect(url).not.toContain('staff_category');
    });
  });

  describe('getChildrenExportUrl', () => {
    it('builds URL without filters', () => {
      const url = apiClient.getChildrenExportUrl(1);
      expect(url).toBe('/api/v1/organizations/1/children/export/excel');
    });

    it('builds URL with filters', () => {
      const url = apiClient.getChildrenExportUrl(1, {
        search: 'Max',
        section_id: '3',
        active_on: '2026-03-01',
      });
      expect(url).toContain('/api/v1/organizations/1/children/export/excel?');
      expect(url).toContain('search=Max');
      expect(url).toContain('section_id=3');
      expect(url).toContain('active_on=2026-03-01');
    });

    it('omits undefined and empty filters', () => {
      const url = apiClient.getChildrenExportUrl(2, {
        search: undefined,
        section_id: undefined,
      });
      expect(url).toBe('/api/v1/organizations/2/children/export/excel');
    });
  });
});
