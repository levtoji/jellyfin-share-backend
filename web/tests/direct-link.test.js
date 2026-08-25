import { test, expect } from '@playwright/test';

const mockShareInfo = {
  title: 'Test Movie',
  itemType: 'Movie',
  expiresAt: '2027-01-01T00:00:00Z',
  requiresPassword: false,
  totalPlays: 0,
  currentConcurrentViewers: 0,
  audioTracks: [
    { index: 1, language: 'deu', displayTitle: 'Deutsch (E-AC3)', codec: 'eac3', channels: 6, isDefault: true },
    { index: 2, language: 'eng', displayTitle: 'English (E-AC3)', codec: 'eac3', channels: 6, isDefault: false },
  ],
  subtitleTracks: [
    { index: 3, language: 'deu', displayTitle: 'Deutsch - SUBRIP', codec: 'subrip', isDefault: true, isForced: true },
    { index: 5, language: 'deu', displayTitle: 'Full SRT - Deutsch', codec: 'subrip', isDefault: false, isForced: false },
  ],
};

test.beforeEach(async ({ page }) => {
  await page.route('**/api/public/shares/testtoken', (route) => {
    route.fulfill({ json: mockShareInfo });
  });
});

test('direct link updates when audio track is selected', async ({ page }) => {
  await page.goto('/s/testtoken');

  const audioSelect = page.locator('#audio-select');
  await audioSelect.waitFor();
  await audioSelect.selectOption('2');

  const input = page.locator('.direct-link-input');
  await expect(input).toHaveValue(/.*\/direct\/testtoken\?AudioStreamIndex=2$/);
});

test('direct link updates when subtitle track is selected', async ({ page }) => {
  await page.goto('/s/testtoken');

  const subtitleSelect = page.locator('#subtitle-select');
  await subtitleSelect.waitFor();
  await subtitleSelect.selectOption('5');

  const input = page.locator('.direct-link-input');
  await expect(input).toHaveValue(/.*\/direct\/testtoken\?SubtitleStreamIndex=5&SubtitleMethod=Encode$/);
});

test('direct link combines audio and subtitle params', async ({ page }) => {
  await page.goto('/s/testtoken');

  await page.locator('#audio-select').selectOption('2');
  await page.locator('#subtitle-select').selectOption('5');

  const input = page.locator('.direct-link-input');
  await expect(input).toHaveValue(/.*\/direct\/testtoken\?AudioStreamIndex=2&SubtitleStreamIndex=5&SubtitleMethod=Encode$/);
});

test('direct link starts with no params (default)', async ({ page }) => {
  await page.goto('/s/testtoken');

  const input = page.locator('.direct-link-input');
  await input.waitFor();
  await expect(input).toHaveValue(/.*\/direct\/testtoken$/);
});
