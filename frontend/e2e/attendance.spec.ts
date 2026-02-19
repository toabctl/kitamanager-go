import { test, expect } from '@playwright/test';
import {
  login,
  createTestOrg,
  deleteTestOrg,
  createChildWithContractViaApi,
  deleteChildViaApi,
  clearAttendanceForDate,
  uniqueName,
} from './utils/test-helpers';

test.use({ locale: 'en-US' });

test.describe('Attendance', () => {
  let orgId: number;
  let childId: number;
  const childFirstName = uniqueName('AttChild');

  test.beforeAll(async ({ browser }) => {
    const page = await browser.newPage();
    await login(page);
    const testOrg = await createTestOrg(page, 'Attendance');
    orgId = testOrg.orgId;
    const child = await createChildWithContractViaApi(page, orgId, {
      first_name: childFirstName,
      last_name: 'Test',
      gender: 'female',
      birthdate: '2021-06-10',
    });
    childId = child.id;
    await page.close();
  });

  test.afterAll(async ({ browser }) => {
    const page = await browser.newPage();
    await login(page);
    await deleteChildViaApi(page, orgId, childId).catch(() => {});
    await deleteTestOrg(page, orgId);
    await page.close();
  });

  test.beforeEach(async ({ page }) => {
    await login(page);
    await page.goto(`/organizations/${orgId}/attendance`);
    await page.waitForLoadState('networkidle');
  });

  test('should display attendance page with heading and week stepper', async ({ page }) => {
    await expect(page.getByRole('heading', { name: /attendance/i }).first()).toBeVisible();
    // Week stepper should show "This week" button
    await expect(page.getByRole('button', { name: /this week/i })).toBeVisible();
  });

  test('should show child in weekly attendance table', async ({ page }) => {
    await expect(page.getByText(childFirstName)).toBeVisible({ timeout: 10000 });

    // Table should have 6 column headers (Name + 5 weekdays)
    const headers = page.locator('thead th');
    await expect(headers).toHaveCount(6);
  });

  test('should check-in child', async ({ page }) => {
    await expect(page.getByText(childFirstName)).toBeVisible({ timeout: 10000 });

    const row = page.getByRole('row').filter({ hasText: childFirstName });

    // Click the first check-in button in the row
    await row
      .getByRole('button', { name: /check-in/i })
      .first()
      .click();

    // After check-in, a check-out button should appear
    await expect(row.getByRole('button', { name: /check-out/i }).first()).toBeVisible({
      timeout: 10000,
    });
  });

  test('should check-out child after check-in', async ({ page }) => {
    await expect(page.getByText(childFirstName)).toBeVisible({ timeout: 10000 });

    const row = page.getByRole('row').filter({ hasText: childFirstName });

    // If not checked in yet, check in first
    const checkOutButton = row.getByRole('button', { name: /check-out/i }).first();
    const checkInButton = row.getByRole('button', { name: /check-in/i }).first();

    if (await checkInButton.isVisible().catch(() => false)) {
      await checkInButton.click();
      await expect(checkOutButton).toBeVisible({ timeout: 10000 });
    }

    // Click check-out
    await checkOutButton.click();

    // After check-out, the times should be displayed (no more buttons)
    // The cell should show a time range like "HH:MM – HH:MM"
    await expect(row.locator('text=/\\d{2}:\\d{2}\\s*–\\s*\\d{2}:\\d{2}/')).toBeVisible({
      timeout: 10000,
    });
  });

  test('should navigate between weeks', async ({ page }) => {
    // Navigate to previous week
    await page.getByRole('button', { name: /previous week/i }).click();
    await page.waitForLoadState('networkidle');

    // Navigate to next week
    await page.getByRole('button', { name: /next week/i }).click();
    await page.waitForLoadState('networkidle');

    // Click "This week" to return
    await page.getByRole('button', { name: /this week/i }).click();
    await page.waitForLoadState('networkidle');

    // Child should still be visible
    await expect(page.getByText(childFirstName)).toBeVisible({ timeout: 10000 });
  });
});

