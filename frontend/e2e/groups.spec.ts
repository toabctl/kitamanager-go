import { test, expect } from '@playwright/test';
import {
  login,
  createTestOrg,
  deleteTestOrg,
  createGroupViaApi,
  deleteGroupViaApi,
  getGroupsViaApi,
  uniqueName,
} from './utils/test-helpers';

test.use({ locale: 'en-US' });

test.describe('Groups', () => {
  let orgId: number;

  test.beforeAll(async ({ browser }) => {
    const page = await browser.newPage();
    await login(page);
    const testOrg = await createTestOrg(page, 'Groups');
    orgId = testOrg.orgId;
    await page.close();
  });

  test.afterAll(async ({ browser }) => {
    const page = await browser.newPage();
    await login(page);
    await deleteTestOrg(page, orgId);
    await page.close();
  });

  test.beforeEach(async ({ page }) => {
    await login(page);
    await page.goto(`/organizations/${orgId}/groups`);
    await page.waitForLoadState('networkidle');
  });

  test('should display groups list', async ({ page }) => {
    await expect(page.getByRole('heading', { name: /group/i }).first()).toBeVisible();
  });

  test('should create a new group via UI', async ({ page }) => {
    const groupName = uniqueName('TestGroup');

    await page.getByRole('button', { name: /new group/i }).click();
    await expect(page.getByRole('dialog')).toBeVisible({ timeout: 5000 });

    await page.getByLabel(/name/i).fill(groupName);
    await page.getByRole('button', { name: /save/i }).click();

    await expect(page.getByRole('dialog')).not.toBeVisible({ timeout: 10000 });
    await expect(page.getByText(groupName)).toBeVisible({ timeout: 10000 });

    // Cleanup via API
    const groups = await getGroupsViaApi(page, orgId);
    const created = groups.find((g) => g.name === groupName);
    if (created) {
      await deleteGroupViaApi(page, orgId, created.id);
    }
  });

  test('should edit a group via UI', async ({ page }) => {
    const origName = uniqueName('EditGroup');
    const group = await createGroupViaApi(page, orgId, { name: origName });

    await page.reload();
    await page.waitForLoadState('networkidle');
    await expect(page.getByText(origName)).toBeVisible({ timeout: 10000 });

    const row = page.getByRole('row').filter({ hasText: origName });
    await row.getByRole('button', { name: /edit/i }).click();

    await expect(page.getByRole('dialog')).toBeVisible({ timeout: 5000 });

    const updatedName = uniqueName('Updated');
    await page.getByLabel(/name/i).clear();
    await page.getByLabel(/name/i).fill(updatedName);

    await page.getByRole('button', { name: /save/i }).click();
    await expect(page.getByRole('dialog')).not.toBeVisible({ timeout: 10000 });
    await expect(page.getByText(updatedName)).toBeVisible({ timeout: 10000 });

    await deleteGroupViaApi(page, orgId, group.id);
  });

  test('should delete a group via UI', async ({ page }) => {
    const groupName = uniqueName('DelGroup');
    await createGroupViaApi(page, orgId, { name: groupName });

    await page.reload();
    await page.waitForLoadState('networkidle');
    await expect(page.getByText(groupName)).toBeVisible({ timeout: 10000 });

    const row = page.getByRole('row').filter({ hasText: groupName });
    await row.getByRole('button', { name: /delete/i }).click();

    await expect(page.getByRole('alertdialog')).toBeVisible({ timeout: 5000 });
    await page.getByRole('button', { name: /delete/i }).click();

    await expect(page.getByText(groupName)).not.toBeVisible({ timeout: 10000 });
  });
});
