import { test, expect } from '@playwright/test';
import {
  login,
  getFirstOrganization,
  createBudgetItemViaApi,
  deleteBudgetItemViaApi,
  uniqueName,
} from './utils/test-helpers';

test.use({ locale: 'en-US' });

test.describe('Budget Item Detail', () => {
  test.describe.configure({ mode: 'serial' });

  let orgId: number;
  let budgetItemId: number;
  let budgetItemName: string;

  test.beforeAll(async ({ browser }) => {
    const page = await browser.newPage();
    await login(page);
    const org = await getFirstOrganization(page);
    orgId = org.id;
    budgetItemName = uniqueName('DetailBudget');
    const item = await createBudgetItemViaApi(page, orgId, {
      name: budgetItemName,
      category: 'income',
    });
    budgetItemId = item.id;
    await page.close();
  });

  test.afterAll(async ({ browser }) => {
    const page = await browser.newPage();
    await login(page);
    await deleteBudgetItemViaApi(page, orgId, budgetItemId).catch(() => {});
    await page.close();
  });

  test('should display budget item detail page', async ({ page }) => {
    await login(page);
    await page.goto(`/organizations/${orgId}/budget-items/${budgetItemId}`);
    await page.waitForLoadState('networkidle');

    // Verify heading shows budget item name
    await expect(page.getByRole('heading', { name: budgetItemName })).toBeVisible({
      timeout: 10000,
    });

    // Verify category badge is visible
    await expect(page.getByText(/income/i)).toBeVisible();
  });

  test('should create an entry via UI', async ({ page }) => {
    await login(page);
    await page.goto(`/organizations/${orgId}/budget-items/${budgetItemId}`);
    await page.waitForLoadState('networkidle');

    // Click "Add Entry" button
    await page.getByRole('button', { name: /add.*entry/i }).click();
    await expect(page.getByRole('dialog')).toBeVisible({ timeout: 5000 });

    // Fill form fields
    await page.locator('#from').fill('2024-01-01');
    await page.locator('#amount_euros').fill('500.00');
    await page.locator('#notes').fill('Test entry notes');

    // Submit
    await page.getByRole('button', { name: /save/i }).click();

    // Dialog should close
    await expect(page.getByRole('dialog')).not.toBeVisible({ timeout: 10000 });

    // Verify entry appears in the table
    await expect(page.getByText('Test entry notes')).toBeVisible({ timeout: 10000 });
  });

  test('should edit an entry via UI', async ({ page }) => {
    await login(page);
    await page.goto(`/organizations/${orgId}/budget-items/${budgetItemId}`);
    await page.waitForLoadState('networkidle');

    // Wait for the entry to be visible
    await expect(page.getByText('Test entry notes')).toBeVisible({ timeout: 10000 });

    // Click the edit (Pencil) icon button - first icon button in the row's action cell
    const row = page.getByRole('row').filter({ hasText: 'Test entry notes' });
    await row.locator('button').first().click();

    // Dialog should open with pre-filled values
    await expect(page.getByRole('dialog')).toBeVisible({ timeout: 5000 });

    // Change the notes
    await page.locator('#notes').clear();
    await page.locator('#notes').fill('Updated entry notes');

    // Submit
    await page.getByRole('button', { name: /save/i }).click();

    // Dialog should close
    await expect(page.getByRole('dialog')).not.toBeVisible({ timeout: 10000 });

    // Verify updated notes appear
    await expect(page.getByText('Updated entry notes')).toBeVisible({ timeout: 10000 });
  });

  test('should delete an entry via UI', async ({ page }) => {
    await login(page);
    await page.goto(`/organizations/${orgId}/budget-items/${budgetItemId}`);
    await page.waitForLoadState('networkidle');

    // Wait for the entry to be visible
    await expect(page.getByText('Updated entry notes')).toBeVisible({ timeout: 10000 });

    // Click the delete (Trash) icon button - last icon button in the row's action cell
    const row = page.getByRole('row').filter({ hasText: 'Updated entry notes' });
    await row.locator('button').last().click();

    // Confirm deletion
    await expect(page.getByRole('alertdialog')).toBeVisible({ timeout: 5000 });
    await page.getByRole('button', { name: /delete/i }).click();

    // Entry should disappear
    await expect(page.getByText('Updated entry notes')).not.toBeVisible({ timeout: 10000 });
  });
});