test.describe('Attendance Status Transitions', () => {
  let orgId: number;
  let childId: number;
  const childFirstName = uniqueName('StatusChild');

  test.beforeAll(async ({ browser }) => {
    const page = await browser.newPage();
    await login(page);
    const testOrg = await createTestOrg(page, 'AttStatus');
    orgId = testOrg.orgId;
    const child = await createChildWithContractViaApi(page, orgId, {
      first_name: childFirstName,
      last_name: 'Test',
      gender: 'male',
      birthdate: '2022-03-15',
    });
    childId = child.id;
    await page.close();
  });

  test.afterAll(async ({ browser }) => {
    const page = await browser.newPage();
    await login(page);
    await deleteChildViaApi(page, orgId, childId).catch(() => {});
    await deleteTestOrg(page, orgId);
    await page.close();
  });

  test.beforeEach(async ({ page }) => {
    await login(page);
    // Clear any attendance for today so each test starts fresh
    const today = new Date().toISOString().slice(0, 10);
    await clearAttendanceForDate(page, orgId, childId, today);
    await page.goto(`/organizations/${orgId}/attendance`);
    await page.waitForLoadState('networkidle');
    await expect(page.getByText(childFirstName)).toBeVisible({ timeout: 10000 });
  });

  test('check-in shows time and check-out button', async ({ page }) => {
    const row = page.getByRole('row').filter({ hasText: childFirstName });

    // Click check-in
    await row
      .getByRole('button', { name: /check-in/i })
      .first()
      .click();

    // Should show a check-out button and a time (HH:MM)
    await expect(row.getByRole('button', { name: /check-out/i }).first()).toBeVisible({
      timeout: 10000,
    });
    await expect(row.locator('text=/\\d{2}:\\d{2}/')).toBeVisible();
  });

  test('mark sick via popover shows Sick status text', async ({ page }) => {
    const row = page.getByRole('row').filter({ hasText: childFirstName });

    // Open the popover (⋯ button)
    await row
      .getByRole('button', { name: /quick mark/i })
      .first()
      .click();

    // Click the Sick status button in the popover
    await page.getByRole('button', { name: /^Sick$/i }).click();

    // Should show "Sick" status text
    await expect(row.getByText('Sick')).toBeVisible({ timeout: 10000 });
  });

  test('check-in then mark sick clears times and shows Sick', async ({ page }) => {
    const row = page.getByRole('row').filter({ hasText: childFirstName });

    // Check-in first
    await row
      .getByRole('button', { name: /check-in/i })
      .first()
      .click();
    await expect(row.getByRole('button', { name: /check-out/i }).first()).toBeVisible({
      timeout: 10000,
    });

    // Open popover and mark sick
    await row
      .getByRole('button', { name: /quick mark/i })
      .first()
      .click();
    await page.getByRole('button', { name: /^Sick$/i }).click();

    // Should show "Sick" text, check-out button from today's cell should be gone
    await expect(row.getByText('Sick')).toBeVisible({ timeout: 10000 });
    // Time range should not be visible (times were cleared)
    await expect(row.locator('text=/\\d{2}:\\d{2}\\s*–\\s*\\d{2}:\\d{2}/')).toBeHidden();
  });

  test('mark sick then mark present auto-sets check-in time', async ({ page }) => {
    const row = page.getByRole('row').filter({ hasText: childFirstName });

    // Mark as sick via popover
    await row
      .getByRole('button', { name: /quick mark/i })
      .first()
      .click();
    await page.getByRole('button', { name: /^Sick$/i }).click();
    await expect(row.getByText('Sick')).toBeVisible({ timeout: 10000 });

    // Now mark as present via popover
    await row
      .getByRole('button', { name: /quick mark/i })
      .first()
      .click();
    await page.getByRole('button', { name: /^Present$/i }).click();

    // Should show check-out button (meaning check-in time was auto-set)
    await expect(row.getByRole('button', { name: /check-out/i }).first()).toBeVisible({
      timeout: 10000,
    });
    // Should NOT show "Present" as italic text (should show time instead)
    await expect(row.getByText('Present')).toBeHidden();
  });

  test('check-in, check-out, then mark vacation clears both times', async ({ page }) => {
    const row = page.getByRole('row').filter({ hasText: childFirstName });

    // Check-in
    await row
      .getByRole('button', { name: /check-in/i })
      .first()
      .click();
    await expect(row.getByRole('button', { name: /check-out/i }).first()).toBeVisible({
      timeout: 10000,
    });

    // Check-out
    await row
      .getByRole('button', { name: /check-out/i })
      .first()
      .click();
    // Should show time range (HH:MM – HH:MM)
    await expect(row.locator('text=/\\d{2}:\\d{2}\\s*–\\s*\\d{2}:\\d{2}/')).toBeVisible({
      timeout: 10000,
    });

    // Mark vacation via popover
    await row
      .getByRole('button', { name: /quick mark/i })
      .first()
      .click();
    await page.getByRole('button', { name: /^Vacation$/i }).click();

    // Should show "Vacation" text, times should be gone
    await expect(row.getByText('Vacation')).toBeVisible({ timeout: 10000 });
    await expect(row.locator('text=/\\d{2}:\\d{2}\\s*–\\s*\\d{2}:\\d{2}/')).toBeHidden();
  });

  test('mark absent then mark present shows check-in time', async ({ page }) => {
    const row = page.getByRole('row').filter({ hasText: childFirstName });

    // Mark absent
    await row
      .getByRole('button', { name: /quick mark/i })
      .first()
      .click();
    await page.getByRole('button', { name: /^Absent$/i }).click();
    await expect(row.getByText('Absent')).toBeVisible({ timeout: 10000 });

    // Mark present
    await row
      .getByRole('button', { name: /quick mark/i })
      .first()
      .click();
    await page.getByRole('button', { name: /^Present$/i }).click();

    // Should show check-out button (check-in auto-set)
    await expect(row.getByRole('button', { name: /check-out/i }).first()).toBeVisible({
      timeout: 10000,
    });
  });
});
