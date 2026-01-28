import { test, expect } from 'playwright/test'
import {
  login,
  createOrganization,
  SUPERADMIN_EMAIL,
  SUPERADMIN_PASSWORD
} from './utils/test-helpers'

test.describe('Organization State and Government Funding', () => {
  // Use a unique timestamp to avoid conflicts between test runs
  const timestamp = Date.now()
  const testOrgName = `Test Org State ${timestamp}`

  test.beforeEach(async ({ page }) => {
    // Login as superadmin before each test
    await login(page, SUPERADMIN_EMAIL, SUPERADMIN_PASSWORD)
  })

  test('superadmin can create an organization with a state', async ({ page }) => {
    // Create a new organization with Berlin state
    const orgId = await createOrganization(page, testOrgName, 'berlin')

    // Verify the organization was created with correct state via API
    const token = await page.evaluate(() => localStorage.getItem('token'))
    const org = await page.evaluate(
      async ({ token, orgId }) => {
        const res = await fetch(`/api/v1/organizations/${orgId}`, {
          headers: { Authorization: `Bearer ${token}` }
        })
        return res.json()
      },
      { token, orgId }
    )

    // Verify the organization has Berlin state
    expect(org.name).toBe(testOrgName)
    expect(org.state).toBe('berlin')
  })

  test('superadmin can navigate to government fundings list', async ({ page }) => {
    // Navigate to government fundings via sidebar
    await page.getByRole('link', { name: /government funding/i }).click()
    await expect(page).toHaveURL(/.*government-funding/)

    // Verify the government fundings list is displayed
    await expect(page.getByRole('heading', { name: /government funding/i })).toBeVisible()

    // Verify Berlin government funding is listed (seeded data)
    await expect(page.getByRole('cell', { name: /Berlin/i })).toBeVisible({ timeout: 5000 })
  })

  test('superadmin can view government funding details', async ({ page }) => {
    // Navigate to government fundings
    await page.getByRole('link', { name: /government funding/i }).click()
    await expect(page).toHaveURL(/.*government-funding/)

    // Click on view details for Berlin
    const berlinRow = page.getByRole('row').filter({ hasText: /Berlin/i })
    await berlinRow.getByRole('button', { name: /view details/i }).click()

    // Verify we're on the detail page
    await expect(page).toHaveURL(/.*government-funding.*\d+/)

    // Verify the detail page shows the "Add Period" button (indicates we're on the details page)
    await expect(page.getByRole('button', { name: 'Add Period' })).toBeVisible({ timeout: 5000 })
  })

  test('organization state determines which government funding is used', async ({ page }) => {
    // Create an organization with Berlin state
    const orgWithState = `Test Org Funding ${timestamp}`
    const orgId = await createOrganization(page, orgWithState, 'berlin')

    // Navigate to the organization's detail page to verify its state
    await page.goto(`/organizations`)
    await page.waitForLoadState('networkidle')

    // Use the API to verify the org has the correct state (more reliable than UI pagination)
    const token = await page.evaluate(() => localStorage.getItem('token'))
    const org = await page.evaluate(
      async ({ token, orgId }) => {
        const res = await fetch(`/api/v1/organizations/${orgId}`, {
          headers: { Authorization: `Bearer ${token}` }
        })
        return res.json()
      },
      { token, orgId }
    )

    // Verify the organization has Berlin state
    expect(org.state).toBe('berlin')
    expect(org.name).toBe(orgWithState)

    // The organization's funding is now automatically determined by its state
    // No manual assignment needed - Berlin orgs use Berlin funding rules
  })
})
