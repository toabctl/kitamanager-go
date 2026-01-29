import { describe, it, expect } from 'vitest'
import type {
  Organization,
  User,
  Employee,
  Child,
  EmployeeContract,
  ChildContract
} from '../api/types'

describe('API Types', () => {
  it('should allow valid Organization object', () => {
    const org: Organization = {
      id: 1,
      name: 'Test Org',
      active: true,
      state: 'berlin',
      created_at: '2024-01-01T00:00:00Z',
      created_by: 'admin@example.com',
      updated_at: '2024-01-01T00:00:00Z'
    }

    expect(org.id).toBe(1)
    expect(org.name).toBe('Test Org')
    expect(org.active).toBe(true)
    expect(org.state).toBe('berlin')
  })

  it('should allow Organization with optional users and groups', () => {
    const org: Organization = {
      id: 1,
      name: 'Test Org',
      active: true,
      state: 'berlin',
      created_at: '2024-01-01T00:00:00Z',
      created_by: 'admin@example.com',
      updated_at: '2024-01-01T00:00:00Z',
      users: [],
      groups: []
    }

    expect(org.users).toEqual([])
    expect(org.groups).toEqual([])
  })

  it('should allow valid User object', () => {
    const user: User = {
      id: 1,
      name: 'John Doe',
      email: 'john@example.com',
      active: true,
      is_superadmin: false,
      created_at: '2024-01-01T00:00:00Z',
      created_by: 'admin@example.com',
      updated_at: '2024-01-01T00:00:00Z'
    }

    expect(user.id).toBe(1)
    expect(user.email).toBe('john@example.com')
    expect(user.is_superadmin).toBe(false)
  })

  it('should allow valid Employee with contracts', () => {
    const contract: EmployeeContract = {
      id: 1,
      employee_id: 1,
      from: '2024-01-01',
      to: '2024-12-31',
      position: 'Erzieher',
      grade: 'S8a',
      step: 3,
      weekly_hours: 40,
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z'
    }

    const employee: Employee = {
      id: 1,
      organization_id: 1,
      first_name: 'Max',
      last_name: 'Mustermann',
      gender: 'male',
      birthdate: '1990-01-01',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
      contracts: [contract]
    }

    expect(employee.contracts).toHaveLength(1)
    expect(employee.contracts![0].position).toBe('Erzieher')
  })

  it('should allow Employee contract with null end date', () => {
    const contract: EmployeeContract = {
      id: 1,
      employee_id: 1,
      from: '2024-01-01',
      to: null,
      position: 'Erzieher',
      grade: 'S8a',
      step: 3,
      weekly_hours: 40,
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z'
    }

    expect(contract.to).toBeNull()
  })

  it('should allow valid Child with contracts', () => {
    const contract: ChildContract = {
      id: 1,
      child_id: 1,
      from: '2024-01-01',
      to: null,
      attributes: ['ganztags', 'ndh'],
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z'
    }

    const child: Child = {
      id: 1,
      organization_id: 1,
      first_name: 'Emma',
      last_name: 'Schmidt',
      gender: 'female',
      birthdate: '2020-03-15',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
      contracts: [contract]
    }

    expect(child.contracts).toHaveLength(1)
    expect(child.contracts![0].attributes).toContain('ganztags')
  })
})
