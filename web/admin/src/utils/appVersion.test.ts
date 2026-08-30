import { describe, expect, it } from 'vitest';
import { formatFeedCraftVersion } from '@/utils/appVersion';

describe('formatFeedCraftVersion', () => {
  it('uses the CI-injected version as-is when present', () => {
    expect(formatFeedCraftVersion('v3.2.0', '1.0.0')).toBe('FeedCraft v3.2.0');
    expect(formatFeedCraftVersion('dev-abc1234', '1.0.0')).toBe(
      'FeedCraft dev-abc1234'
    );
  });

  it('falls back to package.json version for local development', () => {
    expect(formatFeedCraftVersion(undefined, '1.0.0')).toBe('FeedCraft v1.0.0');
    expect(formatFeedCraftVersion('  ', '1.0.0')).toBe('FeedCraft v1.0.0');
  });
});
