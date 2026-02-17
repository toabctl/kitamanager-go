import { test, expect } from '@playwright/test';
import {
  login,
  getFirstOrganization,
  createPayPlanViaApi,
  deletePayPlanViaApi,
  getPayPlansViaApi,
  uniqueName,
} from './utils/test-helpers';

test.use({ locale: 'en-US' });

test.describe('Pay Plans', () => {
  let orgId: number;

  test.beforeEach(async ({ page }) => {
    await login(page);
    const org = await getFirstOrganization(page);
    orgId = org.id;
    await page.goto(`/organizations/${orgId}/payplans`);
    await page.waitForLoadState('networkidle');
  });

  test('should display pay plans list', async ({ page }) => {
    await expect(page.getByRole('heading', { name: /pay plan/i }).first()).toBeVisible();
    await expect(page.locator('table, [role="table"]')).toBeVisible({ timeout: 10000 });
  });

  test('should create a new pay plan via UI', async ({ page }) => {
    const planName = uniqueName('TestPlan');

    // Click "New Pay Plan" button
    await page.getByRole('button', { name: /new pay plan/i }).click();
    await expect(page.getByRole('dialog')).toBeVisible({ timeout: 5000 });

    // Fill form fields
    await page.getByLabel(/name/i).fill(planName);

    // Submit
    await page.getByRole('button', { name: /save/i }).click();

    // Dialog should close
    await expect(page.getByRole('dialog')).not.toBeVisible({ timeout: 10000 });

    // Pay plan should appear in list
    await expect(page.getByText(planName)).toBeVisible({ timeout: 10000 });

    // Cleanup via API
    const plans = await getPayPlansViaApi(page, orgId);
    const created = plans.find((p) => p.name === planName);
    if (created) {
      await deletePayPlanViaApi(page, orgId, created.id);
    }
  });

  test('should edit a pay plan via UI', async ({ page }) => {
    // Setup: create pay plan via API
    const origName = uniqueName('EditPlan');
    const plan = await createPayPlanViaApi(page, orgId, origName);

    // Reload to see the plan
    await page.reload();
    await page.waitForLoadState('networkidle');
    await expect(page.getByText(origName)).toBeVisible({ timeout: 10000 });

    // Click edit button on the plan's row
    const row = page.getByRole('row').filter({ hasText: origName });
    await row.getByRole('button', { name: /edit/i }).click();

    // Dialog should open
    await expect(page.getByRole('dialog')).toBeVisible({ timeout: 5000 });

    // Modify name
    const updatedName = uniqueName('Updated');
    await page.getByLabel(/name/i).clear();
    await page.getByLabel(/name/i).fill(updatedName);

    // Submit
    await page.getByRole('button', { name: /save/i }).click();

    // Dialog should close
    await expect(page.getByRole('dialog')).not.toBeVisible({ timeout: 10000 });

    // Updated name should appear
    await expect(page.getByText(updatedName)).toBeVisible({ timeout: 10000 });

    // Cleanup
    await deletePayPlanViaApi(page, orgId, plan.id);
  });

  test('should delete a pay plan via UI', async ({ page }) => {
    // Setup: create pay plan via API
    const planName = uniqueName('DelPlan');
    await createPayPlanViaApi(page, orgId, planName);

    // Reload to see the plan
    await page.reload();
    await page.waitForLoadState('networkidle');
    await expect(page.getByText(planName)).toBeVisible({ timeout: 10000 });

    // Click delete button on the plan's row
    const row = page.getByRole('row').filter({ hasText: planName });
    await row.getByRole('button', { name: /delete/i }).click();

    // Confirm deletion in alert dialog
    await expect(page.getByRole('alertdialog')).toBeVisible({ timeout: 5000 });
    await page.getByRole('button', { name: /delete/i }).click();

    // Pay plan should disappear
    await expect(page.getByText(planName)).not.toBeVisible({ timeout: 10000 });
  });
});
