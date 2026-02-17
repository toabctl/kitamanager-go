import { test, expect } from '@playwright/test';
import {
  login,
  getFirstOrganization,
  createPayPlanViaApi,
  deletePayPlanViaApi,
  uniqueName,
} from './utils/test-helpers';

test.use({ locale: 'en-US' });

test.describe('Pay Plan Detail', () => {
  test.describe.configure({ mode: 'serial' });

  let orgId: number;
  let payPlanId: number;
  let payPlanName: string;

  test.beforeAll(async ({ browser }) => {
    const page = await browser.newPage();
    await login(page);
    const org = await getFirstOrganization(page);
    orgId = org.id;
    payPlanName = uniqueName('DetailPlan');
    const plan = await createPayPlanViaApi(page, orgId, payPlanName);
    payPlanId = plan.id;
    await page.close();
  });

  test.afterAll(async ({ browser }) => {
    const page = await browser.newPage();
    await login(page);
    await deletePayPlanViaApi(page, orgId, payPlanId).catch(() => {});
    await page.close();
  });

  /** Switch to Panels view and wait for the view to actually render */
  async function switchToPanelsView(page: import('@playwright/test').Page) {
    await page.getByRole('button', { name: /panels/i }).click();
    // "Add Period" button only appears in Panels view - wait for it
    await expect(page.getByRole('button', { name: /add.*period/i })).toBeVisible({
      timeout: 5000,
    });
  }

  test('should display pay plan detail page', async ({ page }) => {
    await login(page);
    await page.goto(`/organizations/${orgId}/payplans/${payPlanId}`);
    await page.waitForLoadState('networkidle');

    // Verify heading shows pay plan name
    await expect(page.getByRole('heading', { name: payPlanName })).toBeVisible({
      timeout: 10000,
    });

    // Verify view toggle buttons exist
    await expect(page.getByRole('button', { name: /panels/i })).toBeVisible();
    await expect(page.getByRole('button', { name: /table/i })).toBeVisible();
  });

  test('should create a period via UI', async ({ page }) => {
    await login(page);
    await page.goto(`/organizations/${orgId}/payplans/${payPlanId}`);
    await page.waitForLoadState('networkidle');

    // Switch to panels view
    await switchToPanelsView(page);

    // Click "Add Period" button
    await page.getByRole('button', { name: /add.*period/i }).click();
    await expect(page.getByRole('dialog')).toBeVisible({ timeout: 5000 });

    // Fill form fields
    await page.locator('#from').fill('2024-01-01');
    await page.locator('#weekly_hours').fill('39');
    await page.locator('#employer_contribution_rate').fill('20');

    // Submit
    await page.getByRole('button', { name: /save/i }).click();

    // Dialog should close
    await expect(page.getByRole('dialog')).not.toBeVisible({ timeout: 10000 });

    // Verify period card appears - "Add Entry" button is unique to panels view period cards
    await expect(page.getByRole('button', { name: /add.*entry/i })).toBeVisible({
      timeout: 10000,
    });
  });

  test('should edit a period via UI', async ({ page }) => {
    await login(page);
    await page.goto(`/organizations/${orgId}/payplans/${payPlanId}`);
    await page.waitForLoadState('networkidle');

    // Switch to panels view and wait for period card buttons to appear
    await switchToPanelsView(page);
    await expect(page.getByRole('button', { name: /add.*entry/i })).toBeVisible({
      timeout: 10000,
    });

    // The period card's action area has: "Add Entry" button, Pencil edit, Trash delete
    // They're siblings in a flex container. Use the "Add Entry" button to find the container,
    // then click the next sibling button (edit Pencil)
    const addEntryBtn = page.getByRole('button', { name: /add.*entry/i });
    // The edit button is the immediate next sibling button after "Add Entry"
    const editBtn = addEntryBtn.locator('xpath=following-sibling::button[1]');
    await editBtn.click();

    // Dialog should open
    await expect(page.getByRole('dialog')).toBeVisible({ timeout: 5000 });

    // Change weekly hours
    await page.locator('#weekly_hours').clear();
    await page.locator('#weekly_hours').fill('38.5');

    // Submit
    await page.getByRole('button', { name: /save/i }).click();

    // Dialog should close
    await expect(page.getByRole('dialog')).not.toBeVisible({ timeout: 10000 });

    // Verify updated value
    await expect(page.getByText(/38\.5h/)).toBeVisible({ timeout: 10000 });
  });

  test('should create an entry within a period via UI', async ({ page }) => {
    await login(page);
    await page.goto(`/organizations/${orgId}/payplans/${payPlanId}`);
    await page.waitForLoadState('networkidle');

    // Switch to panels view
    await switchToPanelsView(page);
    await expect(page.getByRole('button', { name: /add.*entry/i })).toBeVisible({
      timeout: 10000,
    });

    // Click "Add Entry" button on the period card
    await page.getByRole('button', { name: /add.*entry/i }).click();
    await expect(page.getByRole('dialog')).toBeVisible({ timeout: 5000 });

    // Fill form fields
    await page.locator('#grade').fill('S8a');
    await page.locator('#step').fill('3');
    await page.locator('#monthly_amount_euros').fill('3500.00');
    await page.locator('#step_min_years').fill('5');

    // Submit
    await page.getByRole('button', { name: /save/i }).click();

    // Dialog should close
    await expect(page.getByRole('dialog')).not.toBeVisible({ timeout: 10000 });

    // Verify entry appears with grade
    await expect(page.getByText('S8a')).toBeVisible({ timeout: 10000 });
  });

  test('should delete an entry via UI', async ({ page }) => {
    await login(page);
    await page.goto(`/organizations/${orgId}/payplans/${payPlanId}`);
    await page.waitForLoadState('networkidle');

    // Switch to panels view
    await switchToPanelsView(page);

    // Wait for the entry to be visible
    await expect(page.getByText('S8a')).toBeVisible({ timeout: 10000 });

    // Find the entry row and click the delete (last) button
    const row = page.getByRole('row').filter({ hasText: 'S8a' });
    await row.locator('button').last().click();

    // Confirm deletion
    await expect(page.getByRole('alertdialog')).toBeVisible({ timeout: 5000 });
    await page.getByRole('button', { name: /delete/i }).click();

    // Entry should disappear
    await expect(page.getByText('S8a')).not.toBeVisible({ timeout: 10000 });
  });

  test('should delete a period via UI', async ({ page }) => {
    await login(page);
    await page.goto(`/organizations/${orgId}/payplans/${payPlanId}`);
    await page.waitForLoadState('networkidle');

    // Switch to panels view
    await switchToPanelsView(page);
    await expect(page.getByRole('button', { name: /add.*entry/i })).toBeVisible({
      timeout: 10000,
    });

    // Click the delete (Trash) icon button - second button after "Add Entry"
    const addEntryBtn = page.getByRole('button', { name: /add.*entry/i });
    const deleteBtn = addEntryBtn.locator('xpath=following-sibling::button[2]');
    await deleteBtn.click();

    // Confirm deletion
    await expect(page.getByRole('alertdialog')).toBeVisible({ timeout: 5000 });
    await page.getByRole('button', { name: /delete/i }).click();

    // Period should disappear - Add Entry button should no longer be visible
    await expect(page.getByRole('button', { name: /add.*entry/i })).not.toBeVisible({
      timeout: 10000,
    });
  });

  test('should toggle between panels and table view', async ({ page }) => {
    await login(page);
    await page.goto(`/organizations/${orgId}/payplans/${payPlanId}`);
    await page.waitForLoadState('networkidle');

    // Default is table view. Click "Panels" to switch
    await page.getByRole('button', { name: /panels/i }).click();
    // "Add Period" button only appears in Panels view
    await expect(page.getByRole('button', { name: /add.*period/i })).toBeVisible({
      timeout: 5000,
    });

    // Click "Table" button to switch back
    await page.getByRole('button', { name: /table/i }).click();
    // "Add Period" button should disappear in Table view
    await expect(page.getByRole('button', { name: /add.*period/i })).not.toBeVisible({
      timeout: 5000,
    });
  });
});
