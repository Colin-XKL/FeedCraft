import { describe, expect, it } from 'vitest';
import {
  formatFeedCraftVersion,
  resolveDisplayVersion,
} from '@/utils/appVersion';

describe('resolveDisplayVersion', () => {
  it('prefers the version handed down by the release pipeline', () => {
    expect(
      resolveDisplayVersion({
        explicitVersion: 'v3.2.0',
        branch: 'main',
        commitSha: 'becb6a35343b6f1dcac111fb105abddce70175c7',
        packageVersion: '3.1.0',
      })
    ).toBe('v3.2.0');
  });

  it('uses the short commit sha for preview builds off main', () => {
    expect(
      resolveDisplayVersion({
        branch: 'cursor/col-36-release-please-90a2',
        commitSha: 'becb6a35343b6f1dcac111fb105abddce70175c7',
        packageVersion: '3.1.0',
      })
    ).toBe('dev-becb6a3');
  });

  it('uses the package version on main and for local builds', () => {
    expect(
      resolveDisplayVersion({
        branch: 'main',
        commitSha: 'becb6a35343b6f1dcac111fb105abddce70175c7',
        packageVersion: '3.1.0',
      })
    ).toBe('v3.1.0');
    expect(resolveDisplayVersion({ packageVersion: '3.1.0' })).toBe('v3.1.0');
    expect(
      resolveDisplayVersion({
        explicitVersion: '  ',
        branch: '',
        commitSha: '',
        packageVersion: '3.1.0',
      })
    ).toBe('v3.1.0');
  });
});

describe('formatFeedCraftVersion', () => {
  it('prefixes the product name', () => {
    expect(formatFeedCraftVersion('dev-becb6a3')).toBe('FeedCraft dev-becb6a3');
  });
});
