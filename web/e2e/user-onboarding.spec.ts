import { test, expect } from 'playwright/test'

/**
 * User onboarding E2E test:
 * 1. Create a new organization
 * 2. Create a group within that organization
 * 3. Create a new user
 * 4. Assign the user to the group with manager role
 * 5. Login as the new user
 */
test.describe('User Onboarding', () => {
  // Generate unique names for this test run
  const timestamp = Date.now()
  const orgName = `Test Org ${timestamp}`
  const groupName = `Test Group ${timestamp}`
  const userName = `Test Manager ${timestamp}`
  const userEmail = `manager${timestamp}@example.com`
  const userPassword = 'testpassword123'

  // This test performs many operations, increase timeout
  test.setTimeout(120000)

  test('complete user onboarding flow', async ({ page }) => {
    // =====================================
    // Step 1: Login as admin
    // =====================================
    await page.goto('/login')
    await page.getByPlaceholder('Email').fill('admin@example.com')
    await page.getByPlaceholder('Password').fill('adminadmin')
    await page.getByRole('button', { name: 'Sign In' }).click()
    await expect(page).not.toHaveURL(/.*login/, { timeout: 10000 })

    // =====================================
    // Step 2: Create Organization
    // =====================================
    await page.getByRole('link', { name: /organization/i }).first().click()
    await expect(page).toHaveURL(/.*organization/)

    // Click "New Organization" button
    await page.getByRole('button', { name: /new organization/i }).click()

    // Fill organization form
    await page.getByPlaceholder('Organization name').fill(orgName)
    await page.getByRole('button', { name: 'Save' }).click()

    // Wait for dialog to close and success toast
    await expect(page.getByRole('dialog')).not.toBeVisible({ timeout: 5000 })
    await expect(page.getByText('Organization created successfully')).toBeVisible({ timeout: 5000 })

    // =====================================
    // Step 3: Select new organization in sidebar
    // =====================================
    // Reload to refresh sidebar's org list (fetched on mount)
    await page.reload()
    await page.waitForLoadState('networkidle')

    // Wait for the sidebar's organizations list to load (dropdown should not be loading)
    // The dropdown shows loading indicator while fetching orgs
    await page.waitForTimeout(1000)

    // Click the org dropdown in sidebar
    const orgDropdown = page.getByRole('combobox').first()
    await orgDropdown.click()

    // Use the filter input to find our org (dropdown has filter enabled)
    const filterInput = page.getByPlaceholder('Search...')
    await filterInput.fill(timestamp.toString())
    await page.waitForTimeout(500)

    // Now the filtered option should be visible - select it
    const orgOption = page.getByRole('option', { name: orgName })
    await expect(orgOption).toBeVisible({ timeout: 10000 })
    await orgOption.click()

    // Wait for sidebar to update with org-scoped links
    await page.waitForTimeout(500)

    // =====================================
    // Step 4: Create Group for the organization
    // =====================================
    // Navigate to Groups via sidebar (now visible after org selection)
    await page.getByRole('link', { name: /group/i }).first().click()
    await expect(page).toHaveURL(/.*groups/)

    // Click "New Group" button
    await page.getByRole('button', { name: /new group/i }).click()

    // Fill group form
    await page.getByPlaceholder('Group name').fill(groupName)
    await page.getByRole('button', { name: 'Save' }).click()

    // Wait for dialog to close and verify group appears in table
    await expect(page.getByRole('dialog')).not.toBeVisible({ timeout: 5000 })
    await expect(page.getByRole('cell', { name: groupName })).toBeVisible()

    // =====================================
    // Step 5: Create new User
    // =====================================
    // Navigate to Users via sidebar (now org-scoped)
    await page.getByRole('link', { name: /user/i }).first().click()
    await expect(page).toHaveURL(/.*users/)

    // Click "New User" button
    await page.getByRole('button', { name: /new user/i }).click()

    // Fill user form
    await page.getByPlaceholder('Full name').fill(userName)
    await page.getByPlaceholder('Email address').fill(userEmail)
    await page.getByPlaceholder('Password').fill(userPassword)
    await page.getByRole('button', { name: 'Save' }).click()

    // Wait for dialog to close and success toast
    await expect(page.getByRole('dialog')).not.toBeVisible({ timeout: 5000 })
    await expect(page.getByText('User created successfully')).toBeVisible({ timeout: 5000 })

    // =====================================
    // Step 6: Add user to group with manager role via API
    // =====================================
    // Note: The user won't appear in the org-scoped table until added to a group.
    // Since there's no UI to add a user to a group directly, we'll use the API.
    // First, get the user ID and group ID by querying the API.
    const apiToken = await page.evaluate(async () => {
      const response = await fetch('/api/v1/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: 'admin@example.com', password: 'adminadmin' })
      })
      const data = await response.json()
      return data.token
    })

    // Get the newly created user
    const usersResponse = await page.evaluate(async (token) => {
      const response = await fetch('/api/v1/users', {
        headers: { Authorization: `Bearer ${token}` }
      })
      return response.json()
    }, apiToken)
    const newUser = usersResponse.data.find((u: { email: string }) => u.email === userEmail)

    // Get the group in the organization
    const groupsResponse = await page.evaluate(
      async ({ token, orgName }) => {
        // First get the org
        const orgsResponse = await fetch('/api/v1/organizations?limit=100', {
          headers: { Authorization: `Bearer ${token}` }
        })
        const orgsData = await orgsResponse.json()
        const org = orgsData.data.find((o: { name: string }) => o.name === orgName)

        // Then get the groups (returns array directly, not paginated)
        const groupsResponse = await fetch(`/api/v1/organizations/${org.id}/groups`, {
          headers: { Authorization: `Bearer ${token}` }
        })
        const groupsData = await groupsResponse.json()
        // Handle both array and paginated response formats
        const groups = Array.isArray(groupsData) ? groupsData : groupsData.data || []
        return { org, groups }
      },
      { token: apiToken, orgName }
    )
    const newGroup = groupsResponse.groups.find((g: { name: string }) => g.name === groupName)

    // Add user to group with manager role
    await page.evaluate(
      async ({ token, userId, groupId }) => {
        await fetch(`/api/v1/users/${userId}/groups`, {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ group_id: groupId, role: 'manager' })
        })
      },
      { token: apiToken, userId: newUser.id, groupId: newGroup.id }
    )

    // Reload the page to see the user in the org users table
    await page.reload()
    await page.waitForLoadState('networkidle')

    // Navigate to Users page
    await page.getByRole('link', { name: /user/i }).first().click()
    await expect(page).toHaveURL(/.*users/)

    // Verify user now appears in the table
    await expect(page.getByRole('cell', { name: userName })).toBeVisible({ timeout: 10000 })

    // =====================================
    // Step 7: Logout
    // =====================================
    // Click on the user email button in the header to open the menu
    await page.getByRole('button', { name: 'admin@example.com' }).click()
    // Click Logout in the popup menu
    await page.getByRole('menuitem', { name: /logout|sign out|abmelden/i }).click()
    await expect(page).toHaveURL(/.*login/, { timeout: 10000 })

    // =====================================
    // Step 8: Login as the new manager user
    // =====================================
    await page.getByPlaceholder('Email').fill(userEmail)
    await page.getByPlaceholder('Password').fill(userPassword)
    await page.getByRole('button', { name: 'Sign In' }).click()

    // Should redirect to dashboard
    await expect(page).not.toHaveURL(/.*login/, { timeout: 10000 })

    // Verify we're logged in as the new user (check if user name appears somewhere)
    // The new user should have access but may see limited content based on their org access
    await expect(page.locator('body')).toContainText(/dashboard/i)
  })
})
