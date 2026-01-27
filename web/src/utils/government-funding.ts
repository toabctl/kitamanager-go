import type { GovernmentFunding } from '@/api/types'

/**
 * Represents a flattened row in the government funding table view.
 * Each row combines data from period and property levels.
 */
export interface FlattenedGovernmentFundingRow {
  periodId: number
  periodFrom: string
  periodTo: string | null
  periodComment: string
  propertyId: number
  propertyName: string
  payment: number
  requirement: number
  minAge: number | null
  maxAge: number | null
  ageRange: string
  propertyComment: string
  // For visual grouping in table
  isFirstPropertyInPeriod: boolean
  periodRowSpan: number
}

/**
 * Flattens the hierarchical government funding structure into rows for table display.
 *
 * The government funding hierarchy is: GovernmentFunding -> Periods -> Properties
 * Each property can optionally have an age range (min_age, max_age).
 * This function flattens it into one row per property, with flags indicating
 * group boundaries for visual display.
 *
 * @param governmentFunding - The government funding object with nested periods and properties
 * @returns Array of flattened rows suitable for table display
 */
export function flattenGovernmentFundingToRows(
  governmentFunding: GovernmentFunding | null
): FlattenedGovernmentFundingRow[] {
  if (!governmentFunding?.periods) return []

  const rows: FlattenedGovernmentFundingRow[] = []

  for (const period of governmentFunding.periods) {
    const properties = period.properties || []
    const periodRowCount = properties.length || 1
    let isFirstPropertyInPeriod = true

    if (properties.length === 0) {
      // Period with no properties
      rows.push({
        periodId: period.id,
        periodFrom: period.from,
        periodTo: period.to || null,
        periodComment: period.comment || '',
        propertyId: 0,
        propertyName: '-',
        payment: 0,
        requirement: 0,
        minAge: null,
        maxAge: null,
        ageRange: '-',
        propertyComment: '',
        isFirstPropertyInPeriod: true,
        periodRowSpan: 1
      })
    } else {
      for (const prop of properties) {
        const minAge = prop.min_age ?? null
        const maxAge = prop.max_age ?? null
        let ageRange = '-'
        if (minAge !== null && maxAge !== null) {
          ageRange = `${minAge} - ${maxAge}`
        } else if (minAge !== null) {
          ageRange = `${minAge}+`
        } else if (maxAge !== null) {
          ageRange = `< ${maxAge}`
        }

        rows.push({
          periodId: period.id,
          periodFrom: period.from,
          periodTo: period.to || null,
          periodComment: period.comment || '',
          propertyId: prop.id,
          propertyName: prop.name,
          payment: prop.payment,
          requirement: prop.requirement,
          minAge,
          maxAge,
          ageRange,
          propertyComment: prop.comment || '',
          isFirstPropertyInPeriod,
          periodRowSpan: isFirstPropertyInPeriod ? periodRowCount : 0
        })
        isFirstPropertyInPeriod = false
      }
    }
  }

  return rows
}
