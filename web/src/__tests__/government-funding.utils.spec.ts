import { describe, it, expect } from 'vitest'
import { flattenGovernmentFundingToRows } from '../utils/government-funding'
import type { GovernmentFunding } from '../api/types'

describe('flattenGovernmentFundingToRows', () => {
  it('should return empty array for null government funding', () => {
    const result = flattenGovernmentFundingToRows(null)
    expect(result).toEqual([])
  })

  it('should return empty array for government funding without periods', () => {
    const governmentFunding: GovernmentFunding = {
      id: 1,
      name: 'Test',
      state: 'berlin',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z'
    }
    const result = flattenGovernmentFundingToRows(governmentFunding)
    expect(result).toEqual([])
  })

  it('should handle period with no properties', () => {
    const governmentFunding: GovernmentFunding = {
      id: 1,
      name: 'Test',
      state: 'berlin',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
      periods: [
        {
          id: 1,
          government_funding_id: 1,
          from: '2024-01-01',
          to: '2024-12-31',
          comment: 'Test period',
          created_at: '2024-01-01T00:00:00Z',
          properties: []
        }
      ]
    }

    const result = flattenGovernmentFundingToRows(governmentFunding)

    expect(result).toHaveLength(1)
    expect(result[0].periodId).toBe(1)
    expect(result[0].periodFrom).toBe('2024-01-01')
    expect(result[0].periodTo).toBe('2024-12-31')
    expect(result[0].propertyId).toBe(0)
    expect(result[0].propertyName).toBe('-')
    expect(result[0].ageRange).toBe('-')
    expect(result[0].isFirstPropertyInPeriod).toBe(true)
  })

  it('should handle property with age range', () => {
    const governmentFunding: GovernmentFunding = {
      id: 1,
      name: 'Test',
      state: 'berlin',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
      periods: [
        {
          id: 1,
          government_funding_id: 1,
          from: '2024-01-01',
          to: null,
          comment: '',
          created_at: '2024-01-01T00:00:00Z',
          properties: [
            {
              id: 1,
              period_id: 1,
              name: 'ganztag',
              payment: 166847,
              requirement: 0.261,
              min_age: 0,
              max_age: 2,
              comment: 'Full-day',
              created_at: '2024-01-01T00:00:00Z'
            }
          ]
        }
      ]
    }

    const result = flattenGovernmentFundingToRows(governmentFunding)

    expect(result).toHaveLength(1)
    expect(result[0].propertyId).toBe(1)
    expect(result[0].propertyName).toBe('ganztag')
    expect(result[0].minAge).toBe(0)
    expect(result[0].maxAge).toBe(2)
    expect(result[0].ageRange).toBe('0 - 2')
    expect(result[0].periodTo).toBeNull()
  })

  it('should handle property without age range', () => {
    const governmentFunding: GovernmentFunding = {
      id: 1,
      name: 'Test',
      state: 'berlin',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
      periods: [
        {
          id: 1,
          government_funding_id: 1,
          from: '2024-01-01',
          to: null,
          comment: '',
          created_at: '2024-01-01T00:00:00Z',
          properties: [
            {
              id: 1,
              period_id: 1,
              name: 'integration',
              payment: 50000,
              requirement: 0.1,
              min_age: null,
              max_age: null,
              comment: 'Integration supplement',
              created_at: '2024-01-01T00:00:00Z'
            }
          ]
        }
      ]
    }

    const result = flattenGovernmentFundingToRows(governmentFunding)

    expect(result).toHaveLength(1)
    expect(result[0].minAge).toBeNull()
    expect(result[0].maxAge).toBeNull()
    expect(result[0].ageRange).toBe('-')
  })

  it('should handle property with only min_age', () => {
    const governmentFunding: GovernmentFunding = {
      id: 1,
      name: 'Test',
      state: 'berlin',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
      periods: [
        {
          id: 1,
          government_funding_id: 1,
          from: '2024-01-01',
          to: null,
          comment: '',
          created_at: '2024-01-01T00:00:00Z',
          properties: [
            {
              id: 1,
              period_id: 1,
              name: 'school-age',
              payment: 80000,
              requirement: 0.15,
              min_age: 6,
              max_age: null,
              comment: '',
              created_at: '2024-01-01T00:00:00Z'
            }
          ]
        }
      ]
    }

    const result = flattenGovernmentFundingToRows(governmentFunding)

    expect(result).toHaveLength(1)
    expect(result[0].minAge).toBe(6)
    expect(result[0].maxAge).toBeNull()
    expect(result[0].ageRange).toBe('6+')
  })

  it('should handle property with only max_age', () => {
    const governmentFunding: GovernmentFunding = {
      id: 1,
      name: 'Test',
      state: 'berlin',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
      periods: [
        {
          id: 1,
          government_funding_id: 1,
          from: '2024-01-01',
          to: null,
          comment: '',
          created_at: '2024-01-01T00:00:00Z',
          properties: [
            {
              id: 1,
              period_id: 1,
              name: 'infant',
              payment: 200000,
              requirement: 0.3,
              min_age: null,
              max_age: 1,
              comment: '',
              created_at: '2024-01-01T00:00:00Z'
            }
          ]
        }
      ]
    }

    const result = flattenGovernmentFundingToRows(governmentFunding)

    expect(result).toHaveLength(1)
    expect(result[0].minAge).toBeNull()
    expect(result[0].maxAge).toBe(1)
    expect(result[0].ageRange).toBe('< 1')
  })

  it('should flatten full hierarchy correctly', () => {
    const governmentFunding: GovernmentFunding = {
      id: 1,
      name: 'Berlin',
      state: 'berlin',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
      periods: [
        {
          id: 1,
          government_funding_id: 1,
          from: '2024-01-01',
          to: '2024-12-31',
          comment: 'Period 2024',
          created_at: '2024-01-01T00:00:00Z',
          properties: [
            {
              id: 1,
              period_id: 1,
              name: 'ganztag-0-2',
              payment: 166847,
              requirement: 0.261,
              min_age: 0,
              max_age: 2,
              comment: 'Full-day under 2',
              created_at: '2024-01-01T00:00:00Z'
            },
            {
              id: 2,
              period_id: 1,
              name: 'ganztag-2-7',
              payment: 150000,
              requirement: 0.2,
              min_age: 2,
              max_age: 7,
              comment: 'Full-day 2-7',
              created_at: '2024-01-01T00:00:00Z'
            },
            {
              id: 3,
              period_id: 1,
              name: 'integration',
              payment: 50000,
              requirement: 0.1,
              min_age: null,
              max_age: null,
              comment: 'All ages',
              created_at: '2024-01-01T00:00:00Z'
            }
          ]
        }
      ]
    }

    const result = flattenGovernmentFundingToRows(governmentFunding)

    // Should have 3 rows (one per property)
    expect(result).toHaveLength(3)

    // First row
    expect(result[0].periodId).toBe(1)
    expect(result[0].periodFrom).toBe('2024-01-01')
    expect(result[0].periodComment).toBe('Period 2024')
    expect(result[0].propertyName).toBe('ganztag-0-2')
    expect(result[0].payment).toBe(166847)
    expect(result[0].requirement).toBe(0.261)
    expect(result[0].ageRange).toBe('0 - 2')
    expect(result[0].isFirstPropertyInPeriod).toBe(true)
    expect(result[0].periodRowSpan).toBe(3) // Total 3 rows in this period

    // Second row
    expect(result[1].propertyName).toBe('ganztag-2-7')
    expect(result[1].ageRange).toBe('2 - 7')
    expect(result[1].isFirstPropertyInPeriod).toBe(false)
    expect(result[1].periodRowSpan).toBe(0) // Not first row of period

    // Third row
    expect(result[2].propertyName).toBe('integration')
    expect(result[2].ageRange).toBe('-')
    expect(result[2].isFirstPropertyInPeriod).toBe(false)
    expect(result[2].periodRowSpan).toBe(0)
  })

  it('should handle multiple periods', () => {
    const governmentFunding: GovernmentFunding = {
      id: 1,
      name: 'Test',
      state: 'berlin',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
      periods: [
        {
          id: 1,
          government_funding_id: 1,
          from: '2024-01-01',
          to: '2024-06-30',
          comment: '',
          created_at: '2024-01-01T00:00:00Z',
          properties: [
            {
              id: 1,
              period_id: 1,
              name: 'ganztag',
              payment: 100000,
              requirement: 0.25,
              min_age: 0,
              max_age: 3,
              comment: '',
              created_at: '2024-01-01T00:00:00Z'
            }
          ]
        },
        {
          id: 2,
          government_funding_id: 1,
          from: '2024-07-01',
          to: null,
          comment: 'New rates',
          created_at: '2024-07-01T00:00:00Z',
          properties: [
            {
              id: 2,
              period_id: 2,
              name: 'ganztag',
              payment: 110000,
              requirement: 0.27,
              min_age: 0,
              max_age: 3,
              comment: 'Updated rate',
              created_at: '2024-07-01T00:00:00Z'
            }
          ]
        }
      ]
    }

    const result = flattenGovernmentFundingToRows(governmentFunding)

    expect(result).toHaveLength(2)

    // First period
    expect(result[0].periodId).toBe(1)
    expect(result[0].periodFrom).toBe('2024-01-01')
    expect(result[0].periodTo).toBe('2024-06-30')
    expect(result[0].isFirstPropertyInPeriod).toBe(true)
    expect(result[0].periodRowSpan).toBe(1)

    // Second period
    expect(result[1].periodId).toBe(2)
    expect(result[1].periodFrom).toBe('2024-07-01')
    expect(result[1].periodTo).toBeNull() // Ongoing period
    expect(result[1].periodComment).toBe('New rates')
    expect(result[1].isFirstPropertyInPeriod).toBe(true)
    expect(result[1].periodRowSpan).toBe(1)
    expect(result[1].payment).toBe(110000)
  })

  it('should correctly calculate row spans for complex structure', () => {
    // Period with 4 properties
    const governmentFunding: GovernmentFunding = {
      id: 1,
      name: 'Test',
      state: 'berlin',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
      periods: [
        {
          id: 1,
          government_funding_id: 1,
          from: '2024-01-01',
          to: null,
          comment: '',
          created_at: '2024-01-01T00:00:00Z',
          properties: [
            {
              id: 1,
              period_id: 1,
              name: 'a',
              payment: 100,
              requirement: 0.1,
              min_age: 0,
              max_age: 2,
              comment: '',
              created_at: ''
            },
            {
              id: 2,
              period_id: 1,
              name: 'b',
              payment: 200,
              requirement: 0.2,
              min_age: 0,
              max_age: 2,
              comment: '',
              created_at: ''
            },
            {
              id: 3,
              period_id: 1,
              name: 'c',
              payment: 300,
              requirement: 0.3,
              min_age: 2,
              max_age: 7,
              comment: '',
              created_at: ''
            },
            {
              id: 4,
              period_id: 1,
              name: 'd',
              payment: 400,
              requirement: 0.4,
              min_age: null,
              max_age: null,
              comment: '',
              created_at: ''
            }
          ]
        }
      ]
    }

    const result = flattenGovernmentFundingToRows(governmentFunding)

    expect(result).toHaveLength(4)

    // First row should have periodRowSpan = 4 (total rows)
    expect(result[0].isFirstPropertyInPeriod).toBe(true)
    expect(result[0].periodRowSpan).toBe(4)

    // Subsequent rows should have periodRowSpan = 0
    expect(result[1].periodRowSpan).toBe(0)
    expect(result[2].periodRowSpan).toBe(0)
    expect(result[3].periodRowSpan).toBe(0)
  })
})
