import { test, expect } from '@playwright/test';
import {
  login,
  createTestOrg,
  deleteTestOrg,
  createChildWithContractViaApi,
  deleteChildViaApi,
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

  test('should display attendance page with heading and day stepper', async ({ page }) => {
    await expect(page.getByRole('heading', { name: /attendance/i }).first()).toBeVisible();
    // Day stepper should show today's date
    await expect(page.getByRole('button', { name: /today/i })).toBeVisible();
  });

  test('should show child in attendance list', async ({ page }) => {
    await expect(page.getByText(childFirstName)).toBeVisible({ timeout: 10000 });
    // Should show "Not Recorded" initially
    await expect(page.getByText(/not recorded/i).first()).toBeVisible();
  });

  test('should quick-mark child as present', async ({ page }) => {
    await expect(page.getByText(childFirstName)).toBeVisible({ timeout: 10000 });

    // Find the row with the child and click the present (check circle) button
    const row = page.getByRole('row').filter({ hasText: childFirstName });
    // The first quick button is "Present" (CheckCircle)
    const quickButtons = row.getByRole('button');
    // Quick buttons are after status column; find the green check one
    await quickButtons.filter({ hasText: '' }).nth(0).click();

    // Wait for the status to update - should show "Present" badge
    await expect(row.getByText(/present/i)).toBeVisible({ timeout: 10000 });
  });

  test('should quick-mark child as absent after being present', async ({ page }) => {
    await expect(page.getByText(childFirstName)).toBeVisible({ timeout: 10000 });
    const row = page.getByRole('row').filter({ hasText: childFirstName });

    // Click the absent (XCircle) button - second quick button
    const quickButtons = row.getByRole('button');
    await quickButtons.nth(1).click();

    // Should update to "Absent"
    await expect(row.getByText(/absent/i)).toBeVisible({ timeout: 10000 });
  });

  test('should edit attendance record via dialog', async ({ page }) => {
    await expect(page.getByText(childFirstName)).toBeVisible({ timeout: 10000 });
    const row = page.getByRole('row').filter({ hasText: childFirstName });

    // First make sure there's a record to edit (quick mark present)
    await row.getByRole('button').nth(0).click();
    await expect(row.getByText(/present/i)).toBeVisible({ timeout: 10000 });

    // Click the edit (pencil) button
    await row.getByRole('button', { name: /edit/i }).click();

    // Dialog should open
    await expect(page.getByRole('dialog')).toBeVisible({ timeout: 5000 });

    // Change check-in time
    await page.getByLabel(/check-in/i).fill('08:30');

    // Add a note
    const noteField = page.getByLabel(/note/i);
    await noteField.fill('Arrived with father');

    // Save
    await page.getByRole('button', { name: /save/i }).click();

    // Dialog should close
    await expect(page.getByRole('dialog')).not.toBeVisible({ timeout: 10000 });

    // Note should be visible in the row (don't assert on time due to timezone differences)
    await expect(row.getByText('Arrived with father')).toBeVisible({ timeout: 10000 });
  });

  test('should delete attendance record', async ({ page }) => {
    await expect(page.getByText(childFirstName)).toBeVisible({ timeout: 10000 });
    const row = page.getByRole('row').filter({ hasText: childFirstName });

    // First make sure there's a record (quick mark present)
    await row.getByRole('button').nth(0).click();
    await expect(row.getByText(/present/i)).toBeVisible({ timeout: 10000 });

    // Click the delete button
    await row.getByRole('button', { name: /delete/i }).click();

    // Confirm deletion in alert dialog
    await expect(page.getByRole('alertdialog')).toBeVisible({ timeout: 5000 });
    await page.getByRole('button', { name: /delete/i }).click();

    // Should go back to "Not Recorded"
    await expect(row.getByText(/not recorded/i)).toBeVisible({ timeout: 10000 });
  });

  test('should navigate between days', async ({ page }) => {
    // Click previous day button
    await page.getByRole('button', { name: /previous day/i }).click();
    await page.waitForLoadState('networkidle');

    // Click next day button twice (back to today + 1)
    await page.getByRole('button', { name: /next day/i }).click();
    await page.getByRole('button', { name: /next day/i }).click();
    await page.waitForLoadState('networkidle');

    // Click "Today" button to return
    await page.getByRole('button', { name: /today/i }).click();
    await page.waitForLoadState('networkidle');

    // Child should still be visible
    await expect(page.getByText(childFirstName)).toBeVisible({ timeout: 10000 });
  });
});
