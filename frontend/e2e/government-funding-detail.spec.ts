import { test, expect } from '@playwright/test';
import {
  login,
  getGovernmentFundingsViaApi,
  getGovernmentFundingViaApi,
  deleteFundingPeriodViaApi,
  ensureFundingHasProperties,
} from './utils/test-helpers';

test.use({ locale: 'en-US' });

test.describe('Government Funding Detail', () => {
  // Tests must run serially: all operate on the shared Berlin funding's periods/properties
  test.describe.configure({ mode: 'serial' });

  let fundingId: number;
  let fundingName: string;

  test.beforeAll(async ({ browser }) => {
    const page = await browser.newPage();
    await login(page);
    const fundings = await getGovernmentFundingsViaApi(page);
    expect(fundings.length).toBeGreaterThan(0);
    fundingId = fundings[0].id;
    fundingName = fundings[0].name;
    await page.close();
  });

  test.afterAll(async ({ browser }) => {
    const page = await browser.newPage();
    await login(page);

    // Ensure funding always has properties for other tests that depend on it
    await ensureFundingHasProperties(page);
    await page.close();
  });

  test('should display funding detail page', async ({ page }) => {
    await login(page);
    await page.goto(`/government-fundings/${fundingId}`);
    await page.waitForLoadState('networkidle');

    // Verify heading shows funding name
    await expect(page.getByRole('heading', { name: fundingName })).toBeVisible({
      timeout: 10000,
    });
  });

  test('should create a period via UI', async ({ page }) => {
    await login(page);

    // Delete existing periods via API to avoid overlap with the new test period.
    // The new period must be active today so it's visible on the page.
    const details = await getGovernmentFundingViaApi(page, fundingId);
    if (details.periods) {
      for (const period of details.periods) {
        await deleteFundingPeriodViaApi(page, fundingId, period.id).catch(() => {});
      }
    }

    await page.goto(`/government-fundings/${fundingId}`);
    await page.waitForLoadState('networkidle');

    // Click "Add Period" button
    await page.getByRole('button', { name: /add.*period/i }).click();
    await expect(page.getByRole('dialog')).toBeVisible({ timeout: 5000 });

    // Fill form fields - ongoing period starting 2020-01-01 (active today)
    // Hours value must match step=0.5 from min=0.1 (valid: 0.1, 0.6, ..., 39.1)
    await page.locator('#from').fill('2020-01-01');
    await page.locator('#full_time_weekly_hours').fill('39.1');
    await page.locator('#comment').fill('E2E test period');

    // Submit
    await page.getByRole('button', { name: /save/i }).click();

    // Dialog should close
    await expect(page.getByRole('dialog')).not.toBeVisible({ timeout: 10000 });

    // Verify period appears on the page (comment is shown below the date range)
    await expect(page.getByText('E2E test period')).toBeVisible({ timeout: 10000 });
  });

  test('should create a property within a period via UI', async ({ page }) => {
    await login(page);
    await page.goto(`/government-fundings/${fundingId}`);
    await page.waitForLoadState('networkidle');

    // Wait for the test period to be visible
    await expect(page.getByText('E2E test period')).toBeVisible({ timeout: 10000 });

    // Click "Add Property" button on the test period's card
    const periodCard = page.locator('[class*="card"]').filter({
      hasText: 'E2E test period',
    });
    await periodCard.getByRole('button', { name: /add.*property/i }).click();
    await expect(page.getByRole('dialog')).toBeVisible({ timeout: 5000 });

    // Fill property form fields
    await page.locator('#key').fill('care_type');
    await page.locator('#value').fill('e2e_test');
    await page.locator('#payment_euros').fill('1234.56');
    await page.locator('#requirement').fill('0.25');

    // Submit
    await page.getByRole('button', { name: /save/i }).click();

    // Dialog should close
    await expect(page.getByRole('dialog')).not.toBeVisible({ timeout: 10000 });

    // Verify property appears
    await expect(page.getByText('e2e_test')).toBeVisible({ timeout: 10000 });
  });

  test('should delete a property via UI', async ({ page }) => {
    await login(page);
    await page.goto(`/government-fundings/${fundingId}`);
    await page.waitForLoadState('networkidle');

    // Wait for the test property to be visible
    await expect(page.getByText('e2e_test')).toBeVisible({ timeout: 10000 });

    // The property value "e2e_test" is a paragraph label above the property table.
    // Navigate to its parent div to find the delete button in the adjacent table.
    await page
      .getByText('e2e_test', { exact: true })
      .locator('..')
      .locator('table button')
      .click();

    // Confirm deletion
    await expect(page.getByRole('alertdialog')).toBeVisible({ timeout: 5000 });
    await page.getByRole('button', { name: /delete/i }).click();

    // Property should disappear
    await expect(page.getByText('e2e_test')).not.toBeVisible({ timeout: 10000 });
  });

  test('should delete a period via UI', async ({ page }) => {
    await login(page);
    await page.goto(`/government-fundings/${fundingId}`);
    await page.waitForLoadState('networkidle');

    // Wait for the test period to be visible
    await expect(page.getByText('E2E test period')).toBeVisible({ timeout: 10000 });

    // Click the delete (Trash) icon button on the period card
    // Period card header has buttons: "Add Property" and delete trash icon
    const periodCard = page.locator('[class*="card"]').filter({
      hasText: 'E2E test period',
    });
    // The trash button is the second button in the header (after "Add Property")
    const trashButton = periodCard
      .locator('.flex.gap-2')
      .locator('button')
      .filter({ has: page.locator('svg') })
      .last();
    await trashButton.click();

    // Confirm deletion
    await expect(page.getByRole('alertdialog')).toBeVisible({ timeout: 5000 });
    await page.getByRole('button', { name: /delete/i }).click();

    // Period should disappear
    await expect(page.getByText('E2E test period')).not.toBeVisible({ timeout: 10000 });
  });
});
