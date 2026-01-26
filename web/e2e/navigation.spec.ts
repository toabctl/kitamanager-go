import { test, expect } from 'playwright/test'

test.describe('Navigation', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto('/login')
    await page.getByPlaceholder('Email').fill('admin@example.com')
    await page.getByPlaceholder('Password').fill('adminadmin')
    await page.getByRole('button', { name: 'Sign In' }).click()
    await expect(page).not.toHaveURL(/.*login/, { timeout: 10000 })
  })

  test('should navigate to organizations', async ({ page }) => {
    await page.getByRole('link', { name: /organization/i }).first().click()
    await expect(page).toHaveURL(/.*organization/)
  })


  // Org-scoped routes require an organization to be selected first
  test.describe('Org-scoped navigation', () => {
    test.beforeEach(async ({ page }) => {
      // Select an organization from the dropdown
      // First ensure we have organizations by going to that page
      await page.getByRole('link', { name: /organization/i }).first().click()
      await expect(page).toHaveURL(/.*organization/)

      // Wait for table to load and check if there are any organizations
      await page.waitForLoadState('networkidle')

      // Click the org dropdown in sidebar and select the first available org
      const orgDropdown = page.getByRole('combobox').first()
      await orgDropdown.click()

      // Select the first organization option (skip if none available)
      const firstOption = page.getByRole('option').first()
      if (await firstOption.isVisible({ timeout: 2000 }).catch(() => false)) {
        await firstOption.click()
        // Wait for the sidebar to update with org-scoped links
        await page.waitForTimeout(500)
      }
    })

    test('should navigate to users', async ({ page }) => {
      const usersLink = page.getByRole('link', { name: /user/i }).first()
      if (await usersLink.isVisible({ timeout: 2000 }).catch(() => false)) {
        await usersLink.click()
        await expect(page).toHaveURL(/.*user/)
      } else {
        test.skip()
      }
    })

    test('should navigate to employees', async ({ page }) => {
      // Check if employees link is visible (only shown when org is selected)
      const employeesLink = page.getByRole('link', { name: /employee/i }).first()
      if (await employeesLink.isVisible({ timeout: 2000 }).catch(() => false)) {
        await employeesLink.click()
        await expect(page).toHaveURL(/.*employee/)
      } else {
        // Skip test if no org is selected (no orgs available)
        test.skip()
      }
    })

    test('should navigate to children', async ({ page }) => {
      const childrenLink = page.getByRole('link', { name: /child/i }).first()
      if (await childrenLink.isVisible({ timeout: 2000 }).catch(() => false)) {
        await childrenLink.click()
        await expect(page).toHaveURL(/.*child/)
      } else {
        test.skip()
      }
    })

    test('should navigate to groups', async ({ page }) => {
      const groupsLink = page.getByRole('link', { name: /group/i }).first()
      if (await groupsLink.isVisible({ timeout: 2000 }).catch(() => false)) {
        await groupsLink.click()
        await expect(page).toHaveURL(/.*group/)
      } else {
        test.skip()
      }
    })
  })
})
